package gocube

import (
	"errors"
	"strings"
)

// A Move stores a single move on the cube.
//
// The face of a move can be 1 through 6 for U, D, F, B, R, and L respectively.
// In some cases, the face can be 7 (M), 8 (E), or 9 (S) as well.
//
// The turns can be 1, -1, or 2 to indicate the number of times to turn the
// face.
type Move struct {
	Face  int
	Turns int
}

// ParseMove parses a move in WCA notation and returns it.
func ParseMove(moveString string) (Move, error) {
	if len(moveString) == 1 {
		faces := map[byte]int{'U': 1, 'D': 2, 'F': 3, 'B': 4, 'R': 5, 'L': 6}
		if face, ok := faces[moveString[0]]; ok {
			return Move{face, 1}, nil
		}
	} else if len(moveString) == 2 {
		result, err := ParseMove(moveString[0:1])
		if err != nil {
			return Move{}, errors.New("Invalid move: " + moveString)
		} else if moveString[1] == '2' {
			result.Turns = 2
			return result, nil
		} else if moveString[1] == '\'' {
			result.Turns = -1
			return result, nil
		}
	}
	return Move{}, errors.New("Invalid move: " + moveString)
}

// ParseMoves parses a space-delimited list of WCA moves.
func ParseMoves(movesString string) ([]Move, error) {
	parts := strings.Split(movesString, " ")
	res := make([]Move, len(parts))
	for i, s := range parts {
		if move, err := ParseMove(s); err != nil {
			return nil, err
		} else {
			res[i] = move
		}
	}
	return res, nil
}
