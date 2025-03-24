package cache

import (
	"github.com/azeroth-sha/simple/internal"
	"hash/fnv"
	"runtime"
	"time"
)

// Dict 字典
type Dict struct {
	bucket      []*shard
	closed      chan struct{}
	shardNum    uint32
	expHandler  ExpiredHandler
	chkInterval time.Duration
}

// Has 判断key是否存在
func (d *Dict) Has(k string) (has bool) {
	sh := d.shard(k)
	return sh.Has(k)
}

// Set 设置key-value
func (d *Dict) Set(k string, v interface{}, opts ...ItemOption) (ok bool) {
	sh := d.shard(k)
	return sh.Set(k, newItem(v, opts...))
}

// SetX 设置key-value，如果key不存在则设置新值
func (d *Dict) SetX(k string, v interface{}, opts ...ItemOption) (has bool) {
	sh := d.shard(k)
	return sh.SetX(k, newItem(v, opts...))
}

// Get 获取key对应的value
func (d *Dict) Get(k string) (val interface{}, has bool) {
	sh := d.shard(k)
	return sh.Get(k)
}

// Del 删除key
func (d *Dict) Del(k string) (ok bool) {
	sh := d.shard(k)
	return sh.Del(k)
}

// DelExpired 删除key对应的过期数据
func (d *Dict) DelExpired(k string) (ok bool) {
	sh := d.shard(k)
	return sh.DelExpired(k)
}

// GetSet 获取key对应的value并设置新值
func (d *Dict) GetSet(k string, v interface{}, opts ...ItemOption) (val interface{}, has bool) {
	sh := d.shard(k)
	return sh.GetSet(k, newItem(v, opts...))
}

// GetSetX 获取key对应的value并设置新值，如果key不存在则设置新值
func (d *Dict) GetSetX(k string, v interface{}, opts ...ItemOption) (val interface{}, has bool) {
	sh := d.shard(k)
	return sh.GetSetX(k, newItem(v, opts...))
}

// GetDel 获取key对应的value并删除
func (d *Dict) GetDel(k string) (val interface{}, has bool) {
	sh := d.shard(k)
	return sh.GetDel(k)
}

// All 获取所有key-value
func (d *Dict) All() (all map[string]interface{}) {
	all = make(map[string]interface{})
	for i := uint32(0); i < d.shardNum; i++ {
		shardAll := d.bucket[i].All()
		for k, v := range shardAll {
			all[k] = v
		}
	}
	return all
}

// Len 获取key的数量
func (d *Dict) Len() (cnt int) {
	for i := uint32(0); i < d.shardNum; i++ {
		cnt += d.bucket[i].Len()
	}
	return cnt
}

// TTL 获取key的剩余时间
func (d *Dict) TTL(k string) (dur time.Duration) {
	sh := d.shard(k)
	return sh.TTL(k)
}

// ExpireDur 设置key的过期时长
func (d *Dict) ExpireDur(k string, dur time.Duration) (has bool) {
	sh := d.shard(k)
	t := time.Now().Add(dur)
	return sh.Expire(k, t)
}

// ExpireAt 设置key的过期时间
func (d *Dict) ExpireAt(k string, t time.Time) (has bool) {
	sh := d.shard(k)
	return sh.Expire(k, t)
}

// Handle 设置key的过期处理器
func (d *Dict) Handle(k string, h ExpiredHandler) (has bool) {
	sh := d.shard(k)
	return sh.Handle(k, h)
}

// CheckAll 检测所有key的过期
func (d *Dict) CheckAll() {
	for i := uint32(0); i < d.shardNum; i++ {
		d.bucket[i].CheckAll()
	}
}

// New 创建一个字典
func New(opts ...DictOption) *Dict {
	d := &Dict{
		bucket:      nil,
		closed:      make(chan struct{}),
		shardNum:    uint32(runtime.NumCPU() * 8),
		expHandler:  nil,
		chkInterval: time.Second,
	}
	for _, opt := range opts {
		opt(d)
	}
	for i := uint32(0); i < d.shardNum; i++ {
		d.bucket = append(d.bucket, newShard(d.expHandler))
	}
	go func() {
		if d.chkInterval <= 0 {
			return
		}
		tk := time.NewTicker(d.chkInterval)
		defer tk.Stop()
	EXIT:
		for {
			select {
			case <-d.closed:
				break EXIT
			case <-tk.C:
				d.CheckAll()
			}
		}
	}()
	runtime.SetFinalizer(d, func(d *Dict) {
		close(d.closed)
	})
	return d
}

/*
  内部方法
*/

func (d *Dict) shard(k string) *shard {
	return d.bucket[d.sum(k)%d.shardNum]
}

func (d *Dict) sum(k string) uint32 {
	h := fnv.New32()
	_, _ = h.Write(internal.ToBytes(k))
	return h.Sum32()
}
