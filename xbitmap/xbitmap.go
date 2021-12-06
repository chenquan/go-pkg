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

package xbitmap

import "errors"

var ErrOutOfRange = errors.New("out of range")

type BitMap struct {
	store []uint8
}

func NewBitmap(n int) *BitMap {
	return &BitMap{
		store: make([]uint8, n),
	}
}

func (b *BitMap) coordinate(n uint8) (int, uint8) {
	return int(n / 8), uint8(1 << (n & (8 - 1)))
}

func (b *BitMap) Add(n uint8) error {
	index, position := b.coordinate(n)
	if index >= len(b.store) {
		return ErrOutOfRange
	}
	b.store[index] |= position
	return nil
}

func (b *BitMap) Exist(n uint8) bool {
	index, position := b.coordinate(n)
	if index >= len(b.store) {
		return false
	}
	return b.store[index]&position == 1
}

func (b *BitMap) Del(n uint8) error {
	index, position := b.coordinate(n)
	if index >= len(b.store) {
		return ErrOutOfRange
	}
	b.store[index] &= ^position
	return nil
}
