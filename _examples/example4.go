package main

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/cachedfetcher"
	"log"
	"time"
)

// gets all cached result and display
func example4(host string, db *sql.DB) (resp *cachedfetcher.Response, err error) {

	log.Print("# Get old cache by context string")

	url := host + "/example/4"
	c := cachedfetcher.NewSqlCache(db)
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
	l := 10

	for i := 1; i <= l; i++ {
		ctx := cachedfetcher.Context{
			Str:  "example/4",
			Time: t,
		}
		_, err = f.Get(fmt.Sprintf("%s/%d", url, i), ctx)
		if err != nil {
			return
		}
		t = t.Add(d)
	}

	// search the existing url
	rs, err := c.FindIn("example/4").GetAll()
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

	// get cached items and display
	for i, resp := range resps {
		log.Printf("[#%d] (%s) URL: \"%s\"", i,
			resp.ContextTime.Format("2006-01-02"),
			string(resp.URL))
	}

	// check number of response
	if len(resps) != l {
		err = fmt.Errorf("i is %d but expecting %d",
			len(resps), l)
	}
	return
}
