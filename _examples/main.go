package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yookoala/cachedfetcher"
	"log"
	"net/http/httptest"
)

var dbdriver, dbsrc *string

type example func(host string, db *sql.DB) (resp *cachedfetcher.Response, err error)

var examples = map[string]example{
	"example1": example1,
	"example2": example2,
	"example3": example3,
	"example4": example4,
	"example5": example5,
	"example6": example6,
}

func init() {

	// read flags
	dbdriver = flag.String("driver", "sqlite3", "Database driver")
	dbsrc = flag.String("db", "file:./cache.db", "Database source")
	flag.Parse()

}

func main() {

	// test server for examples
	ts := httptest.NewServer(ExampleServer())
	defer ts.Close()

	// open database for test
	db, err := sql.Open(*dbdriver, *dbsrc)
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
