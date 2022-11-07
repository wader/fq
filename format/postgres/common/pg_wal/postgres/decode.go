package postgres

import "github.com/wader/fq/pkg/decode"

func DecodePGWAL(d *decode.D) any {
	wal := &Wal{
		DecodeXLogRecord: decodeXLogRecord,
	}
	return Decode(d, wal)
}
