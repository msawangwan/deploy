package main

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
	"sync"
	"syscall"
	"time"

	"github.com/msawangwan/ci.io/lib/midware"

	"github.com/msawangwan/ci.io/lib/cache"
	"github.com/msawangwan/ci.io/lib/cred"
	"github.com/msawangwan/ci.io/lib/dir"
	"github.com/msawangwan/ci.io/lib/dock"
	"github.com/msawangwan/ci.io/lib/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
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
	errorlog         = "errors.log"
	outputlog        = "out.log"
	statuslog        = "status.log"
)

var (
	credentials cred.Github

	imgCache       cache.KVStorer
	containerCache cache.KVStorer
	wsCache        cache.KVStorer

	dockerClient *http.Client

	outlog  *log.Logger
	errlog  *log.Logger
	statlog *log.Logger

	onetimeCleanup sync.Once

	dockerHostAddr string
	localip        string
	accesstoken    string
)

var (
	routes     = make(map[string]string)
	killsig    = make(chan os.Signal, 1)
	cleanupsig = make(chan struct{}, 1)
)

var (
	errInvalidWebhookEvent = errors.New("not a valid webhook event, expected: push")
	errItemNotInCache      = errors.New("item not in cache")
)

func init() {
	signal.Notify(killsig, syscall.SIGINT, syscall.SIGTERM)

	var (
		err error
	)

	statlog, err = createLogger(statuslog, "[status]", os.Stdout, log.Lshortfile)
	if err != nil {
		log.Fatal(err)
	}

	outlog, err = createLogger(outputlog, "[debug]", nil, log.Lshortfile)
	if err != nil {
		log.Fatal(err)
	}

	errlog, err = createLogger(errorlog, "[err]", nil, log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		log.Fatal(err)
	}

	errlog.SetPrefix("[err][init]")
	defer errlog.SetPrefix("[err]")

	rootdir, _ := os.Getwd()
	pathenv := os.Getenv("PATH")

	os.Setenv("PATH", fmt.Sprintf("%s:%s/bin", pathenv, rootdir))

	err = jsonutil.FromFilepath("secret/github.auth.json", &credentials)
	if err != nil {
		errlog.Fatalf("%s", err)
	}

	outlog.Printf("loaded credentials: %+v", credentials)

	err = os.Mkdir(wsdir, 655)
	if err != nil {
		errlog.Printf("%s", err)
	}

	err = os.Chdir(wsdir)
	if err != nil {
		errlog.Fatalf("%s", err)
	}

	d, _ := os.Getwd()
	outlog.Printf("working dir: %s", d)

	defer func() {
		onetimeCleanup.Do(cleanup)
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
		errlog.Fatalf("%s", err)
	}

	imgCache = dock.NewIDCache()
	containerCache = dock.NewIDCache()
	wsCache = dir.NewWorkspaceCache()

	outlog.Printf("server container ip: %s\n", localip)
	outlog.Printf("docker host container ip: %s\n", dockerHostAddr)

	routes["webhook"] = "/webhooks/payload"
	routes["service"] = "/service/stats"
	routes["images"] = "/docker/images"
	routes["containers"] = "/docker/containers"
}

func cleanup() {
	apply := func(cid string) error {
		return killContainer(dockerClient, cid)
	}

	containerCache.Map(apply)
	imgCache.Map(apply)

	containerCache.Flush()
	imgCache.Flush()

	_ = os.Chdir("../")
	_ = os.Remove(wsdir)

	wsCache.Flush()

	close(cleanupsig)
}

type logger struct {
	*log.Logger
	out    io.Writer
	file   string
	prefix string
	flags  int
}

func createLogger(logfile, logprefix string, logout io.Writer, logflags int) (*log.Logger, error) {
	var (
		f *os.File
		w io.Writer
	)

	if len(logfile) > 0 {
		o, e := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if e != nil {
			return nil, e
		}

		f = o
		w = f
	}

	if logout != nil {
		if f != nil {
			w = io.MultiWriter(f, logout)
		} else {
			w = logout
		}
	}

	return log.New(w, logprefix, logflags), nil
}

// func printStats(debug bool) {
// 	if debug {
// 		statlog.Printf("%d", runtime.NumGoroutine())
// 	}
// }

func serviceStats() string {
	return fmt.Sprintf("number of goroutines: %d", runtime.NumGoroutine())
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
		errlog.Print(stderr.String())
		return er
	}

	outlog.Print(stdout.String())

	return nil
}

func buildTar(target string) (arch string, er error) {
	arch = target + ".tar"

	var stdout, stderr bytes.Buffer

	cmd := exec.Command("buildtar", arch, target)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er = cmd.Run(); er != nil {
		errlog.Print(stderr.String())
		return
	}

	outlog.Print(stdout.String())

	return
}

