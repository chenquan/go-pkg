/*
 *
 *     Copyright 2021 chenquan
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */

package xstream

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"sort"
	"sync"
)

var (
	// empty an empty Stream.
	empty *Stream
	// pool a  ants pool.
	pool, _ = ants.NewPool(-1)
	// ErrNoElement no element error.
	ErrNoElement = errors.New("no element")
)

func init() {
	source := make(chan interface{})
	close(source)
	empty = &Stream{source: source}
}

type (
	// FilterFunc defines the method to filter a Stream.
	FilterFunc func(item interface{}) bool
	// ForAllFunc defines the method to handle all elements in a Stream.
	ForAllFunc func(pipe <-chan interface{})
	// ForEachFunc defines the method to handle each element in a Stream.
	ForEachFunc func(item interface{})
	// GenerateFunc defines the method to send elements into a Stream.
	GenerateFunc func(source chan<- interface{})
	// KeyFunc defines the method to generate keys for the elements in a Stream.
	KeyFunc func(item interface{}) interface{}
	// LessFunc defines the method to compare the elements in a Stream.
	LessFunc func(a, b interface{}) bool
	// MapFunc defines the method to map each element to another object in a Stream.
	MapFunc func(item interface{}) interface{}
	// ParallelFunc defines the method to handle elements parallelly.
	ParallelFunc func(item interface{})
	// ReduceFunc defines the method to reduce all the elements in a Stream.
	ReduceFunc func(pipe <-chan interface{}) (interface{}, error)
	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc func(item interface{}, pipe chan<- interface{})
	// Collector represents a stream collector to collect items
	Collector interface {
		Input(c <-chan interface{})
	}
	// CollectorFunc represents a stream collector.
	CollectorFunc func(c <-chan interface{})
)

// Stream Represents a stream.
type Stream struct {
	source <-chan interface{}
}

// Input implements Collector.
func (cf CollectorFunc) Input(c <-chan interface{}) {
	cf(c)
}

// -------------

func startGoroutine(f func()) {
	_ = pool.Submit(f)
}

// -------------

// Empty Returns an empty stream.
func Empty() *Stream {
	return empty
}

// Range Returns a Stream from source channel.
func Range(source <-chan interface{}) *Stream {
	return &Stream{
		source: source,
	}
}

// Of Returns a Stream based any element
func Of(items ...interface{}) *Stream {
	n := len(items)
	if n == 0 {
		return Empty()
	}

	source := make(chan interface{}, n)
	startGoroutine(func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	})
	return Range(source)
}

// Concat Returns a concat Stream.
func Concat(a *Stream, others ...*Stream) *Stream {
	return a.Concat(others...)
}

// From Returns a Stream from generate function.
func From(generate GenerateFunc) *Stream {
	source := make(chan interface{})
	startGoroutine(func() {
		defer close(source)
		generate(source)
	})
	return Range(source)
}

// Distinct Returns a distinct Stream.
func (s *Stream) Distinct(f KeyFunc) *Stream {
	source := make(chan interface{})

	startGoroutine(func() {
		defer close(source)

		unique := make(map[interface{}]struct{})
		for item := range s.source {
			k := f(item)
			if _, ok := unique[k]; !ok {
				source <- item
				unique[k] = struct{}{}
			}
		}
	})
	return Range(source)
}

// Count Returns a number that the elements total size.
func (s *Stream) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// Buffer Returns a buffer Stream.
func (s *Stream) Buffer(n int) *Stream {
	if n < 0 {
		n = 0
	}
	source := make(chan interface{}, n)
	startGoroutine(func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	})

	return Range(source)
}

// Done Stream.
func (s *Stream) Done() {
	drain(s.source)
}

// Chan Returns a channel of Stream.
func (s *Stream) Chan() <-chan interface{} {
	return s.source
}

// Split Returns a split Stream that contains multiple slices of chunk size n.
func (s *Stream) Split(n int) *Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}
	source := make(chan interface{})
	startGoroutine(func() {
		var chunk []interface{}
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	})
	return Range(source)
}

