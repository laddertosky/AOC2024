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

const Start rune = 'S'
const End rune = 'E'
const Wall rune = '#'
const Empty rune = '.'
const Chair rune = 'O'

const Turn int = 1000
const Forward int = 1

var DIRECTIONS [][2]int = [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type State struct {
	d, r, c int
	move    int
}

func dfs(maze [][]rune, current State, dist [][][4]int) {
	m := len(maze)
	n := len(maze[0])

	if dist[current.r][current.c][current.move] < current.d {
		return
	}

	for i, dir := range DIRECTIONS {
		nr := current.r
		nc := current.c

		nextDist := current.d
		if i == current.move {
			nextDist += Forward
			nr += dir[0]
			nc += dir[1]
		} else {
			nextDist += Turn
		}

		if nr < 0 || nc < 0 || nr >= m || nc >= n {
			continue
		}

		if maze[nr][nc] == Wall {
			continue
		}

		if dist[nr][nc][i] <= nextDist {
			continue
		}

		dist[nr][nc][i] = nextDist
		nextState := State{
			d:    nextDist,
			r:    nr,
			c:    nc,
			move: i,
		}
		dfs(maze, nextState, dist)
	}
}

func solveMaze(maze [][]rune, start [2]int) [][][4]int {
	m := len(maze)
	n := len(maze[0])

	firstState := State{d: 0, r: start[0], c: start[1], move: 1}
	dist := make([][][4]int, m)
	for r := range m {
		dist[r] = make([][4]int, n)
		for c := range n {
			for move := range DIRECTIONS {
				dist[r][c][move] = 1e9
			}
		}
	}
	dfs(maze, firstState, dist)

	return dist
}

func reverseWalk(dist [][][4]int, start [2]int, end [2]int, maze [][]rune) int {
	m := len(dist)
	n := len(dist[0])
	best := map[int]bool{}
	best[end[0]*n+end[1]] = true
	maze[end[0]][end[1]] = Chair

	minMove := 0
	minDist := dist[end[0]][end[1]][0]
	for move := range DIRECTIONS {
		if dist[end[0]][end[1]][move] < minDist {
			minDist = dist[end[0]][end[1]][move]
			minMove = move
		}
	}

	queue := [][3]int{}
	queue = append(queue, [3]int{end[0], end[1], minMove})

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		if front[0] == start[0] && front[1] == start[1] {
			continue
		}

		for move, dir := range DIRECTIONS {
			nr := front[0] - dir[0]
			nc := front[1] - dir[1]
			if nr < 0 || nc < 0 || nr >= m || nr >= n {
				continue
			}

			diff := dist[front[0]][front[1]][front[2]] - dist[nr][nc][move]
			if diff == Forward || diff == Turn+Forward {
				best[nr*n+nc] = true
				maze[nr][nc] = Chair
				queue = append(queue, [3]int{nr, nc, move})
			}
		}
	}

	// debug purpose
	for r := range m {
		for c := range n {
			if maze[r][c] == Chair {
				minDist := dist[r][c][0]
				for move := range DIRECTIONS {
					if dist[r][c][move] < minDist {
						minDist = dist[r][c][move]
						minMove = move
					}
				}
				fmt.Print(minDist % 10)
			} else {
				fmt.Print(string(maze[r][c]))
			}
		}
		fmt.Println()
	}
	fmt.Println()

	return len(best)
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)

	sc := bufio.NewScanner(f)
	maze := [][]rune{}

	var start [2]int
	var end [2]int

	for sc.Scan() {
		line := sc.Text()
		row := []rune(line)
		if strings.ContainsRune(line, Start) {
			start = [2]int{len(maze), strings.IndexRune(line, Start)}
		}

		if strings.ContainsRune(line, End) {
			end = [2]int{len(maze), strings.IndexRune(line, End)}
		}
		maze = append(maze, row)
	}

	dist := solveMaze(maze, start)

	partOneResult := dist[end[0]][end[1]][0]
	for move := range DIRECTIONS {
		partOneResult = min(partOneResult, dist[end[0]][end[1]][move])
	}
	partTwoResult := reverseWalk(dist, start, end, maze)

	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)
	return nil
}

func main() {
	filepath.Walk(".", run)
}
