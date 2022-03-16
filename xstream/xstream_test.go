/*
 *
 *     Copyright 2020 chenquan
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */

package xstream

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func equal[T any](t *testing.T, stream *Stream[T], data []T) {
	items := make([]T, 0)
	for item := range stream.source {
		items = append(items, item)
	}
	if !reflect.DeepEqual(items, data) {
		t.Errorf(" %v, want %v", items, data)
	}
}

func assertEqual(t *testing.T, except any, data any) {
	if !reflect.DeepEqual(except, data) {
		t.Errorf(" %v, want %v", data, except)
	}

}

func TestEmpty(t *testing.T) {
	empty := Empty[int]()
	assertEqual(t, len(empty.source), 0)
	assertEqual(t, cap(empty.source), 0)

}

//
// func TestRange(t *testing.T) {
//	stream1 := Range(make(chan any))
//	assertEqual(t, len(stream1.source), 0)
//
//	stream2 := Range(make(chan any, 2))
//	assertEqual(t, len(stream2.source), 0)
//	assertEqual(t, cap(stream2.source), 2)
// }
//
func TestOf(t *testing.T) {
	ints := []int{1, 2, 3, 4}
	of := Of(ints...).Sort(func(a, b int) bool {
		return a < b
	})
	var items []int
	for item := range of.source {
		items = append(items, item)
	}
	assertEqual(t, items, ints)
}

func TestConcat(t *testing.T) {
	a1 := []int{1, 2, 3}
	a2 := []int{4, 5, 6}
	s1 := Of(a1...)
	s2 := Of(a2...)
	stream := Concat(s1, s2)
	var items []int
	for item := range stream.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
	ints := make([]int, 0)
	ints = append(ints, a1...)
	ints = append(ints, a2...)
	assertEqual(t, ints, items)

	of := Of(1)
	equal(t, of.Concat(of), []int{1})

}

func TestFrom(t *testing.T) {
	ints := make([]int, 0)
	stream := From(func(source chan<- int) {
		for i := 0; i < 10; i++ {
			source <- i
			ints = append(ints, i)
		}
	})
	items := make([]int, 0)
	for item := range stream.source {
		items = append(items, item)
	}
	assertEqual(t, items, ints)
}

func TestStream_Distinct(t *testing.T) {
	stream := Of(1, 2, 3, 4, 4, 22, 2, 1, 4).Distinct(func(item int) any {
		return item
	})
	equal(t, stream, []int{1, 2, 3, 4, 22})
}

func TestStream_Count(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 22, 2, 1, 4}
	assertEqual(t, Of(data...).Count(), len(data))
}

func TestStream_Buffer(t *testing.T) {
	stream := Of(1, 2, 4)
	assertEqual(t, cap(stream.source), 3)
	stream = stream.Buffer(10)
	assertEqual(t, cap(stream.source), 10)
	stream = stream.Buffer(-1)
	assertEqual(t, cap(stream.source), 0)

}

//
// func TestStream_Split(t *testing.T) {
//
//	stream := Of(1, 2, 444, 441, 1).Split(3)
//	assertEqual(t, (<-stream.source).([]any), []any{1, 2, 444})
//	assertEqual(t, (<-stream.source).([]any), []any{441, 1})
//	assert.Panics(t, func() {
//		Of(1, 2, 444, 441, 1).Split(-1)
//	})
//	assert.Panics(t, func() {
//		Of(1, 2, 444, 441, 1).Split(0)
//	})
// }
//
// func TestStream_SplitSteam2(t *testing.T) {
//	streams := Of(1, 2, 444, 441, 1).SplitSteam(3)
//
//	equal(t, (<-streams.source).(*Stream), []any{1, 2, 444})
//	equal(t, (<-streams.source).(*Stream), []any{441, 1})
// }

func TestStream_Sort(t *testing.T) {
	ints := []int{4, 2, 1, 441, 23, 14, 1, 23}
	stream := Of(ints...).Sort(func(a, b int) bool {
		return a < b
	})
	sort.Slice(ints, func(i, j int) bool {
		return ints[i] < ints[j]
	})
	equal(t, stream, ints)
}

