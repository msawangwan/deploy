package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/msawangwan/ci.io/api/ciio"
	"github.com/msawangwan/ci.io/api/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/lib/netutil"
	"github.com/msawangwan/ci.io/types/cred"
)

const (
	version       = "1.30"
	port          = ":80"
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

var commands = struct {
	cloneRemoteRepo string
}{
	"clrep",
}

var pwd = func() { d, _ := os.Getwd(); log.Printf("current working dir: %s", d) }
var route = func(adr, ver, src string) string { return fmt.Sprintf("http://%s/v%s/%s", adr, ver, src) }

func init() {
	rootdir, _ := os.Getwd()
	pathenv := os.Getenv("PATH")
	newpathenv := fmt.Sprintf("%s:%s/bin", pathenv, rootdir)

	os.Setenv("PATH", newpathenv)

	log.Printf("old path: %s", pathenv)
	log.Printf("new path: %s", os.Getenv("PATH"))

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
		pwd()
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

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		/* some stats */

		gocount = runtime.NumGoroutine()

		log.Printf("goroutine count: %d\n", gocount)

		/* parse request header */

		ename := r.Header.Get("x-github-event")

		log.Printf("incoming webhook: %s\n", r.URL.Path)
		log.Printf("payload event name: %s\n", ename)

		if ename != "push" {
			log.Printf("cannot handle event: %s", ename)
			return
		}

		/* parse webhook json payload */

		var (
			res     *http.Response
			payload *github.PushEvent
			body    []byte
			err     error
		)

		body, err = ioutil.ReadAll(r.Body)

		if err != nil {
			log.Printf("%s", err)
		}

		err = json.Unmarshal(body, &payload)

		if err != nil {
			log.Printf("%s", err)
		}

		/* parse project name from payload */

		var (
			repo     string
			projname string
		)

		repo = payload.Repository.HTMLURL
		projname = payload.Repository.Name

		log.Printf("project name: %s", projname)
		log.Printf("repo: %s", repo)

		/* pull the repo into a tmp dir */

		var (
			tmpdir       string
			tmpdirpath   string
			tmpdirprefix string
		)

		tmpdirpath = "./"
		tmpdirprefix = projname

		tmpdir, err = ioutil.TempDir(tmpdirpath, tmpdirprefix)
		if err != nil {
			log.Printf("%s", err)
		}

		//		defer os.RemoveAll(tmpdir)

		log.Printf("created tmp workspace: %s", tmpdir)

		/* clone the remote repo into temp workspace */

		pwd()

		var (
			cmdout bytes.Buffer
			cmderr bytes.Buffer
		)

		clone := exec.Command(
			commands.cloneRemoteRepo,
			credential.User,
			credential.OAuth,
			payload.Repository.CloneURL,
		)

		clone.Dir = tmpdir
		clone.Stdout = &cmdout
		clone.Stderr = &cmderr

		if err = clone.Run(); err != nil {
			log.Printf("err when calling exec: %s", err)
			log.Printf("%s", cmderr.String())
		} else {
			log.Printf("no err when calling exec: %+v", commands.cloneRemoteRepo)
			log.Printf("%s", cmdout.String())
		}

		/* create the container command */

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
				"name": "SOME_CONTAINER",
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

		/* find the local buildfile */

		var (
			buildfile *os.File
		)

		buildfile, err = os.Open(strings.ToLower(buildfilename))
		if err != nil {
			log.Printf("%s", err)
		}

		parsed := json.NewDecoder(buildfile)

		var buildfilepayload ciio.Buildfile

		if err = parsed.Decode(&buildfilepayload); err != nil {
			log.Printf("%s", err)
		}

		log.Printf("project build params: %+v", buildfilepayload)

		jsonbuf := []byte(
			`{
				"Image":"golang:1.9.0-alpine3.6",
				"WorkingDir": "/app",
				"Cmd": ["date"]
			 }`,
		)

		/* query the docker host and create container */

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
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
