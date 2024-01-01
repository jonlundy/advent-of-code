package aoc

import (
	"fmt"
	"sort"
)

type priorityQueue[T any, U []T] struct {
	elems        U
	less         func(a, b T) bool
	maxDepth     int
	totalEnqueue int
}

func PriorityQueue[T any, U []T](less func(a, b T) bool) *priorityQueue[T, U] {
	return &priorityQueue[T, U]{less: less}
}
func (pq *priorityQueue[T, U]) Enqueue(elem T) {
	pq.elems = append(pq.elems, elem)
	pq.totalEnqueue++
	pq.maxDepth = max(pq.maxDepth, len(pq.elems))
	sort.Slice(pq.elems, func(i, j int) bool { return pq.less(pq.elems[i], pq.elems[j]) })
}
func (pq *priorityQueue[T, I]) IsEmpty() bool {
	return len(pq.elems) == 0
}
func (pq *priorityQueue[T, I]) Dequeue() (T, bool) {
	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	pq.elems, elem = pq.elems[:len(pq.elems)-1], pq.elems[len(pq.elems)-1]
	return elem, true
}

func ManhattanDistance[T integer](a, b Point[T]) T {
	return ABS(a[1]-b[1]) + ABS(a[0]-b[0])
}

type pather[C number, N any] interface {
	Neighbors(N) []N
	Cost(a, b N) C
	Potential(a, b N) C

	// OPTIONAL:
	// Seen modify value used by seen pruning.
	// Seen(N) N
	// Target returns true if target reached.
	// Target(N) bool
}

type Path[C number, N any] []N

func FindPath[C integer, N comparable](g pather[C, N], start, end N) (C, Path[C, N]) {
	var zero C
	closed := make(map[N]bool)

	type node struct {
		cost      C
		potential C
		parent    *node
		position  N
	}

	NewPath := func(n *node) []N {
		var path []N
		for n.parent != nil {
			path = append(path, n.position)
			n = n.parent
		}
		path = append(path, n.position)

		Reverse(path)
		return path
	}

	less := func(a, b node) bool {
		return b.cost+b.potential < a.cost+a.potential
	}

	pq := PriorityQueue(less)
	pq.Enqueue(node{position: start})

	defer func() {
		fmt.Println("queue max depth = ", pq.maxDepth, "total enqueue = ", pq.totalEnqueue)
	}()

	var seenFn = func(a N) N { return a }
	if s, ok := g.(interface{ Seen(N) N }); ok {
		seenFn = s.Seen
	}

	var targetFn = func(a N) bool { return true }
	if s, ok := g.(interface{ Target(N) bool }); ok {
		targetFn = s.Target
	}

	for !pq.IsEmpty() {
		current, _ := pq.Dequeue()
		cost, potential, n := current.cost, current.potential, current.position

		seen := seenFn(n)
		if closed[seen] {
			continue
		}
		closed[seen] = true

		if cost > 0 && potential == zero && targetFn(current.position) {
			return cost, NewPath(&current)
		}

		for _, nb := range g.Neighbors(n) {
			seen := seenFn(nb)
			if closed[seen] {
				continue
			}

			cost := g.Cost(n, nb) + current.cost
			nextPath := node{
				position:  nb,
				parent:    &current,
				cost:      cost,
				potential: g.Potential(nb, end),
			}
			pq.Enqueue(nextPath)
		}
	}
	return zero, nil
}