func TestStream_Tail(t *testing.T) {
	equal(t, Of(1, 232, 3, 2, 3).Tail(1), []int{3})
	equal(t, Of(1, 232, 3, 2, 3).Tail(2), []int{2, 3})
	equal(t, Of(1, 232, 3, 2, 3).Tail(8), []int{1, 232, 3, 2, 3})
}

func TestTailZero(t *testing.T) {
	Of(1, 2, 3, 4).Tail(0).Done()

	assert.Panics(t, func() {
		Of(1, 2, 3, 4).Tail(-1).Done()
	})
}

func TestStream_Skip(t *testing.T) {
	assertEqual(t, 3, Of(1, 2, 3, 4).Skip(1).Count())
	assertEqual(t, 1, Of(1, 2, 3, 4).Skip(3).Count())
	equal(t, Of(1, 2, 3, 4).Skip(3), []int{4})
	equal(t, Of(1, 2, 3).Skip(0), []int{1, 2, 3})
	assert.Panics(t, func() {
		Of(1, 2, 3).Skip(-1)
	})

}

func TestStream_Limit(t *testing.T) {

	equal(t, Of(1, 2, 3, 4).Limit(3), []int{1, 2, 3})
	equal(t, Of(1, 2, 3, 4).Limit(4), []int{1, 2, 3, 4})
	equal(t, Of(1, 2, 3, 4).Limit(5), []int{1, 2, 3, 4})
	equal(t, Of(1, 2, 3, 4).Limit(0), []int{})
	assert.Panics(t, func() {
		Of(1, 2, 3, 4).Limit(-1)
	})

}

func TestStream_Foreach(t *testing.T) {
	var items []any
	Of(1, 2, 3, 4).Foreach(func(item int) {
		items = append(items, item)
	})
	assertEqual(t, []any{1, 2, 3, 4}, items)
}

func TestStream_ForeachOrdered(t *testing.T) {
	var items []any
	Of(1, 2, 3, 4).ForeachOrdered(func(item int) {
		items = append(items, item)
	})
	assertEqual(t, []any{4, 3, 2, 1}, items)
}

func TestStream_Concat(t *testing.T) {
	stream := Of(1).Concat(Of(2), Of(3))
	var items []any
	for item := range stream.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].(int) < items[j].(int)
	})
	assertEqual(t, []any{1, 2, 3}, items)
}

func TestStream_Filter(t *testing.T) {
	equal(t, Of(1, 2, 3, 4).Filter(func(item int) bool {
		return item > 3
	}), []int{4})
	equal(t, Of(1, 2, 3, 4).Filter(func(item int) bool {
		return item > 2
	}).Sort(func(a, b int) bool {
		return a < b
	}), []int{3, 4})
}

func TestStream_Map(t *testing.T) {
	equal(t, Of(1, 2, 3).Map(func(item int) int {
		return item + 1
	}).Sort(func(a, b int) bool {
		return a < b
	}), []int{2, 3, 4})
}

// func TestStream_FlatMap(t *testing.T) {
//	equal(t,
//		Of([]any{1, 2}, []any{3, 4}).FlatMap(func(item any) any {
//			return item
//		}).Sort(func(a, b any) bool {
//			return a.(int) < b.(int)
//		}),
//		[]any{1, 2, 3, 4},
//	)
//	equal(t,
//		Of(1, 2, 3, 4, 5, 6).FlatMap(func(item any) any {
//			return item
//		}).Sort(func(a, b any) bool {
//			return a.(int) < b.(int)
//		}),
//		[]any{1, 2, 3, 4, 5, 6},
//	)
// }
//
// func TestStream_Group(t *testing.T) {
//	equal(t,
//		Of(1, 2, 3, 4).Group(func(item any) any {
//			return item.(int) % 2
//		}).Map(func(item any) any {
//			return item.([]any)
//		}).FlatMap(func(item any) any {
//			return item
//		}).Sort(func(a, b any) bool {
//			return a.(int) < b.(int)
//		}),
//		[]any{1, 2, 3, 4},
//	)
// }
//
func TestStream_Merge(t *testing.T) {
	// equal(t, Of(1, 2, 3, 4).Merge(), []int{
	//
	//	[]int{1, 2, 3, 4},
	// })
	// Of(1, 2, 3, 4).Merge().Done()
}

