package gocube

import (
	"errors"
	"strconv"
)

// CornerIndexes contains 8 sets of 3 values which corresponds to the x, y, and
// z sticker indexes for each corner piece.
var CornerIndexes = []int{
	51, 15, 35,
	44, 17, 33,
	45, 0, 29,
	38, 2, 27,
	53, 9, 24,
	42, 11, 26,
	47, 6, 18,
	36, 20, 8,
}

// CornerPieces contains 8 sets of 3 values which correspond to the x, y, and
// z stickers for each corner piece.
var CornerPieces = []int{
	6, 2, 4,
	5, 2, 4,
	6, 1, 4,
	5, 1, 4,
	6, 2, 3,
	5, 2, 3,
	6, 1, 3,
	5, 1, 3,
	
}

// EdgeIndexes contains 12 pairs of values which correspond to the sticker
// indexes of each edge.
var EdgeIndexes = []int{
	7, 19,
	23, 39,
	10, 25,
	21, 50,
	3, 46,
	5, 37,
	1, 28,
	30, 41,
	16, 34,
	32, 48,
	12, 52,
	14, 43,
}

// EdgePieces contains 12 pairs of values which correspond to the stickers of
// each edge.
var EdgePieces = []int{
	1, 3,
	3, 5,
	2, 3,
	3, 6,
	1, 6,
	1, 5,
	1, 4,
	4, 5,
	2, 4,
	4, 6,
	2, 6,
	2, 5,
}

// CubieToSticker converts a CubieCube to a StickerCube
func CubieToSticker(c CubieCube) StickerCube {
	var res StickerCube
	// TODO: this
	return res
}

// StickerToCubie converts a StickerCube to a CubieCube.
func StickerToCubie(s StickerCube) (*CubieCube, error) {
	var result CubieCube

	// Translate corner pieces.
	for i := 0; i < 8; i++ {
		idx := i * 3
		stickers := [3]int{s[CornerIndexes[idx]], s[CornerIndexes[idx+1]],
			s[CornerIndexes[idx+2]]}
		piece, orientation, err := findCorner(stickers)
		if err != nil {
			return nil, err
		}
		result.Corners[i].Piece = piece
		result.Corners[i].Orientation = orientation
	}

	// Translate edge pieces.
	for i := 0; i < 12; i++ {
		idx := i * 2
		stickers := [2]int{s[EdgeIndexes[idx]], s[EdgeIndexes[idx+1]]}
		piece, flip, err := findEdge(stickers)
		if err != nil {
			return nil, err
		}
		result.Edges[i].Piece = piece
		result.Edges[i].Flip = flip
	}

	return &result, nil
}

// findCorner finds the physical corner given its three colors.
func findCorner(stickers [3]int) (idx int, orientation int, err error) {
	for i := 0; i < 8; i++ {
		start := i * 3
		if !setsEqual(stickers[:], CornerPieces[start:start+3]) {
			continue
		}
		orientation = listIndex(stickers[:], 5)
		if orientation == -1 {
			orientation = listIndex(stickers[:], 6)
		}
		return i, orientation, nil
	}
	return 0, 0, errors.New("Unrecognized corner: " +
		strconv.Itoa(stickers[0]) + "," + strconv.Itoa(stickers[1]) + "," +
		strconv.Itoa(stickers[2]))
}

// findEdge finds the physical edge given its two colors.
func findEdge(stickers [2]int) (idx int, flip bool, err error) {
	for i := 0; i < 12; i++ {
		start := i * 2
		if !setsEqual(stickers[:], EdgePieces[start:start+2]) {
			continue
		}

		// Using the EO rules, we can tell if the edge is good or bad.
		flip = false
		if stickers[1] == 1 || stickers[1] == 2 {
			// Top/bottom color in the wrong direction.
			flip = true
		} else if stickers[1] == 3 || stickers[1] == 4 {
			if stickers[0] != 1 && stickers[0] != 2 {
				// E-Slice edge with left/right color facing wrong direction.
				flip = true
			}
		}
		return i, flip, nil
	}
	return 0, false, errors.New("Unrecognized edge: " +
		strconv.Itoa(stickers[0]) + "," + strconv.Itoa(stickers[1]))
}

func listContains(list []int, num int) bool {
	for _, x := range list {
		if x == num {
			return true
		}
	}
	return false
}

func listIndex(list []int, num int) int {
	for i, x := range list {
		if x == num {
			return i
		}
	}
	return -1
}

func setsEqual(set1 []int, set2 []int) bool {
	if len(set1) != len(set2) {
		return false
	}
	for _, x := range set1 {
		if !listContains(set2, x) {
			return false
		}
	}
	return true
}
