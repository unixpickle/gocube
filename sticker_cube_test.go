package gocube

import (
	"testing"
)

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
