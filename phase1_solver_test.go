package gocube

import "testing"

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
