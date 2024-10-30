// Copyright 2020 github.com/pschou/go-bloom-worm
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

package bwdb

import (
	"fmt"
	"reflect"
	"unsafe"

	zxxh3 "github.com/zeebo/xxh3"
)

type Filter struct {
	Data []byte
}

// Test if the string may be in the filter
func (f *Filter) TestString(s string) bool {
	hash := zxxh3.Hash(s2b(s))
	return f.Data[int(hash>>3)%len(f.Data)]&(1<<(hash&0x7)) > 0
}

// Test if a byte slice may be in the filter
func (f *Filter) Test(d []byte) bool {
	hash := zxxh3.Hash(d)
	return f.Data[int(hash>>3)%len(f.Data)]&(1<<(hash&0x7)) > 0
}

// Add a string to the filter
func (f *Filter) AddString(s string) (hash uint64) {
	hash = zxxh3.Hash(s2b(s))
	f.Data[int(hash>>3)%len(f.Data)] |= 1 << (hash & 0x7)
	return
}

// Add a byte slice to the filter
func (f *Filter) Add(d []byte) (hash uint64) {
	hash = zxxh3.Hash(d)
	f.Data[int(hash>>3)%len(f.Data)] |= 1 << (hash & 0x7)
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

// Fold will reduce the memory resident size by a factor n
func (w *Filter) Fold(n int) error {
	if n == 1 { // Do nothing
		return nil
	} else if n < 1 {
		return fmt.Errorf("Folding n (%d) has to be a positive value", n)
	} else if len(w.Data)%n > 0 {
		return fmt.Errorf("Folding n (%d) has to be a multiple of current filter size (%d)", n, len(w.Data))
	}
	sz := len(w.Data) / n
	dat := make([]byte, sz)
	for i, v := range w.Data {
		dat[i%sz] |= v
	}
	w.Data = dat
	return nil
}
