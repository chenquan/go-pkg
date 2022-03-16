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
	"sort"
	"sync"
)

type (
	// FilterFunc defines the method to filter a Stream.
	FilterFunc[T any] func(item T) bool
	// ForAllFunc defines the method to handle all elements in a Stream.
	ForAllFunc[T any] func(pipe <-chan T)
	// ForEachFunc defines the method to handle each element in a Stream.
	ForEachFunc[T any] func(item T)
	// GenerateFunc defines the method to send elements into a Stream.
	GenerateFunc[T any] func(source chan<- T)
	// KeyFunc defines the method to generate keys for the elements in a Stream.
	KeyFunc[V any] func(v V) interface{}
	// LessFunc defines the method to compare the elements in a Stream.
	LessFunc[T any] func(a, b T) bool
	// MapFunc defines the method to map each element to another object in a Stream.
	MapFunc[T any] func(item T) T
	// ParallelFunc defines the method to handle elements parallelly.
	ParallelFunc[T any] func(item T)
	// ReduceFunc defines the method to reduce all the elements in a Stream.
	ReduceFunc func(pipe <-chan any) (any, error)
	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc[T any] func(item T, pipe chan<- T)
	// Collector represents a stream collector to collect items
	Collector[T any] interface {
		Input(c <-chan T)
	}
	CollectorFunc[T any] func(c <-chan T)
)

// Input implements Collector.
func (cf CollectorFunc[T]) Input(c <-chan T) {
	cf(c)
}

// Stream Represents a stream.
type Stream[T any] struct {
	source <-chan T
}

// Empty Returns an empty stream.
func Empty[T any]() *Stream[T] {
	source := make(chan T)
	close(source)
	return &Stream[T]{source}
}

// Range Returns a Stream from source channel.
func Range[T any](source <-chan T) *Stream[T] {
	return &Stream[T]{
		source: source,
	}
}

// Of Returns a Stream based any element
func Of[T any](items ...T) *Stream[T] {
	n := len(items)
	if n == 0 {
		return Empty[T]()
	}

	source := make(chan T, n)
	go func() {
		for _, item := range items {
			source <- item
		}
		close(source)
	}()
	return Range[T](source)
}

// Concat Returns a concat Stream.
func Concat[T any](a *Stream[T], others ...*Stream[T]) *Stream[T] {
	return a.Concat(others...)
}

// From Returns a Stream from generate function.
func From[T any](generate GenerateFunc[T]) *Stream[T] {
	source := make(chan T)

	go func() {
		defer close(source)
		generate(source)
	}()

	return Range[T](source)
}

// Distinct Returns a distinct Stream.
func (s *Stream[T]) Distinct(f KeyFunc[T]) *Stream[T] {
	source := make(chan T)

	go func() {
		defer close(source)

		unique := make(map[interface{}]struct{})
		for item := range s.source {
			k := f(item)
			if _, ok := unique[k]; !ok {
				source <- item
				unique[k] = struct{}{}
			}
		}
	}()
	return Range[T](source)
}

