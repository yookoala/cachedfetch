package cachedfetcher

import (
	"database/sql"
)

type Cache interface {
	Add(url string, ctx Context, r *Response) (err error)
}

func NewSqlCache(db *sql.DB) *SqlCache {
	return &SqlCache{
		DB: db,
	}
}

type SqlCache struct {
	DB *sql.DB
}

func (c *SqlCache) Add(url string, ctx Context, r *Response) (err error) {
	// TODO: store the response into database
	return
}
