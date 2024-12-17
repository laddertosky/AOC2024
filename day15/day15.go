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

const Box rune = 'O'
const BoxL rune = '['
const BoxR rune = ']'
const Robot rune = '@'
const Wall rune = '#'
const Empty rune = '.'

const Up rune = '^'
const Left rune = '<'
const Down rune = 'v'
const Right rune = '>'

type Move struct {
	dr, dc int
}

func movement(move rune) Move {
	switch move {
	case Up:
		return Move{-1, 0}
	case Left:
		return Move{0, -1}
	case Down:
		return Move{1, 0}
	case Right:
		return Move{0, 1}
	}

	panic("Unknown move: " + string(move))
}

func printBoard(board [][]rune, r int, c int, move rune) {
	fmt.Println("move: ", string(move))
	for i := range len(board) {
		for j := range len(board[0]) {
			if i == r && j == c {
				fmt.Print(string(Robot))
			} else {
				fmt.Print(string(board[i][j]))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func calCoordinates(board [][]rune) int {
	result := 0
	for r := range len(board) {
		for c := range len(board[0]) {
			if board[r][c] == Box || board[r][c] == BoxL {
				result += r*100 + c
			}
		}
	}
	return result
}

func moveBox(r int, c int, board [][]rune, moves []rune) {
	for _, m := range moves {
		move := movement(m)

		nr := r + move.dr
		nc := c + move.dc
		if board[nr][nc] == Empty {
			r, c = nr, nc
		} else if board[nr][nc] == Wall {
			// do nothing
		} else {
			for board[nr][nc] == Box {
				nr += move.dr
				nc += move.dc
			}

			if board[nr][nc] == Empty {
				board[nr][nc] = Box
				r += move.dr
				c += move.dc
				board[r][c] = Empty
			}
		}
		//printBoard(board, r, c, m)
	}
}

type BoxPoint struct {
	r, c int
}

func moveBoxWider(r int, c int, board [][]rune, moves []rune) {
	for _, m := range moves {
		move := movement(m)

		nr := r + move.dr
		nc := c + move.dc
		if board[nr][nc] == Empty {
			r, c = nr, nc
		} else if board[nr][nc] == Wall {
			// do nothing
		} else {
			boxes := map[BoxPoint]bool{}
			current := BoxPoint{nr, nc}
			fronts := []BoxPoint{current}
			stack := []BoxPoint{}
			boxes[current] = true

			if board[nr][nc] == BoxL {
				pair := BoxPoint{nr, nc + 1}
				boxes[pair] = true
				fronts = append(fronts, pair)
			} else {
				pair := BoxPoint{nr, nc - 1}
				boxes[pair] = true
				fronts = append(fronts, pair)
			}

			allMovable := true
			for len(fronts) > 0 && allMovable {
				front := fronts[0]
				fronts = fronts[1:]
				stack = append(stack, front)

				nnr := front.r + move.dr
				nnc := front.c + move.dc
				next := BoxPoint{nnr, nnc}
				if boxes[next] {
					continue
				}

				boxes[next] = true

				if board[nnr][nnc] == BoxL {
					pair := BoxPoint{nnr, nnc + 1}

					boxes[pair] = true
					if m == Left {
						fronts = append(fronts, pair)
						fronts = append(fronts, next)
					} else {
						fronts = append(fronts, next)
						fronts = append(fronts, pair)
					}
				} else if board[nnr][nnc] == BoxR {
					pair := BoxPoint{nnr, nnc - 1}
					boxes[pair] = true

					if m == Right {
						fronts = append(fronts, pair)
						fronts = append(fronts, next)
					} else {
						fronts = append(fronts, next)
						fronts = append(fronts, pair)
					}
				} else if board[nnr][nnc] == Wall {
					allMovable = false
				}
			}

			if allMovable {
				// move from frontmost element
				// printBoard(board, r, c, '0')
				for i := range len(stack) {
					box := stack[len(stack)-1-i]
					board[box.r+move.dr][box.c+move.dc] = board[box.r][box.c]
					board[box.r][box.c] = Empty
					// printBoard(board, r, c, m)
				}
				r += move.dr
				c += move.dc
				// printBoard(board, r, c, '1')
				board[r][c] = Empty
			}
		}
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
	widerBoard := [][]rune{}
	moves := []rune{}
	var startR int
	var startC int

	boardCompleted := false
	for sc.Scan() {
		line := sc.Text()
		if boardCompleted {
			for _, c := range line {
				moves = append(moves, c)
			}
		} else if len(line) == 0 {
			boardCompleted = true
		} else {
			row := make([]rune, len(line))
			widerRow := make([]rune, len(line)*2)
			for i, c := range line {
				if c == Robot {
					startR = len(board)
					startC = i
					row[i] = Empty
					widerRow[2*i] = Empty
					widerRow[2*i+1] = Empty
				} else if c == Box {
					row[i] = Box
					widerRow[2*i] = BoxL
					widerRow[2*i+1] = BoxR
				} else {
					row[i] = c
					widerRow[2*i] = c
					widerRow[2*i+1] = c
				}
			}
			board = append(board, row)
			widerBoard = append(widerBoard, widerRow)
		}
	}

	moveBox(startR, startC, board, moves)
	partOneResult := calCoordinates(board)

	moveBoxWider(startR, 2*startC, widerBoard, moves)
	partTwoResult := calCoordinates(widerBoard)

	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)

	return nil
}

func main() {
	filepath.Walk(".", run)
}
