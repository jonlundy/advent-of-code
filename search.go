package aoc

import (
	"sort"
)

type PQElem[T any, I integer] struct {
	Value    T
	Priority I
}
type PQList[T any, I integer] []PQElem[T, I]

func (pq PQList[T, I]) Len() int {
	return len(pq)
}
func (pq PQList[T, I]) Less(i int, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PQList[T, I]) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

var _ sort.Interface = (*PQList[rune, int])(nil)

type PriorityQueue[T any, I integer] struct {
	elem PQList[T, I]
}

func (pq *PriorityQueue[T, I]) Enqueue(elem T, priority I) {
	pq.elem = append(pq.elem, PQElem[T, I]{elem, priority})
	sort.Sort(pq.elem)
}
func (pq *PriorityQueue[T, I]) IsEmpty() bool {
	return len(pq.elem) == 0
}
func (pq *PriorityQueue[T, I]) Dequeue() (T, bool) {
	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	elem, pq.elem = pq.elem[0].Value, pq.elem[1:]
	return elem, true
}

type Vertex[V comparable, I integer] struct {
	to    V
	score I
}
type graph[V comparable, I uinteger] struct {
	adj map[V][]Vertex[V, I]
}

func Graph[V comparable, I uinteger](size int) *graph[V, I] {
	return &graph[V, I]{
		adj: make(map[V][]Vertex[V, I], size),
	}
}
func (g *graph[V, I]) AddEdge(u, v V, w I) {
	g.adj[u] = append(g.adj[u], Vertex[V, I]{to: v, score: w})
	g.adj[v] = append(g.adj[v], Vertex[V, I]{to: u, score: w})
}
func (g *graph[V, I]) Dijkstra(m interface{ Get() }, src V) map[V]I {
	pq := PriorityQueue[V, I]{}
	dist := make(map[V]I, len(g.adj))
	visited := make(map[V]bool, len(g.adj))
	var INF I
	INF = ^INF

	pq.Enqueue(src, 0)
	dist[src] = 0

	for !pq.IsEmpty() {
		u, _ := pq.Dequeue()

		if _, ok := visited[u]; ok {
			continue
		}
		visited[u] = true

		for _, v := range g.adj[u] {
			_, ok := visited[v.to]
			var du, dv I
			if d, inf := dist[u]; !inf {
				du = INF
			} else {
				du = d
			}
			if d, inf := dist[v.to]; !inf {
				dv = INF
			} else {
				dv = d
			}

			if !ok && du+v.score < dv {
				dist[v.to] = du + v.score
				pq.Enqueue(v.to, du+v.score)
			}
		}
	}

	return dist
}
