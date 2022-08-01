package postgres14

import (
	"context"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	HEAP_HASNULL          = 0x0001 /* has null attribute(s) */
	HEAP_HASVARWIDTH      = 0x0002 /* has variable-width attribute(s) */
	HEAP_HASEXTERNAL      = 0x0004 /* has external stored attribute(s) */
	HEAP_HASOID_OLD       = 0x0008 /* has an object-id field */
	HEAP_XMAX_KEYSHR_LOCK = 0x0010 /* xmax is a key-shared locker */
	HEAP_COMBOCID         = 0x0020 /* t_cid is a combo CID */
	HEAP_XMAX_EXCL_LOCK   = 0x0040 /* xmax is exclusive locker */
	HEAP_XMAX_LOCK_ONLY   = 0x0080 /* xmax, if valid, is only a locker */

	HEAP_XMAX_SHR_LOCK = HEAP_XMAX_EXCL_LOCK | HEAP_XMAX_KEYSHR_LOCK

	HEAP_LOCK_MASK = HEAP_XMAX_SHR_LOCK | HEAP_XMAX_EXCL_LOCK | HEAP_XMAX_KEYSHR_LOCK

	HEAP_XMIN_COMMITTED = 0x0100 /* t_xmin committed */
	HEAP_XMIN_INVALID   = 0x0200 /* t_xmin invalid/aborted */
	HEAP_XMIN_FROZEN    = HEAP_XMIN_COMMITTED | HEAP_XMIN_INVALID
	HEAP_XMAX_COMMITTED = 0x0400 /* t_xmax committed */
	HEAP_XMAX_INVALID   = 0x0800 /* t_xmax invalid/aborted */
	HEAP_XMAX_IS_MULTI  = 0x1000 /* t_xmax is a MultiXactId */
	HEAP_UPDATED        = 0x2000 /* this is UPDATEd version of row */
	HEAP_MOVED_OFF      = 0x4000 /* moved to another place by pre-9.0
	 * VACUUM FULL; kept for binary
	 * upgrade support */
	HEAP_MOVED_IN = 0x8000 /* moved from another place by pre-9.0
	 * VACUUM FULL; kept for binary
	 * upgrade support */
	HEAP_MOVED = HEAP_MOVED_OFF | HEAP_MOVED_IN
)

const (
	HEAP_KEYS_UPDATED = 0x2000 /* tuple was updated and key cols modified, or tuple deleted */
	HEAP_HOT_UPDATED  = 0x4000 /* tuple was HOT-updated */
	HEAP_ONLY_TUPLE   = 0x8000 /* this is heap-only tuple */
)

// type = struct PageHeaderData
/*    0      |     8 */ // PageXLogRecPtr pd_lsn;
/*    8      |     2 */ // uint16 pd_checksum;
/*   10      |     2 */ // uint16 pd_flags;
/*   12      |     2 */ // LocationIndex pd_lower;
/*   14      |     2 */ // LocationIndex pd_upper;
/*   16      |     2 */ // LocationIndex pd_special;
/*   18      |     2 */ // uint16 pd_pagesize_version;
/*   20      |     4 */ // TransactionId pd_prune_xid;
/*   24      |     0 */ // ItemIdData pd_linp[];
//
/* total size (bytes):   24 */

// type = struct PageXLogRecPtr {
/*    0      |     4 */ // uint32 xlogid;
/*    4      |     4 */ // uint32 xrecoff;

/* total size (bytes):    8 */

// type = struct ItemIdData {
/*    0: 0   |     4 */ // unsigned int lp_off: 15
/*    1: 7   |     4 */ // unsigned int lp_flags: 2
/*    2: 1   |     4 */ // unsigned int lp_len: 15

/* total size (bytes):    4 */

// typedef uint16 LocationIndex;
// #define SizeOfPageHeaderData (offsetof(PageHeaderData, pd_linp))

