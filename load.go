package bloom

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	gzip "github.com/klauspost/pgzip"
)

// Load a filter from a reader.  Use n = 1 to read the whole file into memory,
// and use n > 1 to read the filter in to memory but only use 1/n the space in
// memory (with a higher likelyhood of false positive hits).
func Load(fh io.Reader, n int) (*Filter, error) {
	hdr := make([]byte, 16)
	c, err := fh.Read(hdr)
	if err != nil {
		return nil, err
	}
	if c != 16 {
		return nil, fmt.Errorf("Bloom: Unable to read header (read %d)", c)
	} else if string(hdr[:8]) != "BLOOMFLT" {
		return nil, fmt.Errorf("Bloom: Invalid header %q", hdr[:8])
	}
	size := int(binary.BigEndian.Uint64(hdr[8:]))

	if n == 1 { // Do nothing
	} else if n < 1 {
		return nil, fmt.Errorf("Folding n (%d) has to be a positive value", n)
	} else if int(size)%n > 0 {
		return nil, fmt.Errorf("Folding n (%d) has to be a multiple of current filter size (%d)", n, size)
	}

	sz := size / n
	dat := make([]byte, sz)

	gzr, err := gzip.NewReader(fh)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(gzr)
	var b byte
	for i := 0; i < size; i++ {
		b, err = buf.ReadByte()
		if err != nil {
			return nil, err
		}
		dat[i%sz] |= b
	}
	return &Filter{dat: dat, size: uint64(sz)}, nil
}
