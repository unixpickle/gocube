package gocube

import (
	"math/rand"
	"testing"
)

func BenchmarkPhase2HeuristicComplete(b *testing.B) {
	moves := NewPhase2Moves()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPhase2Heuristic(moves, true)
	}
}

func BenchmarkPhase2HeuristicIncomplete(b *testing.B) {
	moves := NewPhase2Moves()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPhase2Heuristic(moves, false)
	}
}

func BenchmarkPhase2Move(b *testing.B) {
	moves := NewPhase2Moves()
	cube := SolvedPhase2Cube()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cube.Move(Phase2Move(i%10), moves)
	}
}

func BenchmarkSolvePhase2(b *testing.B) {
	moves := NewPhase2Moves()
	heuristic := NewPhase2Heuristic(moves, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cube := SolvedPhase2Cube()
		for j := 0; j < 30; j++ {
			cube.Move(Phase2Move(rand.Intn(10)), moves)
		}
		b.StartTimer()
		SolvePhase2(cube, 18, heuristic, moves)
	}
}

func TestSolvePhase2(t *testing.T) {
	table := NewPhase2Moves()
	heuristic := NewPhase2Heuristic(table, false)

	// Do a bunch of random move sequences and make sure a solution is found.
	for length := 1; length <= 18; length++ {
		for i := 0; i < 10; i++ {
			cube := SolvedPhase2Cube()
			moves := make([]Phase2Move, length)
			for j := 0; j < length; j++ {
				move := Phase2Move(rand.Intn(10))
				cube.Move(move, table)
				moves[j] = move
			}
			solution := SolvePhase2(cube, 18, heuristic, table)

			if solution == nil {
				t.Error("No solution found for:", moves)
				continue
			}

			// Make sure the solution is short enough.
			if len(solution) > len(moves) {
				t.Error("Solution is too long:", solution, "for scramble",
					moves)
			}

			// Make sure the solution actually works.
			for _, m := range solution {
				cube.Move(m, table)
			}
			if !cube.Solved() {
				t.Error("Solution", solution, "did not work for scramble",
					moves)
			}
		}
	}
}
