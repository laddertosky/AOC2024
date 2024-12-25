package codec

type MoveGroup map[string]map[string]int

type Tranisition map[string]map[string][]string

var MoveTransition Tranisition = Tranisition{
	Up: {
		Left:     {Down, Left, Activate},
		Up:       {Activate},
		Right:    {Down, Right, Activate},
		Down:     {Down, Activate},
		Activate: {Right, Activate},
	},
	Left: {
		Left:     {Activate},
		Up:       {Right, Up, Activate},
		Right:    {Right, Right, Activate},
		Down:     {Right, Activate},
		Activate: {Right, Right, Up, Activate},
	},
	Down: {
		Left:     {Left, Activate},
		Up:       {Up, Activate},
		Right:    {Right, Activate},
		Down:     {Activate},
		Activate: {Right, Up, Activate},
	},
	Right: {
		Left:     {Left, Left, Activate},
		Up:       {Left, Up, Activate},
		Right:    {Activate},
		Down:     {Left, Activate},
		Activate: {Up, Activate},
	},
	Activate: {
		Left:     {Down, Left, Left, Activate},
		Up:       {Left, Activate},
		Right:    {Down, Activate},
		Down:     {Down, Left, Activate},
		Activate: {Activate},
	},
}

func Decode(instruction string, indirect int) int {
	return DecodeNumeric(instruction, indirect)
}

func DecodeNumeric(instruction string, indirect int) int {
	var result int
	DecodeNumericOneStep(&result, instruction, 0, "", Activate, indirect)
	return result
}

func DecodeNumericOneStep(result *int, instruction string, index int, current string, pos string, indirect int) {
	if index == len(instruction) {
		candidate := DecodeArrow(current, indirect)
		if *result == 0 {
			*result = candidate
		} else if *result > candidate {
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

func initialPattern() MoveGroup {
	return MoveGroup{
		Up: {
			Left:     0,
			Up:       0,
			Right:    0,
			Down:     0,
			Activate: 0,
		},
		Left: {
			Left:     0,
			Up:       0,
			Right:    0,
			Down:     0,
			Activate: 0,
		},
		Down: {
			Left:     0,
			Up:       0,
			Right:    0,
			Down:     0,
			Activate: 0,
		},
		Right: {
			Left:     0,
			Up:       0,
			Right:    0,
			Down:     0,
			Activate: 0,
		},
		Activate: {
			Left:     0,
			Up:       0,
			Right:    0,
			Down:     0,
			Activate: 0,
		},
	}
}

func DecodeArrow(instruction string, indirect int) int {
	pattern := initialPattern()
	current := Activate
	for _, next := range instruction {
		pattern[current][string(next)]++
		current = string(next)
	}

	return DecodeArrowOneStep(pattern, indirect)
}

func DecodeArrowOneStep(pattern MoveGroup, indirect int) int {
	if indirect == 0 {
		result := 0
		for _, next := range pattern {
			for _, count := range next {
				result += count
			}
		}
		return result
	}

	nextPattern := initialPattern()
	for first, next := range pattern {
		for second, count := range next {
			prev := Activate
			for _, transition := range MoveTransition[first][second] {
				nextPattern[prev][transition] += count
				prev = transition
			}
		}
	}

	return DecodeArrowOneStep(nextPattern, indirect-1)
}

func move(direction string, count int) string {
	var result string

	for range count {
		result += direction
	}
	return result
}
