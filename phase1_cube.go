package gocube

// xCornerIndices are the indexes of the corners on the Y axis cube which
// correspond to the corners on the X axis cube. An index in this array
// corresponds to the physical slot in the X axis cube. A value in this array
// corresponds to the physical slot in the Y axis cube.
var xCornerIndices []int = []int{1, 3, 0, 2, 5, 7, 4, 6}

// xEdgeIndices are the indexes of the edges on the Y axis cube which correspond
// to edges on the X axis cube. An index in this array corresponds to the
// physical slot in the X axis cube. A value in this array corresponds to the
// physical slot in the Y axis cube.
var xEdgeIndices []int = []int{3, 0, 1, 2, 10, 4, 9, 6, 7, 8, 11, 5}

// xMoveTranslation maps moves from the Y axis phase-1 cube to moves on the X
// axis cube. The mapping is: F->F, B->B, U->R, D->L, L->U, R->D.
// For example, doing U on a Y-axis cube is like doing R on the X-axis version
// of that cube.
// This mapping is kind of like doing a "z" rotation before the move.
var xMoveTranslation []Move = []Move{4, 5, 2, 3, 1, 0, 10, 11, 8, 9, 7, 6, 16,
	17, 14, 15, 13, 12}

// zCornerIndices are like xCornerIndices but for the Z axis cube.
var zCornerIndices []int = []int{2, 3, 6, 7, 0, 1, 4, 5}

// zEdgeIndices are like xEdgeIndices but for the Z axis cube.
var zEdgeIndices []int = []int{2, 11, 8, 10, 3, 1, 0, 5, 6, 4, 9, 7}

// zMoveTranslation is like xMoveTranslation, but it's for doing an "x" rotation
// before applying a move. The mapping is: R->R, L->L, F->U, B->D, U->B, D->F.
var zMoveTranslation []Move = []Move{3, 2, 0, 1, 4, 5, 9, 8, 6, 7, 10, 11, 15,
	14, 12, 13, 16, 17}

// A Phase1Cube is an efficient way to represent the parts of a cube which
// matter for the first phase of Kociemba's algorithm.
// The FB edge orientation can be used for both Y and X phase-1 goals, and the
// UD edge orientation can be used for the Z phase-1 goal. Thus, no RL edge
// orientations are needed.
type Phase1Cube struct {
	XCornerOrientation int
	YCornerOrientation int
	ZCornerOrientation int

	FBEdgeOrientation int
	UDEdgeOrientation int

	MSlicePermutation int
	ESlicePermutation int
	SSlicePermutation int
}

// SolvedPhase1Cube returns a solved phase1 cube.
func SolvedPhase1Cube() Phase1Cube {
	return Phase1Cube{
		1093, 1093, 1093,
		0, 0,
		220, 220, 220,
	}
}

// Move applies a move to a Phase1Cube.
func (p *Phase1Cube) Move(m Move, moves *Phase1Moves) {
	// Apply the move to the y-axis cube.
	p.YCornerOrientation = moves.COMoves[p.YCornerOrientation][m]
	p.FBEdgeOrientation = moves.EOMoves[p.FBEdgeOrientation][m]
	p.ESlicePermutation = moves.ESliceMoves[p.ESlicePermutation][m]

	// Apply the move to the z-axis cube.
	zMove := zMoveTranslation[m]
	p.ZCornerOrientation = moves.COMoves[p.ZCornerOrientation][zMove]
	p.UDEdgeOrientation = moves.EOMoves[p.UDEdgeOrientation][zMove]
	p.SSlicePermutation = moves.ESliceMoves[p.SSlicePermutation][zMove]

	// Apply the move to the x-axis cube.
	xMove := xMoveTranslation[m]
	p.XCornerOrientation = moves.COMoves[p.XCornerOrientation][xMove]
	p.MSlicePermutation = moves.ESliceMoves[p.MSlicePermutation][xMove]
}

// Solved returns whether the phase-1 cube is solved in all three axes.
func (p *Phase1Cube) Solved() (x bool, y bool, z bool) {
	x = true
	y = true
	z = true
	if p.XCornerOrientation != 1093 {
		x = false
	} else if p.MSlicePermutation != 220 {
		x = false
	} else if p.FBEdgeOrientation != 0 {
		x = false
	}
	if p.YCornerOrientation != 1093 {
		y = false
	} else if p.ESlicePermutation != 220 {
		y = false
	} else if p.FBEdgeOrientation != 0 {
		y = false
	}
	if p.ZCornerOrientation != 1093 {
		z = false
	} else if p.SSlicePermutation != 220 {
		z = false
	} else if p.UDEdgeOrientation != 0 {
		z = false
	}
	return
}

