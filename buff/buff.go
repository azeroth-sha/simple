package buff

import (
	"bytes"
	"sync"
)

var (
	buffPool = &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
)

// Buff is an interface for a bytes.Buffer pool
type Buff interface {
	// GetBuff returns a new bytes.Buffer
	GetBuff() *bytes.Buffer
	// PutBuff puts a bytes.Buffer back into the pool
	PutBuff(buff *bytes.Buffer)
}

// buff is a wrapper around a bytes.Buffer
type buff struct {
	pool *sync.Pool
}

// GetBuff returns a new bytes.Buffer
func (b *buff) GetBuff() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

// PutBuff puts a bytes.Buffer back into the pool
func (b *buff) PutBuff(buff *bytes.Buffer) {
	buff.Reset()
	b.pool.Put(buff)
}

/*
  Package method
*/
