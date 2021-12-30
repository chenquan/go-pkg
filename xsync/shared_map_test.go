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
	"hash/fnv"
	"sort"
	"strconv"
	"testing"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("map is null.")
	}

}

//
func TestInsert(t *testing.T) {
	m := New()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.Store("elephant", elephant)
	m.Store("monkey", monkey)
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

//
func TestInsertAbsent(t *testing.T) {
	m := New()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.LoadOrStore("elephant", elephant)
	if ok := m.LoadOrStore("elephant", monkey); ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {
	m := New()

	// Get a missing element.
	val, ok := m.Load("Money")

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal{"elephant"}
	m.Store("elephant", elephant)

	// Retrieve inserted element.
	tmp, ok := m.Load("elephant")
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	elephant, ok = tmp.(Animal) // Type assertion.
	if !ok {
		t.Error("expecting an element, not null.")
	}

	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {
	m := New()

	// Get a missing element.
	if m.Has("Money") == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal{"elephant"}
	m.Store("elephant", elephant)

	if m.Has("elephant") == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	m := New()

	monkey := Animal{"monkey"}
	m.Store("monkey", monkey)

	m.Delete("monkey")

	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Load("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Delete("noone")
}

func TestClear(t *testing.T) {
	m := New()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Store(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	m.Clear()
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Error("We should have 0 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	m := New()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Store(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Load(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Store(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Load(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])
	count := 0
	// Make sure map contains 1000 elements.
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestMInsert(t *testing.T) {
	animals := map[string]interface{}{
		"elephant": Animal{"elephant"},
		"monkey":   Animal{"monkey"},
	}
	m := New()
	m.MStore(animals)
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestFnv32(t *testing.T) {
	key := []byte("ABC")

	hasher := fnv.New32()
	_, err := hasher.Write(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if fnv32(string(key)) != hasher.Sum32() {
		t.Errorf("Bundled fnv32 produced %d, expected result from hash/fnv32 is %d", fnv32(string(key)), hasher.Sum32())
	}

}
