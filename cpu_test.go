package chippy_test

import (
	"encoding/binary"
	"fmt"
	"github.com/jayshaffer/chippy"
	"testing"
)

func TestJP(t *testing.T) {
	cpu := new(chippy.CPU)
	cpu.Boot()
	if cpu.PC != 0 {
		t.Error()
	}
	bytes := binary.BigEndian.Uint16([]byte{0xfe, 0xee})
	cpu.JP(bytes)
	if cpu.PC != 0x0eee {
		t.Error(fmt.Sprintf("CPU Jump failed, result was %x, should have been %x", cpu.PC, 0x0eee))
	}
}

func TestCALL(t *testing.T) {
	cpu := new(chippy.CPU)
	cpu.Boot()
	cpu.PC = 0x00bb
	bytes := binary.BigEndian.Uint16([]byte{0x00, 0xaa})
	cpu.CALL(bytes)
	if cpu.PC != 0xaa {
		t.Error(fmt.Sprintf("Program Counter wasn't! Value is at: %x", cpu.PC))
	}
	pcState := cpu.ProgStack.Pop()
	if pcState != uint16(0x00bb) {
		t.Error(
			fmt.Sprintf(
				"Program stack didn't store the PC correctly! Got: %x, should be: %x",
				pcState,
				0xbb,
			))
	}
}

func TestSEKKPos(t *testing.T) {
	cpu := new(chippy.CPU)
	cpu.Boot()
	cpu.Registers[0x01] = 0x1b
	bytes := binary.BigEndian.Uint16([]byte{0x01, 0x1b})
	cpu.SEKK(bytes)
	if cpu.PC != 2 {
		t.Error("PC did not increment")
	}
}

func TestSEKKNeg(t *testing.T) {
	cpu := new(chippy.CPU)
	cpu.Boot()
	cpu.Registers[0x01] = 0x1b
	bytes := binary.BigEndian.Uint16([]byte{0x01, 0x1c})
	cpu.SEKK(bytes)
	if cpu.PC > 0 {
		t.Error("PC incremented incorrectly")
	}
}
