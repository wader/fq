package bplist

import (
	"embed"
	"math"
	"math/big"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bplist.jq bplist.md
var bplistFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.BPLIST,
		ProbeOrder:  format.ProbeOrderBinUnique,
		Description: "Apple Binary Property List",
		Groups:      []string{format.PROBE},
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

var elementTypeMap = scalar.UToScalar{
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

// decodes the number of bits required to store the following object
func decodeSize(d *decode.D, sms ...scalar.Mapper) uint64 {
	n := d.FieldU4("size_bits")
	if n != 0x0f {
		return n
	}

	d.FieldU4("large_size_marker", d.AssertU(0b0001))

	// get the exponent value
	n = d.FieldU4("exponent")

	// calculate the number of bytes encoding the size
	n = 1 << n

	// decode that many bytes as big endian
	n = d.FieldUFn(
		"size_bigint",
		func(d *decode.D) uint64 {
			v := d.UBigInt(int(n * 8))
			d.AssertBigIntRange(big.NewInt(1), big.NewInt(math.MaxInt64))
			return v.Uint64()
		}, sms...)

	return n
}

func decodeItem(d *decode.D, p *plist) {
	m := d.FieldU4("type", elementTypeMap)
	switch m {
	case elementTypeNullOrBoolOrFill:
		t := d.U4()
		switch t {
		case null:
			d.FieldValueNil("value")
		case boolTrue:
			d.FieldValueBool("value", true)
		case boolFalse:
			d.FieldValueBool("value", false)
		}
	case elementTypeInt:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		if n*8 <= 64 {
			d.FieldU("value", int(n*8))
		} else {
			d.FieldUBigInt("value", int(n))
		}
	case elementTypeReal:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldF("value", int(n))
	case elementTypeDate:
		n := 1 << decodeSize(d, d.AssertU(4, 8))
		d.FieldValueU("size", uint64(n))
		d.FieldF("value", n*8, scalar.DescriptionActualFCocoaDate)
	case elementTypeData:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldRawLen("value", int64(n*8))
	case elementTypeASCIIString:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldUTF8("value", int(n))
	case elementTypeUnicodeString:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldUTF16("value", int(n))
	case elementTypeUID:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldUBigInt("value", int(n)).Uint64()
	case elementTypeArray:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				idx := d.FieldU8("object_index")
				decodeReference(d, p, idx)
			})
	case elementTypeSet:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				idx := d.FieldU8("object_index")
				decodeReference(d, p, idx)
			})
	case elementTypeDict:
		n := decodeSize(d)
		d.FieldValueU("size", n)
		d.FieldStructNArray("entries", "entry", int64(n),
			func(d *decode.D) {
				var ki, vi uint64
				ki = d.FieldU8("key_index")
				d.SeekRel(int64((n-1)*p.t.objRefSize)*8, func(d *decode.D) {
					vi = d.FieldU8("value_index")
				})
				d.FieldStruct("key", func(d *decode.D) {
					decodeReference(d, p, ki)
				})
				d.FieldStruct("value", func(d *decode.D) {
					decodeReference(d, p, vi)
				})
			})
	default:
		d.Errorf("unknown type marker: %d", m)
	}
}

func decodeReference(d *decode.D, p *plist, idx uint64) {
	if idx > uint64(len(p.o)) {
		// prevent a panic
		d.Errorf("index %d out of bounds for object table size %d", idx, len(p.o))
		return
	}

	d.SeekAbs(int64(p.o[idx]*8), func(d *decode.D) {
		decodeItem(d, p)
	})
}

type trailer struct {
	offTblOffSize    uint64
	objRefSize       uint64
	nObjects         uint64
	topObjectOffset  uint64
	offsetTableStart uint64
}

type plist struct {
	t trailer
	o []uint64
}

func bplistDecode(d *decode.D, _ any) any {
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 6, d.AssertStr("bplist"))
		d.FieldUTF8("version", 2, d.AssertStr("00"))
	})

	p := new(plist)

	d.SeekAbs(d.Len()-32*8, func(d *decode.D) {
		d.FieldStruct("trailer", func(d *decode.D) {
			d.FieldU40("unused")
			d.FieldS8("sort_version")
			p.t.offTblOffSize = d.FieldU8("offset_table_offset_size", d.AssertURange(1, 8))
			p.t.objRefSize = d.FieldU8("object_reference_size", d.AssertURange(1, 8))
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
			decodeReference(d, p, 0)
		})

	return nil
}
