package edgesearch

import (
	"github.com/unixpickle/gocube"
)

type SolvedGoal struct{}

func (_ SolvedGoal) IsGoal(c gocube.CubieEdges) bool {
	return c.Solved()
}
