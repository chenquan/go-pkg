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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBallast(t *testing.T) {
	ballast := NewBallast(100)
	assert.Equal(t, ballast.GetSize(), 0)

	err := ballast.SetSize(100)
	assert.NoError(t, err)
	assert.Equal(t, ballast.GetSize(), 100)

	err = ballast.SetSize(1000)
	assert.Error(t, err)
	assert.Equal(t, ballast.GetSize(), 100)

	err = ballast.SetSize(-1)
	assert.Error(t, err)

}
