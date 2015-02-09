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

const (
	SQL_INIT_MYSQL = `
		CREATE TABLE IF NOT EXISTS cachedfetch_cache (
			--
			-- context and fetch information
			--
			url               VARCHAR(255) DEFAULT '',
			context_str       VARCHAR(255) DEFAULT '',
			context_time      INT(11) DEFAULT 0,
			fetched_time      INT(11) DEFAULT 0,

			--
			-- response meta information
			--
			status            TEXT DEFAULT '',
			status_code       INT(5) DEFAULT 200,
			proto             TEXT DEFAULT '',
			content_length    INT(11) DEFAULT 0,
			transfer_encoding TEXT DEFAULT '',
			header            TEXT DEFAULT '',
			trailer           TEXT DEFAULT '',
			request           TEXT DEFAULT '',
			tls               TEXT DEFAULT '',

			--
			-- response body
			--
			body              MEDIUMBLOB,

			PRIMARY KEY (url, context_str, context_time)
		) engine=InnoDB CHARACTER SET utf8;

		CREATE INDEX url ON cachedfetch_cache(url);
		CREATE INDEX context ON cachedfetch_cache(context_str, context_time);
		CREATE INDEX context_str  ON cachedfetch_cache(context_str);
		CREATE INDEX context_time ON cachedfetch_cache(context_time);
	`

	SQL_INIT_SQLITE3 = `
		CREATE TABLE IF NOT EXISTS cachedfetch_cache (
			--
			-- context and fetch information
			--
			url               VARCHAR(255) DEFAULT '',
			context_str       VARCHAR(255) DEFAULT '',
			context_time      INT(11) DEFAULT 0,
			fetched_time      INT(11) DEFAULT 0,

			--
			-- response meta information
			--
			status            TEXT DEFAULT '',
			status_code       INT(5) DEFAULT 200,
			proto             TEXT DEFAULT '',
			content_length    INT(11) DEFAULT 0,
			transfer_encoding TEXT DEFAULT '',
			header            TEXT DEFAULT '',
			trailer           TEXT DEFAULT '',
			request           TEXT DEFAULT '',
			tls               TEXT DEFAULT '',

			--
			-- response body
			--
			body              MEDIUMBLOB,

			PRIMARY KEY (url, context_str, context_time)
		);

		CREATE INDEX url ON cachedfetch_cache(url);
		CREATE INDEX context ON cachedfetch_cache(context_str, context_time);
		CREATE INDEX context_str  ON cachedfetch_cache(context_str);
		CREATE INDEX context_time ON cachedfetch_cache(context_time);
	`

	SQL_INIT_PSQL = `
		CREATE TABLE IF NOT EXISTS cachedfetch_cache (
			--
			-- context and fetch information
			--
			url               VARCHAR(255) DEFAULT '',
			context_str       VARCHAR(255) DEFAULT '',
			context_time      INTEGER DEFAULT 0,
			fetched_time      INTEGER DEFAULT 0,

			--
			-- response meta information
			--
			status            TEXT DEFAULT '',
			status_code       SMALLINT DEFAULT 200,
			proto             TEXT DEFAULT '',
			content_length    INTEGER DEFAULT 0,
			transfer_encoding TEXT DEFAULT '',
			header            TEXT DEFAULT '',
			trailer           TEXT DEFAULT '',
			request           TEXT DEFAULT '',
			tls               TEXT DEFAULT '',

			--
			-- response body
			--
			body              BYTEA,

			PRIMARY KEY (url, context_str, context_time)
		);

		--
		-- Add extra index
		--
		CREATE INDEX url ON cachedfetch_cache(url);
		CREATE INDEX context ON cachedfetch_cache(context_str, context_time);
		CREATE INDEX context_str  ON cachedfetch_cache(context_str);
		CREATE INDEX context_time ON cachedfetch_cache(context_time);
	`
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

func (c *Cache) Init() (err error) {
	var s string
	switch c.Type {
	case SQL_MYSQL:
		s = SQL_INIT_MYSQL
	case SQL_PSQL:
		s = SQL_INIT_PSQL
	default:
		s = SQL_INIT_SQLITE3
	}
	stmt, err := c.Prepare(s)
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	return
}

func (c *Cache) Rebuild() (err error) {
	stmt, err := c.Prepare("DROP TABLE IF EXISTS cachedfetch_cache")
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		return
	}
	err = c.Init()
	return
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
