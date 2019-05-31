package chippy

type Memory struct {
	ProgData map[uint16]uint8
}

func Load(bytes []byte) *Memory {
	mem := new(Memory)
	mem.ProgData = make(map[uint16]uint8, 3895)
	var index uint16
	index = 0x200
	for i := 0; i < len(bytes); i++ {
		mem.ProgData[index] = bytes[i]
		index++
	}
	return mem
}
