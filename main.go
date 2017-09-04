package main

// TODO:
// need to run an external script which
// - pulls the latest code from a repository
// - makes a build

import (
	"bytes"
	"encoding/json"
	"html"
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

func main() {
	log.Printf("listen on %s", port)

	exec.Command("pwd").Run()

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("webhook recieved: %q", html.EscapeString(r.URL.Path))

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
		}

		var data interface{}

		err = json.Unmarshal([]byte(body), &data)

		if err != nil {
			log.Println(err)
		} else {
			recurseAndPrintJSON(data.(map[string]interface{}), "")
		}

		var out bytes.Buffer
		var stderr bytes.Buffer

		cmd := exec.Command("/bin/sh", "./scripts/webhook.sh")
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()

		if err != nil {
			log.Printf("%s\n", err)
			log.Printf("%s\n", stderr.String())
		}

		log.Printf("command executed with result: %s\n", out.String())
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func recurseAndPrintJSON(m map[string]interface{}, indent string) {
	for k, v := range m {
		switch cur := v.(type) {
		case map[string]interface{}:
			debug.Println(indent, k, ":")
			recurseAndPrintJSON(cur, indent+"\t")
		case []interface{}:
			debug.Println(indent, k, ":")
			for _, u := range cur {
				nested, isNested := u.(map[string]interface{})

				if isNested {
					recurseAndPrintJSON(nested, indent+"\t")
				} else {
					debug.Println(indent+"\t", u)
				}
			}
		default:
			debug.Println(indent, k, ":", v)
		}
	}
}
