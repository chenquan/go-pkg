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
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetNBytesPool(t *testing.T) {
	m := sync.Map{}
	n := 100
	waitGroup := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func(j int) {
			defer waitGroup.Done()
			k := j % 10
			t.Run(fmt.Sprintf("%d bytes", k), func(t *testing.T) {
				pool := GetNBytesPool(k)
				value, _ := m.LoadOrStore(k, atomic.Value{})
				v := value.(atomic.Value)
				swap := v.CompareAndSwap(nil, pool)
				b := pool.Get()
				assert.EqualValues(t, len(b), k)
				if !swap {
					assert.EqualValues(t, v.Load(), pool)
				}
				pool.Put(b)
			})

		}(i)
	}
	waitGroup.Wait()
}
func TestGetNBytesPool_Panic(t *testing.T) {
	assert.Panics(t, func() {
		GetNBytesPool(-1)
	})
}
func BenchmarkMakeBytes(b *testing.B) {
	b.ReportAllocs()
	//n := 20
	rand.Seed(2021)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			waitGroup := sync.WaitGroup{}
			for j := 0; j < 20; j++ {
				b.StopTimer()
				waitGroup.Add(1)
				k := 2 << j
				b.StartTimer()
				go func() {
					a := make([]byte, k)
					_ = a
					b.StopTimer()
					waitGroup.Done()
					b.StartTimer()
				}()
			}
			waitGroup.Wait()

		}
	})

}

func BenchmarkGetNBytesPool(b *testing.B) {
	//n := 20
	rand.Seed(2021)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			waitGroup := sync.WaitGroup{}
			for j := 0; j < 20; j++ {
				b.StopTimer()
				k := 2 << j
				waitGroup.Add(1)
				b.StartTimer()
				pool := GetNBytesPool(k)
				go func() {
					a := pool.Get()
					b.StopTimer()
					pool.Put(a)
					waitGroup.Done()
					b.StartTimer()
				}()
			}
			waitGroup.Wait()

		}
	})

}
