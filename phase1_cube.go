package gocube

// A Phase1Axis represents the y-axis corner orientations, ZZ edge orientations,
// and the permutation of the E slice.
type Phase1Axis struct {
	CornerOrientations int
	EdgeOrientations   int
	SlicePerm          int
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
	COMoves [2048][18]int
}

// NewPhase1Moves generates tables for applying phase-1 moves.
func NewPhase1Moves() *Phase1Moves {
	res := &Phase1Moves{}

	// TODO: generate the E-slice moves
	// TODO: generate the EO moves
	// TODO: generate the CO moves

	return res
}
