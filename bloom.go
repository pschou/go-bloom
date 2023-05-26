// Copyright 2020 github.com/pschou/go-bloom
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A fast bloom filter.

package bloom

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"unsafe"

	zxxh3 "github.com/zeebo/xxh3"
)

type Filter struct {
	size uint64
	dat  []byte
}

// Size will create a new bloom filter in memory, note that if the size is not a multiple of 8 it will be rounded down to the multiple of 8.
func New(size int) *Filter {
	return &Filter{size: uint64(size >> 3), dat: make([]byte, size>>3)}
}

// Test if the string may be in the filter, and return hash used
func (f *Filter) TestStringHash(s string) (hash uint64, ok bool) {
	hash = zxxh3.Hash(s2b(s))
	return hash, f.dat[(hash>>3)%f.size]&(1<<(hash&0x7)) > 0
}

// Test if a string may be in the filter, and return hash used
func (f *Filter) TestHash(d []byte) (hash uint64, ok bool) {
	hash = zxxh3.Hash(d)
	return hash, f.dat[(hash>>3)%f.size]&(1<<(hash&0x7)) > 0
}

// Test if the string may be in the filter
func (f *Filter) TestString(s string) bool {
	hash := zxxh3.Hash(s2b(s))
	return f.dat[(hash>>3)%f.size]&(1<<(hash&0x7)) > 0
}

// Test if a string may be in the filter
func (f *Filter) Test(d []byte) bool {
	hash := zxxh3.Hash(d)
	return f.dat[(hash>>3)%f.size]&(1<<(hash&0x7)) > 0
}

// Add a string to the filter
func (f *Filter) AddString(s string) (hash uint64) {
	hash = zxxh3.Hash(s2b(s))
	f.dat[(hash>>3)%f.size] |= 1 << (hash & 0x7)
	return
}

// Add a byte slice to the filter
func (f *Filter) Add(d []byte) (hash uint64) {
	hash = zxxh3.Hash(d)
	//fmt.Println("hash", h%f.size)
	f.dat[(hash>>3)%f.size] |= 1 << (hash & 0x7)
	return
}

func s2b(value string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&value))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

type save struct {
	Size uint64
	Dat  []byte
}

// Save the filter into a writer
func (w *Filter) Save(fh io.Writer) (err error) {
	_, err = io.Copy(fh, bytes.NewReader(w.dat))
	return
}

// Return the filter size in bytes
func (w Filter) Size() int {
	return int(w.size)
}

// Fold will reduce the memory resident size by a factor n
func (w *Filter) Fold(n int) error {
	if n == 1 { // Do nothing
		return nil
	} else if n < 1 {
		return errors.New("Folding n has to be a positive value")
	} else if int(w.size)%n > 0 {
		return errors.New("Folding n has to be a multiple of current filter size")
	}
	sz := int(w.size) / n
	dat := make([]byte, sz)
	for i, v := range w.dat {
		dat[i%sz] |= v
	}
	w.size = uint64(sz)
	w.dat = dat
	return nil
}

// Load a filter from a reader
func Load(fh io.Reader) (*Filter, error) {
	switch f := fh.(type) {
	case *os.File:
		var size int
		if info, err := f.Stat(); err == nil {
			size64 := info.Size()
			if int64(int(size64)) == size64 {
				size = int(size64)
			}
		}
		size++ // one byte for final read at EOF

		// If a file claims a small size, read at least 512 bytes.
		// In particular, files in Linux's /proc claim size 0 but
		// then do not work right if read in small pieces,
		// so an initial read of 1 byte would not work correctly.
		if size < 512 {
			size = 512
		}

		data := make([]byte, 0, size)
		for {
			if len(data) >= cap(data) {
				d := append(data[:cap(data)], 0)
				data = d[:len(data)]
			}
			n, err := f.Read(data[len(data):cap(data)])
			data = data[:len(data)+n]
			if err != nil {
				if err == io.EOF {
					return &Filter{dat: data, size: uint64(len(data))}, nil
				}
				return nil, err
			}
		}
	default:
		var buf bytes.Buffer
		n, err := buf.ReadFrom(f)
		if err != nil {
			return nil, err
		}
		return &Filter{dat: buf.Bytes(), size: uint64(n)}, nil
	}
	return nil, nil
}
