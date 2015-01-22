package gocube

import (
	"testing"
)

func BenchmarkMakeEOPruner(b *testing.B) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	for i := 0; i < b.N; i++ {
		MakeEOPruner(moves)
	}
}

func BenchmarkEdgeSearch(b *testing.B) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	scramble, _ := ParseMoves("U D F B R")
	start := SolvedCubieEdges()
	for _, move := range scramble {
		start.Move(move)
	}
	for i := 0; i < b.N; i++ {
		s := start.Search(SolveEdgesGoal{}, nil, moves, 5, 0)
		for {
			if _, ok := <-s.Solutions(); !ok {
				break
			}
		}
	}
}

