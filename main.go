package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
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

	err = jsonutil.FromFile("secret/github.auth.json", &credential)
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
	var (
		gocount int
	)

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
		var (
			err error
		)

		/* some stats */

		gocount = runtime.NumGoroutine()

		log.Printf("goroutine count: %d\n", gocount)

		/* parse request header */

		eventname := r.Header.Get("x-github-event")

		log.Printf("incoming webhook: %s\n", r.URL.Path)
		log.Printf("payload event name: %s\n", eventname)

		if eventname != "push" {
			panic(errInvalidWebhookEvent)
		}

		/* parse webhook json payload */

		var (
			res     *http.Response
			payload *github.PushEvent
		)

		if err = jsonutil.FromReader(r.Body, payload); err != nil {
			panic(err)
		}

		//err error
		//body    []byte

		//body, err = ioutil.ReadAll(r.Body)

		//if err != nil {
		//log.Printf("%s", err)
		//}

		//err = json.Unmarshal(body, &payload)

		//if err != nil {
		//	log.Printf("%s", err)
		//}

		/* parse project name from payload */

		var (
			//repo     string
			projname string
		)

		//repo = payload.Repository.HTMLURL
		projname = payload.Repository.Name

		log.Printf("project name: %s", projname)
		//log.Printf("repo: %s", repo)

		/* create and cache (or read from cache) the tmp workspace */

		var (
			tmpdir string
			//tmpdirpath   string
			//tmpdirprefix string
		)

		//tmpdirpath = "./"
		//tmpdirprefix = projname

		dircache.Lock()

		if cachedDir, ok := dircache.m[projname]; ok {
			log.Printf("already exists in cache %s", projname)
			tmpdir = cachedDir
		} else {
			log.Printf("creating a cache entry for: %s", projname)

			//tmpdir, err = ioutil.TempDir(tmpdirpath, tmpdirprefix)
			tmpdir, err = ioutil.TempDir("./", projname)
			if err != nil {
				log.Printf("%s", err)
			}

			dircache.m[projname] = tmpdir
		}

		dircache.Unlock()

		log.Printf("workspace: %s", tmpdir)

		/* clone/pull the remote repo into temp workspace */

		pwd("pulling remote repo")

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

		/* find the project buildfile */

		var (
			buildfile        *os.File
			buildfilepath    string
			buildfilepayload ciio.Buildfile
		)

		buildfilepath = filepath.Join(tmpdir, strings.ToLower(buildfilename))

		if err = jsonutil.FromFile(buildfilepath, buildfilepayload); err != nil {
			panic(err)
		}

		//buildfile, err = os.Open(buildfilepath)
		//if err != nil {
		//	log.Printf("%s", err)
		//}

		//parsed := json.NewDecoder(buildfile)

		//var buildfilepayload ciio.Buildfile

		//if err = parsed.Decode(&buildfilepayload); err != nil {
		//	log.Printf("%s", err)
		//}

		log.Printf("project build params: %+v", buildfilepayload)

		var (
			containername = buildfilepayload.ContainerName
		)

		/* find any previous images and replace them! */

		if cachedID, ok := containercache.m[containername]; ok {
			//inspect := dock.ContainerCommandByID{
			//	URLComponents: dock.URLComponents{
			//		Command: "containers",
			//		Option:  "inspect",
			//	},
			//	ID: cachedID,
			//}
			inspect := dock.NewContainerCommandByID("containers", "inspect", cachedID)

			url, err := dock.BuildAPIURLString(inspect)
			if err != nil {
				panic(err)
			}

			log.Printf("executing: %+v", inspect)
			log.Printf("command url: %s", url)

			res, err = dockerClient.Get(route(dockerHostAddr, version, url))
			if err != nil {
				if netutil.IsTimeOutError(err) {
					log.Println("timeout error") // todo: fix
				}
				panic(err)
			}

			// TODO: finish from here
			stop := dock.NewContainerCommandByID("containers", "stop", cachedID)
			remove := dock.NewContainerCommandByID("containers", "remove", cachedID)

			//stop := dock.ContainerCommandByID{
			//	URLComponents: dock.URLComponents{
			//		Command: "containers",
			//		Option:  "stop",
			//	},
			//	ID: cachedID,
			//}

			//remove := dock.ContainerCommandByID{
			//	URLComponents: dock.URLComponents{
			//		Command: "containers",
			//		Option:  "remove",
			//	},
			//	ID: cachedID,
			//}
		} else {
			log.Printf("no previous container found")
		}

		/* create the container url TODO: replace with dock.ContainerCommand struct*/

		var (
			tmpl    *template.Template
			tmplbuf bytes.Buffer
			tmplres string
			tmplurl string
		)

		tmpldata := struct {
			Endpoint     string
			QueryStrings map[string]string
		}{
			"containers/create",
			map[string]string{
				"name": containername,
			},
		}

		tmplurl = `{{ .Endpoint }}?{{ range $k, $v := .QueryStrings }}{{ $k }}={{ $v }}{{ end }}`
		tmpl = template.New("docker_url")

		tmpl, err = tmpl.Parse(tmplurl)
		if err != nil {
			log.Printf("%s", err)
		}

		if err = tmpl.Execute(&tmplbuf, tmpldata); err != nil {
			log.Printf("%s", err)
		}

		tmplres = tmplbuf.String()

		log.Printf("command: %s", tmplres)

		/* query the docker host and create container from parameters */

		jsonbuf := []byte(
			`{
				"Image":"golang:1.9.0-alpine3.6",
				"WorkingDir": "/app",
				"Cmd": ["date"]
			 }`,
		)

		res, err = dockerClient.Post(
			route(dockerHostAddr, version, tmplres),
			"application/json; charset=utf-8",
			bytes.NewBuffer(jsonbuf),
		)

		if err != nil {
			if netutil.IsTimeOutError(err) {
				log.Println("timeout error")
			}
			panic(err)
		}

		buf, err := jsonutil.BufPretty(res.Body, "", "  ")
		if err != nil {
			panic(err)
		}

		res.Body.Close()
		io.Copy(os.Stdout, &buf)
	}))

	log.Fatal(http.ListenAndServe(port, nil))
}
