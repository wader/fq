package postgres14

import (
	"fmt"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

//nolint:revive
const (
	XLOG_BLCKSZ     = 8192
	XLP_LONG_HEADER = 2
)

//nolint:revive
const (
	BKPBLOCK_FORK_MASK = 0x0F
	BKPBLOCK_FLAG_MASK = 0xF0
	BKPBLOCK_HAS_IMAGE = 0x10 /* block data is an XLogRecordBlockImage */
	BKPBLOCK_HAS_DATA  = 0x20
	BKPBLOCK_WILL_INIT = 0x40 /* redo will re-init the page */
	BKPBLOCK_SAME_REL  = 0x80 /* RelFileNode omitted, same as previous */
)

/* Information stored in bimg_info */
//nolint:revive
const (
	BKPIMAGE_HAS_HOLE      = 0x01 /* page image has "hole" */
	BKPIMAGE_IS_COMPRESSED = 0x02 /* page image is compressed */
	BKPIMAGE_APPLY         = 0x04 /* page image should be restored during replay */
)

var rmgrIds = scalar.UToScalar{
	0:  {Sym: "XLOG", Description: "RM_XLOG_ID"},
	1:  {Sym: "Transaction", Description: "RM_XACT_ID"},
	2:  {Sym: "Storage", Description: "RM_SMGR_ID"},
	3:  {Sym: "CLOG", Description: "RM_CLOG_ID"},
	4:  {Sym: "Database", Description: "RM_DBASE_ID"},
	5:  {Sym: "Tablespace", Description: "RM_TBLSPC_ID"},
	6:  {Sym: "MultiXact", Description: "RM_MULTIXACT_ID"},
	7:  {Sym: "RelMap", Description: "RM_RELMAP_ID"},
	8:  {Sym: "Standby", Description: "RM_STANDBY_ID"},
	9:  {Sym: "Heap2", Description: "RM_HEAP2_ID"},
	10: {Sym: "Heap", Description: "RM_HEAP_ID"},
	11: {Sym: "Btree", Description: "RM_BTREE_ID"},
	12: {Sym: "Hash", Description: "RM_HASH_ID"},
	13: {Sym: "Gin", Description: "RM_GIN_ID"},
	14: {Sym: "Gist", Description: "RM_GIST_ID"},
	15: {Sym: "Sequence", Description: "RM_SEQ_ID"},
	16: {Sym: "SPGist", Description: "RM_SPGIST_ID"},
	17: {Sym: "BRIN", Description: "RM_BRIN_ID"},
	18: {Sym: "CommitTs", Description: "RM_COMMIT_TS_ID"},
	19: {Sym: "ReplicationOrigin", Description: "RM_REPLORIGIN_ID"},
	20: {Sym: "Generic", Description: "RM_GENERIC_ID"},
	21: {Sym: "LogicalMessage", Description: "RM_LOGICALMSG_ID"},
}

const (
	XLOG_PAGE_MAGIC_MASK       = 0xD000
	XLOG_PAGE_MAGIC_POSTGRES14 = 0xD10D
)

const (
	XLR_MAX_BLOCK_ID          = 32
	XLR_BLOCK_ID_DATA_SHORT   = 255
	XLR_BLOCK_ID_DATA_LONG    = 254
	XLR_BLOCK_ID_ORIGIN       = 253
	XLR_BLOCK_ID_TOPLEVEL_XID = 252
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

func decodeXLogPageHeaderData(d *decode.D) {
	/*    0      |     2 */ // uint16 xlp_magic;
	/*    2      |     2 */ // uint16 xlp_info;
	/*    4      |     4 */ // TimeLineID xlp_tli;
	/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
	/*   16      |     4 */ // uint32 xlp_rem_len;
	/* XXX  4-byte padding  */
	xlpMagic := d.FieldU16("xlp_magic")
	xlpInfo := d.FieldU16("xlp_info")
	d.FieldU32("xlp_timeline")
	d.FieldU64("xlp_pageaddr")
	d.FieldU32("xlp_rem_len")
	d.FieldU32("padding0")

	if (xlpMagic & XLOG_PAGE_MAGIC_MASK) == 0 {
		d.Fatalf("invalid xlp_magic = %X\n", xlpMagic)
	}

	if (xlpInfo & XLP_LONG_HEADER) != 0 {
		// Long header
		d.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}
}

type walD struct {
	maxOffset int64

	pages   *decode.D
	records *decode.D

	pageRecords *decode.D

	record            *decode.D
	recordRemLenBytes int64
}

func DecodePgwal(d *decode.D, maxOffset uint32) any {
	pages := d.FieldArrayValue("Pages")
	wal := &walD{
		maxOffset:         int64(maxOffset),
		pages:             pages,
		records:           d.FieldArrayValue("Records"),
		recordRemLenBytes: -1, // -1 means not initialized
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

	xLogPage := d.FieldStructValue("Page")

	// type = struct XLogPageHeaderData {
	/*    0      |     2 */ // uint16 xlp_magic;
	/*    2      |     2 */ // uint16 xlp_info;
	/*    4      |     4 */ // TimeLineID xlp_tli;
	/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
	/*   16      |     4 */ // uint32 xlp_rem_len;
	/* XXX  4-byte padding  */
	header := xLogPage.FieldStructValue("XLogPageHeaderData")

	xlpMagic := header.FieldU16("xlp_magic")
	xlpInfo := header.FieldU16("xlp_info")
	header.FieldU32("xlp_tli")
	header.FieldU64("xlp_pageaddr")
	remLenBytes := header.FieldU32("xlp_rem_len")
	header.FieldU32("padding0")

	if (xlpMagic & XLOG_PAGE_MAGIC_MASK) == 0 {
		d.Fatalf("invalid xlp_magic = %X\n", xlpMagic)
	}

	if (xlpInfo & XLP_LONG_HEADER) != 0 {
		// Long header
		header.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}

	if wal.recordRemLenBytes >= 0 { // check recordRemLenBytes is initialized
		if wal.recordRemLenBytes != int64(remLenBytes) {
			d.Fatalf("incorrect wal.recordRemLenBytes = %d, remLenBytes = %d", wal.recordRemLenBytes, remLenBytes)
		}
	}

	remLenBytesAligned := int64(common.TypeAlign8(remLenBytes))
	remLen := remLenBytesAligned * 8

	pos1 := header.Pos()
	xLogPage.SeekAbs(pos1)

	// parted XLogRecord
	if remLen > 0 {
		if wal.record == nil {
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

	if wal.record != nil {
		wal.record.SeekAbs(pos1)
	}

	xLogPage.SeekAbs(pos2)
	pageRecords := xLogPage.FieldArrayValue("Records")

	wal.pageRecords = pageRecords

	decodeXLogRecords(wal, d)
}

func decodeXLogRecords(wal *walD, d *decode.D) {
	pageRecords := wal.pageRecords

	posBytes := d.Pos() / 8
	posMaxOfPageBytes := int64(common.TypeAlign(XLOG_BLCKSZ, uint64(posBytes)))
	fmt.Printf("posMaxOfPageBytes = %d\n", posMaxOfPageBytes)

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
			wal.record = nil
			wal.recordRemLenBytes = 0
			return
		}

		d.SeekAbs(posBytes1Aligned * 8)

		record := pageRecords.FieldStructValue("XLogRecord")
		wal.record = record
		wal.records.AddChild(record.Value)

		xlTotLen := record.FieldU32("xl_tot_len")
		xlTotLen1Bytes := xlTotLen - 4
		pos2Bytes := d.Pos() / 8

		remOnPage := posMaxOfPageBytes - pos2Bytes
		if remOnPage <= 0 {
			d.Fatalf("remOnPage is negative\n")
		}

		if remOnPage < int64(xlTotLen1Bytes) {
			//record.FieldRawLen("xLogBody", remOnPage*8)
			decodeXLogRecord(wal, remOnPage)
			wal.recordRemLenBytes = int64(xlTotLen1Bytes) - remOnPage
			break
		}

		xLogBodyLen := int64(xlTotLen1Bytes) * 8
		if xLogBodyLen <= 0 {
			d.Fatalf("xlTotLen1Bytes is negative, xLogBodyLen = %d\n", xLogBodyLen)
		}

		//record.FieldRawLen("xLogBody", xLogBodyLen)
		decodeXLogRecord(wal, int64(xlTotLen1Bytes))
		wal.record = nil
		wal.recordRemLenBytes = 0
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

func fieldTryGetScalarActualU(d *decode.D, name string, posMax int64, bitsCount int64) (value uint64, end bool) {
	if ok, val := d.FieldTryGetScalarActualU("block_id"); ok {
		value = val
	} else {
		if isEnd(d, posMax, bitsCount) {
			return 0, true
		}
		switch bitsCount {
		case 8:
			value = d.FieldU8(name)
		case 16:
			value = d.FieldU16(name)
		case 24:
			value = d.FieldU24(name)
		case 32:
			value = d.FieldU32(name)
		case 64:
			value = d.FieldU64(name)
		default:
			d.Fatalf("not implemented bitsCount = %d\n", bitsCount)
		}
	}
	return value, false
}

func decodeXLogRecord(wal *walD, maxBytes int64) {
	record := wal.record

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
		record.FieldU64("xl_prev")
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

	//blockId := uint64(0)
	//if ok, val := record.FieldTryGetScalarActualU("block_id"); ok {
	//	blockId = val
	//} else {
	//	if isEnd(record, posMax, 8) {
	//		return
	//	}
	//	blockId = record.FieldU8("block_id")
	//}
	blockId, end := fieldTryGetScalarActualU(record, "block_id", posMax, 8)
	if end {
		return
	}

	if blockId == XLR_BLOCK_ID_DATA_SHORT {
		//typedef struct XLogRecordDataHeaderShort
		//{
		//	uint8		id;				/* XLR_BLOCK_ID_DATA_SHORT */
		//	uint8		data_length;	/* number of payload bytes */
		//}
		//
		/* total size (bytes):   24 */
	}

	//XLR_BLOCK_ID_DATA_SHORT   = 255
	//XLR_BLOCK_ID_DATA_LONG    = 254
	//XLR_BLOCK_ID_ORIGIN       = 253
	//XLR_BLOCK_ID_TOPLEVEL_XID = 252

	mainDataLen := uint64(0)
	recordOrigin := uint64(0)
	toplevelXid := uint64(0)
	if blockId == XLR_BLOCK_ID_DATA_SHORT {
		// COPY_HEADER_FIELD(&main_data_len, sizeof(uint8));
		mainDataLen, end = fieldTryGetScalarActualU(record, "main_data_len", posMax, 8)
		if end {
			return
		}
	} else if blockId == XLR_BLOCK_ID_DATA_LONG {
		// COPY_HEADER_FIELD(&main_data_len, sizeof(uint32));
		mainDataLen, end = fieldTryGetScalarActualU(record, "main_data_len", posMax, 32)
		if end {
			return
		}
	} else if blockId == XLR_BLOCK_ID_ORIGIN {
		// COPY_HEADER_FIELD(&state->record_origin, sizeof(RepOriginId));
		// unsigned short - 2 bytes
		recordOrigin, end = fieldTryGetScalarActualU(record, "record_origin", posMax, 16)
		if end {
			return
		}
	} else if blockId == XLR_BLOCK_ID_TOPLEVEL_XID {
		// COPY_HEADER_FIELD(&state->toplevel_xid, sizeof(TransactionId));
		// 4 bytes
		toplevelXid, end = fieldTryGetScalarActualU(record, "record_origin", posMax, 32)
		if end {
			return
		}
	} else if blockId >= XLR_MAX_BLOCK_ID {
		record.Fatalf("catched blockId = %d\n", blockId)
	} else if blockId < XLR_MAX_BLOCK_ID {
		// COPY_HEADER_FIELD(&fork_flags, sizeof(uint8));
		//forkFlags := uint64(0)
		//if ok, val := record.FieldTryGetScalarActualU("fork_flags"); ok {
		//	forkFlags = val
		//} else {
		//	if isEnd(record, posMax, 8) {
		//		return
		//	}
		//	forkFlags = record.FieldU8("fork_flags")
		//}
		forkFlags, end := fieldTryGetScalarActualU(record, "fork_flags", posMax, 8)
		if end {
			return
		}

		// blk->forknum = fork_flags & BKPBLOCK_FORK_MASK;
		// blk->flags = fork_flags;
		// blk->has_image = ((fork_flags & BKPBLOCK_HAS_IMAGE) != 0);
		// blk->has_data = ((fork_flags & BKPBLOCK_HAS_DATA) != 0);
		hasImage := uint64(0)
		hasData := uint64(0)
		forkNum := forkFlags & BKPBLOCK_FORK_MASK
		if (forkFlags & BKPBLOCK_HAS_IMAGE) != 0 {
			hasImage = 1
		}
		if (forkFlags & BKPBLOCK_HAS_DATA) != 0 {
			hasData = 1
		}
		if record.FieldGet("forknum") == nil {
			record.FieldValueU("forknum", forkNum)
		}
		if record.FieldGet("has_image") == nil {
			record.FieldValueU("has_image", hasImage)
		}
		if record.FieldGet("has_data") == nil {
			record.FieldValueU("has_data", hasData)
		}

		// COPY_HEADER_FIELD(&blk->data_len, sizeof(uint16));
		//dataLen := uint64(0)
		//if ok, val := record.FieldTryGetScalarActualU("data_len"); ok {
		//	dataLen = val
		//} else {
		//	if isEnd(record, posMax, 8) {
		//		return
		//	}
		//	dataLen = record.FieldU8("data_len")
		//}
		dataLen, end := fieldTryGetScalarActualU(record, "data_len", posMax, 8)
		if end {
			return
		}

		// if (blk->has_data && blk->data_len == 0)
		if hasData != 0 && dataLen == 0 {
			record.Fatalf("invalid record with hasData = %d, dataLen = %d\n", hasData, dataLen)
		}
		// if (!blk->has_data && blk->data_len != 0)
		if hasData == 0 && dataLen != 0 {
			record.Fatalf("invalid record with hasData = %d, dataLen = %d\n", hasData, dataLen)
		}

		// if (blk->has_image)
		if hasImage != 0 {
			// COPY_HEADER_FIELD(&blk->bimg_len, sizeof(uint16));
			bimgLen, end := fieldTryGetScalarActualU(record, "bimg_len", posMax, 16)
			if end {
				return
			}

			// COPY_HEADER_FIELD(&blk->hole_offset, sizeof(uint16));
			holeOffset, end := fieldTryGetScalarActualU(record, "hole_offset", posMax, 16)
			if end {
				return
			}

			// COPY_HEADER_FIELD(&blk->bimg_info, sizeof(uint8));
			bimgInfo, end := fieldTryGetScalarActualU(record, "bimg_info", posMax, 8)
			if end {
				return
			}

			// if (blk->bimg_info & BKPIMAGE_IS_COMPRESSED)
			bimgIsCompressed := uint64(0)
			if (bimgInfo & BKPIMAGE_IS_COMPRESSED) != 0 {
				bimgIsCompressed = 1
			}
			if record.FieldGet("bimg_is_compressed") == nil {
				record.FieldValueU("bimg_is_compressed", bimgIsCompressed)
			}

			holeLength := uint64(0)
			bimgHasHole := uint64(0)
			if bimgIsCompressed != 0 {
				if (bimgInfo & BKPIMAGE_HAS_HOLE) != 0 {
					bimgHasHole = 1
				}
				if record.FieldGet("bimg_has_hole") == nil {
					record.FieldValueU("bimg_has_hole", bimgHasHole)
				}
				if bimgHasHole != 0 {
					// COPY_HEADER_FIELD(&blk->hole_length, sizeof(uint16));
					holeLength, end = fieldTryGetScalarActualU(record, "hole_length", posMax, 16)
					if end {
						return
					}
				}
			} else { // bimgIsCompressed is false
				holeLength = XLOG_BLCKSZ - bimgLen
			}
			if record.FieldGet("hole_length") == nil {
				record.FieldValueU("hole_length", holeLength)
			}

			if bimgHasHole != 0 && (holeOffset != 0 || holeLength != 0 || bimgLen == XLOG_BLCKSZ) {
				record.Fatalf("check failed 1")
			}
			if (bimgInfo&BKPIMAGE_HAS_HOLE) == 0 && (holeOffset != 0 || holeLength != 0) {
				record.Fatalf("check failed 2")
			}
			if (bimgInfo&BKPIMAGE_IS_COMPRESSED) != 0 && bimgLen == XLOG_BLCKSZ {
				record.Fatalf("check failed 3")
			}
			if (bimgInfo&BKPIMAGE_HAS_HOLE) == 0 && (bimgInfo&BKPIMAGE_IS_COMPRESSED) == 0 && bimgLen != XLOG_BLCKSZ {
				record.Fatalf("check failed 4")
			}

			if (forkFlags & BKPBLOCK_SAME_REL) == 0 {
				// COPY_HEADER_FIELD(&blk->rnode, sizeof(RelFileNode));

			}

		}
	}

	fmt.Printf("mainDataLen = %d, recordOrigin = %d, toplevelXid = %d\n", mainDataLen, recordOrigin, toplevelXid)

	record.SeekAbs(posMax)
}

func DecodePgwalOri(d *decode.D, in any) any {
	d.SeekAbs(0)

	pageHeaders := d.FieldArrayValue("XLogPageHeaders")
	header := pageHeaders.FieldStruct("XLogPageHeaderData", decodeXLogPageHeaderData)

	xlpRemLen, ok := header.FieldGet("xlp_rem_len").V.(uint32)
	if !ok {
		d.Fatalf("can't get xlp_rem_len\n")
	}

	d.FieldRawLen("prev_file_rec", int64(xlpRemLen*8))
	d.FieldRawLen("prev_file_rec_padding", int64(d.AlignBits(64)))

	d.FieldArray("XLogRecords", func(d *decode.D) {
		for {
			d.FieldStruct("XLogRecord", func(d *decode.D) {
				recordPos := uint64(d.Pos()) >> 3
				recordLen := d.FieldU32("xl_tot_len")
				recordEnd := recordPos + recordLen
				headerPos := recordEnd - recordEnd%XLOG_BLCKSZ
				d.FieldU32("xl_xid")
				d.FieldU64("xl_prev", scalar.ActualHex)
				d.FieldU8("xl_info")
				d.FieldU8("xl_rmid", rmgrIds)
				d.FieldRawLen("padding", int64(d.AlignBits(32)))
				d.FieldU32("xl_crc", scalar.ActualHex)

				var lengths []uint64

				d.FieldArray("XLogRecordBlockHeader", func(d *decode.D) {
					for blkheaderid := uint64(0); d.PeekBits(8) == blkheaderid; blkheaderid++ {
						d.FieldStruct("XlogRecordBlockHeader", func(d *decode.D) {
							/* block reference ID */
							d.FieldU8("id", d.AssertU(blkheaderid))
							/* fork within the relation, and flags */
							forkFlags := d.FieldU8("fork_flags")
							/* number of payload bytes (not including page image) */
							lengths = append(lengths, d.FieldU16("data_length"))
							if forkFlags&BKPBLOCK_HAS_IMAGE != 0 {
								d.FieldStruct("XLogRecordBlockImageHeader", func(d *decode.D) {
									/* number of page image bytes */
									d.FieldU16("length")
									/* number of bytes before "hole" */
									d.FieldU16("hole_offset")
									/* flag bits, see below */
									bimgInfo := d.FieldU8("bimg_info")
									d.FieldRawLen("padding", int64(d.AlignBits(16)))
									if bimgInfo&BKPIMAGE_HAS_HOLE != 0 &&
										bimgInfo&BKPIMAGE_IS_COMPRESSED != 0 {
										d.FieldU16("hole_length")
									}
								})
							}
							if forkFlags&BKPBLOCK_SAME_REL == 0 {
								d.FieldStruct("RelFileNode", func(d *decode.D) {
									/* tablespace */
									d.FieldU32("spcNode")
									/* database */
									d.FieldU32("dbNode")
									/* relation */
									d.FieldU32("relNode")
								})
								d.FieldU32("BlockNum")
							}
						})
					}
				})
				if d.PeekBits(8) == 0xff {
					d.FieldStruct("XLogRecordDataHeaderShort", func(d *decode.D) {
						d.FieldU8("id", d.AssertU(0xff))
						lengths = append(lengths, d.FieldU8("data_length"))
					})
				}

				d.FieldArray("data", func(d *decode.D) {
					for _, x := range lengths {
						pos := uint64(d.Pos()) >> 3
						if pos < headerPos && (headerPos < pos+x) {
							d.FieldRawLen("data", int64((headerPos-pos)*8))
							header := pageHeaders.FieldStruct("XLogPageHeaderData", decodeXLogPageHeaderData)
							_ = header.FieldGet("xlp_rem_len").TryScalarFn(d.ValidateU(recordEnd - headerPos))
							d.FieldRawLen("data", int64((x+pos-headerPos)*8))
						} else {
							d.FieldRawLen("data", int64(x*8))
						}
					}
				})

				d.FieldRawLen("ending_padding", int64(d.AlignBits(64)))
			})
		}
	})

	return nil
}
