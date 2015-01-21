package gocube

import (
	"sync"
)

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

// An EOLineGoal is satisfied when a CubieEdges has the edge-orientation line as
// defined by the ZZ method.
type EOLineGoal struct{}

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

// An OrientEdgesGoal is satisfied when a CubieEdges has no flipped edges.
type OrientEdgesGoal struct{}

// IsGoal returns true if no edges are flipped.
func (_ OrientEdgesGoal) IsGoal(edges CubieEdges) bool {
	for _, x := range edges {
		if x.Flip {
			return false
		}
	}
	return true
}

// An OrientEdgesPruner stores the number of moves required to solve every EO
// configuration.
type OrientEdgesPruner [0x800]int

// MakeOrientEdgesPruner uses breadth-first search to generate an
// OrientEdgesHeuristic.
func MakeOrientEdgesPruner(moves []Move) *OrientEdgesPruner {
	var res OrientEdgesPruner

	for i := 0; i < 0x800; i++ {
		res[i] = -1
	}

	// Perform a very basic breadth-first search.
	nodes := []edgesSearchNode{edgesSearchNode{SolvedCubieEdges(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := encodeEO(node.State)

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
func (o *OrientEdgesPruner) MinMoves(e CubieEdges) int {
	return o[encodeEO(e)]
}

// A SolveEdgesGoal is satisfied when a CubieEdges is completely solved.
type SolveEdgesGoal struct{}

// IsGoal returns edges.Solved().
func (_ SolveEdgesGoal) IsGoal(edges CubieEdges) bool {
	return edges.Solved()
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

func encodeEO(e CubieEdges) int {
	res := 0
	for i := 0; i < 11; i++ {
		if e[i].Flip {
			res |= (1 << uint(i))
		}
	}
	return res
}
