package main

import (
	"github.com/yookoala/cachedfetcher"
	"log"
)

func example1(host string) (err error) {
	url := host + "/example/1"
	resp, err := cachedfetcher.Get(url)
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
