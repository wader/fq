package leveldb

// https://github.com/google/leveldb/blob/main/doc/table_format.md
// https://github.com/google/leveldb/blob/main/doc/impl.md
// https://github.com/google/leveldb/blob/main/doc/index.md
//
// Files in LevelDB using this format include:
//  - *.ldb

import (
	"bytes"
	"embed"
	"hash/crc32"

	"github.com/golang/snappy"
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed leveldb_table.md
var leveldbTableFS embed.FS

func init() {
	interp.RegisterFormat(
		format.LevelDB_LDB,
		&decode.Format{
			Description: "LevelDB Table",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    ldbTableDecode,
		})
	interp.RegisterFS(leveldbTableFS)
}

const (
	// four varints (each max 10 bytes) + magic number (8 bytes)
	// https://github.com/google/leveldb/blob/main/table/format.h#L53
	footerEncodedLength = (4*10 + 8) * 8
	magicNumberLength   = 8 * 8
	// leading 64 bits of
	//     echo http://code.google.com/p/leveldb/ | sha1sum
	// https://github.com/google/leveldb/blob/main/table/format.h#L76
	tableMagicNumber = 0xdb4775248b80fb57
	uint32Size       = 32
	uint64Size       = 64
)

// https://github.com/google/leveldb/blob/main/include/leveldb/options.h#L25
const (
	compressionTypeNone      = 0x0
	compressionTypeSnappy    = 0x1
	compressionTypeZstandard = 0x2
)

var compressionTypes = scalar.UintMapSymStr{
	compressionTypeNone:      "none",
	compressionTypeSnappy:    "snappy",
	compressionTypeZstandard: "zstd",
}

// https://github.com/google/leveldb/blob/main/db/dbformat.h#L54
var valueTypes = scalar.UintMapSymStr{
	0x0: "deletion",
	0x1: "value",
}

type blockHandle struct {
	offset uint64
	size   uint64
}

func ldbTableDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	// footer

	var indexOffset int64
	var indexSize int64
	var metaIndexOffset int64
	var metaIndexSize int64

	d.FieldStruct("footer", func(d *decode.D) {
		// check for magic number and fail fast if it isn't there
		d.SeekAbs(d.Len() - magicNumberLength)
		d.FieldU64("magic_number", d.UintAssert(tableMagicNumber), scalar.UintHex)

		d.SeekAbs(d.Len() - footerEncodedLength)
		d.LimitedFn(footerEncodedLength-magicNumberLength, func(d *decode.D) {
			d.FieldStruct("metaindex_handle", func(d *decode.D) {
				metaIndexOffset = int64(d.FieldULEB128("offset"))
				metaIndexSize = int64(d.FieldULEB128("size"))
			})
			d.FieldStruct("index_handle", func(d *decode.D) {
				indexOffset = int64(d.FieldULEB128("offset"))
				indexSize = int64(d.FieldULEB128("size"))
			})
			d.FieldRawLen("padding", d.BitsLeft())
		})
	})

	// metaindex

	d.SeekAbs(metaIndexOffset * 8)
	var metaHandles []blockHandle
	readTableBlock("metaindex", metaIndexSize, keyValueContentsReader(nil, func(size int, d *decode.D) {
		// blockHandle
		// https://github.com/google/leveldb/blob/main/table/format.cc#L24
		handle := blockHandle{
			offset: d.FieldULEB128("offset"),
			size:   d.FieldULEB128("size"),
		}
		metaHandles = append(metaHandles, handle)
	}), d)

	// index

	d.SeekAbs(indexOffset * 8)
	var dataHandles []blockHandle
	readTableBlock("index", indexSize, keyValueContentsReader(readInternalKey, func(size int, d *decode.D) {
		// blockHandle
		// https://github.com/google/leveldb/blob/main/table/format.cc#L24
		handle := blockHandle{
			offset: d.FieldULEB128("offset"),
			size:   d.FieldULEB128("size"),
		}
		dataHandles = append(dataHandles, handle)
	}), d)

	// meta

	if len(metaHandles) > 0 {
		d.FieldArray("meta", func(d *decode.D) {
			for _, handle := range metaHandles {
				d.SeekAbs(int64(handle.offset) * 8)
				readTableBlock("meta_block", int64(handle.size), readMetaContent, d)
			}
		})
	}

	// data

	if len(dataHandles) > 0 {
		d.FieldArray("data", func(d *decode.D) {
			for _, handle := range dataHandles {
				d.SeekAbs(int64(handle.offset) * 8)
				readTableBlock(
					"data_block",
					int64(handle.size),
					keyValueContentsReader(readInternalKey, nil),
					d,
				)
			}
		})
	}

	return nil
}

// Readers

