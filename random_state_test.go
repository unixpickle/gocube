package gocube

import "testing"

func TestRandomLLCorners(t *testing.T) {
	seen := map[CubieCorners]bool{}

	for i := 0; i < 100000; i++ {
		corns, _ := RandomLLCorners()
		seen[corns] = true
	}

	expectedCount := factorial(4) * 3 * 3 * 3
	if len(seen) != expectedCount {
		t.Errorf("expected %d cases but got %d", expectedCount, len(seen))
	}
}

func TestRandomZBLL(t *testing.T) {
	seen := map[CubieCube]bool{}

	for i := 0; i < 500000; i++ {
		seen[RandomZBLL()] = true
	}

	expectedCount := (factorial(4) * factorial(4) / 2) * 3 * 3 * 3
	if len(seen) != expectedCount {
		t.Errorf("expected %d cases but got %d", expectedCount, len(seen))
	}
}
