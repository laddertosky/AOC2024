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

type TrieNode struct {
	isEnd    bool
	children map[rune]*TrieNode
}

func buildTrie(patterns []string) *TrieNode {
	root := TrieNode{}

	for _, pattern := range patterns {
		pattern = reverse(pattern)
		current := &root

		for _, c := range pattern {
			if current.children == nil {
				current.children = make(map[rune]*TrieNode)
			}
			if current.children[c] == nil {
				current.children[c] = &TrieNode{}
			}
			current = current.children[c]
		}
		current.isEnd = true
	}

	return &root
}

func dfs(cache map[string]map[*TrieNode]int, towel string, index int, node *TrieNode, root *TrieNode) int {
	if index == len(towel) {
		if node.isEnd {
			return 1
		}
		return 0
	}

	after := towel[index:]
	if cache[after] != nil && cache[after][node] > 0 {
		return cache[after][node]
	}

	c := rune(towel[index])
	if node.children[c] == nil {
		return 0
	}

	result := 0
	next := node.children[c]
	result += dfs(cache, towel, index+1, next, root)

	if next.isEnd {
		result += dfs(cache, towel, index+1, root, root)
	}

	if cache[after] == nil {
		cache[after] = map[*TrieNode]int{}
	}
	cache[after][node] = result
	return result
}

// for input1, reverse the input is much much faster
func do_dfs(root *TrieNode, workers <-chan string, found chan<- int, done chan<- struct{}) {
	cache := map[string]map[*TrieNode]int{}
	for towel := range workers {
		result := dfs(cache, towel, 0, root, root)

		found <- result
	}
	done <- struct{}{}
}

func collectResult(found <-chan int) <-chan int {
	result := make(chan int)
	go func() {
		partOne := 0
		partTwo := 0
		for count := range found {
			if count > 0 {
				partOne++
			}
			partTwo += count
		}
		result <- partOne
		result <- partTwo
	}()
	return result
}

func filterPossible(towels []string, root *TrieNode) <-chan int {
	WORKERS := 16

	workers := make(chan string)
	done := make(chan struct{}, WORKERS)
	found := make(chan int)
	result := collectResult(found)

	for range WORKERS {
		go do_dfs(root, workers, found, done)
	}

	for _, towel := range towels {
		workers <- towel
	}

	close(workers)
	for range WORKERS {
		<-done
	}

	close(found)
	return result
}

func reverse(input string) string {
	result := make([]rune, len(input))
	for i := range len(input)/2 + 1 {
		tmp := input[i]
		j := len(input) - 1 - i
		result[i] = rune(input[j])
		result[j] = rune(tmp)
	}
	output := string(result)
	return output
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)

	sc := bufio.NewScanner(f)

	sc.Scan()
	header := sc.Text()
	patterns := strings.Split(header, ", ")
	trie := buildTrie(patterns)

	towels := []string{}
	sc.Scan()
	for sc.Scan() {
		towel := sc.Text()
		towel = reverse(towel)
		towels = append(towels, towel)
	}

	result := filterPossible(towels, trie)
	partOneResult := <-result
	partTwoResult := <-result

	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)

	return nil
}

func main() {
	filepath.Walk(".", run)
}
