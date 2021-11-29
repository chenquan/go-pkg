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

package xtask

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDoWithPanic(t *testing.T) {

	assert.Panics(t, func() {
		_ = DoWithTimeout(time.Second, func() (err error) {
			panic("")

			return nil
		})
	})
}

func TestDoWithTimeout(t *testing.T) {
	assert.Equal(t, context.DeadlineExceeded, DoWithTimeout(time.Millisecond, func() error {
		time.Sleep(time.Millisecond * 50)
		return nil
	}))
}

func TestDoWithoutTimeout(t *testing.T) {
	assert.Nil(t, DoWithTimeout(time.Second, func() error {
		return nil
	}))
}

func TestDoWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()
	err := Do(ctx, func() error {
		time.Sleep(time.Minute)
		return errors.New("err")
	})
	assert.Equal(t, context.Canceled, err)
}
