package cachedfetcher

import (
	"sync"
)

const (
	MUTEX_SYNC = iota
	MUTEX_ASYNC
)

func NewMutex(t int) sync.Locker {
	if t == MUTEX_SYNC {
		return &sync.Mutex{}
	}
	return &fakeMutex{}
}

// a fake mutex that doesn't really lock
// works as async mutex solution
type fakeMutex struct {
}

func (m *fakeMutex) Lock() {
}

func (m *fakeMutex) Unlock() {
}
