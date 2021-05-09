package ebml

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

type Enum struct {
	Value      string
	Label      string
	Definition string
}

type Attribute struct {
	Name          string
	Type          Type
	Tag           Tag
	Definition    string
	IntegerEnums  map[int64]Enum
	UintegerEnums map[uint64]Enum
	StringEnums   map[string]Enum
}

type Tag map[uint64]Attribute

const (
	CRC32              = 0xbf
	Void               = 0xec
	EBMLVersion        = 0x4286
	EBMLReadVersion    = 0x42f7
	EBMLMaxIDLength    = 0x42f2
	EBMLMaxSizeLength  = 0x42f3
	DocType            = 0x4282
	DocTypeVersion     = 0x4287
	DocTypeReadVersion = 0x4285
)

var Global = Tag{
	CRC32: {Name: "CRC-32", Type: Binary},
	Void:  {Name: "Void", Type: Binary},
}

var Header = Tag{
	EBMLVersion:        {Name: "EBMLVersion", Type: Uinteger},
	EBMLReadVersion:    {Name: "EBMLReadVersion", Type: Uinteger},
	EBMLMaxIDLength:    {Name: "EBMLMaxIDLength", Type: Uinteger},
	EBMLMaxSizeLength:  {Name: "EBMLMaxSizeLength", Type: Uinteger},
	DocType:            {Name: "DocType", Type: String},
	DocTypeVersion:     {Name: "DocTypeVersion", Type: Uinteger},
	DocTypeReadVersion: {Name: "DocTypeReadVersion", Type: Uinteger},
}