// Count Returns a number that the elements total size.
func (s *Stream[T]) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// Buffer Returns a buffer Stream.
func (s *Stream[T]) Buffer(n int) *Stream[T] {
	if n < 0 {
		n = 0
	}
	source := make(chan T, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range[T](source)
}

// Done Stream.
func (s *Stream[T]) Done() {
	drain[T](s.source)
}

// Chan Returns a channel of Stream.
func (s *Stream[T]) Chan() <-chan T {
	return s.source
}

// Split Returns a split Stream that contains multiple slices of chunk size n.
// func (s *Stream[T]) Split(n int) *Stream[[]T] {
//	if n < 1 {
//		panic("n should be greater than 0")
//	}
//	source := make(chan []T)
//	go func() {
//		var chunk []T
//		for item := range s.source {
//			chunk = append(chunk, item)
//			if len(chunk) == n {
//				source <- chunk
//				chunk = nil
//			}
//		}
//		if chunk != nil {
//			source <- chunk
//		}
//		close(source)
//	}()
//	return Range[[]T](source)
// }

// SplitSteam Returns a split Stream that contains multiple stream of chunk size n.
// func (s *Stream[T]) SplitSteam(n int) *Stream[chan T] {
//	if n < 1 {
//		go drain(s.source)
//		panic("n should be greater than 0")
//	}
//	source := make(chan *Stream[T])
//
//	go func() {
//
//		var chunkSource = make(chan T, n)
//		for item := range s.source {
//			chunkSource <- item
//			if len(chunkSource) == n {
//
//				source <- Range[T](chunkSource)
//				close(chunkSource)
//
//				chunkSource = make(chan T, n)
//			}
//		}
//		if len(chunkSource) != 0 {
//			source <- Range[T](chunkSource)
//			close(chunkSource)
//		}
//		close(source)
//	}()
//
//	return Range[chan T](source)
// }

// Sort Returns a sorted Stream.
func (s *Stream[T]) Sort(less LessFunc[T]) *Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Of(items...)
}

// Tail Returns a Stream that has n element at the end.
func (s *Stream[T]) Tail(n int) *Stream[T] {
	if n <= 0 {
		go drain[T](s.source)
		if n == 0 {
			return Empty[T]()
		}
		panic("n should be greater than 0")
	}

	source := make(chan T)

	go func() {
		defer close(source)

		r := newRing[T](uint(n))
		for item := range s.source {
			r.add(item)
		}
		for _, item := range r.take() {
			source <- item
		}
	}()

	return Range[T](source)
}

// Skip Returns a Stream that skips size elements.
func (s *Stream[T]) Skip(size int) *Stream[T] {
	if size < 0 {
		go drain[T](s.source)
		panic("size should be greater than 0")
	}

	if size == 0 {
		return s
	}

	source := make(chan T)

	go func() {
		defer close(source)

		i := 0
		for item := range s.source {
			if i >= size {
				source <- item
			}
			i++
		}
	}()

	return Range[T](source)
}

// Limit Returns a Stream that contains size elements.
func (s *Stream[T]) Limit(size int) *Stream[T] {
	if size == 0 {
		go drain(s.source)
		return Empty[T]()
	}
	if size < 0 {
		panic("size must be greater than -1")
	}
	source := make(chan T)

	go func() {

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
	}()

	return Range[T](source)
}

// Foreach Traversals all elements.
func (s *Stream[T]) Foreach(f ForEachFunc[T]) {
	for item := range s.source {
		f(item)
	}
}

// ForeachOrdered Traversals all elements in reverse order.
func (s *Stream[T]) ForeachOrdered(f ForEachFunc[T]) {
	items := make([]T, 0)
	for item := range s.source {
		items = append(items, item)
	}
	n := len(items)
	for i := n - 1; i >= 0; i-- {
		f(items[i])
	}
}

// Concat Returns a Stream that concat others streams
func (s *Stream[T]) Concat(others ...*Stream[T]) *Stream[T] {
	source := make(chan T)
	wg := sync.WaitGroup{}

	go func() {
		for _, other := range others {

			wg.Add(1)
			go func(s *Stream[T]) {
				for item := range s.source {
					source <- item
				}
				wg.Done()
			}(other)

		}

		wg.Add(1)
		go func() {
			for item := range s.source {
				source <- item
			}
			wg.Done()
		}()

		wg.Wait()
		close(source)

	}()

	return Range[T](source)
}

// Filter Returns a Stream that
func (s *Stream[T]) Filter(fn FilterFunc[T], opts ...Option) *Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// Walk Returns a Stream that lets the callers handle each item, the caller may write zero,
// one or more items base on the given item.
func (s *Stream[T]) Walk(f WalkFunc[T], opts ...Option) *Stream[T] {
	option := loadOptions(opts...)
	pipe := make(chan T, option.workSize)
	go func() {
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
			go func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				f(item, pipe)
			}()
		}
		wg.Wait()
		close(pipe)
	}()

	return Range[T](pipe)
}

