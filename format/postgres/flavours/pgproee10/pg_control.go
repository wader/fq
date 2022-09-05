package pgproee10

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
/*   40      |     8 */ // XLogRecPtr prevCheckPoint;
/*   48      |   120 */ // CheckPoint checkPointCopy;
/*  168      |     8 */ // XLogRecPtr unloggedLSN;
/*  176      |     8 */ // XLogRecPtr minRecoveryPoint;
/*  184      |     4 */ // TimeLineID minRecoveryPointTLI;
/* XXX  4-byte hole */
/*  192      |     8 */ // XLogRecPtr backupStartPoint;
/*  200      |     8 */ // XLogRecPtr backupEndPoint;
/*  208      |     1 */ // bool backupEndRequired;
/* XXX  3-byte hole */
/*  212      |     4 */ // int wal_level
/*  216      |     1 */ // bool wal_log_hints
/* XXX  3-byte hole */
/*  220      |     4 */ // int MaxConnections
/*  224      |     4 */ // int max_worker_processes
/*  228      |     4 */ // int max_prepared_xacts
/*  232      |     4 */ // int max_locks_per_xact
/*  236      |     1 */ // bool track_commit_timestamp
/* XXX  3-byte hole */
/*  240      |     4 */ // uint32 maxAlign
/* XXX  4-byte hole */
/*  248      |     8 */ // double floatFormat
/*  256      |     4 */ // uint32 blcksz
/*  260      |     4 */ // uint32 relseg_size
/*  264      |     4 */ // uint32 xlog_blcksz
/*  268      |     4 */ // uint32 xlog_seg_size
/*  272      |     4 */ // uint32 nameDataLen
/*  276      |     4 */ // uint32 indexMaxKeys
/*  280      |     4 */ // uint32 toast_max_chunk_size
/*  284      |     4 */ // uint32 loblksize
/*  288      |     1 */ // bool float4ByVal
/*  289      |     1 */ // bool float8ByVal
/* XXX  2-byte hole */
/*  292      |     4 */ // uint32 data_checksum_version
/*  296      |    32 */ // char mock_authentication_nonce[32]
/*  328      |     4 */ // pg_icu_version icu_version
/*  332      |     4 */ // uint32 pg_old_version
/*  336      |     4 */ // pg_crc32c crc
/* XXX  4-byte padding */
//
/* total size (bytes):  344 */

// type = struct CheckPoint {
/*    0      |     8 */ // XLogRecPtr redo
/*    8      |     4 */ // TimeLineID ThisTimeLineID
/*   12      |     4 */ // TimeLineID PrevTimeLineID
/*   16      |     1 */ // bool fullPageWrites
/* XXX  7-byte hole */
/*   24      |     8 */ // TransactionId nextXid
/*   32      |     4 */ // Oid nextOid
/* XXX  4-byte hole */
/*   40      |     8 */ // MultiXactId nextMulti
/*   48      |     8 */ // MultiXactOffset nextMultiOffset
/*   56      |     8 */ // TransactionId oldestXid
/*   64      |     4 */ // Oid oldestXidDB
/* XXX  4-byte hole */
/*   72      |     8 */ // MultiXactId oldestMulti
/*   80      |     4 */ // Oid oldestMultiDB
/* XXX  4-byte hole */
/*   88      |     8 */ // pg_time_t time
/*   96      |     8 */ // TransactionId oldestCommitTsXid
/*  104      |     8 */ // TransactionId newestCommitTsXid
/*  112      |     8 */ // TransactionId oldestActiveXid
//
/* total size (bytes):  120 */

