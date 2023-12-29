package main

import (
	"bufio"
	_ "embed"
	"fmt"

	aoc "go.sour.is/advent-of-code"
	"golang.org/x/exp/maps"
)

// var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {
	var m aoc.Map[rune]

	for scan.Scan() {
		text := scan.Text()
		m = append(m, []rune(text))
	}

	result := result{}
	result.valuePT1 = search(m, 1, 3)
	result.valuePT2 = search(m, 4, 10)

	return &result, nil
}

func search(m aoc.Map[rune], minSteps, maxSteps int) int {
	type direction int8
	type rotate int8

	const (
		CW  rotate = 1
		CCW rotate = -1
	)

	var (
		U = aoc.Point{-1, 0}
		R = aoc.Point{0, 1}
		D = aoc.Point{1, 0}
		L = aoc.Point{0, -1}
	)

	var Directions = map[aoc.Point]direction{
		U: 0, // U
		R: 1, // R
		D: 2, // D
		L: 3, // L
	}
	var DirectionIDX = maps.Keys(Directions)

	rows, cols := m.Size()
	target := aoc.Point{rows - 1, cols - 1}

	type position struct {
		loc       aoc.Point
		direction aoc.Point
		steps     int
	}

	step := func(p position) position {
		return position{p.loc.Add(p.direction), p.direction, p.steps + 1}
	}
	rotateAndStep := func(p position, towards rotate) position {
		d := DirectionIDX[(int8(Directions[p.direction])+int8(towards)+4)%4]
		return position{p.loc.Add(d), d, 1}
	}

	type memo struct {
		cost int
		position
	}
	less := func(a, b memo) bool {
		if a.cost != b.cost {
			return a.cost < b.cost
		}
		if a.position.loc != b.position.loc {
			return b.position.loc.Less(a.position.loc)
		}
		if a.position.direction != b.position.direction {
			return b.position.direction.Less(a.position.direction)
		}
		return b.steps < a.steps
	}

	pq := aoc.PriorityQueue(less)
	pq.Enqueue(memo{position: position{direction: D}})
	pq.Enqueue(memo{position: position{direction: R}})
	visited := aoc.Set[position]()

	for !pq.IsEmpty() {
		current, _ := pq.Dequeue()

		if current.loc == target && current.steps >= minSteps {
			return current.cost
		}

		seen := position{loc: current.loc, steps: current.steps}
		if visited.Has(seen) {
			continue
		}
		visited.Add(seen)

		if left := rotateAndStep(current.position, CCW); current.steps >= minSteps && m.Valid(left.loc) {
			_, cost, _ := m.Get(left.loc)
			pq.Enqueue(memo{cost: current.cost + int(cost-'0'), position: left})
		}

		if right := rotateAndStep(current.position, CW); current.steps >= minSteps && m.Valid(right.loc) {
			_, cost, _ := m.Get(right.loc)
			pq.Enqueue(memo{cost: current.cost + int(cost-'0'), position: right})
		}

		if forward := step(current.position); current.steps < maxSteps && m.Valid(forward.loc) {
			_, cost, _ := m.Get(forward.loc)
			pq.Enqueue(memo{cost: current.cost + int(cost-'0'), position: forward})
		}
	}

	return -1
}
