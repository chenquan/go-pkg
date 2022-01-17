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

package xbytes

import (
	"github.com/chenquan/go-pkg/xsync"
	"golang.org/x/sync/singleflight"
	"strconv"
	"sync"
)

var (
	bytesPoolMap = &xsync.Map{}
	singleFlight = &singleflight.Group{}
)

// Malloc returns a bytes of slice.
func Malloc(size, capacity int) []byte {
	c := size

	if capacity > size {
		c = capacity
	}

	pool := getOrCreatePool(c)
	s := pool.Get()
	data := *((s).(*[]byte))
	//return [:size]
	return data[:size]
}

// MallocSize returns a bytes of slice.
func MallocSize(size int) []byte {
	return Malloc(size, size)
}

func getOrCreatePool(c int) *sync.Pool {
	pool, _, _ := singleFlight.Do(strconv.Itoa(c), func() (interface{}, error) {
		actual, _ := bytesPoolMap.ComputeIfAbsent(c, func(key interface{}) interface{} {
			p := &sync.Pool{New: func() interface{} {
				s := make([]byte, 0, c)
				return &s
			}}
			return p
		})
		return actual, nil
	})
	return pool.(*sync.Pool)
}

// Free recovers a bytes of slice.
func Free(buf []byte) {
	c := cap(buf)
	pool := getOrCreatePool(c)
	s := buf[:0]
	pool.Put(&s)
}
