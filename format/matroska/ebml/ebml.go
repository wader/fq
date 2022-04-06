package ebml

import "github.com/wader/fq/pkg/scalar"

type Type int

const (
	Integer Type = iota
	Uinteger
	Float
	String
	UTF8
	Date
	Binary
	Master
)

var TypeNames = map[Type]string{
	Integer:  "integer",
	Uinteger: "uinteger",
	Float:    "float",
	String:   "string",
	UTF8:     "UTF8",
	Date:     "data",
	Binary:   "binary",
	Master:   "master",
}

type Attribute struct {
	Name          string
	Type          Type
	Tag           Tag
	Definition    string
	IntegerEnums  scalar.SToScalar
	UintegerEnums scalar.UToScalar
	StringEnums   scalar.StrToScalar
}

type Tag map[uint64]Attribute

const (
	CRC32ID              = 0xbf
	VoidID               = 0xec
	HeaderID             = 0x1a45dfa3
	EBMLVersionID        = 0x4286
	EBMLReadVersionID    = 0x42f7
	EBMLMaxIDLengthID    = 0x42f2
	EBMLMaxSizeLengthID  = 0x42f3
	DocTypeID            = 0x4282
	DocTypeVersionID     = 0x4287
	DocTypeReadVersionID = 0x4285
)

var Global = Tag{
	CRC32ID: {Name: "crc32", Type: Binary},
	VoidID:  {Name: "void", Type: Binary},
}

var Header = Tag{
	EBMLVersionID:        {Name: "ebml_version", Type: Uinteger},
	EBMLReadVersionID:    {Name: "ebml_read_version", Type: Uinteger},
	EBMLMaxIDLengthID:    {Name: "ebml_max_id_length", Type: Uinteger},
	EBMLMaxSizeLengthID:  {Name: "ebml_max_size_length", Type: Uinteger},
	DocTypeID:            {Name: "doc_type", Type: String},
	DocTypeVersionID:     {Name: "doc_type_version", Type: Uinteger},
	DocTypeReadVersionID: {Name: "doc_type_read_version", Type: Uinteger},
}
