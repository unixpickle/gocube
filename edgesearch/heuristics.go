package edgesearch

import (
	"github.com/unixpickle/gocube"
)

// OrientHeuristic stores the number of moves required to solve every EO
// configuration.
type OrientHeuristic [0x800]int

// MakeOrientHeuristic uses breadth-first search to generate an OrientHeuristic.
// This should run pretty fast.
func MakeOrientHeuristic(moves []gocube.Move) *OrientHeuristic {
	var res OrientHeuristic
	
	for i := 0; i < 0x800; i++ {
		res[i] = -1
	}
	
	// Perform a very basic breadth-first search.
	nodes := []searchNode{searchNode{gocube.SolvedCubieEdges(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := encodeEdges(node.State)
		
		if res[idx] >= 0 {
			// Node was already visited
			continue
		}
		
		// Expand the node
		res[idx] = node.Depth
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, searchNode{state, node.Depth+1})
		}
	}
	
	return &res
}

// MinMoves returns the number of moves required to orient all the edges.
func (o *OrientHeuristic) MinMoves(e gocube.CubieEdges) int {
	return o[encodeEdges(e)]
}

func encodeEdges(e gocube.CubieEdges) int {
	encodedEdges := 0
	for i := 0; i < 11; i++ {
		if e[i].Flip {
			encodedEdges |= (1 << uint(i))
		}
	}
	return encodedEdges
}

type searchNode struct {
	State gocube.CubieEdges
	Depth int
}
