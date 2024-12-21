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

const Corrupt rune = '#'
const Visited rune = 'O'
const Empty rune = '.'

const Unreachable int = 1e9

var DIRECTIONS [][]int = [][]int{
	{0, 1}, {0, -1}, {-1, 0}, {1, 0},
}

func reachExit(corrupts [][]int, size int, r int, c int) int {
	board := make([][]rune, size+1)
	for i := range size + 1 {
		board[i] = make([]rune, size+1)
	}

	for _, corrupt := range corrupts {
		board[corrupt[0]][corrupt[1]] = Corrupt
	}

	q := [][3]int{}
	q = append(q, [3]int{0, r, c})

	m := len(board)
	n := len(board[0])

	dist := make([][]int, m)
	for i := range m {
		dist[i] = make([]int, n)
		for j := range n {
			dist[i][j] = Unreachable
		}
	}
	dist[r][c] = 0

	for len(q) > 0 {
		front := q[0]
		q = q[1:]

		d := front[0]
		row := front[1]
		col := front[2]

		if row == m-1 && col == n-1 {
			continue
		}

		if dist[row][col] < d {
			continue
		}

		for _, dir := range DIRECTIONS {
			nr := row + dir[0]
			nc := col + dir[1]
			if nr < 0 || nc < 0 || nr >= m || nc >= n {
				continue
			}

			if board[nr][nc] == Corrupt {
				continue
			}

			if dist[nr][nc] <= dist[row][col]+1 {
				continue
			}
			dist[nr][nc] = dist[row][col] + 1
			q = append(q, [3]int{dist[nr][nc], nr, nc})
		}
	}

	return dist[m-1][n-1]
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
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Scan()
	header := sc.Text()
	var count, size int
	fmt.Sscanf(header, "%d,%d", &count, &size)

	corrupts := [][]int{}
	for sc.Scan() {
		var r, c int
		fmt.Sscanf(sc.Text(), "%d,%d", &c, &r)
		corrupts = append(corrupts, []int{r, c})
	}

	partOneResult := reachExit(corrupts[:count], size, 0, 0)
	stoppable := 0

	for corrupt := range corrupts {
		steps := reachExit(corrupts[:corrupt+1], size, 0, 0)
		if steps == Unreachable {
			stoppable = corrupt
			break
		}
	}

	partTwoResult := strconv.Itoa(corrupts[stoppable][1]) + "," + strconv.Itoa(corrupts[stoppable][0])

	fmt.Printf("%s, partOne: %d, partTwo: %s\n", info.Name(), partOneResult, partTwoResult)
	return nil
}

func main() {
	filepath.Walk(".", run)
}
