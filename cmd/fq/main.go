package main

import (
	"flag"
	"fq/internal/decode"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	// s, err := ioutil.ReadFile(flag.Arg(0))
	// if err != nil {
	// 	panic(err)
	// }

	buf, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	f := &decode.Field{Name: "root"}
	p := decode.FLAC{Common: decode.Common{Current: f, Buf: buf}}
	func() {
		defer func() {
			e := recover()
			log.Printf("e: %#+v\n", e)

		}()
		p.Unmarshl()
	}()

	decode.Dump(f, 0)

}
