package gocube

import "testing"

func TestSolverRandomCubes(t *testing.T) {
	p1Moves := NewPhase1Moves()
	p1Heuristic := NewPhase1Heuristic(p1Moves)
	p2Moves := NewPhase2Moves()
	p2Heuristic := NewPhase2Heuristic(p2Moves, false)
	tables := SolverTables{
		P1Heuristic: p1Heuristic,
		P1Moves:     p1Moves,
		P2Heuristic: p2Heuristic,
		P2Moves:     p2Moves,
	}
	for i := 0; i < 20; i++ {
		cube := RandomCubieCube()
		solver := NewSolverTables(cube, 24, tables)
		solution := <-solver.Solutions()
		solver.Stop()
		for _, m := range solution {
			cube.Move(m)
		}
		if !cube.Solved() {
			t.Errorf("%d: solution %v resulted in cube %v", i, solution, cube)
		}
	}
}
