package format

import "github.com/wader/fq/pkg/decode"

type ProtoBufType int

const (
	ProtoBufTypeInt32 = iota
	ProtoBufTypeInt64
	ProtoBufTypeUInt32
	ProtoBufTypeUInt64
	ProtoBufTypeSInt32
	ProtoBufTypeSInt64
	ProtoBufTypeBool
	ProtoBufTypeEnum
	ProtoBufTypeFixed64
	ProtoBufTypeSFixed64
	ProtoBufTypeDouble
	ProtoBufTypeString
	ProtoBufTypeBytes
	ProtoBufTypeMessage
	ProtoBufTypePackedRepeated
	ProtoBufTypeFixed32
	ProtoBufTypeSFixed32
	ProtoBufTypeFloat
)

var ProtoBufTypeNames = decode.UToStr{
	ProtoBufTypeInt32:          "Int32",
	ProtoBufTypeInt64:          "Int64",
	ProtoBufTypeUInt32:         "UInt32",
	ProtoBufTypeUInt64:         "UInt64",
	ProtoBufTypeSInt32:         "SInt32",
	ProtoBufTypeSInt64:         "SInt64",
	ProtoBufTypeBool:           "Bool",
	ProtoBufTypeEnum:           "Enum",
	ProtoBufTypeFixed64:        "Fixed64",
	ProtoBufTypeSFixed64:       "SFixed64",
	ProtoBufTypeDouble:         "Double",
	ProtoBufTypeString:         "String",
	ProtoBufTypeBytes:          "Bytes",
	ProtoBufTypeMessage:        "Message",
	ProtoBufTypePackedRepeated: "PackedRepeated",
	ProtoBufTypeFixed32:        "Fixed32",
	ProtoBufTypeSFixed32:       "SFixed32",
	ProtoBufTypeFloat:          "Float",
}

type ProtoBufField struct {
	Type    int
	Name    string
	Message ProtoBufMessage
	Enums   map[uint64]string
}

type ProtoBufMessage map[int]ProtoBufField
