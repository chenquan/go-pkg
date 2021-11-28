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
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type dummyResource struct {
	age int
}

func (dr *dummyResource) Close() error {
	return errors.New("close")
}

func TestResourceManager_GetResource(t *testing.T) {
	manager := NewResourceManager()
	defer func() {
		_ = manager.Close()
	}()

	var age int
	for i := 0; i < 10; i++ {
		val, err := manager.Get("key", func() (io.Closer, error) {
			age++
			return &dummyResource{
				age: age,
			}, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, val.(*dummyResource).age)
	}
}

func TestResourceManager_GetResourceError(t *testing.T) {
	manager := NewResourceManager()
	defer func() {
		_ = manager.Close()
	}()

	for i := 0; i < 10; i++ {
		_, err := manager.Get("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)
	}
}

func TestResourceManager_Close(t *testing.T) {
	manager := NewResourceManager()
	for i := 0; i < 10; i++ {
		_, err := manager.Get("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)
	}

	if assert.NoError(t, manager.Close()) {
		assert.Equal(t, 0, len(manager.resources))
	}
}

func TestResourceManager_UseAfterClose(t *testing.T) {
	manager := NewResourceManager()
	_, err := manager.Get("key", func() (io.Closer, error) {
		return nil, errors.New("fail")
	})
	assert.NotNil(t, err)
	if assert.NoError(t, manager.Close()) {
		_, err = manager.Get("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)
	}
}

func TestResourceManager_Remove(t *testing.T) {
	manager := NewResourceManager()
	closer, err := manager.Get("key", func() (io.Closer, error) {
		return &dummyResource{}, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, closer)
	assert.True(t, manager.Remove("key"))
}
