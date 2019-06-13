package chippy

import (
	"encoding/binary"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"math/rand"
	"time"
)

type CPU struct {
	PC         uint16
	ProgStack  stack.Stack
	Registers  map[uint8]uint8
	PRM        *Memory
	I          uint16
	SP         uint8
	Dt         uint8
	St         uint8
	Waiting    bool
	Jumped     bool
	DelayTimer *time.Ticker
}

func (cpu *CPU) Tick() {
	//cpu.LogStatus()
	cpu.HandleTimerTick()
	cpu.command(cpu.LoadCommandBytes())
	if !cpu.Waiting && !cpu.Jumped {
		cpu.PC += 2
	}
	cpu.Jumped = false
}

func (cpu *CPU) HandleTimerTick() {
	select {
	case <-cpu.DelayTimer.C:
		if cpu.Dt > 0 {
			cpu.Dt--
		}
	default:
		return
	}
}

func (cpu *CPU) LoadCommandBytes() uint16 {
	first := cpu.PRM.ProgData[cpu.PC]
	second := cpu.PRM.ProgData[cpu.PC+1]
	return binary.BigEndian.Uint16([]byte{first, second})
}

func (cpu *CPU) LogStatus() {
	fmt.Printf("Command: 0x%x\n", cpu.LoadCommandBytes())
	fmt.Printf("PC: 0x%x | I: 0x%x | Dt: 0x%x | St: 0x%x | SP: 0x%x\n", cpu.PC, cpu.I, cpu.Dt, cpu.St, cpu.SP)
	fmt.Printf("Registers: %x\n", cpu.Registers)
}

func (cpu *CPU) Boot(memory *Memory) {
	cpu.DelayTimer = time.NewTicker(time.Second / 60)
	cpu.Dt = 0
	cpu.PC = uint16(0x0200)
	cpu.I = 0
	cpu.SP = 0
	cpu.Dt = 0
	cpu.St = 0
	cpu.Jumped = false
	cpu.ProgStack = *stack.New()
	cpu.PRM = memory
	cpu.Waiting = false
	cpu.Registers = map[uint8]uint8{
		0x00: 0,
		0x01: 0,
		0x02: 0,
		0x03: 0,
		0x04: 0,
		0x05: 0,
		0x06: 0,
		0x07: 0,
		0x08: 0,
		0x09: 0,
		0x0A: 0,
		0x0B: 0,
		0x0C: 0,
		0x0D: 0,
		0x0E: 0,
		0x0F: 0}
}

func (cpu *CPU) command(instruction uint16) {
	masked := instruction & 0xf000
	switch masked {
	case 0x0000:
		switch instruction {
		case 0x00e0:
			cpu.CLS(instruction)
		case 0x00ee:
			cpu.RET(instruction)
		}
	case 0x1000:
		cpu.JP(instruction)
	case 0x2000:
		cpu.CALL(instruction)
	case 0x3000:
		cpu.SEKK(instruction)
	case 0x4000:
		cpu.SNE(instruction)
	case 0x5000:
		cpu.SEVY(instruction)
	case 0x6000:
		cpu.LD(instruction)
	case 0x7000:
		cpu.ADDKK(instruction)
	case 0x8000:
		var command uint16 = instruction & 0x000f
		switch command {
		case 0x0000:
			cpu.LDVXVY(instruction)
		case 0x0001:
			cpu.OR(instruction)
		case 0x0002:
			cpu.AND(instruction)
		case 0x0003:
			cpu.XOR(instruction)
		case 0x0004:
			cpu.ADDVY(instruction)
		case 0x0005:
			cpu.SUB(instruction)
		case 0x0006:
			cpu.SHR(instruction)
		case 0x0007:
			cpu.SUBN(instruction)
		case 0x000e:
			cpu.SHL(instruction)
		}
	case 0x9000:
		cpu.SNEVY(instruction)
	case 0xA000:
		cpu.LDI(instruction)
	case 0xB000:
		cpu.JP0(instruction)
	case 0xC000:
		cpu.RND(instruction)
	case 0xD000:
		cpu.DRW(instruction)
	case 0xE000:
		var command uint16 = instruction & 0xf0ff
		switch command {
		case 0xe09e:
			cpu.SKP(instruction)
		case 0xe0a1:
			cpu.SKNP(instruction)
		}
	case 0xF000:
		var command16 uint16 = instruction & 0x00ff
		switch command16 {
		case 0x0007:
			cpu.LD_VX_DT(instruction)
		case 0x000A:
			cpu.LD_VX_K(instruction)
		case 0x0015:
			cpu.LDDT(instruction)
		case 0x0018:
			cpu.LDST(instruction)
		case 0x001E:
			cpu.ADDI(instruction)
		case 0x0029:
			cpu.LDF(instruction)
		case 0x0033:
			cpu.LDB(instruction)
		case 0x0055:
			cpu.LDIVX(instruction)
		case 0x0065:
			cpu.LDVXI(instruction)
		}
	}
}

