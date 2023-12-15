package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
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
	r := &result{}

	var ops []string
	for scan.Scan() {
		text := scan.Text()
		ops = strings.Split(text, ",")

		r.valuePT1 = aoc.Reduce(func(i int, t string, sum int) int {
			sum += hash(t)
			return sum
		}, 0, ops...)
	}

	var boxen boxes

	boxen = aoc.Reduce(func(i int, op string, b boxes) boxes {
		return b.Op(op)
	}, boxen, ops...)

	r.valuePT2 = boxen.Sum()

	log(boxen)

	return r, nil
}

func hash(s string) int {
	var sum int
	for _, a := range s {
		sum += int(a)
		sum *= 17
		sum %= 256
	}

	return sum
}

type lens struct {
	label string
	value int
}

func (l lens) String() string {
	return fmt.Sprintf("[%s %d]", l.label, l.value)
}

type box []lens

func (lis box) String() string {
	var buf strings.Builder
	if len(lis) > 0 {
		buf.WriteString(lis[0].String())
	}
	for _, l := range lis[1:] {
		buf.WriteString(l.String())
	}
	return buf.String()
}

type boxes [256]box

func (lis boxes) String() string {
	var buf strings.Builder
	buf.WriteString("Boxes:\n")
	for i, b := range lis {
		if len(b) > 0 {
			fmt.Fprintf(&buf, "Box %d: %v\n", i, b)
		}
	}
	return buf.String()
}

func (lis boxes) Op(op string) boxes {
	if a, _, ok := strings.Cut(op, "-"); ok {
		i := hash(a)

		pos := slices.IndexFunc(lis[i], func(l lens) bool { return l.label == a })

		if pos >= 0 {
			lis[i] = append(lis[i][:pos], lis[i][pos+1:]...)
		}
	} else if a, b, ok := strings.Cut(op, "="); ok {
		i := hash(a)
		v := aoc.Atoi(b)

		pos := slices.IndexFunc(lis[i], func(l lens) bool { return l.label == a })

		if pos == -1 {
			lis[i] = append(lis[i], lens{a, v})
		} else {
			lis[i][pos].value = v
		}
	}

	return lis
}
func (lis boxes) Sum() int {
	// return aoc.Reduce(func(b int, box box, sum int) int {
	// 	return aoc.Reduce(
	// 		func(s int, lens lens, sum int) int {
	// 			return sum + (b+1)*(s+1)*lens.value
	// 		}, sum, box...)
	// 	}, 0, lis[:]...)

	var sum int
	for b := range lis {
		for s := range lis[b] {
			sum += (b+1) * (s+1) * lis[b][s].value
		}
	}
	return sum
}
