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
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTryLock(t *testing.T) {
	var lock Spinlock
	assert.True(t, lock.TryLock())
	assert.False(t, lock.TryLock())
	lock.Unlock()
	assert.True(t, lock.TryLock())
}

func TestSpinLock(t *testing.T) {
	var lock Spinlock
	lock.Lock()
	assert.False(t, lock.TryLock())
	lock.Unlock()
	assert.True(t, lock.TryLock())
}

func TestSpinLockRace(t *testing.T) {
	var lock Spinlock
	lock.Lock()
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		wait.Done()
	}()
	time.Sleep(time.Millisecond * 100)
	lock.Unlock()
	wait.Wait()
	assert.True(t, lock.TryLock())
}

func TestSpinLock_TryLock(t *testing.T) {
	var lock Spinlock
	var count int32
	var wait sync.WaitGroup
	wait.Add(2)
	sig := make(chan struct{})

	go func() {
		lock.TryLock()
		sig <- struct{}{}
		atomic.AddInt32(&count, 1)
		runtime.Gosched()
		lock.Unlock()
		wait.Done()
	}()

	go func() {
		<-sig
		lock.Lock()
		atomic.AddInt32(&count, 1)
		lock.Unlock()
		wait.Done()
	}()

	wait.Wait()
	assert.Equal(t, int32(2), atomic.LoadInt32(&count))
}
