package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/view", snippetView)
	mux.HandleFunc("/create", snippetCreate)

	log.Print("start listening on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
