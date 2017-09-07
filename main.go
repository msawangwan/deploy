package main

import (
	"log"
	"net/http"
)

const (
	endpoint = "/webhooks/payload"
	port     = ":80"
)

func main() {
	http.HandleFunc(endpoint, server.HandleWebhook)
	log.Fatal(http.ListenAndServe(port, nil))
}
