package postgres_test

import (
	"github.com/wader/fq/format/postgres"

	"testing"
)

func TestParseLsn(t *testing.T) {
	lsn1, err := postgres.ParseLsn("0/4E394440")
	if err != nil {
		t.Fatalf("TestParseLsn 1, err = %v\n", err)
	}
	if lsn1 != 0x4E394440 {
		t.Fatalf("TestParseLsn 2, invalid lsn value\n")
	}

	lsn2, err := postgres.ParseLsn("0/4469E930")
	if err != nil {
		t.Fatalf("TestParseLsn 3, err = %v\n", err)
	}
	if lsn2 != 0x4469E930 {
		t.Fatalf("TestParseLsn 4, invalid lsn value\n")
	}
}

func TestXLogSegmentOffset(t *testing.T) {
	offset := postgres.XLogSegmentOffset(0x4E394440)
	if offset == 0 {
		t.Fatalf("TestXLogSegmentOffset 1, invalid offset\n")
	}
}
