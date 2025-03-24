package event

import "runtime"

var object *Event

func init() {
	num := int32(runtime.NumCPU())
	object = New(num*4, true)
}

// Add a handler to event.
func Add(name string, h Handler) {
	object.Add(name, h)
}

// SetDefaultHandler set default handler.
func SetDefaultHandler(h Handler) {
	object.SetDefaultHandler(h)
}

// SetPanicHandler set panic handler.
func SetPanicHandler(h func()) {
	object.SetPanicHandler(h)
}

// Push event, return error if event queue is full.
func Push(name string, data interface{}, attr map[string]any) error {
	return object.Push(name, data, attr)
}
