package main

import (
	"bytes"
	"context"
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
	"syscall"
	"time"

	"github.com/msawangwan/ci.io/lib/dir"
	"github.com/msawangwan/ci.io/lib/dock"
	"github.com/msawangwan/ci.io/lib/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
	"github.com/msawangwan/ci.io/types/cred"
)

const (
	dockerAPIVersion = "1.30"
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
	credentials    cred.Github
	dirCache       dir.WorkspaceCacher
	dockerClient   *http.Client
	dockerHostAddr string
	localip        string
	accesstoken    string
)

var (
	errInvalidWebhookEvent = errors.New("not a valid webhook event, expected: push")
	errDoesNotExistInCache = errors.New("does not exist in cache")
	errIDMismatch          = errors.New("expected id doesnt match id")
)

var killsig = make(chan os.Signal, 1)

// type responseCodeMismatchError struct {
// 	Expected int
// 	Actual   int
// 	Message  string
// }

// func (rcme responseCodeMismatchError) Error() string {
// 	return fmt.Sprintf(
// 		"[response_code_err][expected: %d][actual: %d] %s",
// 		rcme.Expected,
// 		rcme.Actual,
// 		rcme.Message,
// 	)
// }

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

	dirCache = dir.NewWorkspaceCache()

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

func extractWebhookPayload(r io.Reader) (payload *github.PushEvent, e error) {
	if e = jsonutil.FromReader(r, &payload); e != nil {
		return
	}

	return
}

func createWorkspace(cache dir.WorkspaceCacher, name string) (ws string, er error) {
	ws, er = cache.MkTempDir(name)
	if er != nil {
		return
	}

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

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er = cmd.Run(); er != nil {
		log.Printf("%s", stderr.String())
		return
	}

	log.Printf("%s", stdout.String())

	return
}

func buildImage(dockfile, imgtar, tag string, client *http.Client) (imgname string, er error) {
	imgname = tag

	params := map[string]string{
		"t":          tag,
		"dockerfile": dockfile,
	}

	// TODO: deperecate
	// cmd := dock.NewBuildDockerfileCommand(params)
	// concat, er := dock.BuildAPIURLString(cmd)
	// if er != nil {
	// 	return
	// }

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

	defer f.Close()

	res, er := client.Post(uri, "application/x-tar", f)
	if er != nil {
		return
	}

	if !isExpectedResponseCode(200, res.StatusCode) {
		return "", parseDockerAPIErrorResponse(200, res)
	}

	// if req.StatusCode != 200 {
	// 	var m dock.ErrorResponse

	// 	if er = jsonutil.FromReader(req.Body, &m); er != nil {
	// 		return
	// 	}

	// 	er = responseCodeMismatchError{200, req.StatusCode, m.Message}
	// }

	return
}

func createContainer(client *http.Client, imgname, containerport string) error {
	// payload, er := dock.CreateContainerPayload{
	// 	Image: imgname,
	// 	Port:  containerport,
	// }.Build()
	// if er != nil {
	// 	return er
	// }

	// log.Printf("create container payload: %s", payload)

	// cmd, er := dock.CreateContainerAPICall{}.Build()
	// if er != nil {
	// 	return er
	// }

	// log.Printf("create container cmd: %s", cmd)

	// uri := buildAPIURL(string(cmd))
	// bufread := bytes.NewReader(payload)

	// res, er := client.Post(uri, "application/json", bufread)
	// if er != nil {
	// 	return er
	// }

	// if !isExpectedResponseCode(res.StatusCode, 201) {
	// 	return parseDockerAPIErrorResponse(201, res)
	// }

	// return nil
	req := dock.APIRequest{
		Endpoint: dock.CreateContainerAPICall{},
		Data: dock.CreateContainerPayload{
			Image: imgname,
			Port:  containerport,
		},
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 201,
	}

	return makeAPIRequest(req, client)
}

func runContainer(client *http.Client) error {
	req := dock.APIRequest{
		Endpoint:    dock.StartContainerAPICall{},
		Data:        dock.StartContainerPayload{},
		Method:      "POST",
		ContentType: "application/json",
		SuccessCode: 204,
	}

	return makeAPIRequest(req, client)
}

func makeAPIRequest(req dock.APIRequest, c *http.Client) error {
	var res *http.Response

	endpoint, er := req.Endpoint.Build()
	if er != nil {
		return er
	}

	uri := buildAPIURL(string(endpoint))

	switch {
	case req.Method == "GET":
		res, er = c.Get(uri)
	case req.Method == "POST":
		payload, er := req.Data.Build()

		if er != nil {
			return er
		}

		res, er = c.Post(
			uri,
			req.ContentType,
			bytes.NewReader(payload),
		)
	}

	if er != nil {
		return er
	}

	if !isExpectedResponseCode(res.StatusCode, req.SuccessCode) {
		return parseDockerAPIErrorResponse(req.SuccessCode, res)
	}

	return nil
}

func main() {
	var panicHandler = func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var e error

			defer func() {
				r := recover()

				if r != nil {
					switch t := r.(type) {
					case string:
						e = errors.New(t)
					case error:
						e = t
					default:
						e = errors.New("unknown error")
					}

					http.Error(w, e.Error(), http.StatusInternalServerError)

					log.Printf("[panic_handler] %s", e)
				}
			}()

			h(w, r)
		}
	}

	http.HandleFunc(endpoint, panicHandler(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling incoming webhook")

		printStats(true)

		if !isPushEvent(r) {
			panic(errInvalidWebhookEvent)
		}

		log.Printf("webhook is a valid push event, extracting payload")

		webhook, er := extractWebhookPayload(r.Body)
		if er != nil {
			panic(er)
		}

		repoName := webhook.Repository.Name

		log.Printf("payload extracted")
		log.Printf("creating workspace")

		tempws, er := createWorkspace(dirCache, repoName)
		if er != nil {
			panic(er)
		}

		cwd, er := os.Getwd()
		if er != nil {
			panic(er)
		}

		workspacePath := filepath.Join(cwd, tempws)

		log.Printf("current working dir: %s", cwd)
		log.Printf("created workspace: %s", tempws)
		log.Printf("pulling repo into: %s", workspacePath)

		if er = buildRepo(credentials, repoName, tempws); er != nil {
			panic(er)
		}

		log.Printf("repo built")
		log.Printf("building a tar file from: %s", tempws)

		archName, er := buildTar(tempws)
		if er != nil {
			panic(er)
		}

		log.Printf("uploading tar of img: %s", archName)

		imgName, er := buildImage(tempws+"/Dockerfile", archName, repoName, dockerClient)
		if er != nil {
			panic(er)
		}

		log.Printf("img name: %s", imgName)
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
