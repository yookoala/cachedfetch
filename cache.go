package cachedfetcher

import (
	"database/sql"
	"net/http"
)

type Cache interface {
	Add(url string, ctx Context, resp *http.Response) (err error)
}

func NewSqlCache(db *sql.DB) *SqlCache {
	return &SqlCache{
		DB: db,
	}
}

type SqlCache struct {
	DB *sql.DB
}

func (c *SqlCache) Add(url string, ctx Context,
	resp *http.Response) (err error) {
	return
}
