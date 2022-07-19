package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_CONTROL,
		Description: "PostgreSQL 14 control file",
		DecodeFn:    decodePgControl,
	})
}

const (
	PG_CONTROL_VERSION_11 = 1100
	PG_CONTROL_VERSION_14 = 1300
)

func decodePgControl(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	d.FieldU64("system_identifier")
	pgControlVersion := d.FieldU32("pg_control_version")

	switch pgControlVersion {
	case PG_CONTROL_VERSION_11:
		return decodePgControl11(d, in)
	case PG_CONTROL_VERSION_14:
		return decodePgControl14(d, in)
	default:
		d.Fatalf("unsupported PG_CONTROL_VERSION = %d\n", pgControlVersion)
	}
	return nil
}

// type = struct ControlFileData {
/*    0      |     8 */ // uint64 system_identifier;
/*    8      |     4 */ // uint32 pg_control_version;
/*   12      |     4 */ // uint32 catalog_version_no;
/*   16      |     4 */ // DBState state;
/* XXX  4-byte hole */
/*   24      |     8 */ // pg_time_t time;
/*   32      |     8 */ // XLogRecPtr checkPoint;
/*   40      |    80 */ // CheckPoint checkPointCopy;
/*  120      |     8 */ // XLogRecPtr unloggedLSN;
/*  128      |     8 */ // XLogRecPtr minRecoveryPoint;
/*  136      |     4 */ // TimeLineID minRecoveryPointTLI;
/* XXX  4-byte hole */
/*  144      |     8 */ // XLogRecPtr backupStartPoint;
/*  152      |     8 */ // XLogRecPtr backupEndPoint;
/*  160      |     1 */ // _Bool backupEndRequired;
/* XXX  3-byte hole */
/*  164      |     4 */ // int wal_level;
/*  168      |     1 */ // _Bool wal_log_hints;
/* XXX  3-byte hole */
/*  172      |     4 */ // int MaxConnections;
/*  176      |     4 */ // int max_worker_processes;
/*  180      |     4 */ // int max_prepared_xacts;
/*  184      |     4 */ // int max_locks_per_xact;
/*  188      |     1 */ // _Bool track_commit_timestamp;
/* XXX  3-byte hole */
/*  192      |     4 */ // uint32 maxAlign;
/* XXX  4-byte hole */
/*  200      |     8 */ // double floatFormat;
/*  208      |     4 */ // uint32 blcksz;
/*  212      |     4 */ // uint32 relseg_size;
/*  216      |     4 */ // uint32 xlog_blcksz;
/*  220      |     4 */ // uint32 xlog_seg_size;
/*  224      |     4 */ // uint32 nameDataLen;
/*  228      |     4 */ // uint32 indexMaxKeys;
/*  232      |     4 */ // uint32 toast_max_chunk_size;
/*  236      |     4 */ // uint32 loblksize;
/*  240      |     1 */ // _Bool float4ByVal;
/*  241      |     1 */ // _Bool float8ByVal;
/* XXX  2-byte hole */
/*  244      |     4 */ // uint32 data_checksum_version;
/*  248      |    32 */ // char mock_authentication_nonce[32];
/*  280      |     4 */ // pg_crc32c crc;
/* XXX  4-byte padding */
//
/* total size (bytes):  288 */

