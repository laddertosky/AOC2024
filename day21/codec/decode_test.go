package codec_test

import (
	"main/codec"
	"testing"
)

func TestDecode(t *testing.T) {
	input := "v<<A>>^A<A>AvA<^AA>A<vAAA>^A"
	decoded := codec.DecodeArrow(input, 1)
	actual := codec.EncodeArrow(decoded)
	if actual != input {
		t.Fatalf("Obtained: %s, expected: %s, decoded: %s\n", actual, input, decoded)
	}
}

func TestDecodeNumeric(t *testing.T) {
	input := "029A"
	decoded := codec.DecodeNumeric(input, 0)
	actual := codec.EncodeNumeric(decoded)
	if actual != input {
		t.Fatalf("Obtained: %s, expected: %s, decoded: %s\n", actual, input, decoded)
	}
}

func TestDecodeCombination(t *testing.T) {
	input := "379A"
	decoded1 := codec.DecodeNumeric(input, 0)
	t.Logf("%s, %d\n", decoded1, len(decoded1))
	decoded2 := codec.DecodeArrow(decoded1, 1)
	t.Logf("%s, %d\n", decoded2, len(decoded2))
	decoded3 := codec.DecodeArrow(decoded2, 1)
	t.Logf("%s, %d\n", decoded3, len(decoded3))
}
