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
