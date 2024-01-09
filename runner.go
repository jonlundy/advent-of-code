package aoc

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func Runner[R any, F func(*bufio.Scanner) (R, error)](run F) (R, error) {
	if len(os.Args) < 2 {
		Log("Usage:", filepath.Base(os.Args[0]), "FILE")
		os.Exit(22)
	}

	inputFilename := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)

	flag.Parse()
	Log(cpuprofile, memprofile, *cpuprofile, *memprofile)
	if *cpuprofile != "" {
		Log("enabled cpu profile")
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		Log("write cpu profile to", f.Name())
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		Log("enabled mem profile")
		defer func() {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			Log("write mem profile to", f.Name())
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}()
	}


	input, err := os.Open(inputFilename)
	if err != nil {
		Log(err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(input)

	return run(scan)
}

func MustResult[T any](result T, err error) {
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}

	Log("result", result)
}

func Log(v ...any) {
	fmt.Fprint(os.Stderr, time.Now(), ": ")
	fmt.Fprintln(os.Stderr, v...)
}
func Logf(format string, v ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, v...)
}

func ReadStringToInts(fields []string) []int {
	return SliceMap(Atoi, fields...)
}
