package common14

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// HeapPage used in tables, indexes...

// type = struct ItemIdData {
/*    0: 0   |     4 */ // unsigned int lp_off: 15
/*    1: 7   |     4 */ // unsigned int lp_flags: 2
/*    2: 1   |     4 */ // unsigned int lp_len: 15
//
/* total size (bytes):    4 */

type ItemId struct {
	Off   uint32 // unsigned int lp_off: 15
	Flags uint32 // unsigned int lp_flags: 2
	Len   uint32 // unsigned int lp_len: 15
}

type HeapPage struct {
	// PageHeaderData fields
	PdLower           uint16
	PdUpper           uint16
	PdSpecial         uint16
	PdPageSizeVersion uint16

	// calculated bytes positions
	BytesPosBegin   int64 // bytes pos of page's beginning
	BytesPosEnd     int64 // bytes pos of page's ending
	BytesPosSpecial int64 // bytes pos of page's special

	// calculated bits positions
	PosItemsEnd     int64 // bits pos of items end
	PosFreeSpaceEnd int64 // bits pos free space end

	// parsed items positions
	ItemIds []ItemId
}

func DecodePageHeader(page *HeapPage, d *decode.D) {
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
	page.PdPageSizeVersion = uint16(d.FieldU16("pd_pagesize_version"))
	d.FieldU32("pd_prune_xid")

	page.BytesPosSpecial = page.BytesPosBegin + int64(page.PdSpecial)
	page.PosItemsEnd = (page.BytesPosBegin * 8) + int64(page.PdLower*8)
	page.PosFreeSpaceEnd = (page.BytesPosBegin * 8) + int64(page.PdUpper*8)
}

func DecodeItemIds(page *HeapPage, d *decode.D) {
	d.FieldArray("pd_linp", func(d *decode.D) {
		decodeItemIdsInternal(page, d)
	})

	pos0 := d.Pos()
	if pos0 > page.PosFreeSpaceEnd {
		d.Fatalf("items overflows free space")
	}
	freeSpaceLen := page.PosFreeSpaceEnd - pos0
	d.FieldRawLen("free_space", freeSpaceLen, scalar.RawHex)
}

func decodeItemIdsInternal(page *HeapPage, d *decode.D) {
	for {
		checkPos := d.Pos()
		if checkPos >= page.PosItemsEnd {
			break
		}
		/*    0: 0   |     4 */ // unsigned int lp_off: 15
		/*    1: 7   |     4 */ // unsigned int lp_flags: 2
		/*    2: 1   |     4 */ // unsigned int lp_len: 15
		d.FieldStruct("item_id", func(d *decode.D) {
			itemID := ItemId{}

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
