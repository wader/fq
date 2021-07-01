package main

import (
	"fq/format"
	_ "fq/format/all"
	"fq/pkg/cli"
)

var Version = "dev"

func main() {
	defer cli.MaybeProfile()()
	cli.MaybeLogFile()
	cli.Main(format.DefaultRegistry, Version)
}
