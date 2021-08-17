package protobuf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.PROTOBUF,
		Description: "Protobuf",
		DecodeFn:    protobufDecode,
	})
}

const (
	wireTypeVarint          = 0
	wireType64Bit           = 1
	wireTypeLengthDelimited = 2
	wireType32Bit           = 5
)

var wireTypeNames = map[uint64]string{
	0: "Varint",
	1: "64-bit",
	2: "Length-delimited",
	5: "32-bit",
}

func varInt(d *decode.D) uint64 {
	var n uint64
	for i := 0; ; i++ {
		b := d.U8()
		n = n | (b&0x7f)<<(7*i)
		if b&0x80 == 0 {
			break
		}
	}

	return n
}

func fieldVarInt(d *decode.D, name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return varInt(d), decode.NumberDecimal, ""
	})
}

func protobufDecodeField(d *decode.D, pbm *format.ProtoBufMessage) {
	d.FieldStructFn("field", func(d *decode.D) {
		keyN := fieldVarInt(d, "key_n")
		fieldNumber := keyN >> 3
		wireType := keyN & 0x7
		d.FieldValueU("field_number", fieldNumber, "")
		d.FieldValueU("wire_type", wireType, wireTypeNames[wireType])

		var value uint64
		var length uint64
		var valueStart int64
		switch wireType {
		case wireTypeVarint:
			value = fieldVarInt(d, "wire_value")
		case wireType64Bit:
			value = d.FieldU64("wire_value")
		case wireTypeLengthDelimited:
			length = fieldVarInt(d, "length")
			valueStart = d.Pos()
			d.FieldBitBufLen("wire_value", int64(length)*8)
		case wireType32Bit:
			value = d.FieldU32("wire_value")
		}

		if pbm != nil {
			if pbf, ok := (*pbm)[int(fieldNumber)]; ok {
				d.FieldValueStr("name", pbf.Name, "")
				d.FieldValueStr("type", format.ProtoBufTypeNames[uint64(pbf.Type)], "")

				switch pbf.Type {
				case format.ProtoBufTypeInt32, format.ProtoBufTypeInt64:
					v := num.ZigZag(value)
					d.FieldValueS("value", v, "")
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[uint64(v)], "")
					}
				case format.ProtoBufTypeUInt32, format.ProtoBufTypeUInt64:
					d.FieldValueU("value", value, "")
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[value], "")
					}
				case format.ProtoBufTypeSInt32, format.ProtoBufTypeSInt64:
					// TODO: correct? 32 different?
					v := num.TwosComplement(64, value)
					d.FieldValueS("value", v, "")
					if len(pbf.Enums) > 0 {
						d.FieldValueStr("enum", pbf.Enums[uint64(v)], "")
					}
				case format.ProtoBufTypeBool:
					d.FieldValueBool("value", value != 0, "")
				case format.ProtoBufTypeEnum:
					d.FieldValueStr("enum", pbf.Enums[value], "")
				case format.ProtoBufTypeFixed64:
					// TODO:
				case format.ProtoBufTypeSFixed64:
					// TODO:
				case format.ProtoBufTypeDouble:
					// TODO:
				case format.ProtoBufTypeString:
					d.FieldValueStr("value", string(d.BytesRange(valueStart, int(length))), "")
				case format.ProtoBufTypeBytes:
					d.FieldValueBytes("value", d.BytesRange(valueStart, int(length)), "")
				case format.ProtoBufTypeMessage:
					// TODO: test
					d.DecodeLenFn(int64(length)*8, func(d *decode.D) {
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
	d.FieldArrayFn("fields", func(d *decode.D) {
		for d.BitsLeft() > 0 {
			protobufDecodeField(d, pbm)
			//break
		}
	})
}

func protobufDecode(d *decode.D, in interface{}) interface{} {
	var pbm *format.ProtoBufMessage
	pbi, ok := in.(format.ProtoBufIn)
	if ok {
		pbm = &pbi.Message
	}

	protobufDecodeFields(d, pbm)

	return nil
}
