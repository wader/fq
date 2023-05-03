package pgpro10

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func DecodePgControl(d *decode.D) any {
	d.FieldU64("system_identifier")
	d.FieldU32("pg_control_version", common.VersionMapper)
	d.FieldU32("catalog_version_no")
	d.FieldU32("state", common.DBState)
	d.FieldU32("hole0")

	d.FieldS64("time", common.TimeMapper)
	d.FieldU64("check_point", common.XLogRecPtrMapper)
	d.FieldU64("prev_check_point", common.XLogRecPtrMapper)
	d.FieldStruct("check_point_copy", func(d *decode.D) {

		d.FieldU64("redo", common.XLogRecPtrMapper)
		d.FieldU32("this_time_line_id")
		d.FieldU32("prev_time_line_id")
		d.FieldU8("full_page_writes")
		d.FieldU24("hole1")

		d.FieldU32("next_xid_epoch")
		d.FieldU32("next_xid")
		d.FieldU32("next_oid")
		d.FieldU32("next_multi")
		d.FieldU32("next_multi_offset")
		d.FieldU32("oldest_xid")
		d.FieldU32("oldest_xid_db")
		d.FieldU32("oldest_multi")
		d.FieldU32("oldest_multi_db")
		d.FieldS64("time", common.TimeMapper)
		d.FieldU32("oldest_commit_ts_xid")
		d.FieldU32("newest_commit_ts_xid")
		d.FieldU32("oldest_active_xid")
		d.FieldU32("padding1")
	})

	d.FieldU64("unlogged_lsn", common.LocPtrMapper)
	d.FieldU64("min_recovery_point", common.LocPtrMapper)
	d.FieldU32("min_recovery_point_tli")
	d.FieldU32("hole2")

	d.FieldU64("backup_start_point", common.LocPtrMapper)
	d.FieldU64("backup_end_point", common.LocPtrMapper)
	d.FieldU8("backup_end_required")
	d.FieldU24("hole3")

	d.FieldS32("wal_level", common.WalLevel)
	d.FieldU8("wal_log_hints")
	d.FieldU24("hole4")

	d.FieldS32("max_connections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.FieldU24("hole5")

	d.FieldU32("max_align")
	d.FieldU32("hole6")

	d.FieldF64("float_format")
	d.FieldU32("blcksz")
	d.FieldU32("relseg_size")
	d.FieldU32("xlog_blcksz")
	d.FieldU32("xlog_seg_size")
	d.FieldU32("name_data_len")
	d.FieldU32("index_max_keys")
	d.FieldU32("toast_max_chunk_size")
	d.FieldU32("loblksize")
	d.FieldU8("float4_by_val")
	d.FieldU8("float8_by_val")
	d.FieldU16("hole7")

	d.FieldU32("data_checksum_version")
	d.FieldRawLen("mock_authentication_nonce", 32*8, scalar.RawHex)
	d.FieldU32("icu_version", common.IcuVersionMapper)
	d.FieldU32("crc")

	d.AssertPos(296 * 8)
	d.FieldRawLen("unused", d.BitsLeft())

	return nil
}
