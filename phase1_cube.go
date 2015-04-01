package gocube


// A Phase1Cube is an efficient way to represent the parts of a cube which
// matter for the first phase of Kociemba's algorithm.
type Phase1Cube struct {
	MSlicePerm int
	ESlicePerm int
	SSlicePerm int
	
	RLEdgeOrientation int
	UDEdgeOrientation int
	FBEdgeOrientation int
	
	RLCornerOrientation int
	UDCornerOrientation int
	FBCornerOrientation int
	
	CornerPerm int
}

// Move applies a move on a Phase1Cube using a moves table.
func (p *Phase1Cube) Move(m Move, table *Phase1Moves) {
	// TODO: this
}

// Phase1Moves is a table containing the necessary data to efficiently perform
// moves on a Phase1Cube.
type Phase1Moves struct {
	ESliceMoves [495][18]int
	UDEOMoves [2048][18]int
	UDCOMoves [2048][18]int
}

func NewPhase1Moves() *Phase1Moves {
	res := &Phase1Moves{}

	// 

	return res
}
