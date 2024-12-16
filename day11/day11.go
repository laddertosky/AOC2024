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

func blink(stones []int, repeats int) int {
	for range repeats {
		next := []int{}
		for _, stone := range stones {
			str := strconv.Itoa(stone)
			if stone == 0 {
				next = append(next, 1)
			} else if len(str)%2 == 0 {
				left, err := strconv.Atoi(str[:len(str)/2])
				panicIf(err)
				right, err := strconv.Atoi(str[len(str)/2:])
				panicIf(err)

				next = append(next, left, right)
			} else {
				next = append(next, stone*2024)
			}
		}
		stones = next
	}

	return len(stones)
}

type Cache map[int]map[int]int

var cache Cache = Cache{}

func dfs(stone int, level int) int {
	if level == 0 {
		return 1
	}

	if cache[stone] != nil && cache[stone][level] != 0 {
		return cache[stone][level]
	}

	var result int
	if stone == 0 {
		result = dfs(1, level-1)
	} else {
		str := strconv.Itoa(stone)
		if len(str)%2 == 0 {
			left, err := strconv.Atoi(str[:len(str)/2])
			panicIf(err)
			right, err := strconv.Atoi(str[len(str)/2:])
			panicIf(err)

			result = dfs(left, level-1) + dfs(right, level-1)
		} else {
			result = dfs(stone*2024, level-1)
		}
	}

	if cache[stone] == nil {
		cache[stone] = make(map[int]int)
	}
	cache[stone][level] = result
	return result
}

func optimizedBlink(stones []int, repeats int) int {
	result := 0
	for _, stone := range stones {
		tmp := dfs(stone, repeats)
		result += tmp
	}
	return result
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	sc.Scan()
	buff := sc.Text()
	stonesStr := strings.Split(buff, " ")

	stones := []int{}
	for i := range stonesStr {
		stone, err := strconv.Atoi(stonesStr[i])
		panicIf(err)
		stones = append(stones, stone)
	}

	partOneResult := blink(stones, 25)
	partTwoResult := optimizedBlink(stones, 75)
	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)

}
