// This package is useful for generating large amounts of random data.
// For example if you want to generate 1GB of random data you can grab
// data in smaller chunks instead of generating all of it at once.
// This library is not appropriate for security sensitive work.
package randbytes

import (
	"io"
	"math/rand"
)

var AlphaChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var NumericChars = []byte("0123456789")
var AlphaNumChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Generate a new Reader with an alphanumberic character set.
func NewReader(bytes uint64) *Reader {
	return NewReaderChars(bytes, AlphaNumChars)
}

// Generate a new Reader with a specific character set.
func NewReaderChars(bytes uint64, chars []byte) *Reader {
	return &Reader{bytes: bytes, chars: chars}
}

// The Reader struct implements the io.Reader interface and is
// useful for creating a reader that needs to create large files.
// For example if you want to create a 5GB file but don't want to
// allocate all 5GB of data in memory you can use this to grab data
// at much smaller chunks.
type Reader struct {
	bytes uint64 // total bytes to generate
	read  uint64 // number of bytes already read
	chars []byte // character set to use
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, io.EOF
	}

	size := len(p)
	remaining := int(r.bytes - r.read)

	// If the remaining bytes are smaller than what the caller asked for
	// only fill the remaining bytes and set err to io.EOF to signify
	// that we're done.
	if remaining < size {
		size = remaining
		err = io.EOF
	}

	for i := 0; i < size; i++ {
		p[i] = r.chars[rand.Intn(len(r.chars))]
		n++
		r.read++
	}

	return
}
