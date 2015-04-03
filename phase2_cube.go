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
		res.CornerPermutation = encodeXCornerPerm(&c.Edges)
	} else if axis == 1 {
		res.CornerPermutation = encodeYCornerPerm(&c.Edges)
	} else {
		res.CornerPermutation = encodeZCornerPerm(&c.Edges)
	}
	// TODO: figure out edge permutation and slice permutation.
	return res, nil
}

func encodeXCornerPerm(e *CubieEdges) int {
	var perm [8]int
	for i, idx := range xCornerIndices {
		perm[i] = inverseXCornerIndices[(*e)[idx].Piece]
	}
	return encodePermutationInPlace(perm[:])
}

func encodeYCornerPerm(e *CubieEdges) int {
	var perm [8]int
	for i := 0; i < 8; i++ {
		perm[i] = (*e)[i].Piece
	}
	return encodePermutationInPlace(perm[:])
}

func encodeZCornerPerm(e *CubieEdges) int {
	var perm [8]int
	for i, idx := range zCornerIndices {
		perm[i] = inverseZCornerIndices[(*e)[idx].Piece]
	}
	return encodePermutationInPlace(perm[:])
}
