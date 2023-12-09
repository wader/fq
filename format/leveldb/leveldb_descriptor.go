package leveldb

// https://github.com/google/leveldb/blob/main/doc/impl.md#manifest
// https://github.com/google/leveldb/blob/main/db/version_edit.cc
//
// Files in LevelDB using this format include:
//  - MANIFEST-*

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed leveldb_descriptor.jq
//go:embed leveldb_descriptor.md
var leveldbDescriptorFS embed.FS

func init() {
	interp.RegisterFormat(
		format.LevelDB_Descriptor,
		&decode.Format{
			Description: "LevelDB Descriptor",
			DecodeFn:    ldbDescriptorDecode,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(leveldbDescriptorFS)
}

const (
	tagTypeComparator     = 1
	tagTypeLogNumber      = 2
	tagTypeNextFileNumber = 3
	tagTypeLastSequence   = 4
	tagTypeCompactPointer = 5
	tagTypeDeletedFile    = 6
	tagTypeNewFile        = 7
	// 8 not used anymore
	tagTypePrevLogNumber = 9
)

var tagTypes = scalar.UintMapSymStr{
	tagTypeComparator:     "comparator",
	tagTypeLogNumber:      "log_number",
	tagTypeNextFileNumber: "next_file_number",
	tagTypeLastSequence:   "last_sequence",
	tagTypeCompactPointer: "compact_pointer",
	tagTypeDeletedFile:    "deleted_file",
	tagTypeNewFile:        "new_file",
	tagTypePrevLogNumber:  "previous_log_number",
}

func ldbDescriptorDecode(d *decode.D) any {
	rro := recordReadOptions{readDataFn: func(size int64, recordType int, d *decode.D) {
		if recordType == recordTypeFull {
			d.FieldStruct("data", func(d *decode.D) {
				d.LimitedFn(size, readManifest)
			})
		} else {
			d.FieldRawLen("data", size)
		}
	}}
	readBlockSequence(rro, d)

	return nil
}

// List of sorted tables for each level involving key ranges and other metadata.
func readManifest(d *decode.D) {
	d.FieldArray("tags", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("tag", func(d *decode.D) {
				tag := d.FieldULEB128("key", tagTypes)
				switch tag {
				case tagTypeComparator:
					readLengthPrefixedString("value", d)
				case tagTypeLogNumber,
					tagTypePrevLogNumber,
					tagTypeNextFileNumber,
					tagTypeLastSequence:
					d.FieldULEB128("value")
				case tagTypeCompactPointer:
					d.FieldStruct("value", func(d *decode.D) {
						d.FieldULEB128("level")
						readTagInternalKey("internal_key", d)
					})
				case tagTypeDeletedFile:
					d.FieldStruct("value", func(d *decode.D) {
						d.FieldULEB128("level")
						d.FieldULEB128("file_number")
					})
				case tagTypeNewFile:
					d.FieldStruct("value", func(d *decode.D) {
						d.FieldULEB128("level")
						d.FieldULEB128("file_number")
						d.FieldULEB128("file_size")
						readTagInternalKey("smallest_internal_key", d)
						readTagInternalKey("largest_internal_key", d)
					})
				default:
					d.Fatalf("unknown tag: %d", tag)
				}
			})
		}
	})
}

func readTagInternalKey(name string, d *decode.D) {
	d.FieldStruct(name, func(d *decode.D) {
		length := d.FieldULEB128("length")
		d.FieldStruct("data", func(d *decode.D) {
			err := readInternalKey(nil, int(length), d)
			if err != nil {
				d.Errorf("%v", err)
			}
		})
	})
}
