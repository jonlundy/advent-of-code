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
		if text[0] == '{' {
			var p part
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
			parts = append(parts, p)
			continue
		}

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
		workflows[name] = r
	}

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

	var result result

	for _, p := range accepted {
		result.valuePT1 += p.x
		result.valuePT1 += p.m
		result.valuePT1 += p.a
		result.valuePT1 += p.s
	}

	return &result, nil
}

type part struct {
	x, m, a, s int
}
func (p part) String() string {
	return fmt.Sprintf("{x:%v m:%v a:%v s:%v}", p.x,p.m,p.a,p.s)
}
type rule struct {
	match string
	op    string
	value int
	queue string
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
		return  true		
	} 
	return  false // no match
}


func in(n string, haystack ...string) bool {
	for _, h := range haystack {
		if n == h {
			return true
		}
	}
	return false
}
