package main

import (
	"bufio"
	"fmt"
	"io/fs"
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

func div(src *int, dst *int, operand int) {
	*dst = *src / (1 << operand)
}

func bxl(B *int, operand int) {
	*B ^= operand
}

func bst(B *int, operand int) {
	*B = operand & 0x7
}

func bxc(B *int, C *int) {
	*B ^= *C
}

func combo(operand int, A *int, B *int, C *int) int {
	switch operand {
	case 4:
		operand = *A
	case 5:
		operand = *B
	case 6:
		operand = *C
	}
	return operand
}

func out(operand int) string {
	return strconv.Itoa(operand&0x7) + ","
}

func execute(program []int, pc int, A *int, B *int, C *int, raw string, earlyBreak bool) string {
	var output string
	for pc < len(program)-1 && !(earlyBreak && !strings.Contains(raw, output)) {
		opcode := program[pc]
		operand := program[pc+1]
		pc += 2

		switch opcode {
		case 0:
			div(A, A, combo(operand, A, B, C))
		case 1:
			bxl(B, operand)
		case 2:
			bst(B, combo(operand, A, B, C))
		case 3:
			if *A != 0 {
				pc = combo(operand, A, B, C)
			}
		case 4:
			bxc(B, C)
		case 5:
			output += out(combo(operand, A, B, C))
		case 6:
			div(A, B, combo(operand, A, B, C))
		case 7:
			div(A, C, combo(operand, A, B, C))
		default:
			panic("unknown opcode: " + string(program[pc]))
		}
	}

	if len(output) > 0 && output[len(output)-1] == ',' {
		output = output[:len(output)-1]
	}
	return output
}

func parseInput(sc *bufio.Scanner) (A int, B int, C int, buf string) {
	sc.Scan()
	fmt.Sscanf(sc.Text(), "Register A: %d", &A)
	sc.Scan()
	fmt.Sscanf(sc.Text(), "Register B: %d", &B)
	sc.Scan()
	fmt.Sscanf(sc.Text(), "Register C: %d", &C)
	sc.Scan()
	sc.Scan()

	buf = sc.Text()[len("Program: "):]
	return
}

// this works for input1 and input3
func findOutputEqualToProgram(B int, C int, program []int, raw string) int {
	var output string
	candidates := []int{}
	for i := range 8 {
		candidates = append(candidates, i)
	}
	a := -1
	visited := map[int]bool{}

	for len(candidates) > 0 && output != raw && len(output) <= len(raw) {
		a = candidates[0]
		candidates = candidates[1:]

		regA := a
		regB := B
		regC := C
		output = execute(program, 0, &regA, &regB, &regC, raw+",", true)

		if strings.Contains(raw+"$", output+"$") {
			for i := range 8 {
				next := 8*a + i
				if visited[next] {
					continue
				}
				visited[next] = true
				candidates = append(candidates, next)
			}
		}
	}

	regA := a
	regB := B
	regC := C
	output = execute(program, 0, &regA, &regB, &regC, raw+",", true)
	if output != raw {
		a = -1
	}
	return a
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	A, B, C, buf := parseInput(sc)

	instructions := strings.Split(buf, `,`)
	program := []int{}
	for _, instruction := range instructions {
		operand, err := strconv.Atoi(instruction)
		panicIf(err)
		program = append(program, operand)
	}

	regA := A
	regB := B
	regC := C
	partOneResult := execute(program, 0, &regA, &regB, &regC, buf, false)
	partTwoResult := findOutputEqualToProgram(B, C, program, buf)

	fmt.Printf("%s, partOne: %s, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)

	return nil
}

func main() {
	filepath.Walk(".", run)
}
