package cachedfetcher

import (
	"testing"
	"time"
)

func TestSqlCache(t *testing.T) {
	var c Cache = &SqlCache{}
	t.Log("SqlCache implements Cache: %#v", c)
}

func TestSqlCacheSqlDefault(t *testing.T) {
	c := &SqlCache{}
	if c.Type != SQL_MYSQL {
		t.Errorf("Default CacheQuery.Type should be SQL_MYSQL")
	}
}

func TestSqlCacheSqlMySQL(t *testing.T) {
	c := &SqlCache{
		Type: SQL_MYSQL,
	}
	raw := "SELECT * FROM t1 WHERE a=? AND b=? AND c=?"
	s := c.Sql(raw)
	if s != raw {
		t.Errorf("SqlCache.Sql return original SQL. Returned: %s", s)
	}
}

func TestSqlCacheSqlPSQL(t *testing.T) {
	c := &SqlCache{
		Type: SQL_PSQL,
	}
	raw := "SELECT * FROM t1 WHERE a=? AND b=? AND c=?"
	expected := "SELECT * FROM t1 WHERE a=$1 AND b=$2 AND c=$3"
	s := c.Sql(raw)
	if s != expected {
		t.Errorf("SqlCache.Sql did not format SQL into PSQL "+
			"position parameter. Returned: %s", s)
	}
}

func TestSqlCacheQuery(t *testing.T) {
	var c CacheQuery = &SqlCacheQuery{}
	t.Log("SqlCacheQuery implements CacheQuery: %#v", c)
}

func TestSqlCacheQueryWhere(t *testing.T) {
	t1 := time.Now()
	t2 := t1.AddDate(0, 0, 1) // add 1 day to t1
	Ctx := Context{
		Str:     "context 1",
		Time:    t1,
		Fetched: t2,
	}
	var q = &SqlCacheQuery{
		URL:     "test url",
		Context: Ctx,
	}
	sql, args := q.sqlWhere()
	sqlE := "WHERE url = ? AND context_str = ? AND " +
		"context_time = ? AND fetched_time = ?"

	// assert result parameters
	if sql != sqlE {
		t.Errorf("SQL generated is different from expected.\n"+
			"Expect: \"%s\"\n"+
			"Get:    \"%s\"",
			sql, sqlE)
	}
	if len(args) != 4 {
		t.Errorf("Number of SQL arguemnts is not as expected.\n"+
			"Expect: \"%d\"\n"+
			"Get:    \"%d\"\n"+
			"Args:   %#v", 4,
			len(args), args)
	} else {
		if args[0] != "test url" {
			t.Errorf("Argument 0 is not expected\n"+
				"Expect: \"%s\"\n"+
				"Get:    \"%s\"",
				"test url", args[0])
		}
		if args[1] != "context 1" {
			t.Errorf("Argument 1 is not expected\n"+
				"Expect: \"%s\"\n"+
				"Get:    \"%s\"",
				"context 1", args[1])
		}
		if args[2] != t1.Unix() {
			t.Errorf("Argument 2 is not expected\n"+
				"Expect: \"%s\"\n"+
				"Get:    \"%s\"",
				t1.Unix(), args[2])
		}
		if args[3] != t2.Unix() {
			t.Errorf("Argument 3 is not expected\n"+
				"Expect: \"%s\"\n"+
				"Get:    \"%s\"",
				t2.Unix(), args[3])
		}
	}
}

func TestSqlCacheQueryOrder(t *testing.T) {
	var q = &SqlCacheQuery{
		Order: []int{
			OrderContextTime,
			OrderFetchedTimeDesc,
		},
	}
	sql := q.sqlOrder()
	sqlE := "ORDER BY context_time, fetched_time DESC"
	if sql != sqlE {
		t.Errorf("SQL generated is different from expected.\n"+
			"Expect: \"%s\"\n"+
			"Get:    \"%s\"",
			sql, sqlE)
	}
}

func TestSqlCacheQueryLimit(t *testing.T) {
	var q = &SqlCacheQuery{
		L: 1123,
	}
	sql := q.sqlLimit()
	sqlE := "LIMIT 1123"
	if sql != sqlE {
		t.Errorf("SQL generated is different from expected.\n"+
			"Expect: \"%s\"\n"+
			"Get:    \"%s\"",
			sql, sqlE)
	}
}

func TestSqlCacheQueryLimit0(t *testing.T) {
	var q = &SqlCacheQuery{}
	sql := q.sqlLimit()
	sqlE := ""
	if sql != sqlE {
		t.Errorf("SQL generated is different from expected.\n"+
			"Expect: \"%s\"\n"+
			"Get:    \"%s\"",
			sql, sqlE)
	}
}

func TestSqlResponseColl(t *testing.T) {
	var rc ResponseColl
	rc = &SqlResponseColl{}
	t.Logf("SqlResponseColl implements ResponseColl: %#v", rc)
}

func TestSqlResponseCollRoutines(t *testing.T) {
	rc := SqlResponseColl{
		col: []Response{
			Response{
				URL:        "Response 1",
				StatusCode: 1,
			},
			Response{
				URL:        "Response 2",
				StatusCode: 2,
			},
			Response{
				URL:        "Response 3",
				StatusCode: 3,
			},
			Response{
				URL:        "Response 4",
				StatusCode: 4,
			},
			Response{
				URL:        "Response 5",
				StatusCode: 5,
			},
		},
	}

	count := 0
	for rc.Next() {
		count++
		resp, err := rc.Get()
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		if resp.StatusCode != count {
			t.Errorf("StatusCode not correct. Expecting %d but get %d", count, resp.StatusCode)
		}
	}
	if err := rc.Close(); err != nil {
		if err != nil {
			t.Errorf("Error closing ResponseCollection: %s", err.Error())
		}
	}
}
