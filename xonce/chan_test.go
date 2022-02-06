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
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestChan(t *testing.T) {

	value := NewChan()
	N := 1000
	c := make(chan struct{})
	flag := int32(-1)
	for i := 0; i < N; i++ {
		go func(i int) {

			if value.Write(i) {
				atomic.StoreInt32(&flag, int32(i))
			}

			c <- struct{}{}
		}(i)
	}

	for i := 0; i < N; i++ {
		<-c
	}

	i := 0
	for n := range value.Chan() {
		if i != 0 {
			t.Errorf("only written once.")
		}
		assert.EqualValues(t, flag, n)

		i++
	}
}
