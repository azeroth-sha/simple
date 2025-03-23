package buff

import (
	"bytes"
	"sync"
)

var defaultBuff = NewBuff(nil)

// NewBuff returns a new Buff
func NewBuff(p *sync.Pool) Buff {
	if p == nil {
		p = buffPool
	}
	return &buff{pool: p}
}

// GetBuff returns a new bytes.Buffer
func GetBuff() *bytes.Buffer {
	return defaultBuff.GetBuff()
}

// PutBuff puts a bytes.Buffer back into the pool
func PutBuff(bf *bytes.Buffer) {
	defaultBuff.PutBuff(bf)
}
