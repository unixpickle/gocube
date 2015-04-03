package gocube

// Phase1Heuristic stores the data needed to effectively prune the search for a
// solution for phase-1.
type Phase1Heuristic struct {
	// This stores the number of moves needed to orient the corners. All the
	// fields in this array should be filled in.
	CO [2187]int8

	// This stores the number of moves needed to solve a given EO + E Slice
	// case. The fields are combined using slice*2048 + eo. If a value in this
	// array is -1, it should be interpreted as 8. This is because depths
	// greater than 7 will not be searched if the index was made under time
	// pressure.
	EOSlice [1013760]int8
}

// NewPhase1Heuristic generates a heuristic for the phase-1 solver.
// The complete argument specifies if the entire EOSlice table should be
// generated.
func NewPhase1Heuristic(moves *Phase1Moves, complete bool) *Phase1Heuristic {
	res := new(Phase1Heuristic)
	res.computeCO(moves)
	res.computeEOSlice(moves, complete)
	return res
}

// LowerBound returns the minimum number of moves needed to solve at least one
// axis of a Phase1Cube.
func (p *Phase1Heuristic) LowerBound(c *Phase1Cube) int {
	var result int8
	
	// Corner orientation heuristic.
	if r := p.CO[c.XCornerOrientation]; r > result {
		result = r
	}
	if r := p.CO[c.YCornerOrientation]; r > result {
		result = r
	}
	if r := p.CO[c.ZCornerOrientation]; r > result {
		result = r
	}
	
	// EOSlice heuristic.
	sliceValues := []int{
		c.MSlicePermutation*2048 + c.XEdgeOrientation(),
		c.ESlicePermutation*2048 + c.FBEdgeOrientation,
		c.SSlicePermutation*2048 + c.UDEdgeOrientation,
	}
	for _, eoSlice := range sliceValues {
		if r := p.EOSlice[eoSlice]; r > result {
			result = r
		} else if r < 0 && result < 8 {
			result = 8
		}
	}
	
	return int(result)
}

func (p *Phase1Heuristic) computeCO(moves *Phase1Moves) {
	for i := 0; i < 2187; i++ {
		p.CO[i] = -1
	}
	nodes := []phase1CONode{phase1CONode{1093, 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		if p.CO[node.corners] != -1 {
			continue
		}
		p.CO[node.corners] = node.depth
		for move := 0; move < 18; move++ {
			applied := moves.COMoves[node.corners][move]
			nodes = append(nodes, phase1CONode{applied, node.depth + 1})
		}
	}
}

func (p *Phase1Heuristic) computeEOSlice(moves *Phase1Moves, complete bool) {
	for i := 0; i < 1013760; i++ {
		p.EOSlice[i] = -1
	}
	nodes := []phase1EOSliceNode{phase1EOSliceNode{0, 220, 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		hash := node.slice*2048 + node.eo
		if p.EOSlice[hash] != -1 {
			continue
		}
		p.EOSlice[hash] = node.depth

		// Stop searching after 7 moves, which makes indexing faster and should
		// not greatly affect search.
		if !complete && node.depth == 7 {
			continue
		}

		for move := 0; move < 18; move++ {
			newEO := moves.EOMoves[node.eo][move]
			newSlice := moves.ESliceMoves[node.slice][move]
			node := phase1EOSliceNode{newEO, newSlice, node.depth + 1}
			nodes = append(nodes, node)
		}
	}
}

type phase1CONode struct {
	corners int
	depth   int8
}

type phase1EOSliceNode struct {
	eo    int
	slice int
	depth int8
}
