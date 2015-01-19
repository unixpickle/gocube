package cubiesearch

import (
	"github.com/unixpickle/gocube"
	"sync"
)

// A Search can perform a cancellable depth-first search.
type Search struct {
	start     gocube.CubieCube
	goal      Goal
	heuristic Heuristic
	moves     []gocube.Move

	lock      sync.RWMutex
	cancelled bool
	running   bool
}

// NewSearch creates an unstarted search with the given parameters.
func NewSearch(c gocube.CubieCube, g Goal, h Heuristic,
	m []gocube.Move) *Search {
	mCopy := make([]gocube.Move, len(m))
	copy(mCopy, m)
	return &Search{c, g, h, mCopy, sync.RWMutex{}, false, false}
}

// Cancel cancels the current search.
// This is useful if you want to cancel the search from a different goroutine
// than the one its running on.
func (s *Search) Cancel() {
	s.lock.Lock()
	s.cancelled = true
	s.lock.Unlock()
}

// Run starts a depth-first search and waits for it to complete or be cancelled.
// If the search is already running on a different thread, this returns
// ErrAlreadySearching.
// If no solution was found either because the search was exhausted or it was
// cancelled, ErrNoSolution is returned.
func (s *Search) Run(maxDepth int, distribute int) ([]gocube.Move, error) {
	// Set the flags so that no other search can be run simultaneously.
	s.lock.Lock()
	if s.running {
		s.lock.Unlock()
		return nil, ErrAlreadySearching
	}
	s.running = true
	s.cancelled = false
	s.lock.Unlock()

	// Perform the search itself.
	var res []gocube.Move
	if distribute > 0 {
		res = s.distSearch(s.start, maxDepth, distribute)
	} else {
		res = s.regularSearch(s.start, maxDepth)
	}

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

func (s *Search) distSearch(st gocube.CubieCube, max int, d int) []gocube.Move {
	// If we can't search any further, check if it's the goal.
	if max == 0 {
		if s.goal.IsGoal(st) {
			return []gocube.Move{}
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
	var solution []gocube.Move
	for _, move := range s.moves {
		wg.Add(1)
		newState := st
		newState.Move(move)
		go func(m gocube.Move) {
			var res []gocube.Move
			if d == 1 {
				res = s.regularSearch(newState, max-1)
			} else {
				res = s.distSearch(newState, max-1, d-1)
			}
			if res != nil {
				solutionLock.Lock()
				solution = append([]gocube.Move{m}, res...)
				solutionLock.Unlock()
				s.Cancel()
			}
			wg.Done()
		}(move)
	}

	wg.Wait()
	return solution
}

func (s *Search) regularSearch(st gocube.CubieCube, max int) []gocube.Move {
	// If we can't search any further, check if it's the goal.
	if max == 0 {
		if s.goal.IsGoal(st) {
			return []gocube.Move{}
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
			return append([]gocube.Move{move}, res...)
		}
	}

	return nil
}

func (s *Search) shouldStop() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cancelled
}
