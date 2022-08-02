package pgproee11

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
/*  220      |     4 */ // int max_prepared_xacts;
/*  224      |     4 */ // int max_locks_per_xact;
/*  228      |     1 */ // _Bool track_commit_timestamp;
/* XXX  3-byte hole */
/*  232      |     4 */ // uint32 maxAlign;
/* XXX  4-byte hole */
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
/*   24      |     8 */ // TransactionId nextXid;
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