// Phase1Moves is a table containing the necessary data to efficiently perform
// moves on a Phase1Cube.
// Note that only one move table is needed for all 3 axes (i.e. all three
// phase-1 goals). Thus, the move tables apply directly to the Y-oriented
// phase-1 goal. Moves much be translated for the X-oriented and Z-oriented
// goals.
type Phase1Moves struct {
	ESliceMoves [495][18]int
	EOMoves     [2048][18]int
	COMoves     [2187][18]int
}

// NewPhase1Moves generates tables for applying phase-1 moves.
func NewPhase1Moves() *Phase1Moves {
	res := &Phase1Moves{}

	// Generate the CO cases and do moves on them.
	for i := 0; i < 2187; i++ {
		corners := decodeCO(i)
		for m := 0; m < 18; m++ {
			aCase := corners
			aCase.Move(Move(m))
			res.COMoves[i][m] = encodeCO(&aCase)
		}
	}

	// Generate the EO cases and do moves on them.
	for i := 0; i < 2048; i++ {
		edges := decodeEO(i)
		for m := 0; m < 18; m++ {
			aCase := edges
			aCase.Move(Move(m))
			res.EOMoves[i][m] = encodeEO(&aCase)
		}
	}

	// Generate the E-slice cases and do moves on them.
	eSliceCase := 0
	for w := 0; w < 12; w++ {
		for x := w + 1; x < 12; x++ {
			for y := x + 1; y < 12; y++ {
				for z := y + 1; z < 12; z++ {
					// The state is bogus, but moves work on it.
					var edges CubieEdges
					edges[w].Piece = 1
					edges[x].Piece = 1
					edges[y].Piece = 1
					edges[z].Piece = 1
					for m := 0; m < 18; m++ {
						aCase := edges
						aCase.Move(Move(m))
						encoded := encodeBogusESlice(&aCase)
						res.ESliceMoves[eSliceCase][m] = encoded
					}
					eSliceCase++
				}
			}
		}
	}

	return res
}

// Phase1Cube generates a Phase1Cube which reflects the state of a CubieCube.
func (c *CubieCube) Phase1Cube() Phase1Cube {
	var res Phase1Cube

	// Encode FB edge orientations
	for i := uint(0); i < 11; i++ {
		if c.Edges[i].Flip {
			res.FBEdgeOrientation |= (1 << i)
		}
	}

	// Encode the UD corner orientations
	scaler := 1
	for i := 0; i < 7; i++ {
		res.YCornerOrientation += scaler * c.Corners[i].Orientation
		scaler *= 3
	}

	// Encode the E slice permutation
	var eChoice [12]bool
	for i := 0; i < 12; i++ {
		piece := c.Edges[i].Piece
		if piece == 1 || piece == 3 || piece == 7 || piece == 9 {
			eChoice[i] = true
		}
	}
	res.ESlicePermutation = encodeChoice(eChoice[:])

	// Translated stuff is too much code to keep in this method.
	res.UDEdgeOrientation = udEdgeOrientations(&c.Edges)
	res.XCornerOrientation, res.ZCornerOrientation =
		xzCornerOrientations(&c.Corners)
	res.MSlicePermutation, res.SSlicePermutation = xzSlices(&c.Edges)

	return res
}

func decodeCO(co int) CubieCorners {
	corners := SolvedCubieCorners()

	// Compute the orientations of the first 7 corners.
	scaler := 1
	for x := 0; x < 7; x++ {
		corners[x].Orientation = (co / scaler) % 3
		scaler *= 3
	}

	// Apply sune combos to orient all the corners except the last one.
	ordering := []int{0, 1, 5, 4, 6, 2, 3, 7}
	orientations := make([]int, 8)
	for i := 0; i < 8; i++ {
		orientations[i] = corners[ordering[i]].Orientation
	}
	for i := 0; i < 7; i++ {
		thisOrientation := orientations[i]
		nextOrientation := orientations[i+1]
		// Twist thisOrientation to be solved, affecting the next corner in the
		// sequence.
		if thisOrientation == 2 {
			// y -> x, x -> z, z -> y
			orientations[i+1] = (nextOrientation + 2) % 3
		} else if thisOrientation == 0 {
			// z -> x, x -> y, y -> z
			orientations[i+1] = (nextOrientation + 1) % 3
		}
	}

	// The twist of the last corner is the inverse of what it should be in the
	// scramble.
	if orientations[7] == 0 {
		corners[7].Orientation = 2
	} else if orientations[7] == 2 {
		corners[7].Orientation = 0
	}

	return corners
}

