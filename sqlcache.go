package cachedfetcher

import (
	"database/sql"
	"strings"
	"time"
)

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
		"(`url`, `context_str`, `context_time`, `fetched_time`, " +
		"`status`, `status_code`, `proto`, `content_length`, " +
		"`transfer_encoding`, `header`, `trailer`, `request`, `body`)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	_, err = stmt.Exec(
		r.URL,
		r.ContextStr,
		r.ContextTime.Unix(),
		r.FetchedTime.Unix(),
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

func (c *SqlCache) Find(url string) CacheQuery {
	return &SqlCacheQuery{
		URL:   url,
		Cache: c,
	}
}

type SqlCacheQuery struct {
	URL     string
	Context Context
	Cache   *SqlCache
}

func (q *SqlCacheQuery) ContextStr(Str string) CacheQuery {
	q.Context.Str = Str
	return q
}

func (q *SqlCacheQuery) ContextTime(t time.Time) CacheQuery {
	q.Context.Time = t
	return q
}

func (q *SqlCacheQuery) FetchedTime(t time.Time) CacheQuery {
	q.Context.Fetched = t
	return q
}

func (q *SqlCacheQuery) Get() (resps []Response, err error) {
	sql := "SELECT `url`, `context_str`, `context_time`, " +
		"`fetched_time`, `status`, `status_code`, `proto`, " +
		"`content_length`, `transfer_encoding`, " +
		"`header`, `trailer`, `request`, " +
		"`body` " +
		"FROM `cachedfetch_cache` "
	sortClause := " ORDER BY `context_time` DESC, `fetched_time` DESC"

	var t time.Time

	// parameters to build query
	w := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)

	if q.URL != "" {
		w = append(w, "`url` = ?")
		args = append(args, q.URL)
	}
	if q.Context.Str != "" {
		w = append(w, "`context_str` = ?")
		args = append(args, q.Context.Str)
	}
	if q.Context.Time != t {
		w = append(w, "`context_time` = ?")
		args = append(args, q.Context.Time.Unix())
	}
	if q.Context.Fetched != t {
		w = append(w, "`fetched_time` = ?")
		args = append(args, q.Context.Fetched.Unix())
	}
	sql += " WHERE " + strings.Join(w, " AND ")
	sql += sortClause

	// query
	stmt, err := q.Cache.DB.Prepare(sql)
	if err != nil {
		return
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()

	// retrieve result
	resps = make([]Response, 0)
	for rows.Next() {

		resp := Response{}
		var ctime, ftime int64

		err = rows.Scan(
			&resp.URL,
			&resp.ContextStr,
			&ctime,
			&ftime,
			&resp.Status,
			&resp.StatusCode,
			&resp.Proto,
			&resp.ContentLength,
			&resp.TransferEncodingJson,
			&resp.HeaderJson,
			&resp.TrailerJson,
			&resp.RequestJson,
			&resp.Body)

		resp.ContextTime = time.Unix(ctime, 0)
		resp.FetchedTime = time.Unix(ftime, 0)
		resps = append(resps, resp)

		if err != nil {
			return
		}
	}
	err = rows.Err()
	return
}
