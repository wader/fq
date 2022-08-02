package fit

// TODO: chained files
// TODO: developer message?
// TODO: filed number mapping, xls file?
// TODO: Compressed Timestamp Header

// https://developer.garmin.com/fit/protocol/
import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.FIT,
		&decode.Format{
			Description: "Flexible and Interoperable Data Transfer Protocol",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeFIT,
		})
}

const (
	architectureTypeLittleEndian = 0
	architectureTypeBigEndian    = 1
)

var architectureTypeMap = scalar.UintMapSymStr{
	architectureTypeLittleEndian: "little_endian",
	architectureTypeBigEndian:    "big_endian",
}

const (
	messageTypeData       = 0
	messageTypeDefinition = 1
)

var messageTypeMap = scalar.UintMapSymStr{
	messageTypeData:       "data",
	messageTypeDefinition: "definition",
}

// Base Type #	Endian Ability	Base Type Field	Type Name	Invalid Value	Size (Bytes)	Comment
// 0	0	0x00	enum	0xFF	1
// 1	0	0x01	sint8	0x7F	1	2’s complement format
// 2	0	0x02	uint8	0xFF	1
// 3	1	0x83	sint16	0x7FFF	2	2’s complement format
// 4	1	0x84	uint16	0xFFFF	2
// 5	1	0x85	sint32	0x7FFFFFFF	4	2’s complement format
// 6	1	0x86	uint32	0xFFFFFFFF	4
// 7	0	0x07	string	0x00	1	Null terminated string encoded in UTF-8 format
// 8	1	0x88	float32	0xFFFFFFFF	4
// 9	1	0x89	float64	0xFFFFFFFFFFFFFFFF	8
// 10	0	0x0A	uint8z	0x00	1
// 11	1	0x8B	uint16z	0x0000	2
// 12	1	0x8C	uint32z	0x00000000	4
// 13	0	0x0D	byte	0xFF	1	Array of bytes. Field is invalid if all bytes are invalid.
// 14	1	0x8E	sint64	0x7FFFFFFFFFFFFFFF	8	2’s complement format
// 15	1	0x8F	uint64	0xFFFFFFFFFFFFFFFF	8
// 16	1	0x90	uint64z	0x0000000000000000	8

const (
	baseTypeEnum    = 0
	baseTypeSint8   = 1
	baseTypeUint8   = 2
	baseTypeSint16  = 3
	baseTypeUint16  = 4
	baseTypeSint32  = 5
	baseTypeUint32  = 6
	baseTypeString  = 7
	baseTypeFloat32 = 8
	baseTypeFloat64 = 9
	baseTypeUint8z  = 10
	baseTypeUint16z = 11
	baseTypeUint32z = 12
	baseTypeByte    = 13
	baseTypeSint64  = 14
	baseTypeUint64  = 15
	baseTypeUint64z = 16
)

var baseTypeMap = scalar.UintMapSymStr{
	baseTypeEnum:    "enum",
	baseTypeSint8:   "sint8",
	baseTypeUint8:   "uint8",
	baseTypeSint16:  "sint16",
	baseTypeUint16:  "uint16",
	baseTypeSint32:  "sint32",
	baseTypeUint32:  "uint32",
	baseTypeString:  "string",
	baseTypeFloat32: "float32",
	baseTypeFloat64: "float64",
	baseTypeUint8z:  "uint8z",
	baseTypeUint16z: "uint16z",
	baseTypeUint32z: "uint32z",
	baseTypeByte:    "byte",
	baseTypeSint64:  "sint64",
	baseTypeUint64:  "uint64",
	baseTypeUint64z: "uint64z",
}

var baseTypeSize = map[int]int{
	baseTypeEnum:    1,
	baseTypeSint8:   1,
	baseTypeUint8:   1,
	baseTypeSint16:  2,
	baseTypeUint16:  2,
	baseTypeSint32:  4,
	baseTypeUint32:  4,
	baseTypeString:  1,
	baseTypeFloat32: 4,
	baseTypeFloat64: 8,
	baseTypeUint8z:  1,
	baseTypeUint16z: 2,
	baseTypeUint32z: 4,
	baseTypeByte:    1,
	baseTypeSint64:  8,
	baseTypeUint64:  8,
	baseTypeUint64z: 8,
}

type field struct {
	number        uint8
	size          uint8
	endianAbility uint8
	baseType      uint8
}

type developerField struct {
	number         uint8
	size           uint8
	developerIndex uint8
}

type definition struct {
	s                   scalar.Uint
	architecture        int
	globalMessageNumber int
	fields              []field
	developerFields     []developerField
}

type definitionEntries map[uint64]definition

func (fes definitionEntries) MapUint(s scalar.Uint) (scalar.Uint, error) {
	u := s.Actual
	if fe, ok := fes[u]; ok {
		s = fe.s
		s.Actual = u
	}
	return s, nil
}

