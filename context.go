package crawler

import (
	"time"
)

type Context struct {
	Str     string
	Time    time.Time
	Fetched time.Time
}

func (ctx *Context) Equal(ctx2 *Context) bool {
	return ctx.Str == ctx2.Str &&
		ctx.Time.Equal(ctx2.Time) &&
		ctx.Fetched.Equal(ctx2.Fetched)
}
