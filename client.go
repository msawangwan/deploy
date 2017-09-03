package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func sendreq() {
	log.Printf("hello mate\n")

	// res, err := http.Get("https://github.com/repos/msawangwan/tyrant-index.io/hooks")
	res, err := http.Get("https://github.com/repos/msawangwan/poker.io/hooks")

	if err != nil {
		log.Printf("%s", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("%s", err)
	}

	log.Printf("%s", body)
}
