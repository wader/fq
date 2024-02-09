package sqlite3

// https://www.sqlite.org/fileformat.html
// https://sqlite.org/src/file?name=src/btreeInt.h&ci=trunk
// https://sqlite.org/schematab.html
// showdb from sqlite tools is also very helpful
// Note that sqlite count pages start at 1 which is at byte 0 to pagesize-1

// . as $r | sqlite3_schema.message.sql | [capture("\\((?<s>.*)\\)").s | split(", ")[] | split(" ")[0] | ascii_downcase] as $c | $r | torepr.message[0] as $m | [$c,$m] | transpose | map({key: .[0], value: .[1]}) | from_entries
// fqlite:
// go run fq.go --arg name "R.E.M." 'torepr as $db | first($db.Artist[] | select(.[1]==$name)) as $artist | $db.Album[] | select(.[2] == $artist[0]) | .[1] | tovalue' format/sqlite3/testdata/chinook.db

// TODO: dont sort array or need "external" decode values?
// TODO: lower case?
// TODO: tovalue?
// TODO: split out cell decode? deep now
// TODO: dummy 0 page to get page 0 at [1] to follow sqlite documentation
// TODO: array/struct external to get sorting correct?
// TODO: overflow pages, two pass?
// TODO: format version
// TODO: table/column names
// TODO: assert version and schema version?
// TODO: ptrmap
// TDOO: wal/journal files? combine?
// TODO: header.unused_space

// > A table with the name "sqlite_sequence" that is used to keep track of the maximum historical INTEGER PRIMARY KEY for a table using AUTOINCREMENT.
// CREATE TABLE sqlite_sequence(name,seq);
// > Tables with names of the form "sqlite_statN" where N is an integer. Such tables store database statistics gathered by the ANALYZE command and used by the query planner to help determine the best algorithm to use for each query.
// CREATE TABLE sqlite_stat1(tbl,idx,stat);
// Only if compiled with SQLITE_ENABLE_STAT2:
// CREATE TABLE sqlite_stat2(tbl,idx,sampleno,sample);
// Only if compiled with SQLITE_ENABLE_STAT3:
// CREATE TABLE sqlite_stat3(tbl,idx,nEq,nLt,nDLt,sample);
// Only if compiled with SQLITE_ENABLE_STAT4:
// CREATE TABLE sqlite_stat4(tbl,idx,nEq,nLt,nDLt,sample);
// TODO: sqlite_autoindex_TABLE_N index

