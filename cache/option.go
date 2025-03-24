package cache

import "time"

// ItemOption 元素选项
type ItemOption func(*item)

// ItemExDur 设置过期时长
func ItemExDur(d time.Duration) ItemOption {
	return func(i *item) {
		i.SetExpireAt(time.Now().Add(d))
	}
}

// ItemExAt 设置过期时间
func ItemExAt(t time.Time) ItemOption {
	return func(i *item) {
		i.SetExpireAt(t)
	}
}

// ItemExHand 设置过期处理器
func ItemExHand(h ExpiredHandler) ItemOption {
	return func(i *item) {
		i.SetExpiredHandler(h)
	}
}

// DictOption 词典选项
type DictOption func(*Dict)

// DictCheckInterval 设置检查间隔
func DictCheckInterval(dur time.Duration) DictOption {
	return func(d *Dict) {
		d.chkInterval = dur
	}
}

// DictExpireHandler 设置过期处理器
func DictExpireHandler(h ExpiredHandler) DictOption {
	return func(d *Dict) {
		d.expHandler = h
	}
}

// DictShardNum 设置分片数量
func DictShardNum(n int) DictOption {
	return func(d *Dict) {
		if n <= 0 {
			return
		}
		d.shardNum = uint32(n)
	}
}
