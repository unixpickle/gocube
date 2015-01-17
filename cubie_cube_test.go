package gocube

import (
	"testing"
)

func BenchmarkCubieMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	cubie := SolvedCubieCube()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			cubie.Move(move)
		}
	}
}

func BenchmarkCubieCornersMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	corners := SolvedCubieCorners()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			corners.Move(move)
		}
	}
}

func BenchmarkCubieEdgesMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	edges := SolvedCubieEdges()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			edges.Move(move)
		}
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
