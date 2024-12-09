package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func panicIf(err error) {

	if err == nil {
		return
	}

	panic(err)
}

func validReport(skip int, report []int) bool {
	if len(report) < 3 {
		return true
	}

	valid := true
	var increasing bool
	switch skip {
	case 1:
		increasing = report[2] > report[0]
	case 0:
		increasing = report[2] > report[1]
	default:
		increasing = report[1] > report[0]
	}

	for i := 1; i < len(report); i++ {
		if i == skip {
			continue
		}

		prev := i - 1
		if prev == skip {
			if i == 1 {
				continue
			}
			prev = i - 2
		}

		if report[i] == report[prev] {
			valid = false
		} else if (report[i] > report[prev]) != increasing {
			valid = false
		} else if report[i]-report[prev] > 3 || report[i]-report[prev] < -3 {
			valid = false
		}
	}

	return valid
}

func parseReport(sc *bufio.Scanner) []int {
	buff := sc.Text()
	reportStr := strings.Fields(buff)
	if len(reportStr) == 0 {
		return nil
	}
	report := make([]int, len(reportStr))

	var err error
	report[0], err = strconv.Atoi(reportStr[0])
	panicIf(err)

	for i := 1; i < len(report); i++ {
		report[i], err = strconv.Atoi(reportStr[i])
		panicIf(err)
	}
	return report
}

func main() {
	f, err := os.OpenFile("./input2.txt", os.O_RDONLY, 660)
	panicIf(err)
	defer f.Close()
	sc := bufio.NewScanner(f)

	partOneCount := 0
	partTwoCount := 0
	for sc.Scan() {
		report := parseReport(sc)
		if len(report) < 2 {
			continue
		}

		for skip := -1; skip < len(report); skip++ {
			valid := validReport(skip, report)

			if valid {
				if skip == -1 {
					partOneCount++
				}

				partTwoCount++
				break
			}
		}
	}
	fmt.Println(partOneCount)
	fmt.Println(partTwoCount)
}
