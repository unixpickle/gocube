package fmc

import "github.com/unixpickle/gocube"

// CornerSticker is an "address" of a corner sticker on the cube.
type CornerSticker struct {
	// Corner is the physical slot of the corner to which this sticker belongs.
	Corner int

	// Axis is 0 for x, 1 for y, and 2 for z, and indicates the axis normal to
	// the sticker in question. This field distinguishes one sticker from
	// another on a given corner piece.
	Axis int
}

// ThreeCycle represents a corner three-cycle. The ThreeCycle takes the first
// sticker and moves it to the second slot, and so on.
type ThreeCycle struct {
	Stickers [3]CornerSticker
}

// Apply performs this three-cycle on a set of corners and returns the resulting
// corners.
func (t *ThreeCycle) Apply(c gocube.CubieCorners) gocube.CubieCorners {
	// TODO: this.
	return gocube.SolvedCubieCorners()
}

// Insertion generates an 8-move insertion for the ThreeCycle. If the cycle
// cannot be expressed in an insertion, this will return nil.
func (t *ThreeCycle) Insertion() []gocube.Move {
	// TODO: this.
	return nil
}

// Move applies a move to this three-cycle.
//
// For example, suppose a ThreeCycle cycles the UFR, UFL, and UBL corners.
// Applying a B' move would result in a ThreeCycle which cycles UFR, UFL and
// RUB.
func (t *ThreeCycle) Move(m gocube.Move) ThreeCycle {
	// TODO: this.
	return *t
}
