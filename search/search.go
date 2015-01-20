package search

import (
	"errors"
	"github.com/unixpickle/gocube"
	"github.com/unixpickle/gocube/cubiesearch"
	"github.com/unixpickle/gocube/edgesearch"
)

var (
	ErrBadGoal      = errors.New("Goal was not the correct type")
	ErrBadHeuristic = errors.New("Heuristic was not the correct type")
	ErrBadStart     = errors.New("Start state was not a known type")
)

type Search interface {
	Cancel()
	Run(maxDepth int, distribute int) ([]gocube.Move, error)
}

func NewSearch(start interface{}, goal interface{}, heuristic interface{},
	moves []gocube.Move) (Search, error) {
	if x, ok := start.(gocube.CubieCube); ok {
		g, ok := goal.(cubiesearch.Goal)
		if !ok {
			return nil, ErrBadGoal
		}
		h, ok := heuristic.(cubiesearch.Heuristic)
		if !ok {
			return nil, ErrBadHeuristic
		}
		return cubiesearch.NewSearch(x, g, h, moves), nil
	}
	if x, ok := start.(gocube.CubieEdges); ok {
		g, ok := goal.(edgesearch.Goal)
		if !ok {
			return nil, ErrBadGoal
		}
		h, ok := heuristic.(edgesearch.Heuristic)
		if !ok {
			return nil, ErrBadHeuristic
		}
		return edgesearch.NewSearch(x, g, h, moves), nil
	}
	return nil, ErrBadStart
}
