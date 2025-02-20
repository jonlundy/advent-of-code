package main

import (
	"bufio"
	_ "embed"
	"fmt"

	aoc "go.sour.is/advent-of-code"
)

var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 int
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {
	var m aoc.Map[int16, rune]

	for scan.Scan() {
		text := scan.Text()
		m = append(m, []rune(text))
	}
	log("start day 17")

	result := result{}
	result.valuePT1 = search(m, 1, 3, seenFn)
	log("result from part 1 = ", result.valuePT1)

	result.valuePT2 = search(m, 4, 10, nil)
	log("result from part 2 = ", result.valuePT2)

	return &result, nil
}

type Point = aoc.Point[int16]
type Map = aoc.Map[int16, rune]

// rotate for changing direction
type rotate int8

const (
	CW  rotate = 1
	CCW rotate = -1
)

// diretion of path steps
type direction int8

var (
	U = Point{-1, 0}
	R = Point{0, 1}
	D = Point{1, 0}
	L = Point{0, -1}
)

var directions = []Point{U, R, D, L}

var directionIDX = func() map[Point]direction {
	m := make(map[Point]direction, len(directions))
	for k, v := range directions {
		m[v] = direction(k)
	}
	return m
}()

// position on the map
type position struct {
	loc       Point
	direction Point
	steps     int8
}

func (p position) step() position {
	return position{p.loc.Add(p.direction), p.direction, p.steps + 1}
}
func (p position) rotateAndStep(towards rotate) position {
	d := directions[(int8(directionIDX[p.direction])+int8(towards)+4)%4]
	return position{p.loc.Add(d), d, 1}
}

// implements FindPath graph interface
type graph struct {
	min, max int8
	m        Map
	target   Point
	reads    int
	seenFn   func(a position) position
}

// Neighbors returns valid steps from given position. if at target returns none.
func (g *graph) Neighbors(current position) []position {
	var nbs []position

	if current.steps == 0 {
		return []position{
			{R, R, 1},
			{D, D, 1},
		}
	}

	if current.loc == g.target {
		return nil
	}

	if left := current.rotateAndStep(CCW); current.steps >= g.min && g.m.Valid(left.loc) {
		nbs = append(nbs, left)
	}

	if right := current.rotateAndStep(CW); current.steps >= g.min && g.m.Valid(right.loc) {
		nbs = append(nbs, right)
	}

	if forward := current.step(); current.steps < g.max && g.m.Valid(forward.loc) {
		nbs = append(nbs, forward)
	}

	return nbs
}

// Cost calculates heat cost to neighbor from map
func (g *graph) Cost(a, b position) int16 {
	g.reads++
	_, r, _ := g.m.Get(b.loc)
	return int16(r - '0')
}

// Potential calculates distance to target
// func (g *graph) Potential(a position) int16 {
// 	return aoc.ManhattanDistance(a.loc, g.target)
// }

// Target returns true when target reached. receives node and cost.
func (g *graph) Target(a position, c int16) bool {
	if a.loc == g.target && a.steps >= g.min && a.steps <= g.max {
		return true
	}
	return false
}

// Seen attempt at simplifying the seen to use horizontal/vertical and no steps.
// It returns correct for part1 but not part 2..
func (g *graph) Seen(a position) position {
	if g.seenFn != nil {
		return g.seenFn(a)
	}
	return a
}

func seenFn(a position) position {
	if a.direction == U {
		a.direction = D
	}
	if a.direction == L {
		a.direction = R
	}
	// a.steps = 0
	return a
}

func search(m Map, minSteps, maxSteps int8, seenFn func(position) position) int {
	rows, cols := m.Size()
	start := Point{}
	target := Point{rows - 1, cols - 1}

	g := graph{min: minSteps, max: maxSteps, m: m, target: target, seenFn: seenFn}

	cost, path, closed := aoc.FindPath[int16, position](&g, position{loc: start}, position{loc: target})

	log("total map reads = ", g.reads, "cost = ", cost)
	printGraph(m, path, closed, g.seenFn)

	return int(cost)
}

// printGraph with the path/cost overlay
func printGraph(m Map, path []position, closed map[position]int16, seenFn func(a position) position) {
	pts := make(map[Point]position, len(path))
	for _, pt := range path {
		pts[pt.loc] = pt
	}

	clpt := make(map[position]position, len(closed))
	for pt := range closed {
		clpt[position{loc: pt.loc, steps: pt.steps}] = pt
	}

	for r, row := range m {
		// if r == 0 {
		// 	for c := range row {
		// 		if c == 0 {
		// 			fmt.Print("     ")
		// 		}
		// 		fmt.Printf("% 5d", c)
		// 	}
		// 	fmt.Println("")
		// }
		for c := range row {
			// if c == 0 {
			// 	fmt.Printf("% 5d", r)
			// }

			if pt, ok := pts[Point{int16(r), int16(c)}]; ok {
				if seenFn != nil {
					pt = seenFn(pt)
				}
				_ = pt
				// fmt.Printf("% 5d", closed[pt])
				fmt.Print("*")

				continue
			}

			// fmt.Print(" ....")
			fmt.Print(" ")
		}
		fmt.Println("")
	}
	fmt.Println("")
}
