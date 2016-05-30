package gocube

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// XRotationPerm stores sticker indices for the
// faces involved in an "x" rotation.
//
// The first four sets of 9 stickers are for the
// four faces which are not normal to the x-axis.
// Each set of 9 consecutive stickers corresponds
// to a face, such that the stickers indexed by
// the first face are moved to the indices of the
// second face, and those to the third face, etc.
//
// The two normal faces are represented by two 8-perms
// at the end of the list.
// These 8-perms are applied twice to permutate the
// non-center pieces of the normal faces.
// The sticker at the first index of the 8-perm is moved
// to the next index, etc.
var XRotationPerm = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	35, 34, 33, 32, 31, 30, 29, 28, 27,
	9, 10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, 21, 22, 23, 24, 25, 26,

	36, 37, 38, 41, 44, 43, 42, 39,
	47, 46, 45, 48, 51, 52, 53, 50,
}

// YRotationPerm is like XRotationPerm, but for
// "y" rotations.
var YRotationPerm = []int{
	18, 19, 20, 21, 22, 23, 24, 25, 26,
	45, 46, 47, 48, 49, 50, 51, 52, 53,
	27, 28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 38, 39, 40, 41, 42, 43, 44,

	0, 1, 2, 5, 8, 7, 6, 3,
	11, 10, 9, 12, 15, 16, 17, 14,
}

// ZRotationPerm is like XRotationPerm, but for
// "z" rotations.
var ZRotationPerm = []int{
	51, 48, 45, 52, 49, 46, 53, 50, 47,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	38, 41, 44, 37, 40, 43, 36, 39, 42,
	17, 16, 15, 14, 13, 12, 11, 10, 9,

	18, 19, 20, 23, 26, 25, 24, 21,
	29, 28, 27, 30, 33, 34, 35, 32,
}

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

// Rotate applies a rotation to the stickers.
// It moves the centers as well, so it might be
// necessary to call s.ReinterpretCenters() before
// converting s to a CubieCube.
func (s *StickerCube) Rotate(r Rotation) {
	var applyCount int
	var perm []int
	switch r.Turns() {
	case 1:
		applyCount = 1
	case 2:
		applyCount = 2
	case -1:
		applyCount = 3
	}
	switch r.Axis() {
	case 0:
		perm = XRotationPerm
	case 1:
		perm = YRotationPerm
	case 2:
		perm = ZRotationPerm
	}

	for j := 0; j < applyCount; j++ {
		for i := 0; i < 9; i++ {
			indexes := [4]int{perm[i], perm[i+9], perm[i+18], perm[i+27]}
			s[indexes[0]], s[indexes[1]], s[indexes[2]], s[indexes[3]] =
				s[indexes[3]], s[indexes[0]], s[indexes[1]], s[indexes[2]]
		}
		for k := 0; k < 2; k++ {
			listIdx := 4*9 + k*8
			for l := 0; l < 2; l++ {
				temp := s[perm[listIdx+7]]
				for i := 7; i > 0; i-- {
					idx := perm[listIdx+i]
					lastIdx := perm[listIdx+i-1]
					s[idx] = s[lastIdx]
				}
				s[perm[listIdx]] = temp
			}
		}
	}
}

// ReinterpretCenters changes the meaning of different
// stickers so that the cube is oriented with 1 on top,
// 3 in front.
// In other words, it recolors the cube to be in the
// standard color scheme.
func (s *StickerCube) ReinterpretCenters() {
	colorMap := map[int]int{
		s[4]:    1,
		s[9+4]:  2,
		s[18+4]: 3,
		s[27+4]: 4,
		s[36+4]: 5,
		s[45+4]: 6,
	}
	for i, x := range s {
		s[i] = colorMap[x]
	}
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
