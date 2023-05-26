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
	"fmt"
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

// Return the filter size in bytes
func (w Filter) Size() int {
	return int(w.size)
}

// Fold will reduce the memory resident size by a factor n
func (w *Filter) Fold(n int) error {
	if n == 1 { // Do nothing
		return nil
	} else if n < 1 {
		return fmt.Errorf("Folding n (%d) has to be a positive value", n)
	} else if int(w.size)%n > 0 {
		return fmt.Errorf("Folding n (%d) has to be a multiple of current filter size (%d)", n, w.size)
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
