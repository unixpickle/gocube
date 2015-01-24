package gocube

import (
	"testing"
)

func BenchmarkCubeSearch(b *testing.B) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	scramble, _ := ParseMoves("U D F B R")
	start := SolvedCubieCube()
	for _, move := range scramble {
		start.Move(move)
	}
	for i := 0; i < b.N; i++ {
		s := start.Search(SolveCubeGoal{}, nil, moves, 5, 0)
		for {
			if _, ok := <-s.Solutions(); !ok {
				break
			}
		}
	}
}

func BenchmarkCubieHalfTurn(b *testing.B) {
	moves, _ := ParseMoves("U2 D2 F2 B2 R2 L2")
	cubie := SolvedCubieCube()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			cubie.Move(move)
		}
	}
}

func BenchmarkCubieMove(b *testing.B) {
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	cubie := SolvedCubieCube()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			cubie.Move(move)
		}
	}
}

func BenchmarkCubieQuarterTurn(b *testing.B) {
	moves, _ := ParseMoves("B U D B' D' R' L F F' R L' U'")
	cubie := SolvedCubieCube()
	for i := 0; i < b.N/len(moves); i++ {
		for _, move := range moves {
			cubie.Move(move)
		}
	}
}

func TestCubeSearch(t *testing.T) {
	moves, _ := ParseMoves("U D F B R L U2 D2 F2 B2 R2 L2 U' D' F' B' R' L'")
	scramble, _ := ParseMoves("R U2 B D F2")
	start := SolvedCubieCube()
	for _, move := range scramble {
		start.Move(move)
	}
	s := start.Search(SolveCubeGoal{}, nil, moves, len(scramble), 0)
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
