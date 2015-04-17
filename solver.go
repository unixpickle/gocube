package gocube

// A Solver finds shorter and shorter solutions in the background.
type Solver struct {
	stopper   chan struct{}
	solutions chan []Move
	phase1    *Phase1Solver
}

// NewSolver creates a new solver.
func NewSolver(c CubieCube, max int) *Solver {
	p1Moves := NewPhase1Moves()
	p1Heuristic := NewPhase1Heuristic(p1Moves)

	res := new(Solver)
	res.stopper = make(chan struct{})
	res.solutions = make(chan []Move)
	res.phase1 = NewPhase1Solver(c.Phase1Cube(), p1Heuristic, p1Moves)

	go res.backgroundLoop(c, nil, max)
	return res
}

// NewSolverTables creates a new solver using a set of pre-generated tables.
func NewSolverTables(c CubieCube, max int, tables SolverTables) *Solver {
	res := new(Solver)
	res.stopper = make(chan struct{})
	res.solutions = make(chan []Move)
	res.phase1 = NewPhase1Solver(c.Phase1Cube(), tables.P1Heuristic,
		tables.P1Moves)

	go res.backgroundLoop(c, &tables, max)
	return res
}

// Solutions is a channel over which shorter and shorter solutions are
// delivered.
func (s *Solver) Solutions() <-chan []Move {
	return s.solutions
}

// Stop stops the solver.
func (s *Solver) Stop() {
	s.phase1.Stop()
	close(s.stopper)
}

func (s *Solver) backgroundLoop(c CubieCube, tables *SolverTables, max int) {
	// Get the tables.
	var p2Moves *Phase2Moves
	var p2Heuristic *Phase2Heuristic
	if tables != nil {
		p2Moves = tables.P2Moves
		p2Heuristic = tables.P2Heuristic
	} else {
		p2Moves = NewPhase2Moves()
		p2Heuristic = NewPhase2Heuristic(p2Moves, false)
	}

OuterLoop:
	for p1Solution := range s.phase1.Solutions() {
		if len(p1Solution.Moves) > max {
			s.Stop()
			break
		}
		
		// Generate the cube after solving phase1.
		cube := c
		for _, m := range p1Solution.Moves {
			cube.Move(m)
		}
		// The phase-1 solution could be in the x, y, or z axis. We will go
		// through each axis and try solving it.
		x, y, z := p1Solution.Cube.Solved()
		for axis, solved := range []bool{x, y, z} {
			if !solved {
				continue
			}

			// Create the phase-2 cube and solve it.
			cube, err := NewPhase2Cube(cube, axis)
			if err != nil {
				continue
			}
			p2Solution := SolvePhase2(cube, max-len(p1Solution.Moves),
				p2Heuristic, p2Moves)
			if p2Solution == nil {
				continue
			}

			// Join the two solutions.
			solution := make([]Move, len(p1Solution.Moves))
			copy(solution, p1Solution.Moves)
			for _, move := range p2Solution {
				solution = append(solution, move.Move(axis))
			}
			max = len(solution) - 1
			select {
			case <-s.stopper:
				break OuterLoop
			case s.solutions <- solution:
			}
			if len(p2Solution) == 0 {
				break OuterLoop
			}
		}
	}
	close(s.solutions)
}

type SolverTables struct {
	P1Heuristic *Phase1Heuristic
	P1Moves     *Phase1Moves
	P2Heuristic *Phase2Heuristic
	P2Moves     *Phase2Moves
}