func (cpu *CPU) SYS(instruction uint16) {
}

func (cpu *CPU) CLS(instruction uint16) {
	cpu.PRM.ClearDisplay()
}

func (cpu *CPU) RET(instruction uint16) {
	cpu.PC = cpu.ProgStack.Pop().(uint16)
}

func (cpu *CPU) JP(instruction uint16) {
	cpu.PC = (instruction & 0x0fff)
	cpu.Jumped = true
}

func (cpu *CPU) CALL(instruction uint16) {
	cpu.ProgStack.Push(cpu.PC)
	cpu.JP(instruction)
}

func (cpu *CPU) SEKK(instruction uint16) {
	if cpu.getVx(instruction) == getKk(instruction) {
		cpu.PC += 2
	}
}

func (cpu *CPU) SNE(instruction uint16) {
	if cpu.getVx(instruction) != getKk(instruction) {
		cpu.PC += 2
	}
}

func (cpu *CPU) SEVY(instruction uint16) {
	if cpu.getVx(instruction) == cpu.getVy(instruction) {
		cpu.PC += 2
	}
}

func (cpu *CPU) LD(instruction uint16) {
	cpu.Registers[getVxAddress(instruction)] = getKk(instruction)
}

func (cpu *CPU) ADDKK(instruction uint16) {
	total := cpu.getVx(instruction) + getKk(instruction)
	cpu.setRegisterFromVx(instruction, total)
}

func (cpu *CPU) LDVY(instruction uint16) {
	cpu.Registers[getVxAddress(instruction)] = cpu.getVy(instruction)
}

func (cpu *CPU) OR(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.getVx(instruction)|cpu.getVy(instruction))
}

func (cpu *CPU) AND(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.getVx(instruction)&cpu.getVy(instruction))
}

func (cpu *CPU) XOR(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.getVx(instruction)^cpu.getVy(instruction))
}

func (cpu *CPU) ADDVY(instruction uint16) {
	vx := cpu.getVx(instruction)
	added := vx + cpu.getVy(instruction)
	cpu.setVF(added > 255)
	cpu.setRegisterFromVx(instruction, added&0xff)
}

func (cpu *CPU) SUB(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	cpu.setVF(vx > vy)
	cpu.setRegisterFromVx(instruction, vx-vy)
}

func (cpu *CPU) SHR(instruction uint16) {
	vx := cpu.getVx(instruction)
	cpu.setVF((vx & 0x01) == 0x01)
	cpu.setRegisterFromVx(instruction, vx>>1)
}

func (cpu *CPU) SHL(instruction uint16) {
	vx := cpu.getVx(instruction)
	cpu.setVF(((vx >> 7) & 0x01) == 0x01)
	cpu.setRegisterFromVx(instruction, vx<<1)
}

func (cpu *CPU) SNEVY(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	if vx != vy {
		cpu.PC += 2
	}
}

func (cpu *CPU) LDI(instruction uint16) {
	cpu.I = instruction & 0x0fff
}

func (cpu *CPU) JP0(instruction uint16) {
	cpu.PC = instruction&0x0fff + uint16(cpu.Registers[0])
	cpu.Jumped = true
}

func (cpu *CPU) DRW(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	nib := getNib(instruction)
	address := cpu.I
	col := false
	if vy >= 32 {
		vy = 0
	}
	if vx >= 64 {
		vx = 0
	}
	for i := 0; i < int(nib); i++ {
		addressByte := cpu.PRM.ProgData[address]
		y := vy
		for j := 0; j < 8; j++ {
			x := (uint16(vx) + uint16(j))
			if x >= 64 || y >= 32 {
				continue
			}
			spriteBit := (addressByte >> uint(7-j)) & 0x01
			i := cpu.PRM.DisplayMem[x][y]
			if !col {
				col = i == 1 && spriteBit == 1
			}
			cpu.PRM.DisplayMem[x][y] = i ^ spriteBit
		}
		address++
		vy++
	}
	cpu.setVF(col)
}