// type = struct HeapTupleHeaderData {
/*    0      |    12 */ // union {
/*                12 */ //     HeapTupleFields t_heap;
/*                12 */ //     DatumTupleFields t_datum;
//
/* total size (bytes):   12  */
/* } t_choice;             //
/*   12      |     6 */ // ItemPointerData t_ctid;
/*   18      |     2 */ // uint16 t_infomask2;
/*   20      |     2 */ // uint16 t_infomask;
/*   22      |     1 */ // uint8 t_hoff;
/*   23      |     0 */ // bits8 t_bits[];
/* XXX  1-byte padding  */
//
/* total size (bytes):   24 */
const SizeOfHeapTupleHeaderData = 24

// type = struct HeapTupleFields {
/*    0      |     4 */ // TransactionId t_xmin;
/*    4      |     4 */ // TransactionId t_xmax;
/*    8      |     4 */ // union {
/*                 4 */ //    CommandId t_cid;
/*                 4 */ //    TransactionId t_xvac;
//                         } t_field3;
/*                      total size (bytes):    4 */
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

type heapD struct {
	pageSize uint64

	// current page
	page *heapPageD
}

type heapPageD struct {
	pdLower             uint16
	pdUpper             uint16
	pdSpecial           uint16
	pd_pagesize_version uint16

	itemIds      []itemIdDataD
	pagePosBegin uint64
	itemsEnd     int64
}

func (hp *heapPageD) getItemId(offset uint32) (bool, itemIdDataD) {
	for _, id := range hp.itemIds {
		if id.lpOff == offset {
			return true, id
		}
	}
	return false, itemIdDataD{}
}

type itemIdDataD struct {
	/*    0: 0   |     4 */ // unsigned int lp_off: 15
	/*    1: 7   |     4 */ // unsigned int lp_flags: 2
	/*    2: 1   |     4 */ // unsigned int lp_len: 15
	lpOff                   uint32
	lpFlags                 uint32
	lpLen                   uint32
}

func getHeapD(d *decode.D) *heapD {
	val := d.Ctx.Value("heap")
	return val.(*heapD)
}

func DecodeHeap(d *decode.D) any {
	heap := &heapD{
		pageSize: common.HeapPageSize,
	}
	parentCtx := d.Ctx
	ctx := context.WithValue(parentCtx, "heap", heap)
	d.Ctx = ctx

	d.SeekAbs(0)
	d.FieldArray("Pages", decodeHeapPages)

	return nil
}

func decodeHeapPages(d *decode.D) {
	heap := getHeapD(d)

	for {
		page := &heapPageD{}
		heap.page = page

		d.FieldStruct("HeapPage", decodeHeapPage)

		// end of page
		endLen := uint64(d.Pos() / 8)
		pageEnd := common.TypeAlign(heap.pageSize, endLen)
		d.SeekAbs(int64(pageEnd) * 8)
	}
}

func decodeHeapPage(d *decode.D) {
	heap := getHeapD(d)

	page := &heapPageD{}
	heap.page = page

	pagePosBegin := common.RoundDown(heap.pageSize, uint64(d.Pos()/8))
	page.pagePosBegin = pagePosBegin

	// PageHeader
	d.FieldStruct("PageHeaderData", decodePageHeaderData)

	// free space
	freeSpaceEnd := int64(pagePosBegin*8) + int64(page.pdUpper*8)
	freeSpaceNBits := freeSpaceEnd - d.Pos()
	d.FieldRawLen("FreeSpace", freeSpaceNBits, scalar.RawHex)

	// Tuples
	d.FieldArray("Tuples", decodeTuples)
}

/*    0      |     8 */ // PageXLogRecPtr pd_lsn;
/*    8      |     2 */ // uint16 pd_checksum;
/*   10      |     2 */ // uint16 pd_flags;
/*   12      |     2 */ // LocationIndex pd_lower;
/*   14      |     2 */ // LocationIndex pd_upper;
/*   16      |     2 */ // LocationIndex pd_special;
/*   18      |     2 */ // uint16 pd_pagesize_version;
/*   20      |     4 */ // TransactionId pd_prune_xid;
/*   24      |     0 */ // ItemIdData pd_linp[];
func decodePageHeaderData(d *decode.D) {
	heap := getHeapD(d)
	page := heap.page

	d.FieldStruct("pd_lsn", func(d *decode.D) {
		/*    0      |     4 */ // uint32 xlogid;
		/*    4      |     4 */ // uint32 xrecoff;
		d.FieldU32("xlogid", common.HexMapper)
		d.FieldU32("xrecoff", common.HexMapper)
	})
	d.FieldU16("pd_checksum")
	d.FieldU16("pd_flags")
	page.pdLower = uint16(d.FieldU16("pd_lower"))
	page.pdUpper = uint16(d.FieldU16("pd_upper"))
	page.pdSpecial = uint16(d.FieldU16("pd_special"))
	page.pd_pagesize_version = uint16(d.FieldU16("pd_pagesize_version"))
	d.FieldU32("pd_prune_xid")

	// ItemIdData pd_linp[];
	page.itemsEnd = int64(page.pagePosBegin*8) + int64(page.pdLower*8)
	d.FieldArray("pd_linp", decodeItemIds)
}

