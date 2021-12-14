package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/fmc"
)

var thirdF2LCorner int = -1

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
	heuristic := fmc.NewEdgesHeuristic(5)

	if cc.Corners[1].Piece == 1 && cc.Corners[1].Orientation == 1 {
		thirdF2LCorner = 1
	} else if cc.Corners[5].Piece == 5 && cc.Corners[5].Orientation == 1 {
		thirdF2LCorner = 5
	}

	for depth := 0; depth <= 20; depth++ {
		fmt.Println("Searching depth", depth)
		if solution := Search(*cc, heuristic, depth); solution != nil {
			fmt.Println("Got a solution:", solution)
			break
		}
	}
}

func Search(start gocube.CubieCube, heuristic fmc.EdgesHeuristic, d int) []gocube.Move {
	if d == 0 {
		if IsSolved(start) {
			return []gocube.Move{}
		} else {
			return nil
		}
	} else if heuristic.Lookup(start.Edges) > d {
		return nil
	}
	for m := 0; m < 18; m++ {
		newCube := start
		newCube.Move(gocube.Move(m))
		if solution := Search(newCube, heuristic, d-1); solution != nil {
			return append([]gocube.Move{gocube.Move(m)}, solution...)
		}
	}
	return nil
}

func IsSolved(state gocube.CubieCube) bool {
	if !state.Edges.Solved() {
		return false
	}
	isCornerSolved := func(idx int) bool {
		return state.Corners[idx].Piece == idx &&
			state.Corners[idx].Orientation == 1
	}
	if thirdF2LCorner < 0 {
		if !isCornerSolved(0) || !isCornerSolved(4) {
			return false
		}
		if !isCornerSolved(1) && !isCornerSolved(5) {
			return false
		}
	} else if !isCornerSolved(thirdF2LCorner) {
		return false
	}
	return true
}
