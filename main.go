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

	"github.com/msawangwan/ci.io/api/ciio"
	"github.com/msawangwan/ci.io/api/github"
	"github.com/msawangwan/ci.io/lib/dock"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
	"github.com/msawangwan/ci.io/types/cred"
)

type cache struct {
	sync.Mutex
	m map[string]string
}

const (
	version       = "1.30"
	port          = ":80"
	mime          = "application/json; charset=utf-8"
	endpoint      = "/webhooks/payload"
	mountpoint    = "/var/run/docker.sock"
	envipaddr     = "CIIO_ROOT_IPADDR"
	socktype      = "unix"
	scratchdir    = "__ws"
	buildfilename = "buildfile.json"
)

var (
	credential     *cred.Github
	dockerClient   *http.Client
	dockerHostAddr string
	localip        string
	accesstoken    string
)

var (
	errInvalidWebhookEvent = errors.New("not a valid webhook event, expected: push")
	errJSONParseErr        = errors.New("encountered an error while read/writing json")
)

//var cache = struct {
//	sync.Mutex
//	m map[string]string
//}{m: make(map[string]string)}

var containercache = struct {
	sync.Mutex
	m map[string]string
}{m: make(map[string]string)}

var dircache = struct {
	sync.Mutex
	m map[string]string
}{m: make(map[string]string)}

var commands = struct {
	cloneRemoteRepo string
}{"clrep"}

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

