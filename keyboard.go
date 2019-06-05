package chippy

import (
	"github.com/veandco/go-sdl2/sdl"
)

var KeyMap = map[int]uint8{
	30: 0x01,
	31: 0x02,
	32: 0x03,
	33: 0x0c,
	20: 0x04,
	26: 0x05,
	8:  0x06,
	21: 0x0d,
	4:  0x07,
	22: 0x08,
	7:  0x09,
	9:  0x0e,
	29: 0x0a,
	27: 0x00,
	6:  0x0b,
	25: 0x0f,
}

type Keyboard struct {
	MemChunk *Memory
}

func NewKeyboard(memory *Memory) *Keyboard {
	keyboard := new(Keyboard)
	keyboard.MemChunk = memory
	return keyboard
}

func (keyboard *Keyboard) HandleKeyPress(event *sdl.KeyboardEvent) {
	if val, ok := KeyMap[int(event.Keysym.Scancode)]; ok {
		keyboard.MemChunk.KeyboardMem[val] = event.State
	}
}
