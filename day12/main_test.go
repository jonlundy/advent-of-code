package main

import (
	"bufio"
	"bytes"
	"testing"

	_ "embed"

	aoc "go.sour.is/advent-of-code"

	"github.com/matryer/is"
)

//go:embed example.txt
var example []byte

//go:embed input.txt
var input []byte

func TestExample(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.valuePT1, 21)
	is.Equal(result.valuePT2, 525152)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.valuePT1, 8193)
	is.Equal(result.valuePT2, 45322533163795)
}

func TestPower2(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Power2(1), 2)
	is.Equal(aoc.Power2(2), 4)
	is.Equal(aoc.Power2(3), 8)
	is.Equal(aoc.Power2(4), 16)
	is.Equal(aoc.Power2(5), 32)
	is.Equal(aoc.Power2(6), 64)
}

func TestCountGroupings(t *testing.T) {
	is := is.New(t)
	is.Equal([]int{1, 3, 1}, countGroupings([]rune(".#.###.#")))
	is.Equal([]int{1, 3, 1}, countGroupings([]rune(".#.###...#.")))
	is.Equal([]int{1, 3, 1}, countGroupings([]rune("#.###...#.")))
}

func TestCombination(t *testing.T) {
	s := spring{
		pattern:   []rune("???"),
		grouping:  []int{1},
		missingNo: 3,
	}
	s.findMatches()
}
