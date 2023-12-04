package ldb

// https://github.com/google/leveldb/blob/main/doc/table_format.md
// https://github.com/google/leveldb/blob/main/doc/impl.md
// https://github.com/google/leveldb/blob/main/doc/index.md

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.LDB,
		&decode.Format{
			Description: "LevelDB Table",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    ldbDecode,
		})
}

const (
	// four varints (each max 10 bytes) + magic number (8 bytes)
	// https://github.com/google/leveldb/blob/main/table/format.h#L53
	footerEncodedLength = 4*10 + 8
	// leading 64 bits of
	//     echo http://code.google.com/p/leveldb/ | sha1sum
	// https://github.com/google/leveldb/blob/main/table/format.h#L76
	tableMagicNumber = 0xdb4775248b80fb57
	// 1-byte compression type + 4-bytes CRC
	// https://github.com/google/leveldb/blob/main/table/format.h#L79
	blockTrailerSize = 5
)

// https://github.com/google/leveldb/blob/main/include/leveldb/options.h#L25
var compressionTypes = scalar.UintMapSymStr{
	0x0: "none",
	0x1: "Snappy",
	0x2: "Zstandard",
}

// https://github.com/google/leveldb/blob/main/db/dbformat.h#L54
var valueTypes = scalar.UintMapSymStr{
	0x0: "deletion",
	0x1: "value",
}

type BlockHandle struct {
	Offset uint64
	Size   uint64
}

func ldbDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	// Read the footer (last 48 bytes)
	d.SeekAbs(d.Len() - footerEncodedLength*8)
	var indexOffset int64
	var indexSize int64
	var metaIndexOffset int64
	var metaIndexSize int64

	d.FieldStruct("footer", func(d *decode.D) {
		// Extract varints for metaindex offset and size, index offset and size
		d.FieldStruct("metaindex_handle", func(d *decode.D) {
			metaIndexOffset = int64(d.FieldUintFn("offset", decodeVarInt))
			metaIndexSize = int64(d.FieldUintFn("size", decodeVarInt))
		})
		d.FieldStruct("index_handle", func(d *decode.D) {
			indexOffset = int64(d.FieldUintFn("offset", decodeVarInt))
			indexSize = int64(d.FieldUintFn("size", decodeVarInt))
		})
		d.FieldRawLen("padding", d.Len()-d.Pos()-8*8)
		d.FieldU64("magic_number", d.UintAssert(tableMagicNumber), scalar.UintHex)
	})

	d.SeekAbs(metaIndexOffset * 8)
	fieldStructBlock("metaindex_block", metaIndexSize, nil, d)

	d.SeekAbs(indexOffset * 8)
	var dataHandles []BlockHandle
	fieldStructBlock("index_block", indexSize, func(d *decode.D) {
		// BlockHandle
		// https://github.com/google/leveldb/blob/main/table/format.cc#L24
		handle := BlockHandle{
			Offset: d.FieldUintFn("offset", decodeVarInt),
			Size:   d.FieldUintFn("size", decodeVarInt),
		}
		dataHandles = append(dataHandles, handle)
	}, d)

	fmt.Println("total handles", len(dataHandles))
	d.FieldArray("data_blocks", func(d *decode.D) {
		for _, handle := range dataHandles {
			d.SeekAbs(int64(handle.Offset) * 8)
			fieldStructBlock("data_block", int64(handle.Size), nil, d)
		}
	})

	return nil
}

// Helpers

func fieldStructBlock(name string, size int64, valueCallbackFn func(d *decode.D), d *decode.D) *decode.D {
	// ReadBlock: https://github.com/google/leveldb/blob/main/table/format.cc#L69
	uint32Size := int64(32)
	uint64Size := int64(64)
	return d.FieldStruct(name, func(d *decode.D) {
		start := d.Pos()
		br := d.RawLen(size * 8)
		end := d.Pos()
		compressionType := d.FieldU8("compression", compressionTypes, scalar.UintHex)
		// validate crc
		data := d.ReadAllBits(br)
		bytesToCheck := append(data, uint8(compressionType))
		maskedCRCInt := maskedCrc32(bytesToCheck)
		d.FieldU32("crc", d.UintAssert(uint64(maskedCRCInt)), scalar.UintHex)
		d.FieldStruct("data", func(d *decode.D) {
			// https://github.com/google/leveldb/blob/main/table/block_builder.cc#L16
			// https://github.com/google/leveldb/blob/main/table/block.cc
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
			// Right now, its omited if empty.
			if restartOffset <= 0 {
				return
			}
			d.SeekAbs(start)
			d.FieldArray("entries", func(d *decode.D) {
				for d.Pos() < start+restartOffset {
					d.FieldStruct("entry", func(d *decode.D) {
						d.FieldUintFn("shared_bytes", decodeVarInt)
						unshared := int64(d.FieldUintFn("unshared_bytes", decodeVarInt))
						valueLength := d.FieldUintFn("value_length", decodeVarInt)
						// InternalKey
						// https://github.com/google/leveldb/blob/main/db/dbformat.h#L171
						d.FieldStruct("key_delta", func(d *decode.D) {
							d.FieldUTF8("user_key", int(unshared-uint64Size/8))
							d.FieldU8("type", valueTypes, scalar.UintHex)
							d.FieldU56("sequence_number")
						})
						if valueCallbackFn == nil {
							d.FieldUTF8("value", int(valueLength))
						} else {
							d.FieldStruct("value", valueCallbackFn)
						}
					})
				}
			})
		})
	})
}

func decodeVarInt(d *decode.D) uint64 {
	var value uint64 = 0
	var shift uint64 = 0

	for {
		b := d.U8()
		value |= (b & 0b01111111) << shift
		shift += 7
		if b&0b10000000 == 0 {
			break
		}
	}

	return value
}

// Return a masked representation of the crc.
// https://github.com/google/leveldb/blob/main/util/crc32c.h#L29
func mask(crc uint32) uint32 {
	const kMaskDelta = 0xa282ead8
	// Rotate right by 15 bits and add a constant.
	return ((crc >> 15) | (crc << 17)) + kMaskDelta
}

func maskedCrc32(bytes []uint8) uint32 {
	crc32C := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	crc32C.Write(bytes)
	return mask(crc32C.Sum32())
}

// Print the hexadecimal representation in little-endian format.
func printLE(name string, value uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	fmt.Printf("%s: % x\n", name, buf)
}
