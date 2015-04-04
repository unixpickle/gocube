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

func BenchmarkPhase1Solver(b *testing.B) {
	moves := NewPhase1Moves()
	heuristic := NewPhase1Heuristic(moves, false)

	b.ResetTimer()

	// Do random move sequences and solve each one.
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cube := SolvedPhase1Cube()
		for j := 0; j < 50; j++ {
			cube.Move(Move(rand.Intn(18)), moves)
		}
		b.StartTimer()
		findPhase1Solution(cube, heuristic, moves)
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

func TestPhase1Solver(t *testing.T) {
	table := NewPhase1Moves()
	heuristic := NewPhase1Heuristic(table, false)

	// Do a bunch of random move sequences and make sure a solution is found.
	for length := 1; length <= 12; length++ {
		for i := 0; i < 50; i++ {
			cube := SolvedPhase1Cube()
			moves := make([]Move, length)
			for j := 0; j < length; j++ {
				move := Move(rand.Intn(18))
				cube.Move(move, table)
				moves[j] = move
			}
			solution := findPhase1Solution(cube, heuristic, table)

			// Make sure the solution is short enough.
			if len(solution) > len(moves) {
				t.Error("Solution is too long:", solution, "for scramble",
					moves)
			}

			// Make sure the solution actually works.
			for _, m := range solution {
				cube.Move(m, table)
			}
			if !cube.AnySolved() {
				t.Error("Solution", solution, "did not work for scramble",
					moves)
			}
		}
	}
}

func findPhase1Solution(c Phase1Cube, h *Phase1Heuristic,
	m *Phase1Moves) []Move {
	solver := NewPhase1Solver(c, h, m)
	res := <-solver.Solutions()
	solver.Stop()
	return res.Moves
}
