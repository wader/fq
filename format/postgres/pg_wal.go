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

// partial parsing of WAL

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
