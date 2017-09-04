package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	endpoint = "/webhooks/payload"
	port     = ":80"
)

func main() {
	log.Printf("listen on %s", port)

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))

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
			recurseJSON(data.(map[string]interface{}), "")
		}
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func prettyPrintJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	return out.Bytes(), err
}

func recurseJSON(m map[string]interface{}, level string) {
	for k, v := range m {
		switch cur := v.(type) {
		case map[string]interface{}:
			log.Println(k, ":")
			recurseJSON(cur, level+"\t")
		default:
			log.Println(level, k, ":", v)
		}
	}
}
