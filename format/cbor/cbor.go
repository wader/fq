package cbor

// https://en.wikipedia.org/wiki/CBOR
// https://www.rfc-editor.org/rfc/rfc8949.html

// TODO: streaming bytes test?
// TODO: decode some sematic tags

import (
	"bytes"
	"embed"
	"math/big"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed cbor.jq
var cborFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.CBOR,
		Description: "Concise Binary Object Representation",
		DecodeFn:    decodeCBOR,
		Functions:   []string{"torepr", "_help"},
	})
	interp.RegisterFS(cborFS)
}

type majorTypeEntry struct {
	s scalar.S
	d func(d *decode.D, shortCount uint64, count uint64) any
}

type majorTypeEntries map[uint64]majorTypeEntry

func (mts majorTypeEntries) MapScalar(s scalar.S) (scalar.S, error) {
	u := s.ActualU()
	if fe, ok := mts[u]; ok {
		s = fe.s
		s.Actual = u
	}
	return s, nil
}

const (
	shortCountVariable8Bit  = 24
	shortCountVariable16Bit = 25
	shortCountVariable32Bit = 26
	shortCountVariable64Bit = 27
	shortCountIndefinite    = 31

	shortCountSpecialFalse     = 20
	shortCountSpecialTrue      = 21
	shortCountSpecialNull      = 22
	shortCountSpecialUndefined = 23

	shortCountSpecialFloat16Bit = 25
	shortCountSpecialFloat32Bit = 26
	shortCountSpecialFloat64Bit = 27
)

var shortCountMap = scalar.UToSymStr{
	shortCountVariable8Bit:  "8bit",
	shortCountVariable16Bit: "16bit",
	shortCountVariable32Bit: "32bit",
	shortCountVariable64Bit: "64bit",
	shortCountIndefinite:    "indefinite",
}

var tagMap = scalar.UToSymStr{
	0:     "date_time",
	1:     "epoch_date_time",
	2:     "unsigned_bignum",
	3:     "negative_bignum",
	4:     "decimal_fraction",
	5:     "bigfloat",
	21:    "base64url",
	22:    "base64",
	23:    "base16",
	24:    "encoded_cbor",
	32:    "uri",
	33:    "base64url",
	34:    "base64",
	36:    "mime_message",
	55799: "self_described_cbor",
}

const (
	majorTypePositiveInt  = 0
	majorTypeNegativeInt  = 1
	majorTypeBytes        = 2
	majorTypeUTF8         = 3
	majorTypeArray        = 4
	majorTypeMap          = 5
	majorTypeSematic      = 6
	majorTypeSpecialFloat = 7
)

const (
	breakMarker = 0xff
)

