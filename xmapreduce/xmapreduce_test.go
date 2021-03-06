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

package xmapreduce

import (
	"context"
	"github.com/chenquan/go-pkg/xbarrier"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		c := Map(context.Background(), func(source chan<- interface{}) {
			for i := 0; i < 10; i++ {
				source <- i
			}
		}, func(item interface{}, writer xbarrier.Writer) {
			writer.Write(item)
		}, WithWorkerSize(1))
		i := 0
		for range c {
			i++
		}
		assert.Equal(t, 10, i)
	})

	t.Run("cancel", func(t *testing.T) {
		ctx, cancelFunc := context.WithCancel(context.Background())
		c := Map(ctx, func(source chan<- interface{}) {
			for i := 0; i < 11; i++ {
				source <- i
				if i == 9 {
					// Wait for data to be read.
					time.Sleep(time.Second)
					cancelFunc()
				}
			}
		}, func(item interface{}, writer xbarrier.Writer) {
			writer.Write(item)
		}, WithWorkerSize(1))
		i := 0
		for range c {
			i++
		}
		assert.Equal(t, 10, i)
	})

}

func TestMapStream(t *testing.T) {

	ctx, cancelFunc := context.WithCancel(context.Background())
	count := MapStream(ctx,
		func(source chan<- interface{}) {
			for i := 0; i < 10; i++ {
				source <- i
				if i == 3 {
					time.Sleep(time.Second)
					cancelFunc()
				}
			}
		},
		func(item interface{}, writer xbarrier.Writer) {
			i := item.(int)
			writer.Write(i)
		},
		WithWorkerSize(1),
	).Count()

	assert.Equal(t, 4, count)
}
