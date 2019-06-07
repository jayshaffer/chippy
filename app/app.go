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
	window, err := sdl.CreateWindow(Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		Width, Height, sdl.WINDOW_SHOWN)
	surface, _ := window.GetSurface()
	defer window.Destroy()
	defer sdl.Quit()
	if err != nil {
		panic(err)
	}

	running := true
	keyboard := chippy.NewKeyboard(mem)
	var rects []sdl.Rect

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				keyboard.HandleKeyPress(t)
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		rects = []sdl.Rect{}

		for i := 0; i < 64; i++ {
			func(i int) {
				for j := 0; j < 32; j++ {
					if mem.DisplayMem[i][j] > 0 {
						rects = append(rects, sdl.Rect{int32(i * XScale), int32(j * YScale), XScale, YScale})
					}
				}
			}(i)
		}

		if len(rects) > 0 {
			surface.FillRect(nil, 0)
			surface.FillRects(rects, 0xffffffff)
		}
		window.UpdateSurface()
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
