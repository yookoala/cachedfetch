package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yookoala/cachedfetcher"
	"log"
	"net/http"
	"net/http/httptest"
)

var dbfile = flag.String("db", "cache.db", "The SQLite3 database file name")

type example func(host string, db *sql.DB) (resp *cachedfetcher.Response, err error)

var examples = map[string]example{
	"example1": example1,
	"example2": example2,
}

func ExampleServer() (mux *http.ServeMux) {

	mux = http.NewServeMux()

	// produce count with a channel
	counts := make(chan int64)
	go func() {
		for i := int64(1); ; i++ {
			counts <- i
		}
	}()

	// bind example paths
	mux.Handle("/example/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, example 1")
	}))
	mux.Handle("/example/2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := <-counts
		fmt.Fprintf(w, "Counter: %03d", count)
	}))
	return
}

func main() {

	// parse to get db file name
	flag.Parse()

	// test server for examples
	ts := httptest.NewServer(ExampleServer())
	defer ts.Close()

	// open database for test
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}

	// run examples with test server
	for name, exp := range examples {
		log.Printf("#### %s ####", name)
		resp, err := exp(ts.URL, db)
		if err != nil {
			log.Printf("*** Error")
			if resp != nil {
				log.Printf("Response Size: %d", resp.ContentLength)
				log.Printf("Response Body: %s", resp.Body)
			}
			log.Fatalf("Error Message: %s", err)
		}
	}
	log.Printf("##################")
}
