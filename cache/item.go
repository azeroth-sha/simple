package cache

import (
	"time"
)

type ExpiredHandler func(k string, v interface{})

type item struct {
	v interface{}
	e time.Time
	h ExpiredHandler
}

func (i *item) Expired() bool {
	if i.e.IsZero() {
		return false
	}
	return time.Now().After(i.e)
}

func (i *item) SetExpireAt(t time.Time) {
	i.e = t
}

func (i *item) SetExpiredHandler(h ExpiredHandler) {
	i.h = h
}

func newItem(v interface{}, opts ...ItemOption) *item {
	i := &item{v: v}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
