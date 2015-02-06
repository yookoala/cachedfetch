package crawler

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

	// find cached response with context time
	FindAt(t time.Time) CacheQuery
}

type CacheQuery interface {

	// add context string as condition
	In(Str string) CacheQuery

	// add context time as condition
	At(t time.Time) CacheQuery

	// add fetch time as condition
	FetchedAt(t time.Time) CacheQuery

	// add sorting requirement(s)
	SortBy(crits ...int) CacheQuery

	// limit the number of cached response to retrieve
	Limit(int) CacheQuery

	// execute the query
	GetAll() (resps ResponseColl, err error)
}

type ResponseColl interface {

	// load the next response
	// must be called before each Get() call
	Next() bool

	// get 1 single response
	Get() (resp *Response, err error)

	// close the response collection
	Close() (err error)
}
