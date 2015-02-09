package sqlcache

import (
	"fmt"
	"github.com/yookoala/crawler"
)

type ResponseColl struct {
	col []crawler.Response
	cur int
}

func (rc *ResponseColl) Next() bool {
	rc.cur++
	if rc.cur <= len(rc.col) {
		return true
	}
	return false
}

func (rc *ResponseColl) Get() (resp *crawler.Response, err error) {
	if rc.cur <= len(rc.col) {
		resp = &rc.col[rc.cur-1]
	} else {
		err = fmt.Errorf("Getting item out of range")
	}
	return
}

func (rc *ResponseColl) Close() (err error) {
	// place holder
	return
}
