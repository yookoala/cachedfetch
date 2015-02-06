package crawler

import (
	"testing"
	"time"
)

func TestResponse(t *testing.T) {
	now := time.Now()
	r1 := &Response{
		// context and fetch information
		URL:         "dummy url 1",
		ContextStr:  "dummy context 1",
		ContextTime: now,
		FetchedTime: now,

		// response meta information
		Status:               "220 Dummy Status",
		StatusCode:           220,
		Proto:                "Dummy Proto",
		ContentLength:        10,
		TransferEncodingJson: []byte("{\"Transfer\":\"encoding\"}"),
		HeaderJson:           []byte("{\"Header\":\"value\"}"),
		TrailerJson:          []byte("{\"Trailer\":\"value\"}"),
		RequestJson:          []byte("{\"Request\":\"value\"}"),
		TlsJson:              []byte("{\"TLS\":\"value\"}"),

		// response body
		Body: []byte("Dummy Body"),
	}

	// copy r1 values to r2
	var r2 Response
	r2 = *r1

	if !r1.Equals(r2) {
		t.Error("Response are not evaluated as the same")
	}

	r2.URL = "dummy url 1 changed"
	if r1.Equals(r2) {
		t.Error("r2 has changed and Response.Equals doesn't recgonize it")
	}

}
