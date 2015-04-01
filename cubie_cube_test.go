package gocube

import (
	"testing"
)

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
