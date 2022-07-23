package json

import (
	"embed"

	"github.com/wader/fq/pkg/interp"
)

//go:embed jq.jq
var jqFS embed.FS

func init() {
	interp.RegisterFS(jqFS)
}
