package cachedfetcher

import (
	"testing"
	"time"
)

func TestSqlCache(t *testing.T) {
	var c Cache = &SqlCache{}
	t.Log("SqlCache implements Cache: %#v", c)
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
	sqlE := "WHERE `url` = ? AND `context_str` = ? AND " +
		"`context_time` = ? AND `fetched_time` = ?"

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
	sqlE := "ORDER BY `context_time`, `fetched_time` DESC"
	if sql != sqlE {
		t.Errorf("SQL generated is different from expected.\n"+
			"Expect: \"%s\"\n"+
			"Get:    \"%s\"",
			sql, sqlE)
	}
}
