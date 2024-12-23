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

func dfs(towel string, index int, node *TrieNode, root *TrieNode) bool {
	if index == len(towel) {
		return node.isEnd
	}

	c := rune(towel[index])
	if node.children[c] == nil {
		return false
	}

	result := false
	node = node.children[c]
	if dfs(towel, index+1, node, root) {
		result = true
	} else if node.isEnd && dfs(towel, index+1, root, root) {
		result = true
	}
	return result
}

// for input1, reverse the input is much much faster
func do_dfs(root *TrieNode, workers <-chan string, found chan<- bool, done chan<- struct{}) {
	for towel := range workers {
		result := dfs(towel, 0, root, root)
		towel = reverse(towel)

		found <- result
	}
	done <- struct{}{}
}

func collectResult(found <-chan bool) <-chan int {
	result := make(chan int)
	go func() {
		count := 0
		for f := range found {
			if f {
				count++
			}
		}
		result <- count
	}()
	return result
}

func filterPossible(towels []string, root *TrieNode) int {
	WORKERS := 16

	workers := make(chan string)
	done := make(chan struct{}, WORKERS)
	found := make(chan bool)
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

	return <-result
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

	partOneResult := filterPossible(towels, trie)

	fmt.Printf("%s, partOne: %d\n", info.Name(), partOneResult)

	return nil
}

func main() {
	filepath.Walk(".", run)
}
