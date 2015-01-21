package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/args"
	"os"
)

func main() {
	a := args.NewArgs(flag.CommandLine)
	flag.Parse()
	scramble, err := a.Scramble()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	goal := gocube.SolveCubeGoal{}
	moves, _ := gocube.ParseMoves("R U L D R' U' L' D' R2 U2 L2 D2 F B F' B' " +
		"F2 B2")
	for i := 0; i < 20; i++ {
		fmt.Println("Exploring depth", i, "...")
		search := scramble.Search(goal, nil, moves, i, 1)
		solution, ok := <-search.Solutions()
		if ok {
			fmt.Println("Found solution:", solution)
			search.Cancel()
			break
		}
	}
}
