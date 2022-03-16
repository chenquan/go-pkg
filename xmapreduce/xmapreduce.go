/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package xmapreduce

import (
	"context"
	"sync"

	"github.com/chenquan/go-pkg/xbarrier"
	"github.com/chenquan/go-pkg/xstream"
	"github.com/chenquan/go-pkg/xworker"
)

type (
	// GenerateFunc is used to let callers send elements into source.
	GenerateFunc[T any] func(source chan<- T)

	// MapFunc is used to do element processing and write the output to writer.
	MapFunc[T any] func(item T, writer xbarrier.Writer[T])

	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	ReducerFunc[T any] func(pipe <-chan interface{}, writer xbarrier.Writer[T], cancel func(error))

	options struct {
		workerSize int
	}

	// Option defines the method to customize the mapreduce.
	Option func(opts *options)
)

// WithWorkerSize customizes a mapreduce processing with given workers.
func WithWorkerSize(workerSize int) Option {
	return func(opts *options) {
		opts.workerSize = workerSize
	}
}

func loadOption(opts ...Option) *options {
	opt := &options{workerSize: 16}

	for _, option := range opts {
		option(opt)
	}

	return opt
}

// Map maps all elements generated from given generate func, and returns an output channel.
func Map[T any](ctx context.Context, generateFunc GenerateFunc[T], mapFunc MapFunc[T], opts ...Option) <-chan T {
	option := loadOption(opts...)
	source := buildSource[T](generateFunc)

	collector := make(chan T, option.workerSize)
	go doMap[T](ctx, mapFunc, source, collector, option)

	return collector
}

// MapStream maps all elements generated from given generate func, and returns a xstream.Stream.
func MapStream[T any](ctx context.Context, generateFunc GenerateFunc[T], mapFunc MapFunc[T], opts ...Option) *xstream.Stream[T] {
	return xstream.Range[T](Map[T](ctx, generateFunc, mapFunc, opts...))
}

func buildSource[T any](generateFunc GenerateFunc[T]) chan T {
	source := make(chan T)

	go func() {
		defer close(source)
		generateFunc(source)
	}()

	return source
}

func doMap[T any](ctx context.Context, mapFunc MapFunc[T], source <-chan T, collector chan<- T, option *options) {
	waitGroup := sync.WaitGroup{}

	defer func() {
		waitGroup.Wait()
		close(collector)
	}()
	worker := xworker.NewWorker(option.workerSize)
	writer := xbarrier.NewWriteBarrier(ctx, collector)

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-source:
			if !ok {
				return
			}
			waitGroup.Add(1)
			worker.Run(ctx, func() {
				mapFunc(item, writer)
			}, func() {
				waitGroup.Done()
			})
		}
	}
}
