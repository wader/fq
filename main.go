package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"
	_ "fq/pkg/format/all"
)

var Version = "dev"

func main() {
	defer cli.MaybeProfile()()
	cli.MaybeLogFile()
	cli.Main(format.DefaultRegistry, Version)
}
