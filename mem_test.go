package chippy_test

import (
	"github.com/jayshaffer/chippy"
	"testing"
)

func TestMemory(t *testing.T) {
	bytes := []byte{0x11, 0x22}
	mem := chippy.Load(bytes)
	if len(mem.ProgData) != 2 {
		t.Error("Memory not initialized correctly")
	}
	if mem.ProgData[0x200] == 0 {
		t.Error("Memory initialized at the correct address")
	}
}