import (
	"bytes"
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed *.jq
var sqlite3FS embed.FS

func init() {
	interp.RegisterFormat(
		format.SQLite3,
		&decode.Format{
			Description: "SQLite v3 database",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    sqlite3Decode,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(sqlite3FS)
}

type intStack struct {
	s []int
}

func (s *intStack) Push(n int) { s.s = append(s.s, n) }

func (s *intStack) Pop() (int, bool) {
	if len(s.s) == 0 {
		return 0, false
	}
	var n int
	n, s.s = s.s[0], s.s[1:]
	return n, true
}

const sqlite3HeaderSize = 100

const (
	serialTypeNULL       = 0
	serialTypeS8         = 1
	serialTypeSBE16      = 2
	serialTypeSBE24      = 3
	serialTypeSBE32      = 4
	serialTypeSBE48      = 5
	serialTypeSBE64      = 6
	serialTypeFloatBE64  = 7
	serialTypeInteger0   = 8
	serialTypeInteger1   = 9
	serialTypeInternal10 = 10
	serialTypeInternal11 = 11
)

var serialTypeMap = scalar.UintMapSymStr{
	serialTypeNULL:       "null",
	serialTypeS8:         "int8",
	serialTypeSBE16:      "int16",
	serialTypeSBE24:      "int24",
	serialTypeSBE32:      "int32",
	serialTypeSBE48:      "int48",
	serialTypeSBE64:      "int64",
	serialTypeFloatBE64:  "float64",
	serialTypeInteger0:   "zero",
	serialTypeInteger1:   "one",
	serialTypeInternal10: "internal10",
	serialTypeInternal11: "internal11",
}

var serialTypeMapper = scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
	typ := uint64(s.Actual)
	if st, ok := serialTypeMap[typ]; ok {
		s.Description = st
	} else if typ >= 12 && typ%2 == 0 {
		s.Description = "blob"
	} else if typ >= 13 && typ%2 != 0 {
		s.Description = "text"
	}
	return s, nil
})

type pageType int

const (
	pageTypePtrmap             pageType = 0x00
	pageTypeBTreeIndexInterior          = 0x02
	pageTypeBTreeTableInterior          = 0x05
	pageTypeBTreeIndexLeaf              = 0x0a
	pageTypeBTreeTableLeaf              = 0x0d
)

var pageTypeMap = scalar.UintMapSymStr{
	// pageTypePtrmap:             "ptrmap",
	pageTypeBTreeIndexInterior: "index_interior",
	pageTypeBTreeTableInterior: "table_interior",
	pageTypeBTreeIndexLeaf:     "index_leaf",
	pageTypeBTreeTableLeaf:     "table_leaf",
}

var ptrmapTypeMap = scalar.UintMapSymStr{
	1: "rootpage",
	2: "freepage",
	3: "overflow1",
	4: "overflow2",
	5: "btree",
}

const (
	textEncodingUTF8    = 1
	textEncodingUTF16LE = 2
	textEncodingUTF16BE = 3
)

var textEncodingMap = scalar.UintMapSymStr{
	textEncodingUTF8:    "utf8",
	textEncodingUTF16LE: "utf16le",
	textEncodingUTF16BE: "utf16be",
}

var versionMap = scalar.UintMapSymStr{
	1: "legacy",
	2: "wal",
}

type sqlite3Header struct {
	pageSize          int64
	databaseSizePages int
	textEncoding      int
}

// TODO: all bits if nine bytes?
// TODO: two complement on bit read count
func varintDecode(d *decode.D) int64 {
	var n uint64
	for i := 0; i < 9; i++ {
		v := d.U8()
		n = n<<7 | v&0b0111_1111
		if v&0b1000_0000 == 0 {
			break
		}
	}
	return mathex.TwosComplement(64, n)
}

func sqlite3DecodeSerialType(d *decode.D, h sqlite3Header, typ int64) {
	switch typ {
	case serialTypeNULL:
		d.FieldValueAny("value", nil)
	case serialTypeS8:
		d.FieldS8("value")
	case serialTypeSBE16:
		d.FieldS16("value")
	case serialTypeSBE24:
		d.FieldS24("value")
	case serialTypeSBE32:
		d.FieldS32("value")
	case serialTypeSBE48:
		d.FieldS48("value")
	case serialTypeSBE64:
		d.FieldS64("value")
	case serialTypeFloatBE64:
		d.FieldF64("value")
	case serialTypeInteger0:
		d.FieldValueAny("value", 0)
	case serialTypeInteger1:
		d.FieldValueAny("value", 1)
	case 10, 11:
		// internal, should not appear in wellformed file
	default:
		if typ%2 == 0 {
			// N => 12 and even: (N-12)/2 bytes blob.
			d.FieldRawLen("value", (typ-12)/2*8)
		} else {
			// N => 13 and odd: (N-13)/2 bytes text
			l := int(typ-13) / 2
			switch h.textEncoding {
			case textEncodingUTF8:
				d.FieldUTF8("value", l)
			case textEncodingUTF16LE:
				d.FieldUTF16LE("value", l)
			case textEncodingUTF16BE:
				d.FieldUTF16BE("value", l)
			}
		}
	}
}

func sqlite3DecodeCellFreeblock(d *decode.D) uint64 {
	nextOffset := d.FieldU16("next_offset")
	if nextOffset == 0 {
		return 0
	}
	// TODO: "header" is size bytes or offset+size? seems to be just size
	// "size of the freeblock in bytes, including the 4-byte header"
	size := d.FieldU16("size")
	// TODO: really?
	if size == 0 {
		return 0
	}
	d.FieldRawLen("space", int64(size-4)*8)
	return nextOffset
}

func sqlite3CellPayloadDecode(d *decode.D, h sqlite3Header) {
	lengthStart := d.Pos()
	length := d.FieldSintFn("length", varintDecode)
	lengthBits := d.Pos() - lengthStart
	var serialTypes []int64
	d.FramedFn((length)*8-lengthBits, func(d *decode.D) {
		d.FieldArray("serials", func(d *decode.D) {
			for !d.End() {
				serialTypes = append(
					serialTypes,
					d.FieldSintFn("serial", varintDecode, serialTypeMapper),
				)
			}
		})
	})
	d.FieldArray("contents", func(d *decode.D) {
		for _, s := range serialTypes {
			sqlite3DecodeSerialType(d, h, s)
		}
	})
}

func sqlite3DecodeTreePage(d *decode.D, h sqlite3Header, x int64, payLoadLen int64) {
	// formulas from sqlite format spec
	u := h.pageSize
	p := payLoadLen
	m := ((u - 12) * 32 / 255) - 23
	k := m + ((p - m) % (u - 4))

	var firstPayLoadLen int64
	if k <= x {
		firstPayLoadLen = k
	} else {
		firstPayLoadLen = m
	}

	if p <= x {
		// payload fits in page
		d.FramedFn(firstPayLoadLen*8, func(d *decode.D) {
			d.FieldStruct("payload", func(d *decode.D) { sqlite3CellPayloadDecode(d, h) })
		})
	} else {
		// payload overflows, collect payload parts
		payLoadBB := &bytes.Buffer{}

		d.FieldArray("overflow_pages", func(d *decode.D) {
			var nextPage int64
			d.FieldStruct("overflow_page", func(d *decode.D) {
				br := d.FieldRawLen("data", firstPayLoadLen*8)
				nextPage = d.FieldS32("next_page")
				d.CopyBits(payLoadBB, br)
			})

			payLoadLenLeft := payLoadLen - firstPayLoadLen
			for nextPage != 0 {
				d.SeekAbs(((nextPage - 1) * h.pageSize) * 8)
				d.FieldStruct("overflow_page", func(d *decode.D) {
					nextPage = d.FieldS32("next_page")
					overflowSize := mathex.Min(h.pageSize-4, payLoadLenLeft)
					br := d.FieldRawLen("data", overflowSize*8)
					payLoadLenLeft -= overflowSize
					d.CopyBits(payLoadBB, br)
				})
			}
		})

		d.FieldStructRootBitBufFn("payload",
			bitio.NewBitReader(payLoadBB.Bytes(), -1),
			func(d *decode.D) { sqlite3CellPayloadDecode(d, h) },
		)
	}
}

func sqlite3SeekPage(d *decode.D, h sqlite3Header, i int) {
	pageOffset := h.pageSize * int64(i)
	if i == 0 {
		pageOffset += sqlite3HeaderSize
	}
	d.SeekAbs(pageOffset * 8)
}

func sqlite3Decode(d *decode.D) any {
	var h sqlite3Header

	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 16, d.StrAssert("SQLite format 3\x00"))
		pageSizeS := d.FieldScalarU16("page_size", scalar.UintMapSymUint{1: 65536}) // in bytes. Must be a power of two between 512 and 32768 inclusive, or the value 1 representing a page size of 65536.
		d.FieldU8("write_version", versionMap)                                      // 1 for legacy; 2 for WAL.
		d.FieldU8("read_version", versionMap)                                       // . 1 for legacy; 2 for WAL.
		d.FieldU8("unused_space")                                                   // at the end of each page. Usually 0.
		d.FieldU8("maximum_embedded_payload_fraction")                              // . Must be 64.
		d.FieldU8("minimum_embedded_payload_fraction")                              // . Must be 32.
		d.FieldU8("leaf_payload_fraction")                                          // . Must be 32.
		d.FieldU32("file_change_counter")                                           //
		databaseSizePages := int(d.FieldU32("database_size_pages"))                 // . The "in-header database size".
		d.FieldU32("page_number_freelist")                                          // of the first freelist trunk page.
		d.FieldU32("total_number_freelist")                                         // pages.
		d.FieldU32("schema_cookie")                                                 // .
		d.FieldU32("schema_format_number")                                          // . Supported schema formats are 1, 2, 3, and 4.
		d.FieldU32("default_page_cache_size")                                       // .
		d.FieldU32("page_number_largest_root_btree")                                // page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.
		textEncoding := int(d.FieldU32("text_encoding", textEncodingMap))
		d.FieldU32("user_version")                       // " as read and set by the user_version pragma.
		d.FieldU32("incremental_vacuum_mode")            // False (zero) otherwise.
		d.FieldU32("application_id")                     // " set by PRAGMA application_id.
		d.FieldRawLen("reserved", 160, d.BitBufIsZero()) // for expansion. Must be zero.
		d.FieldU32("version_valid_for")                  // number.
		d.FieldU32("sqlite_version_number")              //

		// TODO: nicer API for fallback?
		pageSize := int64(pageSizeS.Actual)
		if pageSizeS.Sym != nil {
			pageSize = int64(pageSizeS.SymUint())
		}

		h = sqlite3Header{
			pageSize:          pageSize,
			databaseSizePages: databaseSizePages,
			textEncoding:      textEncoding,
		}
	})

	// pageTypes := map[int]pageType{}
	// pageVisitStack := &intStack{}
	// pageVisitStack.Push(0)

	// for {
	// 	i, ok := pageVisitStack.Pop()
	// 	if !ok {
	// 		break
	// 	}
	// 	if _, ok := pageTypes[i]; ok {
	// 		d.Fatalf("page %d already visited", i)
	// 	}

	// 	sqlite3SeekPage(d, h, i)
	// 	typ := d.U8()

	// 	switch typ {
	// 	case pageTypeBTreeIndexInterior,
	// 		pageTypeBTreeTableInterior:

	// 		d.U16() // start_free_blocks
	// 		d.U16() // cell_start
	// 		d.U8()  // cell_fragments
	// 		rightPointer := d.U32()

	// 		pageCells := d.U16()
	// 		for i := uint64(0); i < pageCells; i++ {

	// 		}

	// 		switch typ {
	// 		case pageTypeBTreeIndexInterior:

	// 		}

	// 	default:
	// 		d.Fatalf("asd")
	// 	}

	// }

	// return nil

	d.FieldArray("pages", func(d *decode.D) {
		// add a filler entry to make real pages start at index 1
		d.FieldStruct("page", func(d *decode.D) {
			d.FieldValueStr("type", "page0_index_fill")
		})

		// for {
		// i, ok := pageStack.Pop()
		// if !ok {
		// 	break
		// }
		// if _, ok := pageSeen[i]; ok {
		// 	d.Fatalf("page %d already visited", i)
		// }
		// pageSeen[i] = struct{}{}

		for i := 0; i < h.databaseSizePages; i++ {
			pageOffset := h.pageSize * int64(i)
			d.SeekAbs(pageOffset * 8)
			// skip header for first page
			if i == 0 {
				d.SeekRel(sqlite3HeaderSize * 8)
			}
			sqlite3SeekPage(d, h, i)

			d.FieldStruct("page", func(d *decode.D) {
				typ := d.FieldU8("type", pageTypeMap)
				switch typ {
				// case pageTypePtrmap:
				// TODO: how to know if just a overflow page?
				// log.Printf("ptrmap i: %#+v\n", i)
				// d.FieldArray("entries", func(d *decode.D) {
				// 	for j := int64(0); j < h.pageSize/5; j++ {
				// 		d.FieldStruct("entry", func(d *decode.D) {
				// 			d.FieldU8("type", ptrmapTypeMap)
				// 			d.FieldU32("page_number")
				// 		})
				// 	}
				// })
				default:
					d.FieldRawLen("data", (h.pageSize-4)*8)

				case pageTypeBTreeIndexInterior,
					pageTypeBTreeIndexLeaf,
					pageTypeBTreeTableInterior,
					pageTypeBTreeTableLeaf:

					startFreeblocks := d.FieldU16("start_freeblocks") // The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
					pageCells := d.FieldU16("page_cells")             // The two-byte integer at offset 3 gives the number of cells on the page.
					d.FieldU16("cell_start")                          // sThe two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
					d.FieldU8("cell_fragments")                       // The one-byte integer at offset 7 gives the number of fragmented free bytes within the cell content area.

					switch typ {
					case pageTypeBTreeIndexInterior,
						pageTypeBTreeTableInterior:
						d.FieldU32("right_pointer") // The four-byte page number at offset 8 is the right-most pointer. This value appears in the header of interior b-tree pages only and is omitted from all other pages.
					}
					var cellPointers []uint64
					d.FieldArray("cells_pointers", func(d *decode.D) {
						for j := uint64(0); j < pageCells; j++ {
							cellPointers = append(cellPointers, d.FieldU16("pointer"))
						}
					})
					if startFreeblocks != 0 {
						d.FieldArray("freeblocks", func(d *decode.D) {
							nextOffset := startFreeblocks
							for nextOffset != 0 {
								d.SeekAbs((pageOffset + int64(nextOffset)) * 8)
								d.FieldStruct("freeblock", func(d *decode.D) {
									nextOffset = sqlite3DecodeCellFreeblock(d)
								})
							}
						})
					}
					d.FieldArray("cells", func(d *decode.D) {
						for _, p := range cellPointers {
							d.FieldStruct("cell", func(d *decode.D) {
								// TODO: SeekAbs with fn later?
								d.SeekAbs((pageOffset + int64(p)) * 8)
								switch typ {
								case pageTypeBTreeIndexInterior:
									d.FieldU32("left_child")
									payLoadLen := d.FieldSintFn("payload_len", varintDecode)
									// formula for x from sqlite format spec
									sqlite3DecodeTreePage(d, h, ((h.pageSize-12)*64/255)-23, payLoadLen)
								case pageTypeBTreeTableInterior:
									d.FieldU32("left_child")
									d.FieldSintFn("rowid", varintDecode)
								case pageTypeBTreeIndexLeaf:
									payLoadLen := d.FieldSintFn("payload_len", varintDecode)
									sqlite3DecodeTreePage(d, h, ((h.pageSize-12)*64/255)-23, payLoadLen)
								case pageTypeBTreeTableLeaf:
									payLoadLen := d.FieldSintFn("payload_len", varintDecode)
									d.FieldSintFn("rowid", varintDecode)
									sqlite3DecodeTreePage(d, h, h.pageSize-35, payLoadLen)
								}
							})
						}
					})
				}
			})
		}
	})

	return nil
}
