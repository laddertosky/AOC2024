package main

import (
	"bufio"
	"fmt"
	"os"
)

type Antinas map[byte][][]int

func parseAntenas(board []string) Antinas {

	m := len(board)
	n := len(board[0])

	antinas := Antinas{}

	for r := range m {
		for c := range n {
			if board[r][c] == '.' {
				continue
			}

			freq := board[r][c]
			if antinas[freq] == nil {
				antinas[freq] = [][]int{}
			}

			antinas[freq] = append(antinas[freq], []int{r, c})
		}
	}
	return antinas
}

func isValid(pair []int, m int, n int) bool {
	if len(pair) != 2 {
		panic("input size is not 2, got " + string(len(pair)))
	}
	return pair[0] >= 0 && pair[0] < m && pair[1] >= 0 && pair[1] < n
}

func countAntiNodeAdvanced(antinas Antinas, m int, n int) int {
	antinodes := map[int]bool{}
	for _, nodes := range antinas {
		if len(nodes) < 2 {
			continue
		}

		for i := range len(nodes) {
			for j := range i {
				diffX := nodes[i][0] - nodes[j][0]
				diffY := nodes[i][1] - nodes[j][1]

				k := 0
				node := []int{nodes[i][0] + k*diffX, nodes[i][1] + k*diffY}
				for isValid(node, m, n) {
					antinodes[node[0]*n+node[1]] = true
					k++
					node = []int{nodes[i][0] + k*diffX, nodes[i][1] + k*diffY}
				}

				k = 0
				node = []int{nodes[j][0] - k*diffX, nodes[j][1] - k*diffY}
				for isValid(node, m, n) {
					antinodes[node[0]*n+node[1]] = true
					k++
					node = []int{nodes[j][0] - k*diffX, nodes[j][1] - k*diffY}
				}
			}
		}
	}

	return len(antinodes)
}

func countAntiNode(antinas Antinas, m int, n int) int {
	antinodes := map[int]bool{}

	for _, nodes := range antinas {
		if len(nodes) < 2 {
			continue
		}

		for i := range len(nodes) {
			for j := range i {
				afterI := []int{nodes[i][0]*2 - nodes[j][0], nodes[i][1]*2 - nodes[j][1]}
				afterJ := []int{nodes[j][0]*2 - nodes[i][0], nodes[j][1]*2 - nodes[i][1]}

				if isValid(afterI, m, n) {
					antinodes[afterI[0]*n+afterI[1]] = true
				}

				if isValid(afterJ, m, n) {
					antinodes[afterJ[0]*n+afterJ[1]] = true
				}
			}
		}
	}

	return len(antinodes)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)

	sc := bufio.NewScanner(f)
	board := []string{}

	for sc.Scan() {
		line := sc.Text()
		board = append(board, line)
	}

	antenas := parseAntenas(board)
	partOneResult := countAntiNode(antenas, len(board), len(board[0]))
	partTwoResult := countAntiNodeAdvanced(antenas, len(board), len(board[0]))
	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)
}
