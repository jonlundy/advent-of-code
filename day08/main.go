package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	aoc "go.sour.is/advent-of-code"
)

func main() {
	result, err := aoc.Runner(run)
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}

	fmt.Println("result", result)
}

type result struct {
	stepsPT1 uint64
	stepsPT2 uint64
}

func (r result) String() string {
	return fmt.Sprintf("solution 1: %v\nsolution 2: %v\n", r.stepsPT1, r.stepsPT2)
}

func run(scan *bufio.Scanner) (*result, error) {
	var path []rune
	m := make(nodeMap)

	for scan.Scan() {
		text := scan.Text()
		if len(text) == 0 {
			continue
		}

		if len(path) == 0 {
			fmt.Println("path", text)
			path = []rune(strings.TrimSpace(text))
			continue
		}

		n := &node{}
		i, err := fmt.Sscanf(text, "%s = (%s %s", &n.value, &n.lvalue, &n.rvalue)
		if err != nil {
			return nil, err
		}
		n.lvalue = strings.TrimRight(n.lvalue, ",)")
		n.rvalue = strings.TrimRight(n.rvalue, ",)")
		m[n.value] = n

		fmt.Println("value", i, n.value, n.lvalue, n.rvalue)
	}
	if err := m.mapNodes(); err != nil {
		return nil, err
	}

	steps1 := m.SolvePT1(path)
	steps2 := m.SolvePT2(path)

	return &result{steps1, steps2}, nil
}

type node struct {
	value          string
	lvalue, rvalue string
	left, right    *node
}

type nodeMap map[string]*node

func (m nodeMap) mapNodes() error {
	for k, v := range m {
		if ptr, ok := m[v.lvalue]; ok {
			v.left = ptr
		} else {
			return fmt.Errorf("%s L-> %s not found", k, v.lvalue)
		}
		if ptr, ok := m[v.rvalue]; ok {
			v.right = ptr
		} else {
			return fmt.Errorf("%s R-> %s not found", k, v.rvalue)
		}

		m[k] = v
	}
	return nil
}

func (m nodeMap) solver(start string, isEnd func(string) bool, path []rune) uint64 {
	position, ok := m[start]
	if !ok {
		return 0
	}
	var i int
	var steps uint64

	for steps < ^uint64(0) {
		steps++
		if path[i] == 'R' {
			// fmt.Println("step", steps, position.value, "R->", position.rvalue)
			position = position.right
		} else {
			// fmt.Println("step", steps, position.value, "L->", position.lvalue)
			position = position.left
		}

		if isEnd(position.value) {
			break
		}

		i++
		if i > len(path)-1 {
			i = 0
		}
	}
	return steps
}

func (m nodeMap) SolvePT1(path []rune) uint64 {
	fmt.Println("---- PART 1 BEGIN ----")
	defer fmt.Println("---- PART 1 END ----")

	return m.solver("AAA", func(s string) bool { return s == "ZZZ" }, path)
}

func (m nodeMap) SolvePT2(path []rune) uint64 {
	fmt.Println("---- PART 2 BEGIN ----")
	defer fmt.Println("---- PART 2 END ----")

	var starts []*node

	for k, n := range m {
		if strings.HasSuffix(k, "A") {
			fmt.Println("start", k)
			starts = append(starts, n)
		}
	}

	loops := make([]uint64, len(starts))
	for i, n := range starts {
		loops[i] = m.solver(n.value, func(s string) bool { return strings.HasSuffix(s, "Z") }, path)
	}
	return aoc.LCM(loops...)
}