// Map Returns a Stream consisting of the results of applying the given
// function to the elements of this stream.
func (s *Stream[T]) Map(fn MapFunc[T], opts ...Option) *Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		pipe <- fn(item)
	}, opts...)
}

// // FlatMap Returns a Stream consisting of the results of replacing each element of this stream with the contents of
// // a mapped stream produced by applying the provided mapping function to each element. Each mapped stream is closed
// // after its contents have been placed into this stream. (If a mapped stream is null an empty stream is used, instead.
// func (s *Stream) FlatMap(fn MapFunc, opts ...Option) *Stream {
//	return s.Walk(func(item any, pipe chan<- any) {
//		switch v := item.(type) {
//		case []any:
//			for _, x := range v {
//				pipe <- fn(x)
//			}
//		case any:
//			pipe <- fn(v)
//		}
//	}, opts...)
// }
//
// // Group Returns a Stream that groups the elements into different groups based on their keys.
// func (s *Stream) Group(f KeyFunc) *Stream {
//	groups := make(map[any][]any)
//	for item := range s.source {
//		key := f(item)
//		groups[key] = append(groups[key], item)
//	}
//
//	source := make(chan any)
//	go func() {
//		for _, group := range groups {
//			source <- group
//		}
//		close(source)
//	}()
//
//	return Range(source)
// }
//

// Merge Returns a Stream that merges all the items into a slice and generates a new stream.
func (s *Stream[T]) Merge() *Stream[interface{}] {
	source := make(chan interface{})

	go func() {
		var items []T
		for item := range s.source {
			items = append(items, item)
		}

		source <- items
		close(source)
	}()

	return Range[interface{}](source)
}

func To[T any](from *Stream[interface{}]) *Stream[T] {
	return From[T](func(source chan<- T) {
		for v := range from.source {
			source <- v.(T)
		}
	})
}

// Reverse Returns a Stream that reverses the elements.
func (s *Stream[T]) Reverse() *Stream[T] {
	var items []T
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
func (s *Stream[T]) ParallelFinish(fn ParallelFunc[T], opts ...Option) {
	s.Walk(func(item T, pipe chan<- T) {
		fn(item)
	}, opts...).Done()
}

// AnyMach Returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s *Stream[T]) AnyMach(f func(item T) bool) (isFind bool) {
	for item := range s.source {
		if f(item) {
			isFind = true
			go drain(s.source)

			return
		}
	}
	return
}

// AllMach Returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s *Stream[T]) AllMach(f func(item T) bool) (isFind bool) {
	isFind = true
	for item := range s.source {
		if !f(item) {
			isFind = false
			go drain(s.source)

			return
		}
	}
	return
}

// FindFirst Returns an interface{} the first element of this stream, or a nil and a error if the stream is empty.
// If the stream has no encounter order, then any element may be returned
func (s *Stream[T]) FindFirst() (result T, err error) {

	for result = range s.source {
		go drain(s.source)
		return
	}

	err = errors.New("no element")
	return
}

// FindLast Returns an interface{} the last element of this stream, or a nil and a error if the stream is empty.
// If the stream has no encounter order, then any element may be returned
func (s *Stream[T]) FindLast() (result T, err error) {
	flag := true
	for result = range s.source {
		flag = false
	}

	if flag {
		err = errors.New("no element")
	}

	return
}

// Peek Returns a Stream consisting of the elements of this stream,
// additionally performing the provided action on each element as elements are consumed from the resulting stream.
func (s *Stream[T]) Peek(f ForEachFunc[T]) *Stream[T] {
	source := make(chan T)
	go func() {
		for item := range s.source {
			source <- item
			f(item)
		}
		close(source)
	}()

	return Range[T](source)
}

func drain[T any](channel <-chan T) {
	for range channel {
	}
}

// Collection collects a Stream.
func (s *Stream[T]) Collection(collector Collector[T]) {
	collector.Input(s.source)
}
