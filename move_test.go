package gocube

import (
	"testing"
)

func TestParseMove(t *testing.T) {
	moves := map[string]Move{
		"U": Move{1, 1},
		"D": Move{2, 1},
		"F": Move{3, 1},
		"B": Move{4, 1},
		"R": Move{5, 1},
		"L": Move{6, 1},
		"U'": Move{1, -1},
		"D'": Move{2, -1},
		"F'": Move{3, -1},
		"B'": Move{4, -1},
		"R'": Move{5, -1},
		"L'": Move{6, -1},
		"U2": Move{1, 2},
		"D2": Move{2, 2},
		"F2": Move{3, 2},
		"B2": Move{4, 2},
		"R2": Move{5, 2},
		"L2": Move{6, 2},
	}
	for str, expect := range moves {
		move, err := ParseMove(str)
		if err != nil {
			t.Error(err)
		} else if expect.Face != move.Face || expect.Turns != move.Turns {
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
	moves := []Move{Move{5, 2}, Move{4, 1}, Move{2, -1}, Move{3, 1}}
	if len(parsed) != len(moves) {
		t.Error("Invalid move count for:", movesString)
		return
	}
	for i, x := range parsed {
		if x.Turns != moves[i].Turns || x.Face != moves[i].Face {
			t.Error("Invalid move at index", i)
		}
	}
}
