package cachedfetcher

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {

	// context and fetch information
	URL         string
	Context     string
	ContextTime time.Time
	Fetched     time.Time

	// response meta information
	Status               string
	StatusCode           int
	Proto                string
	ContentLength        int64
	TransferEncodingJson []byte
	HeaderJson           []byte
	TrailerJson          []byte
	RequestJson          []byte
	TlsJson              []byte

	// response body
	Body []byte
}

func (r *Response) ReadRaw(raw *http.Response) (err error) {

	// read basic information
	r.Status = raw.Status
	r.StatusCode = raw.StatusCode
	r.Proto = raw.Proto
	r.ContentLength = raw.ContentLength

	// read complex fields to JSON
	r.TransferEncodingJson, err = json.Marshal(raw.TransferEncoding)
	if err != nil {
		return
	}
	r.HeaderJson, err = json.Marshal(raw.Header)
	if err != nil {
		return
	}
	r.TrailerJson, err = json.Marshal(raw.Trailer)
	if err != nil {
		return
	}
	r.RequestJson, err = json.Marshal(*raw.Request)
	if err != nil {
		return
	}
	/*
		// no such field or method in Golang 1.2
		r.TlsJson, err = json.Marshal(raw.TLS)
		if err != nil {
			return
		}
	*/

	// read body
	r.Body = make([]byte, raw.ContentLength)
	_, err = raw.Body.Read(r.Body)
	if err != nil {
		return
	}

	return
}
