package pgproee14

import "github.com/wader/fq/pkg/decode"

// type = struct ControlFileData {
/*    0      |     8 */ // uint64 system_identifier;
/*    8      |     4 */ // uint32 pg_control_version;
/*   12      |     4 */ // uint32 catalog_version_no;
/*   16      |     4 */ // DBState state;
/* XXX  4-byte hole  */
/*   24      |     8 */ // pg_time_t time;
/*   32      |     8 */ // XLogRecPtr checkPoint;
/*   40      |   120 */ // CheckPoint checkPointCopy;
/*  160      |     8 */ // XLogRecPtr unloggedLSN;
/*  168      |     8 */ // XLogRecPtr minRecoveryPoint;
/*  176      |     4 */ // TimeLineID minRecoveryPointTLI;
/* XXX  4-byte hole  */
/*  184      |     8 */ // XLogRecPtr backupStartPoint;
/*  192      |     8 */ // XLogRecPtr backupEndPoint;
/*  200      |     1 */ // _Bool backupEndRequired;
/* XXX  3-byte hole  */
/*  204      |     4 */ // int wal_level;
/*  208      |     1 */ // _Bool wal_log_hints;
/* XXX  3-byte hole  */
/*  212      |     4 */ // int MaxConnections;
/*  216      |     4 */ // int max_worker_processes;
/*  220      |     4 */ // int max_wal_senders;
/*  224      |     4 */ // int max_prepared_xacts;
/*  228      |     4 */ // int max_locks_per_xact;
/*  232      |     1 */ // _Bool track_commit_timestamp;
/* XXX  3-byte hole  */
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
/*  280      |     1 */ // _Bool float8ByVal;
/* XXX  3-byte hole  */
/*  284      |     4 */ // uint32 data_checksum_version;
/*  288      |    32 */ // char mock_authentication_nonce[32];
/*  320      |     4 */ // pg_icu_version icu_version;
/*  324      |     4 */ // uint32 pg_old_version;
/*  328      |     4 */ // pg_crc32c crc;
/* XXX  4-byte padding  */
//
/* total size (bytes):  336 */

// type = struct CheckPoint {
/*    0      |     8 */ // XLogRecPtr redo;
/*    8      |     4 */ // TimeLineID ThisTimeLineID;
/*   12      |     4 */ // TimeLineID PrevTimeLineID;
/*   16      |     1 */ // _Bool fullPageWrites;
/* XXX  7-byte hole  */
/*   24      |     8 */ // FullTransactionId nextXid;
/*   32      |     4 */ // Oid nextOid;
/* XXX  4-byte hole  */
/*   40      |     8 */ // MultiXactId nextMulti;
/*   48      |     8 */ // MultiXactOffset nextMultiOffset;
/*   56      |     8 */ // TransactionId oldestXid;
/*   64      |     4 */ // Oid oldestXidDB;
/* XXX  4-byte hole  */
/*   72      |     8 */ // MultiXactId oldestMulti;
/*   80      |     4 */ // Oid oldestMultiDB;
/* XXX  4-byte hole  */
/*   88      |     8 */ // pg_time_t time;
/*   96      |     8 */ // TransactionId oldestCommitTsXid;
/*  104      |     8 */ // TransactionId newestCommitTsXid;
/*  112      |     8 */ // TransactionId oldestActiveXid;
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
	d.FieldU32("pg_control_version")
	d.FieldU32("catalog_version_no")
	d.FieldU32("state")
	d.U32()

	/*   24      |     8 */ // pg_time_t time;
	/*   32      |     8 */ // XLogRecPtr checkPoint;
	/*   40      |   120 */ // CheckPoint checkPointCopy;
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
		/* XXX  4-byte hole  */
		d.FieldU64("nextXid")
		d.FieldU32("nextOid")
		d.U32()

		/*   40      |     8 */ // MultiXactId nextMulti;
		/*   48      |     8 */ // MultiXactOffset nextMultiOffset;
		/*   56      |     8 */ // TransactionId oldestXid;
		/*   64      |     4 */ // Oid oldestXidDB;
		/* XXX  4-byte hole  */
		d.FieldU64("nextMulti")
		d.FieldU64("nextMultiOffset")
		d.FieldU64("oldestXid")
		d.FieldU32("oldestXidDB")
		d.U32()

		/*   72      |     8 */ // MultiXactId oldestMulti;
		/*   80      |     4 */ // Oid oldestMultiDB;
		/* XXX  4-byte hole  */
		d.FieldU64("oldestMulti")
		d.FieldU32("oldestMultiDB")
		d.U32()

		/*   88      |     8 */ // pg_time_t time;
		/*   96      |     8 */ // TransactionId oldestCommitTsXid;
		/*  104      |     8 */ // TransactionId newestCommitTsXid;
		/*  112      |     8 */ // TransactionId oldestActiveXid;
		d.FieldS64("time")
		d.FieldU64("oldestCommitTsXid")
		d.FieldU64("newestCommitTsXid")
		d.FieldU64("oldestActiveXid")
	})

	/*  160      |     8 */ // XLogRecPtr unloggedLSN;
	/*  168      |     8 */ // XLogRecPtr minRecoveryPoint;
	/*  176      |     4 */ // TimeLineID minRecoveryPointTLI;
	/* XXX  4-byte hole  */
	d.FieldU64("unloggedLSN")
	d.FieldU64("minRecoveryPoint")
	d.FieldU32("minRecoveryPointTLI")
	d.U32()

	/*  184      |     8 */ // XLogRecPtr backupStartPoint;
	/*  192      |     8 */ // XLogRecPtr backupEndPoint;
	/*  200      |     1 */ // _Bool backupEndRequired;
	/* XXX  3-byte hole  */
	d.FieldU64("backupStartPoint")
	d.FieldU64("backupEndPoint")
	d.FieldU8("backupEndRequired")
	d.U24()

	/*  204      |     4 */ // int wal_level;
	/*  208      |     1 */ // _Bool wal_log_hints;
	/* XXX  3-byte hole  */
	d.FieldS32("wal_level")
	d.FieldU8("wal_log_hints")
	d.U24()

	/*  212      |     4 */ // int MaxConnections;
	/*  216      |     4 */ // int max_worker_processes;
	/*  220      |     4 */ // int max_wal_senders;
	/*  224      |     4 */ // int max_prepared_xacts;
	/*  228      |     4 */ // int max_locks_per_xact;
	/*  232      |     1 */ // _Bool track_commit_timestamp;
	/* XXX  3-byte hole  */
	d.FieldS32("MaxConnections")
	d.FieldS32("max_worker_processes")
	d.FieldS32("max_wal_senders")
	d.FieldS32("max_prepared_xacts")
	d.FieldS32("max_locks_per_xact")
	d.FieldU8("track_commit_timestamp")
	d.U24()

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
	/*  280      |     1 */ // _Bool float8ByVal;
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

	/*  284      |     4 */ // uint32 data_checksum_version;
	/*  288      |    32 */ // char mock_authentication_nonce[32];
	/*  320      |     4 */ // pg_icu_version icu_version;
	/*  324      |     4 */ // uint32 pg_old_version;
	/*  328      |     4 */ // pg_crc32c crc;
	/* XXX  4-byte padding  */
	d.FieldU32("data_checksum_version")
	d.FieldUTF8ShortStringFixedLen("mock_authentication_nonce", 32)
	d.FieldU32("icu_version")
	d.FieldU32("pg_old_version")
	d.FieldU32("crc")
	d.U32()
	/* total size (bytes):  336 */

	d.AssertPosBytes(336)

	return nil
}
