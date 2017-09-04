package main

import (
	"encoding/json"
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
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func recurseAndPrintJSON(m map[string]interface{}, indent string) {
	for k, v := range m {
		switch cur := v.(type) {
		case map[string]interface{}:
			log.Println(indent, k, ":")
			recurseAndPrintJSON(cur, indent+"\t")
		default:
			log.Println(indent, k, ":", v)
		}
	}
}
