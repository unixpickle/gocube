package gocube

import (
	"sync"
)

// A CubeGoal represents an abstract goal state for a depth-first search of the
// cube.
type CubeGoal interface {
	IsGoal(c CubieCube) bool
}

// A CubePruner is used as a lower-bound heuristic for a depth-first search of
// the cube.
type CubePruner interface {
	MinMoves(c CubieCube) int
}

// A SolveCubeGoal is satisfied when a CubieCube is completely solved.
type SolveCubeGoal struct{}

// IsGoal returns cube.Solved().
func (_ SolveCubeGoal) IsGoal(cube CubieCube) bool {
	return cube.Solved()
}

// Search starts a search using the receiver as the starting state.
// If the specified CubePruner is nil, no pruning will be performed.
// The depth argument specifies the maximum depth for the search.
// The branch argument specifies how many levels of the search to parallelize.
func (c *CubieCube) Search(g CubeGoal, p CubePruner, moves []Move,
	depth, branch int) Search {
	res := &cubeSearch{newSimpleSearch(moves), g, p}
	go func(st CubieCube) {
		prefix := make([]Move, 0, depth)
		searchCubeBranch(st, res, depth, branch, prefix)
		close(res.channel)
	}(*c)
	return res
}

func searchCube(st CubieCube, s *cubeSearch, depth int, prefix []Move) {
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
		searchCube(newState, s, depth-1, append(prefix, move))
	}
}

func searchCubeBranch(st CubieCube, s *cubeSearch, depth, branch int,
	prefix []Move) {
	// If we shouldn't branch, do a regular search.
	if branch == 0 || depth == 0 {
		searchCube(st, s, depth, prefix)
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
		go func(m Move, newState CubieCube) {
			// Apply the move
			newState.Move(m)

			// Create the new prefix by copying the old one.
			pref := make([]Move, len(prefix)+1, len(prefix)+depth)
			copy(pref, prefix)
			pref[len(prefix)] = m

			// Branch out and search.
			searchCubeBranch(newState, s, depth-1, branch-1, pref)

			wg.Done()
		}(move, st)
	}
	wg.Wait()
}

type cubeSearch struct {
	*simpleSearch
	goal   CubeGoal
	pruner CubePruner
}

func (c *cubeSearch) prune(state CubieCube) int {
	if c.pruner == nil {
		return 0
	}
	return c.pruner.MinMoves(state)
}
