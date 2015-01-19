package main

import (
	"flag"
	"fmt"
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

	// run examples with test server
	log.Printf("----")
	example1(ts.URL)
	log.Printf("----")
	example2(ts.URL)
	log.Printf("----")
}
