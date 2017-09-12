package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
)

const (
	version    = "1.30"
	endpoint   = "/webhooks/payload"
	mountpoint = "/var/run/docker.sock"
	controller = "CIIO_ROOT_IPADDR"
	port       = ":80"
	socktype   = "unix"
)

var (
	dockerClient   *http.Client
	dockerHostAddr string
	localip        string
)

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
	dockerClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(socktype, mountpoint)
			},
		},
	}

	dockerHostAddr = os.Getenv(controller)
	localip, err := localIP("eth0")

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
		gocount = runtime.NumGoroutine()
		ename := r.Header.Get("x-github-event")

		log.Printf("incoming webhook: %s\n", r.URL.Path)
		log.Printf("payload event name: %s\n", ename)
		log.Printf("goroutine count: %d\n", gocount)

		var (
			res *http.Response
			err error
		)

		// res, err = dockerClient.Get(route(dockerHostAddr, version, "containers/json"))

		// create a container

		// payload := &dockr.Container{
		// 	Image: "alpine",
		// 	CMD: []string{
		// 		"date",
		// 	},
		// }

		payload := []byte(`{"Image":"alpine", "Cmd": ["date"]}`)

		// var b bytes.Buffer
		// jsonpayload, _ := json.Marshal(payload)

		// json.NewEncoder(b).Encode(payload)

//        var query struct = {
//            endpoint string
//            params string
//        }, {
//            "containers/create",
//            "SOME_CONTAINER",
//        }

//        tmpl, err := template.New("create_endpoint").Parse("{{.endpoint}}?name={{.params}}")

//        if err != nil {
//            panic(err)
//        }

//        err = tmpl.Execute()

		res, err = dockerClient.Post(
			route(dockerHostAddr, version, "containers/create?name=SOME_CONTAINER"),
			"application/json; charset=utf-8",
			bytes.NewBuffer(payload),
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
