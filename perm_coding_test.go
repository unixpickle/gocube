package gocube

import "testing"

func TestEncodeChoose(t *testing.T) {
	answer := 0
	for i := 0; i < 12; i++ {
		for j := i+1; j < 12; j++ {
			for k := j+1; k < 12; k++ {
				for l := k+1; l < 12; l++ {
					perm := make([]bool, 12)
					perm[i] = true
					perm[j] = true
					perm[k] = true
					perm[l] = true
					if answer != encodeChoice(perm) {
						t.Error("Expected", answer, "got",
							encodeChoice(perm))
					}
					answer++
				}
			}
		}
	}
}