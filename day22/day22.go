package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func prune(input int) int {
	return input & 0xffffff
}

func mix(a int, b int) int {
	return a ^ b
}

func generate(secret int, level int) (int, []int) {
	prices := []int{secret % 10}
	for range level {
		secret = prune(mix(secret, secret<<6))
		secret = prune(mix(secret, secret>>5))
		secret = prune(mix(secret, secret<<11))

		prices = append(prices, secret%10)
	}

	return secret, prices
}

// find the best 4 sequence
func sell(buyers [][]int) int {
	sequences := map[string]int{}
	for _, prices := range buyers {
		bestSold := map[string]int{}

		changes := []string{
			strconv.Itoa((prices[1] - prices[0])),
			strconv.Itoa((prices[2] - prices[1])),
			strconv.Itoa((prices[3] - prices[2])),
			strconv.Itoa((prices[4] - prices[3])),
		}
		key := strings.Join(changes, ",")
		bestSold[key] = prices[4]

		for i := 4; i < len(prices); i++ {
			change := strconv.Itoa((prices[i] - prices[i-1]))
			changes = changes[1:]
			changes = append(changes, change)
			key := strings.Join(changes, ",")

			if bestSold[key] == 0 {
				bestSold[key] = prices[i]
			}
		}

		for k, price := range bestSold {
			sequences[k] += price
		}
	}

	result := 0
	maxKey := ""
	for k, gain := range sequences {
		if gain > result {
			result = gain
			maxKey = k
		}
	}
	fmt.Printf("key: %s, gain: %d\n", maxKey, result)
	return result
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)

	sc := bufio.NewScanner(f)

	secrets := []int{}
	partOneResult := 0
	for sc.Scan() {
		secret, err := strconv.Atoi(sc.Text())
		panicIf(err)
		secrets = append(secrets, secret)
	}

	buyers := [][]int{}
	for _, secret := range secrets {
		lastSecret, buyer := generate(secret, 2000)
		buyers = append(buyers, buyer)
		partOneResult += lastSecret
	}

	partTwoResult := sell(buyers)
	fmt.Printf("%s, partOne: %d, partTwo: %d\n", info.Name(), partOneResult, partTwoResult)
	return nil
}

func main() {
	filepath.Walk(".", run)
}
