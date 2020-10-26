package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"
	_ "fq/pkg/format/all"
)

func main() {
	// format.Dot(os.Stdout)
	// os.Exit(0)

	cli.StandardOSMain(format.DefaultRegistry)
}