func DecodePgControl(d *decode.D, in any) any {
	d.SeekAbs(0)
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

	/* XXX  4-byte hole */
	/*   24      |     8 */ // pg_time_t time;
	/*   32      |     8 */ // XLogRecPtr checkPoint;
	/*   40      |     8 */ // XLogRecPtr prevCheckPoint;
	/*   48      |   120 */ // CheckPoint checkPointCopy;
	d.FieldS64("time", common.TimeMapper)
	d.FieldU64("checkPoint", common.XLogRecPtrMapper)
	d.FieldU64("prevCheckPoint", common.XLogRecPtrMapper)
	d.FieldStruct("checkPointCopy", func(d *decode.D) {
		/*    0      |     8 */ // XLogRecPtr redo;
		/*    8      |     4 */ // TimeLineID ThisTimeLineID;
		/*   12      |     4 */ // TimeLineID PrevTimeLineID;
		/*   16      |     1 */ // _Bool fullPageWrites;
		/* XXX  7-byte hole  */
		d.FieldU64("redo", common.XLogRecPtrMapper)
		d.FieldU32("ThisTimeLineID")
		d.FieldU32("PrevTimeLineID")
		d.FieldU8("fullPageWrites")
		d.FieldU56("hole1")

		/*   24      |     8 */ // FullTransactionId nextXid;
		/*   32      |     4 */ // Oid nextOid;
		/* XXX  4-byte hole  */
		d.FieldU64("nextXid")
		d.FieldU32("nextOid")
		d.FieldU32("hole2")

		/*   40      |     8 */ // MultiXactId nextMulti;
		/*   48      |     8 */ // MultiXactOffset nextMultiOffset;
		/*   56      |     8 */ // TransactionId oldestXid;
		/*   64      |     4 */ // Oid oldestXidDB;
		/* XXX  4-byte hole  */
		d.FieldU64("nextMulti")
		d.FieldU64("nextMultiOffset")
		d.FieldU64("oldestXid")
		d.FieldU32("oldestXidDB")
		d.FieldU32("hole3")

		/*   72      |     8 */ // MultiXactId oldestMulti;
		/*   80      |     4 */ // Oid oldestMultiDB;
		/* XXX  4-byte hole  */
		d.FieldU64("oldestMulti")
		d.FieldU32("oldestMultiDB")
		d.FieldU32("hole4")

		/*   88      |     8 */ // pg_time_t time;
		/*   96      |     8 */ // TransactionId oldestCommitTsXid;
		/*  104      |     8 */ // TransactionId newestCommitTsXid;
		/*  112      |     8 */ // TransactionId oldestActiveXid;
		d.FieldS64("time", common.TimeMapper)
		d.FieldU64("oldestCommitTsXid")
		d.FieldU64("newestCommitTsXid")
		d.FieldU64("oldestActiveXid")
	})

	/*  168      |     8 */ // XLogRecPtr unloggedLSN;
	/*  176      |     8 */ // XLogRecPtr minRecoveryPoint;
	/*  184      |     4 */ // TimeLineID minRecoveryPointTLI;
	/* XXX  4-byte hole */
	d.FieldU64("unloggedLSN", common.LocPtrMapper)
	d.FieldU64("minRecoveryPoint", common.LocPtrMapper)
	d.FieldU32("minRecoveryPointTLI")
	d.FieldU32("hole5")

	/*  192      |     8 */ // XLogRecPtr backupStartPoint;
	/*  200      |     8 */ // XLogRecPtr backupEndPoint;
	/*  208      |     1 */ // bool backupEndRequired;
	/* XXX  3-byte hole */
	d.FieldU64("backupStartPoint", common.LocPtrMapper)
	d.FieldU64("backupEndPoint", common.LocPtrMapper)
	d.FieldU8("backupEndRequired")
	d.FieldU24("hole6")

	/*  212      |     4 */ // int wal_level
	/*  216      |     1 */ // bool wal_log_hints
	/* XXX  3-byte hole */
	d.FieldS32("wal_level", common.WalLevel)
	d.FieldU8("wal_log_hints")
	d.FieldU24("hole7")

	/*  220      |     4 */ // int MaxConnections
	/*  224      |     4 */ // int max_worker_processes
	/*  228      |     4 */ // int max_prepared_xacts
	/*  232      |     4 */ // int max_locks_per_xact
	/*  236      |     1 */ // bool track_commit_timestamp
	/* XXX  3-byte hole */
	d.FieldS32("MaxConnections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.FieldU24("hole8")

	/*  240      |     4 */ // uint32 maxAlign
	/* XXX  4-byte hole */
	d.FieldU32("maxAlign")
	d.FieldU32("hole9")

	/*  248      |     8 */ // double floatFormat
	/*  256      |     4 */ // uint32 blcksz
	/*  260      |     4 */ // uint32 relseg_size
	/*  264      |     4 */ // uint32 xlog_blcksz
	/*  268      |     4 */ // uint32 xlog_seg_size
	/*  272      |     4 */ // uint32 nameDataLen
	/*  276      |     4 */ // uint32 indexMaxKeys
	/*  280      |     4 */ // uint32 toast_max_chunk_size
	/*  284      |     4 */ // uint32 loblksize
	/*  288      |     1 */ // bool float4ByVal
	/*  289      |     1 */ // bool float8ByVal
	/* XXX  2-byte hole */
	d.FieldF64("floatFormat")
	d.FieldU32("blcksz")
	d.FieldU32("relseg_size")
	d.FieldU32("xlog_blcksz")
	d.FieldU32("xlog_seg_size")
	d.FieldU32("nameDataLen")
	d.FieldU32("indexMaxKeys")
	d.FieldU32("toast_max_chunk_size")
	d.FieldU32("loblksize")
	d.FieldU8("float4ByVal")
	d.FieldU8("float8ByVal")
	d.FieldU16("hole10")

	/*  292      |     4 */ // uint32 data_checksum_version
	/*  296      |    32 */ // char mock_authentication_nonce[32]
	/*  328      |     4 */ // pg_icu_version icu_version
	/*  332      |     4 */ // uint32 pg_old_version
	/*  336      |     4 */ // pg_crc32c crc
	/* XXX  4-byte padding */
	d.FieldU32("data_checksum_version")
	d.FieldRawLen("mock_authentication_nonce", 32*8, scalar.RawHex)
	d.FieldU32("icu_version", common.IcuVersionMapper)
	d.FieldU32("pg_old_version")
	d.FieldU32("crc")
	d.FieldU32("hole11")
	/* total size (bytes):  344 */

	d.AssertPosBytes(344)
	d.FieldRawLen("unused", d.BitsLeft())

	return nil
}
