package codec

const Up string = "^"
const Left string = "<"
const Down string = "v"
const Right string = ">"
const Activate string = "A"
const Empty string = "E"

var Movement map[string][2]int = map[string][2]int{
	Up:       {-1, 0},
	Activate: {0, 0},
	Left:     {0, -1},
	Down:     {1, 0},
	Right:    {0, 1},
}

type Pad map[string][2]int
type PadPos [][]string

// Definition of Arrow Pad: start from A
// ....+---+---+
// ....| ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+
var ArrowPadPos PadPos = PadPos{
	{Empty, Up, Activate},
	{Left, Down, Right},
}

var ArrowPad Pad = Pad{
	Up:       {0, 1},
	Activate: {0, 2},
	Left:     {1, 0},
	Down:     {1, 1},
	Right:    {1, 2},
}

// Definition of Numeric Pad: start from A
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
// ....| 0 | A |
// ....+---+---+
var NumericPadPos PadPos = PadPos{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{Empty, "0", Activate},
}

var NumericPad Pad = Pad{
	"0": [2]int{3, 1},
	"A": [2]int{3, 2},
	"1": [2]int{2, 0},
	"2": [2]int{2, 1},
	"3": [2]int{2, 2},
	"4": [2]int{1, 0},
	"5": [2]int{1, 1},
	"6": [2]int{1, 2},
	"7": [2]int{0, 0},
	"8": [2]int{0, 1},
	"9": [2]int{0, 2},
}
