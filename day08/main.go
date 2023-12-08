package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: day08 FILE")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	scan := bufio.NewScanner(input)

	result, err := run(scan)
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}

	fmt.Println("result", result)
}

type result struct {
	stepsPT1 int
	stepsPT2 uint64
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

	steps1 := SolutionPT1(m, path)
	steps2 := SolutionPT2(m, path)

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

func SolutionPT1(m nodeMap, path []rune) int {
	fmt.Println("---- PART 1 BEGIN ----")
	position, ok := m["AAA"]
	if !ok {
		return 0
	}
	var i int
	var steps int

	for steps < 100000 {
		steps++
		if path[i] == 'R' {
			fmt.Println("step", steps, position.value, "R->", position.rvalue)
			position = position.right
		} else {
			fmt.Println("step", steps, position.value, "L->", position.lvalue)
			position = position.left
		}

		if position.value == "ZZZ" {
			break
		}

		i++
		if i > len(path)-1 {
			i = 0
		}
	}
	fmt.Println("---- PART 1 END ----")
	return steps
}

func SolutionPT2(m nodeMap, path []rune) uint64 {
	fmt.Println("---- PART 2 BEGIN ----")

	type loop struct {
		start, position, end *node
		steps                uint64
	}
	loops := make(map[*node]loop)

	endpoints := make(map[*node]struct{})
	for k, n := range m {
		if strings.HasSuffix(k, "A") {
			fmt.Println("start", k)
			loops[n] = loop{start: n, position: n}
		}

		if strings.HasSuffix(k, "Z") {
			fmt.Println("stop", k)
			endpoints[n] = struct{}{}
		}
	}

	var i int
	var steps uint64
	var stops int
	maxUint := ^uint64(0)
	loopsFound := 0

	for steps < maxUint {
		steps++
		if path[i] == 'R' {
			for k, loop := range loops {
				// fmt.Println("step", steps, position.value, "R->", position.rvalue)
				loop.position = loop.position.right
				loops[k] = loop
			}
		} else {
			for k, loop := range loops {
				// fmt.Println("step", steps, position.value, "L->", position.lvalue)
				loop.position = loop.position.left
				loops[k] = loop
			}
		}

		done := true
		s := 0
		for k, loop := range loops {
			if _, ok := endpoints[loop.position]; !ok {
				// fmt.Println("no stop", i, position.value)
				done = false
				// break
			} else {
				// fmt.Println("stop", i, position.value)
				if loop.end == nil {
					loop.end = loop.position
					loop.steps = steps
					fmt.Println("loop found", loop.position.value, "steps", steps)
					loops[k] = loop
					loopsFound++
				}
				s++
			}
		}

		if loopsFound == len(loops) {
			var values []uint64
			for _, loop := range loops {
				values = append(values, loop.steps)
			}
			return LCM(values...)
		}

		if s > stops {
			stops = s
			fmt.Println("stops", stops, "steps", steps)
		}

		if done {
			break
		}

		i++
		if i > len(path)-1 {
			i = 0
		}
	}

	fmt.Println("---- PART 2 END ----")
	return steps
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...uint64) uint64 {
	if len(integers) == 0 {
		return 0
	}
	if len(integers) == 1 {
		return integers[0]
	}

	a, b := integers[0], integers[1]
	result := a * b / GCD(a, b)

	for _, c := range integers[2:] {
		result = LCM(result, c)
	}

	return result
}
