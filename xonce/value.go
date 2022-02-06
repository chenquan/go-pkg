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

// Value represents a value that can only be written too once.
type Value struct {
	v    interface{}
	once sync.Once
}

// NewValue returns a Value.
func NewValue() *Value {
	return &Value{}
}

// Write writes a v.
func (val *Value) Write(v interface{}) (success bool) {
	val.once.Do(func() {
		val.v = v
		success = true
	})
	return
}

// Value return a v.
func (val *Value) Value() interface{} {
	return val.v
}
