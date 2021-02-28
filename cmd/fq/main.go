package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"
	_ "fq/pkg/format/all"
	"io"
	"log"
	"os"
)

// TODO: remove?
func init() {
	if lf := os.Getenv("LOGFILE"); lf != "" {
		log.SetOutput(func() io.Writer { f, _ := os.Create(lf); return f }())
	}
}

func main() {
	cli.Main(format.DefaultRegistry)
}
