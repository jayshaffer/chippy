package main

import (
	"flag"
	"github.com/jayshaffer/chippy"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const x_scale = 8
const y_scale = 6

func main() {
	var filename = flag.String("filename", "", "full path to CHIP-8 ROM")
	flag.Parse()
	file, error := os.Open(*filename)
	if error != nil {
		panic(error)
	}
	bytes := make([]byte, 5000)
	file.Read(bytes)
	cpu := new(chippy.CPU)
	mem := chippy.Load(bytes)
	cpu.Boot(mem)
	go cpu.Run()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Chippy", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	surface, err := window.GetSurface()
	surface.FillRect(nil, 0)
	if err != nil {
		panic(err)
	}
	running := true
	for running {
		for i := 0; i < 64; i++ {
			for j := 0; j < 32; j++ {
				if mem.DisplayMem[i][j] > 0 {
					rect := sdl.Rect{int32(i * x_scale), int32(j * y_scale), x_scale, y_scale}
					surface.FillRect(&rect, 0xffff0000)
				}
				window.UpdateSurface()
			}
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}
	}
}
