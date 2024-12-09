package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Pair struct {
	val   int
	index int
}

func partOne() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	var n1 int
	var n2 int

	l1 := []Pair{}
	l2 := []Pair{}

	for index := 0; ; index++ {
		count, err := fmt.Fscan(f, &n1, &n2)
		if err != io.EOF {
			panicIf(err)
		}

		if count == 0 {
			break
		}
		l1 = append(l1, Pair{n1, index})
		l2 = append(l2, Pair{n2, index})
	}

	sort.Slice(l1, func(i int, j int) bool {
		return l1[i].val < l1[j].val
	})

	sort.Slice(l2, func(i int, j int) bool {
		return l2[i].val < l2[j].val
	})

	result := 0
	for i := range len(l1) {
		diff := l1[i].val - l2[i].val
		if diff < 0 {
			diff *= -1
		}
		result += diff
	}

	fmt.Println(result)
}

func partTwo() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	defer f.Close()

	l1 := []int{}
	m2 := map[int]int{}
	var n1 int
	var n2 int

	for index := 0; ; index++ {
		count, err := fmt.Fscan(f, &n1, &n2)
		if err != io.EOF {
			panicIf(err)
		}

		if count == 0 {
			break
		}

		l1 = append(l1, n1)
		m2[n2]++
	}

	result := 0
	for i := range len(l1) {
		result += m2[l1[i]] * l1[i]
	}

	fmt.Println(result)
}

func main() {
	partOne()
	partTwo()
}
