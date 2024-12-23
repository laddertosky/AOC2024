package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const Unreachable int = 1e9
const Wall rune = '#'

var DIRECTIONS [][]int = [][]int{
	{0, 1}, {0, -1}, {1, 0}, {-1, 0},
}

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

	sc := bufio.NewScanner(f)

	board := [][]rune{}
	start := [2]int{}
	end := [2]int{}
	for sc.Scan() {
		buff := sc.Text()

		row := make([]rune, len(buff))
		for c, ch := range buff {
			row[c] = ch
			if ch == 'S' {
				start[0] = len(board)
				start[1] = c
			} else if ch == 'E' {
				end[0] = len(board)
				end[1] = c
			}
		}
		board = append(board, row)
	}

	saveAtLeast := 100
	dist := traverse(board, start, end)
	partOneResult := findCheatSaving(dist, saveAtLeast)

	fmt.Printf("%s, partOne: %d (>= %d)\n", info.Name(), partOneResult, saveAtLeast)

	return nil
}

func traverse(board [][]rune, start [2]int, end [2]int) [][]int {
	m := len(board)
	n := len(board[0])
	dist := make([][]int, m)
	for r := range m {
		dist[r] = make([]int, n)
		for c := range n {
			dist[r][c] = Unreachable
		}
	}

	queue := [][]int{{0, start[0], start[1]}}
	dist[start[0]][start[1]] = 0

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		d := front[0]
		r := front[1]
		c := front[2]

		if d > dist[r][c] {
			continue
		}

		if r == end[0] && c == end[1] {
			continue
		}

		for _, dir := range DIRECTIONS {
			nr := r + dir[0]
			nc := c + dir[1]
			if nr < 0 || nc < 0 || nr >= m || nc >= n {
				continue
			}

			if board[nr][nc] == Wall {
				continue
			}

			if dist[r][c]+1 >= dist[nr][nc] {
				continue
			}

			dist[nr][nc] = dist[r][c] + 1
			queue = append(queue, []int{dist[nr][nc], nr, nc})
		}

	}

	return dist
}

func printRecord(record map[int]int) {
	keys := []int{}
	for k := range record {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, save := range keys {
		fmt.Printf("save: %d, count: %d\n", save, record[save])
	}
}

func findCheatSaving(dist [][]int, saveAtLeast int) int {
	result := 0

	visited := map[int]map[int]int{}
	record := map[int]int{}
	m := len(dist)
	n := len(dist[0])

	for r := range m {
		for c := range n {
			if dist[r][c] == Unreachable {
				continue
			}

			for _, stepOne := range DIRECTIONS {
				r1 := r + stepOne[0]
				c1 := c + stepOne[1]
				if r1 < 0 || c1 < 0 || r1 >= m || c1 >= n {
					continue
				}

				h1 := r1*n + c1

				for _, stepTwo := range DIRECTIONS {
					r2 := r1 + stepTwo[0]
					c2 := c1 + stepTwo[1]
					if r2 < 0 || c2 < 0 || r2 >= m || c2 >= n {
						continue
					}

					if r2 == r && c2 == c {
						continue
					}

					if dist[r2][c2] == Unreachable {
						continue
					}

					h2 := r2*n + c2

					if visited[h1] != nil && visited[h1][h2] > 0 {
						continue
					}

					save := dist[r2][c2] - dist[r][c] - 2

					if save > 0 {
						if visited[h1] == nil {
							visited[h1] = map[int]int{}
						}
						visited[h1][h2] = save
						record[save]++
						if save >= saveAtLeast {
							result++
						}
					}
				}
			}

		}
	}

	printRecord(record)
	return result
}

func main() {
	filepath.Walk(".", run)
}
