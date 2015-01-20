package twophase

import (
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/edgesearch"
)

// A PhaseOneHeuristic reports a good estimate as to the number of moves to
// complete the first phase of Kociemba's two-phase algorithm.
type PhaseOneHeuristic struct {
	Edges *edgesearch.PhaseOneHeuristic
}

// MakePhaseOneHeuristic uses search to generate a PhaseOneHeuristic.
func MakePhaseOneHeuristic(moves []gocube.Move) *PhaseOneHeuristic {
	return &PhaseOneHeuristic{edgesearch.MakePhaseOneHeuristic(moves)}
}

// MinMoves returns a lower-bound for completing the first phase of Kociemba's
// two-phase algorithm.
func (h *PhaseOneHeuristic) MinMoves(c gocube.CubieCube) int {
	return h.Edges.MinMoves(c.Edges)
}
