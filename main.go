package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	endpoint = "/webhooks/payload"
	port     = ":80"
)

func main() {
	http.HandleFunc(endpoint, server.HandleWebhook)
	log.Fatal(http.ListenAndServe(port, nil))
}