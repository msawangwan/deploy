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
	"strings"
	"time"
)

const (
	version    = "1.30"
	endpoint   = "/webhooks/payload"
	controller = "CIIO_ROOT_IPADDR"
	port       = ":80"
)

var (
	dockerClient   *http.Client
	dockerHostAddr string
)

func route(adr, ver, src string) string {
	return fmt.Sprintf("http://%s/v%s/%s", adr, ver, src)
}

func pretty(b []byte, delim, indent string) (bytes.Buffer, error) {
	var (
		out bytes.Buffer
		err error
	)

	err = json.Indent(&out, buf, delim, indent)

	return out, err
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
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
}

func main() {
	dockerHostAddr = os.Getenv(controller)

	log.Printf("docker host addr, from env var: %s", dockerHostAddr)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming webhook: %s", r.URL.Path)

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

		// io.Copy(os.Stdout, res.Body)
		res.Body.Close()
		buf, err := pretty([]byte(body), "", "  ")
		io.Copy(os.Stdout, &buf)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
