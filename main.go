package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	log.Printf("listen on 9001")

	http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
