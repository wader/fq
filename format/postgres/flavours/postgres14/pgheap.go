package postgres14

import (
	"context"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
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

	itemIds []itemIdDataD
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

		pagePosBegin := common.RoundDown(heap.pageSize, uint64(d.Pos()/8))

		d.FieldStruct("HeapPage", func(d *decode.D) {
			/*    0      |     8 */ // PageXLogRecPtr pd_lsn;
			/*    8      |     2 */ // uint16 pd_checksum;
			/*   10      |     2 */ // uint16 pd_flags;
			/*   12      |     2 */ // LocationIndex pd_lower;
			/*   14      |     2 */ // LocationIndex pd_upper;
			/*   16      |     2 */ // LocationIndex pd_special;
			/*   18      |     2 */ // uint16 pd_pagesize_version;
			/*   20      |     4 */ // TransactionId pd_prune_xid;
			/*   24      |     0 */ // ItemIdData pd_linp[];
			d.FieldStruct("PageHeaderData", func(d *decode.D) {
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
				itemsEnd := int64(pagePosBegin*8) + int64(page.pdLower*8)
				d.FieldArray("pd_linp", func(d *decode.D) {
					for {
						checkPos := d.Pos()
						if checkPos >= itemsEnd {
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
				}) // pd_linp in PageHeaderData

			}) // PageHeaderData, PageHeader

			// free space
			freeSpaceEnd := int64(pagePosBegin*8) + int64(page.pdUpper*8)
			freeSpaceNBits := freeSpaceEnd - d.Pos()
			d.FieldRawLen("FreeSpace", freeSpaceNBits, scalar.RawHex)

			// Tuples
			d.FieldArray("Tuples", func(d *decode.D) {
				for i := 0; i < len(page.itemIds); i++ {
					id := page.itemIds[i]
					if id.lpOff == 0 || id.lpLen == 0 {
						continue
					}

					pos := int64(pagePosBegin)*8 + int64(page.itemIds[i].lpOff)*8
					tupleDataLen := id.lpLen - SizeOfHeapTupleHeaderData

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
						d.FieldStruct("t_choice", func(d *decode.D) {
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

							d.SeekAbs(pos1)
							// type = struct DatumTupleFields {
							/*    0      |     4 */ // int32 datum_len_;
							/*    4      |     4 */ // int32 datum_typmod;
							/*    8      |     4 */ // Oid datum_typeid;
							//
							/* total size (bytes):   12 */

							// t_datum
							d.FieldStruct("t_datum", func(d *decode.D) {
								d.FieldS32("datum_len_")
								d.FieldS32("datum_typmod")
								d.FieldU32("datum_typeid")
							}) // DatumTupleFields t_datum
						}) // t_choice

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
						d.FieldU16("t_infomask")
						d.FieldU8("t_hoff")
						d.U8()

						d.FieldRawLen("t_bits", int64(tupleDataLen*8), scalar.RawHex)

					}) // HeapTupleHeaderData
				}
			}) // Tuples

			// end of page
			endLen := uint64(d.Pos() / 8)
			pageEnd := common.TypeAlign(heap.pageSize, endLen)
			d.SeekAbs(int64(pageEnd) * 8)
		}) // HeapPage

	} // for Heap pages
}
