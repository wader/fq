package postgres14

import (
	"context"
	"fmt"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

//func init() {
//	interp.RegisterFormat(decode.Format{
//		Name:        format.PGWAL,
//		Description: "PostgreSQL write-ahead log file",
//		DecodeFn:    pgwalDecode,
//	})
//}

const XLOG_BLCKSZ = 8192

const XLP_LONG_HEADER = 2

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
const (
	/* page image has "hole" */
	BKPIMAGE_HAS_HOLE = 0x01
	/* page image is compressed */
	BKPIMAGE_IS_COMPRESSED = 0x02
	/* page image should be restored during replay */
	BKPIMAGE_APPLY = 0x04
)

var expected_rem_len uint64 = 0

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

// type = struct XLogPageHeaderData {
/*    0      |     2 */ // uint16 xlp_magic;
/*    2      |     2 */ // uint16 xlp_info;
/*    4      |     4 */ // TimeLineID xlp_tli;
/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
/*   16      |     4 */ // uint32 xlp_rem_len;
/* XXX  4-byte padding  */
//
/* total size (bytes):   24 */

// type = struct XLogRecord {
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
	remLen      uint32

	record *decode.D
}

func getWalD(d *decode.D) *walD {
	val := d.Ctx.Value("wald")
	return val.(*walD)
}

func DecodePgwal(d *decode.D, in any) any {
	walD := &walD{
		pages:   d.FieldArrayValue("pages"),
		records: d.FieldArrayValue("records"),
	}
	parentCtx := d.Ctx
	ctx := context.WithValue(parentCtx, "wald", walD)
	d.Ctx = ctx

	d.SeekAbs(0)
	d.FieldArray("XLogPages", decodeXLogPage)

	return nil
}

func decodeXLogPage(d *decode.D) {

	wal := getWalD(d)

	// type = struct XLogPageHeaderData {
	/*    0      |     2 */ // uint16 xlp_magic;
	/*    2      |     2 */ // uint16 xlp_info;
	/*    4      |     4 */ // TimeLineID xlp_tli;
	/*    8      |     8 */ // XLogRecPtr xlp_pageaddr;
	/*   16      |     4 */ // uint32 xlp_rem_len;
	/* XXX  4-byte padding  */
	page := wal.pages.FieldStructValue("XLogPageHeaderData")

	page.FieldU16("xlp_magic")
	xlpInfo := page.FieldU16("xlp_info")
	page.FieldU32("xlp_tli")
	page.FieldU64("xlp_pageaddr")
	remLen := page.FieldU32("xlp_rem_len")
	page.U32()

	if xlpInfo&XLP_LONG_HEADER != 0 {
		// Long header
		d.FieldStruct("XLogLongPageHeaderData", func(d *decode.D) {
			d.FieldU64("xlp_sysid")
			d.FieldU32("xlp_seg_size")
			d.FieldU32("xlp_xlog_blcksz")
		})
	}

	remLen = 40
	wal.remLen = uint32(remLen)

	record := wal.record
	if record == nil {
		rawLen := int64(common.TypeAlign8(remLen))
		page.FieldRawLen("prev_file_rec", rawLen*8)
	}

	pageRecords := page.FieldArrayValue("records")
	wal.pageRecords = pageRecords

	decodeXLogRecords(d)

	//page.Pos()
	//for {
	//
	//}
	//fmt.Printf("d pos = %d\n", d.Pos())
}

