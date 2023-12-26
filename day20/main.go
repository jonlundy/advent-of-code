package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

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
	m := &machine{}
	receivers := make(map[string][]string)

	for scan.Scan() {
		text := scan.Text()

		name, text, _ := strings.Cut(text, " -> ")
		dest := strings.Split(text, ", ")

		switch {
		case name == "broadcaster":
			m.Add(name, &broadcaster{dest: dest})
		case strings.HasPrefix(name, "%"):
			name = strings.TrimPrefix(name, "%")
			m.Add(name, &flipflop{name: name, dest: dest})

		case strings.HasPrefix(name, "&"):
			name = strings.TrimPrefix(name, "&")
			m.Add(name, &conjunction{name: name, dest: dest})
		}

		for _, d := range dest {
			// rx is present so enable pt 2
			if d == "rx" {
				m.Add("rx", &rx{})
			}
			receivers[d] = append(receivers[d], name)
		}
	}

	m.setup(receivers)

	result := &result{}

	for i := 0; i < 10_000; i++ { // need enough presses to find the best LCM values for each conjunction
		if i == 1000 {
			result.valuePT1 = m.highPulses * m.lowPulses
		}
		m.Push(i)
	}

	// rx is present.. perform part 2.
	if rx, ok := receivers["rx"]; ok {
		tip := m.m[rx[0]].(*conjunction) // panic if missing!

		var lvalues []int
		for k, v := range tip.pushes {
			for i, h := range makeHistory(v) {
				if i == 1 && len(h) > 0 && h[0] > 0 {
					fmt.Println(tip.name, k, "frequency", h[0])
					lvalues = append(lvalues, h[0])
				}
			}
		}

		result.valuePT2 = aoc.LCM(lvalues...)
		fmt.Println(tip.name, "LCM", result.valuePT2, lvalues)
	}
	return result, nil
}

type signal bool

const (
	LOW  signal = false
	HIGH signal = true
)

type message struct {
	signal
	from, to string
}

type machine struct {
	m map[string]pulser

	queue []message

	press      int
	highPulses int
	lowPulses  int
}

func (m *machine) Add(name string, p pulser) {
	if m.m == nil {
		m.m = make(map[string]pulser)
	}
	p.SetMachine(m)
	m.m[name] = p
}
func (m *machine) Send(msgs ...message) {
	m.queue = append(m.queue, msgs...)
	for _, msg := range msgs {
		// fmt.Println(msg)
		if msg.signal {
			m.highPulses++
		} else {
			m.lowPulses++
		}
	}
}
func (m *machine) Push(i int) {
	m.press = i
	m.Send(generate(LOW, "button", "broadcaster")...)
	m.processQueue(i)
}
func (m *machine) processQueue(i int) {
	// look for work and process up to the queue length. repeat.
	hwm := 0
	for hwm < len(m.queue) {
		end := len(m.queue)

		for ; hwm < end; hwm++ {
			msg := m.queue[hwm]

			if p, ok := m.m[msg.to]; ok {
				// fmt.Println(i, "S:", m.m[msg.from], msg.signal, "R:", p)
				p.Pulse(msg)
			}
		}

		hwm = 0
		copy(m.queue, m.queue[end:])
		m.queue = m.queue[:len(m.queue)-end]
		// fmt.Println("")
	}
}
func (m *machine) setup(receivers map[string][]string) {
	for name, recv := range receivers {
		if p, ok := m.m[name]; ok {
			if p, ok := p.(interface{ Receive(...string) }); ok {
				p.Receive(recv...)
			}
		}
	}
}

type pulser interface {
	Pulse(message)
	SetMachine(*machine)
}

// IsModule implements the machine registration for each module.
type IsModule struct {
	*machine
}

func (p *IsModule) SetMachine(m *machine) { p.machine = m }

type broadcaster struct {
	dest []string
	IsModule
}

func (b *broadcaster) Pulse(msg message) {
	b.Send(generate(msg.signal, "broadcaster", b.dest...)...)
}

type flipflop struct {
	name  string
	state signal
	dest  []string

	IsModule
}

func (b *flipflop) Pulse(msg message) {
	if !msg.signal {
		b.state = !b.state
		b.Send(generate(b.state, b.name, b.dest...)...)
	}
}

type conjunction struct {
	name  string
	state map[string]signal
	dest  []string

	pushes map[string][]int

	IsModule
}

func (b *conjunction) Receive(names ...string) {
	if b.state == nil {
		b.state = make(map[string]signal)
		b.pushes = make(map[string][]int)
	}
	for _, name := range names {
		b.state[name] = false
		b.pushes[name] = []int{}
	}
}
func (b *conjunction) Pulse(msg message) {
	b.state[msg.from] = msg.signal

	if msg.signal {
		// collect frequency of pushes to esti,ate rate
		b.pushes[msg.from] = append(b.pushes[msg.from], b.press)
	}

	if all(HIGH, maps.Values(b.state)...) {
		b.Send(generate(LOW, b.name, b.dest...)...)
		return
	}
	b.Send(generate(HIGH, b.name, b.dest...)...)
}

type rx struct {
	IsModule
}

func (rx *rx) Pulse(msg message) {
	if !msg.signal {
		panic("pulse received") // will never happen...
	}
}

// helper funcs
func all[T comparable](match T, lis ...T) bool {
	for _, b := range lis {
		if b != match {
			return false
		}
	}

	return true
}
func generate(t signal, from string, destinations ...string) []message {
	msgs := make([]message, len(destinations))
	for i, to := range destinations {
		msgs[i] = message{signal: t, from: from, to: to}
	}
	return msgs
}

// makeHistory from day 9
func makeHistory(in []int) [][]int {
	var history [][]int
	history = append(history, in)

	// for {
	var diffs []int

	current := history[len(history)-1]

	for i := range current[1:] {
		diffs = append(diffs, current[i+1]-current[i])
	}

	history = append(history, diffs)

	// if len(diffs) == 0 || aoc.Max(diffs[0], diffs[1:]...) == 0 && aoc.Min(diffs[0], diffs[1:]...) == 0 {
	// 	break
	// }
	// }
	return history
}
