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

package xballast

import (
	"fmt"
	"sync"
)

// Ballast is a Ballast object.
type Ballast struct {
	ballast     []byte
	ballastLock sync.Mutex
	maxSize     int
}

// NewBallast returns a Ballast.
func NewBallast(maxSize int) *Ballast {
	b := new(Ballast)
	b.maxSize = 1024 * 1024 * 1024 * 2
	if maxSize > 0 {
		b.maxSize = maxSize
	}
	return b
}

// GetSize get the size of ballast object
func (b *Ballast) GetSize() int {
	var sz int
	b.ballastLock.Lock()
	sz = len(b.ballast)
	b.ballastLock.Unlock()
	return sz
}

// SetSize set the size of ballast object
func (b *Ballast) SetSize(newSize int) error {
	if newSize < 0 {
		return fmt.Errorf("newSize cannot be negative: %d", newSize)
	}
	if newSize > b.maxSize {
		return fmt.Errorf("newSize cannot be bigger than %d but it has value %d", b.maxSize, newSize)
	}
	b.ballastLock.Lock()
	b.ballast = make([]byte, newSize)
	b.ballastLock.Unlock()
	return nil
}
