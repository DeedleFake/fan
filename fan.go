package fan

import "sync"

// A Fan coordinates Futures in a way that makes it easy to broadcast
// data to multiple other places. Essentially, it's a one-to-many,
// infinitely-buffered channel.
//
// A zero-valued Fan is ready to use.
type Fan struct {
	m   sync.RWMutex
	cur *Future
}

// Send sends data to the Fan.
func (f *Fan) Send(val interface{}) {
	f.m.Lock()
	defer f.m.Unlock()

	if f.cur == nil {
		f.cur = new(Future)
	}

	var next Future
	f.cur.Send(link{
		Val:  val,
		Next: &next,
	})

	f.cur = &next
}

// Receiver returns a Receiver that gets data from the Fan.
func (f *Fan) Receiver() *Receiver {
	f.m.RLock()
	switch f.cur {
	case nil:
		f.m.RUnlock()
		f.m.Lock()
		defer f.m.Unlock()

		if f.cur == nil {
			f.cur = new(Future)
		}

	default:
		defer f.m.RUnlock()
	}

	return &Receiver{
		cur: f.cur,
	}
}

// A Receiver receives values from a Fan. Use of it is not
// thread-safe; each client of the Fan should have its own, separate
// Receiver.
type Receiver struct {
	cur *Future
}

// Get gets the next value sent by the Fan, blocking until one is
// available. Values are not skipped, meaning that if the Fan has sent
// multiple values since the last call to Get, the next call to Get
// will return the first value sent since then.
func (r *Receiver) Get() interface{} {
	w := r.cur.Get().(link)
	r.cur = w.Next

	return w.Val
}

// Sent returns a channel that will be closed the next time a value is
// available to receive with Get. This works the same way as the
// method of the same name on Future.
//
// Note that the channel returned will no longer be valid after the
// next call to Get, so this function should be called again at that
// point.
func (r *Receiver) Sent() <-chan struct{} {
	return r.cur.Sent()
}

// A link is a simple struct for wrapping both data that has been sent
// and a pointer to the next place to which data will be sent.
type link struct {
	// Val is the data that was sent.
	Val interface{}

	// Next is the next place that data will be sent to.
	Next *Future
}
