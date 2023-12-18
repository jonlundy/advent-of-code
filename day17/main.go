package main

import (
	"bufio"
	_ "embed"
	"fmt"

	aoc "go.sour.is/advent-of-code"
)

// var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {
	var m Map
	var pq aoc.PriorityQueue[int, uint]

	for scan.Scan() {
		text := scan.Text()
		m = append(m, []rune(text))
	}

	rows := len(m)
	cols := len(m[0])

	END := [2]int{rows-1, cols-1}

	return &result{}, nil
}

var (
	ZERO = [2]int{0, 0}

	UP = [2]int{-1, 0}
	DN = [2]int{1, 0}
	LF = [2]int{0, -1}
	RT = [2]int{0, 1}
)

type Map [][]rune

func (m *Map) Get(p [2]int) rune {
	if p[0] < 0 || p[0] >= len((*m)) {
		return 0
	}
	if p[1] < 0 || p[1] >= len((*m)[0]) {
		return 0
	}

	return (*m)[p[0]][p[1]]
}
