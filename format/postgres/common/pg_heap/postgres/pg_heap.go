package postgres

import (
	"fmt"

	"github.com/wader/fq/format"
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
/*   18      |     2 */ // uint16 PdPagesizeVersion;
/*   20      |     4 */ // TransactionId pd_prune_xid;
/*   24      |     0 */ // ItemIdData pd_linp[];
//
/* total size (bytes):   24 */

// type = struct PageXLogRecPtr {
/*    0      |     4 */ // uint32 xlogid;
/*    4      |     4 */ // uint32 xrecoff;

/* total size (bytes):    8 */

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

type Heap struct {
	Args format.Pg_Heap_In

	// current Page
	Page *HeapPage
	// Page special data
	Special *PageSpecial

	// current tuple
	Tuple *TupleD

	DecodePageHeaderData func(page *HeapPage, d *decode.D)
	DecodePageSpecial    func(heap *Heap, d *decode.D)
}

type PageSpecial struct {
	// pgproee
	PdXidBase   uint64 // 8 TransactionId pd_xid_base;
	PdMultiBase uint64 // 8 TransactionId pd_multi_base;
	PdPruneXid  uint64 // 4 ShortTransactionId pd_prune_xid;
	PdMagic     uint64 // 4 uint32 pd_magic;
}

type TupleD struct {
	IsMulti bool
}

func Decode(heap *Heap, d *decode.D) any {
	decodeHeapPages(heap, d)
	return nil
}

func decodeHeapPages(heap *Heap, d *decode.D) {
	blockNumber := uint32(heap.Args.Page + heap.Args.Segment*common.RelSegSize)
	count := int64(0)
	for {
		if d.End() {
			return
		}

		d.FieldStruct("page", func(d *decode.D) {
			decodeHeapPage(heap, d, blockNumber)
		})
		blockNumber++
		count++

		// end of Page
		endLen := uint64(d.Pos() / 8)
		pageEnd := int64(common.TypeAlign(common.PageSize, endLen))
		pageEnd0 := count * common.PageSize
		if pageEnd0 != pageEnd {
			d.Errorf("invalid page %d end expected %d, actual %d, endLen  %d\n", count-1, pageEnd0, pageEnd, endLen)
		}
		d.SeekAbs(pageEnd0 * 8)
	}
}

func decodeHeapPage(heap *Heap, d *decode.D, blockNumber uint32) {
	page := &HeapPage{}
	if heap.Page != nil {
		// use prev page
		page.BytesPosBegin = heap.Page.BytesPosEnd
	}
	page.BytesPosEnd = int64(common.TypeAlign(common.PageSize, uint64(page.BytesPosBegin)+1))
	heap.Page = page
	heap.Special = &PageSpecial{}

	checkSum := calcCheckSum(d, blockNumber)

	d.FieldStruct("page_header", func(d *decode.D) {
		heap.DecodePageHeaderData(page, d)

		d.FieldValueUint("pd_checksum_check", uint64(checkSum))
		sumEqual := page.PdChecksum == checkSum
		d.FieldValueBool("pd_checksum_check_equal", sumEqual)
	})

	DecodeItemIds(page, d)

	if uint64(page.PdSpecial) != common.PageSize && heap.DecodePageSpecial != nil {
		heap.DecodePageSpecial(heap, d)
	}

	// Tuples
	d.FieldArray("tuples", func(d *decode.D) {
		decodeTuples(heap, d)
	})
}

func calcCheckSum(d *decode.D, blockNumber uint32) uint16 {
	pos0 := d.Pos()
	pageBuffer := make([]byte, common.PageSize)
	rdrPage := d.RawLen(int64(common.PageSize * 8))
	_, err := rdrPage.ReadBits(pageBuffer, int64(common.PageSize*8))
	if err != nil {
		d.Fatalf("can't read page, err = %v", err)
	}
	sum2 := common.CheckSumBlock(pageBuffer, blockNumber)
	d.SeekAbs(pos0)
	return sum2
}

func decodeTuples(heap *Heap, d *decode.D) {
	page := heap.Page
	for i := 0; i < len(page.ItemIds); i++ {
		id := page.ItemIds[i]
		if id.Off == 0 || id.Len == 0 {
			continue
		}
		if id.Flags != common.LP_NORMAL {
			continue
		}

		pos := (page.BytesPosBegin * 8) + int64(id.Off)*8
		if id.Len < SizeOfHeapTupleHeaderData {
			d.Fatalf("item len = %d, is less than %d HeapTupleHeaderData", id.Len, SizeOfHeapTupleHeaderData)
		}
		tupleDataLen := id.Len - SizeOfHeapTupleHeaderData

		// seek to tuple with ItemID offset
		d.SeekAbs(pos)

		// type = struct HeapTupleHeaderData {
		/*    0      |    12 */ // union {
		/*                12 */ //     HeapTupleFields t_heap;
		/*                12 */ //     DatumTupleFields t_datum;
		//						} t_choice;
		/* total size (bytes):   12  */
		//
		/*   12      |     6 */ // ItemPointerData t_ctid;
		/*   18      |     2 */ // uint16 t_infomask2;
		/*   20      |     2 */ // uint16 t_infomask;
		/*   22      |     1 */ // uint8 t_hoff;
		/*   23      |     0 */ // bits8 t_bits[];
		/* XXX  1-byte padding  */
		//
		/* total size (bytes):   24 */
		d.FieldStruct("tuple", func(d *decode.D) {
			heap.Tuple = &TupleD{}

			d.FieldStruct("header", func(d *decode.D) {

				pos1 := d.Pos()
				// we need infomask before t_xmin, t_xmax
				d.SeekAbs(pos1 + 18*8)
				infomask2 := d.FieldU16("t_infomask2")
				d.FieldStruct("infomask2", func(d *decode.D) {
					decodeInfomask2(d, infomask2)
				})
				infomask := d.FieldU16("t_infomask")
				d.FieldStruct("infomask", func(d *decode.D) {
					decodeInfomask(heap, d, infomask)
				})

				// restore pos and continue
				d.SeekAbs(pos1)
				d.FieldStruct("t_choice", func(d *decode.D) {
					decodeTChoice(heap, d)
				})
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
				//d.FieldU16("t_infomask2")
				//d.FieldStruct("Infomask2", decodeInfomask2)
				//d.FieldU16("t_infomask")
				//d.FieldStruct("Infomask", decodeInfomask)
				// already done
				d.SeekRel(32)

				d.FieldU8("t_hoff")
				d.FieldU8("padding0")
			}) // HeapTupleHeaderData

			d.FieldRawLen("data", int64(tupleDataLen*8), scalar.RawHex)

			// data alignment
			pos2 := uint64(d.Pos() / 8)
			pos1Aligned := common.TypeAlign8(pos2)
			if pos2 != pos1Aligned {
				alignedLen := (pos1Aligned - pos2) * 8
				d.FieldRawLen("padding1", int64(alignedLen), scalar.RawHex)
			}
			pos3 := uint64(d.Pos() / 8)
			pos2Aligned := common.TypeAlign8(pos3)
			if pos3 != pos2Aligned {
				d.Fatalf("pos3 isn't aligned, pos2 = %d, pos3 = %d", pos2, pos3)
			}

		})

	} // for ItemsIds
}

func decodeInfomask2(d *decode.D, infomask2 uint64) {
	d.FieldValueBool("heap_keys_updated", common.IsMaskSet0(infomask2, HEAP_KEYS_UPDATED))
	d.FieldValueBool("heap_hot_updated", common.IsMaskSet0(infomask2, HEAP_HOT_UPDATED))
	d.FieldValueBool("heap_only_tuple", common.IsMaskSet0(infomask2, HEAP_ONLY_TUPLE))
}

func decodeInfomask(heap *Heap, d *decode.D, infomask uint64) {
	tuple := heap.Tuple

	isMulti := common.IsMaskSet0(infomask, HEAP_XMAX_IS_MULTI)
	tuple.IsMulti = isMulti

	d.FieldValueBool("heap_hasnull", common.IsMaskSet0(infomask, HEAP_HASNULL))
	d.FieldValueBool("heap_hasvarwidth", common.IsMaskSet0(infomask, HEAP_HASVARWIDTH))
	d.FieldValueBool("heap_hasexternal", common.IsMaskSet0(infomask, HEAP_HASEXTERNAL))
	d.FieldValueBool("heap_hasoid_old", common.IsMaskSet0(infomask, HEAP_HASOID_OLD))
	d.FieldValueBool("heap_xmax_keyshr_lock", common.IsMaskSet0(infomask, HEAP_XMAX_KEYSHR_LOCK))
	d.FieldValueBool("heap_combocid", common.IsMaskSet0(infomask, HEAP_COMBOCID))
	d.FieldValueBool("heap_xmax_excl_lock", common.IsMaskSet0(infomask, HEAP_XMAX_EXCL_LOCK))
	d.FieldValueBool("heap_xmax_lock_only", common.IsMaskSet0(infomask, HEAP_XMAX_LOCK_ONLY))
	d.FieldValueBool("heap_xmax_shr_lock", common.IsMaskSet0(infomask, HEAP_XMAX_SHR_LOCK))
	d.FieldValueBool("heap_lock_mask", common.IsMaskSet0(infomask, HEAP_LOCK_MASK))
	d.FieldValueBool("heap_xmin_committed", common.IsMaskSet0(infomask, HEAP_XMIN_COMMITTED))
	d.FieldValueBool("heap_xmin_invalid", common.IsMaskSet0(infomask, HEAP_XMIN_INVALID))
	d.FieldValueBool("heap_xmin_frozen", common.IsMaskSet0(infomask, HEAP_XMIN_FROZEN))
	d.FieldValueBool("heap_xmax_committed", common.IsMaskSet0(infomask, HEAP_XMAX_COMMITTED))
	d.FieldValueBool("heap_xmax_invalid", common.IsMaskSet0(infomask, HEAP_XMAX_INVALID))
	d.FieldValueBool("heap_xmax_is_multi", isMulti)
	d.FieldValueBool("heap_updated", common.IsMaskSet0(infomask, HEAP_UPDATED))
	d.FieldValueBool("heap_moved_off", common.IsMaskSet0(infomask, HEAP_MOVED_OFF))
	d.FieldValueBool("heap_moved_in", common.IsMaskSet0(infomask, HEAP_MOVED_IN))
	d.FieldValueBool("heap_moved", common.IsMaskSet0(infomask, HEAP_MOVED))
}

/*    0      |    12 */ // union {
/*                12 */ //     HeapTupleFields t_heap;
/*                12 */ //     DatumTupleFields t_datum;
//						} t_choice;
func decodeTChoice(heap *Heap, d *decode.D) {
	special := heap.Special
	tuple := heap.Tuple

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
		d.FieldU32("t_xmin", TransactionMapper{Heap: heap, Special: special, Tuple: tuple})
		d.FieldU32("t_xmax", TransactionMapper{Heap: heap, Special: special, Tuple: tuple})
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

type TransactionMapper struct {
	Heap    *Heap
	Special *PageSpecial
	Tuple   *TupleD
}

func (m TransactionMapper) MapUint(s scalar.Uint) (scalar.Uint, error) {
	xid := s.Actual

	if m.Special.PdXidBase != 0 && m.Tuple.IsMulti && common.TransactionIDIsNormal(xid) {
		xid64 := xid + m.Special.PdXidBase
		s.Sym = fmt.Sprintf("%d", xid64)
	}

	if m.Special.PdMultiBase != 0 && !m.Tuple.IsMulti && common.TransactionIDIsNormal(xid) {
		xid64 := xid + m.Special.PdMultiBase
		s.Sym = fmt.Sprintf("%d", xid64)
	}

	return s, nil
}
