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

const Wall byte = '#'
const Empty byte = '.'
const Start byte = '^'
const Visited byte = 'x'

var DIRECTIONS [4][2]int = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func isValid(r int, c int, world [][]byte) bool {
	return r >= 0 && r < len(world) && c >= 0 && c < len(world[0])
}

func findStart(world [][]byte) (r int, c int) {
	for i := range len(world) {
		for j := range len(world[0]) {
			if world[i][j] == Start {
				r = i
				c = j
				world[i][j] = Empty
			}
		}
	}
	return
}

func patrol(r int, c int, world [][]byte) int {
	count := 0

	dir := 0
	for isValid(r, c, world) {
		if world[r][c] == Wall {
			r -= DIRECTIONS[dir][0]
			c -= DIRECTIONS[dir][1]
			dir = (dir + 1) % 4
			continue
		}

		if world[r][c] != Visited {
			count++
		}

		world[r][c] = Visited

		r += DIRECTIONS[dir][0]
		c += DIRECTIONS[dir][1]
	}

	return count
}

func duplicate(world [][]byte) [][]byte {
	dup := [][]byte{}

	for r := range len(world) {
		tmp := make([]byte, len(world[r]))
		copy(tmp, world[r])
		dup = append(dup, tmp)
	}
	return dup
}

func isVisited(world [][]byte, r int, c int, dir int) bool {
	if world[r][c] >= 1<<len(DIRECTIONS) {
		return false
	}
	return (world[r][c] & (1 << dir)) > 0
}

func setVisited(world [][]byte, r int, c int, dir int) {
	if world[r][c] == Empty {
		world[r][c] = 0
	}

	world[r][c] += (1 << dir)
}

func makeLoops(startR int, startC int, world [][]byte) int {
	count := 0

	for i := range len(world) {
		for j := range len(world[0]) {
			if world[i][j] == Empty {
				parallelWorld := duplicate(world)

				parallelWorld[i][j] = Wall

				dir := 0
				r := startR
				c := startC
				for isValid(r, c, world) && !isVisited(parallelWorld, r, c, dir) {
					if parallelWorld[r][c] == Wall {
						r -= DIRECTIONS[dir][0]
						c -= DIRECTIONS[dir][1]
						dir = (dir + 1) % 4
						continue
					}

					setVisited(parallelWorld, r, c, dir)
					r += DIRECTIONS[dir][0]
					c += DIRECTIONS[dir][1]
				}

				if isValid(r, c, world) {
					count++
				}
			}
		}
	}

	return count
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	partOneResult := 0
	partTwoResult := 0

	world := [][]byte{}
	for sc.Scan() {
		buff := sc.Text()
		world = append(world, []byte(buff))
	}

	r, c := findStart(world)
	partOneResult = patrol(r, c, duplicate(world))
	partTwoResult = makeLoops(r, c, duplicate(world))

	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)
}
