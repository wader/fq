package postgres14_test

import (
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"testing"
)

func TestTypeAlign8(t *testing.T) {
	expected39 := postgres14.TypeAlign8(39)
	if expected39 != 40 {
		t.Errorf("must be 40")
	}
	expected41 := postgres14.TypeAlign8(41)
	if expected41 != 48 {
		t.Errorf("must be 40")
	}
}
