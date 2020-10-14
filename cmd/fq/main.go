package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"
	_ "fq/pkg/format/all"
)

func main() {
	cli.StandardOSMain(format.MustAll())
}
