## PostgreSQL formats

### Supported PostgreSQL data types (`fq -d pg_heap`):
 - `pg_heap` - heap, table, relation
 - `pg_control` - control file data
 - `pg_wal` - wal, write ahead log (not implemented yet)
 - `pg_btree` - btree index (not implemented yet)
 
### Supported PostgreSQL flavours (`fq -o flavour=postgres14`):
1. Standart (vanilla) PostgreSQL 
   - `postgres11`
   - `postgres12`
   - `postgres13`
   - `postgres14`
2. Postgres Pro Enterprise
   - `pgproee10`
   - `pgproee11`
   - `pgproee12`
   - `pgproee13`
   - `pgproee14`
3. Postgres Professional Standard
   - `pgpro11`
   - `pgpro12`
   - `pgpro13`
   - `pgpro14`

### Supported OS
Postgres for x64 Linux deb, rpm based OS.

### How to run?
Need to specify format, flavour for file and expression:
```
fq -d pg_control -o flavour=postgres14 d pg_control
fq -d pg_heap -o flavour=postgres14 ".[]" 16397
```

### pg_control file format

Source code:
 - [pg_control.h](https://github.com/postgres/postgres/blob/master/src/include/catalog/pg_control.h)

To see content of pg_control run:
```
$ fq -d pg_control -o flavour=postgres14 d pg_control
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|.{}: pg_control (pg_control)
0x000|68 e7 dd 05 b3 88 b3 61                        |h......a        |  system_identifier: 7040120944989169512
0x000|                        14 05 00 00            |        ....    |  pg_control_version: 1300
0x000|                                    2d e9 0b 0c|            -...|  catalog_version_no: 202107181
0x010|06 00 00 00                                    |....            |  state: "DB_IN_PRODUCTION" (6)
0x010|                        46 a1 c2 62 00 00 00 00|        F..b....|  time: "Mon, 04 Jul 2022 08:13:58 UTC" (1656922438)
0x020|c0 88 d6 9a 00 00 00 00                        |........        |  checkPoint: "0/9AD688C0" (2597750976)
     |                                               |                |  checkPointCopy{}:
0x020|                        c0 88 d6 9a 00 00 00 00|        ........|    redo: "0/9AD688C0" (2597750976)
0x030|01 00 00 00                                    |....            |    ThisTimeLineID: 1
0x030|            01 00 00 00                        |    ....        |    PrevTimeLineID: 1
0x030|                        01                     |        .       |    fullPageWrites: 1
0x040|4b ab 1c 00 00 00 00 00                        |K.......        |    nextXid: 1878859
0x040|                        dd 81 00 00            |        ....    |    nextOid: 33245
0x040|                                    01 00 00 00|            ....|    nextMulti: 1
0x050|00 00 00 00                                    |....            |    nextMultiOffset: 0
0x050|            d6 02 00 00                        |    ....        |    oldestXid: 726
0x050|                        01 00 00 00            |        ....    |    oldestXidDB: 1
0x050|                                    01 00 00 00|            ....|    oldestMulti: 1
0x060|01 00 00 00                                    |....            |    oldestMultiDB: 1
0x060|                        46 a1 c2 62 00 00 00 00|        F..b....|    time: "Mon, 04 Jul 2022 08:13:58 UTC" (1656922438)
0x070|00 00 00 00                                    |....            |    oldestCommitTsXid: 0
0x070|            00 00 00 00                        |    ....        |    newestCommitTsXid: 0
0x070|                        00 00 00 00            |        ....    |    oldestActiveXid: 0
0x080|e8 03 00 00 00 00 00 00                        |........        |  unloggedLSN: "0/3E8" (1000)
0x080|                        00 00 00 00 00 00 00 00|        ........|  minRecoveryPoint: "0/0" (0)
0x090|00 00 00 00                                    |....            |  minRecoveryPointTLI: 0
0x090|                        00 00 00 00 00 00 00 00|        ........|  backupStartPoint: "0/0" (0)
0x0a0|00 00 00 00 00 00 00 00                        |........        |  backupEndPoint: "0/0" (0)
0x0a0|                        00                     |        .       |  backupEndRequired: 0
0x0a0|                                    01 00 00 00|            ....|  wal_level: "WAL_LEVEL_REPLICA" (1)
0x0b0|00                                             |.               |  wal_log_hints: 0
0x0b0|            64 00 00 00                        |    d...        |  MaxConnections: 100
0x0b0|                        08 00 00 00            |        ....    |  max_worker_processes: 8
0x0b0|                                    0a 00 00 00|            ....|  max_wal_senders: 10
0x0c0|00 00 00 00                                    |....            |  max_prepared_xacts: 0
0x0c0|            40 00 00 00                        |    @...        |  max_locks_per_xact: 64
0x0c0|                        00                     |        .       |  track_commit_timestamp: 0
0x0c0|                                    08 00 00 00|            ....|  maxAlign: 8
0x0d0|00 00 00 00 87 d6 32 41                        |......2A        |  floatFormat: 1.234567e+06
0x0d0|                        00 20 00 00            |        . ..    |  blcksz: 8192
0x0d0|                                    00 00 02 00|            ....|  relseg_size: 131072
0x0e0|00 20 00 00                                    |. ..            |  xlog_blcksz: 8192
0x0e0|            00 00 00 01                        |    ....        |  xlog_seg_size: 16777216
0x0e0|                        40 00 00 00            |        @...    |  nameDataLen: 64
0x0e0|                                    20 00 00 00|             ...|  indexMaxKeys: 32
0x0f0|cc 07 00 00                                    |....            |  toast_max_chunk_size: 1996
0x0f0|            00 08 00 00                        |    ....        |  loblksize: 2048
0x0f0|                        01                     |        .       |  float8ByVal: 1
0x0f0|                                    00 00 00 00|            ....|  data_checksum_version: 0
0x100|00 45 fd 64 7e d4 d3 53 82 75 0a b7 d6 be c1 9a|.E.d~..S.u......|  mock_authentication_nonce: "0045fd647ed4d35382750ab7d6bec19a77af72bae00f728..." (raw bits)
0x110|77 af 72 ba e0 0f 72 80 4a 57 43 fb 76 c8 98 8c|w.r...r.JWC.v...|
0x120|4b 76 27 eb                                    |Kv'.            |  crc: 3945231947
```

Specific fields can be got by request:
```
$ fq -d pg_control -o flavour=postgres14 ".state, .checkPointCopy.redo, .wal_level" pg_control
    |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|
0x10|06 00 00 00                                    |....            |.state: "DB_IN_PRODUCTION" (6)
    |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|
0x20|                        c0 88 d6 9a 00 00 00 00|        ........|.checkPointCopy.redo: "0/9AD688C0" (2597750976)
    |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|
0xa0|                                    01 00 00 00|            ....|.wal_level: "WAL_LEVEL_REPLICA" (1)
```

### Heap data format

Heap page structure:
 - [https://www.postgresql.org/docs/current/storage-page-layout.html](https://www.postgresql.org/docs/current/storage-page-layout.html)
 - [https://postgrespro.ru/docs/postgresql/14/storage-page-layout?lang=en](https://postgrespro.ru/docs/postgresql/14/storage-page-layout?lang=en)

Heap consists of pages. You can see page content:
```
$ fq -d pg_heap -o flavour=postgres14 ".[0]" 16397
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0]{}: HeapPage
0x0000|00 00 00 00 f0 aa 72 01 00 00 04 00 0c 01 80 01 00 20 04 20|......r.......... . |  PageHeaderData{}:
*     |until 0x10b.7 (268)                                        |                    |
0x0104|                        00 00 00 00 00 00 00 00 00 00 00 00|        ............|  FreeSpace: "00000000000000000000000000000000000000000000000..." (raw bits)
0x0118|00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00|....................|
*     |until 0x17f.7 (116)                                        |                    |
0x017c|            e2 02 00 00 00 00 00 00 50 04 00 00 00 00 00 00|    ........P.......|  Tuples[0:61]:
0x0190|3d 00 04 00 02 09 18 00 3d 00 00 00 01 00 00 00 00 00 00 00|=.......=...........|
*     |until 0x1ff8.7 (7801)
```

To get PageHeaderData you can use this request:
```
$ fq -d pg_heap -o flavour=postgres14 ".[0].PageHeaderData" 16397
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].PageHeaderData{}:
0x000|00 00 00 00 f0 aa 72 01                                    |......r.            |  pd_lsn{}:
0x000|                        00 00                              |        ..          |  pd_checksum: 0
0x000|                              04 00                        |          ..        |  pd_flags: 4
0x000|                                    0c 01                  |            ..      |  pd_lower: 268
0x000|                                          80 01            |              ..    |  pd_upper: 384
0x000|                                                00 20      |                .   |  pd_special: 8192
0x000|                                                      04 20|                  . |  pd_pagesize_version: 8196
0x014|00 00 00 00                                                |....                |  pd_prune_xid: 0
0x014|            80 9f f2 00 00 9f f2 00 80 9e f2 00 00 9e f2 00|    ................|  pd_linp[0:61]:
0x028|80 9d f2 00 00 9d f2 00 80 9c f2 00 00 9c f2 00 80 9b f2 00|....................|
*    |until 0x10b.7 (244)
```

To get first and last item pointers on first page:
```
$ fq -d pg_heap -o flavour=postgres14 ".[0].PageHeaderData.pd_linp[0, -1]" 16397
    |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].PageHeaderData.pd_linp[0]{}: ItemIdData
0x14|            80 9f f2 00                                    |    ....            |  lp_off: 8064
0x14|            80 9f f2 00                                    |    ....            |  lp_flags: "LP_NORMAL" (1)
0x14|            80 9f f2 00                                    |    ....            |  lp_len: 121
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].PageHeaderData.pd_linp[60]{}: ItemIdData
0x104|            80 81 f2 00                                    |    ....            |  lp_off: 384
0x104|            80 81 f2 00                                    |    ....            |  lp_flags: "LP_NORMAL" (1)
0x104|            80 81 f2 00                                    |    ....            |  lp_len: 121
```

First and last tuple on first page:
```
$ fq -d pg_heap -o flavour=postgres14 ".[0].Tuples[0, -1]" 16397
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].Tuples[0]{}: HeapTupleHeaderData
0x1f7c|            e2 02 00 00 00 00 00 00 50 04 00 00            |    ........P...    |  t_choice{}:
0x1f7c|                                                00 00 00 00|                ....|  t_ctid{}:
0x1f90|01 00                                                      |..                  |
0x1f90|      04 00                                                |  ..                |  t_infomask2: 4
0x1f90|      04 00                                                |  ..                |  Infomask2{}:
0x1f90|            02 09                                          |    ..              |  t_infomask: 2306
0x1f90|            02 09                                          |    ..              |  Infomask{}:
0x1f90|                  18                                       |      .             |  t_hoff: 24
0x1f90|                        01 00 00 00 01 00 00 00 00 00 00 00|        ............|  t_bits: "010000000100000000000000ab202020202020202020202..." (raw bits)
0x1fa4|ab 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20|.                   |
*     |until 0x1ff8.7 (97)                                        |                    |
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].Tuples[60]{}: HeapTupleHeaderData
0x17c|            e2 02 00 00 00 00 00 00 50 04 00 00            |    ........P...    |  t_choice{}:
0x17c|                                                00 00 00 00|                ....|  t_ctid{}:
0x190|3d 00                                                      |=.                  |
0x190|      04 00                                                |  ..                |  t_infomask2: 4
0x190|      04 00                                                |  ..                |  Infomask2{}:
0x190|            02 09                                          |    ..              |  t_infomask: 2306
0x190|            02 09                                          |    ..              |  Infomask{}:
0x190|                  18                                       |      .             |  t_hoff: 24
0x190|                        3d 00 00 00 01 00 00 00 00 00 00 00|        =...........|  t_bits: "3d0000000100000000000000ab202020202020202020202..." (raw bits)
0x1a4|ab 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20|.                   |
*    |until 0x1f8.7 (97)
```

Tuple contains a lot of field:
```
$ fq -d pg_heap -o flavour=postgres14 ".[0].Tuples[0] | d" 16397
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].Tuples[0]{}: HeapTupleHeaderData
      |                                                           |                    |  t_choice{}:
      |                                                           |                    |    t_heap{}:
0x1f7c|            e2 02 00 00                                    |    ....            |      t_xmin: 738
0x1f7c|                        00 00 00 00                        |        ....        |      t_xmax: 0
      |                                                           |                    |      t_field3{}:
0x1f7c|                                    50 04 00 00            |            P...    |        t_cid: 1104
0x1f7c|                                    50 04 00 00            |            P...    |        t_xvac: 1104
      |                                                           |                    |    t_datum{}:
0x1f7c|            e2 02 00 00                                    |    ....            |      datum_len_: 738
0x1f7c|                        00 00 00 00                        |        ....        |      datum_typmod: 0
0x1f7c|                                    50 04 00 00            |            P...    |      datum_typeid: 1104
      |                                                           |                    |  t_ctid{}:
0x1f7c|                                                00 00 00 00|                ....|    ip_blkid: 0
0x1f90|01 00                                                      |..                  |    ip_posid: 1
0x1f90|      04 00                                                |  ..                |  t_infomask2: 4
      |                                                           |                    |  Infomask2{}:
0x1f90|      04 00                                                |  ..                |    HEAP_KEYS_UPDATED: 0
0x1f90|      04 00                                                |  ..                |    HEAP_HOT_UPDATED: 0
0x1f90|      04 00                                                |  ..                |    HEAP_ONLY_TUPLE: 0
0x1f90|            02 09                                          |    ..              |  t_infomask: 2306
      |                                                           |                    |  Infomask{}:
0x1f90|            02 09                                          |    ..              |    HEAP_HASNULL: 0
0x1f90|            02 09                                          |    ..              |    HEAP_HASVARWIDTH: 1
0x1f90|            02 09                                          |    ..              |    HEAP_HASEXTERNAL: 0
0x1f90|            02 09                                          |    ..              |    HEAP_HASOID_OLD: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_KEYSHR_LOCK: 0
0x1f90|            02 09                                          |    ..              |    HEAP_COMBOCID: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_EXCL_LOCK: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_LOCK_ONLY: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_SHR_LOCK: 0
0x1f90|            02 09                                          |    ..              |    HEAP_LOCK_MASK: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMIN_COMMITTED: 1
0x1f90|            02 09                                          |    ..              |    HEAP_XMIN_INVALID: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMIN_FROZEN: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_COMMITTED: 0
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_INVALID: 1
0x1f90|            02 09                                          |    ..              |    HEAP_XMAX_IS_MULTI: 0
0x1f90|            02 09                                          |    ..              |    HEAP_UPDATED: 0
0x1f90|            02 09                                          |    ..              |    HEAP_MOVED_OFF: 0
0x1f90|            02 09                                          |    ..              |    HEAP_MOVED_IN: 0
0x1f90|            02 09                                          |    ..              |    HEAP_MOVED: 0
0x1f90|                  18                                       |      .             |  t_hoff: 24
0x1f90|                        01 00 00 00 01 00 00 00 00 00 00 00|        ............|  t_bits: "010000000100000000000000ab202020202020202020202..." (raw bits)
0x1fa4|ab 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20|.                   |
*     |until 0x1ff8.7 (97)
```

Some of tuple fields are C unions:
```
$ time fq -d pg_heap -o flavour=postgres14 ".[0].Tuples[0].t_choice | d" 16397
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].Tuples[0].t_choice{}:
      |                                                           |                    |  t_heap{}:
0x1f7c|            e2 02 00 00                                    |    ....            |    t_xmin: 738
0x1f7c|                        00 00 00 00                        |        ....        |    t_xmax: 0
      |                                                           |                    |    t_field3{}:
0x1f7c|                                    50 04 00 00            |            P...    |      t_cid: 1104
0x1f7c|                                    50 04 00 00            |            P...    |      t_xvac: 1104
      |                                                           |                    |  t_datum{}:
0x1f7c|            e2 02 00 00                                    |    ....            |    datum_len_: 738
0x1f7c|                        00 00 00 00                        |        ....        |    datum_typmod: 0
0x1f7c|                                    50 04 00 00            |            P...    |    datum_typeid: 1104
```

Some tuple fields are flags:
```
$ time fq -d pg_heap -o flavour=postgres14 ".[0].Tuples[0] | .t_infomask, .Infomask" 16397
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|
0x1f90|            02 09                                          |    ..              |.[0].Tuples[0].t_infomask: 2306
      |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13|0123456789abcdef0123|.[0].Tuples[0].Infomask{}:
0x1f90|            02 09                                          |    ..              |  HEAP_HASNULL: 0
0x1f90|            02 09                                          |    ..              |  HEAP_HASVARWIDTH: 1
0x1f90|            02 09                                          |    ..              |  HEAP_HASEXTERNAL: 0
0x1f90|            02 09                                          |    ..              |  HEAP_HASOID_OLD: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_KEYSHR_LOCK: 0
0x1f90|            02 09                                          |    ..              |  HEAP_COMBOCID: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_EXCL_LOCK: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_LOCK_ONLY: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_SHR_LOCK: 0
0x1f90|            02 09                                          |    ..              |  HEAP_LOCK_MASK: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMIN_COMMITTED: 1
0x1f90|            02 09                                          |    ..              |  HEAP_XMIN_INVALID: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMIN_FROZEN: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_COMMITTED: 0
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_INVALID: 1
0x1f90|            02 09                                          |    ..              |  HEAP_XMAX_IS_MULTI: 0
0x1f90|            02 09                                          |    ..              |  HEAP_UPDATED: 0
0x1f90|            02 09                                          |    ..              |  HEAP_MOVED_OFF: 0
0x1f90|            02 09                                          |    ..              |  HEAP_MOVED_IN: 0
0x1f90|            02 09                                          |    ..              |  HEAP_MOVED: 0
```