package sqlcache

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/crawler"
	"time"
)

const (
	SQL_MYSQL = iota
	SQL_SQLITE3
	SQL_PSQL
)

func New(driver string, db *sql.DB) *Cache {

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
	return &Cache{
		DB:   db,
		Type: t,
	}
}

type Cache struct {
	DB   *sql.DB
	Type int
}

func (c *Cache) Sql(s string) string {
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

func (c *Cache) Prepare(s string) (stmt *sql.Stmt, err error) {
	stmt, err = c.DB.Prepare(c.Sql(s))
	return
}

func (c *Cache) Lock() {
	lockers[c.DB].Lock()
}

func (c *Cache) Unlock() {
	lockers[c.DB].Unlock()
}

func (c *Cache) Add(url string, ctx crawler.Context, r *crawler.Response) (err error) {

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

func (c *Cache) Find(url string) crawler.CacheQuery {
	return &CacheQuery{
		URL:   url,
		Cache: c,
		Order: make([]int, 0),
	}
}

// find with context string
func (c *Cache) FindIn(str string) crawler.CacheQuery {
	ctx := crawler.Context{
		Str: str,
	}
	return &CacheQuery{
		Context: ctx,
		Cache:   c,
		Order:   make([]int, 0),
	}
}

// find with context time
func (c *Cache) FindAt(t time.Time) crawler.CacheQuery {
	ctx := crawler.Context{
		Time: t,
	}
	return &CacheQuery{
		Context: ctx,
		Cache:   c,
		Order:   make([]int, 0),
	}
}
