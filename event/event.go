package event

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 定义事件处理相关的错误
var (
	ErrClosed = errors.New(`event closed`)        // 事件处理器已关闭
	ErrFull   = errors.New(`event queue is full`) // 事件队列已满
)

// Handler 定义了事件处理函数的类型
type Handler func(t *Task)

// Handles 是事件处理函数的切片类型
type Handles []Handler

// Task 定义了事件任务的结构体，包含事件名称、触发时间、事件属性和事件数据
type Task struct {
	Name string         // 事件名称
	Time time.Time      // 触发时间
	Attr map[string]any // 事件属性
	Data interface{}    // 事件数据
}

// Event 定义了事件处理器的结构体，包含事件队列、关闭信号、阻塞模式、队列大小、任务数量、任务处理函数等
type Event struct {
	ch             chan *Task         // 事件队列
	closed         chan struct{}      // 关闭信号
	block          bool               // 是否阻塞模式
	size           int32              // 队列大小
	quantity       int32              // 当前队列中的任务数量
	group          *sync.Pool         // 用于管理任务组的 sync.Pool
	defaultHandler Handler            // 默认事件处理函数
	taskHandler    map[string]Handles // 按事件名称分类的处理函数映射
	panicHandler   func()             // 异常处理函数
}

// Close 关闭事件处理器，释放资源
func (e *Event) Close() {
	if e.closed != nil {
		close(e.closed)
		e.closed = nil
	}
}

// Add 为指定事件名称添加处理函数
func (e *Event) Add(name string, h Handler) {
	e.taskHandler[name] = append(e.taskHandler[name], h)
}

// SetDefaultHandler 设置默认事件处理函数
func (e *Event) SetDefaultHandler(h Handler) {
	e.defaultHandler = h
}

// SetPanicHandler 设置异常处理函数
func (e *Event) SetPanicHandler(h func()) {
	e.panicHandler = h
}

// Push 将事件任务推入事件队列，返回可能的错误
func (e *Event) Push(name string, data interface{}, attr map[string]any) error {
	if e.closed == nil {
		return ErrClosed
	}
	if index := atomic.AddInt32(&e.quantity, 1); e.block || index <= e.size {
		e.ch <- &Task{
			Name: name,
			Time: time.Now(),
			Attr: attr,
			Data: data,
		}
		return nil
	} else {
		atomic.AddInt32(&e.quantity, -1)
		return ErrFull
	}
}

// New 创建并初始化一个事件处理器
func New(size int32, block bool) *Event {
	e := &Event{
		ch:             make(chan *Task, size),
		closed:         make(chan struct{}),
		block:          block,
		size:           size,
		quantity:       0,
		group:          &sync.Pool{New: newGroup},
		defaultHandler: nil,
		taskHandler:    make(map[string]Handles),
		panicHandler:   nil,
	}
	go e.start()
	runtime.SetFinalizer(e, func(e *Event) {
		e.Close()
	})
	return e
}

/*
  Package method
*/

// start 启动事件处理器的监听循环
func (e *Event) start() {
EXIT:
	for {
		select {
		case t, ok := <-e.ch:
			if t != nil {
				e.handle(t)
			} else if !ok {
				break EXIT
			}
		case <-e.closed:
			break EXIT
		}
	}
}

// handle 处理单个事件任务
func (e *Event) handle(t *Task) {
	wait := e.group.Get().(*sync.WaitGroup)
	defer e.group.Put(wait)
	defer atomic.AddInt32(&e.quantity, -1)
	if jobs, ok := e.taskHandler[t.Name]; ok && len(jobs) > 0 {
		wait.Add(len(jobs))
		for _, job := range jobs {
			go e.job(wait, job, t)
		}
		wait.Wait()
	} else if e.defaultHandler != nil {
		wait.Add(1)
		go e.job(wait, e.defaultHandler, t)
		wait.Wait()
	}
}

// job 执行单个事件处理函数
func (e *Event) job(w *sync.WaitGroup, h Handler, t *Task) {
	defer w.Done()
	if e.panicHandler != nil {
		defer e.panicHandler()
	}
	h(t)
}

// newGroup 创建一个新的 sync.WaitGroup 对象
func newGroup() interface{} {
	return new(sync.WaitGroup)
}
