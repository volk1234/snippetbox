package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      *string
	staticDir *string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	config := config{}
	config.addr = flag.String("addr", ":4000", "HTTP address to run web server")
	config.staticDir = flag.String("dir", "./ui/static/", "Path to static assests")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(*config.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/view", app.snippetView)
	mux.HandleFunc("/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *config.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("start listening on port %s", *config.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
