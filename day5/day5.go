package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Graph map[string]map[string]bool

func addToGraph(buff string, graph Graph) {
	nodes := strings.Split(buff, "|")
	if graph[nodes[0]] == nil {
		graph[nodes[0]] = map[string]bool{}
	}
	graph[nodes[0]][nodes[1]] = true
}

func midPage(buff string, graph Graph) (bool, int) {
	nodes := strings.Split(buff, ",")

	valid := true

	for i := 0; i < len(nodes); i++ {
		for j := range i {
			if graph[nodes[i]] != nil && graph[nodes[i]][nodes[j]] == true {
				valid = false

				tmp := nodes[i]
				nodes[i] = nodes[j]
				nodes[j] = tmp

				i = 0
				break
			}
		}
	}

	result, err := strconv.Atoi(nodes[len(nodes)/2])
	panicIf(err)

	return valid, result
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("./input1.txt")
	panicIf(err)
	firstCompleted := false
	graph := Graph{}

	sc := bufio.NewScanner(f)
	partOneResult := 0
	partTwoResult := 0
	for sc.Scan() {
		buff := sc.Text()

		if len(buff) == 0 {
			firstCompleted = true
		} else if !firstCompleted {
			addToGraph(buff, graph)
		} else {
			valid, mid := midPage(buff, graph)
			if valid {
				partOneResult += mid
			} else {
				partTwoResult += mid
			}
		}
	}

	fmt.Println(partOneResult)
	fmt.Println(partTwoResult)
}
