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
	"runtime"
	"sync/atomic"
)

// Spinlock represents a spin lock.
type Spinlock struct {
	lock uint32
}

// Lock locks the Spinlock.
func (lock *Spinlock) Lock() {
	for !lock.TryLock() {
		runtime.Gosched()
	}
}

// TryLock tries to lock the Spinlock.
func (lock *Spinlock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&lock.lock, 0, 1)
}

// Unlock unlocks the Spinlock.
func (lock *Spinlock) Unlock() {
	atomic.StoreUint32(&lock.lock, 0)
}
