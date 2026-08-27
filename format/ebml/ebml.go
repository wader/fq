package ebml

import (
	"time"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// 2001-01-01T00:00:00.000000000 UTC
var EpochDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

type ID int

type Element interface {
	GetType() string
	GetID() ID
	GetParentID() ID
	GetName() string
	GetDefinition() string
	GetMinOccurs() uint64
	GetMaxOccurs() uint64
	GetRange() string
	GetLength() string
	GetDefault() string
	GetUnknownSizeAllowed() bool
	GetRecursive() bool
	GetRecurring() bool
	GetMinVer() uint64
	GetMaxVer() uint64
}

type Enum struct {
	Name        string
	Description string
}

type ElementType struct {
	ID                 ID
	ParentID           ID
	Name               string
	Definition         string
	MinOccurs          uint64
	MaxOccurs          uint64
	Range              string
	Length             string
	Default            string
	UnknownSizeAllowed bool
	Recursive          bool
	Recurring          bool
	MinVer             uint64
	MaxVer             uint64
}

func (e *ElementType) GetType() string             { return "" }
func (e *ElementType) GetID() ID                   { return e.ID }
func (e *ElementType) GetParentID() ID             { return e.ParentID }
func (e *ElementType) GetName() string             { return e.Name }
func (e *ElementType) GetDefinition() string       { return e.Definition }
func (e *ElementType) GetMinOccurs() uint64        { return e.MinOccurs }
func (e *ElementType) GetMaxOccurs() uint64        { return e.MaxOccurs }
func (e *ElementType) GetRange() string            { return e.Range }
func (e *ElementType) GetLength() string           { return e.Length }
func (e *ElementType) GetDefault() string          { return e.Default }
func (e *ElementType) GetUnknownSizeAllowed() bool { return e.UnknownSizeAllowed }
func (e *ElementType) GetRecursive() bool          { return e.Recursive }
func (e *ElementType) GetRecurring() bool          { return e.Recurring }
func (e *ElementType) GetMinVer() uint64           { return e.MinVer }
func (e *ElementType) GetMaxVer() uint64           { return e.MaxVer }

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

// Simplifies the resolution of id to an Element by looking into standard EBML elements
// if id is not found in schema-specific types.
func (e *Master) LookupByID(id ID) (Element, bool) {
	if match, ok := e.Master[id]; ok {
		return match, true
	}
	if match, ok := Global.Master[id]; ok {
		return match, true
	}
	return nil, false
}

// Peek the next element and sets `out`
func (e *Master) PeekNextElement(d *decode.D, out *Element) uint64 {
	n := PeekRawVint(d)
	var ok bool
	*out, ok = e.LookupByID(ID(n))
	if !ok {
		*out = &Unknown{}
	}
	return n
}

const (
	RootID = 0

	CRC32ID = 0xbf
	VoidID  = 0xec
)

var Global = &Master{
	ElementType: ElementType{
		ID:        -1,
		ParentID:  -1,
		Name:      "",
		MinOccurs: 1,
		MaxOccurs: 1,
	},
	Master: map[ID]Element{
		CRC32ID: &Binary{ElementType: ElementType{Name: "crc32", Length: "4", MinOccurs: 0, MaxOccurs: 1}},
		VoidID:  &Binary{ElementType: ElementType{Name: "void"}},
	},
}

const (
	HeaderID                  = 0x1a45dfa3
	EBMLVersionID             = 0x4286
	EBMLReadVersionID         = 0x42f7
	EBMLMaxIDLengthID         = 0x42f2
	EBMLMaxSizeLengthID       = 0x42f3
	DocTypeID                 = 0x4282
	DocTypeVersionID          = 0x4287
	DocTypeReadVersionID      = 0x4285
	DocTypeExtensionID        = 0x4281
	DocTypeExtensionNameID    = 0x4283
	DocTypeExtensionVersionID = 0x4284
)

var Header = &Master{
	ElementType: ElementType{
		ID:        HeaderID,
		Name:      "ebml",
		MinOccurs: 1,
		MaxOccurs: 1,
	},
	Master: map[ID]Element{
		EBMLVersionID:        &Uinteger{ElementType: ElementType{Name: "ebml_version", Definition: "EBML Version", MinOccurs: 1, MaxOccurs: 1, Range: ">0", Default: "1"}},
		EBMLReadVersionID:    &Uinteger{ElementType: ElementType{Name: "ebml_read_version", Definition: "Minimum EBML reader version", MinOccurs: 1, MaxOccurs: 1, Range: "1", Default: "1"}},
		EBMLMaxIDLengthID:    &Uinteger{ElementType: ElementType{Name: "ebml_max_id_length", Definition: "Maximum id length", MinOccurs: 1, MaxOccurs: 1, Range: ">=4", Default: "4"}},
		EBMLMaxSizeLengthID:  &Uinteger{ElementType: ElementType{Name: "ebml_max_size_length", Definition: "Maximum length of encoded size", MinOccurs: 1, MaxOccurs: 1, Range: ">0", Default: "8"}},
		DocTypeID:            &String{ElementType: ElementType{Name: "doc_type", Definition: "Document content type", MinOccurs: 1, MaxOccurs: 1}},
		DocTypeVersionID:     &Uinteger{ElementType: ElementType{Name: "doc_type_version", Definition: "Document type version", MinOccurs: 1, MaxOccurs: 1, Range: ">0", Default: "1"}},
		DocTypeReadVersionID: &Uinteger{ElementType: ElementType{Name: "doc_type_read_version", Definition: "Minimum document reader version", MinOccurs: 1, MaxOccurs: 1, Range: ">0", Default: "1"}},
		DocTypeExtensionID: &Master{
			ElementType: ElementType{
				ID:   DocTypeExtensionID,
				Name: "doc_type_extension",
			},
			Master: map[ID]Element{
				DocTypeExtensionNameID:    &String{ElementType: ElementType{Name: "doc_type_extension_name", Definition: "Extensions of the main doctype", MinOccurs: 1, MaxOccurs: 1}},
				DocTypeExtensionVersionID: &Uinteger{ElementType: ElementType{Name: "doc_type_extension_version", Definition: "Extension version", MinOccurs: 1, MaxOccurs: 1, Range: ">0"}},
			},
		},
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

// TODO: smarter?
func DecodeRawVintWidth(d *decode.D) (uint64, int) {
	n := d.U8()
	w := 1
	for i := 0; i <= 7 && (n&(1<<(7-i))) == 0; i++ {
		w++
	}
	for i := 1; i < w; i++ {
		n = n<<8 | d.U8()
	}
	return n, w
}

func DecodeRawVint(d *decode.D) uint64 {
	n, _ := DecodeRawVintWidth(d)
	return n
}

func PeekRawVint(d *decode.D) uint64 {
	n, w := DecodeRawVintWidth(d)
	d.SeekRel(int64(-w) * 8)
	return n
}

func DecodeVint(d *decode.D) uint64 {
	n, w := DecodeRawVintWidth(d)
	m := (uint64(1<<((w-1)*8+(8-w))) - 1)
	return n & m
}

// returns a scalar.UintMapper for EBML tags' display
func ElementIDMapper(e Element) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		s.DisplayFormat = scalar.NumberHex
		s.Sym = e.GetName()
		s.Description = e.GetDefinition()
		return s, nil
	})
}
