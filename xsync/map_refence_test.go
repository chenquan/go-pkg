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
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsync

import (
	"sync"
	"sync/atomic"
)

// This file contains reference map implementations for unit-tests.

// mapInterface is the interface Map implements.
type mapInterface interface {
	Load(interface{}) (interface{}, bool)
	Store(key, value interface{})
	LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
	LoadAndDelete(key interface{}) (value interface{}, loaded bool)
	Delete(interface{})
	Range(func(key, value interface{}) (shouldContinue bool))
	ComputeIfAbsent(key interface{}, computeFunc func(key interface{}) interface{}) (actual interface{}, loaded bool)
	ComputeIfPresent(key interface{}, computeFunc func(key, value interface{}) interface{}) (actual interface{}, exist bool)
}

var _ mapInterface = (*RWMutexMap)(nil)

// RWMutexMap is an implementation of mapInterface using a sync.RWMutex.
type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[interface{}]interface{}
}

func (m *RWMutexMap) ComputeIfAbsent(key interface{}, computeFunc func(key interface{}) interface{}) (actual interface{}, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.dirty[key]; ok {
		return v, true
	}
	if m.dirty == nil {
		m.dirty = make(map[interface{}]interface{})
	}
	m.dirty[key] = computeFunc(key)
	return m.dirty[key], false
}

func (m *RWMutexMap) ComputeIfPresent(key interface{}, computeFunc func(key interface{}, value interface{}) interface{}) (actual interface{}, exist bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if value, ok := m.dirty[key]; !ok {
		return nil, false
	} else {

		if m.dirty == nil {
			m.dirty = make(map[interface{}]interface{})
		}
		m.dirty[key] = computeFunc(key, value)
		return m.dirty[key], false
	}

}

func (m *RWMutexMap) Load(key interface{}) (value interface{}, ok bool) {
	m.mu.RLock()
	value, ok = m.dirty[key]
	m.mu.RUnlock()
	return
}

func (m *RWMutexMap) Store(key, value interface{}) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[interface{}]interface{})
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[interface{}]interface{})
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *RWMutexMap) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	m.mu.Lock()
	value, loaded = m.dirty[key]
	if !loaded {
		m.mu.Unlock()
		return nil, false
	}
	delete(m.dirty, key)
	m.mu.Unlock()
	return value, loaded
}

func (m *RWMutexMap) Delete(key interface{}) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

func (m *RWMutexMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	m.mu.RLock()
	keys := make([]interface{}, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

var _ mapInterface = (*DeepCopyMap)(nil)

// DeepCopyMap is an implementation of mapInterface using a Mutex and
// atomic.Value.  It makes deep copies of the map on every write to avoid
// acquiring the Mutex in Load.
type DeepCopyMap struct {
	mu    sync.Mutex
	clean atomic.Value
}

func (m *DeepCopyMap) ComputeIfAbsent(key interface{}, computeFunc func(key interface{}) interface{}) (actual interface{}, loaded bool) {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	if value, ok := clean[key]; ok {
		return value, ok
	}
	m.mu.Lock()
	// Reload clean in case it changed while we were waiting on m.mu.
	clean, _ = m.clean.Load().(map[interface{}]interface{})
	actual, loaded = clean[key]
	if !loaded {
		dirty := m.dirty()
		actual = computeFunc(key)
		dirty[key] = actual
		m.clean.Store(dirty)
	}
	m.mu.Unlock()
	return actual, false
}

func (m *DeepCopyMap) ComputeIfPresent(key interface{}, computeFunc func(key interface{}, value interface{}) interface{}) (actual interface{}, exist bool) {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	if value, ok := clean[key]; !ok {
		return nil, false
	} else {
		m.mu.Lock()
		// Reload clean in case it changed while we were waiting on m.mu.
		clean, _ = m.clean.Load().(map[interface{}]interface{})
		actual, ok = clean[key]
		if ok {
			dirty := m.dirty()
			actual = computeFunc(key, value)
			dirty[key] = actual
			exist = true
			m.clean.Store(dirty)
		}
		m.mu.Unlock()
		return

	}

}

func (m *DeepCopyMap) Load(key interface{}) (value interface{}, ok bool) {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	value, ok = clean[key]
	return value, ok
}

func (m *DeepCopyMap) Store(key, value interface{}) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	actual, loaded = clean[key]
	if loaded {
		return actual, loaded
	}

	m.mu.Lock()
	// Reload clean in case it changed while we were waiting on m.mu.
	clean, _ = m.clean.Load().(map[interface{}]interface{})
	actual, loaded = clean[key]
	if !loaded {
		dirty := m.dirty()
		dirty[key] = value
		actual = value
		m.clean.Store(dirty)
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *DeepCopyMap) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	m.mu.Lock()
	dirty := m.dirty()
	value, loaded = dirty[key]
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
	return
}

func (m *DeepCopyMap) Delete(key interface{}) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	for k, v := range clean {
		if !f(k, v) {
			break
		}
	}
}

func (m *DeepCopyMap) dirty() map[interface{}]interface{} {
	clean, _ := m.clean.Load().(map[interface{}]interface{})
	dirty := make(map[interface{}]interface{}, len(clean)+1)
	for k, v := range clean {
		dirty[k] = v
	}
	return dirty
}
