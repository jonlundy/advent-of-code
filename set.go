package aoc

import "golang.org/x/exp/maps"

type set[T comparable] map[T]struct{}

func Set[T comparable](arr ...T) set[T] {
	m := make(set[T], len(arr))
	for _, a := range arr {
		m[a] = struct{}{}
	}
	return m
}
func (m *set[T]) Add(a T) {
	(*m)[a] = struct{}{}
}
func (m *set[T]) Items() []T {
	return maps.Keys(*m)
}
func (m *set[T]) Has(a T) bool {
	var ok bool
	_, ok = (*m)[a]
	return ok
}

type defaultMap[K comparable, V any] struct {
	m map[K]V
	d V
}

func DefaultMap[K comparable, V any](d V) *defaultMap[K, V] {
	return &defaultMap[K, V]{
		make(map[K]V),
		d,
	}
}

func (m *defaultMap[K, V]) Set(k K, v V) {
	m.m[k] = v
}
func (m *defaultMap[K, V]) Get(k K) (V, bool) {
	if v, ok := m.m[k]; ok {
		return v, true
	}
	return m.d, false
}

type pair[K, V any] struct {
	K K
	V V
}

func (m *defaultMap[K, V]) Items() []pair[K, V] {
	var items = make([]pair[K, V], 0, len(m.m))
	for k, v := range m.m {
		items = append(items, pair[K, V]{k, v})
	}
	return items
}
