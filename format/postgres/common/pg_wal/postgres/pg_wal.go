package postgres

import (
	"fmt"

	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
)

//nolint:revive
const (
	XLOG_BLCKSZ     = 8192
	XLP_LONG_HEADER = 2
)

// struct XLogLongPageHeaderData {
//	/*    0      |    24 */    XLogPageHeaderData std;
//	/*   24      |     8 */    uint64 xlp_sysid;
//	/*   32      |     4 */    uint32 xlp_seg_size;
//	/*   36      |     4 */    uint32 xlp_xlog_blcksz;
//
//	/* total size (bytes):   40 */
//}

// struct XLogPageHeaderData {
/*    0      |     2 */ // uint16 xlp_magic;
/*    2      |     2 */ // uint16 xlp_info;
/*    4      |     4 */ // TimeLineID xlp_tli;
/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
/*   16      |     4 */ // uint32 xlp_rem_len;
/* XXX  4-byte padding  */
//
/* total size (bytes):   24 */

// struct XLogRecord {
/*    0      |     4 */ // uint32 xl_tot_len
/*    4      |     4 */ // TransactionId xl_xid
/*    8      |     8 */ // XLogRecPtr xl_prev
/*   16      |     1 */ // uint8 xl_info
/*   17      |     1 */ // RmgrId xl_rmid
/* XXX  2-byte hole  */
/*   20      |     4 */ // pg_crc32c xl_crc
//
/* total size (bytes):   24 */

// struct RelFileNode {
/*    0      |     4 */ // Oid spcNode
/*    4      |     4 */ // Oid dbNode
/*    8      |     4 */ // Oid relNode
//
/* total size (bytes):   12 */

type Wal struct {
	page *walPage

	pageRecords *decode.D

	State *walState

	DecodeXLogRecord func(wal *Wal, maxBytes int64)
}

type walState struct {
	Record            *decode.D
	recordRemLenBytes int64
}

type walPage struct {
	xlpPageAddr uint64
}

func Decode(d *decode.D, wal *Wal) any {
	pages := d

	for {
		decodeXLogPage(wal, pages)

		if pages.End() {
			break
		}

		posBytes := pages.Pos() / 8
		remBytes := posBytes % XLOG_BLCKSZ
		if remBytes != 0 {
			d.Fatalf("invalid page remBytes = %d\n", remBytes)
		}
	}

	return nil
}