func main() {
	var step = func(ms ...string) {
		for _, s := range ms {
			log.Printf("%s", s)
		}
	}

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

	var executeDockCmd = func(c dock.APIStringBuilder) (r *http.Response, e error) {
		log.Printf("executing: %+v", c)

		u, e := dock.BuildAPIURLString(c)
		if e != nil {
			return nil, e
		}

		log.Printf("cmd url: %s", u)

		var (
			method string
		)

		switch t := c.(type) {
		case dock.ContainerCommand:
			method = t.URLComponents.Method
		case dock.ContainerCommandByID:
			method = t.URLComponents.Method
		}

		a := route(dockerHostAddr, version, u)

		switch method {
		case "GET":
			r, e = dockerClient.Get(a)
		case "POST":
			r, e = dockerClient.Post(a, mime, bytes.NewBuffer(c.Build()))
		case "PUT":
		case "PATCH":
		case "DELETE":
			req, e := http.NewRequest("DELETE", a, nil)
			if req != nil {
				return nil, e
			}

			r, e = dockerClient.Do(req)
		}

		return
	}

	http.HandleFunc(endpoint, panicHandler(func(w http.ResponseWriter, r *http.Request) {
		step("stats", fmt.Sprintf("goroutine count: %d\n", runtime.NumGoroutine()))

		var (
			err error
		)

		step("parse headers")

		eventname := r.Header.Get("x-github-event")

		log.Printf("incoming webhook: %s\n", r.URL.Path)
		log.Printf("payload event name: %s\n", eventname)

		if eventname != "push" {
			panic(errInvalidWebhookEvent)
		}

		step("parse webhook json payload")

		var (
			res     *http.Response
			payload *github.PushEvent
		)

		if err = jsonutil.FromReader(r.Body, payload); err != nil {
			panic(err)
		}

		step("extract project name from webhook payload")

		var (
			projname string
		)

		projname = payload.Repository.Name

		log.Printf("project name: %s", projname)

		step("fetch workspace from cache or create if none")

		var (
			tmpdir string
		)

		dircache.Lock()

		if cachedDir, ok := dircache.m[projname]; ok {
			log.Printf("already exists in cache %s", projname)
			tmpdir = cachedDir
		} else {
			log.Printf("creating a cache entry for: %s", projname)

			tmpdir, err = ioutil.TempDir("./", projname)
			if err != nil {
				log.Printf("%s", err)
			}

			dircache.m[projname] = tmpdir
		}

		dircache.Unlock()

		step("fetch remote repository into tmp workspace", fmt.Sprintf("workspace: %s", tmpdir))
		pwd("working dir")

		var (
			cmdout bytes.Buffer
			cmderr bytes.Buffer
		)

		clone := exec.Command(
			commands.cloneRemoteRepo,
			credential.User,
			credential.OAuth,
			tmpdir,
			projname,
		)

		clone.Dir = tmpdir
		clone.Stdout = &cmdout
		clone.Stderr = &cmderr

		if err = clone.Run(); err != nil {
			log.Printf("err: %s", cmderr.String())
		} else {
			log.Printf("succ: %s", cmdout.String())
		}

		step("find the project buildfile")

		var (
			buildfilepath    string
			buildfilepayload ciio.Buildfile
		)

		buildfilepath = filepath.Join(tmpdir, strings.ToLower(buildfilename))

		if err = jsonutil.FromFilepath(buildfilepath, buildfilepayload); err != nil {
			panic(err)
		}

		log.Printf("project build params: %+v", buildfilepayload)

		var (
			containername = buildfilepayload.ContainerName
		)

		step("find previous container instances and remove and replace")

		if cachedID, ok := containercache.m[containername]; ok {
			var (
				res     *http.Response
				inspect dock.ContainerCommandByID
				stop    dock.ContainerCommandByID
				remove  dock.ContainerCommandByID
			)

			printJSON := func(r io.Reader) {
				formatted, e := jsonutil.ExtractBufferFormatted(r, "", "  ")
				if e != nil {
					panic(e)
				}

				io.Copy(os.Stdout, &formatted)
			}

			inspect = dock.NewContainerCommandByID("GET", "containers", "inspect", cachedID)
			stop = dock.NewContainerCommandByID("POST", "containers", "stop", cachedID)
			remove = dock.NewContainerCommandByID("DELETE", "containers", "", cachedID)

			var (
				inspectPayload dock.InspectResponse
			)

			res, err = executeDockCmd(inspect)
			if err != nil {
				panic(err)
			}

			if err = jsonutil.FromReader(res.Body, &inspectPayload); err != nil {
				panic(err)
			}

			res.Body.Close()
			log.Printf("%+v", inspectPayload)

			if res.StatusCode != 200 {
				panic(fmt.Errorf("expected 200 ok but got something else when inspecting a container"))
			}

			if inspectPayload.ID != cachedID {
				panic(fmt.Errorf("expected id %s but got id %s", inspectPayload.ID, cachedID))
			}

			res, err = executeDockCmd(stop)
			if err != nil {
				panic(err)
			}

			printJSON(res.Body)
			res.Body.Close()

			res, err = executeDockCmd(remove)
			if err != nil {
				panic(err)
			}

			printJSON(res.Body)
			res.Body.Close()
		}

		var (
			create dock.ContainerCommand
			start  dock.ContainerCommandByID
			suc201 dock.CreateResponse
		)

		create = dock.NewContainerCommand("POST", "containers", "create") // TODO: need to pass in query params?

		res, err = executeDockCmd(create)
		if err != nil {
			panic(err)
		}

		if err = jsonutil.FromReader(res.Body, &suc201); err != nil {
			panic(err)
		}

		id := suc201.ID

		log.Printf("created new container with id: %s", id)

		start = dock.NewContainerCommandByID("POST", "containers", "start", "")

		res, err = executeDockCmd(start)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != 204 {
			panic(errors.New("expected 204 but got something else when starting a container"))
		}

		/* create the container url */

		// var (
		// 	tmpl    *template.Template
		// 	tmplbuf bytes.Buffer
		// 	tmplres string
		// 	tmplurl string
		// )

		// tmpldata := struct {
		// 	Endpoint     string
		// 	QueryStrings map[string]string
		// }{
		// 	"containers/create",
		// 	map[string]string{
		// 		"name": containername,
		// 	},
		// }

		// tmplurl = `{{ .Endpoint }}?{{ range $k, $v := .QueryStrings }}{{ $k }}={{ $v }}{{ end }}`
		// tmpl = template.New("docker_url")

		// tmpl, err = tmpl.Parse(tmplurl)
		// if err != nil {
		// 	log.Printf("%s", err)
		// }

		// if err = tmpl.Execute(&tmplbuf, tmpldata); err != nil {
		// 	log.Printf("%s", err)
		// }

		// tmplres = tmplbuf.String()

		// log.Printf("command: %s", tmplres)

		// /* query the docker host and create container from parameters */

		// jsonbuf := []byte(
		// 	`{
		// 		"Image":"golang:1.9.0-alpine3.6",
		// 		"WorkingDir": "/app",
		// 		"Cmd": ["date"]
		// 	 }`,
		// )

		// res, err = dockerClient.Post(
		// 	route(dockerHostAddr, version, tmplres),
		// 	mime,
		// 	bytes.NewBuffer(jsonbuf),
		// )

		// if err != nil {
		// 	if netutil.IsTimeOutError(err) {
		// 		log.Println("timeout error")
		// 	}
		// 	panic(err)
		// }

		// buf, err := jsonutil.ExtractBufferFormatted(res.Body, "", "  ")
		// if err != nil {
		// 	panic(err)
		// }

		// res.Body.Close()
		// io.Copy(os.Stdout, &buf)
	}))

	log.Fatal(http.ListenAndServe(port, nil))
}
