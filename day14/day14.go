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

func parsePosition(line string) (px int, py int, vx int, vy int) {
	fmt.Sscanf(line, "p=%d,%d v=%d,%d", &px, &py, &vx, &vy)
	return
}

func parseSize(line string) (w int, h int) {
	fmt.Sscanf(line, "%d,%d", &w, &h)
	return
}

type Quadrant struct {
	one, two, three, four int
}

func (q Quadrant) safetyFactor() int {
	return q.one * q.two * q.three * q.four
}

func (q *Quadrant) move(px int, py int, vx int, vy int, w int, h int, time int) {
	px += time * vx
	px %= w
	px = (px + w) % w
	py += time * vy
	py %= h
	py = (py + h) % h

	if px < w/2 && py < h/2 {
		q.one++
	} else if px < w/2 && py > h/2 {
		q.two++
	} else if px > w/2 && py < h/2 {
		q.three++
	} else if px > w/2 && py > h/2 {
		q.four++
	}
}

type Robot struct {
	px, py, vx, vy int
}

func move(robots []Robot, w int, h int) {
	for i := range robots {
		robots[i].px += robots[i].vx
		robots[i].px %= w
		robots[i].px = (robots[i].px + w) % w
		robots[i].py += robots[i].vy
		robots[i].py %= h
		robots[i].py = (robots[i].py + h) % h
	}
}

func updateBoard(robots []Robot, board [][]rune, time int) bool {
	for r := range len(board) {
		for c := range len(board[0]) {
			board[r][c] = '.'
		}
	}

	for i := range robots {
		board[robots[i].py][robots[i].px] = '#'
	}

	count := 0
	maxCount := 0
	for r := range len(board) {
		for c := range len(board[0]) {
			if board[r][c] == '#' {
				count++
				maxCount = int(max(maxCount, count))
			} else {
				count = 0
			}
		}
	}

	// Arbitrary number
	if maxCount > 10 {
		fmt.Printf("\nTrial: %d\n", time)
		for r := range len(board) {
			for c := range len(board[0]) {
				fmt.Print(string(board[r][c]))
			}
			fmt.Println()
		}
		return true
	}
	return false
}

func run(_ string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Scan()
	header := sc.Text()
	w, h := parseSize(header)
	q := Quadrant{}
	board := make([][]rune, h)
	for r := range h {
		board[r] = make([]rune, w)
	}

	robots := []Robot{}

	for sc.Scan() {
		line := sc.Text()
		px, py, vx, vy := parsePosition(line)
		robots = append(robots, Robot{px, py, vx, vy})
		q.move(px, py, vx, vy, w, h, 100)
	}
	fmt.Printf("%s, partOne: %d\n", info.Name(), q.safetyFactor())

	found := updateBoard(robots, board, 0)
	/* w > 100 to prevent infinite loop in sample test*/
	for time := 0; w > 100 && !found; time++ {
		move(robots, w, h)
		found = updateBoard(robots, board, time+1)
	}

	return nil
}

func main() {
	filepath.Walk(".", run)
}
