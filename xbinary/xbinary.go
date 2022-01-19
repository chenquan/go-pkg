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

package xbinary

import (
	"encoding/binary"
	"errors"
	"github.com/chenquan/go-pkg/xbytes"
	"io"
)

var (
	// ErrInvalidLength invalid length error.
	ErrInvalidLength = errors.New("invalid length error")
)

// WriteUint16 writes unit16.
func WriteUint16(w io.Writer, i uint16) error {
	data := xbytes.MallocSize(2)

	data[0] = byte(i >> 8)
	data[1] = byte(i)

	_, err := w.Write(data)
	xbytes.Free(data)
	return err
}

// ReadUint16 reads unit16.
func ReadUint16(r io.Reader) (uint16, error) {
	data := xbytes.MallocSize(2)
	defer xbytes.Free(data)

	n, err := r.Read(data)
	if err != nil {
		return 0, err
	}
	if n < 2 {
		return 0, ErrInvalidLength
	}
	return binary.BigEndian.Uint16(data), nil
}

//-----------------

// WriteBool writes bool.
func WriteBool(w io.Writer, b bool) error {
	data := xbytes.MallocSize(1)
	if b {
		data[0] = 1
	} else {
		data[0] = 0
	}

	_, err := w.Write(data)
	xbytes.Free(data)
	return err
}

// ReadBool reads bool.
func ReadBool(r io.Reader) (bool, error) {
	b := xbytes.MallocSize(1)
	defer xbytes.Free(b)

	_, err := r.Read(b)
	if err != nil {
		return false, err
	}
	return b[0] == 1, nil
}

//------------

// ReadUint32 reads unit32.
func ReadUint32(r io.Reader) (uint32, error) {
	data := xbytes.MallocSize(4)
	defer xbytes.Free(data)

	n, err := r.Read(data)
	if err != nil {
		return 0, err
	}
	if n < 4 {
		return 0, ErrInvalidLength
	}
	return binary.BigEndian.Uint32(data), nil
}

// WriteUint32 writes unit32.
func WriteUint32(w io.Writer, i uint32) error {
	data := xbytes.MallocSize(4)

	data[0] = byte(i >> 24)
	data[1] = byte(i >> 16)
	data[2] = byte(i >> 8)
	data[3] = byte(i)

	_, err := w.Write(data)
	xbytes.Free(data)
	return err

}

//------------

// WriteBytes writes bytes.
func WriteBytes(w io.Writer, s []byte) (err error) {
	// length
	err = WriteUint16(w, uint16(len(s)))
	if err == nil {
		_, err = w.Write(s)
	}
	return
}

// ReadBytes reads bytes.
func ReadBytes(r io.Reader) (b []byte, err error) {
	length, err := ReadUint16(r)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}
	payload := xbytes.MallocSize(int(length))
	_, err = io.ReadFull(r, payload)
	if err != nil {
		xbytes.Free(payload)
		return nil, err
	}

	return payload, nil
}