func TestStream_Reverse(t *testing.T) {
	equal(t, Of(1, 2, 3, 4, 1).Reverse(), []int{1, 4, 3, 2, 1})
}

func TestStream_ParallelFinish(t *testing.T) {

	Of(1, 23).ParallelFinish(func(item int) {

	}, WithWorkSize(2))
}

func TestStream_AnyMach(t *testing.T) {
	assertEqual(t, false, Of(1, 2, 3).AnyMach(func(item int) bool {
		return item == 4
	}))
	assertEqual(t, true, Of(1, 2, 3).AnyMach(func(item int) bool {
		return item == 2
	}))
}

func TestStream_AllMach(t *testing.T) {
	assertEqual(
		t, true, Of(1, 2, 3).AllMach(func(item int) bool {
			return true
		}),
	)
	assertEqual(
		t, false, Of(1, 2, 3).AllMach(func(item int) bool {
			return false
		}),
	)
	assertEqual(
		t, false, Of(1, 2, 3).AllMach(func(item int) bool {
			return item == 1
		}),
	)
}

func TestStream_Chan(t *testing.T) {
	var items []int

	for item := range Of(1, 2, 3).Chan() {
		items = append(items, item)
	}
	assertEqual(t, items, []int{1, 2, 3})
}

// func TestStream_SplitSteam(t *testing.T) {
//	streams := Of(1, 2, 444, 441, 1).SplitSteam(3)
//	equal(t, (<-streams.source).(*Stream), []any{1, 2, 444})
//	equal(t, (<-streams.source).(*Stream), []any{441, 1})
//	assert.Panics(t, func() {
//		Of(1, 2, 444, 441, 1).SplitSteam(-1)
//	})
// }
//
func TestStream_Peek(t *testing.T) {
	items := make([]int, 0)
	Of(1, 2, 3, 4).Peek(func(item int) {
		items = append(items, item)
	}).Done()
	assertEqual(t, items, []int{1, 2, 3, 4})
}

func TestStream_FindFirst(t *testing.T) {
	result, err := Of(1, 2, 3).FindFirst()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, result)

	result, err = Empty[int]().FindFirst()
	assert.Error(t, err)
	assert.Equal(t, 0, result)
}

// func TestStream_Copy(t *testing.T) {
//	stream := Of(1, 2, 3)
//	s1 := stream.Copy()
//	assert.Equal(t, 3, s1.Count())
//	assert.Equal(t, 3, stream.Count())
// }
//
// func TestStream_Collection(t *testing.T) {
//
//	t.Run("GroupBy", func(t *testing.T) {
//		stream := Of(1, 2, 3)
//		group := GroupBy(func(item any) any {
//			return item.(int) % 2
//		})
//		stream.Collection(group)
//
//		assert.EqualValues(t, map[any][]any{
//			0: {2},
//			1: {1, 3},
//		}, group.Map())
//	})
//
//	t.Run("CollectorFunc", func(t *testing.T) {
//		ints := make([]int, 0, 3)
//		Of(1, 2, 3).Collection(CollectorFunc(func(c <-chan any) {
//			for i := range c {
//				ints = append(ints, i.(int))
//			}
//		}))
//		assert.EqualValues(t, []int{1, 2, 3}, ints)
//	})
// }
//

func TestStream_FindLast(t *testing.T) {
	t.Run("has value", func(t *testing.T) {
		last, err := Of(1, 2, 3).FindLast()
		assert.NoError(t, err)
		assert.EqualValues(t, 3, last)
	})

	t.Run("hasn't value", func(t *testing.T) {
		last, err := Empty[int]().FindLast()
		assert.Error(t, err)
		assert.EqualValues(t, 0, last)
	})
}

func TestTo(t *testing.T) {
	To[int](Of[interface{}](1, 2)).Foreach(func(item int) {
		fmt.Println(item)
	})
}
