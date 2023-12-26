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

	for scan.Scan() {
		_ = scan.Text()

	}

	return &result{}, nil
}
