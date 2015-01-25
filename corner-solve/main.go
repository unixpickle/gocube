package main

import (
	"flag"
	"fmt"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/args"
	"os"
)

type JointHeuristic struct {
	P1 gocube.CornerPermPruner
	P2 *gocube.COPruner
}

func NewJointHeuristic(m []gocube.Move) *JointHeuristic {
	return &JointHeuristic{gocube.MakeCornerPermPruner(m),
		gocube.MakeCOPruner(m)}
}

func (j *JointHeuristic) MinMoves(c gocube.CubieCorners) int {
	m1 := j.P1.MinMoves(c)
	m2 := j.P2.MinMoves(c)
	if m1 > m2 {
		return m1
	} else {
		return m2
	}
}

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
	goal := gocube.SolveCornersGoal{}
	start := scramble.Corners
	fmt.Println("Generating heuristic...")
	pruner := NewJointHeuristic(moves)
	
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
