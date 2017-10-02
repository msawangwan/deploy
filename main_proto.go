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
	version        = "1.30"
	port           = ":80"
	mime           = "application/json; charset=utf-8"
	endpoint       = "/webhooks/payload"
	mountpoint     = "/var/run/docker.sock"
	envipaddr      = "DOCK_MASTERCONTAINER_IPADDR"
	socktype       = "unix"
	scratchdir     = "__ws"
	buildfilename  = "buildfile.json"
	dockerfilepath = "Dockerfile"
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
	errDoesNotExistInCache = errors.New("does not exist in cache")
	errIDMismatch          = errors.New("expected id doesnt match id")
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

var killsig = make(chan os.Signal, 1)

var (
	dirCache dir.WorkspaceCacher
)

func init() {
	signal.Notify(killsig, syscall.SIGINT, syscall.SIGTERM)

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

	dirCache = &dir.WorkspaceTable{cache: map[string]string{}}

	log.Printf("server container ip: %s\n", localip)
	log.Printf("docker host container ip: %s\n", dockerHostAddr)
}

func pwd(s string)                      { d, _ := os.Getwd(); log.Printf("[current working dir %s] %s", d, s) }
func route(adr, ver, src string) string { return fmt.Sprintf("http://%s/v%s/%s", adr, ver, src) }
func apiurl(resource string) string     { return route(dockerHostAddr, version, resource) }

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

type directories struct {
	workspace    string
	scratchspace string
}

func createWorkAndScratchSpace(cache dir.WorkspaceCacher, name string) (ds directories, er error) {
	ws, er := cache.MkTempDir(name)
	if er != nil {
		return er
	}

	sc, er := cache.MkTempDir(name + "scratch")
	if er != nil {
		return er
	}

	ds = directories{ws, sc}

	return
}

func buildRepo(ds directories, c credentials, repo string) error {
	var stdout, stderr bytes.Buffer

	args := []string{
		ds.workspace,
		ds.scratchspace,
		c.User,
		c.OAuth,
		repo,
	}

	cmd := exec.Command("makebarer", args)

	cmd.Dir = "__ws"
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er = cmd.Run(); er != nil {
		log.Printf("%s", stderr.String())
		return er
	}

	log.Printf("%s", stdout.String())

	return nil
}

func buildImage(client *http.Client, tag, remote string) error {
	params := map[string]string{
		"remote": remote,
		"t":      tag,
	}

	cmd := dock.NewBuildDockerfileCommand(params)
	concat, er := dock.BuildAPIURLString(cmd)
	if er != nil {
		return er
	}

	uri := apiurl(concat)

	log.Printf("build api url: %s", uri)

	req, er := client.Post(client, mime, io.Reader(nil))
	if er != nil {
		return er
	}

	if req.StatusCode != 200 {
		var m dock.ErrorResponse

		if er = jsonutil.FromReader(req.Body, &m); er != nil {
			return er
		}

		return responseCodeMismatchError{200, req.StatusCode, m.Message}
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

		if ds, er := createWorkAndScratchSpace(dirCache, repoName); er != nil {
			panic(er)
		}

		// need to get the url to the repo from here
		if er = buildRepo(ds, credential, webhook.Repository.URL); er != nil {
			panic(er)
		}

		if er = buildImage(dockerClient, repoName, ""); er != nil {
			panic(er)
		}

		log.Printf("webhook event, handled")

	}))

	go func() {
		<-killsig

		killContainer := func(id string, c *http.Client) error {
			var e error

			if e = stopContainer(id, c); e != nil {
				return e
			}

			if e = removeContainer(id, c); e != nil {
				return e
			}

			return nil
		}

		log.Printf("kill all running containers")

		containerCache.Lock()
		{
			for k, v := range containerCache.store {
				log.Printf("killing container: %s", k)

				if e := killContainer(v, dockerClient); e != nil {
					panic(e)
				}
			}
		}
		containerCache.Unlock()

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
