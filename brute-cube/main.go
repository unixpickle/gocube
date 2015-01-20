package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/args"
	"github.com/unixpickle/gocube/cubiesearch"
	"os"
)

type EmptyHeuristic struct{}

func (_ EmptyHeuristic) MinMoves(c gocube.CubieCube) int {
	return 0
}

func main() {
	a := args.NewArgs(flag.CommandLine)
	flag.Parse()
	scramble, err := a.Scramble()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	goal := cubiesearch.SolvedGoal{}
	heuristic := EmptyHeuristic{}
	moves, _ := gocube.ParseMoves("R U L D R' U' L' D' R2 U2 L2 D2 F B F' B' " +
		"F2 B2")
	search := cubiesearch.NewSearch(*scramble, goal, heuristic, moves)
	for i := 0; i < 20; i++ {
		fmt.Println("Exploring depth", i, "...")
		res, _ := search.Run(i, 1)
		if res != nil {
			fmt.Println("Found solution:", res)
			break
		}
	}
}