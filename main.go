package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

const (
	kENDPOINT = "/webhooks/payload"
	kPORT     = ":80"
)

func main() {
	log.Printf("listen on %s", kPORT)

	http.HandleFunc(kENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(kPORT, nil))
}
