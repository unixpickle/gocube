package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/gocube"
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

	fmt.Println("Solving...")
	solver := gocube.NewSolver(*cc, 30)
	for solution := range solver.Solutions() {
		fmt.Println("Solution:", solution, "-", len(solution), "moves")
	}
}
