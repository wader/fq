package main

import (
	"fq/pkg/cli"
	"fq/pkg/decode"
	"fq/pkg/format"
	"os"
)

func main() {
	if err := (cli.Main{
		Stdout: os.Stdout,
		Args:   os.Args,
		FormatsList: [][]*decode.Format{
			format.All,
		},
	}).Run(); err != nil {
		os.Exit(1)
	}
}
