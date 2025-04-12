package studio

// Option 是一个函数类型，用于配置 engine 实例
type Option func(*engine)

// WithPipeSize 设置事件队列管道的容量
// 参数 n 必须大于 0，否则配置无效
func WithPipeSize(n int) Option {
	return func(e *engine) {
		if n > 0 {
			e.pipeSize = n
		}
	}
}

// WithJobSize 设置工作线程的数量
// 参数 n 必须大于 0，否则配置无效
func WithJobSize(n int) Option {
	return func(e *engine) {
		if n > 0 {
			e.jobSize = n
		}
	}
}

// WithRecycler 设置全局回收处理器
// 该处理器用于处理未注册事件
func WithRecycler(h Handler) Option {
	return func(e *engine) {
		e.recycler = h
	}
}

// WithPanicFunc 设置恐慌恢复函数
// 该函数在工作线程发生恐慌时被调用
func WithPanicFunc(h func()) Option {
	return func(e *engine) {
		e.panicFunc = h
	}
}
