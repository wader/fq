package format

import "github.com/wader/fq/pkg/scalar"

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

var ProtoBufTypeNames = scalar.UintMapSymStr{
	ProtoBufTypeInt32:          "int32",
	ProtoBufTypeInt64:          "int64",
	ProtoBufTypeUInt32:         "uint32",
	ProtoBufTypeUInt64:         "uint64",
	ProtoBufTypeSInt32:         "sint32",
	ProtoBufTypeSInt64:         "sint64",
	ProtoBufTypeBool:           "bool",
	ProtoBufTypeEnum:           "enum",
	ProtoBufTypeFixed64:        "fixed64",
	ProtoBufTypeSFixed64:       "sfixed64",
	ProtoBufTypeDouble:         "double",
	ProtoBufTypeString:         "string",
	ProtoBufTypeBytes:          "bytes",
	ProtoBufTypeMessage:        "message",
	ProtoBufTypePackedRepeated: "packed_repeated",
	ProtoBufTypeFixed32:        "fixed32",
	ProtoBufTypeSFixed32:       "sfixed32",
	ProtoBufTypeFloat:          "float",
}

type ProtoBufField struct {
	Type    int
	Name    string
	Message ProtoBufMessage
	Enums   map[uint64]string
}

type ProtoBufMessage map[int]ProtoBufField
