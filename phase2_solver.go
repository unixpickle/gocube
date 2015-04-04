package gocube

// A Phase2Heuristic estimates a lower bound for the number of moves to solve a
// Phase2Cube.
type Phase2Heuristic struct {
	Corners [40320]int8
}

// NewPhase2Heuristic generates a Phase2Heuristic.
func NewPhase2Heuristic(moves *Phase2Moves) *Phase2Heuristic {
	res := new(Phase2Heuristic)

	// Generate the corners.
	for i := 0; i < 40320; i++ {
		res.Corners[i] = -1
	}
	nodes := []phase2CornerNode{phase2CornerNode{0, 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		if res.Corners[node.corners] >= 0 {
			continue
		}
		res.Corners[node.corners] = node.depth
		for m := 0; m < 10; m++ {
			state := moves.CornerMoves[node.corners][m]
			nodes = append(nodes, phase2CornerNode{state, node.depth + 1})
		}
	}

	return res
}

// LowerBound returns the heuristic lower bound for a given Phase2Cube.
func (p *Phase2Heuristic) LowerBound(c *Phase2Cube) int {
	return int(p.Corners[c.CornerPermutation])
}

// SolvePhase2 finds the first solution to a Phase2Cube, or gives up after
// maxLen moves.
func SolvePhase2(cube Phase2Cube, maxLen int, moves *Phase2Moves,
	heuristic *Phase2Heuristic) []Phase2Move {
	for depth := 0; depth <= maxLen; depth++ {
		return depthFirstPhase2(cube, depth, moves, heuristic, 0)
	}
	return nil
}

func depthFirstPhase2(cube Phase2Cube, depth int, moves *Phase2Moves,
	heuristic *Phase2Heuristic, lastFace int) []Phase2Move {
	if depth == 0 {
		if cube.Solved() {
			return []Phase2Move{}
		}
	} else if heuristic.LowerBound(&cube) > depth {
		return nil
	}

	// Apply moves and recurse.
	for i := 0; i < 10; i++ {
		m := Phase2Move(i)
		if m.Face() == lastFace {
			continue
		}
		c := cube
		c.Move(m, moves)
		res := depthFirstPhase2(c, depth-1, moves, heuristic, m.Face())
		if res != nil {
			return append([]Phase2Move{m}, res...)
		}
	}

	return nil
}

type phase2CornerNode struct {
	corners int
	depth   int8
}
