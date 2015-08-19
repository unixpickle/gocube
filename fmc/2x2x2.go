package fmc

import "github.com/unixpickle/gocube"

// Is2x2x2Solved returns whether or not any 2x2x2 blocks are solved. If a block
// is solved, its corresponding corner index is also returned.
func Is2x2x2Solved(state gocube.CubieCube) (solved bool, cornerIndex int) {
	var edgesSolved [12]bool
	for i := 0; i < 12; i++ {
		if state.Edges[i].Piece == i && !state.Edges[i].Flip {
			edgesSolved[i] = true
		}
	}

	edgesForCorners := []int{
		8, 9, 10,
		7, 8, 11,
		4, 6, 9,
		5, 6, 7,
		2, 3, 10,
		1, 2, 11,
		0, 3, 4,
		0, 1, 5,
	}

OuterLoop:
	for i := 0; i < 8; i++ {
		if state.Corners[i].Piece == i && state.Corners[i].Orientation == 1 {
			for j := i*3; j < i*3+3; j++ {
				if !edgesSolved[edgesForCorners[j]] {
					continue OuterLoop
				}
			}
			return true, i
		}
	}

	return false, -1
}

// Solve2x2x2 finds solutions to 2x2x2 blocks.
func Solve2x2x2(cube gocube.CubieCube) <-chan []gocube.Move {
	channel := make(chan []gocube.Move, 1)
	go iterativeDeepening2x2x2(cube, channel)
	return channel
}

func iterativeDeepening2x2x2(start gocube.CubieCube,
	output chan<- []gocube.Move) {
	for depth := 0; true; depth++ {
		moves := make([]gocube.Move, 0, depth)
		solve2x2x2Search(start, depth, 0, moves, output)
	}
}

func solve2x2x2Search(start gocube.CubieCube, depth, lastFace int,
	previousMoves []gocube.Move, output chan<- []gocube.Move) {
	if depth == 0 {
		if solved, _ := Is2x2x2Solved(start); solved {
			newArray := make([]gocube.Move, len(previousMoves))
			copy(newArray, previousMoves)
			output <- newArray
		}
		return
	}
	for m := 0; m < 18; m++ {
		move := gocube.Move(m)
		face := move.Face()
		if face == lastFace {
			continue
		}
		newCube := start
		newCube.Move(move)
		previousMoves := append(previousMoves, move)
		solve2x2x2Search(newCube, depth-1, face, previousMoves, output)
		previousMoves = previousMoves[:len(previousMoves)-1]
	}
}
