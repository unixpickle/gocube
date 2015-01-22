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
