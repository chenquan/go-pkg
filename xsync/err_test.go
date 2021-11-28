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
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	err1 = "error1"
	err2 = "error2"
)

func TestBatchErrorNil(t *testing.T) {
	var batch BatchError
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
	batch.Add(nil)
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
}

func TestBatchErrorNilFromFunc(t *testing.T) {
	err := func() error {
		var be BatchError
		return be.Err()
	}()
	assert.True(t, err == nil)
}

func TestBatchErrorOneError(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	assert.NotNil(t, batch)
	assert.Equal(t, err1, batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchErrorWithErrors(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	batch.Add(errors.New(err2))
	assert.NotNil(t, batch)
	assert.Equal(t, fmt.Sprintf("%s\n%s", err1, err2), batch.Err().Error())
	assert.True(t, batch.NotNil())
}
