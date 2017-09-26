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
	errInvalidWebhookEvent  = errors.New("not a valid webhook event, expected: push")
	errDoesNotExistInCache  = errors.New("does not exist in cache")
	errResponseCodeMismatch = errors.New("expected response code success but got fail")
	errIDMismatch           = errors.New("expected id doesnt match id")
)

type responseCodeMismatchError struct {
	Expected int
	Actual   int
	Message  string
}

func (rcme responseCodeMismatchError) Error() string {
	return fmt.Sprintf(
		"[response_code_err][expected: %d][actual: %d] %s",
		rcme.Expected,
		rcme.Actual,
		rcme.Message,
	)
}

type cache struct {
	store map[string]string
	sync.Mutex
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

var (
	dirCache       *cache
	containerCache *cache
)

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

	dirCache = newCache()
	containerCache = newCache()

	log.Printf("server container ip: %s\n", localip)
	log.Printf("docker host container ip: %s\n", dockerHostAddr)
}

func printJSON(l io.Writer, r io.Reader) {
	formatted, e := jsonutil.ExtractBufferFormatted(r, "", "  ")
	if e != nil {
		panic(e)
	}

	io.Copy(l, &formatted)
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

func getWorkspace(c *cache, dirname string) (ws string, e error) {
	ws = ""

	c.Lock()
	defer c.Unlock()
	{
		if dir, found := c.store[dirname]; found {
			ws = dir
		} else {
			e = errDoesNotExistInCache
		}
	}

	return
}

func createTmpWorkspace(c *cache, dirname string) (ws string, e error) {
	ws, e = ioutil.TempDir("./", dirname)
	if e != nil {
		return
	}

	c.Lock()
	defer c.Unlock()
	{
		c.store[dirname] = ws
	}

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
		return responseCodeMismatchError{200, r.StatusCode, p.Message}
	}

	if p.ID != id {
		return errIDMismatch
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
	log.Printf("-- create payload")

	// var postdata dock.CreateRequest

	// postdata.Image = b.Image
	// postdata.WorkingDir = b.WorkingDir
	// postdata.Cmd = []string{b.Cmd.Exec}

	// for _, v := range b.Cmd.Args {
	// 	postdata.Cmd = append(postdata.Cmd, v)
	// }

	// payload, e := jsonutil.ToReader(postdata)
	// if e != nil {
	// 	return
	// }

	var bb bytes.Buffer
	bb.Write([]byte(`{
			"Image": "alpine",
			"Cmd": ["echo", "HELLO, WORLD"]
		}`))

	// j := []byte(
	// 	`{
	// 		"Image": "alpine",
	// 		"Cmd": ["echo", "HELLO, WORLD"]
	// 	}`,
	// )

	// payload := bytes.NewReader(j)

	// o, e := jsonutil.ExtractBufferFormatted(payload, "", "  ")
	// if e != nil {
	// 	return
	// }

	// io.Copy(os.Stdout, &o)

	cmd := dock.NewContainerCommand("POST", "containers", "create")
	cmd.URLComponents.Parameters = map[string]string{
		"name": b.ContainerName,
	}

	url, e := dock.BuildAPIURLString(cmd)
	if e != nil {
		return
	}

	u := route(dockerHostAddr, version, url)
	log.Printf("-- post payload: %s", u)

	r, e := c.Post(u, mime, &bb)
	if e != nil {
		return
	}

	defer r.Body.Close()

	if e = jsonutil.FromReader(r.Body, &p); e != nil {
		return
	}

	jsonutil.PrettyPrintStruct(&p)

	if r.StatusCode != 201 {
		e = responseCodeMismatchError{201, r.StatusCode, p.Message}
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

	defer r.Body.Close()

	if r.StatusCode != 204 {
		p := dock.StartResponse{}

		if e = jsonutil.FromReader(r.Body, &p); e != nil {
			return e
		}
		return responseCodeMismatchError{204, r.StatusCode, p.Message}
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

		webhook, e := extractWebhookPayload(r.Body)
		if e != nil {
			panic(e)
		}

		repoName := webhook.Repository.Name

		log.Printf("payload extracted, repo name: [%s]", repoName)
		log.Printf("fetching workspace")

		ws, e := getWorkspace(dirCache, repoName)
		if e != nil {
			if e != errDoesNotExistInCache {
				panic(e)
			} else {
				log.Printf("workspace not found, creating")
				ws, e = createTmpWorkspace(dirCache, repoName)
				if e != nil {
					panic(e)
				}
			}
		}

		log.Printf("workspace: %s", ws)
		log.Printf("pulling repository latest into workspace")

		if e := pullRepository(credential, ws, repoName); e != nil {
			panic(e)
		}

		log.Printf("pull successful")
		log.Printf("locating and loading buildfile")

		buildfile, e := loadBuildfile(ws, strings.ToLower(buildfilename))
		if e != nil {
			panic(e)
		}

		log.Printf("buildfile loaded")

		containerName := buildfile.ContainerName

		log.Printf("container name: %s", containerName)
		log.Printf("looking for previous containers")

		cid, e := findPreviousContainer(containerCache, containerName)
		if e != nil {
			panic(e)
		}

		if cid != "" {
			log.Printf("found previous container: %s", cid)

			if e = verifyPreviousContainer(cid, dockerClient); e != nil {
				panic(e)
			}

			if e = removePreviousContainer(cid, dockerClient); e != nil {
				panic(e)
			}
		}

		log.Printf("creating new container")

		container, e := createNewContainer(buildfile, dockerClient)
		if e != nil {
			panic(e)
		}

		log.Printf("created a new container: %s", container.ID)

		if e = jsonutil.PrettyPrintStruct(container); e != nil {
			panic(e)
		}

		log.Printf("starting new container: %s", container.ID)

		if e = startNewContainer(container.ID, dockerClient); e != nil {
			panic(e)
		}

		log.Printf("container running: %s", container.ID)
		log.Printf("webhook event, handled")
	}))

	log.Fatal(http.ListenAndServe(port, nil))
}
