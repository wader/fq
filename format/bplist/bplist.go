package bplist

import (
	"math"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.BPLIST,
		ProbeOrder:  format.ProbeOrderBinFuzzy, // after most others (overlap some with webp)
		Description: "Apple Binary Property List",
		Groups:      []string{format.PROBE},
		DecodeFn:    bplistDecode,
	})
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
	elementTypeArray            = 0x0A
	elementTypeSet              = 0x0C
	elementTypeDict             = 0x0D
)

const (
	null      = 0x00
	boolFalse = 0x08
	boolTrue  = 0x09
)

var singleByteElementMap = scalar.UToScalar{
	null:      {Sym: "null", Description: "Null Value"},
	boolFalse: {Sym: "false", Description: "False Value"},
	boolTrue:  {Sym: "true", Description: "True Value"},
}

var elementTypeMap = scalar.UToScalar{
	elementTypeNullOrBoolOrFill: {Sym: "singleton", Description: "Singleton value (null/bool)"},
	elementTypeInt:              {Sym: "int", Description: "Integer"},
	elementTypeReal:             {Sym: "real", Description: "Floating Point Number"},
	elementTypeDate:             {Sym: "date", Description: "Date, 8-byte float"},
	elementTypeData:             {Sym: "data", Description: "Binary data"},
	elementTypeASCIIString:      {Sym: "ascii_string", Description: "ASCII encoded string"},
	elementTypeUnicodeString:    {Sym: "unicode_string", Description: "Unicode string"},
	elementTypeUID:              {Sym: "uid", Description: "UID"},
	elementTypeArray:            {Sym: "array", Description: "Array"},
	elementTypeSet:              {Sym: "set", Description: "Set"},
	elementTypeDict:             {Sym: "dict", Description: "Dictionary"},
}

// decodes the number of bits required to store the following object
func decodeSize(d *decode.D) uint64 {
	n := d.FieldU4("size_bits")
	if n != 0x0F {
		return uint64(n)
	}

	d.FieldU4("large_size_marker", d.AssertU(1)) // TODO: add assertion that this is 0001

	// get the exponent value
	n = d.FieldU4("exponent")

	// calculate the number of bytes encoding the size
	n = uint64(math.Pow(2, float64(n)))

	// decode that many bytes as big endian
	n = d.FieldUBigIntBE("size", int(n*8)).Uint64()
	return n
}

func decodeItem(d *decode.D, p *plist) any {
	d.FieldStruct("object", func(d *decode.D) {
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
			d.FieldSFn("value", func(d *decode.D) int64 {
				return d.UBigIntBE(int(n)).Int64()
			})
		case elementTypeReal:
			n := decodeSize(d)
			d.FieldFFn("value", func(d *decode.D) float64 {
				return d.FE(int(n), decode.BigEndian)
			})
		case elementTypeDate:
			d.FieldStrFn("value", func(d *decode.D) string {
				v := d.F64()
				t := time.Unix(int64(v), 0)
				return t.String()
			})
		case elementTypeData:
			d.FieldStruct("value", func(d *decode.D) {
				n := decodeSize(d)
				d.FieldRawLen("value", int64(n))
			})
		case elementTypeASCIIString:
			d.FieldStruct("value", func(d *decode.D) {
				d.FieldStrFn("value", func(d *decode.D) string {
					n := decodeSize(d)
					return d.UTF8(int(n))
				})
			})
		case elementTypeUnicodeString:
			d.FieldStrFn("value", func(d *decode.D) string {
				n := decodeSize(d)
				return d.UTF16(int(n))
			})
		case elementTypeUID:
			n := decodeSize(d)
			d.FieldUFn("value", func(d *decode.D) uint64 {
				return d.UBigIntBE(int(n)).Uint64()
			})
		case elementTypeArray:
			n := decodeSize(d)
			i := uint64(0)
			d.FieldStructArrayLoop("elements", "item",
				func() bool { return i < n },
				func(d *decode.D) {
					idx := d.FieldU8("object_index")
					decodeReference(d, p, idx)
					i++
				})
		case elementTypeSet:
			n := decodeSize(d)
			i := uint64(0)
			d.FieldArrayLoop("entries",
				func() bool { return i < n },
				func(d *decode.D) {
					idx := d.FieldU8("object_index")
					decodeReference(d, p, idx)
					i++
				})
		case elementTypeDict:
			d.FieldStruct("dictionary", func(d *decode.D) {
				s := decodeSize(d)
				i := uint64(0)
				d.FieldStructArrayLoop("entries", "entry",
					func() bool { return i < s },
					func(d *decode.D) {
						var ki, vi uint64
						ki = d.FieldU8("key_index")
						d.SeekRel(int64((s-1)*uint64(p.T.objRefSize)*8), func(d *decode.D) {
							vi = d.FieldU8("value_index")
						})
						i++
						d.FieldStruct("key", func(d *decode.D) {
							decodeReference(d, p, ki)
						})
						d.FieldStruct("value", func(d *decode.D) {
							decodeReference(d, p, vi)
						})
					})
			})
		default:
			d.Errorf("unknown type marker: %d", m)
		}
	})

	return nil
}

func decodeReference(d *decode.D, p *plist, idx uint64) {
	d.SeekAbs(int64(p.O[idx]*8), func(d *decode.D) {
		decodeItem(d, p)
	})
}

type trailer struct {
	offTblOffSize    int64
	objRefSize       int64
	nObjects         int64
	topObjectOffset  int64
	offsetTableStart int64
}

type plist struct {
	T *trailer
	O []uint64
}

func bplistDecode(d *decode.D, _ any) any {
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 6, d.AssertStr("bplist"))
		d.FieldUTF8("version", 2, d.AssertStr("00"))
	})

	t := new(trailer)

	d.SeekAbs(d.Len()-32*8, func(d *decode.D) {
		d.FieldStruct("trailer", func(d *decode.D) {
			d.FieldU40("padding") // unused
			d.FieldS8("sort_version")
			t.offTblOffSize = d.FieldS8("offset_table_offset_size")
			t.objRefSize = d.FieldS8("object_reference_size")
			t.nObjects = d.FieldS64BE("object_count")
			t.topObjectOffset = d.FieldS64BE("top_object_offset")
			t.offsetTableStart = d.FieldS64BE("offset_table_start")
		})
	})

	p := new(plist)
	p.T = t

	d.SeekAbs(t.offsetTableStart*8, func(d *decode.D) {
		i := int64(0)
		d.FieldArrayLoop("offset_table",
			func() bool { return i < t.nObjects },
			func(d *decode.D) {
				off := d.FieldU("element", 8*int(t.offTblOffSize))
				p.O = append(p.O, off)
				i++
			},
		)
	})

	d.SeekAbs(int64(p.O[0] * 8))

	d.FieldStruct("objects",
		func(d *decode.D) {
			decodeItem(d, p)
		})
	return nil
}
