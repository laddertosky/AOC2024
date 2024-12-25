package codec_test

import (
	"main/codec"
	"testing"
)

func TestDecode(t *testing.T) {
	input := "v<<A>>^A<A>AvA<^AA>A<vAAA>^A"
	expected := len("<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A")
	actual := codec.DecodeArrow(input, 1)
	if actual != expected {
		t.Fatalf("Obtained: %d, expected: %d\n", actual, expected)
	}
}

func TestDecodeNumeric(t *testing.T) {
	input := "029A"
	expected := len("<A^A>^^AvvvA")
	actual := codec.DecodeNumeric(input, 0)
	if actual != expected {
		t.Fatalf("Obtained: %d, expected: %d\n", actual, expected)
	}
}
