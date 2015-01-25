package gocube

import (
	"strconv"
	"sync"
)

// A COGoal is satisfied when the corners are all oriented.
type COGoal struct{}

// IsGoal returns true if the corners are all oriented.
func (_ COGoal) IsGoal(c CubieCorners) bool {
	for _, x := range c {
		if x.Orientation != 0 {
			return false
		}
	}
	return true
}

// A COPruner stores the number of moves required to solve every given
// corner-orientation case.
type COPruner [2187]int

// MakeCOPruner uses breadth-first search to find all the CO cases.
func MakeCOPruner(moves []Move) *COPruner {
	var res COPruner

	for i := 0; i < 2187; i++ {
		res[i] = -1
	}

	// Perform a very basic breadth-first search.
	nodes := []cornersSearchNode{cornersSearchNode{SolvedCubieCorners(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := node.State.EncodeCO()

		if res[idx] >= 0 {
			// Node was already visited
			continue
		}

		// Expand the node
		res[idx] = node.Depth
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, cornersSearchNode{state, node.Depth + 1})
		}
	}

	return &res
}

// MinMoves encodes the corners' orientations and looks it up in the move
// count table.
func (c *COPruner) MinMoves(corners CubieCorners) int {
	return c[corners.EncodeCO()]
}

// A CornerPermPruner records the number of moves to solve every corner
// permutation.
type CornerPermPruner map[int]int

// MakeCornerPermPruner creates a new CornerPermPruner using a breadth-first
// search.
func MakeCornerPermPruner(moves []Move) CornerPermPruner {
	res := CornerPermPruner{}
	
	// Perform a very basic breadth-first search.
	nodes := []cornersSearchNode{cornersSearchNode{SolvedCubieCorners(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := node.State.EncodePerm()

		if _, ok := res[idx]; ok {
			// Node was already visited
			continue
		}

		// Expand the node
		res[idx] = node.Depth
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, cornersSearchNode{state, node.Depth + 1})
		}
	}
	
	return res
}

// MinMoves returns the number of moves required to solve the corners.
func (c CornerPermPruner) MinMoves(corners CubieCorners) int {
	if res, ok := c[corners.EncodePerm()]; !ok {
		panic("Missing corner case.")
	} else {
		return res
	}
}

// A CornersGoal represents an abstract goal state for a depth-first search of
// the cube's corners.
type CornersGoal interface {
	IsGoal(c CubieCorners) bool
}

// A CornersPruner is used as a lower-bound heuristic for a depth-first search
// of the cube's corners.
type CornersPruner interface {
	MinMoves(c CubieCorners) int
}

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

// SolvedCubieCorners generates the corners of a solved cube.
func SolvedCubieCorners() CubieCorners {
	var res CubieCorners
	for i := 0; i < 8; i++ {
		res[i].Piece = i
	}
	return res
}

// EncodeCO encodes the orientations of the first seven corners as an integer
// which ranges between 0 (inclusive) and 3^7 (=2,187; exclusive).
func (c *CubieCorners) EncodeCO() int {
	result := 0
	scale := 1
	for i := 0; i < 7; i++ {
		result += scale * c[i].Orientation
		scale *= 3
	}
	return result
}

// EncodePerm sub-optimally encodes the corner permutations. The result will
// range from 0 (inclusive) to 8^7 (exclusive).
func (c *CubieCorners) EncodePerm() int {
	res := 0
	scaler := 1
	for i := 0; i < 7; i++ {
		res += c[i].Piece * scaler
		scaler *= 8
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

// Search starts a search using the receiver as the starting state.
// If the specified CornersPruner is nil, no pruning will be performed.
// The depth argument specifies the maximum depth for the search.
// The branch argument specifies how many levels of the search to parallelize.
func (c *CubieCorners) Search(g CornersGoal, p CornersPruner, moves []Move,
	depth, branch int) Search {
	res := &cornersSearch{newSimpleSearch(moves), g, p}
	go func(st CubieCorners) {
		prefix := make([]Move, 0, depth)
		searchCornersBranch(st, res, depth, branch, prefix)
		close(res.channel)
	}(*c)
	return res
}

// Solved returns true if all the corners are properly positioned and oriented.
func (c *CubieCorners) Solved() bool {
	for i := 0; i < 8; i++ {
		if c[i].Piece != i || c[i].Orientation != 0 {
			return false
		}
	}
	return true
}

// A SolveCornersGoal is satisfied when a CubieCorners is completely solved.
type SolveCornersGoal struct{}

// IsGoal returns corners.Solved().
func (_ SolveCornersGoal) IsGoal(corners CubieCorners) bool {
	return corners.Solved()
}

func searchCorners(st CubieCorners, s *cornersSearch, depth int,
	prefix []Move) {
	// If we can't search any further, check if it's the goal.
	if depth == 0 {
		if s.goal.IsGoal(st) {
			// We must make a copy of the prefix before sending it as a
			// solution, since it may be modified after we return.
			solution := make([]Move, len(prefix))
			copy(solution, prefix)
			s.channel <- solution
		}
		return
	}

	// Prune the state
	if s.prune(st) > depth {
		return
	}

	// Apply each move and recurse.
	for _, move := range s.moves {
		if depth > 5 && s.cancelled() {
			return
		}
		newState := st
		newState.Move(move)
		searchCorners(newState, s, depth-1, append(prefix, move))
	}
}

func searchCornersBranch(st CubieCorners, s *cornersSearch, depth, branch int,
	prefix []Move) {
	// If we shouldn't branch, do a regular search.
	if branch == 0 || depth == 0 {
		searchCorners(st, s, depth, prefix)
		return
	}

	// Prune the state
	if s.prune(st) > depth {
		return
	}

	// Run each search on a different goroutine
	wg := sync.WaitGroup{}
	for _, move := range s.moves {
		wg.Add(1)
		go func(m Move, newState CubieCorners) {
			// Apply the move
			newState.Move(m)

			// Create the new prefix by copying the old one.
			pref := make([]Move, len(prefix)+1, len(prefix)+depth)
			copy(pref, prefix)
			pref[len(prefix)] = m

			// Branch out and search.
			searchCornersBranch(newState, s, depth-1, branch-1, pref)

			wg.Done()
		}(move, st)
	}
	wg.Wait()
}

type cornersSearch struct {
	*simpleSearch
	goal   CornersGoal
	pruner CornersPruner
}

func (c *cornersSearch) prune(state CubieCorners) int {
	if c.pruner == nil {
		return 0
	}
	return c.pruner.MinMoves(state)
}

type cornersSearchNode struct {
	State CubieCorners
	Depth int
}