func decodeItemIds(d *decode.D) {
	heap := getHeapD(d)
	page := heap.page

	for {
		checkPos := d.Pos()
		if checkPos >= page.itemsEnd {
			break
		}
		/*    0: 0   |     4 */ // unsigned int lp_off: 15
		/*    1: 7   |     4 */ // unsigned int lp_flags: 2
		/*    2: 1   |     4 */ // unsigned int lp_len: 15
		d.FieldStruct("ItemIdData", func(d *decode.D) {
			itemId := itemIdDataD{}

			itemPos := d.Pos()
			itemId.lpOff = uint32(d.FieldU32("lp_off", common.LpOffMapper))
			d.SeekAbs(itemPos)
			itemId.lpFlags = uint32(d.FieldU32("lp_flags", common.LpFlagsMapper))
			d.SeekAbs(itemPos)
			itemId.lpLen = uint32(d.FieldU32("lp_len", common.LpLenMapper))

			page.itemIds = append(page.itemIds, itemId)
		})
	} // for pd_linp
}

func decodeTuples(d *decode.D) {
	heap := getHeapD(d)
	page := heap.page

	for i := 0; i < len(page.itemIds); i++ {
		id := page.itemIds[i]
		if id.lpOff == 0 || id.lpLen == 0 {
			continue
		}

		pos := int64(page.pagePosBegin)*8 + int64(page.itemIds[i].lpOff)*8
		tupleDataLen := id.lpLen - SizeOfHeapTupleHeaderData

		// seek to tuple with ItemId offset
		d.SeekAbs(pos)

		/*    0      |    12 */ // union {
		/*                12 */ //     HeapTupleFields t_heap;
		/*                12 */ //     DatumTupleFields t_datum;
		//						} t_choice;
		/* total size (bytes):   12  */
		/*
			/*   12      |     6 */ // ItemPointerData t_ctid;
		/*   18      |     2 */ // uint16 t_infomask2;
		/*   20      |     2 */ // uint16 t_infomask;
		/*   22      |     1 */ // uint8 t_hoff;
		/*   23      |     0 */ // bits8 t_bits[];
		/* XXX  1-byte padding  */
		//
		/* total size (bytes):   24 */
		d.FieldStruct("HeapTupleHeaderData", func(d *decode.D) {
			d.FieldStruct("t_choice", decodeTChoice)

			d.FieldStruct("t_ctid", func(d *decode.D) {
				/*    0      |     4 */ // BlockIdData ip_blkid;
				/*    4      |     2 */ // OffsetNumber ip_posid;
				d.FieldU32("ip_blkid")
				d.FieldU16("ip_posid")
			}) // ItemPointerData t_ctid

			/*   18      |     2 */ // uint16 t_infomask2;
			/*   20      |     2 */ // uint16 t_infomask;
			/*   22      |     1 */ // uint8 t_hoff;
			/*   23      |     0 */ // bits8 t_bits[];
			/* XXX  1-byte padding  */
			d.FieldU16("t_infomask2")
			d.FieldStruct("Infomask2", decodeInfomask2)
			d.FieldU16("t_infomask")
			d.FieldStruct("Infomask", decodeInfomask)

			d.FieldU8("t_hoff")
			d.U8()

			d.FieldRawLen("t_bits", int64(tupleDataLen*8), scalar.RawHex)

		}) // HeapTupleHeaderData
	} // for ItemsIds
}

