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

// GetNBytesPool returns a bytes sync.Pool.
// It is recommended to use n Byte greater than or equal to 64.
func GetNBytesPool(nBytes int) *Pool {
	if nBytes < 0 {
		panic("nBytes must be greater than or equal to 0")
	}
	pool, _, _ := singleFlight.Do(strconv.Itoa(nBytes), func() (interface{}, error) {
		actual, _ := bytesPoolMap.ComputeIfAbsent(nBytes, func(key interface{}) interface{} {
			return &Pool{n: nBytes, pool: &sync.Pool{New: func() interface{} {
				return make([]byte, nBytes)
			}}}
		})
		return actual, nil
	})
	return pool.(*Pool)
}

// Pool Represents a pool of the same byte size.
type Pool struct {
	n    int
	pool *sync.Pool
}

// Get Returns a bytes of slice.
func (p *Pool) Get() (bytes []byte) {
	return p.pool.Get().([]byte)
}

// Put Recovers a byte slice.
func (p *Pool) Put(b []byte) {
	if len(b) >= p.n {
		p.pool.Put(b[:p.n])
	}
}
