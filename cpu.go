package chippy

import (
	"github.com/golang-collections/collections/stack"
)

type CPU struct {
	pc        uint16
	progStack stack.Stack
	registers map[uint8]uint8
}

func (cpu CPU) sys(instruction uint16) {
	//Unused??
}

func (cpu CPU) cls(instruction uint16) {
	//Clear the screen
}

func (cpu CPU) ret(instruction uint16) {
	//Return from function
}

func (cpu CPU) jp(instruction uint16) {
	cpu.pc = instruction & 0x0fff
}

func (cpu CPU) call(instruction uint16) {
	cpu.pc += 1
	cpu.progStack.Push(cpu.pc)
	cpu.jp(instruction)
}

func (cpu CPU) se(instruction uint16) {
	if getVx(instruction) == getKk(instruction) {
		cpu.pc += 2
	}
}

func getVx(instruction uint16) uint8 {
	return uint8((instruction >> 8) & 0x0f)
}

func getKk(instruction uint16) uint8 {
	return uint8(instruction & 0x00ff)
}
