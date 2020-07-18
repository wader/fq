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
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

var forceFormatNameFlag = flag.String("f", "", "")
var verboseFlag = flag.Bool("v", false, "")

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

	registry := decode.NewRegistryWithFormats(format.All)

	if flag.Arg(0) == "" {
		maxNameLen := 0
		for _, f := range registry.Formats {
			if len(f.Name) > maxNameLen {
				maxNameLen = len(f.Name)
			}
		}

		for _, f := range registry.Formats {
			fmt.Printf("%s%s    %s\n", f.Name, strings.Repeat(" ", maxNameLen-len(f.Name)), strings.Join(f.MIMEs, ", "))
		}
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	var forceFormats []*decode.Format
	if *forceFormatNameFlag != "" {
		forceFormat := registry.FindFormat(*forceFormatNameFlag)
		if forceFormat == nil {
			panic("found not find format " + *forceFormatNameFlag)
		}
		forceFormats = append(forceFormats, forceFormat)
	}
	bb := bitbuf.NewFromBytes(buf)
	d, errs := registry.Probe(nil, flag.Arg(0), decode.Range{Start: 0, Stop: bb.Len}, bitbuf.NewFromBytes(buf), forceFormats)
	if d == nil || *verboseFlag {
		for _, err := range errs {
			fmt.Printf("%s\n", err)
			if pe := err.(*decode.ProbeError); pe != nil {
				// if pe.PanicHandeled {
				fmt.Printf("%s", pe.PanicStack)
				// }
			}
		}
	}

	if d != nil {
		f := d.Root()
		exp := flag.Arg(1)
		if _, err := f.Eval(os.Stdout, exp); err != nil {
			panic(err)
		}
	}
}
