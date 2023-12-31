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
	valuePT2 int
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
	var rejected []part
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
					rejected = append(rejected, p)
					workflow = ""
					break nextStep
				}

				continue nextStep
			}
		}
	}

	fmt.Println("accepted", accepted)
	fmt.Println("rejected", rejected)

	sum := 0
	for _, p := range accepted {
		sum += p.x
		sum += p.m
		sum += p.a
		sum += p.s
	}
	return sum
}

/*
in{s<1351:px, s>=1351:qqz}

--> px{a<2006&&(x<1416||x>2662):A, a<2006&&x>=1416&&x<=2662:R, m>2090:A, a>=2006&&m<=2090&&s<537:R||x>2440:R, a>=2006&&m<=2090&&s<537&x<=2440:A}

--> qqz{s>2770:A, m<1801&&m>838:A, m<1801&&a>1716:R, m<1801&&a<=1716:A, s<=2770&&m>=1801:R}







in [/]
--
s<1351 -> px
s>=1351 -> qqz

px [/]
--
s<1351 -> px

a< 2006  -> qkq
m> 2090  -> A
a>=2006 -> ...
m<=2090 -> rfg

qqz [ ]
--
s>=1351 -> qqz

s> 2770 -> qs
m< 1801 -> hdj
 -> R

qkq [ ]
--
s< 1351 ->
a< 2006 ->

x< 1416 -> A
x>=1416 -> crn

rfg [ ]
--
s< 1351 -> px
a>=2006 ->...
m<=2090 -> rfg

s<  537 -> gd
x> 2440 -> R
s>= 537 ->...
x<=2440 -> A

crn [ ]
--
s< 1351 -> px
a< 2006 -> qkq
x>=1416 -> crn

x> 2662 -> A
x<=2662 -> R

A
--

*/
