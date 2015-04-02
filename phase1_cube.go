package gocube

// xMoveTranslation maps moves from the Y axis phase-1 cube to moves on the X
// axis cube. The mapping is: F->F, B->B, U->R, D->L, L->U, R->D.
// For example, doing U on a Y-axis cube is like doing R on the X-axis version
// of that cube.
// This mapping is kind of like doing a "z" rotation before the move.
var xMoveTranslation []Move = []Move{4, 5, 2, 3, 1, 0, 10, 11, 8, 9, 7, 6, 16,
	17, 14, 15, 13, 12}

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
