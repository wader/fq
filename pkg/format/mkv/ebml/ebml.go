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

var Global = Tag{
	0xbf: {Name: "CRC-32", Type: Binary},
	0xec: {Name: "Void", Type: Binary},
}

var Header = Tag{
	0x4286: {Name: "EBMLVersion", Type: Uinteger},
	0x42f7: {Name: "EBMLReadVersion", Type: Uinteger},
	0x42f2: {Name: "EBMLMaxIDLength", Type: Uinteger},
	0x42f3: {Name: "EBMLMaxSizeLength", Type: Uinteger},
	0x4282: {Name: "DocType", Type: String},
	0x4287: {Name: "DocTypeVersion", Type: Uinteger},
	0x4285: {Name: "DocTypeReadVersion", Type: Uinteger},
}
