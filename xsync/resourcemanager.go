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
	"github.com/chenquan/go-pkg/xerror"
	"golang.org/x/sync/singleflight"
	"io"
	"sync"
)

// ResourceManager is a resource manager.
type ResourceManager struct {
	rw           sync.RWMutex
	resources    map[string]io.Closer
	singleFlight singleflight.Group
}

// NewResourceManager returns a ResourceManager.
func NewResourceManager() *ResourceManager {
	return &ResourceManager{resources: map[string]io.Closer{}}
}

// Close the manager.
// Don't use the ResourceManager after Close() called.
func (m *ResourceManager) Close() error {
	m.rw.Lock()

	var be xerror.BatchError
	for _, resource := range m.resources {
		if err := resource.Close(); err != nil {
			be.Add(err)
		}
	}
	m.resources = nil

	m.rw.Unlock()
	return be.Err()
}

// Get returns the resource associated with given key.
func (m *ResourceManager) Get(key string, create func() (io.Closer, error)) (io.Closer, error) {
	val, err, _ := m.singleFlight.Do(key, func() (interface{}, error) {

		m.rw.RLock()
		resource, ok := m.resources[key]
		m.rw.RUnlock()
		if ok {
			return resource, nil
		}

		resource, err := create()
		if err != nil {
			return nil, err
		}

		m.rw.Lock()
		m.resources[key] = resource
		m.rw.Unlock()

		return resource, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(io.Closer), nil
}

// Remove the resource associated with given key and return it if existed.
func (m *ResourceManager) Remove(key string) (exist bool) {
	m.rw.Lock()
	if _, exist = m.resources[key]; exist {
		delete(m.resources, key)
	}
	m.rw.Unlock()
	return
}
