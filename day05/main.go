package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"

)

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	minLocation int
	minRange    int
}

func run(scan *bufio.Scanner) (result, error) {
	log("begin...")

	var seeds []int
	var seedRanges [][2]int
	lookup := map[string]*Lookup{}

	for scan.Scan() {
		text := scan.Text()
		if strings.HasPrefix(text, "seeds:") && len(seeds) == 0 {
			seeds, seedRanges = readSeeds(text)
			log("seeds", len(seeds), "ranges", len(seedRanges))
		}

		lookup = readMaps(scan)
		log("lookups", len(lookup))
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

	return result{findMinLocation(seeds, find), FindMinRangeLocationMulti(seedRanges, find)}, nil
}

func readSeeds(text string) ([]int, [][2]int) {
	var seeds []int
	var seedRanges [][2]int

	for i, n := range aoc.SliceMap(aoc.Atoi, strings.Fields(strings.TrimPrefix(text, "seeds: "))...) {
		seeds = append(seeds, n)

		if i%2 == 0 {
			seedRanges = append(seedRanges, [2]int{n, 0})
		} else {
			seedRanges[len(seedRanges)-1][1] = n
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

func findMinLocation(seeds []int, find *Finder) int {
	seedLocations := make([]int, len(seeds))
	for i, s := range seeds {
		seedLocations[i] = find.Find(s)
	}
	return min(seedLocations...)
}

func FindMinRangeLocation(ranges [][2]int, find *Finder) int {
	results := 0
	for _, r := range ranges {
		results += r[1]
	}

	seedLocations := make([]int, 0, results)

	for _, s := range ranges {
		for i := 0; i < s[1]; i++ {
			seedLocations = append(seedLocations, find.Find(s[0]+i))
		}
	}
	return min(seedLocations...)
}

func FindMinRangeLocationMulti(ranges [][2]int, find *Finder) int {
	worker := func(id int, jobs <-chan [2]int, results chan<- []int) {
		for s := range jobs {
			res := make([]int, s[1])
			for i := 0; i < s[1]; i++ {
				res[i] = find.Find(s[0] + i)
			}
			results <- res
		}
	}

	numWorkers := 16
	jobsCh := make(chan [2]int, numWorkers)
	resultsCh := make(chan []int, len(ranges))

	for w := 0; w < numWorkers; w++ {
		go worker(w, jobsCh, resultsCh)
	}
	log("started workers", numWorkers)

	go func() {
		for i, s := range ranges {
			log("job", i, "send", s)
			jobsCh <- s
		}
		close(jobsCh)
	}()

	results := 0
	for _, r := range ranges {
		results += r[1]
	}
	log("expecting results", results)

	seedLocations := make([]int, 0, results)
	expectResults := make([]struct{}, len(ranges))
	for range expectResults {
		r := <-resultsCh
		seedLocations = append(seedLocations, r...)
	}

	return min(seedLocations...)
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
	if len(arr) == 0 {
		return 0
	}
	m := arr[0]
	for _, a := range arr[1:] {
		if m > a {
			m = a
		}
	}
	return m
}

func log(v ...any) {
	fmt.Fprintln(os.Stderr, v...)
}
