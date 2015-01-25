package cachedfetcher

import (
	"time"
)

const (
	OrderContextTime     = iota
	OrderContextTimeDesc = iota
	OrderFetchedTime     = iota
	OrderFetchedTimeDesc = iota
)

type Cache interface {
	Add(url string, ctx Context, r *Response) (err error)

	// find cached response with URL field
	Find(url string) CacheQuery

	// find cached response with context string
	FindIn(str string) CacheQuery
}

type CacheQuery interface {
	ContextStr(Str string) CacheQuery
	ContextTime(t time.Time) CacheQuery
	FetchedTime(t time.Time) CacheQuery
	SortBy(crits ...int) CacheQuery
	Get() (resps []Response, err error)
}