func (cpu *CPU) SKP(instruction uint16) {
	vx := cpu.getVx(instruction)
	if val, ok := cpu.PRM.KeyboardMem[vx]; ok && val > 0 {
		cpu.PC += 2
	}
}

func (cpu *CPU) SKNP(instruction uint16) {
	vx := cpu.getVx(instruction)
	if val, ok := cpu.PRM.KeyboardMem[vx]; !ok || val == 0 {
		cpu.PC += 2
	}
}

func (cpu *CPU) LD_VX_DT(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.Dt)
}

func (cpu *CPU) LDDT(instruction uint16) {
	cpu.Dt = cpu.getVx(instruction)
}

func (cpu *CPU) LD_VX_K(instruction uint16) {
	for i := range cpu.PRM.KeyboardMem {
		pressed := cpu.PRM.KeyboardMem[i]
		if pressed > 0 {
			cpu.setRegisterFromVx(instruction, pressed)
			cpu.Waiting = false
			return
		}
	}
	cpu.Waiting = true
}

func (cpu *CPU) LDST(instruction uint16) {
	cpu.St = cpu.getVx(instruction)
}

func (cpu *CPU) ADDI(instruction uint16) {
	cpu.I = cpu.I + uint16(cpu.getVx(instruction))
}

func (cpu *CPU) LDF(instruction uint16) {
	cpu.I = uint16(0x0100) + (uint16(cpu.getVx(instruction)))
}

func (cpu *CPU) LDVXVY(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.getVy(instruction))
}

func (cpu *CPU) LDB(instruction uint16) {
	vx := int(cpu.getVx(instruction))
	cpu.PRM.ProgData[cpu.I] = uint8(vx / 100)
	cpu.PRM.ProgData[cpu.I+1] = uint8(vx % 100 / 10)
	cpu.PRM.ProgData[cpu.I+2] = uint8(vx % 10)
}

func (cpu *CPU) LDIVX(instruction uint16) {
	for i := uint8(0); i <= cpu.getVx(instruction); i++ {
		cpu.PRM.ProgData[cpu.I+uint16(i)] = cpu.Registers[i]
	}
}

func (cpu *CPU) LDVXI(instruction uint16) {
	for i := uint8(0); i <= cpu.getVx(instruction); i++ {
		cpu.Registers[i] = cpu.PRM.ProgData[cpu.I+uint16(i)]
	}
}

func (cpu *CPU) RND(instruction uint16) {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed).Intn(256) & int(getKk(instruction))
	cpu.setRegisterFromVx(instruction, uint8(r))
}

func (cpu *CPU) SUBN(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	cpu.setVF(vy > vx)
	cpu.setRegisterFromVx(instruction, vy-vx)
}

func (cpu *CPU) setRegisterFromVx(instruction uint16, value uint8) {
	cpu.Registers[getVxAddress(instruction)] = value
}

func (cpu *CPU) setRegisterFromVy(instruction uint16, value uint8) {
	cpu.Registers[getVyAddress(instruction)] = value
}

func (cpu *CPU) setRegister(address uint8, value uint8) {
	cpu.Registers[address] = value
}

func (cpu *CPU) setVF(flag bool) {
	if flag {
		cpu.Registers[0x0f] = 0x01
	} else {
		cpu.Registers[0x0f] = 0
	}
}

func getVxAddress(instruction uint16) uint8 {
	return uint8((instruction & 0x0f00) >> 8)
}

func getVyAddress(instruction uint16) uint8 {
	return uint8((instruction & 0x00f0) >> 4)
}

func (cpu *CPU) getVx(instruction uint16) uint8 {
	return cpu.Registers[uint8(getVxAddress(instruction))]
}

func (cpu *CPU) getVy(instruction uint16) uint8 {
	return cpu.Registers[uint8(getVyAddress(instruction))]
}

func getNib(instruction uint16) uint8 {
	return uint8(instruction & 0x000f)
}

func getKk(instruction uint16) uint8 {
	return uint8(instruction & 0x00ff)
}
