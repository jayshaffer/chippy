package chippy

type Memory struct {
	ProgData   map[uint16]uint8
	SystemMem  []bool
	DisplayMem [][]uint8
}

func Load(uint8s []byte) *Memory {
	mem := new(Memory)
	mem.ProgData = make(map[uint16]uint8, 3895)
	mem.SystemMem = make([]bool, 16)
	mem.ClearDisplay()
	var index uint16
	index = 0x200
	mem.LoadCharSprites()
	for i := 0; i < len(uint8s); i++ {
		mem.ProgData[index] = uint8s[i]
		index++
	}
	return mem
}

func (mem *Memory) ClearDisplay() {
	dispSlice := make([][]uint8, 64)
	for i := range dispSlice {
		dispSlice[i] = make([]uint8, 32)
	}
	mem.DisplayMem = dispSlice
}

func (mem *Memory) LoadCharSprites() {
	sprite0 := [5]uint8{
		0xf0,
		0x90,
		0x90,
		0x90,
		0xf0,
	}
	sprite1 := [5]uint8{
		0x20,
		0x60,
		0x20,
		0x20,
		0x70,
	}
	sprite2 := [5]uint8{
		0xf0,
		0x10,
		0xf0,
		0x80,
		0xf0,
	}
	sprite3 := [5]uint8{
		0xf0,
		0x10,
		0xf0,
		0x10,
		0xf0,
	}
	sprite4 := [5]uint8{
		0x90,
		0x90,
		0xf0,
		0x10,
		0x10,
	}
	sprite5 := [5]uint8{
		0xf0,
		0x80,
		0xf0,
		0x10,
		0xf0,
	}
	sprite6 := [5]uint8{
		0xf0,
		0x80,
		0xf0,
		0x90,
		0xf0,
	}
	sprite7 := [5]uint8{
		0xf0,
		0x10,
		0x20,
		0x40,
		0x40,
	}
	sprite8 := [5]uint8{
		0xf0,
		0x10,
		0xf0,
		0x80,
		0xf0,
	}
	sprite9 := [5]uint8{
		0xf0,
		0x90,
		0xf0,
		0x90,
		0xf0,
	}
	spriteA := [5]uint8{
		0xf0,
		0x90,
		0xf0,
		0x90,
		0x90,
	}
	spriteB := [5]uint8{
		0xe0,
		0x90,
		0xe0,
		0x90,
		0xe0,
	}
	spriteC := [5]uint8{
		0xf0,
		0x80,
		0x80,
		0x80,
		0xf0,
	}
	spriteD := [5]uint8{
		0xe0,
		0x90,
		0x90,
		0x90,
		0xe0,
	}
	spriteE := [5]uint8{
		0xf0,
		0x80,
		0xf0,
		0x80,
		0xf0,
	}
	spriteF := [5]uint8{
		0xf0,
		0x80,
		0xf0,
		0x80,
		0x80,
	}
	mem.AddSprite(0x100, sprite0)
	mem.AddSprite(0x105, sprite1)
	mem.AddSprite(0x10a, sprite2)
	mem.AddSprite(0x10f, sprite3)
	mem.AddSprite(0x114, sprite4)
	mem.AddSprite(0x119, sprite5)
	mem.AddSprite(0x11e, sprite6)
	mem.AddSprite(0x123, sprite7)
	mem.AddSprite(0x128, sprite8)
	mem.AddSprite(0x12d, sprite9)
	mem.AddSprite(0x132, spriteA)
	mem.AddSprite(0x137, spriteB)
	mem.AddSprite(0x13c, spriteC)
	mem.AddSprite(0x141, spriteD)
	mem.AddSprite(0x146, spriteE)
	mem.AddSprite(0x14b, spriteF)
}

func (mem *Memory) AddSprite(address uint16, sprite [5]uint8) {
	for i := uint16(0); i < 5; i++ {
		mem.ProgData[address+i] = sprite[i]
	}
}
