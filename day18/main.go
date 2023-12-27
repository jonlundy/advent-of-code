package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	aoc "go.sour.is/advent-of-code"
	"golang.org/x/exp/maps"
)

// var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {

	var vecsPT1 []aoc.Vector
	var vecsPT2 []aoc.Vector

	for scan.Scan() {
		text := scan.Text()

		if len(text) == 0 {
			continue
		}

		v, color := fromLine(text)

		vecsPT1 = append(vecsPT1, v)
		vecsPT2 = append(vecsPT2, fromColor(color))
	}
	return &result{
		valuePT1: findArea(vecsPT1),
		valuePT2: findArea(vecsPT2),
	}, nil
}

var OFFSET = map[string]aoc.Point{
	"R": {0, 1},
	"D": {1, 0},
	"L": {0, -1},
	"U": {-1, 0},
}
var OFFSET_INDEXES = maps.Values(OFFSET)

func fromLine(text string) (aoc.Vector, string) {
	v := aoc.Vector{}
	s, text, _ := strings.Cut(text, " ")
	v.Offset = OFFSET[s]

	s, text, _ = strings.Cut(text, " ")
	v.Scale = aoc.Atoi(s)

	_, text, _ = strings.Cut(text, "(#")
	s, _, _ = strings.Cut(text, ")")
	return v, s
}
 
func fromColor(c string) aoc.Vector {
	scale, _ := strconv.ParseInt(c[:5], 16, 64)
	offset := OFFSET_INDEXES[c[5]-'0']

	return aoc.Vector{
		Offset: offset,
		Scale:  int(scale),
	}
}

func findArea(vecs []aoc.Vector) int {
	shoelace := []aoc.Point{{0,0}}
	borderLength := 0

	for _, vec := range vecs {
		shoelace = append(shoelace, shoelace[len(shoelace)-1].Add(vec.Point()))
		borderLength += vec.Scale
	}

	return aoc.NumPoints(shoelace, borderLength)
}

