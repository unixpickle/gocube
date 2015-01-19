package cubiesearch

import (
	"github.com/unixpickle/gocube"
)

type Goal interface {
	IsGoal(state gocube.CubieCube) bool
}

type Heuristic interface {
	MinMoves(state gocube.CubieCube) int
}
