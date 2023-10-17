package bplist

import (
	"embed"
	"math"
	"math/big"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/apple"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bplist.jq bplist.md ns_keyed_archiver.jq
var bplistFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Bplist,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinUnique,
			Description: "Apple Binary Property List",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    bplistDecode,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(bplistFS)
}

const (
	elementTypeNullOrBoolOrFill = 0x00
	elementTypeInt              = 0x01
	elementTypeReal             = 0x02
	elementTypeDate             = 0x03
	elementTypeData             = 0x04
	elementTypeASCIIString      = 0x05
	elementTypeUnicodeString    = 0x06
	elementTypeUID              = 0x08
	elementTypeArray            = 0x0a
	elementTypeSet              = 0x0c
	elementTypeDict             = 0x0d
)

const (
	null      = 0x00
	boolFalse = 0x08
	boolTrue  = 0x09
)

var elementTypeMap = scalar.UintMap{
	elementTypeNullOrBoolOrFill: {Sym: "singleton", Description: "Singleton value (null/bool)"},
	elementTypeInt:              {Sym: "int", Description: "Integer"},
	elementTypeReal:             {Sym: "real", Description: "Floating Point Number"},
	elementTypeDate:             {Sym: "date", Description: "Date, 4 or 8 byte float"},
	elementTypeData:             {Sym: "data", Description: "Binary data"},
	elementTypeASCIIString:      {Sym: "ascii_string", Description: "ASCII encoded string"},
	elementTypeUnicodeString:    {Sym: "unicode_string", Description: "Unicode string"},
	elementTypeUID:              {Sym: "uid", Description: "UID"},
	elementTypeArray:            {Sym: "array", Description: "Array"},
	elementTypeSet:              {Sym: "set", Description: "Set"},
	elementTypeDict:             {Sym: "dict", Description: "Dictionary"},
}

var cocoaTimeEpochDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

// decodes the number of bits required to store the following object
func decodeSize(d *decode.D, sms ...scalar.UintMapper) uint64 {
	n := d.FieldU4("size_bits")
	if n != 0x0f {
		return n
	}

	d.FieldU4("large_size_marker", d.UintAssert(0b0001))

	// get the exponent value
	n = d.FieldU4("exponent")

	// calculate the number of bytes encoding the size
	n = 1 << n

	// decode that many bytes as big endian
	n = d.FieldUintFn(
		"size_bigint",
		func(d *decode.D) uint64 {
			v := d.UBigInt(int(n * 8))
			d.BigIntAssertRange(big.NewInt(1), big.NewInt(math.MaxInt64))
			return v.Uint64()
		}, sms...)

	return n
}

// decodeItem decodes an object from the plist, and assumes that the current
// seek position of the *decode.D is an object type tag. Returns a bool
// indicating whether or not a string was decoded, which is necssary for
// checking dictionary key type validity.
func decodeItem(d *decode.D, p *plist) bool {
	m := d.FieldU4("type", elementTypeMap)
	switch m {
	case elementTypeNullOrBoolOrFill:
		d.FieldU4("value", scalar.UintMap{
			null:      scalar.Uint{Sym: nil},
			boolTrue:  scalar.Uint{Sym: true},
			boolFalse: scalar.Uint{Sym: false},
		})
	case elementTypeInt, elementTypeUID:
		n := d.FieldUintFn("size", func(d *decode.D) uint64 {
			return 1 << d.U4()
		})
		switch n {
		case 1:
			d.FieldU8("value")
		case 2:
			d.FieldU16("value")
		case 4:
			d.FieldU32("value")
		case 8:
			d.FieldS64("value")
		case 16:
			d.FieldSBigInt("value", int(n*8))
		default:
			d.Errorf("invalid integer size %d", n)
		}
	case elementTypeReal:
		n := 1 << decodeSize(d)
		d.FieldValueUint("size", uint64(n))
		d.FieldF("value", n*8)
	case elementTypeDate:
		n := 1 << decodeSize(d, d.UintAssert(4, 8))
		d.FieldValueUint("size", uint64(n))
		d.FieldF("value", n*8, scalar.FltActualDateDescription(cocoaTimeEpochDate, time.Second, time.RFC3339))
	case elementTypeData:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldRawLen("value", int64(n*8))
	case elementTypeASCIIString:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldUTF8("value", int(n))
		return true
	case elementTypeUnicodeString:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldUTF16BE("value", int(n))
		return true
	case elementTypeArray:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				idx := d.FieldU("object_index", int(p.t.objRefSize)*8)
				p.decodeReference(d, idx)
			})
	case elementTypeSet:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				idx := d.FieldU("object_index", int(p.t.objRefSize)*8)
				p.decodeReference(d, idx)
			})
	case elementTypeDict:
		n := decodeSize(d)
		d.FieldValueUint("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				var ki, vi uint64
				ki = d.FieldU("key_index", int(p.t.objRefSize)*8)
				d.SeekRel(int64((n-1)*p.t.objRefSize)*8, func(d *decode.D) {
					vi = d.FieldU("value_index", int(p.t.objRefSize)*8)
				})
				d.FieldStruct("key", func(d *decode.D) {
					if k := p.decodeReference(d, ki); !k {
						d.Errorf("non-string key in dictionary")
					}
				})
				d.FieldStruct("value", func(d *decode.D) {
					p.decodeReference(d, vi)
				})
			})
	default:
		d.Errorf("unknown type marker: %d", m)
	}

	return false
}