func decodeXLogRecords(d *decode.D) {
	wal := getWalD(d)
	pageRecords := wal.pageRecords

	pos := d.Pos() / 8
	posMaxOfPage := int64(common.TypeAlign(8192, uint64(pos)))
	fmt.Printf("posMaxOfPage = %d\n", posMaxOfPage)

	for {
		/*    0      |     4 */ // uint32 xl_tot_len
		/*    4      |     4 */ // TransactionId xl_xid
		/*    8      |     8 */ // XLogRecPtr xl_prev
		/*   16      |     1 */ // uint8 xl_info
		/*   17      |     1 */ // RmgrId xl_rmid
		/* XXX  2-byte hole  */
		/*   20      |     4 */ // pg_crc32c xl_crc

		//record := page.FieldStructValue("XLogRecord")
		//wal.record = record
		//wal.records.AddChild(record.Value)
		//
		//xLogRecordBegin := record.Pos()
		//xlTotLen := record.FieldU32("xl_tot_len")
		//record.FieldU32("xl_xid")
		//record.FieldU64("xl_prev")
		//record.FieldU8("xl_info")
		//record.FieldU8("xl_rmid")
		//record.U16()
		//record.FieldU32("xl_crc")
		//xLogRecordEnd := record.Pos()
		//sizeOfXLogRecord := (xLogRecordEnd - xLogRecordBegin) / 8
		//
		//xLogRecordBodyLen := xlTotLen - uint64(sizeOfXLogRecord)
		//
		//rawLen := int64(TypeAlign8(xLogRecordBodyLen))
		//page.FieldRawLen("xLogBody", rawLen*8)

		pos := d.Pos() / 8
		if pos >= posMaxOfPage {
			break
		}

		pageRecords.FieldStruct("XLogRecord", func(d *decode.D) {
			record := d
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
			record.FieldRawLen("xLogBody", rawLen*8)
		})

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

	d.FieldRawLen("prev_file_rec", int64(header.FieldGet("xlp_rem_len").V.(uint32)*8))
	d.FieldRawLen("prev_file_rec_padding", int64(d.AlignBits(64)))

	d.FieldArray("XLogRecords", func(d *decode.D) {
		for {
			d.FieldStruct("XLogRecord", func(d *decode.D) {
				record_pos := uint64(d.Pos()) >> 3
				record_len := d.FieldU32("xl_tot_len")
				record_end := record_pos + record_len
				header_pos := record_end - record_end%XLOG_BLCKSZ
				d.FieldU32("xl_xid")
				d.FieldU64("xl_prev", scalar.ActualHex)
				d.FieldU8("xl_info")
				d.FieldU8("xl_rmid", rmgrIds)
				d.FieldRawLen("padding", int64(d.AlignBits(32)))
				d.FieldU32("xl_crc", scalar.ActualHex)

				var lenghts []uint64 = []uint64{}

				d.FieldArray("XLogRecordBlockHeader", func(d *decode.D) {
					for blkheaderid := uint64(0); d.PeekBits(8) == blkheaderid; blkheaderid++ {
						d.FieldStruct("XlogRecordBlockHeader", func(d *decode.D) {
							/* block reference ID */
							d.FieldU8("id", d.AssertU(blkheaderid))
							/* fork within the relation, and flags */
							fork_flags := d.FieldU8("fork_flags")
							/* number of payload bytes (not including page image) */
							lenghts = append(lenghts, d.FieldU16("data_length"))
							if fork_flags&BKPBLOCK_HAS_IMAGE != 0 {
								d.FieldStruct("XLogRecordBlockImageHeader", func(d *decode.D) {
									/* number of page image bytes */
									d.FieldU16("length")
									/* number of bytes before "hole" */
									d.FieldU16("hole_offset")
									/* flag bits, see below */
									bimg_info := d.FieldU8("bimg_info")
									d.FieldRawLen("padding", int64(d.AlignBits(16)))
									if bimg_info&BKPIMAGE_HAS_HOLE != 0 &&
										bimg_info&BKPIMAGE_IS_COMPRESSED != 0 {
										d.FieldU16("hole_length")
									}
								})
							}
							if fork_flags&BKPBLOCK_SAME_REL == 0 {
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
						lenghts = append(lenghts, d.FieldU8("data_length"))
					})
				}

				d.FieldArray("data", func(d *decode.D) {
					for _, x := range lenghts {
						pos := uint64(d.Pos()) >> 3
						if pos < header_pos && (header_pos < pos+x) {
							d.FieldRawLen("data", int64((header_pos-pos)*8))
							header := pageHeaders.FieldStruct("XLogPageHeaderData", decodeXLogPageHeaderData)
							header.FieldGet("xlp_rem_len").TryScalarFn(d.ValidateU(record_end - header_pos))
							d.FieldRawLen("data", int64((x+pos-header_pos)*8))
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
