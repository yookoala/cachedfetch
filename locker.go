package crawler

import (
	"sync"
)

const (
	LOCKER_SYNC = iota
	LOCKER_ASYNC
)

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
