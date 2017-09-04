package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
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

		var data struct {
			String string
		}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)

		if err != nil {
			log.Printf("err %s", err)
		}

		fmt.Println(data.String)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func prettyPrintJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	return out.Bytes(), err
}
