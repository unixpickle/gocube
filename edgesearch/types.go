package edgesearch

import (
	"github.com/unixpickle/gocube"
)

type Goal interface {
	IsGoal(state gocube.CubieEdges) bool
}

type Heuristic interface {
	MinMoves(state gocube.CubieEdges) int
}
