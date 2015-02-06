package cachedfetcher

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	SQL_MYSQL = iota
	SQL_SQLITE3
	SQL_PSQL
)

func NewSqlCache(driver string, db *sql.DB) *SqlCache {

	// determine SQL type
	// and sync type
	t := SQL_MYSQL
	lt := LOCKER_ASYNC

	switch driver {
	case "postgres":
		t = SQL_PSQL
	case "sqlite3":
		t = SQL_SQLITE3
		lt = LOCKER_SYNC
	}

	// add locker to global lockers map
	lockers[db] = NewLocker(lt)

	// create the cache struct
	return &SqlCache{
		DB:   db,
		Type: t,
	}
}

type SqlCache struct {
	DB   *sql.DB
	Type int
}

func (c *SqlCache) Sql(s string) string {
	if c.Type != SQL_PSQL {
		return s
	}
	so := ""
	pos := 0
	for _, ch := range s {
		if ch != '?' {
			so += string(ch)
		} else {
			pos++
			so += fmt.Sprintf("$%d", pos)
		}
	}
	return so
}

func (c *SqlCache) Prepare(s string) (stmt *sql.Stmt, err error) {
	stmt, err = c.DB.Prepare(c.Sql(s))
	return
}

func (c *SqlCache) Lock() {
	lockers[c.DB].Lock()
}

func (c *SqlCache) Unlock() {
	lockers[c.DB].Unlock()
}

func (c *SqlCache) Add(url string, ctx Context, r *Response) (err error) {

	// sync database call sequence, if necessary
	c.Lock()
	defer c.Unlock()

	// prepare and execute the insert call
	stmt, err := c.Prepare("INSERT INTO cachedfetch_cache " +
		"(url, context_str, context_time, fetched_time, " +
		"status, status_code, proto, content_length, " +
		"transfer_encoding, header, trailer, request, body)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()

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
	L       int // limit
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

func (q *SqlCacheQuery) Limit(l int) CacheQuery {
	q.L = l
	return q
}

// generate SQL where clause based on query parameters
func (q *SqlCacheQuery) sqlWhere() (c string, args []interface{}) {
	w := make([]string, 0, 4)
	args = make([]interface{}, 0, 4)
	var t time.Time // empty time for reference

	if q.URL != "" {
		w = append(w, "url = ?")
		args = append(args, q.URL)
	}
	if q.Context.Str != "" {
		w = append(w, "context_str = ?")
		args = append(args, q.Context.Str)
	}
	if q.Context.Time != t {
		w = append(w, "context_time = ?")
		args = append(args, q.Context.Time.Unix())
	}
	if q.Context.Fetched != t {
		w = append(w, "fetched_time = ?")
		args = append(args, q.Context.Fetched.Unix())
	}
	c = "WHERE " + strings.Join(w, " AND ")

	return
}

// generate SQL order clause based on query parameters
func (q *SqlCacheQuery) sqlOrder() (c string) {
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
			sqlf = "context_time"
		case OrderContextTimeDesc:
			sqlf = "context_time DESC"
		case OrderFetchedTime:
			sqlf = "fetched_time"
		case OrderFetchedTimeDesc:
			sqlf = "fetched_time DESC"
		}
		if sqlf != "" {
			o = append(o, sqlf)
		}
	}
	c = "ORDER BY " + strings.Join(o, ", ")
	return
}

// generate SQL limit clause
func (q *SqlCacheQuery) sqlLimit() (c string) {
	if q.L > 0 {
		c = fmt.Sprintf("LIMIT %d", q.L)
	}
	return
}

// generate final SQL
func (q *SqlCacheQuery) genGetSql() (sql string, args []interface{}) {

	// select clause
	sql = "SELECT url, context_str, context_time, " +
		"fetched_time, status, status_code, proto, " +
		"content_length, transfer_encoding, " +
		"header, trailer, request, " +
		"body " +
		"FROM cachedfetch_cache "

	// build other clauses
	var where, order, limit string
	where, args = q.sqlWhere()
	order = q.sqlOrder()
	limit = q.sqlLimit()
	sql += " " + where + " " + order + " " + limit
	return
}

func (q *SqlCacheQuery) GetAll() (resps ResponseColl, err error) {

	// sync database call sequence, if necessary
	q.Cache.Lock()
	defer q.Cache.Unlock()

	// query
	sql, args := q.genGetSql()
	stmt, err := q.Cache.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()

	// retrieve result
	rs := make([]Response, 0)
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
		rs = append(rs, resp)

		if err != nil {
			return
		}
	}
	resps = &SqlResponseColl{
		col: rs,
	}
	err = rows.Err()
	return
}

type SqlResponseColl struct {
	col []Response
	cur int
}

func (rc *SqlResponseColl) Next() bool {
	rc.cur++
	if rc.cur <= len(rc.col) {
		return true
	}
	return false
}

func (rc *SqlResponseColl) Get() (resp *Response, err error) {
	if rc.cur <= len(rc.col) {
		resp = &rc.col[rc.cur-1]
	} else {
		err = fmt.Errorf("Getting item out of range")
	}
	return
}

func (rc *SqlResponseColl) Close() (err error) {
	// place holder
	return
}
