package avro

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var jsonGroup decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AVRO_OCF,
		Description: "Avro object container file",
		Groups:      []string{format.PROBE},
		DecodeFn:    avroDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.JSON}, Group: &jsonGroup},
		},
	})
}

const headerMetadataSchema = `{"type": "map", "values": "bytes"}`
const intMask = byte(127)
const intFlag = byte(128)

// readLong reads a variable length zig zag long from the current position in decoder
// and returns the decoded value and the number of bytes read.
func varZigZag(d *decode.D) int64 {
	var value uint64
	var shift uint
	for d.NotEnd() {
		b := byte(d.U8())
		value |= uint64(b&intMask) << shift
		if b&intFlag == 0 {
			return int64(value>>1) ^ -int64(value&1)
		}
		shift += 7
	}
	panic("unexpected end of data")
}

func avroDecode(d *decode.D, in interface{}) interface{} {
	d.FieldRawLen("magic", 4*8, d.AssertBitBuf([]byte{'O', 'b', 'j', 1}))
	//var schema []byte
	var blockCount int64 = -1
	d.FieldStructArrayLoop("meta", "block",
		func() bool { return blockCount != 0 },
		func(d *decode.D) {
			blockCount = d.FieldSFn("count", varZigZag)
			// If its negative, then theres another long representing byte size
			if blockCount < 0 {
				blockCount *= -1
				d.FieldSFn("size", varZigZag)
			}
			if blockCount == 0 {
				return
			}

			var i int64 = 0
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < blockCount }, func(d *decode.D) {
				keyL := d.FieldSFn("key_length", varZigZag)
				key := d.FieldUTF8("key", int(keyL))
				valL := d.FieldSFn("value_length", varZigZag)
				if key == "avro.schema" {
					d.FieldFormatLen("value", valL*8, jsonGroup, nil)
				} else {
					d.FieldUTF8("value", int(valL))
				}
				i++
			})
		})
	syncbb := d.FieldRawLen("sync", 16*8)
	sync, err := syncbb.BytesLen(16)
	if err != nil {
		d.Fatalf("unable to read sync bytes: %v", err)
	}
	d.FieldStructArrayLoop("blocks", "block", func() bool { return d.NotEnd() }, func(d *decode.D) {
		count := d.FieldSFn("count", varZigZag)
		_ = count
		size := d.FieldSFn("size", varZigZag)
		d.FieldRawLen("data", size*8)
		d.FieldRawLen("sync", 16*8, d.AssertBitBuf(sync))
	})

	return nil
}
