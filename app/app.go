package main

import (
	"flag"
	"github.com/jayshaffer/chippy"
	"os"
)

func main() {
	go chippy.OpenScreen()
	var filename = flag.String("filename", "", "full path to CHIP-8 ROM")
	flag.Parse()
	file, error := os.Open(*filename)
	if error != nil {
		panic("Something broke!")
	}
	bytes := make([]byte, 5000)
	file.Read(bytes)
	cpu := new(chippy.CPU)
	mem := chippy.Load(bytes)
	cpu.Boot()
	cpu.Run(mem)
}
