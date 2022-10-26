package postgres

import (
	"fmt"
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"

	"strconv"
	"strings"
)

// TO DO
// not ready yet

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_WAL,
		Description: "PostgreSQL write-ahead log file",
		DecodeFn:    decodePGWAL,
		DecodeInArg: format.PostgresWalIn{
			Flavour: PG_FLAVOUR_POSTGRES14,
			Lsn:     "",
		},
	})
}

// https://pgpedia.info/x/XLOG_PAGE_MAGIC.html
//nolint:revive
const (
	XLOG_PAGE_MAGIC_15 = uint16(0xD10F)
	XLOG_PAGE_MAGIC_14 = uint16(0xD10D)
	XLOG_PAGE_MAGIC_13 = uint16(0xD106)
	XLOG_PAGE_MAGIC_12 = uint16(0xD101)
	XLOG_PAGE_MAGIC_11 = uint16(0xD098)
	XLOG_PAGE_MAGIC_10 = uint16(0xD097)
	XLOG_PAGE_MAGIC_96 = uint16(0xD093)
)

func ParseLsn(lsn string) (uint32, error) {
	// check for 0/4E394440
	str1 := lsn
	if strings.Contains(lsn, "/") {
		parts := strings.Split(lsn, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid lsn = %s", lsn)
		}
		str1 = parts[1]
	}
	// parse hex to coded file name + file offset
	r1, err := strconv.ParseInt(str1, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint32(r1), err
}

func XLogSegmentOffset(xLogPtr uint32) uint32 {
	const walSegSizeBytes = 16 * 1024 * 1024
	return xLogPtr & (walSegSizeBytes - 1)
}

func decodePGWAL(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	pgIn, ok := in.(format.PostgresWalIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresIn!\n")
	}

	maxOffset := uint32(0xFFFFFFFF)
	if pgIn.Lsn != "" {
		lsn, err := ParseLsn(pgIn.Lsn)
		if err != nil {
			d.Fatalf("Failed to ParseLsn, err = %v\n", err)
		}
		maxOffset = XLogSegmentOffset(lsn)
	}

	return common.DecodePGWAL(d, maxOffset)
}
