package chippy_test

import (
	"fmt"
	"github.com/jayshaffer/chippy"
	"testing"
)

func TestMemory(t *testing.T) {
	bytes := []byte{0x11, 0x22}
	mem := chippy.Load(bytes)
	for i := uint16(0x100); i < uint16(0x200); i++ {
		fmt.Println(mem.ProgData[i])
	}
	if len(mem.ProgData) != 2 {
		t.Error("Memory not initialized correctly")
	}
	if mem.ProgData[0x200] == 0 {
		t.Error("Memory initialized at the correct address")
	}
}
