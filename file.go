package bloom

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	gzip "github.com/klauspost/pgzip"
)

type save struct {
	Size uint64
	Dat  []byte
}

// Save the filter into a writer
func (w *Filter) Save(fh io.Writer) (err error) {
	hdr := []byte("BLOOMFLT        ")
	binary.BigEndian.PutUint64(hdr[8:], w.size)
	_, err = fh.Write(hdr)
	if err != nil {
		return
	}
	gzw := gzip.NewWriter(fh)
	_, err = io.Copy(gzw, bytes.NewReader(w.dat))
	gzw.Close()
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
