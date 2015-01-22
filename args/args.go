package args

import (
	"flag"
	"github.com/unixpickle/gocube"
)

// Args stores the arguments which can be applied to pretty much any search.
type Args struct {
	scrambleStr string
	movesStr    string
	scramble    *gocube.CubieCube
	moves       []gocube.Move
	branch      *int
	maxDepth    *int
	minDepth    *int
	multiple    *bool
}

// NewArgs creates a new Args object given a FlagSet.
// This will not automatically parse the FlagSet.
func NewArgs(f *flag.FlagSet) *Args {
	res := &Args{}
	f.StringVar(&res.scrambleStr, "scramble", "", "The scramble in WCA notation")
	f.StringVar(&res.movesStr, "moves", "", "The basis for the search")
	res.branch = f.Int("branch", 1, "The parallelization depth")
	res.maxDepth = f.Int("maxdepth", 20, "The maximum search depth")
	res.minDepth = f.Int("mindepth", 0, "The minimum search depth")
	res.multiple = f.Bool("multiple", false, "Find multiple solutions")
	return res
}

// Branch returns the branch depth argument.
func (a *Args) Branch() int {
	return *a.branch
}

// MaxDepth returns the maximum depth argument.
func (a *Args) MaxDepth() int {
	return *a.maxDepth
}

// MinDepth returns the minimum depth argument.
func (a *Args) MinDepth() int {
	return *a.minDepth
}

// Moves returns the basis moves.
func (a *Args) Moves() []gocube.Move {
	return a.moves
}

// Multiple returns the multiple solutions argument.
func (a *Args) Multiple() bool {
	return *a.multiple
}

// Parse processes the supplied arguments and (possibly) asks the user for
// clarifying input.
// You should call this after calling parse on the FlagSet.
func (a *Args) Parse() error {
	// Create the scramble
	if len(a.scrambleStr) > 0 {
		// Parse the scramble they used as an argument.
		moves, err := gocube.ParseMoves(a.scrambleStr)
		if err != nil {
			return err
		}
		solved := gocube.SolvedCubieCube()
		a.scramble = &solved
		for _, move := range moves {
			a.scramble.Move(move)
		}
	} else {
		// Input a scramble from the console.
		input, err := gocube.InputStickerCube()
		if err != nil {
			return err
		}
		a.scramble, err = input.CubieCube()
		if err != nil {
			return err
		}
	}

	// Create the basis
	if len(a.movesStr) > 0 {
		// Parse the supplied basis
		var err error
		a.moves, err = gocube.ParseMoves(a.movesStr)
		if err != nil {
			return err
		}
	} else {
		// Order the default moves in terms of comfort.
		a.moves, _ = gocube.ParseMoves("R U L D R' U' L' D' R2 U2 L2 D2 " +
			"F B F' B' F2 B2")
	}

	return nil
}

// Scramble gets the scramble from the arguments.
func (a *Args) Scramble() gocube.CubieCube {
	return *a.scramble
}
