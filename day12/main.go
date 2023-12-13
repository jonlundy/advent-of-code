package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"
)

// var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {
	var matches []int

	for scan.Scan() {
		text := scan.Text()
		status, text, ok := strings.Cut(text, " ")
		if !ok {
			continue
		}

		grouping := aoc.SliceMap(aoc.Atoi, strings.Split(text, ",")...)
		pattern := []rune(status)
		missing := countQuestion(pattern)

		s := spring{pattern: pattern, grouping: grouping, missingNo: missing}

		matches = append(matches, s.findMatches())
	}

	return &result{valuePT1: aoc.Sum(matches...)}, nil
}

type spring struct {
	pattern   []rune
	grouping  []int
	missingNo int
}

func (s *spring) findMatches() int {
	matches := 0
	for _, pattern := range s.genPatterns() {
		pattern := []rune(pattern)
		target := make([]rune, len(s.pattern))
		i := 0
		for j, r := range s.pattern {
			if r == '?' {
				target[j] = pattern[i]
				i++
				continue
			}
			target[j] = r
		}

		if slices.Equal(countGroupings(target), s.grouping) {
			matches++
		}
	}

	return matches
}
func (s *spring) genPatterns() []string {
	buf := &strings.Builder{}
	combinations := aoc.Power2(s.missingNo)
	lis := make([]string, 0, combinations)
	for i := 0; i < combinations; i++ {
		for b := 0; b < s.missingNo; b++ {
			if i>>b&0b1 == 1 {
				buf.WriteRune('#')
			} else {
				buf.WriteRune('.')
			}
		}
		lis = append(lis, buf.String())
		buf.Reset()
	}

	return lis
}

func countQuestion(pattern []rune) int {
	count := 0
	for _, r := range pattern {
		if r == '?' {
			count++
		}
	}
	return count
}
func countGroupings(pattern []rune) []int {
	var groupings []int
	inGroup := false
	for _, r := range pattern {

		if r == '#' {
			if !inGroup {
				groupings = append(groupings, 0)
			}

			inGroup = true
			groupings[len(groupings)-1]++

		}
		if inGroup && r != '#' {
			inGroup = false
		}
	}
	return groupings
}

