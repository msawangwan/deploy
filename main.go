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

func pretty(buf []byte, delim, indent string) (bytes.Buffer, error) {
	var (
		out bytes.Buffer
		err error
	)

	err = json.Indent(&out, buf, delim, indent)

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
					return addrstr, nil
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

		res, err = dockerClient.Get(route(dockerHostAddr, version, "containers/json"))

		if err != nil {
			if timedOut(err) {
				log.Println("timeout error")
			}
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			panic(err)
		}

		res.Body.Close()
		buf, err := pretty([]byte(body), "", "  ")
		io.Copy(os.Stdout, &buf)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
