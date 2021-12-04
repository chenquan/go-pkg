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
	"github.com/chenquan/go-pkg/xstream"
	"sync"
)

type (
	Writer interface {
		Write(v interface{})
	}

	GuardedWriter struct {
		channel chan<- interface{}
		ctx     context.Context
	}
)

func NewGuardedWriter(channel chan<- interface{}, ctx context.Context) GuardedWriter {
	return GuardedWriter{
		channel: channel,
		ctx:     ctx,
	}
}

func (gw GuardedWriter) Write(v interface{}) {
	select {
	case <-gw.ctx.Done():
		return
	default:
		gw.channel <- v
	}
}

type (
	// GenerateFunc is used to let callers send elements into source.
	GenerateFunc func(source chan<- interface{})
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)
	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))
	Options     struct {
		workerSize int
	}
	Option func(opts *Options)
)

func WithWorkerSize(workerSize int) Option {
	return func(opts *Options) {
		opts.workerSize = workerSize
	}
}
func loadOption(opts ...Option) *Options {
	opt := &Options{workerSize: 16}

	for _, option := range opts {
		option(opt)
	}
	return opt
}

func Map(ctx context.Context, generateFunc GenerateFunc, mapFunc MapFunc, opts ...Option) <-chan interface{} {
	option := loadOption(opts...)
	source := buildSource(generateFunc)

	collector := make(chan interface{}, option.workerSize)
	go doMap(ctx, mapFunc, source, collector, option)
	return collector
}

func MapStream(ctx context.Context, generateFunc GenerateFunc, mapFunc MapFunc, opts ...Option) *xstream.Stream {
	return xstream.Range(Map(ctx, generateFunc, mapFunc, opts...))
}

func buildSource(generateFunc GenerateFunc) chan interface{} {
	source := make(chan interface{})

	go func() {
		defer close(source)
		generateFunc(source)
	}()

	return source
}

func doMap(ctx context.Context, mapFunc MapFunc, source <-chan interface{}, collector chan<- interface{}, option *Options) {
	waitGroup := sync.WaitGroup{}
	defer func() {
		waitGroup.Wait()
		close(collector)
	}()

	workerChan := make(chan struct{}, option.workerSize)
	writer := NewGuardedWriter(collector, ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case workerChan <- struct{}{}:
			item, ok := <-source
			if !ok {
				<-workerChan
				return
			}
			waitGroup.Add(1)
			go func(value interface{}) {
				defer func() {
					waitGroup.Done()
					<-workerChan
				}()
				mapFunc(value, writer)
			}(item)
		}

	}
}
