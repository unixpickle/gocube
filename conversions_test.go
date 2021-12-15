package gocube

import (
	"testing"
)

func BenchmarkCubieToSticker(b *testing.B) {
	// Run algorithm to generate a cubie cube
	cubie := SolvedCubieCube()
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		cubie.Move(move)
	}

	// Convert it to a sticker cube N times.
	for i := 0; i < b.N; i++ {
		cubie.StickerCube()
	}
}

func BenchmarkStickerToCubie(b *testing.B) {
	stickers, _ := ParseStickerCube("OGBYWWOOY OWOGYGGBR WBBGGOBRB " +
		"RWYYBRWYR RBWWROWYG GRGBORYOY")

	// Convert it to a sticker cube N times.
	for i := 0; i < b.N; i++ {
		stickers.CubieCube()
	}
}

func TestCubieToSticker(t *testing.T) {
	// Run algorithm for comparison
	cubie := SolvedCubieCube()
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		cubie.Move(move)
	}

	stickers := cubie.StickerCube()
	str := "OGBYWWOOY OWOGYGGBR WBBGGOBRB RWYYBRWYR RBWWROWYG GRGBORYOY"
	if stickers.String() != str {
		t.Error("Invalid stickers:", stickers.String())
	}
}

func TestStickerToCubieIdentity(t *testing.T) {
	stickers, err := ParseStickerCube("111111111 222222222 333333333 " +
		"444444444 555555555 666666666")
	if err != nil {
		t.Error(err)
		return
	}
	cubies, err := stickers.CubieCube()
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 8; i++ {
		if cubies.Corners[i].Piece != i || cubies.Corners[i].Orientation != 1 {
			t.Error("Invalid corner at index", i)
		}
	}

	for i := 0; i < 12; i++ {
		if cubies.Edges[i].Piece != i || cubies.Edges[i].Flip {
			t.Error("Invalid edge at index", i)
		}
	}
}

func TestStickerToCube(t *testing.T) {
	// I did the algorithm B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'
	stickers, err := ParseStickerCube("OGBYWWOOY OWOGYGGBR WBBGGOBRB " +
		"RWYYBRWYR RBWWROWYG GRGBORYOY")
	if err != nil {
		t.Error(err)
		return
	}
	cubies, err := stickers.CubieCube()
	if err != nil {
		t.Error(err)
		return
	}

	// Run algorithm for comparison
	answer := SolvedCubieCube()
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		answer.Move(move)
	}

	// Make sure the cubes are equal
	for i, x := range answer.Corners {
		c := cubies.Corners[i]
		if x.Piece != c.Piece || x.Orientation != c.Orientation {
			t.Error("Invalid corner at index", i)
		}
	}
	for i, x := range answer.Edges {
		e := cubies.Edges[i]
		if x.Piece != e.Piece || x.Flip != e.Flip {
			t.Error("Invalid edge at index", i)
		}
	}
}

func TestStickerCubieCycleConsistency(t *testing.T) {
	for i := 0; i < 10; i++ {
		cubieCube := RandomCubieCube()
		stickerCube := cubieCube.StickerCube()
		stickerStr := stickerCube.String()
		stickerCube1, err := ParseStickerCube(stickerStr)
		if err != nil {
			t.Fatal(err)
		}
		if *stickerCube1 != stickerCube {
			t.Fatalf("sticker cube parse failed: %v", stickerCube.String())
		}
		cubieCube1, err := stickerCube1.CubieCube()
		if err != nil {
			t.Fatal(err)
		}
		if *cubieCube1 != cubieCube {
			t.Fatalf("cubie cube parse failed: %v", stickerCube.String())
		}
	}
}
