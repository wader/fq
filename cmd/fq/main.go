package main

import (
	"flag"
	"fmt"
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

var forceFormatName = flag.String("f", "", "")

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

	buf, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	registry := decode.NewRegistryWithFormats(format.All)
	var forceFormats []*decode.Format
	if *forceFormatName != "" {
		forceFormat := registry.FindFormat(*forceFormatName)
		if forceFormat == nil {
			panic("found not find format " + *forceFormatName)
		}
		forceFormats = append(forceFormats, forceFormat)
	}
	d, errs := registry.Probe(nil, bitbuf.NewFromBytes(buf), forceFormats)
	for _, err := range errs {
		fmt.Printf("%s\n", err)
		if pe := err.(*decode.ProbeError); pe != nil {
			if pe.PanicHandeled {
				fmt.Printf("%s", pe.PanicStack)
			}
		}
	}
	if d != nil {
		decode.Dump(d.GetCommon().Current, 0)
	}
}
