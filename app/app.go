package main

import (
	"flag"
	"github.com/jayshaffer/chippy"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	Title  = "Chippy"
	Width  = 800
	Height = 600
	XScale = 10
	YScale = 10
)

func run() int {
	var filename = flag.String("filename", "", "full path to CHIP-8 ROM")
	flag.Parse()
	mem := StartCPU(filename)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		Width, Height, sdl.WINDOW_SHOWN)
	surface, err := window.GetSurface()
	surface.FillRect(nil, 0)
	if err != nil {
		panic(err)
	}
	running := true
	keyboard := chippy.NewKeyboard(mem)
	for running {
		for i := 0; i < 64; i++ {
			go func(i int) {
				for j := 0; j < 32; j++ {
					if mem.DisplayMem[i][j] > 0 {
						rect := sdl.Rect{int32(i * XScale), int32(j * YScale), XScale, YScale}
						surface.FillRect(&rect, 0xffffffff)
					}
				}
			}(i)
		}
		window.UpdateSurface()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				keyboard.HandleKeyPress(t)
			case *sdl.QuitEvent:
				running = false
				break
			}
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}

func StartCPU(filename *string) (memory *chippy.Memory) {
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
	return mem
}
