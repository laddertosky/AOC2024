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

var DIRECTIONS [][]int = [][]int{
	{0, 1}, {0, -1}, {1, 0}, {-1, 0},
}

var NEIGHBORS [][]int = [][]int{
	{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1},
}

func hasCorner(pattern [3]bool) int {
	if pattern == [3]bool{false, true, false} {
		return 1
	}

	trues := 0
	for _, same := range pattern {
		if same {
			trues++
		}
	}

	if trues%2 == 0 {
		return 1
	}

	return 0
}

func bfs(garden []string, visited []bool, m int, n int, r int, c int) (int, int, int) {
	queue := [][]int{{r, c}}

	area := 0
	perimeter := 0
	side := 0
	corners := make([][]int, m+1)
	for i := range m + 1 {
		corners[i] = make([]int, n+1)
	}

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		row := front[0]
		col := front[1]

		area++

		for i := range 4 {
			sameRegion := [3]bool{}
			for j := range 3 {
				move := NEIGHBORS[(i*2+j)%8]
				nr := row + move[0]
				nc := col + move[1]
				if nr < 0 || nc < 0 || nr >= m || nc >= n {
					sameRegion[j] = false
				} else {
					sameRegion[j] = garden[row][col] == garden[nr][nc]
				}
			}
			switch i {
			case 0:
				corners[row][col] += hasCorner(sameRegion)
			case 1:
				corners[row][col+1] += hasCorner(sameRegion)
			case 2:
				corners[row+1][col+1] += hasCorner(sameRegion)
			case 3:
				corners[row+1][col] += hasCorner(sameRegion)
			}
		}

		for _, dir := range DIRECTIONS {
			nr := row + dir[0]
			nc := col + dir[1]

			if nr < 0 || nc < 0 || nr >= m || nc >= n {
				perimeter++
				continue
			}

			if garden[r][c] != garden[nr][nc] {
				perimeter++
				continue
			}

			if visited[nr*n+nc] {
				continue
			}

			visited[nr*n+nc] = true
			queue = append(queue, []int{nr, nc})
		}
	}

	for r := range m + 1 {
		for c := range n + 1 {
			// one corner might be touch by this region 1/2/3 times,
			// for 2 times, it should look like AE
			//									EA
			//
			// for 3 times, it should look like AA
			//									AE
			if corners[r][c] == 3 {
				side++
			} else {
				side += corners[r][c]
			}
		}
	}

	return area, perimeter, side
}

func fencePrice(garden []string) (int, int) {
	partOneResult := 0
	partTwoResult := 0

	m := len(garden)
	n := len(garden[0])

	visited := make([]bool, m*n)

	for r := range m {
		for c := range n {
			if visited[r*n+c] {
				continue
			}

			visited[r*n+c] = true

			area, perimeter, side := bfs(garden, visited, m, n, r, c)
			partOneResult += area * perimeter
			partTwoResult += area * side
			// fmt.Printf("%s, area: %d, perimeter: %d, side: %d\n", string(garden[r][c]), area, perimeter, side)
		}
	}

	return partOneResult, partTwoResult
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	garden := []string{}

	for sc.Scan() {
		row := sc.Text()
		garden = append(garden, row)
	}

	partOneResult, partTwoResult := fencePrice(garden)
	fmt.Printf("%s, partOne: %d, partTwo: %d\n", path, partOneResult, partTwoResult)
	return nil
}

func main() {
	err := filepath.Walk(".", run)
	panicIf(err)
}
