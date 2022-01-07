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

type (
	// SharedBlockMap is an alias that describes the block map.
	SharedBlockMap = Map
	SharedMap      struct {
		b              []*SharedBlockMap
		shardBlockSize uint32
		n              uint32
	}
	sharedMapOptions struct {
		shardBlockSize uint32
	}
	SharedMapOption func(*sharedMapOptions)
)

func WithShardBlockSize(shardBlockSize int) SharedMapOption {
	return func(sharedMapOptions *sharedMapOptions) {
		sharedMapOptions.shardBlockSize = uint32(shardBlockSize)
	}
}

// NewSharedMap returns a SharedMap.
func NewSharedMap(opts ...SharedMapOption) *SharedMap {
	options := new(sharedMapOptions)
	options.shardBlockSize = 32
	for _, opt := range opts {
		opt(options)
	}

	if options.shardBlockSize <= 0 {
		panic("error")
	}
	blockSize := getShardBlockSize(options.shardBlockSize)
	b := make([]*SharedBlockMap, blockSize)
	for i := uint32(0); i < blockSize; i++ {
		b[i] = &Map{}
	}
	return &SharedMap{
		b:              b,
		n:              blockSize - 1,
		shardBlockSize: blockSize,
	}
}

// ComputeIfAbsent  if the value corresponding to the key does not exist,
// use the recalculated value obtained by remappingFunction and save it as the value of the key,
// otherwise return the value.
func (m *SharedMap) ComputeIfAbsent(key string, computeFunc func(key string) interface{}) (actual interface{}, loaded bool) {
	shard := m.GetShard(key)
	return shard.ComputeIfAbsent(key, func(key interface{}) interface{} {
		return computeFunc(key.(string))
	})
}

// ComputeIfPresent if the value corresponding to the key does not exist,
// the null is returned, and if it exists, the value recalculated by remappingFunction is returned.
func (m *SharedMap) ComputeIfPresent(key string, computeFunc func(key string, value interface{}) interface{}) (actual interface{}, exist bool) {
	shard := m.GetShard(key)
	return shard.ComputeIfPresent(key, func(key, value interface{}) interface{} {
		return computeFunc(key.(string), value)
	})
}

func getShardBlockSize(shardBlockSize uint32) uint32 {
	shardBlockSize = shardBlockSize - 1
	shardBlockSize |= shardBlockSize >> 1
	shardBlockSize |= shardBlockSize >> 2
	shardBlockSize |= shardBlockSize >> 4
	shardBlockSize |= shardBlockSize >> 8
	shardBlockSize |= shardBlockSize >> 16
	return shardBlockSize + 1
}

// GetShard returns shard under given key.
func (m SharedMap) GetShard(key string) *SharedBlockMap {
	return m.b[fnv32(key)&m.n]
}

// MStore sets multiple keys and values.
func (m SharedMap) MStore(data map[string]interface{}) {
	for key, value := range data {
		m.Store(key, value)
	}
}

// Store sets the value for a key.
func (m SharedMap) Store(key string, value interface{}) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Store(key, value)
}

// LoadOrStore the given value under the specified key if no value was associated with it.
func (m *SharedMap) LoadOrStore(key string, value interface{}) bool {
	// Get map shard.
	shard := m.GetShard(key)
	_, loaded := shard.LoadOrStore(key, value)
	return !loaded
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *SharedMap) Load(key string) (interface{}, bool) {
	// Get shard
	shard := m.GetShard(key)
	return shard.Load(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *SharedMap) Range(fn func(key, value interface{}) bool) {

	for i := uint32(0); i < m.shardBlockSize; i++ {
		shard := m.b[i]
		b := true
		shard.Range(func(key, value interface{}) bool {
			b = fn(key, value)
			return b
		})
		if !b {
			break
		}
	}
}

// Has looks up an item under specified key.
func (m *SharedMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	_, ok := shard.Load(key)
	return ok
}

// Delete deletes an element from the map.
func (m SharedMap) Delete(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Delete(key)
}

// Clear removes all items from map.
func (m SharedMap) Clear() {
	for _, b := range m.b {
		b.Clear()
	}
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
