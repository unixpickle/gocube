package gocube

import "strconv"

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
// white or yellow sticker (i.e. the sticker that is usually normal to the y
// axis). The corners of a solved cube all have an orientation of 1.
type CubieCorner struct {
	Piece       int
	Orientation int
}

// CubieCorners represents the corners of a cube.
type CubieCorners [8]CubieCorner

// SolvedCubieCorners generates the corners of a solved cube.
func SolvedCubieCorners() CubieCorners {
	var res CubieCorners
	for i := 0; i < 8; i++ {
		res[i].Piece = i
		res[i].Orientation = 1
	}
	return res
}

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
		c[1], c[7] = c[7], c[1]
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
	if m >= 12 {
		c.HalfTurn(m.Face())
	} else {
		c.QuarterTurn(m.Face(), m.Turns())
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

// Solved returns true if all the corners are properly positioned and oriented.
func (c *CubieCorners) Solved() bool {
	for i := 0; i < 8; i++ {
		if c[i].Piece != i || c[i].Orientation != 1 {
			return false
		}
	}
	return true
}

// EncodeIndex encodes the state of the corners as a unique integer in the
// range [0, 3^7 * 8!).
//
// This assumes that the corner orientation is valid.
func (c *CubieCorners) EncodeIndex() uint32 {
	var result uint32

	base := uint32(1)
	for _, corner := range c[:7] {
		result += base * uint32(corner.Orientation)
		base *= 3
	}
	perm := make([]int, 8)
	for i, corner := range c {
		perm[i] = corner.Piece
	}
	result += uint32(encodePermutationInPlace(perm)) * base

	return result
}

// fixLastOrientation orients the final corner (corner 7),
// assuming all of the other corners are correct.
func (c *CubieCorners) fixLastOrientation() {
	// Start by orienting the final corner upright. Then, solve a series of
	// adjacent corners in order, twisting the next corner in the opposite
	// direction to preserve orientation. The resulting orientation of the
	// final corner tells us the inverse of what the orientation should have
	// been.

	c[7].Orientation = 1

	var orientations [8]int
	for i, x := range []int{0, 1, 5, 4, 6, 2, 3, 7} {
		orientations[i] = c[x].Orientation
	}
	for i := 0; i < 7; i++ {
		thisOrientation := orientations[i]
		nextOrientation := orientations[i+1]
		if thisOrientation == 2 {
			orientations[i+1] = (nextOrientation + 2) % 3
		} else if thisOrientation == 0 {
			orientations[i+1] = (nextOrientation + 1) % 3
		}
	}
	if orientations[7] == 0 {
		c[7].Orientation = 2
	} else if orientations[7] == 2 {
		c[7].Orientation = 0
	}
}
