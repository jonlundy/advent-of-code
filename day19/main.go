package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

	aoc "go.sour.is/advent-of-code"
)

// var log = aoc.Log

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	valuePT1 int
	valuePT2 uint
}

func (r result) String() string { return fmt.Sprintf("%#v", r) }

func run(scan *bufio.Scanner) (*result, error) {

	var workflows = make(map[string][]rule)
	var parts []part

	for scan.Scan() {
		text := strings.TrimSuffix(scan.Text(), "}")
		if len(text) == 0 {
			continue
		}

		// Is Part
		if p, ok := scanPart(text); ok {
			parts = append(parts, p)
			continue
		}

		if name, r, ok := scanRule(text); ok {
			workflows[name] = r
		}
	}

	var result result
	result.valuePT1 = solveWorkflow(parts, workflows)
	result.valuePT2 = solveRanges(workflows)

	return &result, nil
}

type part struct {
	x, m, a, s int
}

func (p part) String() string {
	return fmt.Sprintf("{x:%v m:%v a:%v s:%v}", p.x, p.m, p.a, p.s)
}
func scanPart(text string) (part, bool) {
	var p part

	// Is Part
	if text[0] == '{' {
		for _, s := range strings.Split(text[1:], ",") {
			a, b, _ := strings.Cut(s, "=")
			i := aoc.Atoi(b)
			switch a {
			case "x":
				p.x = i
			case "m":
				p.m = i
			case "a":
				p.a = i
			case "s":
				p.s = i
			}
		}
		return p, true
	}
	return p, false
}

type rule struct {
	match string
	op    string
	value int
	queue string
}

func scanRule(text string) (string, []rule, bool) {
	name, text, _ := strings.Cut(text, "{")
	var r []rule
	for _, s := range strings.Split(text, ",") {
		if a, b, ok := strings.Cut(s, "<"); ok {
			b, c, _ := strings.Cut(b, ":")
			r = append(r, rule{
				match: a,
				op:    "<",
				value: aoc.Atoi(b),
				queue: c,
			})

			continue
		}
		if a, b, ok := strings.Cut(s, ">"); ok {
			b, c, _ := strings.Cut(b, ":")
			r = append(r, rule{
				match: a,
				op:    ">",
				value: aoc.Atoi(b),
				queue: c,
			})
			continue
		}

		// default queue comes last
		r = append(r, rule{queue: s})
		break
	}
	return name, r, len(r) > 0
}
func (r rule) Match(p part) bool {
	var value int

	switch r.match {
	case "x":
		value = p.x
	case "m":
		value = p.m
	case "a":
		value = p.a
	case "s":
		value = p.s
	default:
		return true // default to new queue
	}

	if r.op == ">" && value > r.value {
		return true
	} else if r.op == "<" && value < r.value {
		return true
	}
	return false // no match
}

func solveWorkflow(parts []part, workflows map[string][]rule) int {
	// var rejected []part
	var accepted []part

	for _, p := range parts {
		workflow := "in"

	nextStep:
		for workflow != "" {
			for _, r := range workflows[workflow] {
				if !r.Match(p) {
					continue
				}
				workflow = r.queue

				if workflow == "A" {
					accepted = append(accepted, p)
					workflow = ""
					break nextStep
				}
				if workflow == "R" {
					// rejected = append(rejected, p)
					workflow = ""
					break nextStep
				}

				continue nextStep
			}
		}
	}

	sum := 0
	for _, p := range accepted {
		sum += p.x
		sum += p.m
		sum += p.a
		sum += p.s
	}
	return sum
}

func solveRanges(workflows map[string][]rule) uint {

	pq := aoc.PriorityQueue(func(a, b queue) bool { return false })
	pq.Enqueue(queue{
		"in",
		block{
			ranger{1, 4000},
			ranger{1, 4000},
			ranger{1, 4000},
			ranger{1, 4000},
		}})

	var accepted []block
	// var rejected []block

	for !pq.IsEmpty() {
		current, _ := pq.Dequeue()
		for _, rule := range workflows[current.name] {
			next := queue{name: rule.queue, block: current.block}

			switch rule.match {
			case "x":
				current.x, next.x = split(current.x, rule.value, rule.op == ">")
			case "m":
				current.m, next.m = split(current.m, rule.value, rule.op == ">")
			case "a":
				current.a, next.a = split(current.a, rule.value, rule.op == ">")
			case "s":
				current.s, next.s = split(current.s, rule.value, rule.op == ">")
			}

			switch next.name {
			case "R":
				// rejected = append(rejected, next.block)

			case "A":
				accepted = append(accepted, next.block)

			default:
				pq.Enqueue(next)
			}
		}
	}

	var sum uint
	for _, a := range accepted {
		sum += uint((a.x[1]-a.x[0]+1) * (a.m[1]-a.m[0]+1) * (a.a[1]-a.a[0]+1) * (a.s[1]-a.s[0]+1))
	}

	return sum
}

type ranger [2]int
type block struct {
	x, m, a, s ranger
}
type queue struct {
	name string
	block
}

func split(a ranger, n int, gt bool) (current ranger, next ranger) {
	if gt { // x > N => [0,N] [N++,inf]
		return ranger{a[0], n}, ranger{n + 1, a[1]}
	}

	// x < N => [N,inf] [0,N--]
	return ranger{n, a[1]}, ranger{a[0], n - 1}
}
