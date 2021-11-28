package xbinary

import (
	"encoding/binary"
	"errors"
	"github.com/chenquan/go-pkg/xbytes"
	"io"
)

var (
	InvalidLengthErr = errors.New("invalid length")

	read1BytesPool = xbytes.GetNBytesPool(1)
	read2BytesPool = xbytes.GetNBytesPool(2)
	read4BytesPool = xbytes.GetNBytesPool(4)
)

func WriteUint16(w io.Writer, i uint16) error {
	data := read2BytesPool.Get()
	data[0] = byte(i >> 8)
	data[1] = byte(i)
	_, err := w.Write(data)
	read2BytesPool.Put(data)
	return err
}

func ReadUint16(r io.Reader) (uint16, error) {
	data := read2BytesPool.Get()
	defer read2BytesPool.Put(data)
	n, err := r.Read(data)
	if err != nil {
		return 0, err
	}
	if n < 2 {
		return 0, InvalidLengthErr
	}
	return binary.BigEndian.Uint16(data), nil
}

//-----------------

func WriteBool(w io.Writer, b bool) error {
	data := read1BytesPool.Get()
	if b {
		data[0] = 1
	} else {
		data[0] = 0
	}
	_, err := w.Write(data)
	read1BytesPool.Put(data)
	return err
}

func ReadBool(r io.Reader) (bool, error) {
	b := read1BytesPool.Get()
	defer read1BytesPool.Put(b)
	_, err := r.Read(b)
	if err != nil {
		return false, err
	}
	return b[0] == 1, nil
}

//------------

func ReadUint32(r io.Reader) (uint32, error) {
	data := read4BytesPool.Get()
	defer read4BytesPool.Put(data)
	n, err := r.Read(data)
	if err != nil {
		return 0, err
	}
	if n < 4 {
		return 0, InvalidLengthErr
	}
	return binary.BigEndian.Uint32(data), nil
}

func WriteUint32(w io.Writer, i uint32) error {
	data := read4BytesPool.Get()
	data[0] = byte(i >> 24)
	data[1] = byte(i >> 16)
	data[2] = byte(i >> 8)
	data[3] = byte(i)
	_, err := w.Write(data)
	read4BytesPool.Put(data)
	return err

}

//------------

func WriteBytes(w io.Writer, s []byte) (err error) {
	// length
	err = WriteUint16(w, uint16(len(s)))
	if err == nil {
		_, err = w.Write(s)
	}
	return
}

func ReadBytes(r io.Reader) (b []byte, err error) {
	nBytes := read2BytesPool.Get()
	defer read2BytesPool.Put(nBytes)
	_, err = io.ReadFull(r, nBytes)
	if err != nil {
		return nil, err
	}

	length := int(binary.BigEndian.Uint16(nBytes))
	if length == 0 {
		return nil, nil
	}
	pool := xbytes.GetNBytesPool(length)
	payload := pool.Get()
	_, err = io.ReadFull(r, payload)
	if err != nil {
		pool.Put(payload)
		return nil, err
	}

	return payload, nil
}
