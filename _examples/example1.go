package main

import (
	"database/sql"
	"github.com/yookoala/cachedfetcher"
	"log"
)

func example1(host string, db *sql.DB) (err error) {
	url := host + "/example/1"
	c := cachedfetcher.NewSqlCache(db)
	f := cachedfetcher.New(c)
	resp, err := f.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// read to byte according to size
	var b []byte
	b = make([]byte, resp.ContentLength)
	size, err := resp.Body.Read(b)
	if err != nil {
		return
	}

	// log response
	log.Printf("Host:   %s", host)
	log.Printf("URL:    %s", url)
	log.Printf("Status: %s", resp.Status)
	log.Printf("Size:   %d", size)
	log.Printf("Body:   \"%s\"", string(b))
	return
}
