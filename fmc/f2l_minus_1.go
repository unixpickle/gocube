package fmc

import "github.com/unixpickle/gocube"

// IsF2LMinus1Solved returns true if any F2L-1 is solved. It returns the face of
// the F2L-1 cross and the corner index of the pair which is not solved.
func IsF2LMinus1Solved(state gocube.CubieCube) (solved bool, face, corner int) {
	crossEdges := []int{
		0, 4, 5, 6,
		2, 8, 10, 11,
		0, 1, 2, 3,
		6, 7, 8, 9,
		1, 5, 7, 11,
		3, 4, 9, 10,
	}

	pairEdges := []int{
		1, 3, 7, 9,
		1, 3, 7, 9,
		4, 5, 10, 11,
		4, 5, 10, 11,
		0, 2, 6, 8,
		0, 2, 6, 8,
	}

	pairCorners := []int{
		7, 6, 3, 2,
		5, 4, 1, 0,
		6, 7, 4, 5,
		2, 3, 0, 1,
		7, 5, 3, 1,
		6, 4, 2, 0,
	}

	var edgesSolved [12]bool
	for i := 0; i < 12; i++ {
		if state.Edges[i].Piece == i && !state.Edges[i].Flip {
			edgesSolved[i] = true
		}
	}

	var cornersSolved [8]bool
	for i := 0; i < 8; i++ {
		if state.Corners[i].Piece == i && state.Corners[i].Orientation == 1 {
			cornersSolved[i] = true
		}
	}

	solved = true

CrossLoop:
	for face = 1; face <= 6; face++ {
		faceStartIndex := (face - 1) * 4

		for i := 0; i < 4; i++ {
			if !edgesSolved[crossEdges[faceStartIndex+i]] {
				continue CrossLoop
			}
		}

		corner = -1
		for i := 0; i < 4; i++ {
			if !cornersSolved[pairCorners[faceStartIndex+i]] ||
				!edgesSolved[pairEdges[faceStartIndex+i]] {
				if corner == -1 {
					corner = pairCorners[faceStartIndex+i]
				} else {
					continue CrossLoop
				}
			}
		}

		return
	}

	return false, -1, -1
}
