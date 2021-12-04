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

package xbarrier

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadBarrier(t *testing.T) {
	c := make(chan interface{}, 1)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	reader := NewReadBarrier(ctx, c)
	go func() {
		defer close(c)
		for i := 0; i < 11; i++ {
			c <- i
		}
	}()
	for {
		read, success := reader.Read()
		if read == 9 {
			cancelFunc()
		}
		if read != nil {
			assert.True(t, success)
			assert.LessOrEqual(t, read, 9)
		}

		if !success {
			return
		}
	}

}
