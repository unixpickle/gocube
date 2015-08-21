package fmc

import "github.com/unixpickle/gocube"

// Is2x2x3Solved returns whether or not any 2x2x3 blocks are solved. If a block
// is solved, its corresponding cross edge index is also returned.
func Is2x2x3Solved(state gocube.CubieCube) (solved bool, edgeIndex int) {
	var edgesSolved [12]bool
	for i := 0; i < 12; i++ {
		if state.Edges[i].Piece == i && !state.Edges[i].Flip {
			edgesSolved[i] = true
		}
	}

	var cornersSolved [8]bool
	for i := 0; i < 8; i++ {
		if state.Corners[i].Piece == i && state.Corners[i].Orientation == 1 {
			cornersSolved[i] = true
		}
	}

	cornersForBlocks := []int{
		6, 7,
		5, 7,
		4, 5,
		4, 6,
		2, 6,
		3, 7,
		2, 3,
		1, 3,
		0, 1,
		0, 2,
		0, 4,
		1, 5,
	}

	edgesForBlocks := []int{
		1, 3, 4, 5,
		0, 2, 5, 11,
		1, 3, 10, 11,
		0, 2, 4, 10,
		0, 6, 3, 9,
		0, 6, 1, 7,
		4, 5, 7, 9,
		5, 11, 6, 8,
		7, 9, 10, 11,
		4, 10, 6, 8,
		2, 8, 3, 9,
		1, 7, 2, 8,
	}

OuterLoop:
	for i := 0; i < 12; i++ {
		if edgesSolved[i] {
			if !cornersSolved[cornersForBlocks[i*2]] ||
				!cornersSolved[cornersForBlocks[i*2+1]] {
				continue OuterLoop
			}
			if !edgesSolved[edgesForBlocks[i*4]] ||
				!edgesSolved[edgesForBlocks[i*4+1]] ||
				!edgesSolved[edgesForBlocks[i*4+2]] ||
				!edgesSolved[edgesForBlocks[i*4+3]] {
				continue OuterLoop
			}
			return true, i
		}
	}

	return false, -1
}

// TwoStep2x2x3 finds solutions to 2x2x3 blocks by first solving 2x2x2 blocks
// and going from there.
func TwoStep2x2x3(cube gocube.CubieCube) <-chan []gocube.Move {
	channel := make(chan []gocube.Move, 1)
	go func() {
		for smallBlock := range Solve2x2x2(cube) {
			start := cube
			for _, move := range smallBlock {
				start.Move(move)
			}
			lastFace := -1
			if len(smallBlock) > 0 {
				lastFace = smallBlock[len(smallBlock)-1].Face()
			}
			moves := iterative2x2x3(start, lastFace)
			channel <- append(smallBlock, moves...)
		}
	}()
	return channel
}

func iterative2x2x3(start gocube.CubieCube, lastFace int) []gocube.Move {
	for depth := 0; true; depth++ {
		if solution := solve2x2x3(start, depth, lastFace); solution != nil {
			return solution
		}
	}
	return nil
}

func solve2x2x3(start gocube.CubieCube, depth, lastFace int) []gocube.Move {
	if depth == 0 {
		if solved, _ := Is2x2x3Solved(start); solved {
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
		if solution := solve2x2x3(newCube, depth-1, face); solution != nil {
			return append([]gocube.Move{move}, solution...)
		}
	}
	return nil
}
