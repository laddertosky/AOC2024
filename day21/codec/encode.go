package codec

func EncodeArrow(input string) string {
	return encodeImpl(input, ArrowPad, ArrowPadPos)
}

func EncodeNumeric(input string) string {
	return encodeImpl(input, NumericPad, NumericPadPos)
}

func encodeImpl(input string, pad Pad, padPos PadPos) string {
	currentPos := pad[Activate]

	var result string
	for _, move := range input {
		currentPos[0] += Movement[string(move)][0]
		currentPos[1] += Movement[string(move)][1]
		current := padPos[currentPos[0]][currentPos[1]]
		if string(move) == Activate {
			result += current
		}
	}

	return result
}
