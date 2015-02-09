package sqlcache

import (
	"database/sql"
	"sync"
)

const (
	LOCKER_SYNC = iota
	LOCKER_ASYNC
)

// global lockers map specific to a database
var lockers map[*sql.DB]sync.Locker

// initialize global variables
func init() {
	lockers = make(map[*sql.DB]sync.Locker)
}

func NewLocker(t int) sync.Locker {
	if t == LOCKER_SYNC {
		return &sync.Mutex{}
	}
	return &nonLocker{}
}

// a fake locker that doesn't really block
// works as async solution
type nonLocker struct {
}

func (m *nonLocker) Lock() {
}

func (m *nonLocker) Unlock() {
}
