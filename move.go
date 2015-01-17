package gocube

// A Move stores a single move on the cube.
//
// The face of a move can be 1 through 6 for U, D, F, B, R, and L respectively.
// In some cases, the face can be 7 (M), 8 (E), or 9 (S) as well.
//
// The turns can be 1, -1, or 2 to indicate the number of times to turn the
// face.
type Move struct {
	Face  int
	Turns int
}
