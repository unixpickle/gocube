package gocube

import "testing"

func TestNewPhase2Cube(t *testing.T) {
	scrambles := []string{
		"R2 U F2 D2 L2 D' B2 R2 L2 D2 U' F2 D R2 U R2 D2",
	}
	axes := []int{1}
	states := []Phase2Cube{
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
			t.Error("state", i, "error:", err)
			continue
		}
		if p2cube != state {
			t.Error("For moves", s, "expected", state, "got", p2cube)
		}
	}
}
