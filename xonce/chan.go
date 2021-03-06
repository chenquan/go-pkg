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

package xonce

import (
	"sync/atomic"
)

// Chan represents a channel that can only be written too once.
type Chan struct {
	channel chan interface{}
	wrote   uint32
}

// NewChan returns a Chan.
func NewChan() *Chan {
	return &Chan{channel: make(chan interface{}, 1)}
}

// Write writes a v.
func (c *Chan) Write(v interface{}) (success bool) {
	if success = atomic.CompareAndSwapUint32(&c.wrote, 0, 1); success {
		c.channel <- v
		close(c.channel)
	}

	return
}

// Chan returns a channel.
func (c *Chan) Chan() chan interface{} {
	return c.channel
}
