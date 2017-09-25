package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/msawangwan/ci.io/lib/ciio"
	"github.com/msawangwan/ci.io/lib/dock"
	"github.com/msawangwan/ci.io/lib/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
	"github.com/msawangwan/ci.io/types/cred"
)

const (
	version       = "1.30"
	port          = ":80"
	mime          = "application/json; charset=utf-8"
	endpoint      = "/webhooks/payload"
	mountpoint    = "/var/run/docker.sock"
	envipaddr     = "DOCK_MASTERCONTAINER_IPADDR"
	socktype      = "unix"
	scratchdir    = "__ws"
	buildfilename = "buildfile.json"
)

var (
	credential     cred.Github
	dockerClient   *http.Client
	dockerHostAddr string
	localip        string
	accesstoken    string
)

var (
	errInvalidWebhookEvent = errors.New("not a valid webhook event, expected: push")
	errJSONParseErr        = errors.New("encountered an error while read/writing json")
)

type cacher interface {
	IsCached(s string) bool
}

type cache struct {
	store map[string]string
	*sync.Mutex
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

func (c cache) IsCached(s string) bool {
	c.Lock()
	_, ok := c.store[s]
	c.Unlock()
	return ok
}

var dirCache = newCache()
var containerCache = newCache()

var pwd = func(s string) { d, _ := os.Getwd(); log.Printf("[current working dir %s] %s", d, s) }
var route = func(adr, ver, src string) string { return fmt.Sprintf("http://%s/v%s/%s", adr, ver, src) }

func init() {
	rootdir, _ := os.Getwd()
	pathenv := os.Getenv("PATH")

	os.Setenv("PATH", fmt.Sprintf("%s:%s/bin", pathenv, rootdir))

	var (
		err error
	)

	err = jsonutil.FromFilepath("secret/github.auth.json", &credential)
	if err != nil {
		log.Printf("%s", err)
	} else {
		log.Printf("loaded credentials: %+v", credential)
	}

	err = os.Mkdir(scratchdir, 655)
	if err != nil {
		log.Printf("%s", err)
	}

	err = os.Chdir(scratchdir)
	if err != nil {
		log.Printf("%s", err)
	} else {
		pwd("created scratch dir")
	}

	dockerClient = &http.Client{
		Timeout: time.Second * 10,
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

	log.Printf("server container ip: %s\n", localip)
	log.Printf("docker host container ip: %s\n", dockerHostAddr)
}

func printStats(logger io.Writer, debug bool) {
	if debug {
		logger.Write(
			[]byte(
				fmt.Sprintf("%d", runtime.NumGoroutine()),
			),
		)
	}
}

func printJSON(l io.Writer, r io.Reader) {
	formatted, e := jsonutil.ExtractBufferFormatted(r, "", "  ")
	if e != nil {
		panic(e)
	}

	io.Copy(l, &formatted)
}

func isPushEvent(logger io.Writer, r *http.Request) bool {
	eventname := r.Header.Get("x-github-event")

	logger.Write(
		[]byte(
			fmt.Sprintf("incoming webhook: %s\n", r.URL.Path),
		),
	)
	logger.Write(
		[]byte(
			fmt.Sprintf("payload event name: %s\n", eventname),
		),
	)

	if eventname != "push" {
		return false
	}

	return true
}

func extractWebhookPayload(r io.Reader) (payload *github.PushEvent, e error) {
	if e = jsonutil.FromReader(r, payload); e != nil {
		return
	}

	return
}

func getWorkspace(c *cache, dirname string) (ws string, e error) {
	c.Lock()

	if found, ok := c.store[dirname]; ok {
		ws = found
	} else {
		ws, e = ioutil.TempDir("./", dirname)
		if e != nil {
			return
		}

		c.store[dirname] = ws
	}

	c.Unlock()

	return
}

func pullRepository(c cred.Github, dir, name string) error {
	var (
		cmdout bytes.Buffer
		cmderr bytes.Buffer
	)

	args := []string{c.User, c.OAuth, dir, name}

	clone := exec.Command("clrep", args...)
	clone.Dir = dir
	clone.Stdout = &cmdout
	clone.Stderr = &cmderr

	if err := clone.Run(); err != nil {
		log.Printf("err: %s", cmderr.String())

		return err
	}

	log.Printf("succ: %s", cmdout.String())

	return nil
}

func loadBuildfile(dirpath, filename string) (b ciio.Buildfile, e error) {
	p := filepath.Join(dirpath, filename)

	if e = jsonutil.FromFilepath(p, &b); e != nil {
		return
	}

	return
}

func findPreviousContainer(c *cache, contname string) (id string, e error) {
	c.Lock()

	id = ""

	if found, ok := c.store[contname]; ok {
		id = found
	}

	c.Unlock()

	return
}

func verifyPreviousContainer(id string, c *http.Client) error {
	cmd := dock.NewContainerCommandByID("GET", "containers", "inspect", id)
	u, e := dock.BuildAPIURLString(cmd)
	if e != nil {
		return e
	}

	r, e := c.Get(u)
	if e != nil {
		return e
	}

	var (
		p dock.InspectResponse
	)

	if e = jsonutil.FromReader(r.Body, &p); e != nil {
		return e
	}

	r.Body.Close()

	if r.StatusCode != 200 {
		return fmt.Errorf("expected 200 ok but got something else when inspecting a container")
	}

	if p.ID != id {
		return fmt.Errorf("expected id %s but got id %s", p.ID, id)
	}

	return nil
}

func removePreviousContainer(id string, c *http.Client) error {
	var (
		cmd dock.ContainerCommandByID
		r   *http.Response
		u   string
		e   error
	)

	cmd = dock.NewContainerCommandByID("POST", "containers", "stop", id)
	u, e = dock.BuildAPIURLString(cmd)
	if e != nil {
		return e
	}

	r, e = c.Post(route(dockerHostAddr, version, u), mime, nil)
	if e != nil {
		return e
	}

	printJSON(os.Stdout, r.Body)
	r.Body.Close()

	cmd = dock.NewContainerCommandByID("DELETE", "containers", "", id)
	u, e = dock.BuildAPIURLString(cmd)
	if e != nil {
		return e
	}

	rq, e := http.NewRequest("DELETE", u, nil)
	if e != nil {
		return e
	}

	r, e = c.Do(rq)
	if e != nil {
		return e
	}

	printJSON(os.Stdout, r.Body)
	r.Body.Close()

	return nil
}

func createNewContainer(b ciio.Buildfile, c *http.Client) (p dock.CreateResponse, e error) {
	var postdata dock.CreateRequest

	postdata.Image = b.Image
	postdata.WorkingDir = b.WorkingDir
	postdata.Cmd = []string{b.Cmd.Exec}

	for _, v := range b.Cmd.Args {
		postdata.Cmd = append(postdata.Cmd, v)
	}

	payload, e := jsonutil.ToReader(postdata)
	if e != nil {
		return
	}

	cmd := dock.NewContainerCommand("POST", "containers", "create") // TODO: need to pass in query params?
	cmd.URLComponents.Parameters = map[string]string{
		"name": b.ContainerName,
	}

	url, e := dock.BuildAPIURLString(cmd)
	if e != nil {
		return
	}

	r, e := c.Post(
		route(dockerHostAddr, version, url),
		mime,
		payload,
	)
	if e != nil {
		return
	}

	if e = jsonutil.FromReader(r.Body, &p); e != nil {
		return
	}

	return
}

func startNewContainer(id string, c *http.Client) error {
	cmd := dock.NewContainerCommandByID("POST", "containers", "start", id)

	url, e := dock.BuildAPIURLString(cmd)
	if e != nil {
		return e
	}

	r, e := c.Post(
		route(dockerHostAddr, version, url),
		mime,
		nil,
	)
	if e != nil {
		return e
	}

	if r.StatusCode != 204 {
		return errors.New("expected 204 but got something else when starting a container")
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
				}
			}()

			h(w, r)
		}
	}

	http.HandleFunc(endpoint, panicHandler(func(w http.ResponseWriter, r *http.Request) {
		printStats(os.Stdout, true)

		if !isPushEvent(os.Stdout, r) {
			panic(errInvalidWebhookEvent)
		}

		webhook, e := extractWebhookPayload(r.Body)
		if e != nil {
			panic(e)
		}

		repoName := webhook.Repository.Name

		ws, e := getWorkspace(dirCache, repoName)
		if e != nil {
			panic(e)
		}

		if e := pullRepository(credential, ws, repoName); e != nil {
			panic(e)
		}

		buildfile, e := loadBuildfile(ws, strings.ToLower(buildfilename))
		if e != nil {
			panic(e)
		}

		containerName := buildfile.ContainerName

		cid, e := findPreviousContainer(containerCache, containerName)
		if e != nil {
			panic(e)
		}

		if cid != "" {
			if e = verifyPreviousContainer(cid, dockerClient); e != nil {
				panic(e)
			}

			if e = removePreviousContainer(cid, dockerClient); e != nil {
				panic(e)
			}
		}

		container, e := createNewContainer(buildfile, dockerClient)
		if e != nil {
			panic(e)
		}

		if e = startNewContainer(container.ID, dockerClient); e != nil {
			panic(e)
		}
	}))

	log.Fatal(http.ListenAndServe(port, nil))
}