// SplitSteam Returns a split Stream that contains multiple stream of chunk size n.
func (s *Stream) SplitSteam(n int) *Stream {
	if n < 1 {
		startGoroutine(func() {
			drain(s.source)
		})
		panic("n should be greater than 0")
	}
	source := make(chan interface{})
	startGoroutine(func() {
		var chunkSource = make(chan interface{}, n)
		for item := range s.source {
			chunkSource <- item
			if len(chunkSource) == n {

				source <- Range(chunkSource)
				close(chunkSource)

				chunkSource = make(chan interface{}, n)
			}
		}
		if len(chunkSource) != 0 {
			source <- Range(chunkSource)
			close(chunkSource)
		}
		close(source)
	})

	return Range(source)
}

// Sort Returns a sorted Stream.
func (s *Stream) Sort(less LessFunc) *Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Of(items...)
}

// Tail Returns a Stream that has n element at the end.
func (s *Stream) Tail(n int) *Stream {
	if n <= 0 {
		startGoroutine(func() {
			drain(s.source)
		})
		if n == 0 {
			return empty
		}
		panic("n should be greater than 0")
	}

	source := make(chan interface{})
	startGoroutine(func() {
		defer close(source)

		r := newRing(uint(n))
		for item := range s.source {
			r.add(item)
		}
		for _, item := range r.take() {
			source <- item
		}
	})

	return Range(source)
}

// Skip Returns a Stream that skips size elements.
func (s *Stream) Skip(size int) *Stream {
	if size < 0 {
		startGoroutine(func() {
			drain(s.source)
		})
		panic("size should be greater than 0")
	}

	if size == 0 {
		return s
	}

	source := make(chan interface{})
	startGoroutine(func() {
		defer close(source)

		i := 0
		for item := range s.source {
			if i >= size {
				source <- item
			}
			i++
		}
	})

	return Range(source)
}

// Limit Returns a Stream that contains size elements.
func (s *Stream) Limit(size int) *Stream {
	if size == 0 {
		startGoroutine(func() {
			drain(s.source)
		})
		return Empty()
	}
	if size < 0 {
		panic("size must be greater than -1")
	}
	source := make(chan interface{})
	startGoroutine(func() {

		for item := range s.source {
			if size > 0 {
				source <- item
			}

			size--

			if size == 0 {
				close(source)
			}
		}

		if size > 0 {
			close(source)
		}
	})

	return Range(source)
}

// Foreach Traversals all elements.
func (s *Stream) Foreach(f ForEachFunc) {
	for item := range s.source {
		f(item)
	}
}

// ForeachOrdered Traversals all elements in reverse order.
func (s *Stream) ForeachOrdered(f ForEachFunc) {
	items := make([]interface{}, 0)
	for item := range s.source {
		items = append(items, item)
	}
	n := len(items)
	for i := n - 1; i >= 0; i-- {
		f(items[i])
	}
}

// Concat Returns a Stream that concat others streams
func (s *Stream) Concat(others ...*Stream) *Stream {
	source := make(chan interface{})
	wg := sync.WaitGroup{}
	startGoroutine(func() {
		for _, other := range others {
			other := other
			wg.Add(1)
			startGoroutine(func() {
				for item := range other.source {
					source <- item
				}
				wg.Done()
			})

		}

		wg.Add(1)
		startGoroutine(func() {
			for item := range s.source {
				source <- item
			}
			wg.Done()
		})

		wg.Wait()
		close(source)

	})

	return Range(source)
}

// Filter Returns a Stream that
func (s *Stream) Filter(fn FilterFunc, opts ...Option) *Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// Walk Returns a Stream that lets the callers handle each item, the caller may write zero,
// one or more items base on the given item.
func (s *Stream) Walk(f WalkFunc, opts ...Option) *Stream {
	option := loadOptions(opts...)
	pipe := make(chan interface{}, option.workSize)
	startGoroutine(func() {
		var wg sync.WaitGroup
		pool := make(chan struct{}, option.workSize)

		for {
			pool <- struct{}{}
			item, ok := <-s.source
			if !ok {
				<-pool
				break
			}

			wg.Add(1)
			// better to safely run caller defined method
			startGoroutine(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				f(item, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	})

	return Range(pipe)
}

// Map Returns a Stream consisting of the results of applying the given
// function to the elements of this stream.
func (s *Stream) Map(fn MapFunc, opts ...Option) *Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opts...)
}

// FlatMap Returns a Stream consisting of the results of replacing each element of this stream with the contents of
// a mapped stream produced by applying the provided mapping function to each element. Each mapped stream is closed
// after its contents have been placed into this stream. (If a mapped stream is null an empty stream is used, instead.
func (s *Stream) FlatMap(fn MapFunc, opts ...Option) *Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		switch v := item.(type) {
		case []interface{}:
			for _, x := range v {
				pipe <- fn(x)
			}
		case interface{}:
			pipe <- fn(v)
		}
	}, opts...)
}

