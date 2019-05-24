package chippy

import (
	"github.com/golang-collections/collections/stack"
	"math/rand"
	"time"
)

type CPU struct {
	pc        uint16
	progStack stack.Stack
	registers map[uint8]uint8
	vf        bool
	i         uint16
	sp        uint8
	dt        uint8
	st        uint8
}

func boot() {
	cpu := new(CPU)
	cpu.pc = 0
	cpu.vf = false
	cpu.i = 0
	cpu.sp = 0
	cpu.dt = 0
	cpu.st = 0
	cpu.registers = map[uint8]uint8{
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

func (cpu CPU) command(instruction uint16) {
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
		switch command := instruction; {
		case (command & 0x8001) == 0x8001:
			cpu.OR(instruction)
		case (command & 0x8002) == 0x8002:
			cpu.AND(instruction)
		case (command & 0x8003) == 0x8003:
			cpu.XOR(instruction)
		case (command & 0x8004) == 0x8004:
			cpu.ADDVY(instruction)
		case (command & 0x8005) == 0x8005:
			cpu.SUB(instruction)
		case (command & 0x8006) == 0x8006:
			cpu.SHR(instruction)
		case (command & 0x8007) == 0x8007:
			cpu.SUBN(instruction)
		case (command & 0x800e) == 0x800e:
			cpu.SHL(instruction)
		}
	case 0x9000:
		cpu.SNE(instruction)
	case 0xA000:
		cpu.LDI(instruction)
	case 0xB000:
		cpu.JP0(instruction)
	case 0xC000:
		cpu.RND(instruction)
	case 0xD000:
		cpu.DRW(instruction)
	case 0xE000:
		switch command := instruction; {
		case (command & 0xe09e) == 0xe090:
			cpu.SKP(instruction)
		case (command & 0xe0a1) == 0xe0a1:
			cpu.SKNP(instruction)
		}
	case 0xF000:
		switch command := instruction; {
		case (command & 0xf007) == 0xf007:
			cpu.LD_VX_DT(instruction)
		case (command & 0xf00A) == 0xf00A:
			cpu.LD_VX_K(instruction)
		case (command & 0xf015) == 0xf015:
			cpu.LDDT(instruction)
		case (command & 0xf018) == 0xf018:
			cpu.LDST(instruction)
		case (command & 0xf01E) == 0xf01E:
			cpu.ADDI(instruction)
		case (command & 0xf029) == 0xf029:
			cpu.LDF(instruction)
		case (command & 0xf033) == 0xf033:
			cpu.LDB(instruction)
		case (command & 0xf055) == 0xf055:
			cpu.LDIVX(instruction)
		case (command & 0xf065) == 0xf065:
			cpu.LDVXI(instruction)
		}
	}
}

func (cpu CPU) SYS(instruction uint16) {
	//Unused??
}

func (cpu CPU) CLS(instruction uint16) {
	//Clear the screen
}

func (cpu CPU) RET(instruction uint16) {
	//Return from function
}

func (cpu CPU) JP(instruction uint16) {
	cpu.pc = instruction & 0x0fff
}

func (cpu CPU) CALL(instruction uint16) {
	cpu.pc += 1
	cpu.progStack.Push(cpu.pc)
	cpu.JP(instruction)
}

func (cpu *CPU) SEKK(instruction uint16) {
	if cpu.getVx(instruction) == getKk(instruction) {
		cpu.pc += 2
	}
}

func (cpu *CPU) SNE(instruction uint16) {
	if cpu.getVx(instruction) != getKk(instruction) {
		cpu.pc += 2
	}
}

func (cpu *CPU) SEVY(instruction uint16) {
	if cpu.getVx(instruction) == cpu.getVy(instruction) {
		cpu.pc += 2
	}
}

func (cpu *CPU) LD(instruction uint16) {
	cpu.registers[getVxAddress(instruction)] = getKk(instruction)
}

func (cpu *CPU) ADDKK(instruction uint16) {
	total := cpu.getVx(instruction) + getKk(instruction)
	cpu.setRegisterFromVx(instruction, total)
}

func (cpu *CPU) LDVY(instruction uint16) {
	cpu.registers[getVxAddress(instruction)] = cpu.getVy(instruction)
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
	cpu.vf = vx > added
	cpu.setRegisterFromVx(instruction, added)
}

func (cpu *CPU) SUB(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	cpu.vf = vx > vy
	cpu.setRegisterFromVx(instruction, vx-vy)
}

func (cpu *CPU) SHR(instruction uint16) {
	vx := cpu.getVx(instruction)
	cpu.vf = ((vx << 7) & 0x80) == 0x80
	cpu.setRegisterFromVx(instruction, vx>>2)
}

func (cpu *CPU) SHL(instruction uint16) {
	vx := cpu.getVx(instruction)
	cpu.vf = ((vx >> 7) & 0x01) == 0x01
	cpu.setRegisterFromVx(instruction, vx<<2)
}

func (cpu *CPU) SNEVY(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	if vx != vy {
		cpu.pc += 2
	}
}

func (cpu *CPU) LDI(instruction uint16) {
	cpu.i = instruction & 0x0fff
}

func (cpu *CPU) JP0(instruction uint16) {
	cpu.pc = instruction&0x0fff + uint16(cpu.registers[0])
}

func (cpu *CPU) DRW(instruction uint16) {
	//show sprite
}

func (cpu *CPU) SKP(instruction uint16) {
	//Skip if reg matched key is pressed
}

func (cpu *CPU) SKNP(instruction uint16) {
	//Skip if reg matched key is not pressed
}

func (cpu *CPU) LD_VX_DT(instruction uint16) {
	cpu.setRegisterFromVx(instruction, cpu.dt)
}

func (cpu *CPU) LDDT(instruction uint16) {
	cpu.dt = cpu.getVx(instruction)
}

func (cpu *CPU) LD_VX_K(instruction uint16) {
	//Wait for a key press, store the value of the key in Vx.
}

func (cpu *CPU) LDST(instruction uint16) {
	cpu.st = cpu.getVx(instruction)
}

func (cpu *CPU) ADDI(instruction uint16) {
	cpu.i = cpu.i + uint16(cpu.getVx(instruction))
}

func (cpu *CPU) LDF(instruction uint16) {
	//Set I = location of sprite for digit Vx.
}

func (cpu *CPU) LDB(instruction uint16) {
	//Store BCD representation of Vx in memory locations I, I+1, and I+2.
}

func (cpu *CPU) LDIVX(instruction uint16) {
	//Store registers V0 through Vx in memory starting at location I.
}

func (cpu *CPU) LDVXI(instruction uint16) {
	//Read registers V0 through Vx from memory starting at location I.
}

func (cpu *CPU) RND(instruction uint16) {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	cpu.setRegisterFromVx(instruction, uint8(r.Intn(256)))
}

func (cpu *CPU) SUBN(instruction uint16) {
	vx := cpu.getVx(instruction)
	vy := cpu.getVy(instruction)
	cpu.vf = vy > vx
	cpu.setRegisterFromVx(instruction, vy-vx)
}

func (cpu *CPU) setRegisterFromVx(instruction uint16, value uint8) {
	cpu.registers[getVxAddress(instruction)] = value
}

func (cpu *CPU) setRegisterFromVy(instruction uint16, value uint8) {
	cpu.registers[getVyAddress(instruction)] = value
}

func (cpu *CPU) setRegister(address uint8, value uint8) {
	cpu.registers[address] = value
}

func getVxAddress(instruction uint16) uint8 {
	return uint8((instruction >> 8) & 0x0f)
}

func getVyAddress(instruction uint16) uint8 {
	return uint8((instruction >> 4) & 0x0f)
}

func (cpu *CPU) getVx(instruction uint16) uint8 {
	return cpu.registers[uint8(getVxAddress(instruction))]
}

func (cpu *CPU) getVy(instruction uint16) uint8 {
	return cpu.registers[uint8(getVyAddress(instruction))]
}

func getKk(instruction uint16) uint8 {
	return uint8(instruction & 0x00ff)
}
