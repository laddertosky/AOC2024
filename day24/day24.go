package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Gate func(a int, b int) int
type Instruction struct {
	input1, input2, output string
	gate                   Gate
}
type Graph map[string][]string
type Wire struct {
	name  string
	value int
}

func simulate(instructions map[string]Instruction, registers map[string]int) {
	graph := Graph{}
	indegree := map[string]int{}
	for _, instruction := range instructions {
		graph[instruction.input1] = append(graph[instruction.input1], instruction.output)
		graph[instruction.input2] = append(graph[instruction.input2], instruction.output)
		indegree[instruction.input1] += 0
		indegree[instruction.input2] += 0
		indegree[instruction.output] += 2
	}

	queue := []string{}
	for node, count := range indegree {
		if count == 0 {
			queue = append(queue, node)
			if node[0] != 'x' && node[0] != 'y' {
				panic(node)
			}
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, next := range graph[node] {
			indegree[next]--
			if indegree[next] == 0 {
				instruction := instructions[next]
				input1 := registers[instruction.input1]
				input2 := registers[instruction.input2]
				registers[next] = instruction.gate(input1, input2)
				queue = append(queue, next)
			}
		}
	}

}

func calculateOutput(registers map[string]int) int {
	output := []Wire{}
	for name, value := range registers {
		if name[0] == 'z' {
			output = append(output, Wire{
				name:  name,
				value: value,
			})
		}
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].name > output[j].name
	})

	result := 0
	for _, wire := range output {
		result <<= 1
		result |= wire.value
	}
	return result
}

func AND(a int, b int) int {
	return a & b
}

func XOR(a int, b int) int {
	return a ^ b
}

func OR(a int, b int) int {
	return a | b
}

func run(path string, info fs.FileInfo, err error) error {
	if !strings.Contains(info.Name(), ".txt") {
		return nil
	}

	f, err := os.Open(info.Name())
	panicIf(err)
	defer f.Close()

	sc := bufio.NewScanner(f)

	registers := map[string]int{}
	instructions := map[string]Instruction{}
	sus := []Instruction{}
	direct := map[string]string{}
	firstComplete := false
	for sc.Scan() {
		line := sc.Text()
		if firstComplete {
			instruction := Instruction{}
			words := strings.Split(line, " ")
			instruction.input1 = words[0]
			instruction.input2 = words[2]
			instruction.output = words[4]
			switch words[1] {
			case "AND":
				instruction.gate = AND
			case "XOR":
				instruction.gate = XOR
			case "OR":
				instruction.gate = OR
			default:
				panic("unsupported operator: " + words[1])
			}

			if instruction.output[0] == 'z' && words[1] != "XOR" {
				sus = append(sus, instruction)
			} else if (instruction.input1[0] == 'x' || instruction.input1[0] == 'y') && words[1] == "OR" {
				sus = append(sus, instruction)
			} else if words[1] == "XOR" && instruction.input1[0] != 'x' && instruction.input1[0] != 'y' && instruction.output[0] != 'z' {
				sus = append(sus, instruction)
			}

			if (instruction.input1[0] == 'x' || instruction.input1[0] == 'y') && words[1] == "XOR" {
				key := instruction.input1[1:]
				direct[key] = instruction.output
			}

			instructions[instruction.output] = instruction
		} else if len(line) == 0 {
			firstComplete = true
		} else {
			pos := strings.Index(line, ":")
			name := line[0:pos]
			bit, err := strconv.Atoi(line[pos+2:])
			panicIf(err)
			registers[name] = bit
		}
	}

	for key, instruction := range instructions {
		if key[0] != 'z' {
			continue
		}
		zKey := key[1:]

		found := false
		if instruction.input1 == direct[zKey] || instruction.input2 == direct[zKey] {
			found = true
		}

		if !found {
			sus = append(sus, instruction)
		}
	}

	simulate(instructions, registers)
	partOneResult := calculateOutput(registers)
	partTwoResult := findSwapped(instructions, registers)
	fmt.Println(sus)

	fmt.Printf("%s, partOne: %d, partTwo: %s\n", info.Name(), partOneResult, partTwoResult)
	return nil
}

func calculateExpected(instructions map[string]Instruction, registers map[string]int) (map[string]int, int) {
	outputSize := 0
	for _, instruction := range instructions {
		if instruction.output[0] == 'z' {
			id, err := strconv.Atoi(instruction.output[1:])
			panicIf(err)
			outputSize = max(outputSize, id+1)
		}
	}
	expected := map[string]int{}
	target := 0
	carry := 0
	base := 1
	for id := range outputSize {
		kx := fmt.Sprintf("x%02d", id)
		ky := fmt.Sprintf("y%02d", id)
		kz := fmt.Sprintf("z%02d", id)
		rx := registers[kx]
		ry := registers[ky]
		rz := (rx + ry + carry) & 1
		carry = (rx + ry + carry) >> 1
		expected[kz] = rz

		target += rz * base
		base <<= 1
	}
	return expected, target
}

func swap(instructions map[string]Instruction, a string, b string) map[string]Instruction {
	result := instructions
	insA := instructions[a]
	insB := instructions[b]

	tmp := Instruction{
		input1: insA.input1,
		input2: insA.input2,
		gate:   insA.gate,
	}

	insA.input1 = insB.input1
	insA.input2 = insB.input2
	insA.gate = insB.gate
	insB.input1 = tmp.input1
	insB.input2 = tmp.input2
	insB.gate = tmp.gate
	result[a] = insA
	result[b] = insB

	return result
}

func findSwapped(instructions map[string]Instruction, registers map[string]int) string {
	// obtained from the sus array
	instructions = swap(instructions, "qbw", "z14")
	instructions = swap(instructions, "wcb", "z34")
	instructions = swap(instructions, "mkk", "z10")
	instructions = swap(instructions, "cvp", "wjb")
	simulate(instructions, registers)

	expected, target := calculateExpected(instructions, registers)
	current := calculateOutput(registers)

	if target != current {
		fmt.Println(expected)
	}

	result := []string{
		"qbw", "z14", "wcb", "z34", "mkk", "z10", "cvp", "wjb",
	}

	sort.Strings(result)
	return strings.Join(result, ",")
}

func main() {
	filepath.Walk(".", run)
}
