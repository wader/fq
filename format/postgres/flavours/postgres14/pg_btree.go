package postgres14

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
)

const (
	BTREE_MAGIC = 0x053162
)

// struct BTMetaPageData {
/*    0      |     4 */ // uint32 btm_magic
/*    4      |     4 */ // uint32 btm_version
/*    8      |     4 */ // BlockNumber btm_root
/*   12      |     4 */ // uint32 btm_level
/*   16      |     4 */ // BlockNumber btm_fastroot
/*   20      |     4 */ // uint32 btm_fastlevel
/*   24      |     4 */ // uint32 btm_last_cleanup_num_delpages
/* XXX  4-byte hole */
/*   32      |     8 */ // float8 btm_last_cleanup_num_heap_tuples
/*   40      |     1 */ // _Bool btm_allequalimage
/* XXX  7-byte padding */
//
/* total size (bytes):   48 */

func DecodePgBTree(d *decode.D) any {
	d.SeekAbs(0)

	btree := &BTreeD{
		PageSize: common.HeapPageSize,
	}
	decodeBTreePages(btree, d)

	return nil
}

type BTreeD struct {
	PageSize uint64
	page     *BTreePage
}

type BTreePage struct {
	heap            HeapPage
	bytesPosBegin   uint64 // bytes pos of page's beginning
	bytesPosEnd     uint64 // bytes pos of page's ending
	bytesPosSpecial uint64 // bytes pos of page's special
}

type HeapPage struct {
	PdLower           uint16
	PdUpper           uint16
	PdSpecial         uint16
	PdPagesizeVersion uint16
}

func decodeBTreePages(btree *BTreeD, d *decode.D) {
	for i := 0; ; i++ {
		if end, _ := d.TryEnd(); end {
			return
		}

		page := &BTreePage{}
		if btree.page != nil {
			// use prev page
			page.bytesPosBegin = btree.page.bytesPosEnd
		}
		page.bytesPosEnd = common.TypeAlign(btree.PageSize, page.bytesPosBegin+1)
		btree.page = page

		if i == 0 {
			// first page contains meta information
			d.FieldStruct("heap_page", func(d *decode.D) {
				decodeBTreeMetaPage(btree, d)
			})
			continue
		}

		if i > 0 {
			return
		}
	}
}

func decodeBTreeMetaPage(btree *BTreeD, d *decode.D) {
	d.FieldStruct("page_header", func(d *decode.D) {
		decodePageHeader(btree, d)
	})
	d.FieldStruct("meta_page_data", func(d *decode.D) {
		decodeBTMetaPageData(btree, d)
	})
}

func decodePageHeader(btree *BTreeD, d *decode.D) {
	heap := btree.page.heap

	d.FieldStruct("pd_lsn", func(d *decode.D) {
		/*    0      |     4 */ // uint32 xlogid;
		/*    4      |     4 */ // uint32 xrecoff;
		d.FieldU32("xlogid", common.HexMapper)
		d.FieldU32("xrecoff", common.HexMapper)
	})
	d.FieldU16("pd_checksum")
	d.FieldU16("pd_flags")
	heap.PdLower = uint16(d.FieldU16("pd_lower"))
	heap.PdUpper = uint16(d.FieldU16("pd_upper"))
	heap.PdSpecial = uint16(d.FieldU16("pd_special"))
	heap.PdPagesizeVersion = uint16(d.FieldU16("pd_pagesize_version"))
	d.FieldU32("pd_prune_xid")

	// ItemIdData pd_linp[];
	//page.ItemsEnd = int64(page.PagePosBegin*8) + int64(page.PdLower*8)
	//d.FieldArray("pd_linp", func(d *decode.D) {
	//	DecodeItemIds(heap, d)
	//})
}

func decodeBTMetaPageData(btree *BTreeD, d *decode.D) {
	/*    0      |     4 */ // uint32 btm_magic
	/*    4      |     4 */ // uint32 btm_version
	/*    8      |     4 */ // BlockNumber btm_root
	/*   12      |     4 */ // uint32 btm_level
	/*   16      |     4 */ // BlockNumber btm_fastroot
	/*   20      |     4 */ // uint32 btm_fastlevel
	/*   24      |     4 */ // uint32 btm_last_cleanup_num_delpages
	/* XXX  4-byte hole */
	/*   32      |     8 */ // float8 btm_last_cleanup_num_heap_tuples
	/*   40      |     1 */ // _Bool btm_allequalimage
	/* XXX  7-byte padding */

	btmMagic := d.FieldU32("btm_magic")
	d.FieldU32("btm_version")
	d.FieldU32("btm_root")
	d.FieldU32("btm_level")
	d.FieldU32("btm_fastroot")
	d.FieldU32("btm_fastlevel")
	d.FieldU32("btm_last_cleanup_num_delpages")
	d.FieldU32("hole0")
	d.FieldF64("btm_last_cleanup_num_heap_tuples")
	d.FieldU8("btm_allequalimage")
	d.FieldU56("padding0")

	if btmMagic != BTREE_MAGIC {
		d.Fatalf("invalid btmMagic = %X, must be %X\n", btmMagic, BTREE_MAGIC)
	}
}
