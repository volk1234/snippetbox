package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello Snippetworld!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Showing the snippet!"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create super snippet here!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/view", snippetView)
	mux.HandleFunc("/create", snippetCreate)

	log.Print("start listening on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