func decodeXLogPage(wal *Wal, d *decode.D) {
	pos0 := d.Pos()
	posPageEnd := pos0 + XLOG_BLCKSZ*8

	d.SeekRel(8 * 8)
	xlpPageAddr0 := d.U64()
	d.SeekAbs(pos0)
	if wal.page != nil {
		xlpPageAddr1 := wal.page.xlpPageAddr + XLOG_BLCKSZ
		if xlpPageAddr0 != xlpPageAddr1 {
			d.Fatalf("invalid xlp_pageaddr expected = %d, actual = %d\n", xlpPageAddr1, xlpPageAddr0)
		}
	}
	wal.page = &walPage{}

	// type = struct XLogPageHeaderData {
	/*    0      |     2 */ // uint16 xlp_magic;
	/*    2      |     2 */ // uint16 xlp_info;
	/*    4      |     4 */ // TimeLineID xlp_tli;
	/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
	/*   16      |     4 */ // uint32 xlp_rem_len;
	/* XXX  4-byte padding  */
	xLogPage := d.FieldStructValue("page")
	header := xLogPage.FieldStructValue("xloog_page_header_data")

	header.FieldU16("xlp_magic")
	xlpInfo := header.FieldU16("xlp_info")
	header.FieldU32("xlp_tli")
	wal.page.xlpPageAddr = header.FieldU64("xlp_pageaddr")
	remLenBytes := header.FieldU32("xlp_rem_len")
	header.FieldU32("padding0")

	//if (xlpMagic & XLOG_PAGE_MAGIC_MASK) == 0 {
	//	d.Fatalf("invalid xlp_magic = %X\n", xlpMagic)
	//}

	if (xlpInfo & XLP_LONG_HEADER) != 0 {
		// Long header
		header.FieldStruct("xlog_long_page_header_data", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}

	if wal.State != nil { // check recordRemLenBytes is initialized
		if wal.State.recordRemLenBytes != int64(remLenBytes) {
			d.Fatalf("recordRemLenBytes = %d != remLenBytes = %d", wal.State.recordRemLenBytes, remLenBytes)
		}
	}

	remLenBytesAligned := int64(common.TypeAlign8(remLenBytes))
	remLen := remLenBytesAligned * 8

	pos1 := header.Pos()
	xLogPage.SeekAbs(pos1)

	maxBitOnPage := posPageEnd - pos1
	if remLen > maxBitOnPage {
		// XLogRecord size is more than page size
		remLen = maxBitOnPage
		remLenBytesAligned = remLen / 8
	}

	// parted XLogRecord
	if remLen > 0 {
		if wal.State == nil {
			// record of previous file
			checkPosBytes := xLogPage.Pos() / 8
			if checkPosBytes >= XLOG_BLCKSZ {
				d.Fatalf("invalid pos of raw_bytes_of_prev_wal_file, pos = %d\n", checkPosBytes)
			}
			xLogPage.FieldRawLen("raw_bytes_of_prev_wal_file", remLen)
		} else {
			// record of previous page
			wal.DecodeXLogRecord(wal, remLenBytesAligned)
		}
	}

	pos2 := xLogPage.Pos()

	if wal.State != nil && wal.State.Record != nil {
		wal.State.Record.SeekAbs(pos1)
	}

	xLogPage.SeekAbs(pos2)
	pageRecords := xLogPage.FieldArrayValue("records")

	wal.pageRecords = pageRecords

	decodeXLogRecords(wal, d)
}

func decodeXLogRecords(wal *Wal, d *decode.D) {
	pageRecords := wal.pageRecords

	posBytes := d.Pos() / 8
	posMaxOfPageBytes := int64(common.TypeAlign(XLOG_BLCKSZ, uint64(posBytes)))

	for {
		/*    0      |     4 */ // uint32 xl_tot_len
		/*    4      |     4 */ // TransactionId xl_xid
		/*    8      |     8 */ // XLogRecPtr xl_prev
		/*   16      |     1 */ // uint8 xl_info
		/*   17      |     1 */ // RmgrId xl_rmid
		/* XXX  2-byte hole  */
		/*   20      |     4 */ // pg_crc32c xl_crc
		posBytes1 := d.Pos() / 8
		posBytes1Aligned := int64(common.TypeAlign8(uint64(posBytes1)))

		// check what we cat read xl_tot_len on this page
		if posMaxOfPageBytes < posBytes1Aligned+4 {
			remOnPage := posMaxOfPageBytes - posBytes1
			d.FieldRawLen("page_padding0", remOnPage*8)
			// can't read xl_tot_len on this page
			// can't create row in this page
			// continue on next page
			wal.State = nil
			return
		}

		if posBytes1 != posBytes1Aligned {
			// ensure align
			d.SeekAbs(posBytes1Aligned * 8)
		}

		record := pageRecords.FieldStructValue("xlog_record")
		wal.State = &walState{
			Record: record,
		}

		lsn0 := uint64(d.Pos() / 8)
		lsn1 := lsn0 % XLOG_BLCKSZ
		lsn := lsn1 + wal.page.xlpPageAddr
		record.FieldValueU("lsn", lsn, common.XLogRecPtrMapper)

		xlTotLen := record.FieldU32("xl_tot_len")
		if xlTotLen < 4 {
			d.Fatalf("xl_tot_len is less than 4\n")
		}
		xlTotLen1Bytes := xlTotLen - 4
		pos2Bytes := d.Pos() / 8

		remOnPage := posMaxOfPageBytes - pos2Bytes
		if remOnPage <= 0 {
			d.Fatalf("remOnPage is negative\n")
		}

		if remOnPage < int64(xlTotLen1Bytes) {
			//record.FieldRawLen("xLogBody", remOnPage*8)
			wal.DecodeXLogRecord(wal, remOnPage)
			wal.State.recordRemLenBytes = int64(xlTotLen1Bytes) - remOnPage
			break
		}

		xLogBodyLen := int64(xlTotLen1Bytes) * 8
		if xLogBodyLen <= 0 {
			errPos := record.Pos() / 8
			d.Fatalf("xlTotLen1Bytes is negative, xLogBodyLen = %d, pos = %X\n", xLogBodyLen, errPos)
		}

		wal.DecodeXLogRecord(wal, int64(xlTotLen1Bytes))

		// align record
		posBytes2 := d.Pos() / 8
		posBytes2Aligned := int64(common.TypeAlign8(uint64(posBytes2)))
		if posBytes2 < posBytes2Aligned {
			alignLen := (posBytes2Aligned - posBytes2) * 8
			wal.State.Record.FieldRawLen("align0", alignLen)
		}

		wal.State = nil
	}
}

// IsEnd - check that we can read bitsCount on page (with posMax?)
func IsEnd(d *decode.D, posMax int64, bitsCount int64) bool {
	pos := d.Pos()
	posRead := pos + bitsCount
	result := posRead > posMax
	if result {
		// set reader at and position to continue reading
		d.SeekAbs(posMax)
	}
	return result
}

func decodeXLogRecord(wal *Wal, maxBytes int64) {
	record := wal.State.Record

	pos0 := record.Pos()
	maxLen := maxBytes * 8
	for i := 0; ; i++ {
		fieldName := fmt.Sprintf("xlog_body%d", i)
		if record.FieldGet(fieldName) == nil {
			record.FieldRawLen(fieldName, maxLen)
			break
		}
	}
	pos1 := record.Pos()
	posMax := pos1
	record.SeekAbs(pos0)

	// struct XLogRecord {
	/*    0      |     4 */ // uint32 xl_tot_len
	/*    4      |     4 */ // TransactionId xl_xid
	/*    8      |     8 */ // XLogRecPtr xl_prev
	/*   16      |     1 */ // uint8 xl_info
	/*   17      |     1 */ // RmgrId xl_rmid
	/* XXX  2-byte hole  */
	/*   20      |     4 */ // pg_crc32c xl_crc

	// xl_tot_len already read

	if record.FieldGet("xl_xid") == nil {
		if IsEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("xl_xid")
	}

	if record.FieldGet("xl_prev") == nil {
		if IsEnd(record, posMax, 64) {
			return
		}
		record.FieldU64("xl_prev", common.XLogRecPtrMapper)
	}

	if record.FieldGet("xl_info") == nil {
		if IsEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_info")
	}

	if record.FieldGet("xl_rmid") == nil {
		if IsEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_rmid")
	}

	if record.FieldGet("hole1") == nil {
		if IsEnd(record, posMax, 16) {
			return
		}
		record.FieldU16("hole1")
	}

	if record.FieldGet("xl_crc") == nil {
		if IsEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("xl_crc")
	}

	record.SeekAbs(posMax)
}
