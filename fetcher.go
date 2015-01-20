package cachedfetcher

import (
	"net/http"
)

func New(c Cache) *Fetcher {
	return &Fetcher{
		Cache: c,
	}
}

type Fetcher struct {
	Cache Cache
}

// actually fetch the url with GET method
func (f *Fetcher) Get(url string) (r *http.Response, err error) {
	return http.Get(url)
}
