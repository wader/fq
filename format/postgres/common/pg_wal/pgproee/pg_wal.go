package pgproee

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/common/pg_wal/postgres"
)

func decodeXLogRecord(wal *postgres.Wal, maxBytes int64) {
	record := wal.State.Record

	pos0 := record.Pos()
	maxLen := maxBytes * 8
	if record.FieldGet("xLogBody0") == nil {
		// body on first page
		record.FieldRawLen("xLogBody0", maxLen)
	} else {
		// body on second page
		record.FieldRawLen("xLogBody1", maxLen)
	}
	pos1 := record.Pos()
	posMax := pos1
	record.SeekAbs(pos0)

	// xl_tot_len already read

	if record.FieldGet("hole0") == nil {
		if postgres.IsEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("hole0")
	}

	if record.FieldGet("xl_xid") == nil {
		if postgres.IsEnd(record, posMax, 64) {
			return
		}
		record.FieldU64("xl_xid")
	}

	if record.FieldGet("xl_prev") == nil {
		if postgres.IsEnd(record, posMax, 64) {
			return
		}
		record.FieldU64("xl_prev", common.XLogRecPtrMapper)
	}

	if record.FieldGet("xl_info") == nil {
		if postgres.IsEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_info")
	}

	if record.FieldGet("xl_rmid") == nil {
		if postgres.IsEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_rmid")
	}

	if record.FieldGet("hole1") == nil {
		if postgres.IsEnd(record, posMax, 16) {
			return
		}
		record.FieldU16("hole1")
	}

	if record.FieldGet("xl_crc") == nil {
		if postgres.IsEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("xl_crc")
	}

	record.SeekAbs(posMax)
}
