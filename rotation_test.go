package gocube

import "testing"

func TestRotationNotation(t *testing.T) {
	rotStrings := []string{"x", "y", "z", "x'", "y'", "z'", "x2", "y2", "z2"}

	for i, s := range rotStrings {
		rot, err := ParseRotation(s)
		if err != nil {
			t.Error(err)
			continue
		}
		if rot != Rotation(i) {
			t.Error("invalid rotation for string", s)
			continue
		}
		if rot.String() != s {
			t.Errorf("got %s for %s.String()", rot.String(), s)
		}
	}

	invalStrings := []string{"x1", "a", "a'", "r2", "R", "R'", "R2"}
	for _, s := range invalStrings {
		_, err := ParseRotation(s)
		if err == nil {
			t.Error("expected error for", s)
		}
	}
}
