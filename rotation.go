package gocube

import "errors"

// A Rotation is a way of rotating the entire
// cube around the x, y, or z axes.
// There are 9 total rotations, 0 through 8.
// The first three rotation values are x, y, and z.
// The next three are x', y', and z'.
// The final three are x2, y2, and z2.
type Rotation int

// NewRotation creates a Rotation around a given
// axis (0=x, 1=y, 2=z), for a given number of
// turns, where "x" has 1 turn, "x'" has -1 turn,
// and "x2" has 2 turns.
func NewRotation(axis, turns int) Rotation {
	if turns == 1 {
		return Rotation(axis)
	} else if turns == -1 {
		return Rotation(axis + 3)
	} else if turns == 2 {
		return Rotation(axis + 6)
	} else {
		panic("unsupported turns value")
	}
}

// Axis returns a number 0, 1, or 2, indicating
// the x, y, or z axis respectively.
func (r Rotation) Axis() int {
	return int(r) % 3
}

// Turns returns the number of "turns".
// This is 1 for regular rotations, -1 for inverse
// rotations, and 2 for double rotations.
func (r Rotation) Turns() int {
	return [3]int{1, -1, 2}[int(r)/3]
}

// String returns the string representation
// of this rotation, in WCA notation.
func (r Rotation) String() string {
	axisStr := []string{"x", "y", "z"}[r.Axis()]
	turnsStr := map[int]string{-1: "'", 1: "", 2: "2"}[r.Turns()]
	return axisStr + turnsStr
}

// Inverse returns the inverse of this rotation.
func (r Rotation) Inverse() Rotation {
	if r.Turns() == 2 {
		return r
	}
	return NewRotation(r.Axis(), -r.Turns())
}

// ParseRotation parses a WCA rotation string.
func ParseRotation(s string) (Rotation, error) {
	if len(s) == 1 {
		axis, ok := map[string]int{"x": 0, "y": 1, "z": 2}[s]
		if !ok {
			return 0, errors.New("unknown axis: " + s)
		}
		return Rotation(axis), nil
	} else if len(s) == 2 {
		rot, err := ParseRotation(s[:1])
		if err != nil {
			return 0, err
		}
		switch s[1] {
		case '\'':
			return NewRotation(rot.Axis(), -1), nil
		case '2':
			return NewRotation(rot.Axis(), 2), nil
		default:
			return 0, errors.New("invalid rotation: " + s)
		}
	} else {
		return 0, errors.New("invalid rotation: " + s)
	}
}
