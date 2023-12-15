package main

import (
	"bufio"
	_ "embed"
	"fmt"
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
	var maps []Map
	var m Map

	for scan.Scan() {
		text := scan.Text()
		if len(text) == 0 {
			maps = append(maps, m)
			m = Map{}
		}

		m = append(m, []rune(text))
	}
	maps = append(maps, m)

	score1 := 0
	for _, m := range maps {
		m = aoc.Transpose(reverse(m))
		m.Sort()
		score1 += m.Score()
	}

	score2 := 0
	// type record [5]int
	// var last, current record
	// memo := make(map[record]int)

	// for _, m := range maps {
	// 	fmt.Println(m)

	// 	m = aoc.Transpose(reverse(m))

	// 	for i := 0; i < 1_000_000_000; i++ {
	// 		m.Sort()
	// 		current[0] = m.Score()

	// 		m = aoc.Transpose(reverse(m))
	// 		m.Sort()
	// 		current[1] = m.Score()

	// 		m = aoc.Transpose(reverse(m))
	// 		m.Sort()
	// 		current[2] = m.Score()

	// 		m = aoc.Transpose(reverse(m))
	// 		m.Sort()
	// 		current[3] = m.Score()

	// 		m = aoc.Transpose(reverse(m))
	// 		current[4] = m.Score()

	// 		if slices.Equal(last[:], current[:]) {
	// 			break
	// 		}

	// 		v, ok := memo[current]
	// 		if (i - 104) % 13 == 12 {
	// 			fmt.Println(i, current, ok, v)
	// 		}
	// 		// if ok && v > 2 {
	// 		// 	break
	// 		// }

	// 		memo[current] += 1

	// 		last = current
	// 	}

	// 	score2 += m.Score()
	// }

	return &result{valuePT1: score1, valuePT2: score2}, nil
}

type Map [][]rune

func (m Map) String() string {
	var buf strings.Builder
	for _, row := range m {
		buf.WriteString(string(row))
		buf.WriteRune('\n')
	}
	buf.WriteRune('\n')

	return buf.String()
}
func (m *Map) Sort() {
	if m == nil {
		return
	}
	for _, row := range *m {
		base := 0
		for i, r := range row {
			if r == '#' {
				base = i + 1
				continue
			}
			if r == 'O' {
				if base < i {
					row[base], row[i] = row[i], row[base]
				}
				base++
			}
		}
	}
}
func (m *Map) Score() int {
	if m == nil {
		return 0
	}

	sum := 0
	max := len(*m)
	for _, row := range *m {
		for i, r := range row {
			if r == 'O' {
				sum += max - i
			}
		}
	}
	return sum
}

func reverse(m Map) Map {
	for _, row := range m {
		aoc.Reverse(row)
	}
	return m
}
