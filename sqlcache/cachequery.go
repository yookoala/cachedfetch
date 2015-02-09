package sqlcache

import (
	"fmt"
	"github.com/yookoala/crawler"
	"strings"
	"time"
)

type CacheQuery struct {
	URL     string
	Context crawler.Context
	Cache   *Cache
	Order   []int
	L       int // limit
}

func (q *CacheQuery) In(Str string) crawler.CacheQuery {
	q.Context.Str = Str
	return q
}

func (q *CacheQuery) At(t time.Time) crawler.CacheQuery {
	q.Context.Time = t
	return q
}

func (q *CacheQuery) FetchedAt(t time.Time) crawler.CacheQuery {
	q.Context.Fetched = t
	return q
}

func (q *CacheQuery) SortBy(crits ...int) crawler.CacheQuery {
	for _, crit := range crits {
		q.Order = append(q.Order, crit)
	}
	return q
}

func (q *CacheQuery) Limit(l int) crawler.CacheQuery {
	q.L = l
	return q
}

// generate SQL where clause based on query parameters
func (q *CacheQuery) sqlWhere() (c string, args []interface{}) {
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
func (q *CacheQuery) sqlOrder() (c string) {
	if len(q.Order) == 0 {
		// default sort
		q.Order = []int{
			crawler.OrderContextTimeDesc,
			crawler.OrderFetchedTimeDesc,
		}
	}
	o := make([]string, 0, cap(q.Order))
	for _, field := range q.Order {
		var sqlf string
		switch field {
		case crawler.OrderContextTime:
			sqlf = "context_time"
		case crawler.OrderContextTimeDesc:
			sqlf = "context_time DESC"
		case crawler.OrderFetchedTime:
			sqlf = "fetched_time"
		case crawler.OrderFetchedTimeDesc:
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
func (q *CacheQuery) sqlLimit() (c string) {
	if q.L > 0 {
		c = fmt.Sprintf("LIMIT %d", q.L)
	}
	return
}

// generate final SQL
func (q *CacheQuery) genGetSql() (sql string, args []interface{}) {

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

func (q *CacheQuery) GetAll() (resps crawler.ResponseColl, err error) {

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
	rs := make([]crawler.Response, 0)
	for rows.Next() {

		resp := crawler.Response{}
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
	resps = &ResponseColl{
		col: rs,
	}
	err = rows.Err()
	return
}
