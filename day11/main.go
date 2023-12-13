package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"sort"
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
	m := NewMap()

	for scan.Scan() {
		text := scan.Text()
		m.readLine(text)
	}

	return &result{
		valuePT1: m.expand(1).sumPaths(),
		valuePT2: m.expand(999_999).sumPaths(),
	}, nil
}

type Map struct {
	rows int
	cols int

	emptyRows map[int]bool
	emptyCols map[int]bool

	*aoc.List[rune]
}

func NewMap() *Map {
	return &Map{
		emptyRows: make(map[int]bool),
		emptyCols: make(map[int]bool),
		List:      aoc.NewList[rune](nil),
	}
}
func (m *Map) String() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "Map size %d x %d\n", m.rows, m.cols)
	fmt.Fprintln(buf, "empty rows:", all(m.emptyRows))
	fmt.Fprintln(buf, "empty cols:", all(m.emptyCols))

	n := m.Head()
	for n != nil {
		fmt.Fprintln(buf, toXY(n.Position(), m.cols), n.String())
		n = n.Next()
	}

	for row := 0; row < m.rows; row++ {
		for col := 0; col < m.cols; col++ {
			if n := m.getRC(row, col); n != nil {
				buf.WriteRune('#')
			} else {
				buf.WriteRune('.')
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func (m *Map) readLine(text string) {
	if m.cols == 0 {
		m.cols = len(text)
	}

	emptyRow, ok := m.emptyRows[m.rows]
	if !ok {
		emptyRow = true
	}

	row := []rune(text)
	for col, r := range row {
		emptyCol, ok := m.emptyCols[col]
		if !ok {
			emptyCol = true
		}

		if r == '#' {
			m.Add(r, fromXY(col, m.rows, m.cols))
			emptyCol = false
			emptyRow = false
		}

		m.emptyRows[m.rows] = emptyRow
		m.emptyCols[col] = emptyCol
	}

	m.rows++
}
func (m *Map) getRC(row, col int) *aoc.Node[rune] {
	return m.List.Get(fromXY(col, row, m.cols))
}
func (m *Map) expand(rate int) *Map {
	newM := NewMap()

	newM.rows = m.rows + rate*len(all(m.emptyRows))
	newM.cols = m.cols + rate*len(all(m.emptyCols))

	offsetC := 0
	for col := 0; col < m.cols; col++ {
		if empty, ok := m.emptyCols[col]; ok && empty {
			for r := 0; r <= rate; r++ {
				newM.emptyCols[offsetC+col+r] = true
			}
			offsetC += rate
			continue
		}
	}

	offsetR := 0
	for row := 0; row < m.rows; row++ {
		if empty, ok := m.emptyRows[row]; ok && empty {
			for r := 0; r <= rate; r++ {
				newM.emptyRows[offsetR+row+r] = true
			}
			offsetR += rate
			continue
		}

		offsetC := 0
		for col := 0; col < m.cols; col++ {
			if empty, ok := m.emptyCols[col]; ok && empty {
				offsetC += rate

				continue
			}

			if n := m.getRC(row, col); n != nil {
				newM.Add('#', fromXY(offsetC+col, offsetR+row, newM.cols))
			}
		}
	}

	return newM
}
func (m *Map) sumPaths() int {
	var positions []int

	n := m.Head()
	for n != nil {
		positions = append(positions, n.Position())
		n = n.Next()
	}

	var paths []int

	for i := 0; i < len(positions); i++ {
		p := positions[i]
		pXY := toXY(p, m.cols)

		for j := i; j < len(positions); j++ {
			c := positions[j]
			if c == p {
				continue
			}

			cXY := toXY(c, m.cols)

			path := aoc.ABS(cXY[0]-pXY[0]) + aoc.ABS(cXY[1]-pXY[1])
			paths = append(paths, path)
		}
	}
	return aoc.Sum(paths...)
}

func all(m map[int]bool) []int {
	lis := make([]int, 0, len(m))
	for k, v := range m {
		if v {
			lis = append(lis, k)
		}
	}
	sort.Ints(lis)
	return lis
}
func fromXY(x, y, w int) int { return y*w + x }
func toXY(i, w int) [2]int    { return [2]int{i % w, i / w} }