func decodeBaseType(d *decode.D, f field) {
	switch f.baseType {
	case baseTypeEnum:
		d.FieldU8("value")
	case baseTypeSint8:
		d.FieldS8("value")
	case baseTypeUint8:
		d.FieldU8("value")
	case baseTypeSint16:
		d.FieldS16("value")
	case baseTypeUint16:
		d.FieldU16("value")
	case baseTypeSint32:
		d.FieldU32("value")
	case baseTypeUint32:
		d.FieldU32("value")
	case baseTypeString:
		d.FieldUTF8NullFixedLen("value", int(f.size))
	case baseTypeFloat32:
		d.FieldF32("value")
	case baseTypeFloat64:
		d.FieldF64("value")
	case baseTypeUint8z:
		d.FieldU8("value")
	case baseTypeUint16z:
		d.FieldU16("value")
	case baseTypeUint32z:
		d.FieldU32("value")
	case baseTypeByte:
		d.FieldRawLen("value", int64(f.size)*8)
	case baseTypeSint64:
		d.FieldS64("value")
	case baseTypeUint64:
		d.FieldU64("value")
	case baseTypeUint64z:
		d.FieldU64("value")
	default:
		d.Fatalf("unknown base type %d", f.baseType)
	}
}

func decodeDataMessage(d *decode.D, de definition) {
	d.FieldArray("fields", func(d *decode.D) {
		for _, f := range de.fields {
			baseSize, ok := baseTypeSize[int(f.baseType)]
			if !ok {
				d.Fatalf("unknown base size for base type %d", f.baseType)
			}
			values := int(f.size) / baseSize

			switch {
			case values == 1,
				f.baseType == baseTypeString,
				f.baseType == baseTypeByte:
				decodeBaseType(d, f)
			default:
				d.FieldArray("values", func(d *decode.D) {
					for i := 0; i < values; i++ {
						decodeBaseType(d, f)
					}
				})
			}
		}
	})
	if len(de.developerFields) > 0 {
		d.FieldArray("developer_fields", func(d *decode.D) {
			for _, f := range de.developerFields {
				d.FieldRawLen("filed", int64(f.size)*8)
			}
		})
	}
}

func decodeDefinitionMessage(d *decode.D, messageTypeSpecific uint64) definition {
	var de definition
	d.FieldU8("reserved")
	de.architecture = int(d.FieldU8("architecture", architectureTypeMap))
	de.globalMessageNumber = int(d.FieldU16("global_message_number"))
	numFields := d.FieldU8("fields")
	d.FieldArray("field_definitions", func(d *decode.D) {
		for i := uint64(0); i < numFields; i++ {
			d.FieldStruct("field_definition", func(d *decode.D) {
				var f field
				f.number = uint8(d.FieldU8("field_definition_number"))
				f.size = uint8(d.FieldU8("size"))
				f.endianAbility = uint8(d.FieldU1("endian_ability"))
				d.FieldRawLen("reserved", 2)
				f.baseType = uint8(d.FieldU5("base_type_number", baseTypeMap))

				de.fields = append(de.fields, f)
			})
		}
	})
	if messageTypeSpecific == 1 {
		developerFields := d.FieldU8("developer_fields")
		d.FieldArray("developer_field_definitions", func(d *decode.D) {
			for i := uint64(0); i < developerFields; i++ {
				d.FieldStruct("developer_field_definition", func(d *decode.D) {
					var f developerField
					f.number = uint8(d.FieldU8("field_number"))
					f.size = uint8(d.FieldU8("size"))
					f.developerIndex = uint8(d.FieldU8("developer_data_index"))

					de.developerFields = append(de.developerFields, f)
				})
			}
		})
	}

	return de
}

func decodeFIT(d *decode.D) any {
	d.Endian = decode.LittleEndian

	definitions := definitionEntries{
		0: definition{
			s: scalar.Uint{Sym: "file_id"},
			fields: []field{
				{number: 0, size: 1, baseType: baseTypeEnum},
				{number: 1, size: 2, baseType: baseTypeUint16},
				{number: 2, size: 2, baseType: baseTypeUint16},
				{number: 3, size: 4, baseType: baseTypeUint32z},
				{number: 4, size: 4, baseType: baseTypeUint32},
				{number: 5, size: 2, baseType: baseTypeUint16},
				// {number: 5, baseType: baseTypeString},

			},
		},
	}
	var dataSize uint64

	d.FieldStruct("header", func(d *decode.D) {
		size := d.FieldU8("size")
		if size < 12 {
			d.Fatalf("Header size too small %d < 12", size)
		}
		d.FieldU8("protocol_version")
		d.FieldU16("profile_version")
		dataSize = d.FieldU32("data_size")
		d.FieldUTF8("data_type", 4, d.StrAssert(".FIT"))
		d.FieldU16("crc", scalar.UintHex)
	})

	d.FramedFn(int64(dataSize)*8, func(d *decode.D) {
		d.FieldArray("records", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("record_header", func(d *decode.D) {
					d.FieldU1("normal_header")
					messageType := d.FieldU1("message_type", messageTypeMap)
					messageTypeSpecific := d.FieldU1("message_type_specific")
					d.FieldU1("reserved")
					localMessageType := d.FieldU4("local_message_type", definitions)
					d.FieldStruct("message", func(d *decode.D) {
						switch messageType {
						case messageTypeData:
							if de, ok := definitions[localMessageType]; ok {
								decodeDataMessage(d, de)
							} else {
								d.Fatalf("unknown local message type %d", localMessageType)
							}
						case messageTypeDefinition:
							definitions[localMessageType] = decodeDefinitionMessage(d, messageTypeSpecific)
						default:
							panic("unreachable")
						}
					})
				})
			}
		})
	})
	d.FieldU16("crc", scalar.UintHex)

	return nil
}
