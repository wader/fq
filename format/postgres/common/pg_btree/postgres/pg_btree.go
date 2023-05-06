package postgres

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/common/pg_heap/postgres"
	"github.com/wader/fq/pkg/decode"
)

const (
	BTREE_MAGIC = 0x053162
)

const (
	INDEX_SIZE_MASK = 0x1FFF
	INDEX_VAR_MASK  = 0x4000
	INDEX_NULL_MASK = 0x8000
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

func DecodePgBTree(d *decode.D, pageNr int) {
	var prevPage *postgres.HeapPage

	for i := pageNr; ; i++ {
		page := &postgres.HeapPage{}
		if prevPage != nil {
			// use prev page
			page.BytesPosBegin = prevPage.BytesPosEnd
		}
		page.BytesPosEnd = int64(common.TypeAlign(common.PageSize, uint64(page.BytesPosBegin)+1))
		prevPage = page

		pos0 := page.BytesPosBegin * 8
		d.SeekAbs(pos0)

		if d.End() {
			return
		}

		if i == 0 {
			// first page contains meta information
			d.FieldStruct("page", func(d *decode.D) {
				decodeBTreeMetaPage(page, d)
			})
			continue
		}

		d.FieldStruct("page", func(d *decode.D) {
			decodeBTreePage(page, d)
		})
	}
}

func decodeBTreeMetaPage(page *postgres.HeapPage, d *decode.D) {

	d.FieldStruct("page_header", func(d *decode.D) {
		postgres.DecodePageHeader(page, d)
	})
	d.FieldStruct("meta_page_data", func(d *decode.D) {
		decodeBTMetaPageData(d)
	})

	pos0 := d.Pos()
	pos1 := page.BytesPosSpecial * 8
	d.FieldRawLen("unused0", pos1-pos0)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.BytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on meta page")
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
		d.Fatalf("invalid btmMagic = %X, must be %X", btmMagic, BTREE_MAGIC)
	}
}

// struct BTPageOpaqueData {
/*    0      |     4 */ // BlockNumber btpo_prev;
/*    4      |     4 */ // BlockNumber btpo_next;
/*    8      |     4 */ // uint32 btpo_level;
/*   12      |     2 */ // uint16 btpo_flags;
/*   14      |     2 */ // BTCycleId btpo_cycleid;
func decodeBTPageOpaqueData(d *decode.D) {
	d.FieldU32("btpo_prev")
	d.FieldU32("btpo_next")
	d.FieldU32("btpo_level")

	// bits in uint16 LE: 7 - 0 15 - 8
	d.FieldStruct("btpo_flags", func(d *decode.D) {
		d.FieldBool("is_incomplete_split")
		d.FieldBool("has_garbage")
		d.FieldBool("split_end")
		isHalfDead := d.FieldBool("is_half_dead")
		d.FieldBool("is_meta")
		isDeleted := d.FieldBool("is_deleted")
		d.FieldBool("is_root")
		d.FieldBool("is_leaf")

		d.FieldU7("skip1")
		d.FieldBool("has_full_xid")

		d.FieldValueBool("is_ignore", isDeleted || isHalfDead)
	})

	d.FieldU16("btpo_cycleid")
}

func decodeBTreePage(page *postgres.HeapPage, d *decode.D) {
	d.FieldStruct("page_header", func(d *decode.D) {
		postgres.DecodePageHeader(page, d)
	})

	pos0 := d.Pos()
	pos1 := page.BytesPosSpecial * 8
	d.SeekAbs(pos1)
	d.FieldStruct("page_opaque_data", func(d *decode.D) {
		decodeBTPageOpaqueData(d)
	})
	pos2 := d.Pos()
	bytesPos2 := pos2 / 8
	if bytesPos2 != page.BytesPosEnd {
		d.Fatalf("invalid pos after read page_opaque_data on btree page")
	}

	d.SeekAbs(pos0)
	postgres.DecodeItemIds(page, d)

	d.FieldArray("tuples", func(d *decode.D) {
		decodeIndexTuples(page, d)
	})
}

func decodeIndexTuples(page *postgres.HeapPage, d *decode.D) {

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
				d.FieldValueUint("size", size)
				if size < IndexTupleDataSize {
					d.Fatalf("invalid size of tuple = %d", size)
				}
				dataSize := size - IndexTupleDataSize
				d.FieldRawLen("data", int64(dataSize*8))
			})

		})

	}
}
