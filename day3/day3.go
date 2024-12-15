package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func parseMul(buff string) (int, int, int) {
	pos := strings.Index(buff, "mul(")

	if pos == -1 {
		return math.MaxInt, 0, 0
	}

	pos += 4
	if pos >= len(buff) {
		return math.MaxInt, 0, 0
	}

	n1 := 0
	i := 0
	for ; i < 3 && buff[pos+i] != ','; i++ {
		d, err := strconv.Atoi(string(buff[pos+i]))
		if err != nil {
			break
		}
		n1 *= 10
		n1 += d
	}

	if buff[pos+i] != ',' {
		return pos + i, 0, 0
	}

	n2 := 0
	pos += i + 1
	i = 0
	for ; i < 3 && buff[pos+i] != ')'; i++ {
		d, err := strconv.Atoi(string(buff[pos+i]))
		if err != nil {
			break
		}
		n2 *= 10
		n2 += d
	}

	if buff[pos+i] != ')' {
		return pos + i, 0, 0
	}
	pos += i + 1

	return pos, n1, n2
}

func parseDo(buff string) int {
	pos := strings.Index(buff, "do()")

	if pos == -1 {
		return math.MaxInt
	}
	return pos + len("do()")
}

func parseDont(buff string) int {
	pos := strings.Index(buff, "don't()")

	if pos == -1 {
		return math.MaxInt
	}

	return pos + len("don't()")
}

func main() {
	f, err := os.OpenFile("./input1.txt", os.O_RDONLY, 0660)
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	partOneResult := 0
	partTwoResult := 0

	enabled := true
	for sc.Scan() {
		buff := sc.Text()

		pos := 0
		for {
			buff = buff[pos:]

			next := math.MaxInt
			nextMul, n1, n2 := parseMul(buff)
			if next > nextMul {
				next = nextMul
			}
			nextDo := parseDo(buff)
			if next > nextDo {
				next = nextDo
			}
			nextDont := parseDont(buff)
			if next > nextDont {
				next = nextDont
			}

			switch next {
			case nextMul:
				if enabled {
					partTwoResult += n1 * n2
				}
				partOneResult += n1 * n2
			case nextDo:
				enabled = true
			case nextDont:
				enabled = false
			}
			pos = next

			if pos == math.MaxInt {
				break
			}
		}
	}
	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)
}
