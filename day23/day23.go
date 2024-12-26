package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Graph map[string]map[string]bool

func findPwd(graph Graph) string {
	groups := map[string]map[string][]string{}
	for u, vs := range graph {
		for v := range vs {
			if groups[u] == nil {
				groups[u] = map[string][]string{}
				groups[u][v] = []string{u, v}
			}

			for w := range graph[v] {
				if w == u {
					continue
				}

				allConnected := true
				for _, x := range groups[u][v] {
					if !graph[x][w] {
						allConnected = false
					}
				}
				if allConnected {
					groups[u][v] = append(groups[u][v], w)
				}
			}
		}
	}

	largest := []string{}
	size := 0

	for _, vs := range groups {
		for _, group := range vs {
			if len(group) > size {
				largest = group
				size = len(group)
			} else if len(group) == size {

			}

		}
	}

	sort.Strings(largest)

	return strings.Join(largest, ",")
}

func findSets(graph Graph, prefix string) int {
	result := 0
	for u, vs := range graph {
		for v := range vs {
			for w := range graph[v] {
				if w == u {
					continue
				}

				pu := string(u[0])
				pv := string(v[0])
				pw := string(w[0])
				if pu != prefix && pv != prefix && pw != prefix {
					continue
				}

				if graph[w][u] {
					result++
				}
			}
		}
	}
	return result / 6
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	graph := Graph{}
	for sc.Scan() {
		connection := sc.Text()
		sepPos := strings.Index(connection, "-")
		u := connection[0:sepPos]
		v := connection[sepPos+1:]

		if graph[u] == nil {
			graph[u] = map[string]bool{}
		}
		if graph[v] == nil {
			graph[v] = map[string]bool{}
		}
		graph[u][v] = true
		graph[v][u] = true
	}

	partOneResult := findSets(graph, "t")
	partTwoResult := findPwd(graph)
	fmt.Printf("%s, partOne: %d, partTwo: %s\n", info.Name(), partOneResult, partTwoResult)
	return nil
}

func main() {
	filepath.Walk(".", run)
}