// Read block contents as well as compression + checksum bytes following it.
// The function `readTableBlockContents` gets the uncompressed bytebuffer.
// https://github.com/google/leveldb/blob/main/table/format.cc#L69
func readTableBlock(
	name string,
	size int64,
	readTableBlockContents func(size int64, d *decode.D),
	d *decode.D,
) {
	d.FieldStruct(name, func(d *decode.D) {
		start := d.Pos()
		br := d.RawLen(size * 8)
		// compression (1 byte)
		compressionType := d.FieldU8("compression", compressionTypes, scalar.UintHex)
		// checksum (4 bytes)
		data := d.ReadAllBits(br)
		bytesToCheck := append(data, uint8(compressionType))
		checksum := computeChecksum(bytesToCheck)
		d.FieldU32("checksum", d.UintAssert(uint64(checksum)), scalar.UintHex)
		// decompress if needed
		d.SeekAbs(start)
		if compressionType == compressionTypeNone {
			d.FieldStruct("uncompressed", func(d *decode.D) {
				readTableBlockContents(size, d)
			})
		} else {
			compressedSize := size
			compressed := data
			bb := &bytes.Buffer{}
			switch compressionType {
			case compressionTypeSnappy:
				decompressed, err := snappy.Decode(nil, compressed)
				if err != nil {
					d.Errorf("failed decompressing data: %v", err)
				}
				d.Copy(bb, bytes.NewReader(decompressed))
			default:
				d.Errorf("Unsupported compression type: %x", compressionType)
			}
			if bb.Len() > 0 {
				d.FieldStructRootBitBufFn(
					"uncompressed",
					bitio.NewBitReader(bb.Bytes(), -1),
					func(d *decode.D) {
						readTableBlockContents(int64(bb.Len()), d)
					},
				)
			}
			d.FieldRawLen("compressed", compressedSize*8)
		}

	})
}

// Read content encoded as a sequence of key/value-entries and a trailer of restarts.
// https://github.com/google/leveldb/blob/main/table/block_builder.cc#L16
// https://github.com/google/leveldb/blob/main/table/block.cc
func readKeyValueContents(
	keyCallbackFn func(size int, d *decode.D),
	valueCallbackFn func(size int, d *decode.D),
	size int64,
	d *decode.D,
) {
	start := d.Pos()
	end := start + size*8

	var restartOffset int64
	d.SeekAbs(end - uint32Size)
	d.FieldStruct("trailer", func(d *decode.D) {
		numRestarts := int64(d.FieldU32("num_restarts"))
		restartOffset = size*8 - (1+numRestarts)*uint32Size
		d.SeekAbs(start + restartOffset)
		d.FieldArray("restarts", func(d *decode.D) {
			for i := 0; i < int(numRestarts); i++ {
				d.FieldU32("restart")
			}
		})
	})
	// TK: how do you make an empty entries-array appear _above_ the trailer?
	// Right now, its omitted if empty.
	if restartOffset <= 0 {
		return
	}
	d.SeekAbs(start)
	d.FieldArray("entries", func(d *decode.D) {
		for d.Pos() < start+restartOffset {
			d.FieldStruct("entry", func(d *decode.D) {
				d.FieldULEB128("shared_bytes")
				unshared := int(d.FieldULEB128("unshared_bytes"))
				valueLength := int(d.FieldULEB128("value_length"))
				if keyCallbackFn == nil {
					d.FieldUTF8("key", unshared)
				} else {
					d.FieldStruct("key", func(d *decode.D) {
						keyCallbackFn(unshared, d)
					})
				}
				if valueCallbackFn == nil {
					d.FieldUTF8("value", valueLength)
				} else {
					d.FieldStruct("value", func(d *decode.D) {
						valueCallbackFn(valueLength, d)
					})
				}
			})
		}
	})
}

func readInternalKey(bitSize int, d *decode.D) {
	// InternalKey
	// https://github.com/google/leveldb/blob/main/db/dbformat.h#L171
	d.FieldUTF8("user_key", bitSize-uint64Size/8)
	d.FieldU8("type", valueTypes, scalar.UintHex)
	d.FieldU56("sequence_number")
}

// Read content encoded in the "filter" or "stats" Meta Block format.
// https://github.com/google/leveldb/blob/main/doc/table_format.md#filter-meta-block
// https://github.com/google/leveldb/blob/main/table/filter_block.cc
func readMetaContent(size int64, d *decode.D) {
	// TK(2023-12-04)
	d.FieldRawLen("raw", size*8)
}

// Helpers

func keyValueContentsReader(
	keyCallbackFn func(size int, d *decode.D),
	valueCallbackFn func(size int, d *decode.D),
) func(size int64, d *decode.D) {
	return func(size int64, d *decode.D) {
		readKeyValueContents(keyCallbackFn, valueCallbackFn, size, d)
	}
}

// Compute the checksum: a CRC32 as in RFC3720 + custom mask.
// https://datatracker.ietf.org/doc/html/rfc3720#appendix-B.4
func computeChecksum(bytes []uint8) uint32 {
	crc32C := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	crc32C.Write(bytes)
	return mask(crc32C.Sum32())
}

// Return a masked representation of the CRC.
// https://github.com/google/leveldb/blob/main/util/crc32c.h#L29
func mask(crc uint32) uint32 {
	const kMaskDelta = 0xa282ead8
	// Rotate right by 15 bits and add a constant.
	return ((crc >> 15) | (crc << 17)) + kMaskDelta
}
