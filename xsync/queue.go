package xsync

import (
	"container/list"
	"sync"
)

const (
	MultiRead      QueueMode = 1
	MultiWrite     QueueMode = 2
	MultiReadWrite           = MultiRead | MultiWrite
)

type (
	Queue struct {
		cond *sync.Cond
		l    *list.List
		cap  int
		mode QueueMode
	}
	QueueMode    byte
	QueueOptions struct {
		cap  int
		mode QueueMode
	}
	QueueOption func(opt *QueueOptions)
)

func WithQueueCap(cap int) QueueOption {
	return func(opt *QueueOptions) {
		opt.cap = cap
	}
}

func WithQueueMode(mode QueueMode) QueueOption {
	return func(opt *QueueOptions) {
		opt.mode = mode
	}
}

func NewQueue(opts ...QueueOption) *Queue {
	options := loadQueueOpts(opts...)

	return &Queue{cond: sync.NewCond(&sync.Mutex{}), l: list.New(), cap: options.cap}
}

func loadQueueOpts(opts ...QueueOption) *QueueOptions {
	opt := new(QueueOptions)
	for _, option := range opts {
		option(opt)
	}

	return opt
}

func (q *Queue) Close() error {
	q.cond.L.Lock()
	q.l = nil
	q.cond.L.Unlock()

	return nil
}

func (q *Queue) Read() interface{} {
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()

		// wake
		if q.mode == MultiWrite || q.mode == MultiReadWrite {
			q.cond.Broadcast()
		} else {
			q.cond.Signal()
		}
	}()

	for q.l.Len() == 0 {
		q.cond.Wait()
	}

	element := q.l.Front()
	q.l.Remove(element)

	return element.Value
}

func (q *Queue) Write(v interface{}) {
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()

		// wake
		if q.mode == MultiRead || q.mode == MultiReadWrite {
			q.cond.Broadcast()
		} else {
			q.cond.Signal()
		}
	}()

	for q.cap != 0 && q.l.Len() >= q.cap {
		// Waiting if the current number of elements in the queue is greater than the capacity.
		q.cond.Wait()
	}

	q.l.PushBack(v)
}

func (q *Queue) Cap(cap int) {
	q.cond.L.Lock()
	q.cap = cap
	q.cond.L.Unlock()
}

func (q *Queue) Remove(v interface{}) (ok bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for element := q.l.Front(); element != nil; element = element.Next() {
		if element.Value == v {
			q.l.Remove(element)
			return true
		}
	}

	return false
}
