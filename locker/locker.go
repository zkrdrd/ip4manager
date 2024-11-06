package locker

import (
	"sync"
)

type Locker struct {
	workMX *sync.Mutex
	value  map[string]*sync.Mutex
}

func NewLocker() *Locker {
	return &Locker{
		workMX: &sync.Mutex{},
		value:  make(map[string]*sync.Mutex),
	}
}

func (l *Locker) Lock(address string) {
	l.workMX.Lock()

	val, ok := l.value[address]
	if !ok {
		val = &sync.Mutex{}
		l.value[address] = val
	}

	l.workMX.Unlock()

	val.Lock()
}

func (l *Locker) Unlock(address string) {
	l.workMX.Lock()

	val, ok := l.value[address]
	if !ok {
		return
	}

	l.workMX.Unlock()

	val.Unlock()
}
