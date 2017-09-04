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

		log.Printf("webhook triggered: %q", html.EscapeString(r.URL.Path))

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
		}

		// var data map[string] interface{}

		// err := json.Unmarshal(body, &data)

		// if err != nil {
		// 	fmt.Printf(err)
		// }

		formatted, err := prettyPrintJSON([]byte(body))

		if err != nil {
			log.Println(err)
		} else {
			log.Println(formatted)
		}

		// var data struct {
		// 	String string
		// }

		// decoder := json.NewDecoder(r.Body)

		// err := decoder.Decode(&data)

		// if err != nil {
		// 	log.Printf("err %s", err)
		// }

		// fmt.Println(data.String)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func prettyPrintJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	return out.Bytes(), err
}
