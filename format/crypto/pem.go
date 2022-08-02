package crypto

import (
	"embed"

	"github.com/wader/fq/pkg/interp"
)

//go:embed pem.jq
var pemFS embed.FS

func init() {
	interp.RegisterFS(pemFS)
}
