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

func BenchmarkMakeEOPruner(b *testing.B) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	for i := 0; i < b.N; i++ {
		MakeEOPruner(moves)
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

func TestEdgeSearch(t *testing.T) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	scramble, _ := ParseMoves("R U2 B D F2")
	start := SolvedCubieEdges()
	for _, move := range scramble {
		start.Move(move)
	}
	s := start.Search(SolveEdgesGoal{}, nil, moves, len(scramble), 0)
	solution, ok := <-s.Solutions()
	if !ok {
		t.Fatal("Failed to find a solution.")
	}
	if len(solution) != len(scramble) {
		t.Fatal("Invalid solution found:", solution)
	}
	for i, m := range solution {
		invMe := scramble[len(scramble) - (i+1)]
		if invMe.Turns != 2 {
			invMe.Turns *= -1
		}
		if m.Turns != invMe.Turns || m.Face != invMe.Face {
			t.Fatal("Invalid solution found:", solution)
		}
	}
	for {
		solution, ok := <-s.Solutions()
		if ok {
			t.Error("Found extraneous solution:", solution)
		} else {
			break
		}
	}
}