func decodeCBORValue(d *decode.D) any {
	majorTypeMap := majorTypeEntries{
		majorTypePositiveInt: {s: scalar.S{Sym: "positive_int"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			d.FieldValueU("value", count)
			return nil
		}},
		majorTypeNegativeInt: {s: scalar.S{Sym: "negative_int"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			n := new(big.Int)
			n.SetUint64(count).Neg(n).Sub(n, mathex.BigIntOne)
			d.FieldValueBigInt("value", n)
			return nil
		}},
		majorTypeBytes: {s: scalar.S{Sym: "bytes"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			if shortCount == shortCountIndefinite {
				bb := &bytes.Buffer{}
				d.FieldArray("items", func(d *decode.D) {
					for d.PeekBits(8) != breakMarker {
						d.FieldStruct("item", func(d *decode.D) {
							v := decodeCBORValue(d)
							switch v := v.(type) {
							case []byte:
								bb.Write(v)
							default:
								d.Fatalf("non-bytes in bytes stream %v", v)
							}
						})
					}
				})
				d.FieldRootBitBuf("value", bitio.NewBitReader(bb.Bytes(), -1))
				// nil, nested indefinite bytes is not allowed
				return nil
			}

			buf := d.ReadAllBits(d.FieldRawLen("value", int64(count)*8))

			return buf
		}},
		majorTypeUTF8: {s: scalar.S{Sym: "utf8"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			if shortCount == shortCountIndefinite {
				sb := &strings.Builder{}
				d.FieldArray("items", func(d *decode.D) {
					for d.PeekBits(8) != breakMarker {
						d.FieldStruct("item", func(d *decode.D) {
							v := decodeCBORValue(d)
							switch v := v.(type) {
							case string:
								sb.WriteString(v)
							default:
								d.Fatalf("non-string in string stream %v", v)
							}
						})
					}
				})
				d.FieldValueStr("value", sb.String())
				// nil, nested indefinite string is not allowed
				return nil
			}

			return d.FieldUTF8("value", int(count))
		}},
		majorTypeArray: {s: scalar.S{Sym: "array"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			d.FieldArray("elements", func(d *decode.D) {
				for i := uint64(0); true; i++ {
					if shortCount == shortCountIndefinite && d.PeekBits(8) == breakMarker {
						break
					} else if i >= count {
						break
					}
					d.FieldStruct("element", func(d *decode.D) { decodeCBORValue(d) })
				}
			})
			if shortCount == shortCountIndefinite {
				d.FieldU8("break")
			}
			return nil
		}},
		majorTypeMap: {s: scalar.S{Sym: "map"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			d.FieldArray("pairs", func(d *decode.D) {
				for i := uint64(0); true; i++ {
					if shortCount == shortCountIndefinite && d.PeekBits(8) == breakMarker {
						break
					} else if i >= count {
						break
					}
					d.FieldStruct("pair", func(d *decode.D) {
						d.FieldStruct("key", func(d *decode.D) { decodeCBORValue(d) })
						d.FieldStruct("value", func(d *decode.D) { decodeCBORValue(d) })
					})
				}
			})
			if shortCount == shortCountIndefinite {
				d.FieldU8("break")
			}
			return nil
		}},
		majorTypeSematic: {s: scalar.S{Sym: "semantic"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			d.FieldValueU("tag", count, tagMap)
			d.FieldStruct("value", func(d *decode.D) { decodeCBORValue(d) })
			return nil
		}},
		majorTypeSpecialFloat: {s: scalar.S{Sym: "special_float"}, d: func(d *decode.D, shortCount uint64, count uint64) any {
			switch shortCount {
			// TODO: 0-19
			case shortCountSpecialFalse:
				d.FieldValueBool("value", false)
			case shortCountSpecialTrue:
				d.FieldValueBool("value", true)
			case shortCountSpecialNull:
				d.FieldValueNil("value")
			case shortCountSpecialUndefined:
				// TODO: undefined
			case 24:
				// TODO: future
			case shortCountSpecialFloat16Bit:
				d.FieldF16("value")
			case shortCountSpecialFloat32Bit:
				d.FieldF32("value")
			case shortCountSpecialFloat64Bit:
				d.FieldF64("value")
			case 28, 29, 30:
				// TODO: future
			}
			return nil
		}},
	}

	typ := d.FieldU3("major_type", majorTypeMap)
	shortCount := d.FieldU5("short_count", shortCountMap)
	count := shortCount
	if typ != majorTypeSpecialFloat {
		switch count {
		// 0-23 value in shortCount
		case shortCountVariable8Bit:
			count = d.FieldU8("variable_count")
		case shortCountVariable16Bit:
			count = d.FieldU16("variable_count")
		case shortCountVariable32Bit:
			count = d.FieldU32("variable_count")
		case shortCountVariable64Bit:
			count = d.FieldU64("variable_count")
		case 28, 29, 30:
			d.Fatalf("incorrect shortCount %d", count)
		}
	}

	if mt, ok := majorTypeMap[typ]; ok {
		if mt.d != nil {
			return mt.d(d, shortCount, count)
		}
		return nil
	}

	panic("unreachable")
}

func decodeCBOR(d *decode.D, _ any) any {
	decodeCBORValue(d)
	return nil
}
