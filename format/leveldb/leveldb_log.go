package leveldb

// https://github.com/google/leveldb/blob/main/doc/impl.md#log-files
// https://github.com/google/leveldb/blob/main/db/write_batch.cc
//
// Files in LevelDB using this format include:
//  - *.log

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed leveldb_log.md
var leveldbLogFS embed.FS

func init() {
	interp.RegisterFormat(
		format.LevelDB_LOG,
		&decode.Format{
			Description: "LevelDB Log",
			DecodeFn:    ldbLogDecode,
		})
	interp.RegisterFS(leveldbLogFS)
}

func ldbLogDecode(d *decode.D) any {
	rro := recordReadOptions{readDataFn: func(size int64, recordType int, d *decode.D) {
		if recordType == recordTypeFull {
			d.FieldStruct("data", func(d *decode.D) {
				d.LimitedFn(size, readBatch)
			})
		} else {
			d.FieldRawLen("data", size)
		}
	}}
	readBlockSequence(rro, d)

	return nil
}

// https://github.com/google/leveldb/blob/main/db/write_batch.cc#L5-L14
//
// WriteBatch::rep_ :=
//
//	sequence: fixed64
//	count: fixed32
//	data: record[count]
//
// record :=
//
//	kTypeValue varstring varstring
//	kTypeDeletion varstring
//
// varstring :=
//
//	len: varint32
//	data: uint8[len]
func readBatch(d *decode.D) {
	d.FieldU64("sequence")
	expectedCount := d.FieldU32("count")
	actualCount := uint64(0)
	d.FieldArray("records", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("record", func(d *decode.D) {
				valueType := d.FieldULEB128("type", valueTypes)
				switch valueType {
				case valueTypeDeletion:
					readLengthPrefixedString("key", d)
				case valueTypeValue:
					readLengthPrefixedString("key", d)
					readLengthPrefixedString("value", d)
				default:
					d.Fatalf("unknown value type: %d", valueType)
				}
			})
			actualCount++
		}
	})
	if actualCount != expectedCount {
		d.Errorf("actual record count (%d) does not equal expected count (%d)", actualCount, expectedCount)
	}
}
