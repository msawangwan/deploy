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

	"github.com/msawangwan/ci.io/model/webhook"
)

const (
	endpoint = "/webhooks/payload"
	port     = ":80"
)

var (
	debug *log.Logger
)

func init() {
	log.Println("app init")
	f, e := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if e != nil {
		log.Fatalln("failed to create/open debug log file:", e)
	}

	label := "[DEBUG]"
	flags := 0

	debug = log.New(f, label, flags)
}

func handlePushEvent(payload *webhook.PushEvent) {
	// startup a container
	// pull the latest code
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("/bin/sh", "./bin/webhook", payload.Repository.FullName, payload.Ref)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Printf("%s\n", err)
		log.Printf("%s\n", stderr.String())
	}

	log.Printf("command executed with result:\n%s\n", out.String())
}

type eventType struct {
	name      string
	guid      string
	signature string
}

func parse(headers http.Header) *eventType {
	return &eventType{
		headers.Get("x-github-event"),
		headers.Get("x-github-delivery"),
		headers.Get("x-github-signature"),
	}
}

func main() {
	log.Printf("listening for incoming webhooks @ %s", port)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		eventName := r.Header.Get("x-github-event")
		eventGUID := r.Header.Get("x-github-delivery")
		eventSig := r.Header.Get("x-github-signature")

		event := parse(r.Header)

		log.Printf("webhook triggered:\n%s\n%s\n%s\n", eventName, eventGUID, eventSig)
		log.Println(event)

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

		// var pretty bytes.Buffer

		// err = json.Indent(&pretty, []byte(body), "", "  ")

		// fmt.Println(pretty.String())

		// var out bytes.Buffer
		// var stderr bytes.Buffer

		// cmd := exec.Command("/bin/sh", "./scripts/webhook.sh")
		// cmd.Stdout = &out
		// cmd.Stderr = &stderr

		// err = cmd.Run()

		// if err != nil {
		// 	log.Printf("%s\n", err)
		// 	log.Printf("%s\n", stderr.String())
		// }

		// log.Printf("command executed with result:\n%s\n", out.String())
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