// decodeReference looks up and decodes an object based on its index in the
// offset table. Returns a bool indicating whether or not the decoded item is
// a string (necessary for checking dictionary key validity).
func (pl *plist) decodeReference(d *decode.D, idx uint64) bool {
	if idx > uint64(len(pl.o)) {
		// prevent a panic
		d.Errorf("index %d out of bounds for object table size %d", idx, len(pl.o))
		return false
	}
	pl.consumed[idx] = true

	defer pl.pld.PushAndPop(idx, func() { d.Fatalf("infinite recursion detected") })()

	itemOffset := pl.o[idx]
	if itemOffset >= pl.t.offsetTableStart {
		d.Errorf("attempting to decode object %d at offset 0x%x beyond offset table start 0x%x",
			idx, itemOffset, pl.t.offsetTableStart)
	}

	var isString bool
	d.SeekAbs(int64(itemOffset*8), func(d *decode.D) {
		isString = decodeItem(d, pl)
	})
	return isString
}

type trailer struct {
	offTblOffSize    uint64
	objRefSize       uint64
	nObjects         uint64
	topObjectOffset  uint64
	offsetTableStart uint64
}

type plist struct {
	t        trailer
	o        []uint64
	consumed map[uint64]bool
	pld      apple.PosLoopDetector[uint64]
}

func bplistDecode(d *decode.D) any {
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 6, d.StrAssert("bplist"))
		d.FieldUTF8("version", 2, d.StrAssert("00"))
	})

	p := new(plist)
	p.consumed = make(map[uint64]bool)

	d.SeekAbs(d.Len()-32*8, func(d *decode.D) {
		d.FieldStruct("trailer", func(d *decode.D) {
			d.FieldU40("unused")
			d.FieldS8("sort_version")
			p.t.offTblOffSize = d.FieldU8("offset_table_offset_size", d.UintAssertRange(1, 8))
			p.t.objRefSize = d.FieldU8("object_reference_size", d.UintAssertRange(1, 8))
			p.t.nObjects = d.FieldU64("object_count")
			p.t.topObjectOffset = d.FieldU64("top_object_offset")
			p.t.offsetTableStart = d.FieldU64("offset_table_start")
		})
	})

	d.SeekAbs(int64(p.t.offsetTableStart*8), func(d *decode.D) {
		i := uint64(0)
		d.FieldArrayLoop("offset_table",
			func() bool { return i < p.t.nObjects },
			func(d *decode.D) {
				off := d.FieldU("element", 8*int(p.t.offTblOffSize))
				p.o = append(p.o, off)
				i++
			},
		)
	})

	d.FieldStruct("objects",
		func(d *decode.D) {
			p.decodeReference(d, 0)
		})

	var lost []uint64

	for i := uint64(0); i < p.t.nObjects; i++ {
		if _, isUsed := p.consumed[i]; !isUsed {
			lost = append(lost, i)
		}
	}

	if len(lost) == 0 {
		return nil
	}

	i := 0

	d.FieldStructNArray("lost_and_found", "entry",
		int64(len(lost)),
		func(d *decode.D) {
			p.decodeReference(d, lost[i])
			i++
		})

	return nil
}
