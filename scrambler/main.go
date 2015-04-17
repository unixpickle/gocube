package main

import (
	"fmt"
	"github.com/unixpickle/gocube"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: scrambler <count> [maxlen=30]")
		os.Exit(1)
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	maxLen := 30
	if len(os.Args) > 2 {
		maxLen, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	p1Moves := gocube.NewPhase1Moves()
	p1Heuristic := gocube.NewPhase1Heuristic(p1Moves)
	p2Moves := gocube.NewPhase2Moves()
	p2Heuristic := gocube.NewPhase2Heuristic(p2Moves, false)
	tables := gocube.SolverTables{p1Heuristic, p1Moves, p2Heuristic, p2Moves}
	for i := 0; i < count; i++ {
		state := gocube.RandomCubieCube()
		solver := gocube.NewSolverTables(state, maxLen, tables)
		for solution := range solver.Solutions() {
			str := fmt.Sprint(solution)
			fmt.Println(str[1 : len(str)-1])
			solver.Stop()
			break
		}
	}
}
