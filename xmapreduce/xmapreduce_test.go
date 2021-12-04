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
	"fmt"
	"github.com/chenquan/go-pkg/xstream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGuardedWriter(t *testing.T) {
	c := make(chan interface{}, 1)
	ctx, cancelFunc := context.WithCancel(context.Background())
	writer := NewGuardedWriter(c, ctx)
	go func() {
		for i := 0; i < 11; i++ {
			writer.Write(1)
			if i == 9 {
				cancelFunc()
			}
		}
	}()

	idx := 0
	for i := 0; i < 11; i++ {
		select {
		case v := <-c:
			idx += v.(int)
		default:

		}
	}
	assert.Equal(t, 10, idx)
}

func TestMap(t *testing.T) {
	c := Map(context.Background(), func(source chan<- interface{}) {
		for i := 0; i < 10; i++ {
			source <- i
		}
	}, func(item interface{}, writer Writer) {
		i := item.(int)
		for j := 0; j < i; j++ {
			writer.Write(j)
		}
		//writer.Write(i)
	}, WithWorkerSize(1))

	fmt.Println(xstream.Range(c).Count())

}

func TestMapStream(t *testing.T) {
	//MapStream(context.Background(),
	//	func(source chan<- interface{}) {
	//		for i := 0; i < 10; i++ {
	//			source <- i
	//		}
	//	},
	//	func(item interface{}, writer Writer) {
	//		i := item.(int)
	//		//for j := 0; j < i; j++ {
	//		//	writer.Write(j)
	//		//}
	//		writer.Write(i)
	//	},
	//	WithWorkerSize(1),
	//).ForeachOrdered(func(item interface{}) {
	//	fmt.Println(item)
	//})
	//
	ctx, cancelFunc := context.WithCancel(context.Background())
	MapStream(ctx,
		func(source chan<- interface{}) {
			for i := 0; i < 10; i++ {
				source <- i
				if i == 3 {
					cancelFunc()
				}
			}
		},
		func(item interface{}, writer Writer) {
			i := item.(int)
			//for j := 0; j < i; j++ {
			//	writer.Write(j)
			//}
			writer.Write(i)
		},
		WithWorkerSize(1),
	).ForeachOrdered(func(item interface{}) {
		fmt.Println(item)
	})

}
