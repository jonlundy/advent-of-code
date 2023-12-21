package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
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
	var forever bool
	if os.Getenv("AOC_FOREVER") == "1" {
		forever = true
	}

	m := &machine{}
	receivers := make(map[string][]string)

	m.Add("rx", &rx{})

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
			m.Add(name, &conjuction{name: name, dest: dest})
		}

		for _, d := range dest {
			receivers[d] = append(receivers[d], name)
		}
	}

	m.Setup(receivers)

	result := &result{}

	if forever {
		i := 0

		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("## Press %d FINISH %v ##", i, p)
				os.Exit(1)
			}
		}()

		for {
			if i%12345 == 0 {
				fmt.Printf("## Press %d ##\r", i)
			}
			m.Push(i)
			i++
		}
	}

	for i := 0; i < 4_0000; i++ {
		// fmt.Printf("\n## Press %d ##\n\n", i)
		if i == 1000 {
			result.valuePT1 = m.highPulses * m.lowPulses
		}
		m.Push(i)
		
	}

	// fmt.Println("\n## SUMMARY ##")
	// fmt.Println("Sent", LOW, m.lowPulses)
	// fmt.Println("Sent", HIGH, m.highPulses)
	var lvalues []int

	for _, p := range m.m {
		if p, ok := p.(*conjuction); ok && in(p.name, []string{"bk","tp","pt","vd"}) {
			var values []int
			for k, v := range p.pushes {
				for i, h := range makeHistory(v) {
					if i == 1 && len(h) > 3 && h[0] > 0 { //&& all(h[0], h[1:]...) {
						fmt.Println(p.name, k, i, h[0])
						values = append(values, h[0])
					}
				}
			}
			max := aoc.Max(values[0], values...)
			fmt.Println(p.name, "MAX", max, values)
			lvalues = append(lvalues, max)
		}
	}
	result.valuePT2 = aoc.LCM(lvalues...)
	fmt.Println("tg", "LCM", result.valuePT2, lvalues)

	// trace("rx", receivers)

	return result, nil
}

type signal bool

const (
	LOW  signal = false
	HIGH signal = true
)

func (m signal) String() string {
	if m {
		return " >>-HIGH-> "
	}
	return " >>-LOW-> "
}

type message struct {
	signal
	from, to string
}

func (m message) String() string {
	return fmt.Sprint(m.from, m.signal, m.to)
}

type pulser interface {
	Pulse(message)
	SetMachine(*machine)
}

type machine struct {
	m map[string]pulser
	press int

	queue []message
	hwm   int

	stop bool

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
func (m *machine) Setup(receivers map[string][]string) {
	for name, recv := range receivers {
		if p, ok := m.m[name]; ok {
			if p, ok := p.(interface{ Receive(...string) }); ok {
				p.Receive(recv...)
			}
		}
	}
}

func (m *machine) Stop() {
	m.stop = true
}

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
func (b *broadcaster) String() string { return "br" }

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
func (b *flipflop) String() string {
	return fmt.Sprintf("%s(%v)", b.name, b.state)
}

type conjuction struct {
	name  string
	state map[string]signal
	dest  []string

	pushes map[string][]int
	activate []int
	last map[string]int
	max  map[string]int
	lcm  int

	IsModule
}

func (b *conjuction) Receive(names ...string) {
	if b.state == nil {
		b.state = make(map[string]signal)
		b.last = make(map[string]int)
		b.max = make(map[string]int)
		b.pushes = make(map[string][]int)
	}
	for _, name := range names {
		b.state[name] = false
		b.max[name] = int(^uint(0)>>1)
		b.pushes[name] = []int{}

	}
}
func (b *conjuction) Pulse(msg message) {
	if b.state == nil {
		b.state = make(map[string]signal)
		b.last = make(map[string]int)
		b.max = make(map[string]int)
	}

	b.state[msg.from] = msg.signal

	if msg.signal {
		b.pushes[msg.from] = append(b.pushes[msg.from], b.press)
		b.last[msg.from] = b.press - b.last[msg.from]


		// vals := maps.Values(b.max)
		// if aoc.Min(vals[0], vals...) != 0 {
		// 	if lcm := aoc.LCM(vals...); lcm > 0 {
		// 		fmt.Printf("\nfound loop %s = %d %v\n", b.name, lcm, vals)
		// 	}
		// }
	}

	if all(HIGH, maps.Values(b.state)...) {
		b.activate = append(b.activate, b.press)
		b.Send(generate(LOW, b.name, b.dest...)...)
		return
	}
	b.Send(generate(HIGH, b.name, b.dest...)...)
}
func (b *conjuction) String() string {
	return fmt.Sprintf("%s(%v)", b.name, b.state)
}

type rx struct {
	IsModule
}

func (rx *rx) Pulse(msg message) {
	if !msg.signal {
		panic("pulse received")
	}
}
func (rx *rx) String() string { return "rx" }

func all[T ~int|~bool](match T, lis ...T) bool {
	if len(lis) == 0 {
		return true
	}

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

func in(n string, haystack []string) bool {
	for _, h := range haystack {
		if n == h {
			return true
		}
	}
	return false
}

func makeHistory(in []int) [][]int {
	var history [][]int
	history = append(history, in)

	for {
		var diffs []int
		

		current := history[len(history)-1]
		if len(current) == 0 { return nil }

		for i := range current[1:] {
			diffs = append(diffs, current[i+1]-current[i])
		}

		history = append(history, diffs)

		if len(diffs) == 0 || aoc.Max(diffs[0], diffs[1:]...) == 0 && aoc.Min(diffs[0], diffs[1:]...) == 0 {
			break
		}
	}
	return history
}
