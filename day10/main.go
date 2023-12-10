package main

import (
	"bufio"
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
	m := Path{s: -1}
	for scan.Scan() {
		text := scan.Text()
		_ = text

		m.readLine(text)
	}
	dist := m.buildPath()



	return &result{valuePT1: dist}, nil
}

type node struct {
	value       rune
	pos         int
	whence      int8
	left, right *node
}

func (n *node) add(a *node) *node {
	if a == nil {
		return n
	}
	n.left, a.right = a, n
	return a
}

func (n *node) String() string {
	return fmt.Sprintf("node %s from %s", string(n.value), enum(n.whence))
}

type Path struct {
	m []rune
	w int
	s int
	l int
	n *node
}

func (m *Path) readLine(text string) {
	if m.w == 0 {
		m.w = len(text)
	}
	if m.s == -1 {
		if i := strings.IndexRune(text, 'S'); i != -1 {
			m.s = i + len(m.m)
		}
	}

	m.m = append(m.m, []rune(text)...)
}
func (m *Path) buildPath() int {
	m.start()
	// log(m.n)
	i := 0
	for m.next() {
		// log(m.n)
		i++
	}
	// log(m.n)
	// log("length", (m.l+2)/2)
	return (m.l+2)/2
}
func (m *Path) start() {
	m.n = &node{value: 'S', pos: m.s}
	// log(m.n)

	switch {
	case m.peek(UP) != nil:
		m.n = m.n.add(m.peek(UP))
	case m.peek(DN) != nil:
		m.n = m.n.add(m.peek(DN))
	case m.peek(LF) != nil:
		m.n = m.n.add(m.peek(LF))
	case m.peek(RT) != nil:
		m.n = m.n.add(m.peek(RT))
	}
}
func (m *Path) next() bool {
	var n *node
	switch m.n.value {
	case '7':
		if m.n.whence == LF {
			n = m.peek(DN)
		} else {
			n = m.peek(LF)
		}
	case '|':
		if m.n.whence == UP {
			n = m.peek(DN)
		} else {
			n = m.peek(UP)
		}
	case 'F':
		if m.n.whence == RT {
			n = m.peek(DN)
		} else {
			n = m.peek(RT)
		}
	case '-':
		if m.n.whence == LF {
			n = m.peek(RT)
		} else {
			n = m.peek(LF)
		}
	case 'J':
		if m.n.whence == LF {
			n = m.peek(UP)
		} else {
			n = m.peek(LF)
		}
	case 'L':
		if m.n.whence == RT {
			n = m.peek(UP)
		} else {
			n = m.peek(RT)
		}
	}
	if n == nil {
		return false
	}

	m.n = m.n.add(n)
	m.l++
	return m.n.value != 'S'
}

const (
	ST int8 = iota
	UP
	DN
	LF
	RT
)

func enum(e int8) string {
	switch e {
	case ST:
		return "ST"
	case UP:
		return "UP"
	case DN:
		return "DN"
	case LF:
		return "LF"
	case RT:
		return "RT"
	default:
		return "XX"
	}
}

func (m *Path) peek(d int8) *node {
	switch d {
	case UP:
		x, y := toXY(m.n.pos, m.w)
		if y == 0 {
			return nil
		}

		p := fromXY(x, y-1, m.w)
		r := m.m[p]
		if any(r, '7', '|', 'F') {
			return &node{value: r, whence: DN, pos: p}
		}
	case DN:
		x, y := toXY(m.n.pos, m.w)
		if y == m.w {
			return nil
		}

		p := fromXY(x, y+1, m.w)
		r := m.m[p]
		if any(r, 'J', '|', 'L') {
			return &node{value: r, whence: UP, pos: p}
		}

	case LF:
		x, y := toXY(m.n.pos, m.w)
		if x == 0 {
			return nil
		}

		p := fromXY(x-1, y, m.w)
		r := m.m[p]
		if any(r, 'F', '-', 'L') {
			return &node{value: r, whence: RT, pos: p}
		}

	case RT:
		x, y := toXY(m.n.pos, m.w)
		if x == m.w {
			return nil
		}
		p := fromXY(x+1, y, m.w)
		r := m.m[p]
		if any(r, '7', '-', 'J') {
			return &node{value: r, whence: LF, pos: p}
		}

	}
	return nil
}

func fromXY(x, y, w int) int { return y*w + x }
func toXY(i, w int) (int, int) {
	return i % w, i / w
}
func any[T comparable](n T, stack ...T) bool {
	var found bool
	for _, s := range stack {
		if n == s {
			found = true
			break
		}
	}
	return found
}
