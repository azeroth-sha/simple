package studio

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 全局错误定义
var (
	ErrReleased = errors.New(`studio released`)     // 引擎已释放时返回的错误
	ErrFull     = errors.New(`event queue is full`) // 事件队列满时返回的错误
)

// engine 实现 Studio 接口的事件处理引擎
type engine struct {
	running   int32                // 原子操作标记引擎运行状态（1运行中/0已停止）
	closed    chan struct{}        // 引擎关闭信号通道
	pipeSize  int                  // 事件队列管道容量
	pipeline  chan Event           // 事件队列管道
	jobSize   int                  // 工作线程数量
	jobWait   *sync.WaitGroup      // 等待所有工作线程退出的同步器
	recycler  Handler              // 未匹配事件处理器（回收处理器）
	jobMu     *sync.RWMutex        // 保护 jobMap 和 recycler 的读写锁
	jobMap    map[string][]Handler // 事件名称到处理器的映射表
	panicFunc func()               // 恐慌恢复函数
}

// Release 释放引擎资源，停止所有工作线程
func (eng *engine) Release() {
	if atomic.SwapInt32(&eng.running, 0) != 1 {
		return
	}
	defer close(eng.pipeline) // 确保关闭事件管道
	close(eng.closed)         // 发送关闭信号
	eng.jobWait.Wait()        // 等待所有工作线程退出
}

// Task 将事件加入处理队列，支持三种模式：
// block=true 阻塞等待入队
// block=false 且 exp 存在时：带超时的阻塞等待
// block=false 且无 exp 时：非阻塞立即返回
func (eng *engine) Task(e Event, block bool, exp ...time.Duration) error {
	if !eng.isRunning() {
		return ErrReleased
	}
	if block {
		select {
		case eng.pipeline <- e:
			return nil
		case <-eng.closed:
			return ErrReleased
		}
	}
	if len(exp) > 0 {
		select {
		case eng.pipeline <- e:
			return nil
		case <-eng.closed:
			return ErrReleased
		case <-time.After(exp[0]):
			return ErrFull
		}
	} else {
		select {
		case eng.pipeline <- e:
			return nil
		case <-eng.closed:
			return ErrReleased
		default:
			return ErrFull
		}
	}
}

// SetWorkstation 设置指定名称的工作站处理器（替换原有）
func (eng *engine) SetWorkstation(n string, h Handler) {
	if !eng.isRunning() {
		return
	}
	eng.jobMu.Lock()
	defer eng.jobMu.Unlock()
	eng.jobMap[n] = []Handler{h}
}

// AddWorkstation 添加指定名称的工作站处理器（追加处理器）
func (eng *engine) AddWorkstation(n string, h Handler) {
	if !eng.isRunning() {
		return
	}
	eng.jobMu.Lock()
	defer eng.jobMu.Unlock()
	eng.jobMap[n] = append(eng.jobMap[n], h)
}

// Recycle 设置全局回收处理器（处理未注册事件）
func (eng *engine) Recycle(h Handler) {
	if !eng.isRunning() {
		return
	}
	eng.jobMu.Lock()
	defer eng.jobMu.Unlock()
	eng.recycler = h
}

// Count 获取当前待处理事件数量
func (eng *engine) Count() int {
	if !eng.isRunning() {
		return 0
	}
	return len(eng.pipeline)
}

// New 创建新的工作室引擎实例（采用选项模式配置）
func New(opts ...Option) Studio {
	numCPU := runtime.NumCPU()
	var obj = &engine{
		running:   1,
		closed:    make(chan struct{}),
		pipeSize:  numCPU,     // 默认管道容量=CPU核心数
		pipeline:  nil,        // 事件队列管道
		jobSize:   numCPU * 2, // 默认工作线程数=2*CPU核心数
		jobWait:   new(sync.WaitGroup),
		recycler:  nil,
		jobMu:     new(sync.RWMutex),
		jobMap:    make(map[string][]Handler),
		panicFunc: nil,
	}
	for _, opt := range opts {
		opt(obj)
	}
	obj.pipeline = make(chan Event, obj.pipeSize)
	for i := 0; i < obj.jobSize; i++ {
		obj.jobWait.Add(1)
		go obj.job() // 启动工作线程池
	}
	return obj
}

// job 工作线程主循环（内部方法）
func (eng *engine) job() {
	defer eng.jobWait.Done()
	wait := new(sync.WaitGroup)
EXIT:
	for {
		select {
		case e, ok := <-eng.pipeline:
			if !ok {
				break EXIT
			} else if handles := eng.getHandles(e.Name()); len(handles) > 0 {
				cnt := len(handles)
				for i := 0; i < cnt; i++ {
					wait.Add(1)
					go eng.work(wait, handles[i], e)()
				}
				wait.Wait()
			}
		case <-eng.closed:
			break EXIT
		}
	}
}

// work 执行单个处理器任务（内部方法）
func (eng *engine) work(w *sync.WaitGroup, h Handler, e Event) func() {
	return func() {
		defer w.Done()
		if eng.panicFunc != nil {
			defer eng.panicFunc() // 恐慌恢复机制
		}
		h(e)
	}
}

// getHandles 获取事件对应的处理器列表（内部方法）
func (eng *engine) getHandles(n string) []Handler {
	eng.jobMu.RLock()
	defer eng.jobMu.RUnlock()
	handles := make([]Handler, 0, 4)
	if hs, ok := eng.jobMap[n]; ok {
		handles = append(handles, hs...)
	}
	if len(handles) == 0 && eng.recycler != nil {
		handles = append(handles, eng.recycler)
	}
	return handles
}

// isRunning 检查引擎是否运行中（内部方法）
func (eng *engine) isRunning() bool {
	return atomic.LoadInt32(&eng.running) == 1
}
