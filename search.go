package gocube

import (
	"sync/atomic"
)

// A Search is returned by search functions.
// An instance of Search controls a set of background routines which are finding
// solutions to a given problem.
type Search interface {
	// Cancel cancels the search.
	// This will return once the search has reached a full stop.
	Cancel()

	// Solutions returns a channel to which solutions will be sent.
	// The channel will be closed once the search terminates or is cancelled.
	Solutions() <-chan []Move
}

// WaitSearch waits until a search terminates or is cancelled.
// This is achieved by continually reading from the Solutions() channel until it
// is closed.
// All of the solutions which were read from the Solutions() channel are
// returned.
func WaitSearch(s Search) [][]Move {
	// Read each value until the
	result := make([][]Move, 0)
	solutions := s.Solutions()
	for {
		solution, ok := <-solutions
		if !ok {
			break
		}
		result = append(result, solution)
	}
	return result
}

type simpleSearch struct {
	// cancelFlag is 0 for false and 1 for true.
	// It is a uint32 so that it can be used with the "atomic" package.
	cancelFlag uint32

	moves   []Move
	channel chan []Move
}

func newSimpleSearch(moves []Move) *simpleSearch {
	return &simpleSearch{0, moves, make(chan []Move)}
}

func (s *simpleSearch) Cancel() {
	// Tell the background threads to stop.
	atomic.StoreUint32(&s.cancelFlag, 1)

	// Wait for the solutions channel to be closed.
	for {
		_, ok := <-s.channel
		if !ok {
			return
		}
	}
}

func (s *simpleSearch) Solutions() <-chan []Move {
	return s.channel
}

func (s *simpleSearch) cancelled() bool {
	return atomic.LoadUint32(&s.cancelFlag) != 0
}
