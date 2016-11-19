package gocube

import "math/rand"

// RandomCubieCube generates a random state.
func RandomCubieCube() CubieCube {
	var res CubieCube

	pieces := rand.Perm(8)
	for i, x := range pieces {
		res.Corners[i].Piece = x
	}

	cornerParity := parity(pieces)
	pieces = rand.Perm(12)
	for i, x := range pieces {
		res.Edges[i].Piece = x
	}

	if cornerParity != parity(pieces) {
		res.Edges[11], res.Edges[10] = res.Edges[10], res.Edges[11]
	}

	lastFlip := false
	for i := 0; i < 11; i++ {
		if rand.Intn(2) == 0 {
			lastFlip = !lastFlip
			res.Edges[i].Flip = true
		}
	}
	res.Edges[11].Flip = lastFlip

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

// RandomZBLL generates a cube with a random last layer
// in which the edges are all properly oriented.
func RandomZBLL() CubieCube {
	res := SolvedCubieCube()
	var cornerParity bool
	res.Corners, cornerParity = RandomLLCorners()

	edgePerm := rand.Perm(4)
	if parity(append([]int{}, edgePerm...)) != cornerParity {
		edgePerm[0], edgePerm[1] = edgePerm[1], edgePerm[0]
	}

	topEdges := []int{0, 4, 5, 6}
	for i, j := range edgePerm {
		res.Edges[topEdges[i]].Piece = topEdges[j]
	}

	return res
}

// RandomLLCorners generates random last-layer corners
// and returns the corners as well as their parity.
//
// A parity of true is even, while false is odd.
func RandomLLCorners() (CubieCorners, bool) {
	orientations := make([]int, 4)
	for i := 0; i < 3; i++ {
		orientations[i] = rand.Intn(3)
	}

	o := append([]int{}, orientations...)
	for i := 0; i < 3; i++ {
		this := o[i]
		next := o[i+1]
		if this == 2 {
			o[i+1] = (next + 2) % 3
		} else {
			o[i+1] = (next + 1) % 3
		}
	}
	if o[3] == 0 {
		orientations[3] = 2
	} else if o[3] == 2 {
		orientations[3] = 0
	} else {
		orientations[3] = 1
	}

	cube := SolvedCubieCorners()
	perm := rand.Perm(4)
	pieces := []int{2, 3, 7, 6}
	for i, piece := range pieces {
		cube[piece].Orientation = orientations[i]
		cube[piece].Piece = pieces[perm[i]]
	}
	return cube, parity(perm)
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
