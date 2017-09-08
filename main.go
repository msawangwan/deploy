package main

import (
	"context"
	"fmt"
	"io"
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

func concat(adr, ver, src string) string {
	return fmt.Sprintf("http:/%s/%s/%s", adr, ver, src)
}

func read(r io.Reader) {
	buf := make([]byte, 1024)

	for {
		n, e := r.Read(buf[:])

		if e != nil {
			return
		}

		log.Printf("client response: %s\n", string(buf[0:n]))
	}
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
				return net.Dial("unix", "var/run/docker.sock")
			},
		},
	}
}

func main() {
	dockerHostAddr = os.Getenv(controller)

	log.Printf("docker host addr, from env var: %s", dockerHostAddr)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		//        c, e := net.Dial("unix", "/var/run/docker.sock")

		//        if e != nil {
		//            panic(e)
		//        }

		//        defer c.Close()
		//        go read(c)

		//        _, e = c.write([]byte(concat(dockerHostAddr, version, "containers/json"))

		//        if e != nil {
		//            log.Println(e)
		//        }

		//        c := &http.Client{
		//            Transport: &http.Transport{
		//                DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
		//                    return net.Dial("unix", "/var/run/docker.sock"
		//                }
		//            }
		//        }

		var (
			res *http.Response
			err error
		)

		res, err = dockerClient.Get(concat(dockerHostAddr, version, "containers/json"))

		if err != nil {
			if timedOut(err) {
				log.Println("timeout error")
			}
			panic(err)
		}

		io.Copy(os.Stdout, res.Body)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
