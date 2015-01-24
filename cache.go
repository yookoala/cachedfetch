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
	stmt, err := c.DB.Prepare("INSERT INTO `cachedfetch_cache` " +
		"(`url`, `context`, `context_time`, `fetched`, " +
		"`status`, `status_code`, `proto`, `content_length`, " +
		"`transfer_encoding`, `header`, `trailer`, `request`, `body`)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	_, err = stmt.Exec(
		r.URL,
		r.Context,
		r.ContextTime.Unix(),
		r.Fetched.Unix(),
		r.Status,
		r.StatusCode,
		r.Proto,
		r.ContentLength,
		r.TransferEncodingJson,
		r.HeaderJson,
		r.TrailerJson,
		r.RequestJson,
		r.Body)

	return
}
