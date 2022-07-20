package postgres14

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
/* XXX  4-byte hole  */
/*   24      |     8 */ // pg_time_t time;
/*   32      |     8 */ // XLogRecPtr checkPoint;
/*   40      |    88 */ // CheckPoint checkPointCopy;
/*  128      |     8 */ // XLogRecPtr unloggedLSN;
/*  136      |     8 */ // XLogRecPtr minRecoveryPoint;
/*  144      |     4 */ // TimeLineID minRecoveryPointTLI;
/* XXX  4-byte hole  */
/*  152      |     8 */ // XLogRecPtr backupStartPoint;
/*  160      |     8 */ // XLogRecPtr backupEndPoint;
/*  168      |     1 */ // _Bool backupEndRequired;
/* XXX  3-byte hole  */
/*  172      |     4 */ // int wal_level;
/*  176      |     1 */ // _Bool wal_log_hints;
/* XXX  3-byte hole  */
/*  180      |     4 */ // int MaxConnections;
/*  184      |     4 */ // int max_worker_processes;
/*  188      |     4 */ // int max_wal_senders;
/*  192      |     4 */ // int max_prepared_xacts;
/*  196      |     4 */ // int max_locks_per_xact;
/*  200      |     1 */ // _Bool track_commit_timestamp;
/* XXX  3-byte hole  */
/*  204      |     4 */ // uint32 maxAlign;
/*  208      |     8 */ // double floatFormat;
/*  216      |     4 */ // uint32 blcksz;
/*  220      |     4 */ // uint32 relseg_size;
/*  224      |     4 */ // uint32 xlog_blcksz;
/*  228      |     4 */ // uint32 xlog_seg_size;
/*  232      |     4 */ // uint32 nameDataLen;
/*  236      |     4 */ // uint32 indexMaxKeys;
/*  240      |     4 */ // uint32 toast_max_chunk_size;
/*  244      |     4 */ // uint32 loblksize;
/*  248      |     1 */ // _Bool float8ByVal;
/* XXX  3-byte hole  */
/*  252      |     4 */ // uint32 data_checksum_version;
/*  256      |    32 */ // char mock_authentication_nonce[32];
/*  288      |     4 */ // pg_crc32c crc;
/* XXX  4-byte padding  */
//
/* total size (bytes):  296 */
//
// type = struct CheckPoint {
/*    0      |     8 */ // XLogRecPtr redo;
/*    8      |     4 */ // TimeLineID ThisTimeLineID;
/*   12      |     4 */ // TimeLineID PrevTimeLineID;
/*   16      |     1 */ // _Bool fullPageWrites;
/* XXX  7-byte hole  */
/*   24      |     8 */ // FullTransactionId nextXid;
/*   32      |     4 */ // Oid nextOid;
/*   36      |     4 */ // MultiXactId nextMulti;
/*   40      |     4 */ // MultiXactOffset nextMultiOffset;
/*   44      |     4 */ // TransactionId oldestXid;
/*   48      |     4 */ // Oid oldestXidDB;
/*   52      |     4 */ // MultiXactId oldestMulti;
/*   56      |     4 */ // Oid oldestMultiDB;
/* XXX  4-byte hole  */
/*   64      |     8 */ // pg_time_t time;
/*   72      |     4 */ // TransactionId oldestCommitTsXid;
/*   76      |     4 */ // TransactionId newestCommitTsXid;
/*   80      |     4 */ // TransactionId oldestActiveXid;
/* XXX  4-byte padding  */
//
/* total size (bytes):   88 */
func DecodePgControl(d *decode.D, in any) any {
	d.SeekAbs(0)
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	/*   12      |     4 */ // uint32 catalog_version_no;
	/*   16      |     4 */ // DBState state;
	/* XXX  4-byte hole  */
	d.FieldU64("system_identifier")
	d.FieldU32("pg_control_version")
	d.FieldU32("catalog_version_no")
	d.FieldU32("state", common.DBState)
	d.U32()

	/*   24      |     8 */ // pg_time_t time;
	/*   32      |     8 */ // XLogRecPtr checkPoint;
	/*   40      |    88 */ // CheckPoint checkPointCopy;
	d.FieldS64("time")
	d.FieldU64("checkPoint")
	d.FieldStruct("checkPointCopy", func(d *decode.D) {
		/*    0      |     8 */ // XLogRecPtr redo;
		/*    8      |     4 */ // TimeLineID ThisTimeLineID;
		/*   12      |     4 */ // TimeLineID PrevTimeLineID;
		/*   16      |     1 */ // _Bool fullPageWrites;
		/* XXX  7-byte hole  */
		d.FieldU64("redo")
		d.FieldU32("ThisTimeLineID")
		d.FieldU32("PrevTimeLineID")
		d.FieldU8("fullPageWrites")
		d.U56()

		/*   24      |     8 */ // FullTransactionId nextXid;
		/*   32      |     4 */ // Oid nextOid;
		/*   36      |     4 */ // MultiXactId nextMulti;
		/*   40      |     4 */ // MultiXactOffset nextMultiOffset;
		/*   44      |     4 */ // TransactionId oldestXid;
		/*   48      |     4 */ // Oid oldestXidDB;
		/*   52      |     4 */ // MultiXactId oldestMulti;
		/*   56      |     4 */ // Oid oldestMultiDB;
		/* XXX  4-byte hole  */
		d.FieldU64("nextXid")
		d.FieldU32("nextOid")
		d.FieldU32("nextMulti")
		d.FieldU32("nextMultiOffset")
		d.FieldU32("oldestXid")
		d.FieldU32("oldestXidDB")
		d.FieldU32("oldestMulti")
		d.FieldU32("oldestMultiDB")
		d.U32()

		/*   64      |     8 */ // pg_time_t time;
		/*   72      |     4 */ // TransactionId oldestCommitTsXid;
		/*   76      |     4 */ // TransactionId newestCommitTsXid;
		/*   80      |     4 */ // TransactionId oldestActiveXid;
		/* XXX  4-byte padding  */
		d.FieldS64("time")
		d.FieldU32("oldestCommitTsXid")
		d.FieldU32("newestCommitTsXid")
		d.FieldU32("oldestActiveXid")
		d.U32()
	})

	/*  128      |     8 */ // XLogRecPtr unloggedLSN;
	/*  136      |     8 */ // XLogRecPtr minRecoveryPoint;
	/*  144      |     4 */ // TimeLineID minRecoveryPointTLI;
	/* XXX  4-byte hole  */
	d.FieldU64("unloggedLSN")
	d.FieldU64("minRecoveryPoint")
	d.FieldU32("minRecoveryPointTLI")
	d.U32()

	/*  152      |     8 */ // XLogRecPtr backupStartPoint;
	/*  160      |     8 */ // XLogRecPtr backupEndPoint;
	/*  168      |     1 */ // _Bool backupEndRequired;
	/* XXX  3-byte hole  */
	d.FieldU64("backupStartPoint")
	d.FieldU64("backupEndPoint")
	d.FieldU8("backupEndRequired")
	d.U24()

	/*  172      |     4 */ // int wal_level;
	/*  176      |     1 */ // _Bool wal_log_hints;
	/* XXX  3-byte hole  */
	d.FieldS32("wal_level", common.WalLevel)
	d.FieldU8("wal_log_hints")
	d.U24()

	/*  180      |     4 */ // int MaxConnections;
	/*  184      |     4 */ // int max_worker_processes;
	/*  188      |     4 */ // int max_wal_senders;
	/*  192      |     4 */ // int max_prepared_xacts;
	/*  196      |     4 */ // int max_locks_per_xact;
	/*  200      |     1 */ // _Bool track_commit_timestamp;
	/* XXX  3-byte hole  */
	d.FieldS32("MaxConnections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_wal_senders")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.U24()

	/*  204      |     4 */ // uint32 maxAlign;
	/*  208      |     8 */ // double floatFormat;
	/*  216      |     4 */ // uint32 blcksz;
	/*  220      |     4 */ // uint32 relseg_size;
	/*  224      |     4 */ // uint32 xlog_blcksz;
	/*  228      |     4 */ // uint32 xlog_seg_size;
	/*  232      |     4 */ // uint32 nameDataLen;
	/*  236      |     4 */ // uint32 indexMaxKeys;
	/*  240      |     4 */ // uint32 toast_max_chunk_size;
	/*  244      |     4 */ // uint32 loblksize;
	/*  248      |     1 */ // _Bool float8ByVal;
	/* XXX  3-byte hole  */
	d.FieldU32("maxAlign")
	d.FieldF64("floatFormat")
	d.FieldU32("blcksz")
	d.FieldU32("relseg_size")
	d.FieldU32("xlog_blcksz")
	d.FieldU32("xlog_seg_size")
	d.FieldU32("nameDataLen")
	d.FieldU32("indexMaxKeys")
	d.FieldU32("toast_max_chunk_size")
	d.FieldU32("loblksize")
	d.FieldU8("float8ByVal")
	d.U24()

	/*  252      |     4 */ // uint32 data_checksum_version;
	/*  256      |    32 */ // char mock_authentication_nonce[32];
	/*  288      |     4 */ // pg_crc32c crc;
	/* XXX  4-byte padding  */
	d.FieldU32("data_checksum_version")
	d.FieldRawLen("mock_authentication_nonce", 32*8, scalar.RawHex)
	d.FieldU32("crc")
	d.U32()
	/* total size (bytes):  296 */

	d.AssertPosBytes(296)

	return nil
}