// Group Returns a Stream that groups the elements into different groups based on their keys.
func (s *Stream) Group(f KeyFunc) *Stream {
	groups := make(map[interface{}][]interface{})
	for item := range s.source {
		key := f(item)
		groups[key] = append(groups[key], item)
	}

	source := make(chan interface{})
	startGoroutine(func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	})

	return Range(source)
}

// Merge Returns a Stream that merges all the items into a slice and generates a new stream.
func (s *Stream) Merge() *Stream {
	source := make(chan interface{})
	startGoroutine(func() {
		var items []interface{}
		for item := range s.source {
			items = append(items, item)
		}

		source <- items
		close(source)
	})

	return Range(source)
}

// Reverse Returns a Stream that reverses the elements.
func (s *Stream) Reverse() *Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}

	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Of(items...)
}

// ParallelFinish applies the given ParallelFunc to each item concurrently with given number of workers
func (s *Stream) ParallelFinish(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item interface{}, pipe chan<- interface{}) {
		fn(item)
	}, opts...).Done()
}

// AnyMach Returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s *Stream) AnyMach(f func(item interface{}) bool) (isFind bool) {
	for item := range s.source {
		if f(item) {
			isFind = true
			startGoroutine(func() {
				drain(s.source)
			})

			return
		}
	}
	return
}

// AllMach Returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s *Stream) AllMach(f func(item interface{}) bool) (isFind bool) {
	isFind = true
	for item := range s.source {
		if !f(item) {
			isFind = false
			startGoroutine(func() {
				drain(s.source)
			})

			return
		}
	}
	return
}

// FindFirst Returns an interface{} the first element of this stream, or a nil and a error if the stream is empty.
// If the stream has no encounter order, then any element may be returned
func (s *Stream) FindFirst() (result interface{}, err error) {

	for result = range s.source {
		startGoroutine(func() {
			drain(s.source)
		})
		return
	}

	err = ErrNoElement
	return
}

// FindLast Returns an interface{} the last element of this stream, or a nil and a error if the stream is empty.
// If the stream has no encounter order, then any element may be returned
func (s *Stream) FindLast() (result interface{}, err error) {
	flag := true
	for result = range s.source {
		flag = false
	}

	if flag {
		err = ErrNoElement
	}

	return
}

// Peek Returns a Stream consisting of the elements of this stream,
// additionally performing the provided action on each element as elements are consumed from the resulting stream.
func (s *Stream) Peek(f ForEachFunc) *Stream {
	source := make(chan interface{})
	startGoroutine(func() {
		for item := range s.source {
			source <- item
			f(item)
		}
		close(source)
	})

	return Range(source)
}

// Copy returns multiple streams copied.
// streamParam specifies the name and buffer size of the replicated stream.
func (s *Stream) Copy(streamParam map[string]int) (streamMap map[string]*Stream) {
	streamMap = map[string]*Stream{}
	channels := make([]chan interface{}, 0, len(streamParam))
	for name, bufferSize := range streamParam {
		c := make(chan interface{}, bufferSize)
		streamMap[name] = Range(c)
		channels = append(channels, c)
	}

	sort.Slice(channels, func(i, j int) bool {
		return cap(channels[i]) > cap(channels[j])
	})

	startGoroutine(func() {
		for v := range s.source {
			for _, c := range channels {
				c <- v
			}
		}
		for _, c := range channels {
			close(c)
		}
	})

	return
}

// Collection collects a Stream.
func (s *Stream) Collection(collector Collector) {
	collector.Input(s.source)
}

func drain(channel <-chan interface{}) {
	for range channel {
	}
}
