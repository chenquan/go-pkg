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

type (
	Bitmap struct {
		values []byte
		size   uint64
	}
)

//New 初始化一个Bitmap
func New(size uint64) *Bitmap {
	if remainder := size % 8; remainder != 0 {
		size += 8 - remainder
	}
	return &Bitmap{size: size, values: make([]byte, size>>3)}
}
