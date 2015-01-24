package gocube

import (
	"testing"
)

func BenchmarkCubieCornersMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	corners := SolvedCubieCorners()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			corners.Move(move)
		}
	}
}

func BenchmarkMakeCOPruner(b *testing.B) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	for i := 0; i < b.N; i++ {
		MakeCOPruner(moves)
	}
}

func TestCubieCorners(t *testing.T) {
	corners := SolvedCubieCorners()

	// Perform a full scramble on the corners.
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		corners.Move(move)
	}

	// Verify that the pieces are right.
	pieces := []int{5, 7, 4, 3, 0, 2, 6, 1}
	for i, corner := range corners {
		if corner.Piece != pieces[i] {
			t.Error("Invalid piece at", i)
		}
	}

	// Verify that the orientations are right.
	orientations := []int{2, 1, 1, 2, 1, 1, 1, 0}
	for i, corner := range corners {
		if corner.Orientation != orientations[i] {
			t.Error("Invalid orientation at", i)
		}
	}
}
