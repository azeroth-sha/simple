package lock

import (
	"github.com/azeroth-sha/simple/internal"
	"hash/fnv"
	"runtime"
	"sync"
)

type Mutex interface {
	Lock()
	Unlock()
	TryLock() bool
}

type RWMutex interface {
	Mutex
	RLock()
	RUnlock()
	TryRLock() bool
}

// NewMutexPool returns a mutex pool.
func NewMutexPool(n ...int) *MutexPool {
	cnt := runtime.NumCPU() * 4
	if len(n) > 0 && n[0] > 0 {
		cnt = n[0]
	}
	p := &MutexPool{
		pool: make([]Mutex, cnt),
		size: uint32(cnt),
	}
	for i := 0; i < cnt; i++ {
		p.pool[i] = new(sync.Mutex)
	}
	return p
}

// NewRWMutexPool returns a rw mutex pool.
func NewRWMutexPool(n ...int) *RwMutexPool {
	cnt := runtime.NumCPU() * 4
	if len(n) > 0 && n[0] > 0 {
		cnt = n[0]
	}
	p := &RwMutexPool{
		pool: make([]RWMutex, cnt),
		size: uint32(cnt),
	}
	for i := 0; i < cnt; i++ {
		p.pool[i] = new(sync.RWMutex)
	}
	return p
}

/*
	Package method
*/

type MutexPool struct {
	pool []Mutex
	size uint32
}

func (m *MutexPool) Get(s string) Mutex {
	return m.pool[sum(s)%m.size]
}

type RwMutexPool struct {
	pool []RWMutex
	size uint32
}

func (m *RwMutexPool) Get(s string) RWMutex {
	return m.pool[sum(s)%m.size]
}

func sum(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write(internal.ToBytes(s))
	return h.Sum32()
}
