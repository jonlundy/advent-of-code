package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"strings"

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
	var matches []int
	var matches2 []int

	for scan.Scan() {
		text := scan.Text()
		status, text, ok := strings.Cut(text, " ")
		if !ok {
			continue
		}

		// part 1 - brute force
		grouping := aoc.SliceMap(aoc.Atoi, strings.Split(text, ",")...)
		pattern := []rune(status)
		// sp := spring{pattern: pattern, grouping: grouping, missingNo: countQuestion(pattern)}
		// matches = append(matches, sp.findMatches())
		matches = append(matches, countPossible(pattern, grouping))

		// part 2 - NFA
		b, a := status, text
		bn, an := "", ""
		for i := 0; i < 5; i++ {
			bn, an = bn+b+"?", an+a+","
		}
		b, a = strings.TrimSuffix(bn, "?"), strings.TrimSuffix(an, ",")
		matches2 = append(matches2, countPossible([]rune(b), aoc.SliceMap(aoc.Atoi, strings.Split(a, ",")...)))
	}

	return &result{valuePT1: aoc.Sum(matches...), valuePT2: aoc.Sum(matches2...)}, nil
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
func countPossible(s []rune, c []int) int {
	pos := 0

	cstates := map[state]int{{}: 1} // current state
	nstates := map[state]int{}      // next state

	for len(cstates) > 0 {
		for st, num := range cstates {
			si, ci, cc, expdot := st.springIndex, st.groupIndex, st.continuous, st.expectDot

			// have we reached the end?
			if si == len(s) {
				if ci == len(c) {
					pos += num
				}
				continue
			}

			switch {
			case (s[si] == '#' || s[si] == '?') && ci < len(c) && !expdot:
				// we are still looking for broken springs
				if s[si] == '?' && cc == 0 {
					// we are not in a run of broken springs, so ? can be working
					nstates[state{si + 1, ci, cc, expdot}] += num
				}

				cc++

				if cc == c[ci] {
					// we've found the full next contiguous section of broken springs
					ci++
					cc = 0
					expdot = true // we only want a working spring next
				}

				nstates[state{si + 1, ci, cc, expdot}] += num

			case (s[si] == '.' || s[si] == '?') && cc == 0:
				// we are not in a contiguous run of broken springs
				expdot = false
				nstates[state{si + 1, ci, cc, expdot}] += num
			}
		}

		// swap and clear previous states
		cstates, nstates = nstates, cstates
		clear(nstates)
	}
	return pos
}

type state struct {
	springIndex int
	groupIndex  int
	continuous  int
	expectDot   bool
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
