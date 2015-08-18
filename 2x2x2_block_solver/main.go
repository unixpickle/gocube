package main

import (
	"fmt"
	"github.com/unixpickle/gocube"
	"os"
)

func main() {
	sc, err := gocube.InputStickerCube()
	if err != nil {
		fmt.Println("Failed to read stickers:", err)
		os.Exit(1)
	}
	cc, err := sc.CubieCube()
	if err != nil {
		fmt.Println("Invalid stickers:", err)
		os.Exit(1)
	}

	for depth := 0; depth <= 20; depth++ {
		fmt.Println("Searching depth", depth)
		if solution := Search(*cc, depth); solution != nil {
			fmt.Println("Got a solution:", solution)
			break
		}
	}
}

func Search(start gocube.CubieCube, d int) []gocube.Move {
	if d == 0 {
		if IsSolved(start) {
			return []gocube.Move{}
		} else {
			return nil
		}
	}
	for m := 0; m < 18; m++ {
		newCube := start
		newCube.Move(gocube.Move(m))
		if solution := Search(newCube, d-1); solution != nil {
			return append([]gocube.Move{gocube.Move(m)}, solution...)
		}
	}
	return nil
}

func IsSolved(state gocube.CubieCube) bool {
	isCornerSolved := func(idx int) bool {
		return state.Corners[idx].Piece == idx &&
			state.Corners[idx].Orientation == 1
	}
	isEdgeSolved := func(idx int) bool {
		return state.Edges[idx].Piece == idx && !state.Edges[idx].Flip
	}
	edgesForCorners := []int{
		8, 9, 10,
		7, 8, 11,
		4, 6, 9,
		5, 6, 7,
		2, 3, 10,
		1, 2, 11,
		0, 3, 4,
		0, 1, 5,
	}
OuterLoop:
	for corner := 0; corner < 8; corner++ {
		if isCornerSolved(corner) {
			for i := corner*3; i < corner*3+3; i++ {
				if !isEdgeSolved(edgesForCorners[i]) {
					continue OuterLoop
				}
			}
			fmt.Println("solved corner", corner)
			return true
		}
	}
	return false
}
