package ebml

import "time"

// 2001-01-01T00:00:00.000000000 UTC
var EpochDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

type ID int

type Element interface {
	GetType() string
	GetID() ID
	GetParentID() ID
	GetName() string
	GetDefinition() string
}

type Enum struct {
	Name        string
	Description string
}

type ElementType struct {
	ID         ID
	ParentID   ID
	Name       string
	Definition string
}

func (e *ElementType) GetType() string       { return "" }
func (e *ElementType) GetID() ID             { return e.ID }
func (e *ElementType) GetParentID() ID       { return e.ParentID }
func (e *ElementType) GetName() string       { return e.Name }
func (e *ElementType) GetDefinition() string { return e.Definition }

type ElementScalarType[T comparable] struct {
	ElementType
	Enums map[T]Enum
}

func (e *ElementScalarType[T]) GetEnum() map[T]Enum { return e.Enums }

type Unknown struct{ ElementType }

func (*Unknown) GetType() string { return "unknown" }

type Integer ElementScalarType[int64]

func (*Integer) GetType() string { return "integer" }

type Uinteger ElementScalarType[uint64]

func (*Uinteger) GetType() string { return "uinteger" }

type Float ElementScalarType[float64]

func (*Float) GetType() string { return "float" }

type String ElementScalarType[string]

func (*String) GetType() string { return "string" }

type UTF8 ElementScalarType[string]

func (*UTF8) GetType() string { return "utf8" }

type Date struct{ ElementType }

func (*Date) GetType() string { return "date" }

type Binary struct{ ElementType }

func (*Binary) GetType() string { return "binary" }

type Master struct {
	ElementType
	Master map[ID]Element
}

func (e *Master) GetType() string           { return "master" }
func (e *Master) GetMaster() map[ID]Element { return e.Master }

const (
	RootID = 0

	CRC32ID = 0xbf
	VoidID  = 0xec
)

var Global = &Master{
	ElementType: ElementType{
		ID:       -1,
		ParentID: -1,
		Name:     "",
	},
	Master: map[ID]Element{
		CRC32ID: &Binary{ElementType: ElementType{Name: "crc32"}},
		VoidID:  &Binary{ElementType: ElementType{Name: "void"}},
	},
}

const (
	HeaderID             = 0x1a45dfa3
	EBMLVersionID        = 0x4286
	EBMLReadVersionID    = 0x42f7
	EBMLMaxIDLengthID    = 0x42f2
	EBMLMaxSizeLengthID  = 0x42f3
	DocTypeID            = 0x4282
	DocTypeVersionID     = 0x4287
	DocTypeReadVersionID = 0x4285
)

var Header = &Master{
	ElementType: ElementType{
		ID:       HeaderID,
		ParentID: RootID,
		Name:     "ebml",
	},
	Master: map[ID]Element{
		EBMLVersionID:        &Uinteger{ElementType: ElementType{Name: "ebml_version", Definition: "EBML Version"}},
		EBMLReadVersionID:    &Uinteger{ElementType: ElementType{Name: "ebml_read_version", Definition: "Minimum EBML reader version"}},
		EBMLMaxIDLengthID:    &Uinteger{ElementType: ElementType{Name: "ebml_max_id_length", Definition: "Maximum id length"}},
		EBMLMaxSizeLengthID:  &Uinteger{ElementType: ElementType{Name: "ebml_max_size_length", Definition: "Maximum body length"}},
		DocTypeID:            &String{ElementType: ElementType{Name: "doc_type", Definition: "Document content type"}},
		DocTypeVersionID:     &Uinteger{ElementType: ElementType{Name: "doc_type_version", Definition: "Document type version"}},
		DocTypeReadVersionID: &Uinteger{ElementType: ElementType{Name: "doc_type_read_version", Definition: "Minimum document reader version"}},
	},
}

// FindParentID find id walking parents of startID
func FindParentID(idToElement map[ID]Element, startID ID, id ID) (Element, bool) {
	current := idToElement[startID]
	for {
		if current.GetID() == id {
			return current, true
		}
		var ok bool
		current, ok = idToElement[current.GetParentID()]
		if !ok {
			break
		}
	}
	return nil, false
}
