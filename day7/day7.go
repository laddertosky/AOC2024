package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func backtrack(target int, current int, index int, options []int) bool {
	if index == len(options) {
		return target == current
	}

	add := backtrack(target, current+options[index], index+1, options)
	mul := backtrack(target, current*options[index], index+1, options)
	base := 1
	for base <= options[index] {
		base *= 10
	}
	concat := backtrack(target, current*base+options[index], index+1, options)

	return add || mul || concat
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)

	sc := bufio.NewScanner(f)
	result := 0

	for sc.Scan() {
		buff := sc.Text()

		line := strings.Split(buff, ":")
		target, err := strconv.Atoi(line[0])
		panicIf(err)
		options := []int{}
		numbers := strings.Split(line[1], " ")

		for _, numStr := range numbers {
			if len(numStr) == 0 {
				continue
			}

			num, err := strconv.Atoi(numStr)
			panicIf(err)
			options = append(options, num)
		}

		success := backtrack(target, options[0], 1, options)
		if success {
			result += target
		}
	}

	fmt.Println(result)
}
