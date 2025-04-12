package studio

import (
	"github.com/azeroth-sha/simple"
	"time"
)

// message 是 Event 接口的具体实现，表示一个事件消息。
type message struct {
	t time.Time  // 事件发生的时间
	n string     // 事件的名称
	d any        // 事件的主要参数
	a simple.Map // 事件的附加参数
}

// Occurred 返回事件发生的时间。
// 返回值为 time.Time 类型，表示事件的时间戳。
func (m *message) Occurred() time.Time {
	return m.t
}

// Name 返回事件的名称。
// 返回值为 string 类型，用于标识事件的类型或来源。
func (m *message) Name() string {
	return m.n
}

// Param 返回事件的主要参数。
// 返回值为 any 类型，表示事件的核心数据。
func (m *message) Param() any {
	return m.d
}

// Additive 返回事件的附加参数。
// 返回值为 simple.Map 类型，表示事件的额外上下文信息。
func (m *message) Additive() simple.Map {
	return m.a
}

// NewEvent 创建一个新的事件实例。
// 参数 n 是事件的名称。
// 参数 d 是事件的主要参数。
// 参数 a 是事件的附加参数，类型为 map[string]any。
// 参数 t 是可选参数，表示事件发生的时间；如果未提供，则使用当前时间。
// 返回值为 Event 接口类型，表示创建的事件对象。
func NewEvent(n string, d any, a simple.Map, t ...time.Time) Event {
	msg := &message{
		t: time.Now(), // 默认使用当前时间
		n: n,
		d: d,
		a: a,
	}
	if len(t) != 0 { // 如果提供了时间参数，则覆盖默认值
		msg.t = t[0]
	}
	return msg
}
