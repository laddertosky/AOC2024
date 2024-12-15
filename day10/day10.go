package main

import (
	"bufio"
	"fmt"
	"os"
)

var DIRECTIONS = [][]int{
	{0, 1}, {0, -1}, {1, 0}, {-1, 0},
}

func dfs(visited9 map[int]bool, paths *int, heights [][]int, m int, n int, pr int, pc int, r int, c int, current int) {
	if current == 9 {
		visited9[r*n+c] = true
		*paths++
		return
	}

	for _, dir := range DIRECTIONS {
		nr := r + dir[0]
		nc := c + dir[1]
		if nr == pr && nc == pc {
			continue
		}

		if nr < 0 || nc < 0 || nr >= m || nc >= n {
			continue
		}

		if heights[nr][nc] != current+1 {
			continue
		}

		dfs(visited9, paths, heights, m, n, r, c, nr, nc, current+1)
	}
}

func findTrailHead(heights [][]int) (int, int) {
	nodes := 0
	paths := 0
	m := len(heights)
	n := len(heights[0])

	for r := range m {
		for c := range n {
			if heights[r][c] != 0 {
				continue
			}

			visited9 := map[int]bool{}
			path := 0

			dfs(visited9, &path, heights, m, n, -1, -1, r, c, 0)
			nodes += len(visited9)
			paths += path
		}
	}
	return nodes, paths
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	heights := [][]int{}

	for sc.Scan() {
		buff := sc.Text()
		row := []int{}
		for i := range buff {
			row = append(row, int(buff[i]-'0'))
		}

		heights = append(heights, row)
	}

	partOneResult, partTwoResult := findTrailHead(heights)
	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)

}
