package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"main/codec"
	"os"
	"path/filepath"
	"strconv"
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
	instructions := []string{}

	for sc.Scan() {
		instruction := sc.Text()
		instructions = append(instructions, instruction)
	}

	partOneResult := calculateComplexity(instructions, 2)
	partTwoResult := calculateComplexity(instructions, 25)
	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)

	return nil
}

func calculateComplexity(instructions []string, indirect int) int {
	result := 0

	for _, instruction := range instructions {
		keystrokes := codec.Decode(instruction, indirect)
		num, err := strconv.Atoi(instruction[0 : len(instructions)-2])
		panicIf(err)

		fmt.Printf("input: %s, length: %d, num: %d\n", instruction, keystrokes, num)
		result += keystrokes * num
	}

	return result
}

func main() {
	filepath.Walk(".", run)
}
