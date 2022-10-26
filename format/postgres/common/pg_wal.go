package common

import (
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

type walD struct {
	maxOffset int64
	page      *walPage

	pageRecords *decode.D

	state *walState
}

type walState struct {
	record            *decode.D
	recordRemLenBytes int64
}

type walPage struct {
	xlpPageAddr uint64
}

func DecodePGWAL(d *decode.D, maxOffset uint32) any {
	pages := d.FieldArrayValue("Pages")
	wal := &walD{
		maxOffset: int64(maxOffset),
	}

	for {
		decodeXLogPage(wal, pages)

		if pages.End() {
			break
		}

		posBytes := pages.Pos() / 8
		if posBytes >= wal.maxOffset {
			d.FieldRawLen("unused", d.BitsLeft())
			break
		}

		remBytes := posBytes % XLOG_BLCKSZ
		if remBytes != 0 {
			d.Fatalf("invalid page remBytes = %d\n", remBytes)
		}
	}

	return nil
}

func decodeXLogPage(wal *walD, d *decode.D) {
	pos0 := d.Pos()
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
	xLogPage := d.FieldStructValue("Page")
	header := xLogPage.FieldStructValue("XLogPageHeaderData")

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
		header.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}

	if wal.state != nil { // check recordRemLenBytes is initialized
		if wal.state.recordRemLenBytes != int64(remLenBytes) {
			d.Fatalf("recordRemLenBytes = %d != remLenBytes = %d", wal.state.recordRemLenBytes, remLenBytes)
		}
	}

	remLenBytesAligned := int64(TypeAlign8(remLenBytes))
	remLen := remLenBytesAligned * 8

	pos1 := header.Pos()
	xLogPage.SeekAbs(pos1)

	// parted XLogRecord
	if remLen > 0 {
		if wal.state == nil {
			// record of previous file
			checkPosBytes := xLogPage.Pos() / 8
			if checkPosBytes >= XLOG_BLCKSZ {
				d.Fatalf("invalid pos for RawBytesOfPreviousWalFile, it must be on first page only, pos = %d\n", checkPosBytes)
			}
			xLogPage.FieldRawLen("RawBytesOfPreviousWalFile", remLen)
		} else {
			// record of previous page
			decodeXLogRecord(wal, remLenBytesAligned)
		}
	}

	pos2 := xLogPage.Pos()

	if wal.state != nil && wal.state.record != nil {
		wal.state.record.SeekAbs(pos1)
	}

	xLogPage.SeekAbs(pos2)
	pageRecords := xLogPage.FieldArrayValue("Records")

	wal.pageRecords = pageRecords

	decodeXLogRecords(wal, d)
}

func decodeXLogRecords(wal *walD, d *decode.D) {
	pageRecords := wal.pageRecords

	posBytes := d.Pos() / 8
	posMaxOfPageBytes := int64(TypeAlign(XLOG_BLCKSZ, uint64(posBytes)))

	for {
		/*    0      |     4 */ // uint32 xl_tot_len
		/*    4      |     4 */ // TransactionId xl_xid
		/*    8      |     8 */ // XLogRecPtr xl_prev
		/*   16      |     1 */ // uint8 xl_info
		/*   17      |     1 */ // RmgrId xl_rmid
		/* XXX  2-byte hole  */
		/*   20      |     4 */ // pg_crc32c xl_crc
		posBytes1 := d.Pos() / 8
		posBytes1Aligned := int64(TypeAlign8(uint64(posBytes1)))
		// check aligned - this is correct
		// record header is 8 byte aligned
		if posBytes1Aligned >= wal.maxOffset {
			d.FieldRawLen("unused", d.BitsLeft())
			break
		}

		// check what we cat read xl_tot_len on this page
		if posMaxOfPageBytes < posBytes1Aligned+4 {
			remOnPage := posMaxOfPageBytes - posBytes1
			d.FieldRawLen("page_padding0", remOnPage*8)
			// can't read xl_tot_len on this page
			// can't create row in this page
			// continue on next page
			wal.state = nil
			return
		}

		d.SeekAbs(posBytes1Aligned * 8)

		record := pageRecords.FieldStructValue("XLogRecord")
		wal.state = &walState{
			record: record,
		}

		lsn0 := uint64(d.Pos() / 8)
		lsn1 := lsn0 % XLOG_BLCKSZ
		lsn := lsn1 + wal.page.xlpPageAddr
		record.FieldValueU("lsn", lsn, XLogRecPtrMapper)

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
			decodeXLogRecord(wal, remOnPage)
			wal.state.recordRemLenBytes = int64(xlTotLen1Bytes) - remOnPage
			break
		}

		xLogBodyLen := int64(xlTotLen1Bytes) * 8
		if xLogBodyLen <= 0 {
			errPos := record.Pos() / 8
			d.Fatalf("xlTotLen1Bytes is negative, xLogBodyLen = %d, pos = %X\n", xLogBodyLen, errPos)
		}

		//record.FieldRawLen("xLogBody", xLogBodyLen)
		decodeXLogRecord(wal, int64(xlTotLen1Bytes))
		wal.state = nil
	}
}

// check that we can read bitsCount on page (with posMax?)
func isEnd(d *decode.D, posMax int64, bitsCount int64) bool {
	pos := d.Pos()
	posRead := pos + bitsCount
	result := posRead > posMax
	if result {
		// set reader at and position to continue reading
		d.SeekAbs(posMax)
	}
	return result
}

func decodeXLogRecord(wal *walD, maxBytes int64) {
	record := wal.state.record

	pos0 := record.Pos()
	maxLen := maxBytes * 8
	if record.FieldGet("xLogBody0") == nil {
		// body on first page
		record.FieldRawLen("xLogBody0", maxLen)
	} else {
		// body on second page
		record.FieldRawLen("xLogBody1", maxLen)
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
		if isEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("xl_xid")
	}

	if record.FieldGet("xl_prev") == nil {
		if isEnd(record, posMax, 64) {
			return
		}
		record.FieldU64("xl_prev", XLogRecPtrMapper)
	}

	if record.FieldGet("xl_info") == nil {
		if isEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_info")
	}

	if record.FieldGet("xl_rmid") == nil {
		if isEnd(record, posMax, 8) {
			return
		}
		record.FieldU8("xl_rmid")
	}

	if record.FieldGet("hole1") == nil {
		if isEnd(record, posMax, 16) {
			return
		}
		record.FieldU16("hole1")
	}

	if record.FieldGet("xl_crc") == nil {
		if isEnd(record, posMax, 32) {
			return
		}
		record.FieldU32("xl_crc")
	}

	record.SeekAbs(posMax)
}
