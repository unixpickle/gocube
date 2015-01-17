package gocube

import (
	"errors"
)

// A CubieCorner represents a physical corner of a cube.
type CubieCorner struct {
	Piece       int
	Orientation int
}

// A CubieCube represents a cube's physical construction.
type CubieCube struct {
	Corners [8]CubieCorner
	Edges   [12]CubieEdge
}

// StickerToCubie converts a StickerCube to a CubieCube.
func StickerToCubie(c StickerCube) (*CubieCube, error) {
	return nil, errors.New("NYI")
}

// A CubieEdge represents a physical edge of a cube.
type CubieEdge struct {
	Piece int
	Flip  bool
}
