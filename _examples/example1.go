package main

import (
	"github.com/yookoala/buflog"
	"github.com/yookoala/cachedfetcher"
	"time"
)

func example1(host string, c cachedfetcher.Cache,
	log *buflog.Logger) (resp *cachedfetcher.Response, err error) {

	log.Print("# Fetch a URL and retrieve from cache")

	url := host + "/example/1"
	f := cachedfetcher.New(c)
	now := time.Now()
	ctx := cachedfetcher.Context{
		Str:  "example/1",
		Time: now,
	}
	resp, err = f.Get(url, ctx)
	if err != nil {
		return
	}

	// search for previous cache of the URL in the context
	rs, err := c.
		Find(url).
		In(ctx.Str).
		At(now).
		FetchedAt(now).
		GetAll()
	if err != nil {
		return
	}

	// log original response
	log.Printf("- original response -")
	log.Printf("URL:    %s", resp.URL)
	log.Printf("Status: %s", resp.Status)
	log.Printf("Size:   %d", resp.ContentLength)
	log.Printf("Body:   \"%s\"", string(resp.Body))

	// load response into response slice
	resps := make([]cachedfetcher.Response, 0)
	for rs.Next() {
		resp, err := rs.Get()
		if err != nil {
			log.Fatal("Error getting next response")
		}
		resps = append(resps, *resp)
	}

	// check the cached responses
	if len(resps) == 0 {
		log.Fatal("Could not find example 1 response in cache")
	} else if len(resps) > 1 {
		log.Fatal("More than 1 responses matches but expecting only 1")
	} else if !resp.Equals(resps[0]) {
		log.Fatal("Response found in cache is different from the one stored previously")
	}

	// log response
	log.Printf("- cached response -")
	log.Printf("URL:    %s", resps[0].URL)
	log.Printf("Status: %s", resps[0].Status)
	log.Printf("Size:   %d", resps[0].ContentLength)
	log.Printf("Body:   \"%s\"", string(resps[0].Body))
	return
}
