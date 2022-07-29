package postgres14

import (
	"context"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
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
	page := &heapPageD{}
	heap.page = page

	pagePosBegin := common.RoundDown(heap.pageSize, uint64(d.Pos()/8))

	for {
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
							itemPos := d.Pos()
							d.FieldU32("lp_off", common.LpOffMapper)
							d.SeekAbs(itemPos)
							d.FieldU32("lp_flags", common.LpFlagsMapper)
							d.SeekAbs(itemPos)
							d.FieldU32("lp_len", common.LpLenMapper)
						})
					} // for pd_linp
				}) // pd_linp in PageHeaderData

			}) // PageHeaderData, PageHeader

			// end of page
			endLen := uint64(d.Pos() / 8)
			pageEnd := common.TypeAlign(heap.pageSize, endLen)
			d.SeekAbs(int64(pageEnd) * 8)
		}) // HeapPage

	} // for Heap pages
}
