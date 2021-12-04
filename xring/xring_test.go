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

package xring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRing(t *testing.T) {
	ring := NewRing(1)
	assert.Equal(t, &Ring{elements: make([]interface{}, 1)}, ring)
}
func TestRing_Add(t *testing.T) {
	ring := NewRing(12)
	ring.Add(1)
	assert.Equal(t, ring.index, 1)
	a := make([]interface{}, 12)
	a[0] = 1
	assert.Equal(t, a, ring.elements)
	ring.Add(2)
	a[1] = 2

	assert.Equal(t, a, ring.elements)

}
func TestRing_Take(t *testing.T) {
	ring := NewRing(4)
	ring.Add(1)
	ring.Add(2)
	ring.Add(3)
	assert.Equal(t, []interface{}{1, 2, 3}, ring.Take())
	ring.Add(4)
	ring.Add(5)
	ring.Add(6)
	ring.Add(7)
	ring.Add(8)
	assert.Equal(t, []interface{}{5, 6, 7, 8}, ring.Take())
	assert.Panics(t, func() {
		NewRing(-1)
	})
}
