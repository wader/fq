package pgproee

import (
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/common/pg_heap/postgres"
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
func DecodePageHeaderData(page *postgres.HeapPage, d *decode.D) {
	d.FieldStruct("pd_lsn", func(d *decode.D) {
		/*    0      |     4 */ // uint32 xlogid;
		/*    4      |     4 */ // uint32 xrecoff;
		d.FieldU32("xlogid", common.HexMapper)
		d.FieldU32("xrecoff", common.HexMapper)
	})
	page.PdChecksum = uint16(d.FieldU16("pd_checksum"))
	d.FieldU16("pd_flags")
	page.PdLower = uint16(d.FieldU16("pd_lower"))
	page.PdUpper = uint16(d.FieldU16("pd_upper"))
	page.PdSpecial = uint16(d.FieldU16("pd_special"))
	page.PdPageSizeVersion = uint16(d.FieldU16("pd_pagesize_version"))

	page.BytesPosSpecial = page.BytesPosBegin + int64(page.PdSpecial)
	page.PosItemsEnd = (page.BytesPosBegin * 8) + int64(page.PdLower*8)
	page.PosFreeSpaceEnd = (page.BytesPosBegin * 8) + int64(page.PdUpper*8)
}

// type = struct HeapPageSpecialData {
/*    0      |     8 */ // TransactionId pd_xid_base;
/*    8      |     8 */ // TransactionId pd_multi_base;
/*   16      |     4 */ // ShortTransactionId pd_prune_xid;
/*   20      |     4 */ // uint32 pd_magic;
//
/* total size (bytes):   24 */
func DecodePageSpecial(heap *postgres.Heap, d *decode.D) {
	page := heap.Page
	special := heap.Special

	specialPos := page.BytesPosSpecial * 8
	d.SeekAbs(specialPos)

	d.FieldStruct("special_data", func(d *decode.D) {
		special.PdXidBase = d.FieldU64("pd_xid_base")
		special.PdMultiBase = d.FieldU64("pd_multi_base")
		special.PdPruneXid = d.FieldU32("pd_prune_xid")
		special.PdMagic = d.FieldU32("pd_magic")
	})
}
