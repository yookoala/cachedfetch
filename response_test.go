package crawler

import (
	"math/rand"
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

	if !r1.Equal(r2) {
		t.Error("Response are not evaluated as the same")
	}

	r2.URL = "dummy url 1 changed"
	if r1.Equal(r2) {
		t.Error("r2 has changed and Response.Equal doesn't recgonize it")
	}

}

func TestResponseSetCtx(t *testing.T) {

	resp := &Response{}
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	t1 := time.Date(2015, time.Month(rand.Intn(12)+1), rand.Intn(20)+1,
		0, 0, 0, 0, loc)
	t2 := time.Date(2015, time.Month(rand.Intn(12)+1), rand.Intn(20)+1,
		0, 0, 0, 0, loc)
	ctx := &Context{
		Str:     "dummy context TestResponseSetCtx",
		Time:    t1,
		Fetched: t2,
	}

	// test before setting
	if resp.InContext(ctx) {
		t.Errorf("Empty response should not be in random context.\n"+
			"ctx: %#v\nresp: %#v", ctx, resp)
	}

	// set context
	resp.SetContext(ctx)

	// test after setting
	if !resp.InContext(ctx) {
		t.Errorf("Response should be in context set to it.\n"+
			"ctx: %#v\nresp: %#v", ctx, resp)
	}

}

func TestResponseGetCtx(t *testing.T) {

	resp := &Response{}
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	t1 := time.Date(2015, time.Month(rand.Intn(12)+1), rand.Intn(20)+1,
		0, 0, 0, 0, loc)
	t2 := time.Date(2015, time.Month(rand.Intn(12)+1), rand.Intn(20)+1,
		0, 0, 0, 0, loc)
	ctx1 := &Context{
		Str:     "dummy context TestResponseGetCtx",
		Time:    t1,
		Fetched: t2,
	}

	// set context
	resp.SetContext(ctx1)

	// get context back
	ctx2 := resp.GetContext()

	if !ctx1.Equal(ctx2) {
		t.Errorf("Context set and get are different.\n"+
			"ctx1: %#v\nctx2: %#v", ctx1, ctx2)
	}

}
