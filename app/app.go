package main

import (
	"flag"
	"github.com/jayshaffer/chippy"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"sync"
	"time"
)

const (
	Title  = "Chippy"
	Width  = 800
	Height = 600
	XScale = 10
	YScale = 10
	FPS    = 60
)

func run() int {
	var filename = flag.String("filename", "", "full path to CHIP-8 ROM")
	flag.Parse()
	cpu := StartCPU(filename)
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
	keyboard := chippy.NewKeyboard(cpu.PRM)
	lastRender := time.Now()
	renderSpan := (1.0 / FPS)
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				keyboard.HandleKeyPress(t)
			case *sdl.QuitEvent:
				running = false
				break
			}
			if time.Now().Sub(lastRender).Seconds()/FPS > renderSpan {
				cpu.Tick()
				Render(cpu.PRM, surface)
			}
		}

		window.UpdateSurface()
	}
	return 0
}

func main() {
	os.Exit(run())
}

func Render(mem *chippy.Memory, surface *sdl.Surface) {
	var rects []sdl.Rect
	rects = []sdl.Rect{}
	ch := make(chan sdl.Rect, 2048)

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 32; j++ {
				if mem.DisplayMem[i][j] > 0 {
					ch <- sdl.Rect{int32(i * XScale), int32(j * YScale), XScale, YScale}
				}
			}
		}(i)
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		rects = append(rects, i)
	}

	if len(rects) > 0 {
		surface.FillRect(nil, 0)
		surface.FillRects(rects, 0xffffffff)
	}
}

func StartCPU(filename *string) (cpu *chippy.CPU) {
	file, error := os.Open(*filename)
	if error != nil {
		panic(error)
	}
	bytes := make([]byte, 5000)
	file.Read(bytes)
	newCPU := new(chippy.CPU)
	mem := chippy.Load(bytes)
	newCPU.Boot(mem)
	return newCPU
}
