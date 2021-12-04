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

package xstring

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoiner_WriteString(t *testing.T) {
	join := NewJoiner(WithJoin("(", ",", ")"))
	_, _ = join.WriteString("1")
	_, _ = join.WriteString("2")
	_, _ = join.WriteString("3")
	assert.Equal(t, "(1,2,3)", join.String())

	join = NewJoiner(WithJoin("(", ",", ""))
	_, _ = join.WriteString("1")
	_, _ = join.WriteString("2")
	_, _ = join.WriteString("3")
	assert.Equal(t, "(1,2,3", join.String())

	join = NewJoiner(WithJoinStep("-"))
	_, _ = join.WriteString("1")
	_, _ = join.WriteString("2")
	_, _ = join.WriteString("3")
	assert.Equal(t, "1-2-3", join.String())

	join = NewJoiner(WithJoinStep("-"), WithJoinPrefix("=>"))
	_, _ = join.WriteString("1")
	_, _ = join.WriteString("2")
	_, _ = join.WriteString("3")
	assert.Equal(t, "=>1-2-3", join.String())

	join = NewJoiner(WithJoinStep("-"), WithJoinPrefix("=>"), WithJoinSuffix("<===="))
	_, _ = join.WriteString("1")
	_, _ = join.WriteString("2")
	_, _ = join.WriteString("3")
	assert.Equal(t, "=>1-2-3<====", join.String())
}

func TestJoiner_Write(t *testing.T) {
	join := NewJoiner(WithJoin("(", ",", ")"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "(a,b,c)", join.String())

	join = NewJoiner(WithJoin("(", ",", ""))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "(a,b,c", join.String())

	join = NewJoiner(WithJoinStep("000"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "a000b000c", join.String())

	join = NewJoiner(WithJoinStep("1"), WithJoinPrefix("--"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "--a1b1c", join.String())

	join = NewJoiner(WithJoinStep("1"), WithJoinPrefix("--"), WithJoinSuffix("<---"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "--a1b1c<---", join.String())

}
func TestJoiner_WriteRune(t *testing.T) {
	join := NewJoiner(WithJoin("(", ",", ")"))
	_, _ = join.WriteRune('a')
	_, _ = join.WriteRune('b')
	_, _ = join.WriteRune('c')
	assert.Equal(t, "(a,b,c)", join.String())

	join = NewJoiner(WithJoin("(", ",", ""))
	_, _ = join.WriteRune('a')
	_, _ = join.WriteRune('b')
	_, _ = join.WriteRune('c')
	assert.Equal(t, "(a,b,c", join.String())

	join = NewJoiner(WithJoinStep("000"))
	_, _ = join.WriteRune('a')
	_, _ = join.WriteRune('b')
	_, _ = join.WriteRune('c')
	assert.Equal(t, "a000b000c", join.String())

	join = NewJoiner(WithJoinStep("1"), WithJoinPrefix("--"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "--a1b1c", join.String())

	join = NewJoiner(WithJoinStep("1"), WithJoinPrefix("--"), WithJoinSuffix("<---"))
	_ = join.WriteByte('a')
	_ = join.WriteByte('b')
	_ = join.WriteByte('c')
	assert.Equal(t, "--a1b1c<---", join.String())

}

func TestJoiner_Len(t *testing.T) {
	join := NewJoiner(WithJoin("(", ",", ")"))
	assert.Equal(t, 2, join.Len())
	_, _ = join.WriteRune('a')
	_, _ = join.WriteRune('b')
	_, _ = join.WriteRune('c')
	assert.Equal(t, 7, join.Len())
	assert.Equal(t, len(join.String()), join.Len())
}

func TestJoiner_Grow(t *testing.T) {
	for _, growLen := range []int{0, 100, 1000, 10000, 100000} {
		p := bytes.Repeat([]byte{'a'}, growLen)
		var b = NewJoiner()
		allocs := testing.AllocsPerRun(100, func() {
			b.Reset()
			b.Grow(growLen) // should be only alloc, when growLen > 0
			if b.Cap() < growLen {
				t.Fatalf("growLen=%d: Cap() is lower than growLen", growLen)
			}
			_, _ = b.Write(p)
			if b.String() != string(p) {
				fmt.Println(b.String(), "  ", string(p))
				fmt.Println(len(b.String()), "  ", len(string(p)))
				t.Fatalf("growLen=%d: bad data written after Grow", growLen)
			}
		})
		wantAllocs := 1
		if growLen == 0 {
			wantAllocs = 0
		}
		if g, w := int(allocs), wantAllocs; g != w {
			t.Errorf("growLen=%d: got %d allocs during Write; want %v", growLen, g, w)
		}
	}
}

func BenchmarkNewJoin(b *testing.B) {
	join := NewJoiner(WithJoin("(", ",", ")"))
	for i := 0; i < b.N; i++ {
		_, _ = join.WriteString("1")
		_, _ = join.WriteString("2")
		_, _ = join.WriteString("3")
	}
}

func TestJoiner_Cap(t *testing.T) {
	join := NewJoiner()
	assert.Equal(t, 0, join.Cap())
	_, _ = join.WriteString("1")
	assert.Equal(t, 8, join.Cap())
}

func TestJoiner_Reset(t *testing.T) {
	join := NewJoiner()
	assert.Equal(t, 0, join.Cap())
	_, _ = join.WriteString("111")
	assert.Equal(t, 8, join.Cap())
	join.Reset()
	assert.Equal(t, 0, join.Cap())
}
