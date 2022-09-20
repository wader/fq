package postgres14

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/flavours/postgres14/common14"
	"github.com/wader/fq/pkg/decode"
)

const (
	BTREE_MAGIC = 0x053162
	P_NONE      = 0

	/* Bits defined in btpo_flags */
	BTP_LEAF             = 1 << 0 /* leaf page, i.e. not internal page */
	BTP_ROOT             = 1 << 1 /* root page (has no parent) */
	BTP_DELETED          = 1 << 2 /* page has been deleted from tree */
	BTP_META             = 1 << 3 /* meta-page */
	BTP_HALF_DEAD        = 1 << 4 /* empty, but still in tree */
	BTP_SPLIT_END        = 1 << 5 /* rightmost page of split group */
	BTP_HAS_GARBAGE      = 1 << 6 /* page has LP_DEAD tuples (deprecated) */
	BTP_INCOMPLETE_SPLIT = 1 << 7 /* right sibling's downlink is missing */
	BTP_HAS_FULLXID      = 1 << 8 /* contains BTDeletedPageData */
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

// struct BTPageOpaqueData {
/*    0      |     4 */ // BlockNumber btpo_prev;
/*    4      |     4 */ // BlockNumber btpo_next;
/*    8      |     4 */ // uint32 btpo_level;
/*   12      |     2 */ // uint16 btpo_flags;
/*   14      |     2 */ // BTCycleId btpo_cycleid;
//
/* total size (bytes):   16 */

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
	page     *HeapPage
}

type HeapPage struct {
	// PageHeaderData fields
	PdLower           uint16
	PdUpper           uint16
	PdSpecial         uint16
	PdPagesizeVersion uint16

	// calculated bytes positions
	bytesPosBegin   int64 // bytes pos of page's beginning
	bytesPosEnd     int64 // bytes pos of page's ending
	bytesPosSpecial int64 // bytes pos of page's special

	// calculated bits positions
	posItemsEnd int64 // bits pos of items end

	// parsed items positions
	ItemIds []common14.ItemIdData
}

func decodeBTreePages(btree *BTreeD, d *decode.D) {
	for i := 0; ; i++ {
		if end, _ := d.TryEnd(); end {
			return
		}

		page := &HeapPage{}
		if btree.page != nil {
			// use prev page
			page.bytesPosBegin = btree.page.bytesPosEnd
		}
		page.bytesPosEnd = int64(common.TypeAlign(btree.PageSize, uint64(page.bytesPosBegin)+1))
		btree.page = page

		pos0 := page.bytesPosBegin * 8
		d.SeekAbs(pos0)

		if i == 0 {
			// first page contains meta information
			d.FieldStruct("page", func(d *decode.D) {
				decodeBTreeMetaPage(btree, d)
			})
			continue
		}

		if i > 2 {
			// limit for dev only
			break
		}

		d.FieldStruct("page", func(d *decode.D) {
			decodeBTreePage(btree, d)
		})
	}
}

func decodeBTreeMetaPage(btree *BTreeD, d *decode.D) {
	page := btree.page

	d.FieldStruct("page_header", func(d *decode.D) {
		decodePageHeader(btree, d)
	})
	d.FieldStruct("meta_page_data", func(d *decode.D) {
		decodeBTMetaPageData(btree, d)
	})

	pos0 := d.Pos()
	pos1 := int64(btree.page.bytesPosSpecial) * 8
	d.FieldRawLen("unused0", pos1-pos0)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(btree, d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.bytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on meta page\n")
	}
}

