package xio

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetBufferReaderSize(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	for i := 0; i < 1000; i++ {
		reader := GetBufferReaderSize(buffer, buffer.Len())
		PutBufferReader(reader)
	}

}

func BenchmarkGetBufferReaderSize(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := GetBufferReaderSize(buffer, buffer.Len())
			PutBufferReader(reader)
		}
	})
}
func BenchmarkBufferReaderSize(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bufio.NewReaderSize(buffer, buffer.Len())
		}
	})
}

func BenchmarkGetBufferReader(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := GetBufferReader(buffer)
			PutBufferReader(reader)
		}
	})
}
func BenchmarkBufferReader(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bufio.NewReader(buffer)
		}
	})
}

func TestGetBufferWriter(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	for i := 0; i < 1000; i++ {
		writer := GetBufferWriter(buffer)
		assert.Equal(t, 0, writer.Buffered())
		assert.Equal(t, 4096, writer.Size())
		_ = writer.WriteByte(1)
		assert.Equal(t, 1, writer.Buffered())
		assert.Equal(t, 4095, writer.Available())
		PutBufWriter(writer)
	}

}

func TestGetBufferWriterSize(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	for i := 0; i < 1000; i++ {
		writer := GetBufferWriterSize(buffer, 1000)
		assert.Equal(t, 0, writer.Buffered())
		_ = writer.WriteByte(1)
		assert.Equal(t, 1, writer.Buffered())
		PutBufWriter(writer)
	}
}

func BenchmarkGetBufferWriter(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := GetBufferWriter(buffer)
			_ = writer.WriteByte(1)
			PutBufWriter(writer)
		}
	})
}

func BenchmarkGetBufferWriterSize(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := GetBufferWriterSize(buffer, 1000)
			_ = writer.WriteByte(1)
			PutBufWriter(writer)
		}
	})
}

func BenchmarkBufferWriterSize(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := bufio.NewWriterSize(buffer, 1000)
			_ = writer.WriteByte(1)
		}
	})
}

func BenchmarkBufferWriter(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 1000))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := bufio.NewWriter(buffer)
			_ = writer.WriteByte(1)
		}
	})
}

func TestGetBufferReader(t *testing.T) {
	b := make([]byte, 1000)
	buffer := bytes.NewBuffer(b)
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go func(k int) {
			defer waitGroup.Done()

			reader := GetBufferReader(buffer)
			assert.NotNil(t, reader)
			PutBufferReader(reader)
		}(i)
	}

	waitGroup.Wait()
}
