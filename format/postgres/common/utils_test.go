package common_test

import (
	"testing"

	"github.com/wader/fq/format/postgres/common"
)

func TestTypeAlign(t *testing.T) {
	expected0 := common.TypeAlign(8192, 8192+0)
	if expected0 != 8192 {
		t.Errorf("must be 8192\n")
	}

	expected1 := common.TypeAlign(8192, 8192+100)
	if expected1 != 8192*2 {
		t.Errorf("must be 8192*2\n")
	}

	expected2 := common.TypeAlign(8192, 0)
	if expected2 != 0 {
		t.Errorf("must be 0\n")
	}

	expected3 := common.TypeAlign(8192, 700)
	if expected3 != 8192 {
		t.Errorf("must be 8192\n")
	}

	expected4 := common.TypeAlign(8192, 8192*2+5000)
	if expected4 != 8192*3 {
		t.Errorf("must be 8192*3\n")
	}

	expected5 := common.TypeAlign(8192, 114720)
	if expected5 != 122880 {
		t.Errorf("must be 8192*3\n")
	}
}

func TestTypeAlign8(t *testing.T) {
	expected39 := common.TypeAlign8(39)
	if expected39 != 40 {
		t.Errorf("must be 40\n")
	}
	expected41 := common.TypeAlign8(41)
	if expected41 != 48 {
		t.Errorf("must be 40\n")
	}
}

func TestRoundDown(t *testing.T) {
	const pageSize1 = 8192
	expected1 := common.RoundDown(pageSize1, 7*pageSize1+35)
	if expected1 != 7*pageSize1 {
		t.Errorf("must be %d\n", 7*pageSize1)
	}
	expected2 := common.RoundDown(pageSize1, 7*pageSize1-1)
	if expected2 != 6*pageSize1 {
		t.Errorf("must be %d\n", 6*pageSize1)
	}

	const pageSize2 = 7744
	expected3 := common.RoundDown(pageSize2, 15*pageSize2+61)
	if expected3 != 15*pageSize2 {
		t.Errorf("must be %d\n", 15*pageSize2)
	}
	expected4 := common.RoundDown(pageSize2, 3*pageSize2-15)
	if expected4 != 2*pageSize2 {
		t.Errorf("must be %d\n", 2*pageSize2)
	}

	expected5 := common.RoundDown(pageSize1, 5*pageSize1)
	if expected5 != 5*pageSize1 {
		t.Errorf("must be %d\n", 5*pageSize1)
	}
}

func TestIsMaskSet(t *testing.T) {
	m1 := common.IsMaskSet(0xff+0x1221000, 0xf0)
	if m1 != 1 {
		t.Errorf("mask must be set\n")
	}
	m2 := common.IsMaskSet(0xff+0x1221000, 0xf00)
	if m2 != 0 {
		t.Errorf("mask must be 0\n")
	}
}
