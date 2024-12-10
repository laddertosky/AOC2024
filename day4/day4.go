package main

import (
	"bufio"
	"fmt"
	"os"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func parseGrid(sc *bufio.Scanner) []string {
	rows := []string{}
	for sc.Scan() {
		buff := sc.Text()
		if len(buff) == 0 {
			continue
		}
		rows = append(rows, sc.Text())
	}
	return rows
}

var TARGET string = "XMAS"

type DIRECTION struct {
	dx int
	dy int
}

var UL DIRECTION = DIRECTION{dx: -1, dy: -1}
var U DIRECTION = DIRECTION{dx: 0, dy: -1}
var UR DIRECTION = DIRECTION{dx: 1, dy: -1}
var R DIRECTION = DIRECTION{dx: 1, dy: 0}
var DR DIRECTION = DIRECTION{dx: 1, dy: 1}
var D DIRECTION = DIRECTION{dx: 0, dy: 1}
var DL DIRECTION = DIRECTION{dx: -1, dy: 1}
var L DIRECTION = DIRECTION{dx: -1, dy: 0}

var ALLDIRECTIONS []DIRECTION = []DIRECTION{UL, U, UR, R, DR, D, DL, L}

func dfs(rows []string, r int, c int, index int, direction DIRECTION) bool {
	if index == len(TARGET) {
		return true
	}

	if r < 0 || c < 0 || r >= len(rows) || c >= len(rows[0]) {
		return false
	}

	if rows[r][c] == TARGET[index] {
		return dfs(rows, r+direction.dx, c+direction.dy, index+1, direction)
	}
	return false
}

func findXMAS(rows []string) int {
	m := len(rows)
	n := len(rows[0])

	result := 0
	for r := range m {
		for c := range n {
			if rows[r][c] != TARGET[0] {
				continue
			}

			for _, direction := range ALLDIRECTIONS {
				if dfs(rows, r, c, 0, direction) {
					result++
				}
			}
		}
	}

	return result
}

func findX_MAS(rows []string) int {
	m := len(rows)
	n := len(rows[0])

	result := 0
	for r := 1; r < m-1; r++ {
		for c := 1; c < n-1; c++ {
			if rows[r][c] != 'A' {
				continue
			}
			corners := string(rows[r-1][c-1]) + string(rows[r-1][c+1]) + string(rows[r+1][c+1]) + string(rows[r+1][c-1])
			if corners == "MMSS" || corners == "SMMS" || corners == "SSMM" || corners == "MSSM" {
				result++
			}

		}
	}
	return result
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)

	sc := bufio.NewScanner(f)
	rows := parseGrid(sc)
	partOneResult := findXMAS(rows)
	partTwoResult := findX_MAS(rows)
	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)
}
