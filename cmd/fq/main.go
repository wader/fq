package main

import (
	"fq/pkg/cli"
	"fq/pkg/format"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	cli.StandardOSMain(format.All)
}
