package aoc

import (
	"math/bits"
	"sort"
)

type priorityQueue[T any] struct {
	elems        []*T
	less         func(a, b *T) bool
	maxDepth     int
	totalEnqueue int
	totalDequeue int
}

// PriorityQueue implements a simple slice based queue.
// less is the function for sorting. reverse a and b to reverse the sort.
// T is the item
// U is a slice of T
func PriorityQueue[T any](less func(a, b *T) bool) *priorityQueue[T] {
	return &priorityQueue[T]{less: less}
}
func (pq *priorityQueue[T]) Insert(elem *T) {
	pq.totalEnqueue++

	pq.elems = append(pq.elems, elem)
	pq.maxDepth = max(pq.maxDepth, len(pq.elems))
}
func (pq *priorityQueue[T]) IsEmpty() bool {
	return len(pq.elems) == 0
}
func (pq *priorityQueue[T]) ExtractMin() *T {
	pq.totalDequeue++

	var elem *T
	if pq.IsEmpty() {
		return elem
	}

	sort.Slice(pq.elems, func(i, j int) bool { return pq.less(pq.elems[j], pq.elems[i]) })
	pq.elems, elem = pq.elems[:len(pq.elems)-1], pq.elems[len(pq.elems)-1]
	return elem
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

	less := func(a, b *node) bool {
		return a.cost+a.potential < b.cost+b.potential
	}

	closed := make(map[N]C)
	open := FibHeap(less)

	open.Insert(&node{position: start, potential: potentialFn(start)})
	closed[start] = zero

	for !open.IsEmpty() {
		current := open.ExtractMin()
		for _, nb := range g.Neighbors(current.position) {
			next := &node{
				position:  nb,
				parent:    current,
				cost:      g.Cost(current.position, nb) + current.cost,
				potential: potentialFn(nb),
			}

			seen := seenFn(nb)
			cost, ok := closed[seen]
			if !ok || next.cost < cost {
				open.Insert(next)
				closed[seen] = next.cost
			}

			if next.potential == zero && g.Target(next.position, next.cost) {
				return next.cost, newPath(next), closed
			}
		}
	}
	return zero, nil, closed
}

type fibTree[T any] struct {
	value *T
	parent *fibTree[T]
	child []*fibTree[T]
}

func (t *fibTree[T]) addAtEnd(n *fibTree[T]) {
	n.parent = t
	t.child = append(t.child, n)
}

type fibHeap[T any] struct {
	trees []*fibTree[T]
	least *fibTree[T]
	count uint
	less  func(a, b *T) bool
}

func FibHeap[T any](less func(a, b *T) bool) *fibHeap[T] {
	return &fibHeap[T]{less: less}
}

func (h *fibHeap[T]) GetMin() *T {
	return h.least.value
}

func (h *fibHeap[T]) IsEmpty() bool { return h.least == nil }

func (h *fibHeap[T]) Insert(v *T) {
	ntree := &fibTree[T]{value: v}
	h.trees = append(h.trees, ntree)
	if h.least == nil || h.less(v, h.least.value) {
		h.least = ntree
	}
	h.count++
}

func (h *fibHeap[T]) ExtractMin() *T {
	smallest := h.least
	if smallest != nil {
		// Remove smallest from root trees.
		for i := range h.trees {
			pos := h.trees[i]
			if pos == smallest {
				h.trees[i] = h.trees[len(h.trees)-1]
				h.trees = h.trees[:len(h.trees)-1]
				break
			}
		}

		// Add children to root
		h.trees = append(h.trees, smallest.child...)
		smallest.child = smallest.child[:0]

		h.least = nil
		if len(h.trees) > 0 {
			h.consolidate()
		}

		h.count--
		return smallest.value
	}
	return nil
}

func (h *fibHeap[T]) consolidate() {
	aux := make([]*fibTree[T], bits.Len(h.count))
	for _, x := range h.trees {
		order := len(x.child)

		// consolidate the larger roots under smaller roots of same order until we have at most one tree per order.
		for aux[order] != nil {
			y := aux[order]
			if h.less(y.value, x.value) {
				x, y = y, x
			}
			x.addAtEnd(y)
			aux[order] = nil
			order++
		}
		aux[order] = x
	}

	h.trees = h.trees[:0]
	// move ordered trees to root and find least node.
	for _, k := range aux {
		if k != nil {
			k.parent = nil
			h.trees = append(h.trees, k)
			if h.least == nil || h.less(k.value, h.least.value) {
				h.least = k
			}
		}
	}
}

func (h *fibHeap[T]) Merge(a *fibHeap[T]) {
	h.trees = append(h.trees, a.trees...)
	h.count += a.count
	h.consolidate()
}

// func (h *fibHeap[T]) Find(n *T) *fibTree[T] {
	
// }