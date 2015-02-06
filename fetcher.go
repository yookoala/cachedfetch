package crawler

import (
	"net/http"
	"time"
)

func NewFetcher(c Cache) *Fetcher {
	return &Fetcher{
		Cache: c,
	}
}

type Fetcher struct {
	Cache Cache
}

// actually fetch the url with GET method
func (f *Fetcher) Get(url string, ctx Context) (r *Response, err error) {

	// obtain raw response
	ctx.Fetched = time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	// render crawler response
	r = &Response{
		URL:         url,
		ContextStr:  ctx.Str,
		ContextTime: ctx.Time,
		FetchedTime: ctx.Fetched,
	}
	err = r.ReadRaw(resp)
	if err != nil {
		return
	}

	// add to cache
	err = f.Cache.Add(url, ctx, r)
	return
}
