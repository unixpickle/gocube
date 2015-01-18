package gocube

type SolvedCubieCubeGoal struct{}

func (_ SolvedCubieCubeGoal) IsGoal(c CubieCube) bool {
	return c.Solved()
}
