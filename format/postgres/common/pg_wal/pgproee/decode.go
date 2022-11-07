package pgproee

import (
	"github.com/wader/fq/format/postgres/common/pg_wal/postgres"
	"github.com/wader/fq/pkg/decode"
)

func DecodePGWAL(d *decode.D) any {
	wal := &postgres.Wal{
		DecodeXLogRecord: decodeXLogRecord,
	}
	return postgres.Decode(d, wal)
}
