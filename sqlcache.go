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
		Order: make([]int, 0),
	}
}

// find with context string
func (c *SqlCache) FindIn(str string) CacheQuery {
	ctx := Context{
		Str: str,
	}
	return &SqlCacheQuery{
		Context: ctx,
		Cache:   c,
		Order:   make([]int, 0),
	}
}

// find with context time
func (c *SqlCache) FindAt(t time.Time) CacheQuery {
	ctx := Context{
		Time: t,
	}
	return &SqlCacheQuery{
		Context: ctx,
		Cache:   c,
		Order:   make([]int, 0),
	}
}

type SqlCacheQuery struct {
	URL     string
	Context Context
	Cache   *SqlCache
	Order   []int
}

func (q *SqlCacheQuery) In(Str string) CacheQuery {
	q.Context.Str = Str
	return q
}

func (q *SqlCacheQuery) At(t time.Time) CacheQuery {
	q.Context.Time = t
	return q
}

func (q *SqlCacheQuery) FetchedAt(t time.Time) CacheQuery {
	q.Context.Fetched = t
	return q
}

func (q *SqlCacheQuery) SortBy(crits ...int) CacheQuery {
	for _, crit := range crits {
		q.Order = append(q.Order, crit)
	}
	return q
}

func (q *SqlCacheQuery) GetAll() (resps []Response, err error) {
	sql := "SELECT `url`, `context_str`, `context_time`, " +
		"`fetched_time`, `status`, `status_code`, `proto`, " +
		"`content_length`, `transfer_encoding`, " +
		"`header`, `trailer`, `request`, " +
		"`body` " +
		"FROM `cachedfetch_cache` "

	var t time.Time

	// parameters to build where clause
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
	whereClause := "WHERE " + strings.Join(w, " AND ")

	// parameters to build order clause
	if len(q.Order) == 0 {
		// default sort
		q.Order = []int{
			OrderContextTimeDesc,
			OrderFetchedTimeDesc,
		}
	}
	o := make([]string, 0, cap(q.Order))
	for _, field := range q.Order {
		var sqlf string
		switch field {
		case OrderContextTime:
			sqlf = "`context_time`"
		case OrderContextTimeDesc:
			sqlf = "`context_time` DESC"
		case OrderFetchedTime:
			sqlf = "`fetched_time`"
		case OrderFetchedTimeDesc:
			sqlf = "`fetched_time` DESC"
		}
		if sqlf != "" {
			o = append(o, sqlf)
		}
	}
	orderClause := "ORDER BY " + strings.Join(o, ", ")

	// query
	sql += " " + whereClause + " " + orderClause
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
