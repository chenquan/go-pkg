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

import "sync"

type Chan struct {
	channel chan interface{}
	sync.Once
}

func NewChan() *Chan {
	return &Chan{channel: make(chan interface{}, 1)}
}

func (c *Chan) Write(v interface{}) (success bool) {
	c.Do(func() {
		c.channel <- v
		success = true
	})
	return
}

func (c *Chan) Chan() chan interface{} {
	return c.channel
}
