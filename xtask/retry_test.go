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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {

	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("")
	}))

	var times int
	assert.Nil(t, DoWithRetry(func() error {
		times++
		if times == defaultRetryTimes {
			return nil
		}
		return errors.New("")
	}))

	times = -1
	assert.NotNil(t, DoWithRetry(func() error {
		times++
		if times == defaultRetryTimes+1 {
			return nil
		}
		return errors.New("")
	}))

	total := 2 + defaultRetryTimes
	times = 0
	assert.Nil(t, DoWithRetry(func() error {
		times++
		if times == total {
			return nil
		}
		return errors.New("")
	}, WithRetry(total)))
}
