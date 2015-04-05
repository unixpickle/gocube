package gocube

import (
	"math/rand"
)

func RandomCubieCube() CubieCube {
	var res CubieCube

	// Generate a random permutation for the corners.
	pieces := rand.Perm(8)
	for i, x := range pieces {
		res.Corners[i].Piece = x
	}

	// Generate a random permutation for the edges.
	cornerParity := parity(pieces)
	pieces = rand.Perm(12)
	for i, x := range pieces {
		res.Edges[i].Piece = x
	}

	// Make sure the overall parity is even.
	if cornerParity != parity(pieces) {
		res.Edges[11], res.Edges[10] = res.Edges[10], res.Edges[11]
	}

	// Generate edge orientations.
	lastFlip := false
	for i := 0; i < 11; i++ {
		if rand.Intn(2) == 0 {
			lastFlip = !lastFlip
			res.Edges[i].Flip = true
		}
	}
	res.Edges[11].Flip = lastFlip

	// Generate the corner orientations.
	for i := 0; i < 7; i++ {
		res.Corners[i].Orientation = rand.Intn(3)
	}
	var orientations [8]int
	for i, x := range []int{0, 1, 5, 4, 6, 2, 3, 7} {
		orientations[i] = res.Corners[x].Orientation
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
		res.Corners[7].Orientation = 2
	} else if orientations[7] == 2 {
		res.Corners[7].Orientation = 0
	}

	return res
}

// parity returns true if the parity is even.
func parity(perm []int) bool {
	parity := true
	for i := 0; i < len(perm); i++ {
		if perm[i] == i {
			continue
		}
		// Swap this element with i.
		for j := i + 1; j < len(perm); j++ {
			if perm[j] == i {
				perm[i], perm[j] = perm[j], perm[i]
				parity = !parity
				break
			}
		}
	}
	return parity
}
