package gocube

// A CubieCube represents a cube's physical construction.
type CubieCube struct {
	Corners CubieCorners
	Edges   CubieEdges
}

// SolvedCubieCube returns a solved CubieCube
func SolvedCubieCube() CubieCube {
	return CubieCube{SolvedCubieCorners(), SolvedCubieEdges()}
}

// HalfTurn applies a half-turn to the edges and corners.
func (c *CubieCube) HalfTurn(face int) {
	c.Corners.HalfTurn(face)
	c.Edges.HalfTurn(face)
}

// Move applies a move to the edges and corners.
func (c *CubieCube) Move(m Move) {
	c.Corners.Move(m)
	c.Edges.Move(m)
}

// QuarterTurn applies a quarter-turn to the edges and corners.
func (c *CubieCube) QuarterTurn(face, turns int) {
	c.Corners.QuarterTurn(face, turns)
	c.Edges.QuarterTurn(face, turns)
}

// Solved returns true if the edges and corners are solved.
func (c *CubieCube) Solved() bool {
	return c.Corners.Solved() && c.Edges.Solved()
}