func decodeEO(eo int) CubieEdges {
	edges := SolvedCubieEdges()
	parity := false
	for x := uint(0); x < 11; x++ {
		if (eo & (1 << x)) != 0 {
			parity = !parity
			edges[x].Flip = true
		}
	}
	edges[11].Flip = parity
	return edges
}

func encodeBogusESlice(c *CubieEdges) int {
	list := make([]bool, 12)
	for i := 0; i < 12; i++ {
		list[i] = (*c)[i].Piece == 1
	}
	return encodeChoice(list)
}

func encodeCO(c *CubieCorners) int {
	res := 0
	scaler := 1
	for i := uint(0); i < 7; i++ {
		res += scaler * (*c)[i].Orientation
		scaler *= 3
	}
	return res
}

func encodeEO(c *CubieEdges) int {
	res := 0
	for i := uint(0); i < 11; i++ {
		if (*c)[i].Flip {
			res |= (1 << i)
		}
	}
	return res
}

func udEdgeOrientations(c *CubieEdges) int {
	res := 0
	for i, idx := range zEdgeIndices[:11] {
		edge := (*c)[idx]
		flip := edge.Flip
		if edge.Piece == 0 || edge.Piece == 2 || edge.Piece == 6 ||
			edge.Piece == 8 {
			// This is an M slice edge piece, so it changes orientation if it
			// was on the S slice or the E slice.
			if idx != 0 && idx != 2 && idx != 6 && idx != 8 {
				flip = !flip
			}
		} else {
			// This is an E or S slice edge, so it changes orientation if it
			// was on the M slice.
			if idx == 0 || idx == 2 || idx == 6 || idx == 8 {
				flip = !flip
			}
		}
		if flip {
			res |= 1 << uint(i)
		}
	}
	return res
}

func xzCornerOrientations(c *CubieCorners) (xVal int, zVal int) {
	var x [8]int
	var z [8]int

	// For each corner, find the direction of the x and z stickers.
	for i := 0; i < 8; i++ {
		corner := (*c)[i]

		// If the corner was in its original slot, here's what the directions
		// would be.
		o := corner.Orientation
		if o == 0 {
			x[i] = 2
			z[i] = 1
		} else if o == 1 {
			x[i] = 0
			z[i] = 2
		} else {
			x[i] = 1
			z[i] = 0
		}

		// If it takes an odd number of quarter turns to move the corner back to
		// its original slot, swap x and z.
		d := (corner.Piece ^ i) & 7
		if d == 1 || d == 2 || d == 4 || d == 7 {
			x[i], z[i] = z[i], x[i]
		}
	}

	// Add the information together to generate the final values.
	scaler := 1
	for i := 0; i < 7; i++ {
		xDirection := x[xCornerIndices[i]]
		if xDirection == 1 {
			xDirection = 0
		} else if xDirection == 0 {
			xDirection = 1
		}
		xVal += scaler * xDirection

		zDirection := z[zCornerIndices[i]]
		if zDirection == 1 {
			zDirection = 2
		} else if zDirection == 2 {
			zDirection = 1
		}
		zVal += scaler * zDirection

		scaler *= 3
	}

	return
}

func xzSlices(e *CubieEdges) (x int, z int) {
	var xChoice [12]bool
	var zChoice [12]bool
	for i, idx := range xEdgeIndices {
		// The M slice is the important slice of the X axis cube.
		p := (*e)[idx].Piece
		if p == 0 || p == 2 || p == 6 || p == 8 {
			xChoice[i] = true
		}
	}
	for i, idx := range zEdgeIndices {
		// The S slice is the important slice of the Z axis cube.
		p := (*e)[idx].Piece
		if p == 4 || p == 5 || p == 10 || p == 11 {
			zChoice[i] = true
		}
	}
	return encodeChoice(xChoice[:]), encodeChoice(zChoice[:])
}
