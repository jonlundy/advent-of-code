package main

import (
	"bufio"
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
	m := &Path{s: -1}
	for scan.Scan() {
		text := scan.Text()
		_ = text

		m.readLine(text)
	}
	dist := m.buildPath()
	// log(m)

	r := &Region{List: m.List, w: m.w, l: len(m.m)}
	count := r.Count()

	return &result{valuePT1: dist, valuePT2: count}, nil
}

type node struct {
	value  rune
	pos    int
	whence int8
	left   *node
}

func (n *node) add(a *node) *node {
	if a == nil {
		return n
	}
	n.left = a
	return a
}

func (n *node) String() string {
	if n == nil {
		return "EOL"
	}
	return fmt.Sprintf("node %s from %s", string(n.value), enum(n.whence))
}

type List struct {
	head *node
	n    *node
	p    map[int]*node
}

func NewList(a *node) *List {
	lis := &List{
		head: a,
		n:    a,
		p:    make(map[int]*node),
	}
	lis.add(a)

	return lis
}
func (l *List) add(a *node) {
	l.n = l.n.add(a)
	l.p[a.pos] = a
}

type Path struct {
	m []rune
	w int
	s int

	*List
}

func (m *Path) String() string {
	var buf strings.Builder
	n := m.head

	buf.WriteString(fmt.Sprintf("head %d", len(m.p)))
	for n != nil {
		buf.WriteString("\n ")
		buf.WriteString(n.String())
		n = n.left
	}
	return buf.String()
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
	for m.next() {
	}
	return (len(m.p) + 1) / 2
}
func (m *Path) start() {
	m.List = NewList(&node{value: 'S', pos: m.s})

	switch {
	case m.peek(UP) != nil:
		m.add(m.peek(UP))
	case m.peek(DN) != nil:
		m.add(m.peek(DN))
	case m.peek(LF) != nil:
		m.add(m.peek(LF))
	case m.peek(RT) != nil:
		m.add(m.peek(RT))
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
	if n.value == 'S' {
		last := n.whence
		next := m.head.left.whence

		switch (last<<4)|next {
		case 0x11, 0x22: m.n.value = '|' // UP UP, DN DN

		case 0x13: m.head.value = 'J' // UP LF
		case 0x14: m.head.value = 'L' // UP RT

		case 0x23: m.head.value = '7' // DN LF
		case 0x24: m.head.value = 'F' // DN RT

		case 0x33, 0x44: m.head.value = '-' // LF LF, RT RT

		case 0x31: m.head.value = '7' // LF UP 
		case 0x32: m.head.value = 'J' // LF DN 

		case 0x41: m.head.value = 'F' // RT UP
		case 0x42: m.head.value = 'L' // RT DN 

		}

		return false
	}

	m.add(n)


	return true
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
		if any(r, '7', '|', 'F', 'S') {
			return &node{value: r, whence: DN, pos: p}
		}
	case DN:
		x, y := toXY(m.n.pos, m.w)
		if y == m.w {
			return nil
		}

		p := fromXY(x, y+1, m.w)
		r := m.m[p]
		if any(r, 'J', '|', 'L', 'S') {
			return &node{value: r, whence: UP, pos: p}
		}

	case LF:
		x, y := toXY(m.n.pos, m.w)
		if x == 0 {
			return nil
		}

		p := fromXY(x-1, y, m.w)
		r := m.m[p]
		if any(r, 'F', '-', 'L', 'S') {
			return &node{value: r, whence: RT, pos: p}
		}

	case RT:
		x, y := toXY(m.n.pos, m.w)
		if x == m.w {
			return nil
		}
		p := fromXY(x+1, y, m.w)
		r := m.m[p]
		if any(r, '7', '-', 'J', 'S') {
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

type Region struct {
	*List
	inLoop bool
	count  int
	w      int
	l      int
}

func (r *Region) Count() int {
	for i := 0; i < r.l; i++ {
		if i%r.w == 0 {
			r.inLoop = false
			// fmt.Println(": ", i)
		}

		a, ok := r.p[i]
		if ok && any(a.value, '|', '7', 'F', 'X') {
			r.inLoop = !r.inLoop
			// fmt.Print(string(a.value))
			continue
		}

		if !ok && r.inLoop {
			// fmt.Print("I")
			r.count++
			continue
		}

		if ok {
			// fmt.Print(string(a.value))
			continue
		}
		// fmt.Print(".")

	}
	return r.count
}