func decodePageHeader(btree *BTreeD, d *decode.D) {
	page := btree.page

	d.FieldStruct("pd_lsn", func(d *decode.D) {
		/*    0      |     4 */ // uint32 xlogid;
		/*    4      |     4 */ // uint32 xrecoff;
		d.FieldU32("xlogid", common.HexMapper)
		d.FieldU32("xrecoff", common.HexMapper)
	})
	d.FieldU16("pd_checksum")
	d.FieldU16("pd_flags")
	page.PdLower = uint16(d.FieldU16("pd_lower"))
	page.PdUpper = uint16(d.FieldU16("pd_upper"))
	page.PdSpecial = uint16(d.FieldU16("pd_special"))
	page.PdPagesizeVersion = uint16(d.FieldU16("pd_pagesize_version"))
	d.FieldU32("pd_prune_xid")

	page.bytesPosSpecial = page.bytesPosBegin + int64(page.PdSpecial)
	page.posItemsEnd = (page.bytesPosBegin * 8) + int64(page.PdLower*8)

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

// struct BTPageOpaqueData {
/*    0      |     4 */ // BlockNumber btpo_prev;
/*    4      |     4 */ // BlockNumber btpo_next;
/*    8      |     4 */ // uint32 btpo_level;
/*   12      |     2 */ // uint16 btpo_flags;
/*   14      |     2 */ // BTCycleId btpo_cycleid;
func decodeBTPageOpaqueData(btree *BTreeD, d *decode.D) {
	prev := d.FieldU32("btpo_prev")
	next := d.FieldU32("btpo_next")
	d.FieldU32("btpo_level")
	flags := d.FieldU16("btpo_flags")
	d.FieldU16("btpo_cycleid")

	isLeftMost := prev == P_NONE
	isRightMost := next == P_NONE
	isLeaf := (flags & BTP_LEAF) != 0
	isRoot := (flags & BTP_ROOT) != 0
	isDeleted := (flags & BTP_DELETED) != 0
	isMeta := (flags & BTP_META) != 0
	isHalfDead := (flags & BTP_HALF_DEAD) != 0
	isIgnore := isDeleted || isHalfDead
	hasGarbage := (flags & BTP_HAS_GARBAGE) != 0
	isIncompleteSplit := (flags & BTP_INCOMPLETE_SPLIT) != 0
	hasFullXid := (flags & BTP_HAS_FULLXID) != 0

	d.FieldStruct("flags", func(d *decode.D) {
		d.FieldValueBool("is_leftmost", isLeftMost)
		d.FieldValueBool("is_right_most", isRightMost)
		d.FieldValueBool("is_leaf", isLeaf)
		d.FieldValueBool("is_root", isRoot)
		d.FieldValueBool("is_deleted", isDeleted)
		d.FieldValueBool("is_meta", isMeta)
		d.FieldValueBool("is_half_dead", isHalfDead)
		d.FieldValueBool("is_ignore", isIgnore)
		d.FieldValueBool("has_garbage", hasGarbage)
		d.FieldValueBool("is_incompleteSplit", isIncompleteSplit)
		d.FieldValueBool("has_fullXid", hasFullXid)
	})
}

func decodeBTreePage(btree *BTreeD, d *decode.D) {
	page := btree.page

	d.FieldStruct("page_header", func(d *decode.D) {
		decodePageHeader(btree, d)
	})

	pos0 := d.Pos()
	pos1 := int64(btree.page.bytesPosSpecial) * 8
	d.SeekAbs(pos1)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(btree, d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.bytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on btree page\n")
	}

	d.SeekAbs(pos0)
	d.FieldArray("pd_linp", func(d *decode.D) {
		decodeItemIds(btree, d)
	})
}

func decodeItemIds(btree *BTreeD, d *decode.D) {
	page := btree.page

	for {
		checkPos := d.Pos()
		if checkPos >= page.posItemsEnd {
			break
		}
		/*    0: 0   |     4 */ // unsigned int lp_off: 15
		/*    1: 7   |     4 */ // unsigned int lp_flags: 2
		/*    2: 1   |     4 */ // unsigned int lp_len: 15
		d.FieldStruct("ItemIdData", func(d *decode.D) {
			itemID := common14.ItemIdData{}

			itemPos := d.Pos()
			itemID.Off = uint32(d.FieldU32("lp_off", common.LpOffMapper))
			d.SeekAbs(itemPos)
			itemID.Flags = uint32(d.FieldU32("lp_flags", common.LpFlagsMapper))
			d.SeekAbs(itemPos)
			itemID.Len = uint32(d.FieldU32("lp_len", common.LpLenMapper))

			page.ItemIds = append(page.ItemIds, itemID)
		})
	} // for pd_linp
}
