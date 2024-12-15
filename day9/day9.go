package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func checkSum(input string) int {
	l := 0
	r := len(input) - 1
	if r%2 == 1 {
		r--
	}

	toMoved, err := strconv.Atoi(string(input[r]))
	panicIf(err)

	result := 0
	id := 0

	for l < r {
		files, err := strconv.Atoi(string(input[l]))
		panicIf(err)
		result += (l / 2) * (id + id + files - 1) * files / 2
		id += files

		slots, err := strconv.Atoi(string(input[l+1]))
		panicIf(err)
		for toMoved <= slots {
			result += (r / 2) * (id + id + toMoved - 1) * toMoved / 2
			id += toMoved
			r -= 2
			slots -= toMoved
			tmp, err := strconv.Atoi(string(input[r]))
			panicIf(err)
			toMoved = tmp
		}

		if toMoved > slots {
			result += (r / 2) * (id + id + slots - 1) * slots / 2
			id += slots
			toMoved -= slots
		}
		l += 2
	}
	result += (r / 2) * (id + id + toMoved - 1) * toMoved / 2

	return result
}

func checkSum2(input string) int {

	initial := map[int][]int{}
	spaces := map[int][]int{}
	pos := 0
	for i, d := range input {
		num, err := strconv.Atoi(string(d))
		panicIf(err)

		if i%2 == 0 {
			initial[i/2] = []int{pos, num}
		} else {
			spaces[i] = []int{num, 0}
		}

		pos += num
	}

	result := 0

	for r := len(input) / 2; r >= 0; r-- {
		toMoved := initial[r][1]

		l := 1
		for ; l < 2*r && toMoved > spaces[l][0]; l += 2 {
		}

		var pos int
		if l < 2*r {
			pos = initial[l/2][0] + initial[l/2][1] + spaces[l][1]

			spaces[l][0] -= toMoved
			spaces[l][1] += toMoved
		} else {
			pos = initial[r][0]
		}

		tmp := r * (pos + pos + toMoved - 1) * toMoved / 2

		result += tmp
	}

	return result
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	if sc.Scan() {
		input := sc.Text()
		partOneResult := checkSum(input)
		partTwoResult := checkSum2(input)
		fmt.Println(partOneResult)
		fmt.Println(partTwoResult)
	}

}
