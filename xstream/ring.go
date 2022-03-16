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

package xstream

type ring[T any] struct {
	elements []T
	index    int
}

// newRing returns a ring object with the given size n.
func newRing[T any](n uint) *ring[T] {
	return &ring[T]{elements: make([]T, n)}
}

// add adds v into r.
func (r *ring[T]) add(v T) {

	r.elements[r.index%len(r.elements)] = v
	r.index++
}

// take all items from r.
func (r *ring[T]) take() []T {

	var size int
	var start int
	n := len(r.elements)
	if r.index > n {
		size = n
		start = r.index % n
	} else {
		size = r.index
	}

	elements := make([]T, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%n]
	}

	return elements
}
