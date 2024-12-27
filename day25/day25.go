package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	keys := [][]int{}
	locks := [][]int{}

	isLock := false
	isKey := false
	var block []int
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 {
			if isKey {
				keys = append(keys, block)
				isKey = false
			} else {
				locks = append(locks, block)
				isLock = false
			}
		} else {
			if !isLock && !isKey {
				block = make([]int, len(line))
				isLock = line[0] == '#'
				isKey = line[0] == '.'
			}

			for i, c := range line {
				if c == '#' && isLock {
					block[i]++
				} else if c == '.' && isKey {
					block[i]++
				}
			}
		}
	}

	partOneResult := findMatches(keys, locks)

	fmt.Printf("%s, partOne: %v\n", info.Name(), partOneResult)
	return nil
}

func findMatches(keys [][]int, locks [][]int) int {
	result := 0
	for _, lock := range locks {
		for _, key := range keys {
			allMatch := true
			for i := range len(lock) {
				if lock[i] > key[i] {
					allMatch = false
				}
			}
			if allMatch {
				result++
			}
		}
	}
	return result
}

func main() {
	filepath.Walk(".", run)
}
