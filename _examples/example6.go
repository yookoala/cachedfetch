package main

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/buflog"
	"github.com/yookoala/cachedfetcher"
	"time"
)

// gets all cached result and display
func example6(host string, db *sql.DB,
	log *buflog.Logger) (resp *cachedfetcher.Response, err error) {

	log.Print("# Get caches while limiting the number of records")

	url := host + "/example/6"
	c := cachedfetcher.NewSqlCache(db)
	c.Type = dbtype
	f := cachedfetcher.New(c)

	// render context time
	d, err := time.ParseDuration("24h")
	if err != nil {
		return
	}
	t, err := time.Parse(time.RFC822Z, "01 Apr 10 00:00 +0800")
	if err != nil {
		return
	}

	// limits to use
	l1 := 10 // limit in generating cache
	l2 := 5  // limit in retriving response

	for i := 0; i < l1; i++ {
		ctx := cachedfetcher.Context{
			Str:  "example/6",
			Time: t,
		}
		_, err = f.Get(url, ctx)
		if err != nil {
			return
		}
		t = t.Add(d)
	}

	// search the existing url
	rs, err := c.
		Find(url).
		Limit(l2).
		GetAll()
	if err != nil {
		return
	}

	// load response into response slice
	resps := make([]cachedfetcher.Response, 0)
	for rs.Next() {
		resp, err := rs.Get()
		if err != nil {
			log.Fatal("Error getting next response")
		}
		resps = append(resps, *resp)
	}

	// check number of response
	if len(resps) != l2 {
		err = fmt.Errorf("i is %d but expecting %d",
			len(resps), l2)
		return
	}

	log.Printf("Number of results is limited to l2 (%d)", l2)
	return
}
