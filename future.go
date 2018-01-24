package fan

import (
	"sync"
)

// A Future is a value which can be written to once and read many
// times. A read of a Future blocks until it has been written to.
//
// A zero-valued Future is ready to use.
type Future struct {
	once sync.Once
	sent chan struct{}
	val  interface{}
}

func (f *Future) init() {
	f.once.Do(func() {
		f.sent = make(chan struct{})
	})
}

// Send writes a value to the Future. This should only ever be done
// once. A second call to this function will cause a panic.
func (f *Future) Send(val interface{}) {
	f.init()

	f.val = val
	close(f.sent)
}

// Get returns the value that was given to Send. If Send has not yet
// been called, Get blocks until it has been.
func (f *Future) Get() interface{} {
	f.init()

	<-f.sent
	return f.val
}

// Sent returns a channel which is closed when Send is called. This
// allows a Future to be used in a select statement.
//
// Note that it is *not* necessary to call this before calling Get, as
// Get will block on its own.
func (f *Future) Sent() <-chan struct{} {
	f.init()

	return f.sent
}