func decodeInfomask2(d *decode.D) {
	pos := d.Pos() - 16
	d.SeekAbs(pos)
	d.FieldU16("HEAP_KEYS_UPDATED", common.Mask{Mask: HEAP_KEYS_UPDATED})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_HOT_UPDATED", common.Mask{Mask: HEAP_HOT_UPDATED})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_ONLY_TUPLE", common.Mask{Mask: HEAP_ONLY_TUPLE})
}

func decodeInfomask(d *decode.D) {
	pos := d.Pos() - 16
	d.SeekAbs(pos)
	d.FieldU16("HEAP_HASNULL", common.Mask{Mask: HEAP_HASNULL})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_HASVARWIDTH", common.Mask{Mask: HEAP_HASVARWIDTH})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_HASEXTERNAL", common.Mask{Mask: HEAP_HASEXTERNAL})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_HASOID_OLD", common.Mask{Mask: HEAP_HASOID_OLD})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_KEYSHR_LOCK", common.Mask{Mask: HEAP_XMAX_KEYSHR_LOCK})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_COMBOCID", common.Mask{Mask: HEAP_COMBOCID})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_EXCL_LOCK", common.Mask{Mask: HEAP_XMAX_EXCL_LOCK})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_LOCK_ONLY", common.Mask{Mask: HEAP_XMAX_LOCK_ONLY})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_SHR_LOCK", common.Mask{Mask: HEAP_XMAX_SHR_LOCK})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_LOCK_MASK", common.Mask{Mask: HEAP_LOCK_MASK})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMIN_COMMITTED", common.Mask{Mask: HEAP_XMIN_COMMITTED})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMIN_INVALID", common.Mask{Mask: HEAP_XMIN_INVALID})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMIN_FROZEN", common.Mask{Mask: HEAP_XMIN_FROZEN})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_COMMITTED", common.Mask{Mask: HEAP_XMAX_COMMITTED})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_INVALID", common.Mask{Mask: HEAP_XMAX_INVALID})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_XMAX_IS_MULTI", common.Mask{Mask: HEAP_XMAX_IS_MULTI})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_UPDATED", common.Mask{Mask: HEAP_UPDATED})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_MOVED_OFF", common.Mask{Mask: HEAP_MOVED_OFF})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_MOVED_IN", common.Mask{Mask: HEAP_MOVED_IN})
	d.SeekAbs(pos)
	d.FieldU16("HEAP_MOVED", common.Mask{Mask: HEAP_MOVED})
}

/*    0      |    12 */ // union {
/*                12 */ //     HeapTupleFields t_heap;
/*                12 */ //     DatumTupleFields t_datum;
//						} t_choice;
func decodeTChoice(d *decode.D) {
	pos1 := d.Pos()
	// type = struct HeapTupleFields {
	/*    0      |     4 */ // TransactionId t_xmin;
	/*    4      |     4 */ // TransactionId t_xmax;
	/*    8      |     4 */ // union {
	/*                 4 */ //    CommandId t_cid;
	/*                 4 */ //    TransactionId t_xvac;
	//                         } t_field3;
	/*                      total size (bytes):    4 */
	//
	/* total size (bytes):   12 */
	d.FieldStruct("t_heap", func(d *decode.D) {
		d.FieldU32("t_xmin")
		d.FieldU32("t_xmax")
		d.FieldStruct("t_field3", func(d *decode.D) {
			pos2 := d.Pos()
			d.FieldU32("t_cid")
			d.SeekAbs(pos2)
			d.FieldU32("t_xvac")
		}) // t_field3
	}) // HeapTupleFields t_heap

	// restore position for union
	d.SeekAbs(pos1)
	// type = struct DatumTupleFields {
	/*    0      |     4 */ // int32 datum_len_;
	/*    4      |     4 */ // int32 datum_typmod;
	/*    8      |     4 */ // Oid datum_typeid;
	//
	/* total size (bytes):   12 */
	d.FieldStruct("t_datum", func(d *decode.D) {
		d.FieldS32("datum_len_")
		d.FieldS32("datum_typmod")
		d.FieldU32("datum_typeid")
	}) // DatumTupleFields t_datum
}
