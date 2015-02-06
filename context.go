package crawler

import (
	"time"
)

type Context struct {
	Str     string
	Time    time.Time
	Fetched time.Time
}
