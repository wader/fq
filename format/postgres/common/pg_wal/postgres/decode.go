package postgres

import "github.com/wader/fq/pkg/decode"

func DecodePGWAL(d *decode.D, maxOffset uint32) any {
	wal := &Wal{
		MaxOffset:        int64(maxOffset),
		DecodeXLogRecord: decodeXLogRecord,
	}
	return Decode(d, wal)
}
