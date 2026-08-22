package math

import (
	"embed"

	"github.com/wader/fq/pkg/interp"
)

//go:embed base.jq
var radixFS embed.FS

func init() {
	interp.RegisterFS(radixFS)
}
