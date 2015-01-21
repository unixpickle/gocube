package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/args"
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

	// Setup search variables.
	moves, _ := gocube.ParseMoves("R U L D R' U' L' D' R2 U2 L2 D2 F B F' B' " +
		"F2 B2")
	goal := gocube.EOLineGoal{}
	start := scramble.Edges
	fmt.Println("Generating heuristic...")
	pruner := gocube.MakeOrientEdgesPruner(moves)
	
	for i := 0; i <= 9; i++ {
		fmt.Println("Exploring depth", i, "...")
		search := start.Search(goal, pruner, moves, i, 1)
		solution, ok := <-search.Solutions()
		if ok {
			fmt.Println("Found solution:", solution)
			search.Cancel()
			break
		}
	}
}
