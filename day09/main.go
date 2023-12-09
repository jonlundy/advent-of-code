package main

import (
	"bufio"
	"cmp"
	"fmt"
	"strconv"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"
)

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

var log = aoc.Log

func run(scan *bufio.Scanner) (*result, error) {
	var histories [][]int
	var values []int
	var rvalues []int

	for scan.Scan() {
		text := scan.Text()
		if len(text) == 0 {
			continue
		}
		histories = append(histories, nil)

		for _, s := range strings.Fields(text) {
			if i, err := strconv.Atoi(s); err == nil {
				histories[len(histories)-1] = append(histories[len(histories)-1], i)
			}
		}

		log(last(histories...))

		values = append(values, predictNext(last(histories...)))
		rvalues = append(rvalues, predictPrev(last(histories...)))
	}

	log("values", values)
	log("rvalues", rvalues)

	return &result{valuePT1: sum(values...), valuePT2: sum(rvalues...)}, nil
}

func predictNext(in []int) int {
	log(" ---- PREDICT NEXT ----")
	defer log(" ----------------------")

	history := makeHistory(in)

	aoc.Reverse(history)

	return predict(history, func(a, b int) int { return a + b })
}

func predictPrev(in []int) int {
	log(" ---- PREDICT PREV ----")
	defer log(" ----------------------")

	history := makeHistory(in)

	for i := range history {
		aoc.Reverse(history[i])
	}
	aoc.Reverse(history)

	return predict(history, func(a, b int) int { return b - a })
}

func predict(history [][]int, diff func(a, b int) int) int {
	log(" ---- PREDICT ----")
	defer log(" -----------------")

	for i := range history[1:] {
		lastHistory, curHistory := last(history[i]...), last(history[i+1]...)

		history[i+1] = append(history[i+1], diff(lastHistory, curHistory))
		log(lastHistory, curHistory, last(history[i+1]))
	}

	log("last", last(history...))
	return last(last(history...)...)
}

func makeHistory(in []int) [][]int {
	var history [][]int
	history = append(history, in)

	for {
		var diffs []int

		current := history[len(history)-1]
		for i := range current[1:] {
			diffs = append(diffs, current[i+1]-current[i])
		}

		history = append(history, diffs)
		log(diffs)

		if max(diffs[0], diffs[1:]...) == 0 && min(diffs[0], diffs[1:]...) == 0 {
			break
		}
	}
	return history
}

func max[T cmp.Ordered](a T, v ...T) T {
	for _, b := range v {
		if b > a {
			a = b
		}
	}
	return a
}
func min[T cmp.Ordered](a T, v ...T) T {
	for _, b := range v {
		if b < a {
			a = b
		}
	}
	return a
}
func sum[T cmp.Ordered](v ...T) T {
	var s T
	for _, a := range v {
		s += a
	}
	return s
}
func last[T any](v ...T) T {
	return v[len(v)-1]
}
