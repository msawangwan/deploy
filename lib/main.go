package main

// TODO:
// need to run an external script which
// - pulls the latest code from a repository
// - makes a build

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
   "context"

	"github.com/msawangwan/ci.io/lib/internal/webhook"
	
	"github.com/moby/moby/client"
)

const (
	endpoint = "/webhooks/payload"
	port     = ":80"
)

type eventHeaders struct {
	name      string
	guid      string
	signature string
}

func (eh eventHeaders) String() string { return fmt.Sprintf("webhook event\nname: %s\nguid: %s\nsignature: %s\n", eh.name, eh.guid, eh.signature) }

func parseWebhookRequest(h http.Header) *eventHeaders {
	return &eventHeaders{
		h.Get("x-github-event"),
		h.Get("x-github-delivery"),
		h.Get("x-github-signature"),
	}
}

func printErr(e error, s string) {
	if e == nil {
		e = "error"
	}

	log.Printf("error\n%s\n%s", e, s)
}

func init() {
	log.Println("app init")
}

func handlePushEvent(payload *webhook.PushEvent) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		err error
	)

	cmd := exec.Command("/bin/sh", "./bin/webhook", payload.Repository.FullName, payload.Ref)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		printErr(err, stderr.String())
	}

	log.Printf("command executed with result:\n%s\n", stdout.String())
}

func main() {
	log.Printf("listening for incoming webhooks @ %s", port)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		event := parseWebhookRequest(r.Header)

		log.Printf("incoming:\n%+v", event.String())

		if eventName == "push" {
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println(err)
			}

			var payload *webhook.PushEvent

			err = json.Unmarshal([]byte(body), &payload)

			if err != nil {
				log.Println(err)
			} else {
				handlePushEvent(payload)
			}
		}
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
