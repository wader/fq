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
	"fmt"
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
	uint32BitSize    = 32
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
const (
	valueTypeDeletion = 0x0
	valueTypeValue    = 0x1
)

var valueTypes = scalar.UintMapSymStr{
	valueTypeDeletion: "deletion",
	valueTypeValue:    "value",
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
	readTableBlock(
		"metaindex",
		metaIndexSize,
		keyValueContentsReader(
			nil,
			func(d *decode.D) {
				handle := readBlockHandle(d)
				metaHandles = append(metaHandles, handle)
			},
		),
		d,
	)

	// index

	d.SeekAbs(indexOffset * 8)
	var dataHandles []blockHandle
	readTableBlock(
		"index",
		indexSize,
		keyValueContentsReader(
			readInternalKey,
			func(d *decode.D) {
				handle := readBlockHandle(d)
				dataHandles = append(dataHandles, handle)
			},
		),
		d,
	)

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
// The function `readTableBlockContents` gets the _uncompressed_ bytebuffer.
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
// https://github.com/google/leveldb/blob/main/table/block.cc#L48
func readKeyValueContents(
	keyCallbackFn func(sharedBytes []byte, unsharedSize int, d *decode.D) error,
	valueCallbackFn func(d *decode.D),
	size int64,
	d *decode.D,
) {
	start := d.Pos()
	end := start + size*8

	var restartOffset int64
	d.SeekAbs(end - uint32BitSize)
	d.FieldStruct("trailer", func(d *decode.D) {
		numRestarts := int64(d.FieldU32("num_restarts"))
		restartOffset = size*8 - (1+numRestarts)*uint32BitSize
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
		var lastKey []byte
		for d.Pos() < start+restartOffset {
			d.FieldStruct("entry", func(d *decode.D) {
				// https://github.com/google/leveldb/blob/main/table/block.cc#L48-L75
				shared := int64(d.FieldULEB128("shared_bytes"))
				unshared := int64(d.FieldULEB128("unshared_bytes"))
				valueLength := int64(d.FieldULEB128("value_length"))

				// read key
				// https://github.com/google/leveldb/blob/main/table/block.cc#L261
				if int(shared) > len(lastKey) {
					d.Fatalf("`shared` size is larger than `lastKey` length")
				}
				keyPrefix := lastKey[:shared]
				keySuffix := readBytesWithoutChangingPosition(unshared, d)
				lastKey = append(keyPrefix, keySuffix...)

				if keyCallbackFn == nil && shared == 0 {
					d.FieldUTF8("key", int(unshared))
				} else {
					d.FieldStruct("key", func(d *decode.D) {
						if keyCallbackFn == nil {
							keyCallbackFn = readPrefixedBytes
						}
						err := keyCallbackFn(keyPrefix, int(unshared), d)
						if err != nil {
							d.Errorf("%v", err)
						}
					})
				}

				// read value
				if valueCallbackFn == nil {
					d.FieldUTF8("value", int(valueLength))
				} else {
					d.FieldStruct("value", valueCallbackFn)
				}
			})
		}
	})
}

func readBytesWithoutChangingPosition(nBytes int64, d *decode.D) []byte {
	var result []byte
	d.RangeFn(d.Pos(), nBytes*8, func(d *decode.D) {
		br := d.RawLen(d.BitsLeft())
		result = d.ReadAllBits(br)
	})
	return result
}

// https://github.com/google/leveldb/blob/main/table/format.cc#L24
func readBlockHandle(d *decode.D) blockHandle {
	return blockHandle{
		offset: d.FieldULEB128("offset"),
		size:   d.FieldULEB128("size"),
	}
}

// Read bytes and prefix with given bytes;
// name read bytes "suffix" and the merged bytes "full" (synthetic field).
func readPrefixedBytes(prefixBytes []byte, nBytes int, d *decode.D) error {
	br, err := d.TryFieldRawLen("suffix", int64(nBytes)*8)
	if err != nil {
		return err
	}
	full := append(prefixBytes, d.ReadAllBits(br)...)
	d.FieldValueStr("full", string(full), strInferred)
	return nil
}

// An "internal key" consists of the triple (user_key, type, sequence_number).
// https://github.com/google/leveldb/blob/main/db/dbformat.h#L171
func readInternalKey(sharedBytes []byte, unsharedSize int, d *decode.D) error {
	// In the LevelDB encoding, the internal key can be cut at any byte:
	// including the user_key, type, or sequence_number:
	// https://github.com/google/leveldb/blob/main/table/block_builder.cc#L79-L83
	//
	// The resulting prefix is then shared among subsequent keys and not
	// specified explicitly by them. Here, we handle each cutoff case.
	//
	// All sizes are in bytes unless mentioned otherwise.
	keySize := len(sharedBytes) + unsharedSize
	typeAndSequenceNumberSize := 8
	//                        key
	// +-----------------------------------------------+
	//               user_key
	// +---------------------------------+
	//               ⁞   user_key_suffix   type sequence_number
	// [AAAAAAAAAAAA]⁞[BBBBBBBBBBBBBBBBBB] [T] [SSSSSSS]
	//               ⁞                      1   7 bytes
	// +------------+⁞+--------------------------------+
	//     shared    ⁞             unshared
	//               ⁞
	//             cutoff
	if keySize < typeAndSequenceNumberSize || int64(unsharedSize) > d.BitsLeft()/8 {
		return fmt.Errorf("key size %d or unshared size %d invalid", keySize, unsharedSize)
	}

	// case 1: user_key, type, and sequence_number fit fully in unshared.
	if len(sharedBytes) == 0 {
		d.FieldUTF8("user_key", keySize-typeAndSequenceNumberSize)
		d.FieldU8("type", valueTypes, scalar.UintHex)
		d.FieldU56("sequence_number")
		return nil
	}

	// case 2: type and sequence_number fit fully in unshared: simulate user_key value.
	if unsharedSize >= typeAndSequenceNumberSize {
		suffix := fieldUTF8ReturnBytes("user_key_suffix", unsharedSize-typeAndSequenceNumberSize, d)
		d.FieldValueStr("user_key", stringify(sharedBytes, suffix), strInferred)
		d.FieldU8("type", valueTypes, scalar.UintHex)
		d.FieldU56("sequence_number")
		return nil
	}

	// case 3: sequence_number fits fully in unshared: simulate user_key and type value,
	sequenceNumberSize := typeAndSequenceNumberSize - 1
	if unsharedSize == sequenceNumberSize {
		lastIndex := len(sharedBytes) - 1
		d.FieldValueStr("user_key", string(sharedBytes[:lastIndex]), strInferred)
		d.FieldValueUint("type", uint64(sharedBytes[lastIndex]), valueTypes, scalar.UintHex, uintInferred)
		d.FieldU56("sequence_number")
		return nil
	}

	// case 4: sequence_number cut: simulate user_key, type, and sequence_number value.
	typeByteIndex := keySize - typeAndSequenceNumberSize
	d.FieldValueStr("user_key", string(sharedBytes[:typeByteIndex]), strInferred)
	d.FieldValueUint("type", uint64(sharedBytes[typeByteIndex]), valueTypes, scalar.UintHex, uintInferred)
	var suffixBytes []byte
	if unsharedSize > 0 {
		br := d.FieldRawLen("sequence_number_suffix", int64(unsharedSize)*8)
		suffixBytes = d.ReadAllBits(br)
	}
	sequenceNumberBytes := append(
		sharedBytes[typeByteIndex+1:keySize-unsharedSize],
		suffixBytes...,
	)
	sequenceNumberBE := bitio.Read64(sequenceNumberBytes[:], 0, int64(sequenceNumberSize*8))
	sequenceNumberLE := bitio.ReverseBytes64(56, sequenceNumberBE)
	d.FieldValueUint("sequence_number", sequenceNumberLE, uintInferred)

	return nil
}

// Read content encoded in the "filter" or "stats" Meta Block format.
// https://github.com/google/leveldb/blob/main/doc/table_format.md#filter-meta-block
// https://github.com/google/leveldb/blob/main/table/filter_block.cc
func readMetaContent(nBytes int64, d *decode.D) {
	// TK(2023-12-04)
	d.FieldRawLen("raw", nBytes*8)
}

// Helpers

var strInferred = scalar.StrFn(func(s scalar.Str) (scalar.Str, error) {
	s.Description = "inferred"
	return s, nil
})

var uintInferred = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	s.Description = "inferred"
	return s, nil
})

func keyValueContentsReader(
	keyCallbackFn func(sharedPrefix []byte, unsharedSize int, d *decode.D) error,
	valueCallbackFn func(d *decode.D),
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

// Concatinate byteslices and convert into a string.
func stringify(byteSlices ...[]byte) string {
	totalSize := 0
	for _, b := range byteSlices {
		totalSize += len(b)
	}

	result := make([]byte, 0, totalSize)

	for _, b := range byteSlices {
		result = append(result, b...)
	}

	return string(result)
}

func fieldUTF8ReturnBytes(name string, nBytes int, d *decode.D) []byte {
	var result []byte
	d.RangeFn(d.Pos(), int64(nBytes)*8, func(d *decode.D) {
		result = d.BytesLen(nBytes)
	})
	d.FieldUTF8(name, nBytes)
	return result
}
