package protobuf

// https://developers.google.com/protocol-buffers/docs/encoding

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed protobuf.md
var protobufFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Protobuf,
		&decode.Format{
			Description: "Protobuf",
			DecodeFn:    protobufDecode,
		})
	interp.RegisterFS(protobufFS)
}

const (
	wireTypeVarint          = 0
	wireType64Bit           = 1
	wireTypeLengthDelimited = 2
	wireType32Bit           = 5
)

var wireTypeNames = scalar.UintMapSymStr{
	0: "varint",
	1: "64bit",
	2: "length_delimited",
	5: "32bit",
}

func protobufDecodeField(d *decode.D, pbm *format.ProtoBufMessage) {
	d.FieldStruct("field", func(d *decode.D) {
		keyN := d.FieldULEB128("key_n")
		fieldNumber := keyN >> 3
		wireType := keyN & 0x7
		d.FieldValueUint("field_number", fieldNumber)
		d.FieldValueUint("wire_type", wireType, scalar.UintSym(wireTypeNames[wireType]))

		var value uint64
		var length uint64
		var valuePos int64
		switch wireType {
		case wireTypeVarint:
			value = d.FieldULEB128("wire_value")
		case wireType64Bit:
			value = d.FieldU64("wire_value")
		case wireTypeLengthDelimited:
			length = d.FieldULEB128("length")
			valuePos = d.Pos()
			d.FieldRawLen("wire_value", int64(length)*8)
		case wireType32Bit:
			value = d.FieldU32("wire_value")
		}

		if pbm != nil {
			if pbf, ok := (*pbm)[int(fieldNumber)]; ok {
				d.FieldValueStr("name", pbf.Name)
				d.FieldValueStr("type", format.ProtoBufTypeNames[uint64(pbf.Type)])

				switch pbf.Type {
				case format.ProtoBufTypeInt32, format.ProtoBufTypeInt64:
					v := mathx.ZigZag[uint64, int64](value)
					d.FieldValueSint("value", v)
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[uint64(v)])
					}
				case format.ProtoBufTypeUInt32, format.ProtoBufTypeUInt64:
					d.FieldValueUint("value", value)
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[value])
					}
				case format.ProtoBufTypeSInt32, format.ProtoBufTypeSInt64:
					// TODO: correct? 32 different?
					v := mathx.TwosComplement(64, value)
					d.FieldValueSint("value", v)
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[uint64(v)])
					}
				case format.ProtoBufTypeBool:
					d.FieldValueBool("value", value != 0)
				case format.ProtoBufTypeEnum:
					d.FieldValueStr("enum", pbf.Enums[value])
				case format.ProtoBufTypeFixed64:
					// TODO:
				case format.ProtoBufTypeSFixed64:
					// TODO:
				case format.ProtoBufTypeDouble:
					// TODO:
				case format.ProtoBufTypeString:
					d.SeekAbs(valuePos)
					d.FieldUTF8("value", int(length))
				case format.ProtoBufTypeBytes:
					d.SeekAbs(valuePos)
					d.FieldRawLen("value", int64(length)*8)
				case format.ProtoBufTypeMessage:
					// TODO: test
					d.FramedFn(int64(length)*8, func(d *decode.D) {
						protobufDecodeFields(d, &pbf.Message)
					})
				case format.ProtoBufTypePackedRepeated:
					// TODO:
				case format.ProtoBufTypeFixed32:
					// TODO:
				case format.ProtoBufTypeSFixed32:
					// TODO:
				case format.ProtoBufTypeFloat:
					// TODO:
				}
			}
		}
	})
}

func protobufDecodeFields(d *decode.D, pbm *format.ProtoBufMessage) {
	d.FieldArray("fields", func(d *decode.D) {
		for d.BitsLeft() > 0 {
			protobufDecodeField(d, pbm)
		}
	})
}

func protobufDecode(d *decode.D) any {
	var pbi format.Protobuf_In
	d.ArgAs(&pbi)

	protobufDecodeFields(d, &pbi.Message)

	return nil
}
