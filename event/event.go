package event

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed = errors.New(`event closed`)
	ErrFull   = errors.New(`event queue is full`)
)

type (
	Handler func(t *Task)
	Handles []Handler
	Task    struct {
		Name string         // 事件名称
		Time time.Time      // 触发时间
		Attr map[string]any // 事件属性
		Data interface{}    // 事件数据
	}
)

type Event struct {
	ch             chan *Task
	closed         chan struct{}
	block          bool
	size           int32
	quantity       int32
	group          *sync.Pool
	defaultHandler Handler
	taskHandler    map[string]Handles
	panicHandler   func()
}

func (e *Event) Close() {
	if e.closed != nil {
		close(e.closed)
		e.closed = nil
	}
}

func (e *Event) Add(name string, h Handler) {
	e.taskHandler[name] = append(e.taskHandler[name], h)
}

func (e *Event) SetDefaultHandler(h Handler) {
	e.defaultHandler = h
}

func (e *Event) SetPanicHandler(h func()) {
	e.panicHandler = h
}

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

// New 创建事件处理器
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

func (e *Event) job(w *sync.WaitGroup, h Handler, t *Task) {
	defer w.Done()
	if e.panicHandler != nil {
		defer e.panicHandler()
	}
	h(t)
}

func newGroup() interface{} {
	return new(sync.WaitGroup)
}
