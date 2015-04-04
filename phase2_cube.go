package gocube

import "errors"

// inverseXCornerIndices is the inverse permutation of xCornerIndices.
var inverseXCornerIndices []int = []int{2, 0, 3, 1, 6, 4, 7, 5}

// inverseZCornerIndices is the inverse permutation of zCornerIndices.
var inverseZCornerIndices []int = []int{4, 5, 0, 1, 6, 7, 2, 3}

// A Phase2Cube represents the parts of a cube that are important for phase-2
// solving.
type Phase2Cube struct {
	// CornerPermutation represents the permutation of the corners.
	CornerPermutation int

	// EdgePermutation represents the permutation of the 8 top/bottom edges.
	EdgePermutation int

	// SlicePermutation represents the permutation of the
	SlicePermutation int
}

// NewPhase2Cube generates a Phase2Cube from a CubieCube.
// The axis argument is 0 for X axis, 1 for Y axis, or 2 for Z axis.
// If the cube is not reduced to phase-2 in the given axis, this may return an
// error, but it does not validate everything.
func NewPhase2Cube(c CubieCube, axis int) (Phase2Cube, error) {
	var res Phase2Cube
	if axis == 0 {
		res.CornerPermutation = encodeXCornerPerm(&c.Corners)
		res.EdgePermutation = encodeRLEdges(&c.Edges)
		res.SlicePermutation = encodeMSlicePerm(&c.Edges)
	} else if axis == 1 {
		res.CornerPermutation = encodeYCornerPerm(&c.Corners)
		res.EdgePermutation = encodeUDEdges(&c.Edges)
		res.SlicePermutation = encodeESlicePerm(&c.Edges)
	} else {
		res.CornerPermutation = encodeZCornerPerm(&c.Corners)
		res.EdgePermutation = encodeFBEdges(&c.Edges)
		res.SlicePermutation = encodeSSlicePerm(&c.Edges)
	}
	if res.EdgePermutation < 0 {
		return res, errors.New("invalid edge permutation")
	} else if res.SlicePermutation < 0 {
		return res, errors.New("invalid slice permutation")
	}
	return res, nil
}

// SolvedPhase2Cube returns a solved Phase2Cube.
func SolvedPhase2Cube() Phase2Cube {
	return Phase2Cube{}
}

// Move applies a move to the Phase2Cube.
func (p *Phase2Cube) Move(move Phase2Move, table *Phase2Moves) {
	p.CornerPermutation = table.CornerMoves[p.CornerPermutation][int(move)]
	p.EdgePermutation = table.EdgeMoves[p.EdgePermutation][int(move)]
	p.SlicePermutation = table.SliceMoves[p.SlicePermutation][int(move)]
}

// Solved returns true if the Phase2Cube is solved.
func (p *Phase2Cube) Solved() bool {
	return p.CornerPermutation == 0 && p.EdgePermutation == 0 &&
		p.SlicePermutation == 0
}

// Phase2Move represents a move which can be applied to a Phase2Cube. This is a
// number in the range [0, 10), corresponding to F2 B2 R2 L2 U U' U2 D D' D2
// respectively.
type Phase2Move int

// Face returns a number from [1, 6] corresponding to the face of the move if it
// were applied on the Y axis.
func (p Phase2Move) Face() int {
	return []int{2, 3, 4, 5, 0, 0, 0, 1, 1, 1}[int(p)] + 1
}

// Inverse returns the move's inverse.
func (p Phase2Move) Inverse() Phase2Move {
	return []Phase2Move{0, 1, 2, 3, 5, 4, 6, 8, 7, 9}[int(p)]
}

// Move converts the Phase2Move into a regular Move.
// The axis argument indicates the axis that the move should act on (i.e. the
// axis of the corresponding Phase2Cube). This is a number in [0, 3).
func (p Phase2Move) Move(axis int) Move {
	return [][]Move{
		[]Move{14, 15, 12, 13, 5, 11, 17, 4, 10, 16},
		[]Move{14, 15, 16, 17, 0, 6, 12, 1, 7, 13},
		[]Move{13, 12, 16, 17, 2, 8, 14, 3, 9, 15},
	}[axis][int(p)]
}

// String returns the string representation of the move on the Y axis.
func (p Phase2Move) String() string {
	return p.Move(1).String()
}

// Phase2Moves is a table containing the necessary data to efficiently perform
// moves on a Phase2Cube.
type Phase2Moves struct {
	CornerMoves [40320][10]int
	EdgeMoves   [40320][10]int
	SliceMoves  [24][10]int
}

