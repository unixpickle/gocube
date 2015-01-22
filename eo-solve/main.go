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
	if err := a.Parse(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	scramble := a.Scramble()

	// Setup search variables.
	moves := a.Moves()
	goal := gocube.EOLineGoal{}
	start := scramble.Edges
	fmt.Println("Generating heuristic...")
	pruner := gocube.MakeEOPruner(moves)
	
	for i := a.MinDepth(); i <= a.MaxDepth(); i++ {
		fmt.Println("Exploring depth", i, "...")
		search := start.Search(goal, pruner, moves, i, a.Branch())
		for {
			solution, ok := <-search.Solutions()
			if !ok {
				break
			}
			fmt.Println("Found solution:", solution)
			if !a.Multiple() {
				search.Cancel()
				return
			}
		}
	}
}
