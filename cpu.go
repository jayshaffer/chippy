package chippy

import (
	"github.com/golang-collections/collections/stack"
)

type CPU struct {
	pc        uint16
	progStack stack.Stack
	registers map[uint8]uint8
	vf        uint16
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
	added := cpu.getVx(instruction) ^ cpu.getVy(instruction)
	if added > 255 {
		cpu.vf = 1
	} else {
		cpu.vf = 0
	}
	cpu.setRegisterFromVx(instruction, added)
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
