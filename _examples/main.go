package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yookoala/buflog"
	"github.com/yookoala/cachedfetcher"
	"log"
	"net/http/httptest"
	"sync"
)

var dbdriver, dbsrc *string
var dbtype int

type example func(host string, db *sql.DB,
	log *buflog.Logger) (resp *cachedfetcher.Response, err error)

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

	// set dbtype according to driver name
	switch *dbdriver {
	case "postgres":
		dbtype = cachedfetcher.SQL_PSQL
	case "sqlite3":
		dbtype = cachedfetcher.SQL_SQLITE3
	default:
		dbtype = cachedfetcher.SQL_MYSQL
	}

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
	defer db.Close()

	// test database connection
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Printf("Unable to connect to database")
		log.Fatal(err)
	}

	// initialize wait group
	wg := &sync.WaitGroup{}
	ch := make(chan *buflog.Logger)
	done := make(chan bool)

	// run examples with test server
	for name, exp := range examples {
		wg.Add(1)
		go func(name string, exp example) {
			log.Printf("** %s start", name)
			defer wg.Done()
			lr := buflog.New()
			lr.Printf("#### %s ####", name)
			resp, err := exp(ts.URL, db, lr)
			if err != nil {
				lr.Printf("** %s: error", name)
				if resp != nil {
					lr.Printf("Response Size: %d", resp.ContentLength)
					lr.Printf("Response Body: %s", resp.Body)
				}
				lr.Fatalf("Error Message: %s", err)
			}
			log.Printf("** %s end", name)
			ch <- lr
		}(name, exp)
	}

	// wait for the wait group to finish
	// and send the done signal
	go func() {
		wg.Wait()
		done <- true
	}()

	// loop and wait for all example to end
	finished := false
	for !finished {
		select {
		case lr := <-ch:
			lr.Play()
		case finished = <-done:
			log.Printf("##################")
		}
	}
}
