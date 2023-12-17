package main

import (
	"bufio"
	_ "embed"
	"fmt"

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
	var m Map

	for scan.Scan() {
		text := scan.Text()
		m = append(m, []rune(text))
	}

	rows := len(m)
	cols := len(m[0])

	options := make([]int, 2*(rows+cols)+2)
	i := 0
	for j:=0; j<=rows-1; j++ {
		options[i+0] = runCycle(m, ray{[2]int{j, -1}, RT})
		options[i+1] = runCycle(m, ray{[2]int{j, cols}, LF})
		i+=2
	}
	for j:=0; j<=cols-1; j++ {
		options[i+0] = runCycle(m, ray{[2]int{-1, j}, DN})
		options[i+1] = runCycle(m, ray{[2]int{rows, j}, UP})
		i+=2
	}

	// fmt.Println(options)
	return &result{valuePT1: options[0], valuePT2: aoc.Max(options[0], options[1:]...)}, nil
}

type stack[T any] []T

func (s *stack[T]) Push(v T) {
	if s == nil {
		panic("nil stack")
	}
	*s = append(*s, v)
}
func (s *stack[T]) Pop() T {
	if s == nil || len(*s) == 0 {
		panic("empty stack")
	}
	defer func() { *s = (*s)[:len(*s)-1] }()
	return (*s)[len(*s)-1]
}

var (
	ZERO = [2]int{0, -1}

	UP = [2]int{-1, 0}
	DN = [2]int{1, 0}
	LF = [2]int{0, -1}
	RT = [2]int{0, 1}
)

type ray struct {
	pos [2]int
	dir [2]int
}

func (r *ray) next() [2]int {
	r.pos[0] += r.dir[0]
	r.pos[1] += r.dir[1]
	return r.pos
}

type Map [][]rune

func (m *Map) Get(p [2]int) rune {
	if p[0] < 0 || p[0] >= len((*m)) {
		return 0
	}
	if p[1] < 0 || p[1] >= len((*m)[0]) {
		return 0
	}

	return (*m)[p[0]][p[1]]
}


func runCycle(m Map, r ray) int {
	current := r

	s := stack[ray]{}
	s.Push(current)

	energized := make(map[[2]int]bool)
	// energized[current.pos] = true

	cycle := make(map[[4]int]bool)

	for len(s) > 0 {
		current = s.Pop()

		r := m.Get(current.next())
		// fmt.Println("pos", current.pos, current.dir, string(r), len(s))
		if r == 0 {
			continue
		}
		energized[current.pos] = true
		v := [4]int{
			current.pos[0],
			current.pos[1],
			current.dir[0],
			current.dir[1],
		}

		if _, ok := cycle[v]; ok {
			// fmt.Println("cycle")
			continue
		}
		cycle[v] = true

		switch r {
		case '|':
			switch current.dir {
			case UP, DN:
				// pass
			case LF, RT:
				current.dir = UP
				s.Push(ray{current.pos, DN})
			}
		case '-':
			switch current.dir {
			case LF, RT:
				// pass
			case UP, DN:
				current.dir = LF
				s.Push(ray{current.pos, RT})
			}
		case '/':
			switch current.dir {
			case UP:
				current.dir = RT
			case DN:
				current.dir = LF
			case LF:
				current.dir = DN
			case RT:
				current.dir = UP
			}
		case '\\':
			switch current.dir {
			case UP:
				current.dir = LF
			case DN:
				current.dir = RT
			case LF:
				current.dir = UP
			case RT:
				current.dir = DN
			}
		}

		s.Push(current)
	}

	// for i := range m {
	// 	for j := range m[i] {
	// 		if v := energized[[2]int{i,j}]; v {
	// 			fmt.Print("#")
	// 		} else {
	// 			fmt.Print(".")
	// 		}
	// 	}
	// 	fmt.Println("")
	// }

	return len(energized)
}