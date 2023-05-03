package pgproee12

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// type = struct ControlFileData {
/*    0      |     8 */ // uint64 system_identifier;
/*    8      |     4 */ // uint32 pg_control_version;
/*   12      |     4 */ // uint32 catalog_version_no;
/*   16      |     4 */ // DBState state;
/* XXX  4-byte hole */
/*   24      |     8 */ // pg_time_t time;
/*   32      |     8 */ // XLogRecPtr checkPoint;
/*   40      |   120 */ // CheckPoint checkPointCopy;
/*  160      |     8 */ // XLogRecPtr unloggedLSN;
/*  168      |     8 */ // XLogRecPtr minRecoveryPoint;
/*  176      |     4 */ // TimeLineID minRecoveryPointTLI;
/* XXX  4-byte hole */
/*  184      |     8 */ // XLogRecPtr backupStartPoint;
/*  192      |     8 */ // XLogRecPtr backupEndPoint;
/*  200      |     1 */ // _Bool backupEndRequired;
/* XXX  3-byte hole */
/*  204      |     4 */ // int wal_level;
/*  208      |     1 */ // _Bool wal_log_hints;
/* XXX  3-byte hole */
/*  212      |     4 */ // int MaxConnections;
/*  216      |     4 */ // int max_worker_processes;
/*  220      |     4 */ // int max_wal_senders;
/*  224      |     4 */ // int max_prepared_xacts;
/*  228      |     4 */ // int max_locks_per_xact;
/*  232      |     1 */ // _Bool track_commit_timestamp;
/* XXX  3-byte hole */
/*  236      |     4 */ // uint32 maxAlign;
/*  240      |     8 */ // double floatFormat;
/*  248      |     4 */ // uint32 blcksz;
/*  252      |     4 */ // uint32 relseg_size;
/*  256      |     4 */ // uint32 xlog_blcksz;
/*  260      |     4 */ // uint32 xlog_seg_size;
/*  264      |     4 */ // uint32 nameDataLen;
/*  268      |     4 */ // uint32 indexMaxKeys;
/*  272      |     4 */ // uint32 toast_max_chunk_size;
/*  276      |     4 */ // uint32 loblksize;
/*  280      |     1 */ // _Bool float4ByVal;
/*  281      |     1 */ // _Bool float8ByVal;
/* XXX  2-byte hole */
/*  284      |     4 */ // uint32 data_checksum_version;
/*  288      |    32 */ // char mock_authentication_nonce[32];
/*  320      |     4 */ // pg_icu_version icu_version;
/*  324      |     4 */ // uint32 pg_old_version;
/*  328      |     4 */ // SnapshotId oldest_snapshot;
/*  332      |     4 */ // SnapshotId recent_snapshot;
/*  336      |     4 */ // SnapshotId active_snapshot;
/*  340      |     4 */ // pg_crc32c crc;
//
/* total size (bytes):  344 */

// type = struct CheckPoint {
/*    0      |     8 */ // XLogRecPtr redo;
/*    8      |     4 */ // TimeLineID ThisTimeLineID;
/*   12      |     4 */ // TimeLineID PrevTimeLineID;
/*   16      |     1 */ // _Bool fullPageWrites;
/* XXX  7-byte hole */
/*   24      |     8 */ // FullTransactionId nextFullXid;
/*   32      |     4 */ // Oid nextOid;
/* XXX  4-byte hole */
/*   40      |     8 */ // MultiXactId nextMulti;
/*   48      |     8 */ // MultiXactOffset nextMultiOffset;
/*   56      |     8 */ // TransactionId oldestXid;
/*   64      |     4 */ // Oid oldestXidDB;
/* XXX  4-byte hole */
/*   72      |     8 */ // MultiXactId oldestMulti;
/*   80      |     4 */ // Oid oldestMultiDB;
/* XXX  4-byte hole */
/*   88      |     8 */ // pg_time_t time;
/*   96      |     8 */ // TransactionId oldestCommitTsXid;
/*  104      |     8 */ // TransactionId newestCommitTsXid;
/*  112      |     8 */ // TransactionId oldestActiveXid;
//
/* total size (bytes):  120 */

