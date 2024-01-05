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

type stack[T any] []T

func Stack[T any](a ...T) *stack[T] {
	var s stack[T] = a
	return &s
}
func (s *stack[T]) Push(a ...T) {
	if s == nil {
		return
	}
	*s = append(*s, a...)
}
func (s *stack[T]) IsEmpty() bool {
	return s == nil || len(*s) == 0
}
func (s *stack[T]) Pop() T {
	var a T
	if s.IsEmpty() {
		return a
	}
	a, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return a
}

// ManhattanDistance the distance between two points measured along axes at right angles.
func ManhattanDistance[T integer](a, b Point[T]) T {
	return ABS(a[0]-b[0]) + ABS(a[1]-b[1])
}

type pather[C number, N comparable] interface {
	// Neighbors returns all neighbors to node N that should be considered next.
	Neighbors(N) []N

	// Cost returns 
	Cost(a, b N) C

	// Target returns true when target reached. receives node and cost.
	Target(N, C) bool

	// OPTIONAL:
	// Add heuristic for running as A* search.
	// Potential(N) C

	// Seen modify value used by seen pruning.
	// Seen(N) N

}

// FindPath uses the A* path finding algorithem.
// g is the graph source that implements the pather interface.
//
//	C is an numeric type for calculating cost/potential
//	N is the node values. is comparable for storing in visited table for pruning.
//
// start, end are nodes that dileniate the start and end of the search path.
// The returned values are the calculated cost and the path taken from start to end.
func FindPath[C integer, N comparable](g pather[C, N], start, end N) (C, []N, map[N]C) {
	var zero C

	var seenFn = func(a N) N { return a }
	if s, ok := g.(interface{ Seen(N) N }); ok {
		seenFn = s.Seen
	}

	var potentialFn = func(N) C { var zero C; return zero }
	if p, ok := g.(interface{ Potential(N) C }); ok {
		potentialFn = p.Potential
	}

	type node struct {
		cost      C
		potential C
		parent    *node
		position  N
	}

	newPath := func(n *node) []N {
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
		return  b.cost+b.potential < a.cost+a.potential
	}

	closed := make(map[N]C)
	open := PriorityQueue(less)

	open.Enqueue(node{position: start, potential: potentialFn(start)})
	closed[start] = zero

	// defer func() {
	// 	Log(
	// 		"queue max depth = ", open.maxDepth, 
	// 		"total enqueue = ", open.totalEnqueue, 
	// 		"total dequeue = ", open.totalDequeue,
	// 		"total closed = ", len(closed),
	// 	)
	// }()

	for !open.IsEmpty() {
		current, _ := open.Dequeue()
		for _, nb := range g.Neighbors(current.position) {
			next := node{
				position:  nb,
				parent:    &current,
				cost:      g.Cost(current.position, nb) + current.cost,
				potential: potentialFn(nb),
			}

			seen := seenFn(nb)
			cost, ok := closed[seen]
			if !ok || next.cost < cost {
				open.Enqueue(next)
				closed[seen] = next.cost
			}	

			if next.potential == zero && g.Target(next.position, next.cost) {
				return next.cost, newPath(&next), closed
			}
		}
	}
	return zero, nil, closed
}
