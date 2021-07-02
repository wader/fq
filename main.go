package main

import (
	"fq/format/all"
	"fq/pkg/cli"
)

var Version = "dev"

func main() {
	defer cli.MaybeProfile()()
	cli.MaybeLogFile()
	cli.Main(all.Registry, Version)
}
