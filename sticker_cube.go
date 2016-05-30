package gocube

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
// 1. First, hold the cube so that the side you wish to be in front is in front,
// and the side you wish to be on top is on top. Now perform an x' so that the
// top side is in the front and enter the (now) front side by reading the top
// left color, then the one to its right, then the one to its right, then the
// far left sticker on the second row, etc. This is the same way you would
// read a book.
//
// 2. Now perform an x2 so that the original bottom side is now in the front.
// Enter this side the same way as above.
//
// 3. Now do an x' to reset the orientation of the cube, and type the front side
// the same way as you entered the other two.
//
// 4. Now do a y2 and enter the back side the same way as the other three sides.
//
// 5. Now do a y' so that the original right side is in front, and enter it in.
//
// 6. Now do a y2 and enter the original left side.
type StickerCube [54]int

// InputStickerCube reads user input for a sticker cube.
// The cube is not validated beyond checking that each sticker occurs 9 times.
func InputStickerCube() (*StickerCube, error) {
	// Print out some helpful information.
	fmt.Println("Enter a cube. Color codes (if it helps):")
	fmt.Println(" 1=W=white, 2=Y=yellow, 3=G=green, 4=B=blue,",
		"5=R=red, 6=O=orange")
	fmt.Println("For each face, enter what you see from left to right, top to")
	fmt.Println("bottom like reading a book.")

	prompts := []string{
		"   Top face [x']: ",
		"Bottom face [x2]: ",
		" Front face [x']: ",
		"  Back face [y2]: ",
		" Right face [y']: ",
		"  Left face [y2]: ",
	}

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

	counts := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}
	for _, sticker := range res {
		counts[sticker]++
	}
	for i, v := range counts {
		if v != 9 {
			return nil, errors.New("invalid number of sticker " +
				strconv.Itoa(i))
		}
	}

	return &res, nil
}

// ParseStickerCube parses a space-delimited list of faces.
func ParseStickerCube(str string) (*StickerCube, error) {
	faceStrs := strings.Split(str, " ")
	if len(faceStrs) != 6 {
		return nil, errors.New("input must have six faces.")
	}

	var res StickerCube
	for i, faceStr := range faceStrs {
		list, err := parseFace(faceStr, i+1)
		if err != nil {
			return nil, err
		}
		copy(res[i*9:i*9+9], list)
	}
	return &res, nil
}

// SolvedStickerCube returns a solved sticker cube.
func SolvedStickerCube() StickerCube {
	var result StickerCube
	for i := 0; i < 54; i++ {
		result[i] = i/9 + 1
	}
	return result
}

// String generates a space-delimited list of faces in human-readable form.
func (s *StickerCube) String() string {
	res := ""
	for i := 0; i < 54; i++ {
		if i%9 == 0 && i != 0 {
			res += " "
		}
		num := s[i]
		strs := " WYGBRO"
		res += strs[num : num+1]
	}
	return res
}

func parseFace(str string, fill int) ([]int, error) {
	runes := []rune(str)
	if len(runes) > 9 {
		return nil, errors.New("face string cannot exceed nine characters.")
	}

	res := make([]int, 9)

	// Read each character of their input.
	for i, c := range runes {
		if c == '1' || c == 'W' || c == 'w' {
			res[i] = 1
		} else if c == '2' || c == 'Y' || c == 'y' {
			res[i] = 2
		} else if c == '3' || c == 'G' || c == 'g' {
			res[i] = 3
		} else if c == '4' || c == 'B' || c == 'b' {
			res[i] = 4
		} else if c == '5' || c == 'R' || c == 'r' {
			res[i] = 5
		} else if c == '6' || c == 'O' || c == 'o' {
			res[i] = 6
		} else {
			return nil, errors.New("unexpected character: " + string(c))
		}
	}

	// If they left out the ending, we fill it in with the face number.
	for i := len(runes); i < 9; i++ {
		res[i] = fill
	}

	return res, nil
}

func readFace(number int) ([]int, error) {
	var line string
	if _, err := fmt.Scanln(&line); err != nil {
		return nil, err
	}
	return parseFace(line, number)
}