func DecodePgControl(d *decode.D) any {
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	/*   12      |     4 */ // uint32 catalog_version_no;
	/*   16      |     4 */ // DBState state;
	/* XXX  4-byte hole  */
	d.FieldU64("system_identifier")
	d.FieldU32("pg_control_version", common.VersionMapper)
	d.FieldU32("catalog_version_no")
	d.FieldU32("state", common.DBState)
	d.FieldU32("hole0")

	/*   24      |     8 */ // pg_time_t time;
	/*   32      |     8 */ // XLogRecPtr checkPoint;
	/*   40      |   120 */ // CheckPoint checkPointCopy;
	d.FieldS64("time", common.TimeMapper)
	d.FieldU64("check_point", common.XLogRecPtrMapper)
	d.FieldStruct("check_point_copy", func(d *decode.D) {
		/*    0      |     8 */ // XLogRecPtr redo;
		/*    8      |     4 */ // TimeLineID ThisTimeLineID;
		/*   12      |     4 */ // TimeLineID PrevTimeLineID;
		/*   16      |     1 */ // _Bool fullPageWrites;
		/* XXX  7-byte hole  */
		d.FieldU64("redo", common.XLogRecPtrMapper)
		d.FieldU32("this_time_line_id")
		d.FieldU32("prev_time_line_id")
		d.FieldU8("full_page_writes")
		d.FieldU56("hole1")

		/*   24      |     8 */ // FullTransactionId nextXid;
		/*   32      |     4 */ // Oid nextOid;
		/* XXX  4-byte hole  */
		d.FieldU64("next_xid")
		d.FieldU32("next_oid")
		d.FieldU32("hole2")

		/*   40      |     8 */ // MultiXactId nextMulti;
		/*   48      |     8 */ // MultiXactOffset nextMultiOffset;
		/*   56      |     8 */ // TransactionId oldestXid;
		/*   64      |     4 */ // Oid oldestXidDB;
		/* XXX  4-byte hole  */
		d.FieldU64("next_multi")
		d.FieldU64("next_multi_offset")
		d.FieldU64("oldest_xid")
		d.FieldU32("oldest_xid_db")
		d.FieldU32("hole3")

		/*   72      |     8 */ // MultiXactId oldestMulti;
		/*   80      |     4 */ // Oid oldestMultiDB;
		/* XXX  4-byte hole  */
		d.FieldU64("oldest_multi")
		d.FieldU32("oldest_multi_db")
		d.FieldU32("hole4")

		/*   88      |     8 */ // pg_time_t time;
		/*   96      |     8 */ // TransactionId oldestCommitTsXid;
		/*  104      |     8 */ // TransactionId newestCommitTsXid;
		/*  112      |     8 */ // TransactionId oldestActiveXid;
		d.FieldS64("time", common.TimeMapper)
		d.FieldU64("oldest_commit_ts_xid")
		d.FieldU64("newest_commit_ts_xid")
		d.FieldU64("oldest_active_xid")
	})

	/*  160      |     8 */ // XLogRecPtr unloggedLSN;
	/*  168      |     8 */ // XLogRecPtr minRecoveryPoint;
	/*  176      |     4 */ // TimeLineID minRecoveryPointTLI;
	/* XXX  4-byte hole  */
	d.FieldU64("unlogged_lsn", common.LocPtrMapper)
	d.FieldU64("min_recovery_point", common.LocPtrMapper)
	d.FieldU32("min_recovery_point_tli")
	d.FieldU32("hole5")

	/*  184      |     8 */ // XLogRecPtr backupStartPoint;
	/*  192      |     8 */ // XLogRecPtr backupEndPoint;
	/*  200      |     1 */ // _Bool backupEndRequired;
	/* XXX  3-byte hole  */
	d.FieldU64("backup_start_point", common.LocPtrMapper)
	d.FieldU64("backup_end_point", common.LocPtrMapper)
	d.FieldU8("backup_end_required")
	d.FieldU24("hole6")

	/*  204      |     4 */ // int wal_level;
	/*  208      |     1 */ // _Bool wal_log_hints;
	/* XXX  3-byte hole  */
	d.FieldS32("wal_level", common.WalLevel)
	d.FieldU8("wal_log_hints")
	d.FieldU24("hole7")

	/*  212      |     4 */ // int MaxConnections;
	/*  216      |     4 */ // int max_worker_processes;
	/*  220      |     4 */ // int max_wal_senders;
	/*  224      |     4 */ // int max_prepared_xacts;
	/*  228      |     4 */ // int max_locks_per_xact;
	/*  232      |     1 */ // _Bool track_commit_timestamp;
	/* XXX  3-byte hole */
	d.FieldS32("max_connections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_wal_senders")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.FieldU24("hole8")

	/*  236      |     4 */ // uint32 maxAlign;
	/*  240      |     8 */ // double floatFormat;
	/*  248      |     4 */ // uint32 blcksz;
	/*  252      |     4 */ // uint32 relseg_size;
	/*  256      |     4 */ // uint32 xlog_blcksz;
	/*  260      |     4 */ // uint32 xlog_seg_size;
	/*  264      |     4 */ // uint32 nameDataLen;
	/*  268      |     4 */ // uint32 indexMaxKeys;
	/*  272      |     4 */ // uint32 toast_max_chunk_size;
	/*  276      |     4 */ // uint32 loblksize;
	/*  280      |     1 */ // _Bool float4ByVal;
	/*  281      |     1 */ // _Bool float8ByVal;
	/* XXX  2-byte hole */
	d.FieldU32("max_align")
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
	d.FieldU16("hole9")

	/*  284      |     4 */ // uint32 data_checksum_version;
	/*  288      |    32 */ // char mock_authentication_nonce[32];
	/*  320      |     4 */ // pg_icu_version icu_version;
	/*  324      |     4 */ // uint32 pg_old_version;
	/*  328      |     4 */ // SnapshotId oldest_snapshot;
	/*  332      |     4 */ // SnapshotId recent_snapshot;
	/*  336      |     4 */ // SnapshotId active_snapshot;
	/*  340      |     4 */ // pg_crc32c crc;
	d.FieldU32("data_checksum_version")
	d.FieldRawLen("mock_authentication_nonce", 32*8, scalar.RawHex)
	d.FieldU32("icu_version", common.IcuVersionMapper)
	d.FieldU32("pg_old_version")
	d.FieldU32("oldest_snapshot")
	d.FieldU32("recent_snapshot")
	d.FieldU32("active_snapshot")
	d.FieldU32("crc")
	/* total size (bytes):  344 */

	d.AssertPos(344 * 8)
	d.FieldRawLen("unused", d.BitsLeft())

	return nil
}
