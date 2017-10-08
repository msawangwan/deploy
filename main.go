// +build !testing

package main

/*
todo:
- cleanup old images
-- run prune at the end
- ping service to list stats
-- print out the cache contents
-- num goroutines
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/msawangwan/ci.io/lib/midware"

	"github.com/msawangwan/ci.io/lib/cache"
	"github.com/msawangwan/ci.io/lib/dir"
	"github.com/msawangwan/ci.io/lib/dock"
	"github.com/msawangwan/ci.io/lib/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
	"github.com/msawangwan/ci.io/types/cred"
)

const (
	dockerAPIVersion = "1.32"
	port             = ":80"
	mime             = "application/json; charset=utf-8"
	endpoint         = "/webhooks/payload"
	mountpoint       = "/var/run/docker.sock"
	envipaddr        = "DOCK_MASTERCONTAINER_IPADDR"
	socktype         = "unix"
	wsdir            = "__ws"
	buildfilename    = "buildfile.json"
	dockerfilepath   = "Dockerfile"
)

var (
	credentials cred.Github

	imgCache       cache.KVStorer
	containerCache cache.KVStorer
	wsCache        cache.KVStorer

	dockerClient *http.Client

	dockerHostAddr string
	localip        string
	accesstoken    string
)

var (
	killsig = make(chan os.Signal, 1)
)

var (
	errInvalidWebhookEvent = errors.New("not a valid webhook event, expected: push")
	errItemNotInCache      = errors.New("item not in cache")
)

func init() {
	signal.Notify(killsig, syscall.SIGINT, syscall.SIGTERM)

	rootdir, _ := os.Getwd()
	pathenv := os.Getenv("PATH")

	os.Setenv("PATH", fmt.Sprintf("%s:%s/bin", pathenv, rootdir))

	var (
		err error
	)

	err = jsonutil.FromFilepath("secret/github.auth.json", &credentials)
	if err != nil {
		log.Printf("%s", err)
	} else {
		log.Printf("loaded credentials: %+v", credentials)
	}

	err = os.Mkdir(wsdir, 655)
	if err != nil {
		log.Printf("%s", err)
	}

	err = os.Chdir(wsdir)
	if err != nil {
		log.Printf("%s", err)
	} else {
		d, _ := os.Getwd()
		log.Printf("working dir: %s", d)
	}

	defer func() {
		_ = os.Chdir("../")
		_ = os.Remove(wsdir)
	}()

	dockerClient = &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(socktype, mountpoint)
			},
		},
	}

	dockerHostAddr = os.Getenv(envipaddr)

	localip, err = netutil.LocalIP("eth0")
	if err != nil {
		log.Printf("%s", err)
	}

	imgCache = dock.NewIDCache()
	containerCache = dock.NewIDCache()
	wsCache = dir.NewWorkspaceCache()

	log.Printf("server container ip: %s\n", localip)
	log.Printf("docker host container ip: %s\n", dockerHostAddr)
}

func printStats(debug bool) {
	if debug {
		log.Printf("%d", runtime.NumGoroutine())
	}
}

func isPushEvent(r *http.Request) bool {
	eventname := r.Header.Get("x-github-event")

	if eventname != "push" {
		return false
	}

	return true
}

func extractExposedPort(dockerfile string) (s string, e error) {
	out, e := exec.Command("extractexpose", dockerfile).Output()
	if e != nil {
		return
	}

	s = string(out)

	return strings.TrimSpace(s), nil
}

func getWorkspace(store cache.KVStorer, key string) (ws string, er error) {
	ws, er = store.Fetch(key)
	if er == nil {
		os.Remove(ws)
	}

	er = nil

	ws, er = dir.MkTempWorkspace(key)
	if er != nil {
		return
	}

	store.Store(key, ws)

	return
}

func parseDockerAPIErrorResponse(code int, r *http.Response) error {
	p := dock.APIResponseError{
		ExpectedCode: code,
		ActualCode:   r.StatusCode,
	}

	if e := jsonutil.FromReader(r.Body, &p); e != nil {
		return e
	}

	return p
}

func isExpectedResponseCode(expected, got int) bool {
	return expected == got
}

func buildAPIURL(call string) string {
	return fmt.Sprintf("http://%s/v%s%s", dockerHostAddr, dockerAPIVersion, call)
}

func buildRepo(c cred.Github, repoName, workspace string) error {
	var stdout, stderr bytes.Buffer

	args := []string{
		c.User,
		c.OAuth,
		repoName,
	}

	cmd := exec.Command("clrep", args...)

	cmd.Dir = workspace
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er := cmd.Run(); er != nil {
		log.Printf("%s", stderr.String())
		return er
	}

	log.Printf("%s", stdout.String())

	return nil
}

func buildTar(target string) (arch string, er error) {
	arch = target + ".tar"

	var stdout, stderr bytes.Buffer

	cmd := exec.Command("buildtar", arch, target)

	// TODO: print these outputs to a log!
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er = cmd.Run(); er != nil {
		log.Print(er)
		return
	}

	return
}

func makeAPIRequest(req dock.APIRequest, c *http.Client) (res *http.Response, er error) {
	endpoint, er := req.Endpoint.Build()
	if er != nil {
		return
	}

	uri := buildAPIURL(string(endpoint))

	log.Printf("making API request: %s", endpoint)

	switch {
	case req.Method == "GET":
		res, er = c.Get(uri)
	case req.Method == "POST":
		var payload = io.Reader(nil)

		if req.Data != nil {
			d, er := req.Data.Build()
			if er != nil {
				return nil, er
			}

			payload = bytes.NewReader(d)
		}

		res, er = c.Post(uri, req.ContentType, payload)
	case req.Method == "DELETE":
		r, er := http.NewRequest("DELETE", uri, nil)
		if er != nil {
			return nil, er
		}

		res, er = c.Do(r)
	}

	if er != nil {
		return
	}

	if !isExpectedResponseCode(res.StatusCode, req.SuccessCode) {
		return nil, parseDockerAPIErrorResponse(req.SuccessCode, res)
	}

	return
}

func fetchCachedImage(store cache.KVStorer, imgname string) (imgID string, er error) {
	imgID, er = store.Fetch(imgname)
	if er != nil {
		return imgID, cache.NewItemNotInCacheError(
			fmt.Sprintf("no image found with name: %s", imgname),
		)
	}

	return
}

func buildImage(client *http.Client, builddir, dockfile, imgtar, imgname string) (imgid string, er error) {
	wd, _ := os.Getwd()

	params := map[string]string{
		"t":          imgname,
		"dockerfile": dockfile,
	}

	cmd, er := dock.BuildImageAPICall{Parameters: params}.Build()
	if er != nil {
		return
	}

	uri := buildAPIURL(string(cmd))

	log.Printf("build api uri: %s", uri)
	log.Printf("tarfile archive: %s", imgtar)

	f, er := os.Open(imgtar)
	if er != nil {
		return
	}

	defer func() {
		f.Close()
		os.Remove(imgtar)
		os.Chdir("../")
	}()

	er = os.Chdir(builddir)
	if er != nil {
		return
	}

	// defer os.Chdir("../")

	wd, _ = os.Getwd()

	log.Printf("build dir: %s", wd)

	res, er := client.Post(uri, "application/x-tar", f)
	if er != nil {
		return
	}

	log.Printf("image built")

	if !isExpectedResponseCode(200, res.StatusCode) {
		return "", parseDockerAPIErrorResponse(200, res)
	}

	imgid, er = getImageID(client, imgname)
	if er != nil {
		return
	}

	return
}

func getImageID(client *http.Client, imgName string) (imgID string, er error) {
	req := dock.APIRequest{
		Endpoint:    dock.InspectImageAPICall{Name: imgName},
		Data:        nil,
		Method:      "GET",
		ContentType: "application/json",
		SuccessCode: 200,
	}

	res, er := makeAPIRequest(req, client)
	if er != nil {
		return
	}

	var payload dock.ImageInspectResponse

	if er = json.NewDecoder(res.Body).Decode(&payload); er != nil {
		return
	}

	imgID = payload.ID

	return
}

func cacheImage(client *http.Client, store cache.KVStorer, imgname, imgid string) (result string, er error) {
	result = fmt.Sprintf("cached image [%s][%s]", imgname, imgid)

	id, er := store.Fetch(imgname)
	if er == nil {
		buf, er := removeImage(client, id)
		if er != nil {
			return "", er
		}

		result = string(buf)
	}

	er = nil

	store.Store(imgname, imgid)

	return
}

func removeImage(client *http.Client, imgName string) (buf []byte, er error) {
	req := dock.APIRequest{
		Endpoint: dock.RemoveImageAPICall{
			Name:       imgName,
			Parameters: map[string]string{"force": "true"},
		},
		Data:        nil,
		Method:      "DELETE",
		ContentType: "application/json",
		SuccessCode: 200,
	}

	res, er := makeAPIRequest(req, client)
	if er != nil {
		return
	}

	var payload []dock.ImageDeleteResponseItem

	if er = json.NewDecoder(res.Body).Decode(&payload); er != nil {
		return
	}

	buf, er = json.MarshalIndent(&payload, "", " \t")
	if er != nil {
		return
	}

	return
}

func deleteUnusedImages(client *http.Client) (buf []byte, er error) {
	req := dock.APIRequest{
		Endpoint: dock.DeleteUnusedImagesAPICall{
			Filters: map[string]string{
				"dangling": "true",
			},
		},
		Data:        nil,
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 200,
	}

	res, er := makeAPIRequest(req, client)
	if er != nil {
		return
	}

	var p dock.ImageDeleteUnusedResponse

	if er = json.NewDecoder(res.Body).Decode(&p); er != nil {
		return
	}

	buf, er = json.MarshalIndent(p, "", "\t")
	if er != nil {
		return
	}

	return
}

func createContainer(client *http.Client, fromImg, containerPort, hostIP, hostPort string) (id string, er error) {
	req := dock.APIRequest{
		Endpoint: dock.CreateContainerAPICall{},
		Data: dock.NewCreateContainerPayload(
			fromImg,
			fmt.Sprintf("%s/tcp", containerPort),
			hostIP,
			hostPort,
		),
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 201,
	}

	res, er := makeAPIRequest(req, client)
	if er != nil {
		return
	}

	var resPayload dock.APIResponse

	if er = jsonutil.FromReader(res.Body, &resPayload); er != nil {
		return
	}

	id = resPayload.ID

	return
}

func cacheContainer(client *http.Client, store cache.KVStorer, imgName, containerID string) error {
	// prevID, er := store.Fetch(imgName)
	// if er == nil {
	// 	if er := stopContainer(client, prevID); er != nil {
	// 		return er
	// 	}

	// 	if er := removeContainer(client, prevID); er != nil {
	// 		return er
	// 	}
	// }

	// er = nil

	store.Store(imgName, containerID)

	return nil
}

func runContainer(client *http.Client, containerID string) error {
	req := dock.APIRequest{
		Endpoint: dock.StartContainerAPICall{
			ContainerID: containerID,
		},
		Data:        nil,
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 204,
	}

	_, er := makeAPIRequest(req, client)
	if er != nil {
		return er
	}

	return nil
}

func stopContainer(client *http.Client, containerID string) error {
	req := dock.APIRequest{
		Endpoint: dock.StopContainerAPICall{
			ContainerID: containerID,
		},
		Data:        nil,
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 204,
	}

	_, er := makeAPIRequest(req, client)
	if er != nil {
		return er
	}

	return nil
}

func killContainer(client *http.Client, containerID string) error {
	if er := stopContainer(client, containerID); er != nil {
		return er
	}

	if er := removeContainer(client, containerID); er != nil {
		return er
	}

	return nil
}

func removeContainer(client *http.Client, containerID string) error {
	req := dock.APIRequest{
		Endpoint: dock.RemoveContainerAPICall{
			ContainerID: containerID,
		},
		Data:        nil,
		Method:      "DELETE",
		ContentType: "application/json",
		SuccessCode: 204,
	}

	_, er := makeAPIRequest(req, client)
	if er != nil {
		return er
	}

	return nil
}

func main() {
	http.HandleFunc(endpoint, midware.Catch(func(w http.ResponseWriter, r *http.Request) {
		var (
			repoName    string
			imgName     string
			imgID       string
			containerID string
		)
		log.Printf("handling incoming webhook")

		printStats(true)

		if !isPushEvent(r) {
			panic(errInvalidWebhookEvent)
		}

		log.Printf("webhook is a valid push event, extracting repository name")

		repoName, er := github.ExtractRepositoryName(r.Body)
		if er != nil {
			panic(er)
		}

		imgName = repoName

		log.Printf("creating workspace")

		wsName, er := getWorkspace(wsCache, repoName)
		if er != nil {
			panic(er)
		}

		log.Printf("workspace dir [key: %s][value: %s]", repoName, wsName)

		cwd, er := os.Getwd()
		if er != nil {
			panic(er)
		}

		log.Printf("current working dir: %s", cwd)

		wsPath := filepath.Join(cwd, wsName)

		log.Printf("pulling repo into: %s", wsPath)

		if er = buildRepo(credentials, repoName, wsName); er != nil {
			panic(er)
		}

		log.Printf("repo built")

		dockerfile := filepath.Join(wsPath, "Dockerfile")

		log.Printf("dockerfile: %s", dockerfile)

		exposedPort, er := extractExposedPort(dockerfile)
		if er != nil {
			panic(er)
		}

		log.Printf("extracted expose port from dockerfile: %s", exposedPort)

		readySig := make(chan struct{})

		go func(ready chan struct{}) {
			log.Printf("find previous image with name: %s", imgName)

			imgID, er = fetchCachedImage(imgCache, imgName)
			if er != nil {
				if _, ok := er.(*cache.ItemNotInCacheError); ok {
					log.Printf("%s", er)
				}
			} else {
				containerID, er = containerCache.Fetch(imgName)
				if er == nil {
					killContainer(dockerClient, containerID)
				}

				buf, er := removeImage(dockerClient, imgID)
				if er != nil {
					panic(er)
				}

				log.Printf("found previous image: %s", imgID)
				log.Printf("removed previous image and layers: %s", string(buf))
			}

			log.Printf("building a tar file from: %s", wsName)

			archName, er := buildTar(wsName)
			if er != nil {
				panic(er)
			}

			log.Printf("successfully tar'ed archive: %s", archName)
			log.Printf("uploading tar of img: %s", archName)

			imgID, er = buildImage(dockerClient, wsName, wsName+"/Dockerfile", archName, repoName)
			if er != nil {
				panic(er)
			}

			// imgID, er = getImageID(dockerClient, imgName)
			// if er != nil {
			// 	panic(er)
			// }

			log.Printf("latest img name: %s", imgName)
			log.Printf("latest img ID: %s", imgID)

			close(ready)
		}(readySig)

		go func(ready chan struct{}) {
			<-ready

			containerID, er := createContainer(dockerClient, imgName, exposedPort, "", "9090")
			if er != nil {
				panic(er)
			}

			log.Printf("container created: %s", containerID)
			log.Printf("caching container: %s", containerID)

			// if er = cacheContainer(dockerClient, containerCache, imgName, containerID); er != nil {
			// 	panic(er)
			// }
			containerCache.Store(imgName, containerID)

			log.Printf("container cached: %s", containerID)
			log.Printf("start container: %s", containerID)

			if er = runContainer(dockerClient, containerID); er != nil {
				panic(er)
			}

			log.Printf("container running: %s", containerID)
			log.Printf("cache image: %s", imgID)

			imgCache.Store(imgName, imgID)
			// result, er := cacheImage(dockerClient, imgCache, imgName, imgID)
			// if er != nil {
			// 	panic(er)
			// }

			// log.Printf("%s", result)
			log.Printf("remove and delete unused images")

			deleted, er := deleteUnusedImages(dockerClient)
			if er != nil {
				panic(er)
			}

			log.Printf("cleanup completed: %s", deleted)
		}(readySig)

		log.Printf("webhook event, handled")
	}))

	go func() {
		<-killsig

		// killContainer := func(id string, c *http.Client) error {
		// var e error

		// if e = stopContainer(id, c); e != nil {
		// 	return e
		// }

		// if e = removeContainer(id, c); e != nil {
		// 	return e
		// }

		// 	return nil
		// }

		log.Printf("kill all running containers")

		// containerCache.Lock()
		// {
		// 	for k, v := range containerCache.store {
		// 		log.Printf("killing container: %s", k)

		// 		if e := killContainer(v, dockerClient); e != nil {
		// 			panic(e)
		// 		}
		// 	}
		// }
		// containerCache.Unlock()

		log.Printf("cleanup complete")

		o, e := exec.Command("cleanup").Output()
		if e != nil {
			log.Printf("%s", e)
			os.Exit(1)
		}

		log.Printf("cleanp cmd: %s", string(o))

		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(port, nil))
}
