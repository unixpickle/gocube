package gocube

import (
	"errors"
	"sync"
)

var (
	ErrAlreadySearching = errors.New("Already searching.")
	ErrNoSolution       = errors.New("No solution was found.")
)

type CubieCubeGoal interface {
	IsGoal(state CubieCube) bool
}

type CubieCubeHeuristic interface {
	MinMoves(state CubieCube) int
}

// A CubieCubeSearch can perform a cancellable depth-first search.
type CubieCubeSearch struct {
	start     CubieCube
	goal      CubieCubeGoal
	heuristic CubieCubeHeuristic
	moves     []Move

	lock      sync.RWMutex
	cancelled bool
	running   bool
}

// NewCubieCubeSearch creates an unstarted search with the given parameters.
func NewCubieCubeSearch(c CubieCube, g CubieCubeGoal, h CubieCubeHeuristic,
	m []Move) *CubieCubeSearch {
	mCopy := make([]Move, len(m))
	copy(mCopy, m)
	return &CubieCubeSearch{c, g, h, mCopy, sync.RWMutex{}, false, false}
}

// Cancel cancels the current search.
// This is useful if you want to cancel the search from a different goroutine
// than the one its running on.
func (s *CubieCubeSearch) Cancel() {
	s.lock.Lock()
	s.cancelled = true
	s.lock.Unlock()
}

// Run starts a depth-first search and waits for it to complete or be cancelled.
// If the search is already running on a different thread, this returns
// ErrAlreadySearching.
// If no solution was found either because the search was exhausted or it was
// cancelled, ErrNoSolution is returned.
func (s *CubieCubeSearch) Run(maxDepth int, distribute int) ([]Move, error) {
	// Set the flags so that no other search can be run simultaneously.
	s.lock.Lock()
	if s.running {
		s.lock.Unlock()
		return nil, ErrAlreadySearching
	}
	s.running = true
	s.cancelled = false
	s.lock.Unlock()

	res := s.search(maxDepth, distribute)

	// Reset the running flag so another search can run.
	s.lock.Lock()
	s.running = false
	s.lock.Unlock()

	// Return the appropriate value and error.
	if res != nil {
		return res, nil
	} else {
		return nil, ErrNoSolution
	}
}

func (s *CubieCubeSearch) distSearch(st CubieCube, max int, d int) []Move {
	// If we can't search any further, check if it's the goal.
	if max == 0 {
		if s.goal.IsGoal(st) {
			return []Move{}
		}
		return nil
	}

	// Prune the state
	pruneVal := s.heuristic.MinMoves(st)
	if pruneVal > max {
		return nil
	}

	// Run each search on a different goroutine
	wg := sync.WaitGroup{}
	solutionLock := sync.Mutex{}
	var solution []Move
	for _, move := range s.moves {
		wg.Add(1)
		newState := st
		newState.Move(move)
		go func() {
			var res []Move
			if d == 1 {
				res = s.regularSearch(newState, max-1)
			} else {
				res = s.distSearch(newState, max-1, d-1)
			}
			if res != nil {
				solutionLock.Lock()
				solution = res
				solutionLock.Unlock()
				s.Cancel()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return solution
}

func (s *CubieCubeSearch) regularSearch(st CubieCube, max int) []Move {
	// If we can't search any further, check if it's the goal.
	if max == 0 {
		if s.goal.IsGoal(st) {
			return []Move{}
		}
		return nil
	}

	// Prune the state
	pruneVal := s.heuristic.MinMoves(st)
	if pruneVal > max {
		return nil
	}

	// Apply each move and recurse.
	for _, move := range s.moves {
		if max > 5 && s.shouldStop() {
			return nil
		}
		newState := st
		newState.Move(move)
		if res := s.regularSearch(newState, max-1); res != nil {
			return append([]Move{move}, res...)
		}
	}

	return nil
}

func (s *CubieCubeSearch) search(max int, distribute int) []Move {
	if distribute > 0 {
		return s.distSearch(s.start, max, distribute)
	} else {
		return s.regularSearch(s.start, max)
	}
}

func (s *CubieCubeSearch) shouldStop() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cancelled
}
