package gocube

// A Phase1Axis represents the y-axis corner orientations, ZZ edge orientations,
// and the permutation of the E slice.
type Phase1Axis struct {
	CornerOrientations int
	EdgeOrientations	 int
	SlicePerm					int
}

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

	ESlicePermutation int
	SSlicePermutation int
	MSlicePermutation int
}

// Move applies a move on a Phase1Cube using a moves table.
func (p *Phase1Cube) Move(m Move, table *Phase1Moves) {
	// TODO: apply the move to each axis of the represented data.
}

// Phase1Moves is a table containing the necessary data to efficiently perform
// moves on a Phase1Cube.
// Note that only one move table is needed for all 3 axes (i.e. all three
// phase-1 goals). Thus, the move tables apply directly to the Y-oriented
// phase-1 goal. Moves much be translated for the X-oriented and Z-oriented
// goals.
type Phase1Moves struct {
	ESliceMoves [495][18]int
	EOMoves [2048][18]int
	COMoves [2187][18]int
}

// NewPhase1Moves generates tables for applying phase-1 moves.
func NewPhase1Moves() *Phase1Moves {
	res := &Phase1Moves{}
	
	// Generate the CO cases and do moves on them.
	for i := 0; i < 2187; i++ {
		corners := decodeCO(i);
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
		for x := w+1; x < 12; x++ {
			for y := x+1; y < 12; y++ {
				for z := y+1; z < 12; z++ {
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

func decodeCO(co int) CubieCorners {
	corners := SolvedCubieCorners()
	
	// Compute the orientations of the first 7 corners.
	scaler := 1
	for x := 0; x < 7; x++ {
		corners[x].Orientation = (co/scaler) % 3
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
		nextOrientation := orientations[i + 1]
		// Twist thisOrientation to be solved, affecting the next corner in the
		// sequence.
		if thisOrientation == 1 {
			// y -> x, x -> z, z -> y
			orientations[i + 1] = (nextOrientation+2) % 3
		} else if thisOrientation == 2 {
			// z -> x, x -> y, y -> z
			orientations[i + 1] = (nextOrientation+1) % 3
		}
	}
	
	// The twist of the last corner is the inverse of what it should be in the
	// scramble.
	if orientations[7] == 1 {
		corners[7].Orientation = 2
	} else if orientations[7] == 2 {
		corners[7].Orientation = 1
	}

	return corners
}

func decodeEO(eo int) CubieEdges {
	edges := SolvedCubieEdges()
	parity := false
	for x := uint(0); x < 11; x++ {
		if (x & (1 << x)) != 0 {
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
