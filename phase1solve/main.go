package main

import (
	"fmt"
	"github.com/unixpickle/gocube"
	"os"
)

func main() {
	// Get input.
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
	
	fmt.Println("Generating data...")
	moves := gocube.NewPhase1Moves()
	heuristic := gocube.NewPhase1Heuristic(moves, false)
	
	fmt.Println("Searching...")
	solver := gocube.NewPhase1Solver(cc.Phase1Cube(), heuristic, moves)
	for solution := range solver.Solutions() {
		fmt.Println(solution.Moves, "-", len(solution.Moves), "moves")
	}
}
