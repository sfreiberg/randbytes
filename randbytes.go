// This package is useful for generating large amounts of random data.
// For example if you want to generate 1GB of random data you can grab
// data in smaller chunks instead of generating all of it at once.
//
// Each Reader draws from its own ChaCha8 generator seeded from
// crypto/rand, and characters are selected without modulo bias.
package randbytes

import (
	crand "crypto/rand"
	"errors"
	"io"
	"math/rand/v2"
)

var AlphaChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var NumericChars = []byte("0123456789")
var AlphaNumChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// ErrEmptyCharset is returned by Read when the Reader has no characters to
// choose from.
var ErrEmptyCharset = errors.New("randbytes: empty character set")

// Generate a new Reader with an alphanumberic character set.
func NewReader(bytes uint64) *Reader {
	return NewReaderChars(bytes, AlphaNumChars)
}

// Generate a new Reader with a specific character set.
func NewReaderChars(bytes uint64, chars []byte) *Reader {
	var seed [32]byte
	crand.Read(seed[:]) // never returns an error; crashes the program on failure
	return &Reader{bytes: bytes, chars: chars, rng: rand.New(rand.NewChaCha8(seed))}
}

// The Reader struct implements the io.Reader interface and is
// useful for creating a reader that needs to create large files.
// For example if you want to create a 5GB file but don't want to
// allocate all 5GB of data in memory you can use this to grab data
// at much smaller chunks.
//
// A Reader is not safe for concurrent use.
type Reader struct {
	bytes uint64     // total bytes to generate
	read  uint64     // number of bytes already read
	chars []byte     // character set to use
	rng   *rand.Rand // cryptographically strong generator (ChaCha8)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	remaining := r.bytes - r.read
	if remaining == 0 {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	if len(r.chars) == 0 {
		return 0, ErrEmptyCharset
	}
	if r.rng == nil {
		var seed [32]byte
		crand.Read(seed[:])
		r.rng = rand.New(rand.NewChaCha8(seed))
	}

	// If the remaining bytes are smaller than what the caller asked for
	// only fill the remaining bytes and set err to io.EOF to signify
	// that we're done.
	if remaining <= uint64(len(p)) {
		p = p[:remaining]
		err = io.EOF
	}

	for i := range p {
		p[i] = r.chars[r.rng.IntN(len(r.chars))]
	}
	n = len(p)
	r.read += uint64(n)

	return n, err
}
