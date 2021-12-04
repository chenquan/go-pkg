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

package xbarrier

import "context"

type (
	// Reader is a read option.
	Reader interface {
		Read() (val interface{}, success bool)
	}
	// ReadBarrier is a read barrier.
	ReadBarrier struct {
		ctx         context.Context
		readChannel <-chan interface{}
	}
)

// NewReadBarrier returns a NewReadBarrier.
func NewReadBarrier(ctx context.Context, readChannel <-chan interface{}) *ReadBarrier {
	return &ReadBarrier{ctx: ctx, readChannel: readChannel}
}

// Read data from the readChannel channel.
func (r *ReadBarrier) Read() (val interface{}, success bool) {
	for {
		select {
		case <-r.ctx.Done():
			return nil, false
		default:
			return <-r.readChannel, true
		}
	}
}
