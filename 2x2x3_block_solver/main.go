package main

import (
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/fmc"
	"os"
)

func main() {
	sc, err := gocube.InputStickerCube()
	if err != nil {
		fmt.Println("Failed to read stickers:", err)
		os.Exit(1)
	}
	cc, err := sc.CubieCube()
	if err != nil {
		fmt.Println("Invalid stickers:", err)
		os.Exit(1)
	}
	solutions := fmc.TwoStep2x2x3(*cc)
	for solution := range solutions {
		fmt.Println("Solution (", len(solution), "): ", solution)
	}
}
