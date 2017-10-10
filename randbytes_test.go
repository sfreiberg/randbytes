package randbytes_test

import (
	"io/ioutil"
	"testing"

	"github.com/sfreiberg/randbytes"
)

func TestRead(t *testing.T) {
	var size uint64 = 1000000

	reader := randbytes.NewReader(size)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Received an error while executing ioutil.ReadAll: %s\n", err)
	}

	if len(bytes) != int(size) {
		t.Fatalf("Received %x but expected %x\n", len(bytes), size)
	}
}
