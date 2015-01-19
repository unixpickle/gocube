package cubiesearch

import (
	"errors"
)

var (
	ErrAlreadySearching = errors.New("Already searching.")
	ErrNoSolution       = errors.New("No solution was found.")
)
