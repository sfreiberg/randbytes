package randbytes_test

import (
	"io"
	"testing"

	"github.com/sfreiberg/randbytes"
)

func TestRead(t *testing.T) {
	var size uint64 = 1000000

	reader := randbytes.NewReader(size)
	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Received an error while executing io.ReadAll: %s\n", err)
	}

	if len(bytes) != int(size) {
		t.Fatalf("Received %x but expected %x\n", len(bytes), size)
	}
}

func TestReadUsesOnlyCharset(t *testing.T) {
	chars := []byte("ab")
	data, err := io.ReadAll(randbytes.NewReaderChars(10000, chars))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 10000 {
		t.Fatalf("got %d bytes, want 10000", len(data))
	}
	var seen [256]int
	for _, b := range data {
		seen[b]++
	}
	for b, n := range seen {
		if n > 0 && b != 'a' && b != 'b' {
			t.Fatalf("byte %q is not in the character set", byte(b))
		}
	}
	// With 10000 samples both characters must show up.
	if seen['a'] == 0 || seen['b'] == 0 {
		t.Fatalf("expected both characters to appear, got a=%d b=%d", seen['a'], seen['b'])
	}
}

func TestReadSmallBuffers(t *testing.T) {
	r := randbytes.NewReaderChars(10, randbytes.NumericChars)
	buf := make([]byte, 3)
	total := 0
	for {
		n, err := r.Read(buf)
		total += n
		for _, b := range buf[:n] {
			if b < '0' || b > '9' {
				t.Fatalf("unexpected byte %q", b)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	if total != 10 {
		t.Fatalf("read %d bytes, want 10", total)
	}

	// Further reads keep returning EOF.
	if n, err := r.Read(buf); n != 0 || err != io.EOF {
		t.Fatalf("read after EOF = (%d, %v), want (0, EOF)", n, err)
	}
}

func TestReadZeroLength(t *testing.T) {
	data, err := io.ReadAll(randbytes.NewReader(0))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Fatalf("got %d bytes, want 0", len(data))
	}
}

func TestReadEmptyBuffer(t *testing.T) {
	r := randbytes.NewReader(10)
	if n, err := r.Read(nil); n != 0 || err != nil {
		t.Fatalf("Read(nil) = (%d, %v), want (0, nil)", n, err)
	}
}

func TestReadEmptyCharset(t *testing.T) {
	r := randbytes.NewReaderChars(10, nil)
	if _, err := r.Read(make([]byte, 5)); err != randbytes.ErrEmptyCharset {
		t.Fatalf("err = %v, want ErrEmptyCharset", err)
	}
}

func TestReadersDiffer(t *testing.T) {
	a, _ := io.ReadAll(randbytes.NewReader(64))
	b, _ := io.ReadAll(randbytes.NewReader(64))
	if string(a) == string(b) {
		t.Fatal("two readers produced identical output")
	}
}