// NewPhase2Moves generates a Phase2Moves table.
func NewPhase2Moves() *Phase2Moves {
	res := new(Phase2Moves)

	perm8 := allPermutations(8)

	// We set some states to -1, that way we can tell which states have been
	// found and which have not.
	for i := 0; i < 40320; i++ {
		for j := 0; j < 10; j++ {
			res.CornerMoves[i][j] = -1
			res.EdgeMoves[i][j] = -1
		}
	}

	// Generate corner cases.
	for state, perm := range perm8 {
		corners := SolvedCubieCorners()
		// Permute the UD edges for the current case.
		for j, x := range perm {
			corners[j].Piece = x
		}
		// Apply all 10 moves to the cube.
		for m := 0; m < 10; m++ {
			// The result of this move may have already been computed as an
			// inverse.
			if res.CornerMoves[state][m] >= 0 {
				continue
			}

			// Record the end state of this move.
			c := corners
			c.Move(Phase2Move(m).Move(1))
			endState := encodeYCornerPerm(&c)
			res.CornerMoves[state][m] = endState

			// For the end state, the inverse of this move gets the current
			// state.
			res.CornerMoves[endState][int(Phase2Move(m).Inverse())] = state
		}
	}

	// Generate edge cases.
	for state, perm := range perm8 {
		edges := SolvedCubieEdges()
		// Permute the UD edges for the current case.
		for j, x := range perm {
			slot := []int{6, 5, 0, 4, 8, 11, 2, 10}[j]
			piece := []int{6, 5, 0, 4, 8, 11, 2, 10}[x]
			edges[slot].Piece = piece
		}
		// Apply all 10 moves to the cube.
		for m := 0; m < 10; m++ {
			// The result of this move may have already been computed as an
			// inverse.
			if res.EdgeMoves[state][m] >= 0 {
				continue
			}

			// Record the end state of this move.
			e := edges
			e.Move(Phase2Move(m).Move(1))
			endState := encodeUDEdges(&e)
			res.EdgeMoves[state][m] = endState

			// For the end state, the inverse of this move gets the current
			// state.
			res.EdgeMoves[endState][int(Phase2Move(m).Inverse())] = state
		}
	}

	// Generate slice moves
	for state, perm := range allPermutations(4) {
		edges := SolvedCubieEdges()
		// Permute the E slice for the current case.
		for j, x := range perm {
			slot := []int{1, 3, 7, 9}[j]
			piece := []int{1, 3, 7, 9}[x]
			edges[slot].Piece = piece
		}
		// Apply all 10 moves to the cube.
		for m := 0; m < 10; m++ {
			e := edges
			e.Move(Phase2Move(m).Move(1))
			res.SliceMoves[state][m] = encodeESlicePerm(&e)
		}
	}

	return res
}

func encodeESlicePerm(e *CubieEdges) int {
	// Generate a permutation of {0, 1, 2, 3} that represents the permutation of
	// the E slice.
	var perm [4]int
	for i, slot := range []int{1, 3, 7, 9} {
		// Get the physical piece at the given slot.
		piece := (*e)[slot].Piece
		// Get an index from 0 to 3 for the piece.
		perm[i] = []int{-1, 0, -1, 1, -1, -1, -1, 2, -1, 3, -1, -1}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeFBEdges(e *CubieEdges) int {
	var perm [8]int
	for i, slot := range []int{0, 1, 2, 3, 6, 7, 8, 9} {
		piece := (*e)[slot].Piece
		perm[i] = []int{0, 1, 2, 3, -1, -1, 4, 5, 6, 7, -1, -1}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeMSlicePerm(e *CubieEdges) int {
	// This is like encodeESlicePerm, but for the M slice.
	var perm [4]int
	for i, slot := range []int{0, 2, 6, 8} {
		// Get the physical piece at the given slot.
		piece := (*e)[slot].Piece
		// Get an index from 0 to 3 for the piece.
		perm[i] = []int{0, -1, 1, -1, -1, -1, 2, -1, 3, -1, -1, -1}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeRLEdges(e *CubieEdges) int {
	var perm [8]int
	for i, slot := range []int{9, 4, 3, 10, 7, 5, 1, 11} {
		piece := (*e)[slot].Piece
		perm[i] = []int{-1, 6, -1, 2, 1, 5, -1, 4, -1, 0, 3, 7}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeSSlicePerm(e *CubieEdges) int {
	// This is like encodeESlicePerm, but for the S slice.
	var perm [4]int
	for i, slot := range []int{11, 10, 5, 4} {
		// Get the physical piece at the given slot.
		piece := (*e)[slot].Piece
		// Get an index from 0 to 3 for the piece.
		perm[i] = []int{-1, -1, -1, -1, 3, 2, -1, -1, -1, -1, 1, 0}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeUDEdges(e *CubieEdges) int {
	var perm [8]int
	for i, slot := range []int{6, 5, 0, 4, 8, 11, 2, 10} {
		piece := (*e)[slot].Piece
		perm[i] = []int{2, -1, 6, -1, 3, 1, 0, -1, 4, -1, 7, 5}[piece]
		if perm[i] < 0 {
			return -1
		}
	}
	return encodePermutationInPlace(perm[:])
}

func encodeXCornerPerm(c *CubieCorners) int {
	var perm [8]int
	for i, idx := range xCornerIndices {
		perm[i] = inverseXCornerIndices[(*c)[idx].Piece]
	}
	return encodePermutationInPlace(perm[:])
}

func encodeYCornerPerm(c *CubieCorners) int {
	var perm [8]int
	for i := 0; i < 8; i++ {
		perm[i] = (*c)[i].Piece
	}
	return encodePermutationInPlace(perm[:])
}

func encodeZCornerPerm(c *CubieCorners) int {
	var perm [8]int
	for i, idx := range zCornerIndices {
		perm[i] = inverseZCornerIndices[(*c)[idx].Piece]
	}
	return encodePermutationInPlace(perm[:])
}
