package main

import (
	"log"
	"net/http"
	"os"
)

const (
	version    = "1.30"
	endpoint   = "/webhooks/payload"
	controller = "CIIO_ROOT_IPADDR"
	port       = ":80"
)

var (
	dockerHostAddr string
)

// example:
// curl --unix-socket /var/run/docker.sock http:/172.17.0.2/v1.30/containers/json

func main() {
	dockerHostAddr = os.Getenv(controller)

	log.Printf("docker host addr, from env var: %s", dockerHostAddr)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		// resp, err := client.Get("")
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
