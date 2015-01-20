package edgesearch

import (
	"github.com/unixpickle/gocube"
)

// A PhaseOneHeuristic records the move count for every phase-one state.
// A phase-one state is determined by the edge orientation case and the
// permutation of the E slice edges.
type PhaseOneHeuristic struct {
	ChooseMap map[int]int
	Counts    [1013760]int8
}

// MakePhaseOneHeuristic uses a breadth-first search to generate a
// PhaseOneHeuristic.
func MakePhaseOneHeuristic(moves []gocube.Move) *PhaseOneHeuristic {
	res := NewPhaseOneHeuristic()

	// Set every move count to -1 so we know which one's we haven't visited.
	for i := 0; i < 1013760; i++ {
		res.Counts[i] = -1
	}

	// Perform a very basic breadth-first search.
	nodes := []searchNode{searchNode{gocube.SolvedCubieEdges(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := res.Encode(node.State)

		if res.Counts[idx] >= 0 {
			// Node was already visited
			continue
		}

		// Expand the node
		res.Counts[idx] = int8(node.Depth)
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, searchNode{state, node.Depth + 1})
		}
	}

	return res
}

// NewPhaseOneHeuristic creates a PhaseOneHeuristic with zeroes for every count
// and a generated ChooseMap.
func NewPhaseOneHeuristic() *PhaseOneHeuristic {
	res := PhaseOneHeuristic{ChooseMap: map[int]int{}}

	// Generate the map of all cases for (12 choose 4). Here are some examples,
	// just so you get the idea:
	// [w|x|y|z| | | | | | | | ], [w| |x|y| | | | |z| | | ], etc.
	chooseVal := 0
	for w := uint(0); w < 9; w++ {
		for x := w + 1; x < 10; x++ {
			for y := x + 1; y < 11; y++ {
				for z := y + 1; z < 12; z++ {
					// The sub-optimal but simple encoding for the choice.
					fromVal := (1 << w) | (1 << x) | (1 << y) | (1 << z)
					res.ChooseMap[fromVal] = chooseVal
					chooseVal++
				}
			}
		}
	}

	return &res
}

// Encode returns a hash number between 0 (inclusive) and 1013760 (exclusive)
// corresponding to the given edge configuration.
func (h *PhaseOneHeuristic) Encode(c gocube.CubieEdges) int {
	// Turn the edges into a bitmap with 1's where there are E slice edges.
	chooseFlags := 0
	for i, p := range c {
		if p.Piece == 1 || p.Piece == 3 || p.Piece == 7 || p.Piece == 9 {
			chooseFlags |= 1 << uint(i)
		}
	}

	// Generate the two numerical components of the hash.
	chooseVal := h.ChooseMap[chooseFlags]
	eoVal := encodeEO(c)

	return chooseVal*0x800 + eoVal
}

// MinMoves returns the minimum number of moves needed to reach the solved
// phase one state for a given set of edges.
func (h *PhaseOneHeuristic) MinMoves(c gocube.CubieEdges) int {
	return int(h.Counts[h.Encode(c)])
}

// An OrientHeuristic stores the number of moves required to solve every EO
// configuration.
type OrientHeuristic [0x800]int

// MakeOrientHeuristic uses breadth-first search to generate an OrientHeuristic.
// This should run pretty fast.
func MakeOrientHeuristic(moves []gocube.Move) *OrientHeuristic {
	var res OrientHeuristic

	for i := 0; i < 0x800; i++ {
		res[i] = -1
	}

	// Perform a very basic breadth-first search.
	nodes := []searchNode{searchNode{gocube.SolvedCubieEdges(), 0}}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		idx := encodeEO(node.State)

		if res[idx] >= 0 {
			// Node was already visited
			continue
		}

		// Expand the node
		res[idx] = node.Depth
		for _, move := range moves {
			state := node.State
			state.Move(move)
			nodes = append(nodes, searchNode{state, node.Depth + 1})
		}
	}

	return &res
}

// MinMoves returns the number of moves required to orient all the edges.
func (o *OrientHeuristic) MinMoves(e gocube.CubieEdges) int {
	return o[encodeEO(e)]
}

type searchNode struct {
	State gocube.CubieEdges
	Depth int
}

func encodeEO(e gocube.CubieEdges) int {
	res := 0
	for i := 0; i < 11; i++ {
		if e[i].Flip {
			res |= (1 << uint(i))
		}
	}
	return res
}
