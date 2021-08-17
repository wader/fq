package main

import (
	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/cli"
)

var version = "dev"

func main() {
	defer cli.MaybeProfile()()
	cli.MaybeLogFile()
	cli.Main(registry.Default, version)
}
