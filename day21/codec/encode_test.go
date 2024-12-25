package codec_test

import (
	"main/codec"
	"testing"
)

func TestEncodeUp(t *testing.T) {
	input := "<A"
	expected := "^"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeLeft(t *testing.T) {
	input := "<v<A"
	expected := "<"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeLeft2(t *testing.T) {
	input := "<<vA"
	expected := "<"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeLeft3(t *testing.T) {
	input := "v<<A"
	expected := "<"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeRight(t *testing.T) {
	input := "vA"
	expected := ">"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeDown(t *testing.T) {
	input := "<vA"
	expected := "v"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeDown2(t *testing.T) {
	input := "v<A"
	expected := "v"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeActivate(t *testing.T) {
	input := "A"
	expected := "A"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeMultiple(t *testing.T) {
	input := "<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"
	expected := "v<<A>>^A<A>AvA<^AA>A<vAAA>^A"
	actual := codec.EncodeArrow(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeNumeric(t *testing.T) {
	input := "<A^A>^^AvvvA"
	expected := "029A"
	actual := codec.EncodeNumeric(input)
	if expected != actual {
		t.Fatalf("Obtained: %s, expected: %s\n", actual, expected)
	}
}

func TestEncodeCombination(t *testing.T) {
	input := "<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"
	expected1 := "v<<A>>^A<A>AvA<^AA>A<vAAA>^A"
	actual1 := codec.EncodeArrow(input)
	if actual1 != expected1 {
		t.Fatalf("Obtained: %s, expected: %s\n", actual1, expected1)
	}

	expected2 := "<A^A>^^AvvvA"
	actual2 := codec.EncodeArrow(actual1)
	if actual2 != expected2 {
		t.Fatalf("Obtained: %s, expected: %s\n", actual2, expected2)
	}

	expected3 := "029A"
	actual3 := codec.EncodeNumeric(actual2)
	if actual3 != expected3 {
		t.Fatalf("Obtained: %s, expected: %s\n", actual3, expected3)
	}
}

func TestEncodeDo(t *testing.T) {
	input := "<v<A>>^AvA^A<vA<AA>>^AAvA<^A>AAvA^A<vA>^AA<A>A<v<A>A>^AAAvA<^A>A"
	t.Logf("%s, %d\n", input, len(input))
	actual1 := codec.EncodeArrow(input)
	t.Logf("%s, %d\n", actual1, len(actual1))

	actual2 := codec.EncodeArrow(actual1)
	t.Logf("%s, %d\n", actual2, len(actual2))

	expected3 := "379A"
	actual3 := codec.EncodeNumeric(actual2)
	if actual3 != expected3 {
		t.Fatalf("Obtained: %s, expected: %s\n", actual3, expected3)
	}
}
