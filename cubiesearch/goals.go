package cubiesearch

import (
	"github.com/unixpickle/gocube"
)

type SolvedGoal struct{}

func (_ SolvedGoal) IsGoal(c gocube.CubieCube) bool {
	return c.Solved()
}
