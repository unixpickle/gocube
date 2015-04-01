package gocube

import (
	"errors"
	"strings"
)

// A Move represents a face turn.
//
// A move can occur on the faces U, D, F, B, R, and L. These are the first 6
// values of the Move type. The next 6 values are U', D', F', B', R', L'. The
// final six values are U2, D2, F2, B2, R2, L2. Thus, there are a total of 18
// possible moves in the range [0, 18).
type Move int

// NewMove creates a new move with a face in the range [1, 6] and a number of
// turns 1, -1, or 2.
func NewMove(face int, turns int) Move {
	if turns == 1 {
		return Move(face - 1)
	} else if turns == -1 {
		return Move(face + 5)
	} else {
		return Move(face + 11)
	}
}

// Face returns the face of the move, which is a number from 1 through 6
// indicating U, D, F, B, R and L respectively.
func (m Move) Face() int {
	return (int(m)%6) + 1
}

// Turns returns 1, -1, or 2 to indicate the number of times the face is turned.
func (m Move) Turns() int {
	res := int(m)/6
	if res == 0 {
		return 1
	} else if res == 1 {
		return -1
	} else {
		return 2
	}
}

// String converts a move to a WCA-notation string.
func (m Move) String() string {
	letter := string(" UDFBRL"[m.Face()])
	if m < 6 {
		return letter
	} else if m < 12 {
		return letter + "'"
	} else {
		return letter + "2"
	}
}

// ParseMove parses a move in WCA notation and returns it.
func ParseMove(moveString string) (Move, error) {
	if len(moveString) == 1 {
		moves := map[byte]Move{'U': 0, 'D': 1, 'F': 2, 'B': 3, 'R': 4, 'L': 5}
		if move, ok := moves[moveString[0]]; ok {
			return move, nil
		}
	} else if len(moveString) == 2 {
		result, err := ParseMove(moveString[0:1])
		if err != nil {
			return 0, errors.New("invalid move: " + moveString)
		} else if moveString[1] == '2' {
			return result + 12, nil
		} else if moveString[1] == '\'' {
			return result + 6, nil
		}
	}
	return 0, errors.New("invalid move: " + moveString)
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

