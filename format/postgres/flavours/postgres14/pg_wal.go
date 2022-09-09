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
	/* block data is an XLogRecordBlockImage */
	BKPBLOCK_HAS_IMAGE = 0x10
	BKPBLOCK_HAS_DATA  = 0x20
	/* redo will re-init the page */
	BKPBLOCK_WILL_INIT = 0x40
	/* RelFileNode omitted, same as previous */
	BKPBLOCK_SAME_REL = 0x80
)

/* Information stored in bimg_info */
//nolint:revive
const (
	/* page image has "hole" */
	BKPIMAGE_HAS_HOLE = 0x01
	/* page image is compressed */
	BKPIMAGE_IS_COMPRESSED = 0x02
	/* page image should be restored during replay */
	BKPIMAGE_APPLY = 0x04
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

/* total size (bytes):   24 */

func decodeXLogPageHeaderData(d *decode.D) {
	var info uint64

	//pages := d.FieldArrayValue("pages")

	//pages.SeekAbs()

	//d.FieldStructValue()

	/*    0      |     2 */ // uint16 xlp_magic;
	/*    2      |     2 */ // uint16 xlp_info;
	/*    4      |     4 */ // TimeLineID xlp_tli;
	/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
	/*   16      |     4 */ // uint32 xlp_rem_len;
	d.FieldU16("xlp_magic")
	d.FieldU16("xlp_info")
	d.FieldU32("xlp_timeline")
	d.FieldU64("xlp_pageaddr")
	d.FieldU32("xlp_rem_len")

	//d.FieldRawLen("padding", int64(d.AlignBits(64)))
	d.U32()

	if info&XLP_LONG_HEADER != 0 {
		// Long header
		d.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}
}

type walD struct {
	pages   *decode.D
	records *decode.D

	pageRecords *decode.D

	record            *decode.D
	recordRemLenBytes int64
}

func DecodePgwal(d *decode.D) any {
	pages := d.FieldArrayValue("Pages")
	walD := &walD{
		pages:             pages,
		records:           d.FieldArrayValue("Records"),
		recordRemLenBytes: -1,
	}

	for {
		decodeXLogPage(walD, pages)

		posBytes := pages.Pos() / 8
		remBytes := posBytes % XLOG_BLCKSZ
		if remBytes != 0 {
			d.Fatalf("invalid page remBytes = %d\n", remBytes)
		}

		if pages.End() {
			break
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

	header.FieldU16("xlp_magic")
	xlpInfo := header.FieldU16("xlp_info")
	header.FieldU32("xlp_tli")
	header.FieldU64("xlp_pageaddr")
	remLenBytes := header.FieldU32("xlp_rem_len")
	header.FieldU32("padding0")

	if xlpInfo&XLP_LONG_HEADER != 0 {
		// Long header
		header.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}

	if wal.recordRemLenBytes >= 0 {
		if wal.recordRemLenBytes != int64(remLenBytes) {
			d.Fatalf("incorrect wal.recordRemLenBytes = %d, remLenBytes = %d", wal.recordRemLenBytes, remLenBytes)
		}
	}

	remLenBytesAligned := common.TypeAlign8(remLenBytes)
	remLen := remLenBytesAligned * 8

	pos1 := header.Pos()
	xLogPage.SeekAbs(pos1)
	// TODO
	xLogPage.FieldRawLen("RecordOfPreviousPage", int64(remLen))
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
	posMaxOfPageBytes := int64(common.TypeAlign(8192, uint64(posBytes)))
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
		posBytes1Aligned := common.TypeAlign8(uint64(posBytes1))
		d.SeekAbs(int64(posBytes1Aligned * 8))

		record := pageRecords.FieldStructValue("XLogRecord")
		wal.record = record
		wal.records.AddChild(record.Value)

		xLogRecordBegin := record.Pos()
		xlTotLen := record.FieldU32("xl_tot_len")
		record.FieldU32("xl_xid")
		record.FieldU64("xl_prev")
		record.FieldU8("xl_info")
		record.FieldU8("xl_rmid")
		record.U16()
		record.FieldU32("xl_crc")
		xLogRecordEnd := record.Pos()
		sizeOfXLogRecord := (xLogRecordEnd - xLogRecordBegin) / 8

		xLogRecordBodyLen := xlTotLen - uint64(sizeOfXLogRecord)

		rawLen := int64(common.TypeAlign8(xLogRecordBodyLen))
		pos1Bytes := d.Pos() / 8

		remOnPage := posMaxOfPageBytes - pos1Bytes
		if remOnPage < rawLen {
			record.FieldRawLen("xLogBody", remOnPage*8)
			wal.recordRemLenBytes = rawLen - remOnPage
			break
		}

		record.FieldRawLen("xLogBody", rawLen*8)
		wal.recordRemLenBytes = -1

		//pos1Bytes := d.Pos() / 8
		//if pos1Bytes > posMaxOfPageBytes {
		//	d.Fatalf("out of page, error in logic!")
		//}

		//pos := d.Pos() / 8
		//if pos >= posMaxOfPage {
		//	break
		//}
		//
		//pageRecords.FieldStruct("XLogRecord", func(d *decode.D) {
		//	record := d
		//	wal.record = record
		//	wal.records.AddChild(record.Value)
		//
		//	xLogRecordBegin := record.Pos()
		//	xlTotLen := record.FieldU32("xl_tot_len")
		//	record.FieldU32("xl_xid")
		//	record.FieldU64("xl_prev")
		//	record.FieldU8("xl_info")
		//	record.FieldU8("xl_rmid")
		//	record.U16()
		//	record.FieldU32("xl_crc")
		//	xLogRecordEnd := record.Pos()
		//	sizeOfXLogRecord := (xLogRecordEnd - xLogRecordBegin) / 8
		//
		//	xLogRecordBodyLen := xlTotLen - uint64(sizeOfXLogRecord)
		//
		//	rawLen := int64(common.TypeAlign8(xLogRecordBodyLen))
		//	record.FieldRawLen("xLogBody", rawLen*8)
		//})

		//pos := d.Pos()
		//if pos >= (4000 * 8) {
		//	break
		//}
	}
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
