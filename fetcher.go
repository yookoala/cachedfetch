package cachedfetcher

import (
	"database/sql"
	"net/http"
	"time"
)

func New(db *sql.DB) *Fetcher {
	return &Fetcher{
		DB: db,
	}
}

type Fetcher struct {
	DB      *sql.DB
	Context Context
}

// set default context string
func (f *Fetcher) ContextStr(ctxStr string) *Fetcher {
	f.Context.Str = ctxStr
	return f
}

// set default context time
func (f *Fetcher) ContextTime(ctxTime time.Time) *Fetcher {
	f.Context.Time = ctxTime
	return f
}

// actually fetch the url with GET method
func (f *Fetcher) Get(url string) (r *http.Response, err error) {
	return http.Get(url)
}
