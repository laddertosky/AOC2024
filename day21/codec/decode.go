package codec

import "fmt"

var NextArrowMove map[string]map[string][]string = map[string]map[string][]string{
	Up: {
		Left:     {"v<A"},
		Up:       {"A"},
		Right:    {"v>A", ">vA"},
		Down:     {"vA"},
		Activate: {">A"},
	},
	Left: {
		Left:     {"A"},
		Up:       {">^A"},
		Right:    {">>A"},
		Down:     {">A"},
		Activate: {">>^A", ">^>A"},
	},
	Down: {
		Left:     {"<A"},
		Up:       {"^A"},
		Right:    {">A"},
		Down:     {"A"},
		Activate: {">^A", "^>A"},
	},
	Right: {
		Left:     {"<<A"},
		Up:       {"<^A", "^<A"},
		Right:    {"A"},
		Down:     {"<A"},
		Activate: {"^A"},
	},
	Activate: {
		Left:     {"v<<A", "<v<A"},
		Up:       {"<A"},
		Right:    {"vA"},
		Down:     {"<vA", "v<A"},
		Activate: {"A"},
	},
}

func Decode(instruction string, indirect int) string {
	return DecodeNumeric(instruction, indirect)
}

func DecodeNumeric(instruction string, indirect int) string {
	var result string
	DecodeNumericOneStep(&result, instruction, 0, "", Activate, indirect)
	fmt.Printf("input: %s, decoded numeric: %s\n", instruction, result)
	return result
}

func DecodeNumericOneStep(result *string, instruction string, index int, current string, pos string, indirect int) {
	if index == len(instruction) {
		candidate := DecodeArrow(current, indirect)
		if len(*result) == 0 {
			*result = candidate
		} else if len(*result) > len(candidate) {
			fmt.Printf("Replacing %s (len: %d) with %s (len: %d)\n", *result, len(*result), candidate, len(candidate))
			*result = candidate
		}
		return
	}

	now := string(instruction[index])
	dr := NumericPad[now][0] - NumericPad[pos][0]
	dc := NumericPad[now][1] - NumericPad[pos][1]
	// up-down first
	if NumericPadPos[NumericPad[pos][0]+dr][NumericPad[pos][1]] != Empty {
		next := current
		if dr > 0 {
			next += move(Down, dr)
		} else {
			next += move(Up, -dr)
		}

		if dc > 0 {
			next += move(Right, dc)
		} else {
			next += move(Left, -dc)
		}
		next += Activate
		DecodeNumericOneStep(result, instruction, index+1, next, now, indirect)
	}

	// left-right first
	if NumericPadPos[NumericPad[pos][0]][NumericPad[pos][1]+dc] != Empty {
		next := current
		if dc > 0 {
			next += move(Right, dc)
		} else {
			next += move(Left, -dc)
		}

		if dr > 0 {
			next += move(Down, dr)
		} else {
			next += move(Up, -dr)
		}
		next += Activate
		DecodeNumericOneStep(result, instruction, index+1, next, now, indirect)
	}
}

func DecodeArrowOneStep(result *string, instruction string, index int, current string, pos string) {
	if index == len(instruction) {
		if len(*result) == 0 {
			*result = current
		} else if len(*result) > len(current) {
			fmt.Printf("Replacing %s (len: %d) with %s (len: %d)\n", *result, len(*result), current, len(current))
			*result = current
		}
		return
	}

	now := string(instruction[index])
	for _, next := range NextArrowMove[pos][now] {
		DecodeArrowOneStep(result, instruction, index+1, current+next, now)
	}
}

func DecodeArrow(instruction string, indirect int) string {
	if indirect == 0 {
		return instruction
	}

	var result string
	DecodeArrowOneStep(&result, instruction, 0, "", Activate)

	fmt.Printf("instruction: %s\n", result)
	return DecodeArrow(result, indirect-1)
}

func move(direction string, count int) string {
	var result string

	for range count {
		result += direction
	}
	return result
}

// Priority: Right > Top > Down > Left
func move2(dr int, dc int) string {
	var result string

	if dc < 0 {
		result += move(Left, -dc)
	}

	if dr < 0 {
		result += move(Up, -dr)
	}

	if dr > 0 {
		result += move(Down, dr)
	}

	if dc > 0 {
		result += move(Right, dc)
	}

	result += Activate
	return result
}
