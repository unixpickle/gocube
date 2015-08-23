package fmc

import "github.com/unixpickle/gocube"

// EdgesHeuristic associates a number of moves with many edge configurations.
type EdgesHeuristic struct {
	Mapping map[string]int
	Depth   int
}

// NewEdgesHeuristic generates an EdgesHeuristic which extends to a certain
// depth.
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

// Lookup returns a lower-bound move count for solving the edges of a cube.
func (e EdgesHeuristic) Lookup(state gocube.CubieEdges) int {
	if res, ok := e.Mapping[HashEdges(state)]; ok {
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

// HashEdges generates a hash string for CubieEdges.
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
