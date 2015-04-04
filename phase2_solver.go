package gocube

// A Phase2Heuristic estimates a lower bound for the number of moves to solve a
// Phase2Cube.
type Phase2Heuristic struct {
	// If an element is -1, it should be assumed to have the value 12.
	CornersSlice [967680]int8

	// If an element is -1, it should be assumed to have the value 9.
	EdgesSlice [967680]int8
}

// NewPhase2Heuristic generates a Phase2Heuristic.
// If complete is true, the full index is found. Otherwise, corners will only
// be searched up to depth 11, and edges will only be searched up to depth 8.
func NewPhase2Heuristic(moves *Phase2Moves, complete bool) *Phase2Heuristic {
	res := new(Phase2Heuristic)

	// Make all the move counts -1 by default.
	for i := 0; i < 967680; i++ {
		res.CornersSlice[i] = -1
		res.EdgesSlice[i] = -1
	}

	// Generate CornersSlice
	nodes := []phase2Node{phase2Node{0, 0, 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		if res.CornersSlice[node.hash()] >= 0 {
			continue
		}
		res.CornersSlice[node.hash()] = node.depth
		if !complete && node.depth == 11 {
			continue
		}
		for m := 0; m < 10; m++ {
			p4 := moves.SliceMoves[node.perm4][m]
			p8 := moves.CornerMoves[node.perm8][m]
			nodes = append(nodes, phase2Node{p4, p8, node.depth + 1})
		}
	}

	// Generate EdgesSlice
	nodes = []phase2Node{phase2Node{0, 0, 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		if res.EdgesSlice[node.hash()] >= 0 {
			continue
		}
		res.EdgesSlice[node.hash()] = node.depth
		if !complete && node.depth == 8 {
			continue
		}
		for m := 0; m < 10; m++ {
			p4 := moves.SliceMoves[node.perm4][m]
			p8 := moves.EdgeMoves[node.perm8][m]
			nodes = append(nodes, phase2Node{p4, p8, node.depth + 1})
		}
	}

	return res
}

// LowerBound returns the heuristic lower bound for a given Phase2Cube.
func (p *Phase2Heuristic) LowerBound(c *Phase2Cube) int {
	cornersSlice := c.CornerPermutation*24 + c.SlicePermutation
	edgesSlice := c.EdgePermutation*24 + c.SlicePermutation
	cMoves := p.CornersSlice[cornersSlice]
	eMoves := p.EdgesSlice[edgesSlice]

	if cMoves < 0 {
		cMoves = 12
	}
	if eMoves < 0 {
		eMoves = 9
	}

	if eMoves > cMoves {
		return int(eMoves)
	} else {
		return int(cMoves)
	}
}

// SolvePhase2 finds the first solution to a Phase2Cube, or gives up after
// maxLen moves.
func SolvePhase2(cube Phase2Cube, maxLen int, heuristic *Phase2Heuristic,
	moves *Phase2Moves) []Phase2Move {
	for depth := 0; depth <= maxLen; depth++ {
		if x := depthFirstPhase2(cube, depth, heuristic, moves, 0); x != nil {
			return x
		}
	}
	return nil
}

func depthFirstPhase2(cube Phase2Cube, depth int, heuristic *Phase2Heuristic,
	moves *Phase2Moves, lastFace int) []Phase2Move {
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
		res := depthFirstPhase2(c, depth-1, heuristic, moves, m.Face())
		if res != nil {
			return append([]Phase2Move{m}, res...)
		}
	}

	return nil
}

type phase2Node struct {
	perm4 int
	perm8 int
	depth int8
}

func (n *phase2Node) hash() int {
	return n.perm8*24 + n.perm4
}
