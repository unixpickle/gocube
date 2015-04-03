package gocube

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
		res.SlicePermutation = encodeMSlicePerm(&c.Edges)
	} else if axis == 1 {
		res.CornerPermutation = encodeYCornerPerm(&c.Corners)
		res.SlicePermutation = encodeESlicePerm(&c.Edges)
	} else {
		res.CornerPermutation = encodeZCornerPerm(&c.Corners)
		res.SlicePermutation = encodeSSlicePerm(&c.Edges)
	}
	return res, nil
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
	// TODO: implement this.
	return -1
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
	// TODO: implement this.
	return -1
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
	// TODO: implement this.
	return -1
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
