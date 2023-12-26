package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"sort"

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

	result := result{}
	// result.valuePT1 = search(m, 1, 3)
	result.valuePT2 = search(m, 4, 10)

	return &result, nil
}

var (
	ZERO = point{0, 0}

	UP = point{-1, 0}
	DN = point{1, 0}
	LF = point{0, -1}
	RT = point{0, 1}

	INF = int(^uint(0) >> 1)
)

type Map [][]rune

func (m *Map) Get(p point) (point, int, bool) {
	if !m.Valid(p) {
		return [2]int{0, 0}, 0, false
	}

	return p, int((*m)[p[0]][p[1]] - '0'), true
}
func (m *Map) GetNeighbor(p point, d point) (point, int, bool) {
	return m.Get(p.add(d))
}
func (m *Map) Size() (int, int) {
	if m == nil || len(*m) == 0 {
		return 0, 0
	}
	return len(*m), len((*m)[0])
}
func (m *Map) Neighbors(p point) []point {
	var lis []point
	for _, d := range []point{UP, DN, LF, RT} {
		if p, _, ok := m.GetNeighbor(p, d); ok {
			lis = append(lis, p)
		}
	}
	return lis
}
func (m *Map) NeighborDirections(p point) []point {
	var lis []point
	for _, d := range []point{UP, DN, LF, RT} {
		if m.Valid(p.add(d)) {
			lis = append(lis, d)
		}
	}
	return lis
}
func (m *Map) Valid(p point) bool {
	rows, cols := m.Size()
	return p[0] >= 0 && p[0] < rows && p[1] >= 0 && p[1] < cols
}

type memo struct {
	h int
	s int
	p point
	d point
}

func (memo) sort(a, b memo) bool {
	if a.h != b.h {
		return a.h < b.h
	}

	if a.s != b.s {
		return a.s < b.s
	}

	if a.p != b.p {
		return a.p.less(b.p)
	}

	return a.d.less(b.d)
}

type priorityQueue[T any, U []T] struct {
	elems U
	sort  func(a, b T) bool
}

func PriorityQueue[T any, U []T](sort func(a, b T) bool) *priorityQueue[T, U] {
	return &priorityQueue[T, U]{sort: sort}
}
func (pq *priorityQueue[T, U]) Enqueue(elem T) {
	pq.elems = append(pq.elems, elem)
	sort.Slice(pq.elems, func(i, j int) bool { return pq.sort(pq.elems[i], pq.elems[j]) })
}
func (pq *priorityQueue[T, I]) IsEmpty() bool {
	return len(pq.elems) == 0
}
func (pq *priorityQueue[T, I]) Dequeue() (T, bool) {
	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	elem, pq.elems = pq.elems[0], pq.elems[1:]
	return elem, true
}

func heuristic(m Map, p point) int {
	rows, cols := m.Size()
	return rows - p[0] + cols - p[1]
}

func search(m Map, minSize, maxSize int) int {
	rows, cols := m.Size()
	END := point{rows - 1, cols - 1}

	visited := make(map[vector]int)
	pq := PriorityQueue(memo{}.sort)
	pq.Enqueue(memo{h: heuristic(m, point{0, 0}), p: point{0, 0}, d: DN})
	
	for !pq.IsEmpty() {
		mem, _ := pq.Dequeue()
		fmt.Println(mem)
		if mem.h > dmap(visited, vector{mem.p[0], mem.p[1], mem.d[0], mem.d[1]}, INF) {
			continue
		}

		if mem.p == END {
			return mem.s
		}

		for _, nd := range m.NeighborDirections(mem.p) {
			if nd[0] == 0 && mem.d == RT || nd[1] == 0 && mem.d == DN {
				continue
			}

			dscore := 0

			for _, size := range irange(1, maxSize+1) {
				np := mem.p.add(nd.scale(size))
				_, s, ok := m.Get(np)

				if !ok {
					break
				}

				dscore += s
				pscore := mem.s + dscore

				nh := heuristic(m, np) + pscore
				vec := vector{np[0], np[1], nd[0], nd[1]}

				if size >= minSize && nh < dmap(visited, vec, INF) {
					pq.Enqueue(memo{nh, pscore, np, nd})
					visited[vec] = nh
				}
			}
		}
	}

	return 0
}

func dmap[K comparable, V any](m map[K]V, k K, d V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}
func irange(a, b int) []int {
	lis := make([]int, b-a)
	for i := range lis {
		lis[i] = i + a
	}
	return lis
}

type point [2]int

func (p point) add(a point) point { return point{p[0] + a[0], p[1] + a[1]} }
func (p point) scale(m int) point { return point{p[0] * m, p[1] * m} }
func (p point) less(a point) bool { return p[0] < a[0] || p[1] < a[1] }

type vector [4]int
