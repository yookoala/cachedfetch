package crawler

import (
	"database/sql"
	"sync"
)

// global lockers map specific to a database
var lockers map[*sql.DB]sync.Locker

// initialize global variables
func init() {
	lockers = make(map[*sql.DB]sync.Locker)
}
