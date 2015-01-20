package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/args"
	"github.com/unixpickle/gocube/edgesearch"
	"os"
)

func main() {
	// Input the cube.
	a := args.NewArgs(flag.CommandLine)
	flag.Parse()
	scramble, err := a.Scramble()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	// Solve that EOLine!
	moves, _ := gocube.ParseMoves("R U L D R' U' L' D' R2 U2 L2 D2 F B F' B' " +
		"F2 B2")
	
	fmt.Println("Generating heuristic...")
	heuristic := edgesearch.MakeOrientHeuristic(moves)
	
	goal := edgesearch.EOLineGoal{}
	start := scramble.Edges
	search := edgesearch.NewSearch(start, goal, heuristic, moves)
	for i := 0; i <= 9; i++ {
		fmt.Println("Exploring depth", i, "...")
		res, _ := search.Run(i, 1)
		if res != nil {
			fmt.Println("Found solution:", res)
			break
		}
	}
}
