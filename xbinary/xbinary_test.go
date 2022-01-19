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
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReadBool(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"true",
			args{r: bytes.NewReader([]byte{1})},
			true,
			false,
		}, {
			"false",
			args{r: bytes.NewReader([]byte{0})},
			false,
			false,
		}, {
			"error",
			args{r: bytes.NewReader([]byte{})},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadBool(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type limitWrite struct {
}

func (l *limitWrite) Write(_ []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestWriteUint16(t *testing.T) {
	b := &bytes.Buffer{}

	err := WriteUint16(b, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 1}), b)

	err = WriteUint16(&limitWrite{}, 1)
	assert.Error(t, err)

}

func TestReadUint16(t *testing.T) {
	readUint16, err := ReadUint16(bytes.NewReader([]byte{0, 1}))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, readUint16)

	readUint16, err = ReadUint16(bytes.NewReader([]byte{1}))
	assert.Error(t, err)
	assert.EqualValues(t, 0, readUint16)

	readUint16, err = ReadUint16(bytes.NewReader([]byte{}))
	assert.Error(t, err)
	assert.EqualValues(t, 0, readUint16)
}

func TestWriteBool(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"1",
			args{b: true},
			bytes.NewBuffer([]byte{1}).String(), false,
		}, {
			"1",
			args{b: false},
			bytes.NewBuffer([]byte{0}).String(), false,
		}, {
			"1",
			args{b: false},
			bytes.NewBuffer([]byte{0}).String(), false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := WriteBool(w, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteBool() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestReadUint32(t *testing.T) {
	readUint32, err := ReadUint32(bytes.NewReader([]byte{0, 0, 0, 1}))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, readUint32)

	readUint32, err = ReadUint32(bytes.NewReader([]byte{0, 0, 1}))
	assert.Error(t, err)
	assert.EqualValues(t, 0, readUint32)

	readUint32, err = ReadUint32(bytes.NewReader(nil))
	assert.Error(t, err)
	assert.EqualValues(t, 0, readUint32)
}

func TestWriteUint32(t *testing.T) {
	err := WriteUint32(&limitWrite{}, 1)
	assert.Error(t, err)

	buffer := &bytes.Buffer{}
	err = WriteUint32(buffer, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 0, 0, 1}), buffer)

}

func TestWriteString(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := WriteBytes(buffer, []byte("1"))
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 1, '1'}), buffer)
	err = WriteBytes(&limitWrite{}, []byte(" "))
	assert.Error(t, err)
}

func TestReadString(t *testing.T) {
	readString, err := ReadBytes(bytes.NewBuffer([]byte{0, 1, '1'}))
	assert.NoError(t, err)
	assert.EqualValues(t, "1", readString)

	readString, err = ReadBytes(bytes.NewBuffer([]byte{0, 0}))
	assert.NoError(t, err)
	assert.Len(t, readString, 0)

	readString, err = ReadBytes(bytes.NewBuffer([]byte{0, 2, '1'}))
	assert.Error(t, err)
	assert.Nil(t, readString)

	readString, err = ReadBytes(bytes.NewBuffer([]byte{0}))
	assert.Error(t, err)
	assert.Nil(t, readString)
}
