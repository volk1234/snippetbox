package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.sandx.link/internal/models"
)

type config struct {
	addr *string
	dsn  *string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	config := config{}
	config.addr = flag.String("addr", ":4000", "HTTP address to run web server")
	config.dsn = flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*config.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *config.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("start listening on port %s", *config.addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Printf("connected to DB with %s\n", dsn)
	return db, nil
}
