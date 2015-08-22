package fmc

import "github.com/unixpickle/gocube"

// IsAllButL5CSolved returns true if the cube is almost solved (having up to 5
// unsolved corners).
func IsAllButL5CSolved(cube gocube.CubieCube) bool {
	if !cube.Edges.Solved() {
		return false
	}
	var unsolvedCorners int
	for i := 0; i < 8; i++ {
		if cube.Corners[i].Piece != i || cube.Corners[i].Orientation != 1 {
			unsolvedCorners++
			if unsolvedCorners > 5 {
				return false
			}
		}
	}
	return true
}

// FourStepAllButL5C finds solutions to everything but the last 5 corners by
// solving the F2L-1 and going from there.
func FourStepAllButL5C(cube gocube.CubieCube) <-chan []gocube.Move {
	channel := make(chan []gocube.Move, 1)
	go func() {
		for f2l1Solution := range ThreeStepF2LMinus1(cube) {
			start := cube
			for _, move := range f2l1Solution {
				start.Move(move)
			}
			lastFace := -1
			if len(f2l1Solution) > 0 {
				lastFace = f2l1Solution[len(f2l1Solution)-1].Face()
			}
			if moves := iterativeAllButL5C(start, lastFace); moves != nil {
				channel <- append(f2l1Solution, moves...)
			} else {
				// TODO: here, use a 2-gen solution with unlimited depth.
			}
		}
	}()
	return channel
}

func iterativeAllButL5C(start gocube.CubieCube, lastFace int) []gocube.Move {
	for depth := 0; depth < 7; depth++ {
		if solution := solveAllButL5C(start, depth, lastFace); solution != nil {
			return solution
		}
	}
	return nil
}

func solveAllButL5C(start gocube.CubieCube, depth, lastFace int) []gocube.Move {
	if depth == 0 {
		if solved := IsAllButL5CSolved(start); solved {
			return []gocube.Move{}
		} else {
			return nil
		}
	}
	for m := 0; m < 18; m++ {
		move := gocube.Move(m)
		face := move.Face()
		if face == lastFace {
			continue
		}
		newCube := start
		newCube.Move(move)
		if solution := solveAllButL5C(newCube, depth-1, face); solution != nil {
			return append([]gocube.Move{move}, solution...)
		}
	}
	return nil
}
