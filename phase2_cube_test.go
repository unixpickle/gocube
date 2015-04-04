package gocube

import "testing"

func BenchmarkNewPhase2Cube(b *testing.B) {
	scramble, _ := ParseMoves("U2 L F2 R2 D2 R' B2 U2 D2 R2 L' F2 R U2 L U2 R2")
	cube := SolvedCubieCube()
	for _, m := range scramble {
		cube.Move(m)
	}
	b.ResetTimer()

	// Do half the calls on the first cube.
	for i := 0; i < b.N/2; i++ {
		NewPhase2Cube(cube, 0)
	}

	// Apply some moves to get a different state
	cube.Move(scramble[0])
	cube.Move(scramble[1])
	cube.Move(scramble[2])
	cube.Move(scramble[5])
	cube.Move(scramble[4])
	cube.Move(scramble[3])

	for i := 0; i < b.N/2; i++ {
		NewPhase2Cube(cube, 0)
	}
}

func BenchmarkNewPhase2Moves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPhase2Moves()
	}
}

func TestNewPhase2Cube(t *testing.T) {
	// I did the same scramble, translated for all three directions. It seems to
	// work for all of them.
	scrambles := []string{
		"R2 U F2 D2 L2 D' B2 R2 L2 D2 U' F2 D R2 U R2 D2",
		"U2 L F2 R2 D2 R' B2 U2 D2 R2 L' F2 R U2 L U2 R2",
		"R2 F D2 B2 L2 B' U2 R2 L2 B2 F' D2 B R2 F R2 B2",
	}
	axes := []int{1, 0, 2}
	states := []Phase2Cube{
		Phase2Cube{29024, 14092, 2},
		Phase2Cube{29024, 14092, 2},
		Phase2Cube{29024, 14092, 2},
	}

	for i, s := range scrambles {
		scramble, _ := ParseMoves(s)
		state := states[i]
		cube := SolvedCubieCube()
		for _, m := range scramble {
			cube.Move(m)
		}
		p2cube, err := NewPhase2Cube(cube, axes[i])
		if err != nil {
			t.Error("For moves", s, "got error:", err)
			continue
		}
		if p2cube != state {
			t.Error("For moves", s, "expected", state, "got", p2cube)
		}
	}

	// Make sure that an error occurs if the thing is wrong.
	normalCube := SolvedCubieCube()
	moves, _ := ParseMoves("U2 L2 D U B' U2 B' R2 F2 L D2 L' B F2 L B' L' R")
	for _, m := range moves {
		normalCube.Move(m)
	}
	for axis := 0; axis < 3; axis++ {
		result, err := NewPhase2Cube(normalCube, axis)
		if err == nil {
			t.Error("Expected error but got", result)
		}
	}
}

func TestPhase2Moves(t *testing.T) {
	// Do the algorithm "R2 U F2 D2 L2 D' B2 R2 L2 D2 U' F2 D R2 U R2 D2"
	moves := []Phase2Move{2, 4, 0, 9, 3, 8, 1, 2, 3, 9, 5, 0, 7, 2, 4, 2, 9}
	table := NewPhase2Moves()
	state := SolvedPhase2Cube()
	for _, m := range moves {
		state.Move(m, table)
	}
	c := Phase2Cube{29024, 14092, 2}
	if state != c {
		t.Error("Invalid state", state, "after doing R2 U F2 D2 L2 D' B2 "+
			"R2 L2 D2 U' F2 D R2 U R2 D2")
	}
}
