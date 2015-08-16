package main

import (
	"fmt"
	"github.com/unixpickle/gocube"
	"os"
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
	heuristic := NewEdgesHeuristic(5)
	
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

func Search(start gocube.CubieCube, heuristic EdgesHeuristic,
	d int) []gocube.Move {
	if d == 0 {
		if IsSolved(start) {
			return []gocube.Move{}
		} else {
			return nil
		}
	} else if heuristic.Lookup(start) > d {
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

// EdgesHeuristic associates a number of moves with many edge configurations.
type EdgesHeuristic struct {
	Mapping map[string]int
	Depth   int
}

func NewEdgesHeuristic(maxDepth int) EdgesHeuristic {
	res := EdgesHeuristic{map[string]int{}, maxDepth}
	queue := []EdgesHeuristicNode{EdgesHeuristicNode{gocube.SolvedCubieEdges(),
		HashEdges(gocube.SolvedCubieEdges()), 0}}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if _, ok := res.Mapping[node.Hash]; ok {
			continue
		}
		res.Mapping[node.Hash] = node.Depth
		if node.Depth == maxDepth {
			continue
		}
		for move := 0; move < 18; move++ {
			newEdges := node.Edges
			newEdges.Move(gocube.Move(move))
			hash := HashEdges(newEdges)
			queue = append(queue, EdgesHeuristicNode{newEdges, hash,
				node.Depth + 1})
		}
	}
	return res
}

func (e EdgesHeuristic) Lookup(state gocube.CubieCube) int {
	if res, ok := e.Mapping[HashEdges(state.Edges)]; ok {
		return res
	} else {
		return e.Depth + 1
	}
}

// EdgesHeuristicNode is used for a breadth-first search.
type EdgesHeuristicNode struct {
	Edges gocube.CubieEdges
	Hash  string
	Depth int
}

func HashEdges(e gocube.CubieEdges) string {
	var res [22]byte
	edgeLetters := []byte("ABCDEFGHIJKL")
	for i := 0; i < 11; i++ {
		if e[i].Flip {
			res[i] = 'F'
		} else {
			res[i] = 'T'
		}
	}
	for i := 0; i < 11; i++ {
		res[i+11] = edgeLetters[e[i].Piece]
	}
	return string(res[:])
}
