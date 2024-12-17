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

func parseCost(line string, ch string) (int, int) {
	var x int
	var y int
	fmt.Sscanf(line, "Button "+ch+": X+%d, Y+%d", &x, &y)
	return x, y
}

func parseTarget(line string) (int, int) {
	var x int
	var y int
	fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y)
	return x, y
}

func solve(ax int, ay int, bx int, by int, x int, y int, offset int) (int, int) {
	x += offset
	y += offset
	// a * ax + b * bx == x
	// a * ay + b * by == y
	Det := ax*by - bx*ay
	if Det == 0 {
		return 0, 0
	}

	a := (by*x - bx*y) / Det
	b := (-ay*x + ax*y) / Det

	if a*ax+b*bx != x {
		return 0, 0
	}

	if a*ay+b*by != y {
		return 0, 0
	}
	return a, b
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	partTwoResult := 0
	partOneResult := 0
	for sc.Scan() {
		Aline := sc.Text()
		sc.Scan()
		Bline := sc.Text()
		sc.Scan()
		Pline := sc.Text()

		ax, ay := parseCost(Aline, "A")
		bx, by := parseCost(Bline, "B")
		x, y := parseTarget(Pline)
		a, b := solve(ax, ay, bx, by, x, y, 0)
		a2, b2 := solve(ax, ay, bx, by, x, y, 10000000000000)
		partOneResult += a*3 + b
		partTwoResult += a2*3 + b2
		if sc.Scan() {
			sc.Text() // remove the empty line
		}
	}

	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)

	return nil
}

func main() {
	err := filepath.Walk(".", run)
	panicIf(err)
}
