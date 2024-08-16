package leveldb

// https://github.com/google/leveldb/blob/main/doc/log_format.md
//
// Files in LevelDB using the "log-format" of block sequences include:
//  - *.log
//  - MANIFEST-*

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type recordReadOptions struct {
	// Both .log- and MANIFEST-files use the Log-format,
	// i.e., a sequence of records split into 32KB blocks.
	// However, the format of the data within the records differ.
	// This function specifies how to read said data.
	readDataFn func(size int64, recordType int, d *decode.D)
}

// https://github.com/google/leveldb/blob/main/db/log_format.h
const (
	// checksum (4 bytes) + length (2 bytes) + record type (1 byte)
	headerSize = (4 + 2 + 1) * 8

	blockSize = (32 * 1024) * 8 // 32KB

	recordTypeZero   = 0 // preallocated file regions
	recordTypeFull   = 1
	recordTypeFirst  = 2 // fragments
	recordTypeMiddle = 3
	recordTypeLast   = 4
)

var recordTypes = scalar.UintMapSymStr{
	recordTypeZero:   "zero",
	recordTypeFull:   "full",
	recordTypeFirst:  "first",
	recordTypeMiddle: "middle",
	recordTypeLast:   "last",
}

// Read a sequence of 32KB-blocks (the last one may be less).
// https://github.com/google/leveldb/blob/main/db/log_reader.cc#L189
func readBlockSequence(rro recordReadOptions, d *decode.D) {
	d.Endian = decode.LittleEndian

	d.FieldArray("blocks", func(d *decode.D) {
		for d.BitsLeft() >= headerSize {
			d.LimitedFn(min(blockSize, d.BitsLeft()), func(d *decode.D) {
				d.FieldStruct("block", bind(readLogBlock, rro))
			})
		}
	})

	if d.BitsLeft() > 0 {
		// The reference implementation says:
		// "[...] if buffer_ is non-empty, we have a truncated header at the
		// end of the file, which can be caused by the writer crashing in the
		// middle of writing the header. Instead of considering this an error,
		// just report EOF."
		d.FieldRawLen("truncated_block", d.BitsLeft())
	}
}

// Read a Log-block, consisting of up to 32KB of records and an optional trailer.
//
// block := record* trailer?
func readLogBlock(rro recordReadOptions, d *decode.D) {
	if d.BitsLeft() > blockSize {
		d.Fatalf("Bits left greater than maximum log-block size of 32KB.")
	}
	// record*
	d.FieldArray("records", func(d *decode.D) {
		for d.BitsLeft() >= headerSize {
			d.FieldStruct("record", bind(readLogRecord, rro))
		}
	})
	// trailer?
	if d.BitsLeft() > 0 {
		d.FieldRawLen("trailer", d.BitsLeft())
	}
}

// Read a Log-record.
//
// checksum: uint32     // crc32c of type and data[] ; little-endian
// length: uint16       // little-endian
// type: uint8          // One of FULL, FIRST, MIDDLE, LAST
// data: uint8[length]
//
// via https://github.com/google/leveldb/blob/main/doc/log_format.md
func readLogRecord(rro recordReadOptions, d *decode.D) {
	// header
	var checksumValue *decode.Value
	var length int64
	var recordType int
	d.LimitedFn(headerSize, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) {
			d.FieldU32("checksum", scalar.UintHex)
			checksumValue = d.FieldGet("checksum")
			length = int64(d.FieldU16("length"))
			recordType = int(d.FieldU8("record_type", recordTypes))
		})
	})

	// verify checksum: record type (1 byte) + data (`length` bytes)
	d.RangeFn(d.Pos()-8, (1+length)*8, func(d *decode.D) {
		bytesToCheck := d.Bits(int(d.BitsLeft()))
		actualChecksum := computeChecksum(bytesToCheck)
		_ = checksumValue.TryUintScalarFn(d.UintAssert(uint64(actualChecksum)))
	})

	// data
	dataSize := length * 8
	rro.readDataFn(dataSize, recordType, d)
}

func readLengthPrefixedString(name string, d *decode.D) {
	d.FieldStruct(name, func(d *decode.D) {
		length := d.FieldULEB128("length")
		d.FieldUTF8("data", int(length))
	})
}

// simplified `functools.partial` (Python) or `Function.prototype.bind` (JavaScript)
func bind(f func(recordReadOptions, *decode.D), rro recordReadOptions) func(*decode.D) {
	return func(d *decode.D) {
		f(rro, d)
	}
}
