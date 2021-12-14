package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/fmc"
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
	solutions := fmc.ThreeStepF2LMinus1(*cc)
	for solution := range solutions {
		fmt.Println("Solution (", len(solution), "): ", solution)
	}
}