func makeAPIRequest(req dock.APIRequest, c *http.Client) (res *http.Response, er error) {
	endpoint, er := req.Endpoint.Build()
	if er != nil {
		return
	}

	uri := buildAPIURL(string(endpoint))

	outlog.Printf("making API request: %s", endpoint)

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

	outlog.Printf("build api uri: %s", uri)
	outlog.Printf("tarfile archive: %s", imgtar)

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

	wd, _ = os.Getwd()

	outlog.Printf("build dir: %s", wd)

	res, er := client.Post(uri, "application/x-tar", f)
	if er != nil {
		return
	}

	time.Sleep(5 * time.Second)

	if !isExpectedResponseCode(200, res.StatusCode) {
		return "", parseDockerAPIErrorResponse(200, res)
	}

	outlog.Printf("image built: %s", imgname)

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

func killContainer(client *http.Client, containerID string) error {
	if er := stopContainer(client, containerID); er != nil {
		return er
	}

	if er := removeContainer(client, containerID); er != nil {
		return er
	}

	return nil
}

func main() {
	http.HandleFunc(routes["webhook"], midware.Catch(func(w http.ResponseWriter, r *http.Request) {
		var (
			repoName    string
			imgName     string
			imgID       string
			containerID string
		)

		statlog.Printf("handling incoming webhook")

		// printStats(true)

		if !isPushEvent(r) {
			panic(errInvalidWebhookEvent)
		}

		outlog.Printf("webhook is a valid push event, extracting repository name")

		repoName, er := github.ExtractRepositoryName(r.Body)
		if er != nil {
			panic(er)
		}

		imgName = repoName

		outlog.Printf("creating workspace")

		wsName, er := getWorkspace(wsCache, repoName)
		if er != nil {
			panic(er)
		}

		outlog.Printf("workspace dir [key: %s][value: %s]", repoName, wsName)

		cwd, er := os.Getwd()
		if er != nil {
			panic(er)
		}

		outlog.Printf("current working dir: %s", cwd)

		wsPath := filepath.Join(cwd, wsName)

		outlog.Printf("pulling repo into: %s", wsPath)

		if er = buildRepo(credentials, repoName, wsName); er != nil {
			panic(er)
		}

		outlog.Printf("repo built")

		dockerfile := filepath.Join(wsPath, "Dockerfile")

		outlog.Printf("dockerfile: %s", dockerfile)

		exposedPort, er := extractExposedPort(dockerfile)
		if er != nil {
			panic(er)
		}

		outlog.Printf("extracted expose port from dockerfile: %s", exposedPort)

		readySig := make(chan struct{})

		go func(ready chan struct{}) {
			outlog.Printf("find previous image with name: %s", imgName)

			imgID, er = fetchCachedImage(imgCache, imgName)
			if er != nil {
				if _, ok := er.(*cache.ItemNotInCacheError); ok {
					outlog.Printf("%s", er)
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

				outlog.Printf("found previous image: %s", imgID)
				outlog.Printf("removed previous image and layers: %s", string(buf))
			}

			outlog.Printf("building a tar file from: %s", wsName)

			archName, er := buildTar(wsName)
			if er != nil {
				panic(er)
			}

			outlog.Printf("successfully tar'ed archive: %s", archName)
			outlog.Printf("uploading tar of img: %s", archName)

			imgID, er = buildImage(dockerClient, wsName, wsName+"/Dockerfile", archName, repoName)
			if er != nil {
				panic(er)
			}

			outlog.Printf("latest img name: %s", imgName)
			outlog.Printf("latest img ID: %s", imgID)

			close(ready)
		}(readySig)

		go func(ready chan struct{}) {
			<-ready

			containerID, er := createContainer(dockerClient, imgName, exposedPort, "", "9090")
			if er != nil {
				panic(er)
			}

			statlog.Printf("container created: %s", containerID)
			outlog.Printf("caching container: %s", containerID)

			containerCache.Store(imgName, containerID)

			statlog.Printf("container cached: %s", containerID)
			outlog.Printf("start container: %s", containerID)

			if er = runContainer(dockerClient, containerID); er != nil {
				panic(er)
			}

			statlog.Printf("container running: %s", containerID)
			outlog.Printf("cache image: %s", imgID)

			imgCache.Store(imgName, imgID)

			outlog.Printf("remove and delete unused images")

			deleted, er := deleteUnusedImages(dockerClient)
			if er != nil {
				panic(er)
			}

			outlog.Printf("unused images deleted: %s", deleted)
		}(readySig)

		statlog.Printf("webhook event, handled")
	}))

	http.HandleFunc(routes["service"], midware.Catch(func(w http.ResponseWriter, r *http.Request) {
		statlog.Print("incoming service status request")

		stats := serviceStats()

		statlog.Print(stats)

		w.WriteHeader(200)
		w.Write([]byte(stats))

		statlog.Printf("service status request handled")
	}))

	go func() {
		<-killsig

		log.Printf("run cleanup, remove images and containers")

		go onetimeCleanup.Do(cleanup)

		<-cleanupsig

		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(port, nil))
}
