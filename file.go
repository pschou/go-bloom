package bloom

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type save struct {
	Size uint64
	Dat  []byte
}

// Save the filter into a writer
func (w *Filter) Save(fh io.Writer) (err error) {
	hdr := []byte("BLOOMFLT        ")
	binary.BigEndian.PutUint64(hdr[8:], w.size)
	_, err = io.Copy(fh, io.MultiReader(bytes.NewReader(hdr), bytes.NewReader(w.dat)))
	return
}

// Save bloom filter to disk into file name specified.
func (f *Filter) SaveFile(file string) error {
	flt, err := os.Create(file)
	if err != nil {
		return err
	}
	defer flt.Close()
	return f.Save(flt)
}

// Load a file into the bloom filter, folding by value n.  Use n = 1 to load
// the whole file into memory 1-to-1.
func LoadFile(file string, n int) (*Filter, error) {
	flt, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer flt.Close()
	return Load(flt, n)
}
