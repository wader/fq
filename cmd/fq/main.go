package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"
	_ "fq/pkg/format/all"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// TODO: remove?
func init() {
	if lf := os.Getenv("LOGFILE"); lf != "" {
		log.SetOutput(func() io.Writer { f, _ := os.Create(lf); return f }())
	}
}

func main() {
	var cpuprofile = os.Getenv("CPUPROFILE")
	var memprofile = os.Getenv("MEMPROFILE")

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	cli.Main(format.DefaultRegistry)
}
