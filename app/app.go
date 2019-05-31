package main

import (
	"github.com/jayshaffer/chippy"
)

func main() {
	go chippy.OpenScreen()
	messages := make(chan chippy.ScreenBuffer)
}
