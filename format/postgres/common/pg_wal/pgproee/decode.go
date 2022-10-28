package pgproee

import (
	"github.com/wader/fq/format/postgres/common/pg_wal/postgres"
	"github.com/wader/fq/pkg/decode"
)

func DecodePGWAL(d *decode.D, maxOffset uint32) any {
	wal := &postgres.Wal{
		MaxOffset:        int64(maxOffset),
		DecodeXLogRecord: decodeXLogRecord,
	}
	return postgres.Decode(d, wal)
}
