package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("DEBUG_LEVEL"))); err == nil && level != 0 {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	}

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: day05 FILE")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	scan := bufio.NewScanner(input)

	minLocation, minRangeLocation := run(scan)

	fmt.Println("min location:", minLocation)
	fmt.Println("min range location:", minRangeLocation)

}

func run(scan *bufio.Scanner) (int, int) {
	var seeds []int
	var seedRanges []int
	lookup := map[string]*Lookup{}

	for scan.Scan() {
		text := scan.Text()
		if strings.HasPrefix(text, "seeds:") && len(seeds) == 0 {
			seeds, seedRanges = readSeeds(text)
		}

		lookup = readMaps(scan)
	}

	find := NewFinder(
		lookup["seed-to-soil"],
		lookup["soil-to-fertilizer"],
		lookup["fertilizer-to-water"],
		lookup["water-to-light"],
		lookup["light-to-temperature"],
		lookup["temperature-to-humidity"],
		lookup["humidity-to-location"],
	)

	seedLocations := make([]int, len(seeds))
	for i, s := range seeds {
		seedLocations[i] = find.Find(s)
	}
	minLocation := min(seedLocations...)

	seedRangeLocations := make([]int, len(seedRanges))
	for i, s := range seedRanges {
		seedRangeLocations[i] = find.Find(s)
	}
	minRangeLocation := min(seedRangeLocations...)

	return minLocation, minRangeLocation
}


func readSeeds(text string) ([]int, []int) {
	var seeds, seedRanges []int
	sp := strings.Fields(strings.TrimPrefix(text, "seeds: "))
	for i, s := range sp {
		n, _ := strconv.Atoi(s)
		seeds = append(seeds, n)

		if i%2 == 0 {
			seedRanges = append(seedRanges, n)
		} else {
			lastN := seedRanges[len(seedRanges)-1]
			r := make([]int, n-1)
			for i := range r {
				r[i] = lastN + i + 1
			}
			seedRanges = append(seedRanges, r...)
		}
	}
	return seeds, seedRanges
}

func readMaps(scan *bufio.Scanner) map[string]*Lookup {
	var cur *Lookup
	lookup := make(map[string]*Lookup)
	for scan.Scan() {
		text := scan.Text()

		if strings.HasSuffix(text, "map:") {
			if cur != nil {
				cur.Sort()
			}
			cur = &Lookup{}
			title := strings.TrimSuffix(text, " map:")
			lookup[title] = cur
		}

		numbers := strings.Fields(text)
		if len(numbers) == 3 {
			rng := make([]int, 3)
			for i, s := range numbers {
				n, _ := strconv.Atoi(s)
				rng[i] = n
			}
			cur.Add(rng[1], rng[0], rng[2])
		}
	}
	return lookup
}

type Range struct {
	src  int
	dest int
	len  int
}
type Ranges []Range

func (r Ranges) Len() int           { return len(r) }
func (r Ranges) Less(i, j int) bool { return r[i].src < r[j].src }
func (r Ranges) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Lookup struct {
	ranges Ranges
}

func (l *Lookup) Add(src, dest, len int) {
	if l == nil {
		return
	}
	l.ranges = append(l.ranges, Range{src, dest, len})
}
func (l *Lookup) Sort() {
	sort.Sort(sort.Reverse(l.ranges))
}
func (l *Lookup) Find(n int) int {
	for _, r := range l.ranges {
		if n >= r.src && n <= r.src+r.len {
			diff := n - r.src
			if diff < r.len {
				return r.dest + diff
			}
		}
	}

	return n
}

type Finder struct {
	stack []*Lookup
}

func NewFinder(stack ...*Lookup) *Finder {
	return &Finder{stack: stack}
}
func (f *Finder) Find(n int) int {
	// fmt.Print("Find: ")
	for _, l := range f.stack {
		// fmt.Print(n, "->")
		n = l.Find(n)
		// fmt.Print(n, " ")
	}
	// fmt.Println("")
	return n
}

func min(arr ...int) int {
	m := arr[0]
	for _, a := range arr[1:] {
		if m > a {
			m = a
		}
	}
	return m
}
