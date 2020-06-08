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
	CRC32ID: {Name: "CRC-32", Type: Binary},
	VoidID:  {Name: "Void", Type: Binary},
}

var Header = Tag{
	EBMLVersionID:        {Name: "EBMLVersion", Type: Uinteger},
	EBMLReadVersionID:    {Name: "EBMLReadVersion", Type: Uinteger},
	EBMLMaxIDLengthID:    {Name: "EBMLMaxIDLength", Type: Uinteger},
	EBMLMaxSizeLengthID:  {Name: "EBMLMaxSizeLength", Type: Uinteger},
	DocTypeID:            {Name: "DocType", Type: String},
	DocTypeVersionID:     {Name: "DocTypeVersion", Type: Uinteger},
	DocTypeReadVersionID: {Name: "DocTypeReadVersion", Type: Uinteger},
}
