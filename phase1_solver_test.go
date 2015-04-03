package gocube

import (
	"math/rand"
	"testing"
)

func BenchmarkNewPhase1HeuristicComplete(b *testing.B) {
	moves := NewPhase1Moves()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPhase1Heuristic(moves, true)
	}
}

func BenchmarkNewPhase1HeuristicIncomplete(b *testing.B) {
	moves := NewPhase1Moves()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPhase1Heuristic(moves, false)
	}
}

func TestPhase1Heuristic(t *testing.T) {
	table := NewPhase1Moves()
	heuristic := NewPhase1Heuristic(table, false)
	
	// Do random move sequences and ensure that the lower bound is never too
	// high.
	for length := 1; length < 12; length++ {
		for i := 0; i < 50; i++ {
			cube := SolvedPhase1Cube()
			moves := make([]Move, length)
			for j := 0; j < length; j++ {
				move := Move(rand.Intn(18))
				cube.Move(move, table)
				moves[j] = move
			}
			if heuristic.LowerBound(&cube) > length {
				t.Error("Invalid lower bound", heuristic.LowerBound(&cube),
					"for moves", moves)
			}
		}
	}
}