// type = struct CheckPoint {
/*    0      |     8 */ // XLogRecPtr redo;
/*    8      |     4 */ // TimeLineID ThisTimeLineID;
/*   12      |     4 */ // TimeLineID PrevTimeLineID;
/*   16      |     1 */ // _Bool fullPageWrites;
/* XXX  3-byte hole */
/*   20      |     4 */ // uint32 nextXidEpoch;
/*   24      |     4 */ // TransactionId nextXid;
/*   28      |     4 */ // Oid nextOid;
/*   32      |     4 */ // MultiXactId nextMulti;
/*   36      |     4 */ // MultiXactOffset nextMultiOffset;
/*   40      |     4 */ // TransactionId oldestXid;
/*   44      |     4 */ // Oid oldestXidDB;
/*   48      |     4 */ // MultiXactId oldestMulti;
/*   52      |     4 */ // Oid oldestMultiDB;
/*   56      |     8 */ // pg_time_t time;
/*   64      |     4 */ // TransactionId oldestCommitTsXid;
/*   68      |     4 */ // TransactionId newestCommitTsXid;
/*   72      |     4 */ // TransactionId oldestActiveXid;
/* XXX  4-byte padding */
//
/* total size (bytes):   80 */
//
func decodePgControl11(d *decode.D, in any) any {
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	/*   12      |     4 */ // uint32 catalog_version_no;
	/*   16      |     4 */ // DBState state;
	/* XXX  4-byte hole  */
	//d.FieldU64("system_identifier")
	//d.FieldU32("pg_control_version")
	d.FieldU32("catalog_version_no")
	d.FieldU32("state")
	d.U32()

	/*   24      |     8 */ // pg_time_t time;
	/*   32      |     8 */ // XLogRecPtr checkPoint;
	/*   40      |    80 */ // CheckPoint checkPointCopy;
	d.FieldS64("time")
	d.FieldU64("checkPoint")
	d.FieldStruct("checkPointCopy", func(d *decode.D) {
		/*    0      |     8 */ // XLogRecPtr redo;
		/*    8      |     4 */ // TimeLineID ThisTimeLineID;
		/*   12      |     4 */ // TimeLineID PrevTimeLineID;
		/*   16      |     1 */ // _Bool fullPageWrites;
		/* XXX  3-byte hole */
		d.FieldU64("redo")
		d.FieldU32("ThisTimeLineID")
		d.FieldU32("PrevTimeLineID")
		d.FieldU8("fullPageWrites")
		d.U24()

		/*   20      |     4 */ // uint32 nextXidEpoch;
		/*   24      |     4 */ // TransactionId nextXid;
		/*   28      |     4 */ // Oid nextOid;
		/*   32      |     4 */ // MultiXactId nextMulti;
		/*   36      |     4 */ // MultiXactOffset nextMultiOffset;
		/*   40      |     4 */ // TransactionId oldestXid;
		/*   44      |     4 */ // Oid oldestXidDB;
		/*   48      |     4 */ // MultiXactId oldestMulti;
		/*   52      |     4 */ // Oid oldestMultiDB;
		/*   56      |     8 */ // pg_time_t time;
		/*   64      |     4 */ // TransactionId oldestCommitTsXid;
		/*   68      |     4 */ // TransactionId newestCommitTsXid;
		/*   72      |     4 */ // TransactionId oldestActiveXid;
		/* XXX  4-byte padding */
		d.FieldU32("nextXidEpoch")
		d.FieldU32("nextXid")
		d.FieldU32("nextOid")
		d.FieldU32("nextMulti")
		d.FieldU32("nextMultiOffset")
		d.FieldU32("oldestXid")
		d.FieldU32("oldestXidDB")
		d.FieldU32("oldestMulti")
		d.FieldU32("oldestMultiDB")
		d.FieldS64("time")
		d.FieldU32("oldestCommitTsXid")
		d.FieldU32("newestCommitTsXid")
		d.FieldU32("oldestActiveXid")
		d.U32()
	})

	/*  120      |     8 */ // XLogRecPtr unloggedLSN;
	/*  128      |     8 */ // XLogRecPtr minRecoveryPoint;
	/*  136      |     4 */ // TimeLineID minRecoveryPointTLI;
	/* XXX  4-byte hole */
	d.FieldU64("unloggedLSN")
	d.FieldU64("minRecoveryPoint")
	d.FieldU32("minRecoveryPointTLI")
	d.U32()

	/*  144      |     8 */ // XLogRecPtr backupStartPoint;
	/*  152      |     8 */ // XLogRecPtr backupEndPoint;
	/*  160      |     1 */ // _Bool backupEndRequired;
	/* XXX  3-byte hole */
	d.FieldU64("backupStartPoint")
	d.FieldU64("backupEndPoint")
	d.FieldU8("backupEndRequired")
	d.U24()

	/*  164      |     4 */ // int wal_level;
	/*  168      |     1 */ // _Bool wal_log_hints;
	/* XXX  3-byte hole */
	d.FieldS32("wal_level")
	d.FieldU8("wal_log_hints")
	d.U24()

	/*  172      |     4 */ // int MaxConnections;
	/*  176      |     4 */ // int max_worker_processes;
	/*  180      |     4 */ // int max_prepared_xacts;
	/*  184      |     4 */ // int max_locks_per_xact;
	/*  188      |     1 */ // _Bool track_commit_timestamp;
	/* XXX  3-byte hole  */
	d.FieldS32("MaxConnections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.U24()

	/*  192      |     4 */ // uint32 maxAlign;
	/* XXX  4-byte hole */
	d.FieldU32("maxAlign")
	d.U32()

	/*  200      |     8 */ // double floatFormat;
	/*  208      |     4 */ // uint32 blcksz;
	/*  212      |     4 */ // uint32 relseg_size;
	/*  216      |     4 */ // uint32 xlog_blcksz;
	/*  220      |     4 */ // uint32 xlog_seg_size;
	/*  224      |     4 */ // uint32 nameDataLen;
	/*  228      |     4 */ // uint32 indexMaxKeys;
	/*  232      |     4 */ // uint32 toast_max_chunk_size;
	/*  236      |     4 */ // uint32 loblksize;
	/*  240      |     1 */ // _Bool float4ByVal;
	/*  241      |     1 */ // _Bool float8ByVal;
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
	d.U16()
	
	/*  252      |     4 */ // uint32 data_checksum_version;
	/*  256      |    32 */ // char mock_authentication_nonce[32];
	/*  288      |     4 */ // pg_crc32c crc;
	/* XXX  4-byte padding  */
	d.FieldU32("data_checksum_version")
	d.FieldUTF8ShortStringFixedLen("mock_authentication_nonce", 32)
	d.FieldU32("crc")
	d.U32()
	/* total size (bytes):  288 */

	d.AssertPosBytes(288)

	return nil
}

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
func decodePgControl14(d *decode.D, in any) any {
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	/*   12      |     4 */ // uint32 catalog_version_no;
	/*   16      |     4 */ // DBState state;
	/* XXX  4-byte hole  */
	//d.FieldU64("system_identifier")
	//d.FieldU32("pg_control_version")
	d.FieldU32("catalog_version_no")
	d.FieldU32("state")
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
	d.FieldS32("wal_level")
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
	d.FieldUTF8ShortStringFixedLen("mock_authentication_nonce", 32)
	d.FieldU32("crc")
	d.U32()
	/* total size (bytes):  296 */

	d.AssertPosBytes(296)

	return nil
}
