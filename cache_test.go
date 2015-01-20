package cachedfetcher

import (
	"testing"
)

func TestSqlCache(t *testing.T) {
	var c Cache = &SqlCache{}
	t.Log("SqlCache implements Cache: %#v", c)
}
