package main

import (
	"bufio"
	"bytes"
	"testing"

	_ "embed"

	"github.com/matryer/is"
)

//go:embed example.txt
var example []byte

//go:embed input.txt
var input []byte

func TestExample(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example))

	minLocation, minRangeLocation := run(scan)
	is.Equal(minLocation, 35)
	is.Equal(minRangeLocation, 46)
}

func SkipTestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	minLocation, minRangeLocation := run(scan)
	is.Equal(minLocation, 199602917)
	is.Equal(minRangeLocation, 0)

}

func TestLookup(t *testing.T) {
	is := is.New(t)
	find := &Lookup{ranges: Ranges{
		{98, 50, 2},
		{50, 52, 48},
	}}
	is.Equal(find.Find(79), 81)

	find = &Lookup{ranges: Ranges{
		{77,45,23},
		{45,81,19},
		{64,68,13},
	}}
	is.Equal(find.Find(77), 45)

}

func TestFinder(t *testing.T) {
	is := is.New(t)
	find := NewFinder(
		// seed-to-soil
		&Lookup{ranges: Ranges{
			{98, 50, 2},
			{50, 52, 48},
		}},
		// soil-to-fertilizer
		&Lookup{ranges: Ranges{
			{15, 0,37},
			{52,37,2},
			{0,39,15},
		}},
		// fertilizer-to-water
		&Lookup{ranges: Ranges{
			{53,49,8},
			{11,0,42},
			{0,42,7},
			{7,57,4},
		}},
		// water-to-light
		&Lookup{ranges: Ranges{
			{18,88,7},
			{25,18,70},
		}},
		// light-to-temperature
		&Lookup{ranges: Ranges{
			{77,45,23},
			{45,81,19},
			{64,68,13},
		}},
		// temperature-to-humidity
		&Lookup{ranges: Ranges{
			{69,0,1},
			{0,1,69},
		}},
		// humidity-to-location
		&Lookup{ranges: Ranges{
			{56,60,37},
			{93,56,4},
		}},
	)
	is.Equal(find.Find(82), 46)
}