package main

import (
	_ "fq/format/all"
	"fq/format/registry"
	"fq/pkg/cli"
)

var version = "dev"

func main() {
	defer cli.MaybeProfile()()
	cli.MaybeLogFile()
	cli.Main(registry.Default, version)
}
