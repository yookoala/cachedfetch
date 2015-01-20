package cachedfetcher

import (
	"testing"
	"time"
)

func TestFetcher(t *testing.T) {

	ctxStr := "some_context"
	ctxTime := time.Now()

	f := &Fetcher{}
	f.ContextStr(ctxStr)
	f.ContextTime(ctxTime)

	if f.Context.Str != ctxStr {
		t.Errorf("Failed to set Context.Str: expect \"%s\" get \"%s\"",
			ctxStr, f.Context.Str)
	}

	if f.Context.Time != ctxTime {
		t.Errorf("Failed to set Context.Time: expect \"%s\" get \"%s\"",
			ctxTime, f.Context.Time)
	}

}
