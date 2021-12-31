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

package xsync

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue()
	go func() {
		for i := 0; i < 100; i++ {
			queue.Write(i)
		}
	}()

	for i := 0; i < 100; i++ {
		assert.Equal(t, i, queue.Read())
	}

}

func TestQueue_Close(t *testing.T) {
	queue := NewQueue()
	assert.NoError(t, queue.Close())
	assert.Panics(t, func() {
		queue.Read()
	})
}

func TestQueue_Remove(t *testing.T) {
	queue := NewQueue()
	queue.Write(1)
	queue.Write(2)
	assert.True(t, queue.Remove(2))
	assert.False(t, queue.Remove(3))
	assert.Equal(t, 1, queue.Read())
}

func TestWithQueueCap(t *testing.T) {
	// -----------

	newQueue := NewQueue(WithQueueCap(2))
	go func() {
		for i := 0; i < 100; i++ {
			newQueue.Write(i)
			if i == 50 {
				newQueue.Cap(200)
				assert.Equal(t, 200, newQueue.cap)
			}
		}
	}()

	for i := 0; i < 100; i++ {
		assert.Equal(t, i, newQueue.Read())
	}
}

func TestWithQueueMode(t *testing.T) {
	opt := WithQueueMode(MultiWrite | MultiRead | MultiReadWrite)
	opts := loadQueueOpts(opt)
	assert.Equal(t, MultiReadWrite, opts.mode)

	opt = WithQueueMode(MultiWrite | MultiRead)
	opts = loadQueueOpts(opt)
	assert.Equal(t, MultiReadWrite, opts.mode)

	opt = WithQueueMode(MultiWrite)
	opts = loadQueueOpts(opt)
	assert.Equal(t, MultiWrite, opts.mode)

	opt = WithQueueMode(MultiRead)
	opts = loadQueueOpts(opt)
	assert.Equal(t, MultiRead, opts.mode)
}

func TestQueue_Cap(t *testing.T) {

	t.Run("pipeline", func(t *testing.T) {
		queue := NewQueue(WithQueueCap(2))
		queue.Write(1)
		queue.Write(1)

		n := 10
		awake := make(chan struct{})
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				queue.Write(1)
				awake <- struct{}{}
			}()

		}
		for i := 0; i < n; i++ {
			queue.Read()
			<-awake
			select {
			case <-awake:
				t.Fatal("goroutine not asleep")
			default:
			}
		}
		wg.Wait()
	})

}
