package cachedfetcher

import (
	"net/http"
)

func Get(url string) (r *http.Response, err error) {
	return http.Get(url)
}
