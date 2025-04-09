package event

import "runtime"

// object 是全局的事件处理器实例
var object *Event

// init 初始化全局事件处理器
func init() {
	num := int32(runtime.NumCPU())
	object = New(num*4, true)
}

// Add 为指定事件名称添加处理函数（全局方法）
func Add(name string, h Handler) {
	object.Add(name, h)
}

// SetDefaultHandler 设置默认事件处理函数（全局方法）
func SetDefaultHandler(h Handler) {
	object.SetDefaultHandler(h)
}

// SetPanicHandler 设置异常处理函数（全局方法）
func SetPanicHandler(h func()) {
	object.SetPanicHandler(h)
}

// Push 将事件任务推入事件队列（全局方法），返回可能的错误
func Push(name string, data interface{}, attr map[string]any) error {
	return object.Push(name, data, attr)
}
