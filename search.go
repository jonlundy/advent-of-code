package aoc

import (
	"sort"
)

type priorityQueue[T any] struct {
	elems        []T
	less         func(a, b T) bool
	maxDepth     int
	totalEnqueue int
	totalDequeue int
}

// PriorityQueue implements a simple slice based queue.
// less is the function for sorting. reverse a and b to reverse the sort.
// T is the item
// U is a slice of T
func PriorityQueue[T any](less func(a, b T) bool) *priorityQueue[T] {
	return &priorityQueue[T]{less: less}
}
func (pq *priorityQueue[T]) Enqueue(elem T) {
	pq.totalEnqueue++

	pq.elems = append(pq.elems, elem)
	pq.maxDepth = max(pq.maxDepth, len(pq.elems))
}
func (pq *priorityQueue[T]) IsEmpty() bool {
	return len(pq.elems) == 0
}
func (pq *priorityQueue[T]) Dequeue() (T, bool) {
	pq.totalDequeue++

	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	sort.Slice(pq.elems, func(i, j int) bool { return pq.less(pq.elems[i], pq.elems[j]) })
	pq.elems, elem = pq.elems[:len(pq.elems)-1], pq.elems[len(pq.elems)-1]
	return elem, true
}

// ManhattanDistance the distance between two points measured along axes at right angles.
func ManhattanDistance[T integer](a, b Point[T]) T {
	return ABS(a[1]-b[1]) + ABS(a[0]-b[0])
}

type pather[C number, N comparable] interface {
	Neighbors(N) []N
	Cost(a, b N) C
	Potential(a, b N) C

	// OPTIONAL:
	// Seen modify value used by seen pruning.
	// Seen(N) N

	// Target returns true if target reached.
	// Target(N) bool
}

// FindPath uses the A* path finding algorithem.
// g is the graph source that implements the pather interface.
//   C is an numeric type for calculating cost/potential
//   N is the node values. is comparable for storing in visited table for pruning.
// start, end are nodes that dileniate the start and end of the search path.
// The returned values are the calculated cost and the path taken from start to end.
func FindPath[C integer, N comparable](g pather[C, N], start, end N) (C, []N) {
	var zero C
	visited := make(map[N]bool)

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
		Log("queue max depth = ", pq.maxDepth, "total enqueue = ", pq.totalEnqueue, "total dequeue = ", pq.totalDequeue)
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
		if visited[seen] {
			continue
		}
		visited[seen] = true

		if cost > 0 && potential == zero && targetFn(current.position) {
			return cost, NewPath(&current)
		}

		for _, nb := range g.Neighbors(n) {
			seen := seenFn(nb)
			if visited[seen] {
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
