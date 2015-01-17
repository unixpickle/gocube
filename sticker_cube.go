package gocube

import (
	"errors"
	"fmt"
	"strconv"
)

// A StickerCube is simply a list of 54 stickers.
//
// Each sticker is a number between 1 and 6. Here is a mapping for a standard
// cube: 1=white, 2=yellow, 3=green, 4=blue, 5=red, 6=orange.
//
// The order of the stickers is well defined but slightly tricky to memorize.
// The stickers are grouped by face, so the first 9 correspond to the top, the
// next to the bottom, next to the front, then the back, the right, then the
// left.
//
// The order of the stickers on a given face are well defined as well. If you do
// entry as I describe below, you will always type colors from left to right,
// top to bottom as if you were reading a book:
//
// * First, hold the cube so that the side you wish to be in front is in front,
//   and the side you wish to be on top is on top. Now perform an x' so that the
//   top side is in the front and enter the (now) front side by reading the top
//   left color, then the one to its right, then the one to its right, then the
//   far left sticker on the second row, etc. This is the same way you would
//   read a book.
// * Now perform an x2 so that the original bottom side is now in the front.
//   Enter this side the same way as above.
// * Now do an x' to reset the orientation of the cube, and type the front side
//   the same way as you entered the other two.
// * Now do a y2 and enter the back side the same way as the other three sides.
// * Now do a y' so that the original right side is in front, and enter it in.
// * Now do a y2 and enter the original left side.
type StickerCube [54]int

// InputStickerCube reads user input for a sticker cube.
// The cube is not validated beyond checking that each sticker occurs 9 times.
func (s StickerCube) InputStickerCube() (*StickerCube, error) {
	// Print out some helpful information.
	fmt.Println("Enter a cube. Color codes (if it helps):")
	fmt.Println(" 1=W=white, 2=Y=yellow, 3=G=green, 4=B=blue,",
		"5=R=red, 6=O=orange")
	fmt.Println("For each face, enter what you see from left to right, top to")
	fmt.Println("bottom like reading a book.")

	// These are the six prompts (one for each side)
	prompts := []string{
		"   Top face [x']: ",
		"Bottom face [x2]: ",
		" Front face [x']: ",
		"  Back face [y2]: ",
		" Right face [y']: ",
		"  Left face [y2]: ",
	}

	// Print each prompt and read each face.
	var res StickerCube
	for i, prompt := range prompts {
		fmt.Print(prompt)
		startIdx := i * 9
		face, err := readFace(i + 1)
		if err != nil {
			return nil, err
		}
		copy(res[startIdx:startIdx+9], face)
	}

	// Validate the result
	counts := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}
	for _, sticker := range res {
		counts[sticker]++
	}
	for i, v := range counts {
		if v != 9 {
			return nil, errors.New("Invalid number of sticker " +
				strconv.Itoa(i))
		}
	}

	return &res, nil
}

func readFace(number int) ([]int, error) {
	// Read the line
	var line string
	if _, err := fmt.Scanln(&line); err != nil {
		return nil, err
	}
	runes := []rune(line)
	if len(runes) > 9 {
		return nil, errors.New("Face input was too long.")
	}

	// Process each character in the line
	res := make([]int, 9)
	for i, c := range runes {
		if c == '1' || c == 'W' {
			res[i] = 1
		} else if c == '2' || c == 'Y' {
			res[i] = 2
		} else if c == '3' || c == 'G' {
			res[i] = 3
		} else if c == '4' || c == 'B' {
			res[i] = 4
		} else if c == '5' || c == 'R' {
			res[i] = 5
		} else if c == '6' || c == 'O' {
			res[i] = 6
		} else {
			return nil, errors.New("Unexpected character: " + string(c))
		}
	}
	for i := len(runes); i < 9; i++ {
		res[i] = number
	}
	return res, nil
}
