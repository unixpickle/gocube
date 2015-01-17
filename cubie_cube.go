package gocube

import (
	"strconv"
)

// A CubieCorner represents a physical corner of a cube.
//
// To understand the meaning of a CubieCorner's fields, you must first
// understand the coordinate system. There are there axes, x, y, and z.
// The x axis is 0 at the L face and 1 at the R face.
// The y axis is 0 at the D face and 1 at the U face.
// The z axis is 0 at the B face and 1 at the F face.
//
// A corner piece's index is determined by it's original position on the cube.
// The index is a binary number of the form ZYX, where Z is the most significant
// digit. Thus, the BLD corner is 0, the BRU corner is 3, the FRU corner is 7,
// etc.
//
// The orientation of a corner tells how it is twisted. It is an axis number 0,
// 1, or 2 for x, y, or z respectively. It indicates the direction normal to the
// red or orange sticker (i.e. the sticker that is usually normal to the x
// axis).
type CubieCorner struct {
	Piece       int
	Orientation int
}

// CubieCorners represents the corners of a cube.
type CubieCorners [8]CubieCorner

// HalfTurn performs a 180 degree turn on a given face.
func (c *CubieCorners) HalfTurn(face int) {
	// A double turn is really just two swaps.
	switch face {
	case 1: // Top face
		c[2], c[7] = c[7], c[2]
		c[3], c[6] = c[6], c[3]
	case 2: // Bottom face
		c[0], c[5] = c[5], c[0]
		c[1], c[4] = c[4], c[1]
	case 3: // Front face
		c[5], c[6] = c[6], c[5]
		c[4], c[7] = c[7], c[4]
	case 4: // Back face
		c[0], c[3] = c[3], c[0]
		c[1], c[2] = c[2], c[1]
	case 5: // Right face
		c[4], c[7] = c[7], c[4]
		c[3], c[5] = c[5], c[3]
	case 6: // Left face
		c[0], c[6] = c[6], c[0]
		c[2], c[4] = c[4], c[2]
	default:
		panic("Unsupported half-turn applied to CubieCorners: " +
			strconv.Itoa(face))
	}
}

// Move applies a face turn to the corners.
func (c *CubieCorners) Move(m Move) {
	// Half turns are a simple case.
	if m.Turns == 2 {
		c.HalfTurn(m.Face)
	} else {
		c.QuarterTurn(m.Face, m.Turns)
	}
}

// QuarterTurn performs a 90 degree turn on a given face.
func (c *CubieCorners) QuarterTurn(face, turns int) {
	// This code is not particularly graceful, but it is rather efficient and
	// quite readable compared to a pure array of transformations.
	switch face {
	case 1: // Top face
		if turns == 1 {
			c[2], c[3], c[7], c[6] = c[6], c[2], c[3], c[7]
		} else {
			c[6], c[2], c[3], c[7] = c[2], c[3], c[7], c[6]
		}
		// Swap orientation 0 with orientation 2.
		for _, i := range []int{2, 3, 6, 7} {
			c[i].Orientation = 2 - c[i].Orientation
		}
	case 2: // Bottom face
		if turns == 1 {
			c[4], c[0], c[1], c[5] = c[0], c[1], c[5], c[4]
		} else {
			c[0], c[1], c[5], c[4] = c[4], c[0], c[1], c[5]
		}
		// Swap orientation 0 with orientation 2.
		for _, i := range []int{0, 1, 4, 5} {
			c[i].Orientation = 2 - c[i].Orientation
		}
	case 3: // Front face
		if turns == 1 {
			c[6], c[7], c[5], c[4] = c[4], c[6], c[7], c[5]
		} else {
			c[4], c[6], c[7], c[5] = c[6], c[7], c[5], c[4]
		}
		// Swap orientation 0 with orientation 1.
		for _, i := range []int{4, 5, 6, 7} {
			if c[i].Orientation == 0 {
				c[i].Orientation = 1
			} else if c[i].Orientation == 1 {
				c[i].Orientation = 0
			}
		}
	case 4: // Back face
		if turns == 1 {
			c[0], c[2], c[3], c[1] = c[2], c[3], c[1], c[0]
		} else {
			c[2], c[3], c[1], c[0] = c[0], c[2], c[3], c[1]
		}
		// Swap orientation 0 with orientation 1.
		for _, i := range []int{0, 1, 2, 3} {
			if c[i].Orientation == 0 {
				c[i].Orientation = 1
			} else if c[i].Orientation == 1 {
				c[i].Orientation = 0
			}
		}
	case 5: // Right face
		if turns == 1 {
			c[7], c[3], c[1], c[5] = c[5], c[7], c[3], c[1]
		} else {
			c[5], c[7], c[3], c[1] = c[7], c[3], c[1], c[5]
		}
		// Swap orientation 2 with orientation 1.
		for _, i := range []int{1, 3, 5, 7} {
			if c[i].Orientation == 1 {
				c[i].Orientation = 2
			} else if c[i].Orientation == 2 {
				c[i].Orientation = 1
			}
		}
	case 6: // Left face
		if turns == 1 {
			c[4], c[6], c[2], c[0] = c[6], c[2], c[0], c[4]
		} else {
			c[6], c[2], c[0], c[4] = c[4], c[6], c[2], c[0]
		}
		// Swap orientation 2 with orientation 1.
		for _, i := range []int{0, 2, 4, 6} {
			if c[i].Orientation == 1 {
				c[i].Orientation = 2
			} else if c[i].Orientation == 2 {
				c[i].Orientation = 1
			}
		}
	default:
		panic("Unsupported quarter-turn applied to CubieCorners: " +
			strconv.Itoa(face))
	}
}

// A CubieCube represents a cube's physical construction.
type CubieCube struct {
	Corners CubieCorners
	Edges   CubieEdges
}

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
			c[0], c[1], c[2], c[3] = c[3], c[0], c[1], c[2]
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
