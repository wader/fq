//go:build profile

package cli

import (
	"os"

	"github.com/wader/fq/internal/profile"
)

func maybeProfile() func() {
	return profile.Start(os.Getenv("CPUPROFILE"), os.Getenv("MEMPROFILE"))
}
