package gocube

import "testing"

func BenchmarkPhase1Cube(b *testing.B) {
	moves := NewPhase1Moves()
	cube := SolvedPhase1Cube()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cube.Move(Move(i%18), moves)
	}
}

func BenchmarkPhase1Moves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPhase1Moves()
	}
}

func TestPhase1Cube(t *testing.T) {
	moves := NewPhase1Moves()
	cube := SolvedPhase1Cube()

	// Apply a scramble to the phase-1 cube.
	scramble, _ := ParseMoves("L R2 B2 F2 L2 U' B2 F U R2 F' L2 R' B' F2 D2 " +
		"R U' L' R U2 F2 D U' R2 U B2 F D U")
	for _, x := range scramble {
		cube.Move(x, moves)
	}

	if cube.YCornerOrientation != 881 || cube.FBEdgeOrientation != 358 ||
		cube.ESlicePermutation != 337 {
		t.Error("Invalid Y state:", cube.YCornerOrientation,
			cube.FBEdgeOrientation, cube.ESlicePermutation)
	}

	if cube.XCornerOrientation != 1893 || cube.MSlicePermutation != 476 {
		t.Error("Invalid X state:", cube.XCornerOrientation,
			cube.MSlicePermutation)
	}

	if cube.ZCornerOrientation != 43 || cube.SSlicePermutation != 428 ||
		cube.UDEdgeOrientation != 740 {
		t.Error("Invalid Z state:", cube.ZCornerOrientation,
			cube.SSlicePermutation, cube.UDEdgeOrientation)
	}
}

func TestPhase1Moves(t *testing.T) {
	moves := NewPhase1Moves()

	// The initial state is a solved phase-1 cube.
	slice := 220
	co := 1093
	eo := 0

	// Apply a scramble to the phase-1 cube.
	scramble, _ := ParseMoves("L R2 B2 F2 L2 U' B2 F U R2 F' L2 R' B' F2 D2 " +
		"R U' L' R U2 F2 D U' R2 U B2 F D U")
	for _, x := range scramble {
		slice = moves.ESliceMoves[slice][x]
		co = moves.COMoves[co][x]
		eo = moves.EOMoves[eo][x]
	}

	// Verify that everything is the way it should be.
	goalSlice := 337
	goalEO := 358
	goalCO := 881
	if slice != goalSlice {
		t.Error("Incorrect slice: expected", goalSlice, "got", slice)
	}
	if eo != goalEO {
		t.Error("Incorrect EO: Expected", goalEO, "got", eo)
	}
	if co != goalCO {
		t.Error("Incorrect CO: Expected", goalCO, "got", co)
	}
}
