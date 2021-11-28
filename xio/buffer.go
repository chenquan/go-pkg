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

package xio

import (
	"bufio"
	"io"
	"sync"
)

var (
	bufReaderPool = &sync.Pool{}
	bufWriterPool = &sync.Pool{}
)

// GetBufferReaderSize returns a bufio.Reader.
func GetBufferReaderSize(r io.Reader, size int) *bufio.Reader {
	if v := bufReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReaderSize(r, size)
}

// GetBufferReader returns a bufio.Reader.
func GetBufferReader(r io.Reader) *bufio.Reader {
	if v := bufReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

// PutBufferReader recycles a bufio.Reader.
func PutBufferReader(r *bufio.Reader) {
	r.Reset(nil)
	bufReaderPool.Put(r)
}

// -----------------

// GetBufferWriterSize returns a bufio.Writer.
func GetBufferWriterSize(w io.Writer, size int) *bufio.Writer {
	if v := bufWriterPool.Get(); v != nil {
		bw := v.(*bufio.Writer)
		bw.Reset(w)
		return bw
	}
	return bufio.NewWriterSize(w, size)
}

// GetBufferWriter returns a bufio.Writer.
func GetBufferWriter(w io.Writer) *bufio.Writer {
	if v := bufWriterPool.Get(); v != nil {
		bw := v.(*bufio.Writer)
		bw.Reset(w)
		return bw
	}
	return bufio.NewWriter(w)
}

// PutBufWriter recycles a bufio.Writer.
func PutBufWriter(w *bufio.Writer) {
	w.Reset(nil)
	bufWriterPool.Put(w)
}
