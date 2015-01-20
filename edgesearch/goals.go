package edgesearch

import (
	"github.com/unixpickle/gocube"
)

type SolvedGoal struct{}

func (_ SolvedGoal) IsGoal(c gocube.CubieEdges) bool {
	return c.Solved()
}

type OrientGoal struct{}

func (_ OrientGoal) IsGoal(c gocube.CubieEdges) bool {
	for _, p := range c {
		if p.Flip {
			return false
		}
	}
	return true
}
