package pgproee11

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/pgproee14/ee14"
	"github.com/wader/fq/format/postgres/flavours/postgres14/common14"
	"github.com/wader/fq/pkg/decode"
)

// type = struct PageHeaderData {
/*    0      |     8 */ // PageXLogRecPtr pd_lsn;
/*    8      |     2 */ // uint16 pd_checksum;
/*   10      |     2 */ // uint16 pd_flags;
/*   12      |     2 */ // LocationIndex pd_lower;
/*   14      |     2 */ // LocationIndex pd_upper;
/*   16      |     2 */ // LocationIndex pd_special;
/*   18      |     2 */ // uint16 pd_pagesize_version;
/*   20      |     0 */ // ItemIdData pd_linp[];
//
/* total size (bytes):   20 */

// type = struct PageXLogRecPtr {
/*    0      |     4 */ // uint32 xlogid;
/*    4      |     4 */ // uint32 xrecoff;
//
/* total size (bytes):    8 */

// type = struct HeapTupleHeaderData {
/*    0      |    12 */ // union {
/*                12 */ //     HeapTupleFields t_heap;
/*                12 */ //     DatumTupleFields t_datum;
//							} t_choice;
/* total size (bytes):   12 */
//
/*   12      |     6 */ // ItemPointerData t_ctid;
/*   18      |     2 */ // uint16 t_infomask2;
/*   20      |     2 */ // uint16 t_infomask;
/*   22      |     1 */ // uint8 t_hoff;
/*   23      |     0 */ // bits8 t_bits[];
/* XXX  1-byte padding */
//
/* total size (bytes):   24 */

// type = struct HeapTupleFields {
/*    0      |     4 */ // ShortTransactionId t_xmin;
/*    4      |     4 */ // ShortTransactionId t_xmax;
/*    8      |     4 */ // union {
/*                 4 */ //     CommandId t_cid;
/*                 4 */ //     ShortTransactionId t_xvac;
//							} t_field3;
/* total size (bytes):    4 */
//
/* total size (bytes):   12 */

// type = struct DatumTupleFields {
/*    0      |     4 */ // int32 datum_len_;
/*    4      |     4 */ // int32 datum_typmod;
/*    8      |     4 */ // Oid datum_typeid;
//
/* total size (bytes):   12 */

// type = struct ItemPointerData {
/*    0      |     4 */ // BlockIdData ip_blkid;
/*    4      |     2 */ // OffsetNumber ip_posid;
//
/* total size (bytes):    6 */

// type = struct HeapPageSpecialData {
/*    0      |     8 */ // TransactionId pd_xid_base;
/*    8      |     8 */ // TransactionId pd_multi_base;
/*   16      |     4 */ // ShortTransactionId pd_prune_xid;
/*   20      |     4 */ // uint32 pd_magic;
//
/* total size (bytes):   24 */

func DecodeHeap(d *decode.D, args format.PostgresHeapIn) any {
	heap := &common14.Heap{
		Args:                 args,
		DecodePageHeaderData: ee14.DecodePageHeaderData,
		DecodePageSpecial:    ee14.DecodePageSpecial,
	}
	return common14.DecodeHeap(heap, d)
}
