package main

import (
	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/cli"
)

const version = "0.0.6"

func main() {
	cli.Main(registry.Default, version)
}
