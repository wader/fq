package postgres14

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/common/pg_heap/postgres"
	"github.com/wader/fq/pkg/decode"
)

//nolint:revive
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

//nolint:revive
const (
	INDEX_SIZE_MASK       = 0x1FFF
	INDEX_AM_RESERVED_BIT = 0x2000 /* reserved for index-AM specific usage */
	INDEX_VAR_MASK        = 0x4000
	INDEX_NULL_MASK       = 0x8000
)

const (
	IndexTupleDataSize = 8
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

// struct IndexTupleData {
/*    0      |     6 */ // ItemPointerData t_tid;
/*    6      |     2 */ // unsigned short t_info;
//
// IndexTupleData *IndexTuple;
/* total size (bytes):    8 */

func DecodePgBTree(d *decode.D) any {
	d.SeekAbs(0)

	btree := &BTree{
		PageSize: common.PageSize,
	}
	decodeBTreePages(btree, d)

	return nil
}

type BTree struct {
	PageSize uint64
	page     *postgres.HeapPage
}

func decodeBTreePages(btree *BTree, d *decode.D) {
	for i := 0; ; i++ {

		page := &postgres.HeapPage{}
		if btree.page != nil {
			// use prev page
			page.BytesPosBegin = btree.page.BytesPosEnd
		}
		page.BytesPosEnd = int64(common.TypeAlign(btree.PageSize, uint64(page.BytesPosBegin)+1))
		btree.page = page

		pos0 := page.BytesPosBegin * 8
		d.SeekAbs(pos0)

		if end, _ := d.TryEnd(); end {
			return
		}

		if i == 0 {
			// first page contains meta information
			d.FieldStruct("page", func(d *decode.D) {
				decodeBTreeMetaPage(btree, d)
			})
			continue
		}

		d.FieldStruct("page", func(d *decode.D) {
			decodeBTreePage(btree, d)
		})
	}
}

func decodeBTreeMetaPage(btree *BTree, d *decode.D) {
	page := btree.page

	d.FieldStruct("page_header", func(d *decode.D) {
		postgres.DecodePageHeader(page, d)
	})
	d.FieldStruct("meta_page_data", func(d *decode.D) {
		decodeBTMetaPageData(d)
	})

	pos0 := d.Pos()
	pos1 := btree.page.BytesPosSpecial * 8
	d.FieldRawLen("unused0", pos1-pos0)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.BytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on meta page\n")
	}
}

func decodeBTMetaPageData(d *decode.D) {
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
func decodeBTPageOpaqueData(d *decode.D) {
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
		d.FieldValueBool("is_rightmost", isRightMost)
		d.FieldValueBool("is_leaf", isLeaf)
		d.FieldValueBool("is_root", isRoot)
		d.FieldValueBool("is_deleted", isDeleted)
		d.FieldValueBool("is_meta", isMeta)
		d.FieldValueBool("is_half_dead", isHalfDead)
		d.FieldValueBool("is_ignore", isIgnore)
		d.FieldValueBool("has_garbage", hasGarbage)
		d.FieldValueBool("is_incomplete_split", isIncompleteSplit)
		d.FieldValueBool("has_full_xid", hasFullXid)
	})
}

func decodeBTreePage(btree *BTree, d *decode.D) {
	page := btree.page

	d.FieldStruct("page_header", func(d *decode.D) {
		postgres.DecodePageHeader(page, d)
	})

	pos0 := d.Pos()
	pos1 := btree.page.BytesPosSpecial * 8
	d.SeekAbs(pos1)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.BytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on btree page\n")
	}

	d.SeekAbs(pos0)
	postgres.DecodeItemIds(page, d)

	d.FieldArray("tuples", func(d *decode.D) {
		decodeIndexTuples(btree, d)
	})
}

func decodeIndexTuples(btree *BTree, d *decode.D) {
	page := btree.page

	for i := 0; i < len(page.ItemIds); i++ {
		id := page.ItemIds[i]
		if id.Off == 0 || id.Len == 0 {
			continue
		}
		if id.Flags != common.LP_NORMAL {
			continue
		}

		pos := (page.BytesPosBegin * 8) + int64(id.Off)*8

		// seek to tuple with ItemID offset
		d.SeekAbs(pos)
		d.FieldStruct("tuple", func(d *decode.D) {

			// IndexTupleData
			d.FieldStruct("index_tuple_data", func(d *decode.D) {
				// struct IndexTupleData {
				/*    0      |     6 */ // ItemPointerData t_tid;
				/*    6      |     2 */ // unsigned short t_info;
				//
				d.FieldStruct("t_tid", func(d *decode.D) {
					/*    0      |     4 */ // BlockIdData ip_blkid;
					/*    4      |     2 */ // OffsetNumber ip_posid;
					d.FieldU32("ip_blkid")
					d.FieldU16("ip_posid")
				})
				tInfo := d.FieldU16("t_info")

				size := tInfo & INDEX_SIZE_MASK
				hasNulls := (tInfo & INDEX_NULL_MASK) != 0
				hasVarWidths := (tInfo & INDEX_VAR_MASK) != 0
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldValueBool("has_nulls", hasNulls)
					d.FieldValueBool("has_var_widths", hasVarWidths)
				})
				d.FieldValueU("size", size)
				if size < IndexTupleDataSize {
					d.Fatalf("invalid size of tuple = %d\n", size)
				}
				dataSize := size - IndexTupleDataSize
				d.FieldRawLen("data", int64(dataSize*8))
			})

		})

	}
}
