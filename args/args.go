package args

import (
	"flag"
	"github.com/unixpickle/gocube"
)

// Args stores the arguments which can be applied to pretty much any search.
type Args struct {
	scrambleStr string
	scramble    *gocube.CubieCube
}

// NewArgs creates a new Args object given a FlagSet.
// This will not automatically parse the FlagSet.
func NewArgs(f *flag.FlagSet) *Args {
	res := &Args{}
	f.StringVar(&res.scrambleStr, "scramble", "", "The WCA scramble.")
	return res
}

// Scramble gets the scramble from the arguments or reads it from standard
// input.
// The result is cached, so you may call Scramble() as many times as you wish.
func (a *Args) Scramble() (*gocube.CubieCube, error) {
	if a.scramble != nil {
		return a.scramble, nil
	} else if len(a.scrambleStr) > 0 {
		// Parse the scramble they used as an argument.
		moves, err := gocube.ParseMoves(a.scrambleStr)
		if err != nil {
			return nil, err
		}
		solved := gocube.SolvedCubieCube()
		a.scramble = &solved
		for _, move := range moves {
			a.scramble.Move(move)
		}
		return a.scramble, nil
	} else {
		// Input a scramble from the console.
		input, err := gocube.InputStickerCube()
		if err != nil {
			return nil, err
		}
		a.scramble, err = input.CubieCube()
		if err != nil {
			return nil, err
		}
		return a.scramble, nil
	}
}
