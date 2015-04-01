package gocube

import (
	"testing"
)

func BenchmarkCubieEdgesMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	edges := SolvedCubieEdges()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			edges.Move(move)
		}
	}
}

func TestCubieEdges(t *testing.T) {
	edges := SolvedCubieEdges()

	// Perform a full scramble on the edges.
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		edges.Move(move)
	}

	// Verify that the edges are the way they should be.
	pieces := []int{9, 4, 5, 1, 11, 6, 0, 10, 8, 7, 3, 2}
	for i, edge := range edges {
		if edge.Piece != pieces[i] {
			t.Error("Invalid edge piece at index", i)
		}
	}

	// Verify the orientation map.
	orientations := []bool{true, true, false, false, false, false, true, false,
		true, true, false, true}
	for i, edge := range edges {
		if edge.Flip != orientations[i] {
			t.Error("Invalid orientation at", i)
		}
	}
}
