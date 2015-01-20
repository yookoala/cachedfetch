package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/http/httptest"
)

var dbfile = flag.String("db", "cache.db", "The SQLite3 database file name")

func main() {

	//
	flag.Parse()

	// test server for examples
	mux := http.NewServeMux()
	mux.Handle("/example/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, example 1")
	}))
	mux.Handle("/example/2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, example 2")
	}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// open database for test
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}

	// run examples with test server
	log.Printf("----")
	err = example1(ts.URL, db)
	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("----")
	err = example2(ts.URL, db)
	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("----")
}
