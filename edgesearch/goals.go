package edgesearch

import (
	"github.com/unixpickle/gocube"
)

// A SolvedGoal is achieved when the edges are oriented and properly positioned.
type SolvedGoal struct{}

func (_ SolvedGoal) IsGoal(c gocube.CubieEdges) bool {
	return c.Solved()
}

// An OrientGoal is achieved when the edges are oriented.
type OrientGoal struct{}

func (_ OrientGoal) IsGoal(c gocube.CubieEdges) bool {
	for _, p := range c {
		if p.Flip {
			return false
		}
	}
	return true
}

// An EOLineGoal is achieved when the edges are oriented and the DF and DB edges
// are solved.
type EOLineGoal struct{}

func (_ EOLineGoal) IsGoal(c gocube.CubieEdges) bool {
	// Make sure the edges are oriented.
	for _, p := range c {
		if p.Flip {
			return false
		}
	}

	// Make sure the line edges are solved.
	return c[2].Piece == 2 && c[8].Piece == 8
}

// A PhaseOneGoal is achieved when the edges are oriented and the E slice edges
// are in the E slice.
type PhaseOneGoal struct{}

func (_ PhaseOneGoal) IsGoal(c gocube.CubieEdges) bool {
	// Make sure the edges are oriented.
	for _, p := range c {
		if p.Flip {
			return false
		}
	}

	// Make sure the E slice edges (1, 3, 7, 9) are in the E slice.
	for _, i := range []int{1, 3, 7, 9} {
		piece := c[i].Piece
		if piece != 1 && piece != 3 && piece != 7 && piece != 9 {
			return false
		}
	}

	return true
}
