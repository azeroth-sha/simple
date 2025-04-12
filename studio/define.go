package studio

import (
	"github.com/azeroth-sha/simple"
	"time"
)

// Handler 表示一个事件处理器，用于处理事件。
// 参数 e 是待处理的事件对象。
type Handler func(e Event)

// Event 表示一个事件，包含事件发生的时间、事件名称、事件参数和附加参数。
type Event interface {
	// Occurred 返回事件发生的时间。
	// 返回值为 time.Time 类型，表示事件发生的时间戳。
	Occurred() time.Time

	// Name 返回事件的名称。
	// 返回值为 string 类型，用于标识事件的类型或来源。
	Name() string

	// Param 返回事件的参数。
	// 返回值为 any 类型，表示事件的主要参数，可以是任意类型。
	Param() any

	// Additive 返回事件的附加参数。
	// 返回值为 simple.Map 类型，表示事件的附加参数，通常用于传递额外的上下文信息。
	Additive() simple.Map
}

// Studio 表示一个工作室，用于管理和处理事件。
type Studio interface {
	// Release 释放工作室资源，停止所有工作线程。
	// 调用后，工作室将不再接受新任务，并等待所有任务处理完成。
	Release()

	// Task 发布一个任务到工作室。
	// 参数 e 是待处理的事件对象。
	// 参数 block 指定是否阻塞等待任务入队。
	// 参数 exp 是可选参数，指定阻塞等待的超时时间。
	// 返回值为 error 类型，表示任务发布的结果（成功或失败）。
	Task(e Event, block bool, exp ...time.Duration) error

	// SetWorkstation 设置指定名称的工作站处理器。
	// 参数 n 是工作站的名称。
	// 参数 h 是事件处理器。
	// 如果工作站已存在，则替换原有处理器。
	SetWorkstation(n string, h Handler)

	// AddWorkstation 添加指定名称的工作站处理器。
	// 参数 n 是工作站的名称。
	// 参数 h 是事件处理器。
	// 如果工作站已存在，则追加处理器到现有列表中。
	AddWorkstation(n string, h Handler)

	// Recycle 设置全局回收处理器，用于处理未匹配的事件。
	// 参数 h 是事件处理器。
	Recycle(h Handler)

	// Count 返回当前工作室中待处理的任务数量。
	// 返回值为 int 类型，表示任务队列中的事件数量。
	Count() int
}
