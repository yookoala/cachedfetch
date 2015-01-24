package cachedfetcher

import (
	"time"
)

type Cache interface {
	Add(url string, ctx Context, r *Response) (err error)
	Find(url string) CacheQuery
}

type CacheQuery interface {
	ContextStr(Str string) CacheQuery
	ContextTime(t time.Time) CacheQuery
	FetchedTime(t time.Time) CacheQuery
	Get() (resps []Response, err error)
}
