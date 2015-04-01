package gocube

import "strconv"

// A CubieEdge represents a physical edge of a cube.
// Edges are indexed from 0 through 11 in the following order:
// UF, RF, DF, LF, UL, UR, BU, BR, BD, BL, DL, DR.
// The flip field is true if the edge is "bad" in the ZZ color scheme (i.e. if
// it requires an F or B move to fix).
type CubieEdge struct {
	Piece int
	Flip  bool
}

// CubieEdges represents the edges of a cube.
type CubieEdges [12]CubieEdge

// SolvedCubieEdges returns CubieEdges in their solved state.
func SolvedCubieEdges() CubieEdges {
	var res CubieEdges
	for i := 0; i < 12; i++ {
		res[i].Piece = i
	}
	return res
}

// HalfTurn performs a 180 degree turn on a given face.
func (c *CubieEdges) HalfTurn(face int) {
	// Every half-turn is really just two swaps.
	switch face {
	case 1: // Top face
		c[0], c[6] = c[6], c[0]
		c[4], c[5] = c[5], c[4]
	case 2: // Bottom face
		c[2], c[8] = c[8], c[2]
		c[10], c[11] = c[11], c[10]
	case 3: // Front face
		c[0], c[2] = c[2], c[0]
		c[1], c[3] = c[3], c[1]
	case 4: // Back face
		c[6], c[8] = c[8], c[6]
		c[7], c[9] = c[9], c[7]
	case 5: // Right face
		c[1], c[7] = c[7], c[1]
		c[5], c[11] = c[11], c[5]
	case 6: // Left face
		c[3], c[9] = c[9], c[3]
		c[4], c[10] = c[10], c[4]
	default:
		panic("Unsupported half-turn applied to CubieEdges: " +
			strconv.Itoa(face))
	}
}

// Move applies a face turn to the edges.
func (c *CubieEdges) Move(m Move) {
	// Half turns are a simple case.
	if m.Turns == 2 {
		c.HalfTurn(m.Face)
	} else {
		c.QuarterTurn(m.Face, m.Turns)
	}
}

// QuarterTurn performs a 90 degree turn on a given face.
func (c *CubieEdges) QuarterTurn(face, turns int) {
	switch face {
	case 1: // Top face
		if turns == 1 {
			c[0], c[4], c[6], c[5] = c[5], c[0], c[4], c[6]
		} else {
			c[5], c[0], c[4], c[6] = c[0], c[4], c[6], c[5]
		}
	case 2: // Bottom face
		if turns == 1 {
			c[2], c[11], c[8], c[10] = c[10], c[2], c[11], c[8]
		} else {
			c[10], c[2], c[11], c[8] = c[2], c[11], c[8], c[10]
		}
	case 3: // Front face
		if turns == 1 {
			c[0], c[1], c[2], c[3] = c[3], c[0], c[1], c[2]
		} else {
			c[3], c[0], c[1], c[2] = c[0], c[1], c[2], c[3]
		}
		// Flip edges
		for i := 0; i < 4; i++ {
			c[i].Flip = !c[i].Flip
		}
	case 4: // Back face
		if turns == 1 {
			c[6], c[9], c[8], c[7] = c[7], c[6], c[9], c[8]
		} else {
			c[7], c[6], c[9], c[8] = c[6], c[9], c[8], c[7]
		}
		// Flip edges
		for i := 6; i < 10; i++ {
			c[i].Flip = !c[i].Flip
		}
	case 5: // Right face
		if turns == 1 {
			c[1], c[5], c[7], c[11] = c[11], c[1], c[5], c[7]
		} else {
			c[11], c[1], c[5], c[7] = c[1], c[5], c[7], c[11]
		}
	case 6: // Left face
		if turns == 1 {
			c[3], c[10], c[9], c[4] = c[4], c[3], c[10], c[9]
		} else {
			c[4], c[3], c[10], c[9] = c[3], c[10], c[9], c[4]
		}
	default:
		panic("Unsupported quarter-turn applied to CubieEdges: " +
			strconv.Itoa(face))
	}
}

// Solved returns true if all the edges are properly positioned and oriented.
func (c *CubieEdges) Solved() bool {
	for i := 0; i < 12; i++ {
		if c[i].Piece != i || c[i].Flip {
			return false
		}
	}
	return true
}
