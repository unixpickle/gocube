package gocube

import "testing"

func TestParseStickerCube(t *testing.T) {
	input := "Y2OYWWYW1 656GY3GO4 BGGBGBGBG BRBOBWWWR RR1ORYYOR RGOBOYYRW"
	output := "YYOYWWYWW OROGYGGOB BGGBGBGBG BRBOBWWWR RRWORYYOR RGOBOYYRW"
	parsed, err := ParseStickerCube(input)
	if err != nil {
		t.Error(err)
		return
	}
	if parsed.String() != output {
		t.Error("Invalid output:", parsed.String())
	}
}

func TestStickerCubeRotate(t *testing.T) {
	scramble := "D F L2 D F L' R' U' F D' R' D2 F2 D' B2 R2 D2 L2 U L2 D'"
	moves, _ := ParseMoves(scramble)
	cube := SolvedCubieCube()
	for _, m := range moves {
		cube.Move(m)
	}
	stickers := cube.StickerCube()

	rotations := []string{"x", "y", "z'"}
	outcomes := []string{
		"ROOOGYGYY WYGYBGYGW YRRBYRROW BGYWWWORB GBWWRORWG BBOBOROGB",
		"YWBGWRBWO RRWRYOYBR WOGBRWGWR OBBGOBBRO WGYGBYGYW ROOOGYGYY",
		"GWRORWWBG BBOBOROGB OYYOGYROG GGWYBGWYY RRWRYOYBR OWBRWGBWY",
	}
	for i, s := range rotations {
		rot, _ := ParseRotation(s)
		stickers.Rotate(rot)
		if stickers.String() != outcomes[i] {
			t.Errorf("invalid outcome for %s expected '%s' but got '%s'", s, outcomes[i],
				stickers.String())
		}
		stickers.Rotate(rot.Inverse())
	}
}

func TestStickerReinterpretCenters(t *testing.T) {
	scramble := "D F L2 D F L' R' U' F D' R' D2 F2 D' B2 R2 D2 L2 U L2 D'"
	moves, _ := ParseMoves(scramble)
	cube := SolvedCubieCube()
	for _, m := range moves {
		cube.Move(m)
	}
	stickers := cube.StickerCube()
	stickers.Rotate(NewRotation(2, -1))
	stickers.Rotate(NewRotation(0, -1))

	stickers.ReinterpretCenters()

	outcome := "RROYWROYY BRRBYRGBY YOGBGOOWY WYBGBWBWW OBGGRWGGR WGBOOORYW"
	if outcome != stickers.String() {
		t.Error("invalid outcome:", stickers.String())
	}
}
