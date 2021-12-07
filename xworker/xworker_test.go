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

package xworker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWorker(t *testing.T) {
	worker := NewWorker(1)
	group := sync.WaitGroup{}
	k := 0
	for i := 0; i < 10; i++ {
		group.Add(1)
		worker.Run(context.Background(), func() {

			k++
		}, func() {
			group.Done()
		})
		assert.EqualValues(t, i+1, k)
	}
	group.Wait()

	worker = NewWorker(10)
	group = sync.WaitGroup{}
	j := uint32(0)
	for i := 0; i < 100; i++ {
		group.Add(1)
		worker.Run(context.Background(), func() {

			atomic.AddUint32(&j, 1)
		}, func() {
			group.Done()
		})

	}
	group.Wait()
	assert.Equal(t, uint32(100), atomic.LoadUint32(&j))

	worker = NewWorker(1)
	group = sync.WaitGroup{}
	j = uint32(0)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	for i := 0; i < 100; i++ {
		group.Add(1)
		worker.Run(ctx, func() {

			if atomic.AddUint32(&j, 1) == 50 {
				cancelFunc()
			}
		}, func() {
			group.Done()
		})

	}
	group.Wait()
	assert.Equal(t, uint32(50), atomic.LoadUint32(&j))

}
