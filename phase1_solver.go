package gocube

// Phase1Heuristic stores the data needed to effectively prune the search for a
// solution for phase-1.
type Phase1Heuristic struct {
	// This stores the number of moves needed to orient the corners and edges.
	COEO [4478976]uint8

	// This stores the number of moves needed orient the edges and put the slice
	// edges on the slice.
	EOSlice [1013760]uint8
}

// NewPhase1Heuristic generates a heuristic for the phase-1 solver.
func NewPhase1Heuristic(moves *Phase1Moves) *Phase1Heuristic {
	res := new(Phase1Heuristic)
	res.computeCOEO(moves)
	res.computeEOSlice(moves)
	return res
}

// LowerBound returns the minimum number of moves needed to solve at least one
// phase-1 axis.
func (p *Phase1Heuristic) LowerBound(c *Phase1Cube) int {
	finalResult := uint8(127)

	coEOValues := []int{
		c.XCornerOrientation*2048 + c.XEdgeOrientation(),
		c.YCornerOrientation*2048 + c.FBEdgeOrientation,
		c.ZCornerOrientation*2048 + c.UDEdgeOrientation,
	}
	sliceValues := []int{
		c.MSlicePermutation*2048 + c.XEdgeOrientation(),
		c.ESlicePermutation*2048 + c.FBEdgeOrientation,
		c.SSlicePermutation*2048 + c.UDEdgeOrientation,
	}

	for axis := 0; axis < 3; axis++ {
		result := p.COEO[coEOValues[axis]]
		if r := p.EOSlice[sliceValues[axis]]; r > result {
			result = r
		}
		if result < finalResult {
			finalResult = result
		}
	}

	return int(finalResult)
}

func (p *Phase1Heuristic) computeCOEO(moves *Phase1Moves) {
	for i := 0; i < 4478976; i++ {
		p.COEO[i] = 8
	}

	nodes := []phase1COEONode{{1093, 0, 0}}
	visited := make([]bool, 4478976)
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		hash := node.co*2048 + node.eo
		if p.COEO[hash] != 8 {
			continue
		}
		p.COEO[hash] = node.depth

		if node.depth == 7 {
			continue
		}

		for move := 0; move < 18; move++ {
			newCO := moves.COMoves[node.co][move]
			newEO := moves.EOMoves[node.eo][move]
			newHash := newCO*2048 + newEO
			if !visited[newHash] {
				nodes = append(nodes, phase1COEONode{newCO, newEO,
					node.depth + 1})
				visited[newHash] = true
			}
		}
	}
}

func (p *Phase1Heuristic) computeEOSlice(moves *Phase1Moves) {
	for i := 0; i < 1013760; i++ {
		p.EOSlice[i] = 8
	}
	nodes := []phase1EOSliceNode{{0, 220, 0}}
	visited := make([]bool, 1013760)
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		hash := node.slice*2048 + node.eo
		if p.EOSlice[hash] != 8 {
			continue
		}
		p.EOSlice[hash] = node.depth

		if node.depth == 7 {
			continue
		}

		for move := 0; move < 18; move++ {
			newEO := moves.EOMoves[node.eo][move]
			newSlice := moves.ESliceMoves[node.slice][move]
			newHash := newSlice*2048 + newEO
			if !visited[newHash] {
				newNode := phase1EOSliceNode{newEO, newSlice, node.depth + 1}
				nodes = append(nodes, newNode)
				visited[newHash] = true
			}
		}
	}
}

// A Phase1Solution stores information about a phase-1 solution.
type Phase1Solution struct {
	Cube  Phase1Cube
	Moves []Move
}

// A Phase1Solver finds solutions to a specific phase-1 state.
type Phase1Solver struct {
	stopped   chan struct{}
	solutions <-chan Phase1Solution

	heuristic *Phase1Heuristic
	moves     *Phase1Moves
}

// NewPhase1Solver creates and starts a Phase1Solver.
func NewPhase1Solver(c Phase1Cube, h *Phase1Heuristic,
	m *Phase1Moves) *Phase1Solver {
	solutions := make(chan Phase1Solution)
	res := &Phase1Solver{make(chan struct{}), solutions, h, m}
	go res.search(solutions, c)
	return res
}

func (p *Phase1Solver) Solutions() <-chan Phase1Solution {
	return p.solutions
}

func (p *Phase1Solver) Stop() {
	close(p.stopped)
}

func (p *Phase1Solver) depthFirst(solutions chan<- Phase1Solution, c Phase1Cube,
	moves []Move, depth int, lastFace int) bool {
	// If the depth is zero, we may have a solution.
	if depth == 0 {
		if c.AnySolved() {
			res := make([]Move, len(moves))
			copy(res, moves)
			select {
			case <-p.stopped:
				return false
			case solutions <- Phase1Solution{c, res}:
			}
		}
		return true
	}

	// Check the heuristic.
	if p.heuristic.LowerBound(&c) > depth {
		return true
	}

	// Apply every move and recurse.
	for m := 0; m < 18; m++ {
		move := Move(m)
		if move.Face() == lastFace {
			continue
		}
		cube := c
		cube.Move(move, p.moves)
		moves = append(moves, move)
		if !p.depthFirst(solutions, cube, moves, depth-1, move.Face()) {
			return false
		}
		moves = moves[:len(moves)-1]
		if depth >= 7 && p.isStopped() {
			return false
		}
	}

	return true
}

func (p *Phase1Solver) isStopped() bool {
	select {
	case <-p.stopped:
		return true
	default:
		return false
	}
}

func (p *Phase1Solver) search(solutions chan<- Phase1Solution, c Phase1Cube) {
	depth := 0
	for {
		moves := make([]Move, 0, depth)
		if !p.depthFirst(solutions, c, moves, depth, 0) {
			close(solutions)
			return
		}
		depth++
	}
}

type phase1COEONode struct {
	co    int
	eo    int
	depth uint8
}

type phase1EOSliceNode struct {
	eo    int
	slice int
	depth uint8
}
