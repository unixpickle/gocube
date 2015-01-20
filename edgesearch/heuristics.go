package edgesearch

import (
	"github.com/unixpickle/gocube"
)

// OrientHeuristic stores the number of moves required to solve every EO
// configuration.
type OrientHeuristic [0x800]int

// MakeOrientHeuristic uses search to generate an OrientHeuristic.
// This should run pretty fast.
func MakeOrientHeuristic() *OrientHeuristic {
	var res OrientHeuristic
	// TODO: do a breadth-first seacrh here
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
			encodedEdges[i] |= (1 << uint(i))
		}
	}
	return encodedEdges
}
