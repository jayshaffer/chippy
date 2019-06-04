package chippy_test

import (
	"encoding/binary"
	"fmt"
	"github.com/jayshaffer/chippy"
	"testing"
)

func TestJP(t *testing.T) {
	cpu := new(chippy.CPU)
	mem := new(chippy.Memory)
	cpu.Boot(mem)
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
	mem := new(chippy.Memory)
	cpu.Boot(mem)
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
	mem := new(chippy.Memory)
	cpu.Boot(mem)
	cpu.Registers[0x01] = 0x1b
	bytes := binary.BigEndian.Uint16([]byte{0x01, 0x1b})
	cpu.SEKK(bytes)
	if cpu.PC != 2 {
		t.Error("PC did not increment")
	}
}

func TestSEKKNeg(t *testing.T) {
	cpu := new(chippy.CPU)
	mem := new(chippy.Memory)
	cpu.Boot(mem)
	cpu.Registers[0x01] = 0x1b
	bytes := binary.BigEndian.Uint16([]byte{0x01, 0x1c})
	cpu.SEKK(bytes)
	if cpu.PC > 0 {
		t.Error("PC incremented incorrectly")
	}
}

func TestDRW(t *testing.T) {
	bytes := make([]byte, 5000)
	mem := chippy.Load(bytes)
	cpu := new(chippy.CPU)
	mem.ProgData[0x00] = 0xff
	mem.ProgData[0x01] = 0xff
	mem.ProgData[0x02] = 0xff
	mem.ProgData[0x03] = 0xff
	mem.ProgData[0x04] = 0xff
	mem.ProgData[0x05] = 0xff
	mem.ProgData[0x06] = 0xff
	mem.ProgData[0x07] = 0xff
	mem.ProgData[0x08] = 0x00
	cpu.Boot(mem)
	cpu.I = 0x00
	instruction := binary.BigEndian.Uint16([]byte{0xd0, 0x08})
	cpu.DRW(instruction)
	if mem.DisplayMem[0x00][0x00] != uint8(0x01) {
		fmt.Print(mem.DisplayMem)
		t.Error("Byte pattern for display mem incorrect")
	}
}

func TestLDB(t *testing.T) {
	bytes := make([]byte, 5000)
	cpu := new(chippy.CPU)
	mem := chippy.Load(bytes)
	cpu.Boot(mem)
	cpu.I = 0x00
	cpu.Registers[0] = 0xfe
	instruction := binary.BigEndian.Uint16([]byte{0xf0, 0x33})
	cpu.LDB(instruction)
	if cpu.PRM.ProgData[cpu.I] != 2 {
		t.Error(fmt.Sprintf("LDB command hundreds place incorrect: %d", cpu.PRM.ProgData[cpu.I]))
	}

	if cpu.PRM.ProgData[cpu.I+1] != 5 {
		t.Error(fmt.Sprintf("LDB command tens place incorrect: %d", cpu.PRM.ProgData[cpu.I+1]))
	}

	if cpu.PRM.ProgData[cpu.I+2] != 4 {
		t.Error(fmt.Sprintf("LDB command ones place incorrect: %d", cpu.PRM.ProgData[cpu.I+2]))
	}
}
