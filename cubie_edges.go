package gocube

import (
	"strconv"
	"sync"
)

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

// EncodeEO returns an 11-bit number representing the edge-orientation case of a
// CubieEdges.
func (c *CubieEdges) EncodeEO() int {
	res := 0
	for i := 0; i < 11; i++ {
		if c[i].Flip {
			res |= (1 << uint(i))
		}
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

// Search starts a search using the receiver as the starting state.
// If the specified EdgesPruner is nil, no pruning will be performed.
// The depth argument specifies the maximum depth for the search.
// The branch argument specifies how many levels of the search to parallelize.
func (c *CubieEdges) Search(g EdgesGoal, p EdgesPruner, moves []Move,
	depth, branch int) Search {
	res := &edgesSearch{newSimpleSearch(moves), g, p}
	go func(st CubieEdges) {
		prefix := make([]Move, 0, depth)
		searchEdgesBranch(st, res, depth, branch, prefix)
		close(res.channel)
	}(*c)
	return res
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

// An EdgesGoal represents an abstract goal state for a depth-first search of
// the cube's edges.
type EdgesGoal interface {
	IsGoal(c CubieEdges) bool
}

// An EdgesPruner is used as a lower-bound heuristic for a depth-first search of
// the cube's edges.
type EdgesPruner interface {
	MinMoves(c CubieEdges) int
}

// An EOGoal is satisfied when a CubieEdges has no flipped edges.
type EOGoal struct{}

// IsGoal returns true if no edges are flipped.
func (_ EOGoal) IsGoal(edges CubieEdges) bool {
	for _, x := range edges {
		if x.Flip {
			return false
		}
	}
	return true
}

// An EOLineGoal is satisfied when a CubieEdges has the edge-orientation line as
// defined by the ZZ method.
type EOLineGoal struct{}

// IsGoal returns true if the edges have the EOLine solved.
func (_ EOLineGoal) IsGoal(c CubieEdges) bool {
	// Make sure the edges are oriented.
	for _, p := range c {
		if p.Flip {
			return false
		}
	}

	// Make sure the line edges are solved.
	return c[2].Piece == 2 && c[8].Piece == 8
}

// An EOPruner stores the number of moves required to solve every EO
// configuration.
type EOPruner [0x800]int

// MakeEOPruner uses breadth-first search to generate an EOPruner.
func MakeEOPruner(moves []Move) *EOPruner {
	var res EOPruner

	for i := 0; i < 0x800; i++ {
		res[i] = -1
	}

	// Perform a very basic breadth-first search.
	nodes := []edgesSearchNode{edgesSearchNode{SolvedCubieEdges(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := node.State.EncodeEO()

		if res[idx] >= 0 {
			// Node was already visited
			continue
		}

		// Expand the node
		res[idx] = node.Depth
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, edgesSearchNode{state, node.Depth + 1})
		}
	}

	return &res
}

// MinMoves returns the number of moves required to orient all the edges.
func (o *EOPruner) MinMoves(e CubieEdges) int {
	return o[e.EncodeEO()]
}

// A SolveEdgesGoal is satisfied when a CubieEdges is completely solved.
type SolveEdgesGoal struct{}

// IsGoal returns edges.Solved().
func (_ SolveEdgesGoal) IsGoal(edges CubieEdges) bool {
	return edges.Solved()
}

func searchEdges(st CubieEdges, s *edgesSearch, depth int, prefix []Move) {
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
		searchEdges(newState, s, depth-1, append(prefix, move))
	}
}

func searchEdgesBranch(st CubieEdges, s *edgesSearch, depth, branch int,
	prefix []Move) {
	// If we shouldn't branch, do a regular search.
	if branch == 0 || depth == 0 {
		searchEdges(st, s, depth, prefix)
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
		go func(m Move, newState CubieEdges) {
			// Apply the move
			newState.Move(m)

			// Create the new prefix by copying the old one.
			pref := make([]Move, len(prefix)+1, len(prefix)+depth)
			copy(pref, prefix)
			pref[len(prefix)] = m

			// Branch out and search.
			searchEdgesBranch(newState, s, depth-1, branch-1, pref)

			wg.Done()
		}(move, st)
	}
	wg.Wait()
}

type edgesSearch struct {
	*simpleSearch
	goal   EdgesGoal
	pruner EdgesPruner
}

func (c *edgesSearch) prune(state CubieEdges) int {
	if c.pruner == nil {
		return 0
	}
	return c.pruner.MinMoves(state)
}

type edgesSearchNode struct {
	State CubieEdges
	Depth int
}
