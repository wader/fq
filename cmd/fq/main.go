package main

import (
	"fq/pkg/cli"
	"fq/pkg/decode"
	"fq/pkg/format"
	"os"
)

func main() {
	if err := (cli.Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Args:   os.Args,
		FormatsList: [][]*decode.Format{
			format.All,
		},
	}).Run(); err != nil {
		os.Exit(1)
	}
}
