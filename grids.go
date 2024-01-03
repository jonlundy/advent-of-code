package aoc

import (
	"cmp"
	"sort"
)

type Vector struct {
	Offset Point[int]
	Scale  int
}

func (v Vector) Point() Point[int] {
	return v.Offset.Scale(v.Scale)
}

type Point[T integer] [2]T

func (p Point[T]) Add(a Point[T]) Point[T] {
	return Point[T]{p[0] + a[0], p[1] + a[1]}
}
func (p Point[T]) Scale(m T) Point[T] {
	return Point[T]{p[0] * m, p[1] * m}
}
func (p Point[T]) Less(b Point[T]) bool {
	if p[0] != b[0] {
		return p[0] < b[0]
	}
	return p[1] < b[1]
}

func Transpose[T any](matrix [][]T) [][]T {
	rows, cols := len(matrix), len(matrix[0])

	m := make([][]T, cols)
	for i := range m {
		m[i] = make([]T, rows)
	}

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			m[i][j] = matrix[j][i]
		}
	}
	return m
}

// NumPoints the number of the points inside an outline plus the number of points in the outline
func NumPoints(outline []Point[int], borderLength int) int {
	// shoelace - find the float area in a shape
	sum := 0
	for _, p := range Pairwise(outline) {
		row1, col1 := p[0][0], p[0][1]
		row2, col2 := p[1][0], p[1][1]

		sum += row1*col2 - row2*col1
	}
	area := sum / 2

	// pick's theorem - find the number of points in a shape given its area
	return (ABS(area) - borderLength/2 + 1) + borderLength
}

type Map[I integer, T any] [][]T

func (m *Map[I, T]) Get(p Point[I]) (Point[I], T, bool) {
	var zero T
	if !m.Valid(p) {
		return [2]I{0, 0}, zero, false
	}

	return p, (*m)[p[0]][p[1]], true
}
func (m *Map[I, T]) Size() (I, I) {
	if m == nil || len(*m) == 0 {
		return 0, 0
	}
	return I(len(*m)), I(len((*m)[0]))
}
func (m *Map[I, T]) Valid(p Point[I]) bool {
	rows, cols := m.Size()
	return p[0] >= 0 && p[0] < rows && p[1] >= 0 && p[1] < cols
}

type adjacencyList[V any, C comparable] map[C][]V
type graph[V any, W cmp.Ordered, C comparable] map[C]*vertex[V, W]
type graphOption[V any, W cmp.Ordered, C comparable] func(g *graph[V, W, C])
type vertex[V any, W cmp.Ordered] struct {
	Value V
	Edges edges[V, W]
}

func (v *vertex[V, W]) Neighbors() []V {
	var nbs []V
	sort.Sort(v.Edges)
	for _, e := range v.Edges {
		nbs = append(nbs, e.Vertex.Value)
	}
	return nbs
}

type edge[V any, W cmp.Ordered] struct {
	Vertex *vertex[V, W]
	Weight W
}
type edges[V any, W cmp.Ordered] []edge[V, W]

func (e edges[V, W]) Len() int           { return len(e) }
func (e edges[V, W]) Less(i, j int) bool { return e[i].Weight < e[j].Weight }
func (e edges[V, W]) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func Graph[V any, W cmp.Ordered, C comparable](opts ...graphOption[V, W, C]) *graph[V, W, C] {
	g := make(graph[V, W, C])
	for _, opt := range opts {
		opt(&g)
	}
	return &g
}
func (g *graph[V, W, C]) AddVertex(id C, value V) {
	(*g)[id] = &vertex[V,W]{Value: value}
}
func (g *graph[V, W, C]) AddEdge(from, to C, w W) {
	if g == nil {
		return
	}
	if _, ok := (*g)[from]; !ok {
		return
	}
	if _, ok := (*g)[to]; !ok {
		return
	}

	(*g)[from].Edges = append((*g)[from].Edges, edge[V,W]{(*g)[to], w})
}
func (g *graph[V, W, C]) Neighbors(v C) []V {
	if g == nil {
		return nil
	}

	return (*g)[v].Neighbors()
}
func (g *graph[V, W, C]) AdjacencyList() adjacencyList[V, C] {
	m := make(map[C][]V)
	for id, v := range *g {
		if len(v.Edges) == 0 {
			continue
		}
		m[id] = v.Neighbors()
	}
	return m
}

func WithAdjacencyList[W cmp.Ordered, C comparable](list adjacencyList[C, C]) graphOption[C, W, C] {
	var zeroW W
	return func(g *graph[C, W, C]) {
		for vertex, edges := range list {
			if _, ok := (*g)[vertex]; !ok {
				g.AddVertex(vertex, vertex)
			}

			// add edges to vertex
			for _, edge := range edges {
				// add edge as vertex, if not added
				if _, ok := (*g)[edge]; !ok {
					g.AddVertex(edge, edge)
				}

				g.AddEdge(vertex, edge, zeroW) // no weights in this adjacency list
			}
		}
	}
}

// func GraphFromMap()
