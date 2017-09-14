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
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/msawangwan/ci.io/api/ciio"
	"github.com/msawangwan/ci.io/api/github"
	"github.com/msawangwan/ci.io/lib/jsonutil"
	"github.com/msawangwan/ci.io/types/cred"
	//	"github.com/msawangwan/ci.io/util"
)

const (
	version       = "1.30"
	endpoint      = "/webhooks/payload"
	mountpoint    = "/var/run/docker.sock"
	controller    = "CIIO_ROOT_IPADDR"
	port          = ":80"
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

func route(adr, ver, src string) string {
	return fmt.Sprintf("http://%s/v%s/%s", adr, ver, src)
}

func pretty(r io.Reader, delim, indent string) (bytes.Buffer, error) {
	var (
		out bytes.Buffer
		err error
	)

	src, err := ioutil.ReadAll(r)

	if err != nil {
		return out, err
	}

	err = json.Indent(&out, []byte(src), delim, indent)

	return out, err
}

func localIP(ifname string) (string, error) {
	intfs, e := net.Interfaces()

	if e != nil {
		return "none", e
	}

	for _, intf := range intfs {
		if strings.Contains(intf.Name, ifname) {
			addrs, e := intf.Addrs()

			if e != nil {
				return "none", e
			}

			for _, addr := range addrs {
				addrstr := addr.String()
				if !strings.Contains(addrstr, "[") {
					return strings.Split(addrstr, "/")[0], nil
				}
			}
		}
	}

	return "none", nil
}

func timedOut(e error) bool {
	switch e := e.(type) {
	case *url.Error:
		if e, ok := e.Err.(net.Error); ok && e.Timeout() {
			return true
		}
	case net.Error:
		if e.Timeout() {
			return true
		}
	case *net.OpError:
		if e.Timeout() {
			return true
		}
	}

	if e != nil {
		if strings.Contains(e.Error(), "use of closed network connection") {
			return true
		}
	}

	return false
}

func init() {
	var (
		err error
	)

	err = jsonutil.FromFile("secret/github.auth.json", credentials)
	if err != nil {
		log.Printf("%s", err)
	}

	err := os.Mkdir(scratchdir, 655)
	if err != nil {
		log.Printf("%s", err)
	}

	err = os.Chdir(scratchdir)
	if err != nil {
		log.Printf("%s", err)
	} else {
		wd, _ := os.Getwd()
		log.Printf("working dir: %s", wd)
	}

	dockerClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(socktype, mountpoint)
			},
		},
	}

	dockerHostAddr = os.Getenv(controller)

	localip, err = localIP("eth0")
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

		defer os.RemoveAll(tmpdir)

		log.Printf("created tmp workspace: %s", tmpdir)

		/* clone the remote repo into temp workspace */

		var (
			repouser  string = "user"
			reponame  string = "repository"
			repoowner string = payload.Repository.Owner.Name
			cloneurl  string = payload.Repository.CloneURL
			cmdout    bytes.Buffer
			cmderr    bytes.Buffer
		)

		clone := exec.Command(commands.cloneRemoteRepo, repoowner, cloneurl)
		clone.Dir = tmpdir
		clone.Stdout = &cmdout
		clone.Stderr = &cmderr

		if err = clone.Run(); err != nil {
			log.Printf("%s", err)
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
			if timedOut(err) {
				log.Println("timeout error")
			}
			panic(err)
		}

		buf, err := pretty(res.Body, "", "  ")
		if err != nil {
			panic(err)
		}

		res.Body.Close()
		io.Copy(os.Stdout, &buf)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
