package gocube

import (
	"testing"
)

func TestParseMove(t *testing.T) {
	moves := map[string]Move{
		"U":  NewMove(1, 1),
		"D":  NewMove(2, 1),
		"F":  NewMove(3, 1),
		"B":  NewMove(4, 1),
		"R":  NewMove(5, 1),
		"L":  NewMove(6, 1),
		"U'": NewMove(1, -1),
		"D'": NewMove(2, -1),
		"F'": NewMove(3, -1),
		"B'": NewMove(4, -1),
		"R'": NewMove(5, -1),
		"L'": NewMove(6, -1),
		"U2": NewMove(1, 2),
		"D2": NewMove(2, 2),
		"F2": NewMove(3, 2),
		"B2": NewMove(4, 2),
		"R2": NewMove(5, 2),
		"L2": NewMove(6, 2),
	}
	for str, expect := range moves {
		move, err := ParseMove(str)
		if err != nil {
			t.Error(err)
		} else if expect != move {
			t.Error("Unexpected result for:", str)
		}
	}
	invalid := []string{"R23", "d2", "F3", "B 2", ""}
	for _, str := range invalid {
		if _, err := ParseMove(str); err == nil {
			t.Error("Expected error (but didn't get one) for:", str)
		}
	}
}

func TestParseMoves(t *testing.T) {
	movesString := "R2 B D' F"
	parsed, err := ParseMoves(movesString)
	if err != nil {
		t.Error(err)
		return
	}
	moves := []Move{NewMove(5, 2), NewMove(4, 1), NewMove(2, -1), NewMove(3, 1)}
	if len(parsed) != len(moves) {
		t.Error("Invalid move count for:", movesString)
		return
	}
	for i, x := range parsed {
		if x != moves[i] {
			t.Error("Invalid move at index", i)
		}
	}
}
