package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"
)

var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {
	maps := []Map{}
	m := Map{}
	r := &result{}

	for scan.Scan() {
		text := scan.Text()

		if text == "" {
			maps = append(maps, m)
			m = Map{}
		} else {
			m = append(m, []rune(text))
		}
	}
	maps = append(maps, m)

	for _, m := range maps {
		sum, sum2 := findSmudge(m)

		if sum == -1 || sum2 == -1 {
			mr := Map(aoc.Transpose(m))
			Hsum, Hsum2 := findSmudge(mr)

			if sum2 == -1 {
				sum2 = Hsum2 * 100
			}

			if sum == -1 {
				sum = Hsum * 100
			}
		}

		r.valuePT1 += sum
		r.valuePT2 += sum2
	}

	return r, nil
}

type Map [][]rune

func (m Map) String() string {
	var buf strings.Builder
	for i, row := range m {
		if i == 0 {
			fmt.Fprint(&buf, "   ")
			for j := range row {
				fmt.Fprintf(&buf, "%d", j)
			}
			fmt.Fprint(&buf, "\n")
		}

		fmt.Fprintf(&buf, "%d ", i)
		buf.WriteRune(' ')
		buf.WriteString(string(row))
		buf.WriteRune('\n')
	}
	buf.WriteRune('\n')
	return buf.String()

}

// func findReflection(m Map) (int, bool) {
// 	candidates := make(map[int]bool)
// 	var candidateList []int

// 	for _, row := range m {
// 		for col := 1; col < len(row); col++ {
// 			if v, ok := candidates[col]; !ok || v {
// 				candidates[col], _ = reflects(row[:col], row[col:])
// 			}
// 		}
// 		candidateList = all(candidates)
// 		if len(candidateList) == 0 {
// 			return 0, false
// 		}
// 	}

// 	if len(candidateList) == 1 {
// 		return candidateList[0], true
// 	}

// 	return 0, false
// }

type level struct {
	blips  int
	nequal int
	fail   bool
}

func findSmudge(m Map) (int, int) {
	candidates := make(map[int]level)

	for _, row := range m {
		for col := 1; col < len(row); col++ {
			candidate := candidates[col]
			if candidate.fail {
				continue
			}

			eq, bl := reflects(row[:col], row[col:])
			if !eq {
				candidate.nequal++
			}
			candidate.blips += bl

			if candidate.nequal > 1 || candidate.blips > 1 {
				candidate.fail = true
			}

			candidates[col] = candidate
		}
	}

	a, b := -1, -1
	for i, cand := range candidates {
		if !cand.fail && cand.blips == 1 {
			b = i
		}
		if !cand.fail && cand.blips == 0 {
			a = i
		}
	}

	return a, b
}

func reflects(a, b []rune) (bool, int) {
	c := min(len(a), len(b))

	a = append([]rune{}, a...)
	b = append([]rune{}, b...)
	aoc.Reverse(a)
	a = a[:c]
	b = b[:c]

	blips := 0
	for i := range a {
		if a[i] != b[i] {
			blips++
		}
	}

	return blips == 0, blips
}

func all[T comparable](m map[T]bool) []T {
	lis := make([]T, 0, len(m))
	for k, v := range m {
		if v {
			lis = append(lis, k)
		}
	}
	return lis
}
