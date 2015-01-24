package cachedfetcher

import (
	"testing"
)

func TestSqlCache(t *testing.T) {
	var c Cache = &SqlCache{}
	t.Log("SqlCache implements Cache: %#v", c)
}

func TestSqlCacheQuery(t *testing.T) {
	var c CacheQuery = &SqlCacheQuery{}
	t.Log("SqlCacheQuery implements CacheQuery: %#v", c)
}
