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

import (
	"context"
)

type (
	// Writer is a write option.
	Writer interface {
		Write(v interface{}) (success bool)
	}
	// WriteBarrier is a write barrier.
	WriteBarrier struct {
		ctx          context.Context
		writeChannel chan<- interface{}
	}
)

// NewWriteBarrier returns a WriteBarrier.
func NewWriteBarrier(ctx context.Context, writeChannel chan<- interface{}) *WriteBarrier {
	return &WriteBarrier{
		ctx:          ctx,
		writeChannel: writeChannel,
	}
}

// Write the value to the writeChannel channel.
func (w *WriteBarrier) Write(v interface{}) (success bool) {
	select {
	case <-w.ctx.Done():
		return
	default:
		w.writeChannel <- v
		return true
	}
}
