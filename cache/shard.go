package cache

import (
	"sync"
	"time"
)

type shard struct {
	l *sync.Mutex
	m map[string]*item
	h ExpiredHandler
}

func (s *shard) Has(k string) bool {
	s.l.Lock()
	defer s.l.Unlock()
	return s.has(k)
}

func (s *shard) Set(k string, v *item) bool {
	s.l.Lock()
	defer s.l.Unlock()
	s.delExpired(k)
	s.set(k, v)
	return true
}

func (s *shard) SetX(k string, v *item) bool {
	s.l.Lock()
	defer s.l.Unlock()
	s.delExpired(k)
	if !s.has(k) {
		s.set(k, v)
		return false
	}
	return true
}

func (s *shard) Get(k string) (interface{}, bool) {
	s.l.Lock()
	defer s.l.Unlock()
	if s.delExpired(k) {
		return nil, false
	} else if i, f := s.get(k); f {
		return i.v, true
	}
	return nil, false
}

func (s *shard) Del(k string) bool {
	s.l.Lock()
	defer s.l.Unlock()
	if s.delExpired(k) {
		return false
	}
	return s.del(k)
}

func (s *shard) DelExpired(k string) bool {
	s.l.Lock()
	defer s.l.Unlock()
	return s.delExpired(k)
}

func (s *shard) GetSet(k string, v *item) (interface{}, bool) {
	s.l.Lock()
	defer s.l.Unlock()
	defer s.set(k, v)
	s.delExpired(k)
	if i, f := s.get(k); f {
		return i.v, true
	}
	return nil, false
}

func (s *shard) GetSetX(k string, v *item) (interface{}, bool) {
	s.l.Lock()
	defer s.l.Unlock()
	s.delExpired(k)
	if i, f := s.get(k); f {
		return i.v, true
	}
	s.set(k, v)
	return nil, false
}

func (s *shard) GetDel(k string) (interface{}, bool) {
	s.l.Lock()
	defer s.l.Unlock()
	defer s.del(k)
	if s.delExpired(k) {
		return nil, false
	} else if i, f := s.get(k); f {
		return i.v, true
	}
	return nil, false
}

func (s *shard) All() map[string]interface{} {
	s.l.Lock()
	defer s.l.Unlock()
	return s.all()
}

func (s *shard) Len() int {
	s.l.Lock()
	defer s.l.Unlock()
	return s.len()
}

func (s *shard) TTL(k string) time.Duration {
	s.l.Lock()
	defer s.l.Unlock()
	if s.delExpired(k) {
		return -1
	}
	return s.ttl(k)
}
func (s *shard) Expire(k string, t time.Time) bool {
	s.l.Lock()
	defer s.l.Unlock()
	if s.delExpired(k) {
		return false
	}
	return s.expire(k, t)
}

func (s *shard) Handle(k string, h ExpiredHandler) bool {
	s.l.Lock()
	defer s.l.Unlock()
	if s.delExpired(k) {
		return false
	}
	return s.handle(k, h)
}

func (s *shard) CheckAll() {
	s.l.Lock()
	defer s.l.Unlock()
	s.checkAll()
}

func newShard(h ExpiredHandler) *shard {
	return &shard{
		l: new(sync.Mutex),
		m: make(map[string]*item),
		h: h,
	}
}

/*
  内部方法
*/

func (s *shard) has(k string) bool {
	return s.m[k] != nil
}

func (s *shard) set(k string, i *item) {
	s.m[k] = i
}

func (s *shard) get(k string) (*item, bool) {
	i, f := s.m[k]
	return i, f
}

func (s *shard) del(k string) bool {
	if _, f := s.m[k]; f {
		delete(s.m, k)
		return true
	}
	return false
}

func (s *shard) delExpired(k string) bool {
	if i, f := s.m[k]; f && i.Expired() {
		s.expired(k, i)
		return true
	}
	return false
}

func (s *shard) all() map[string]interface{} {
	m := make(map[string]interface{}, s.len())
	for k, v := range s.m {
		if v.Expired() {
			continue
		}
		m[k] = v.v
	}
	return m
}

func (s *shard) len() int {
	var cnt int
	for _, v := range s.m {
		if v.Expired() {
			continue
		}
		cnt++
	}
	return cnt
}

func (s *shard) ttl(k string) time.Duration {
	var dur time.Duration = -1
	if i, f := s.m[k]; f {
		if i.e.IsZero() {
			dur = 0
		} else {
			dur = time.Until(i.e)
		}
	}
	return dur
}

func (s *shard) expire(k string, t time.Time) bool {
	if i, f := s.m[k]; f {
		i.SetExpireAt(t)
		return true
	}
	return false
}

func (s *shard) handle(k string, h ExpiredHandler) bool {
	if i, f := s.m[k]; f {
		i.SetExpiredHandler(h)
		return true
	}
	return false
}

func (s *shard) checkAll() {
	for k, v := range s.m {
		if v.Expired() {
			s.expired(k, v)
		}
	}
}

func (s *shard) expired(k string, v *item) {
	delete(s.m, k)
	if v.h != nil {
		go func() { v.h(k, v.v) }()
	} else if s.h != nil {
		go func() { s.h(k, v.v) }()
	}
}
