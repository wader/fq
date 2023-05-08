package lazyre_test

import (
	"testing"

	"github.com/wader/fq/internal/lazyre"
)

func TestMust(t *testing.T) {
	if !lazyre.New("a").Must().MatchString("a") {
		t.Fatal("should compile and be non-nil and match a")
	}
}
