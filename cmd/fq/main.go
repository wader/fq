package main

import (
	"flag"
	"fq/internal/bitbuf"
	"fq/internal/decode"
	"fq/internal/format"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

var forceDecoder = flag.String("f", "", "")

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

	// f := &decode.Field{Name: "root"}
	// c := decode.Common{Current: f, Buffer: bitbuf.NewFromBytes(buf)}
	// d := flac.Decoder{Common: c}
	//d := mp3.Decoder{Common: c}
	//d := id3v2.Decoder{Common: c}

	var registers []*decode.Register
	if *forceDecoder != "" {
		for _, r := range format.All {
			if r.Name == *forceDecoder {
				registers = append(registers, r)
			}
		}
		if len(registers) == 0 {
			panic("could not find format")
		}
	} else {
		registers = format.All
	}

	func() {
		defer func() {
			// if e := recover(); e != nil {
			// 	log.Printf("e: %#+v\n", e)
			// }
		}()
		r, d := decode.New(nil, bitbuf.NewFromBytes(buf), registers)
		log.Printf("r: %#+v\n", r)

		decode.Dump(d.Current, 0)

	}()

}
