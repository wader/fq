package asn1

// T-REC-X.690-200811 (BER, DER, CER)
// https://www.itu.int/ITU-T/studygroups/com10/languages/X.690_1297.pdf
// https://cdn.standards.iteh.ai/samples/12285/039296509e8b40f3b25ba025de60365d/ISO-6093-1985.pdf
// https://en.wikipedia.org/wiki/X.690
// https://letsencrypt.org/docs/a-warm-welcome-to-asn1-and-der/
// https://luca.ntop.org/Teaching/Appunti/asn1.html
// https://lapo.it/asn1js/

// TODO: schema
// TODO: der/cer via mode?
// TODO: better torepr
// TODO: utc time
// TODO: validate CER DER
// TODO: bigrat?

import (
	"embed"
	"math"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed asn1_ber.jq
//go:embed asn1_ber.md
var asn1FS embed.FS

func init() {
	interp.RegisterFormat(
		format.ASN1_BER,
		&decode.Format{
			Description: "ASN1 BER (basic encoding rules, also CER and DER)",
			DecodeFn:    decodeASN1BER,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(asn1FS)
}

const (
	classUniversal   = 0b00
	classApplication = 0b01
	classContext     = 0b10
	classPrivate     = 0b11
)

var tagClassMap = scalar.UintMapSymStr{
	classUniversal:   "universal",
	classApplication: "application",
	classContext:     "context",
	classPrivate:     "private",
}

const (
	formPrimitive   = 0
	formConstructed = 1
)

var constructedPrimitiveMap = scalar.UintMapSymStr{
	formConstructed: "constructed",
	formPrimitive:   "primitive",
}

const (
	universalTypeEndOfContent     = 0x00
	universalTypeBoolean          = 0x01
	universalTypeInteger          = 0x02
	universalTypeBitString        = 0x03
	universalTypeOctetString      = 0x04
	universalTypeNull             = 0x05
	universalTypeObjectIdentifier = 0x06
	universalTypeObjectDescriptor = 0x07 // not encoded, just documentation?
	universalTypeExternal         = 0x08
	universalTypeReal             = 0x09
	universalTypeEnumerated       = 0x0a
	universalTypeEmbedded         = 0x0b
	universalTypeUTF8string       = 0x0c
	universalTypeSequence         = 0x10
	universalTypeSet              = 0x11
	universalTypeNumericString    = 0x12
	universalTypePrintableString  = 0x13
	universalTypeTeletexString    = 0x14
	universalTypeVideotexString   = 0x15
	universalTypeIA5String        = 0x16
	universalTypeUTCTime          = 0x17
	universalTypeGeneralizedtime  = 0x18
	universalTypeGraphicString    = 0x19 // not encoded?
	universalTypeVisibleString    = 0x1a
	universalTypeGeneralString    = 0x1b
	universalTypeUniversalString  = 0x1c // not encoded?
)

var universalTypeMap = scalar.UintMapSymStr{
	universalTypeEndOfContent:     "end_of_content",
	universalTypeBoolean:          "boolean",
	universalTypeInteger:          "integer",
	universalTypeBitString:        "bit_string",
	universalTypeOctetString:      "octet_string",
	universalTypeNull:             "null",
	universalTypeObjectIdentifier: "object_identifier",
	universalTypeObjectDescriptor: "object_descriptor",
	universalTypeExternal:         "external",
	universalTypeReal:             "real",
	universalTypeEnumerated:       "enumerated",
	universalTypeEmbedded:         "embedded",
	universalTypeUTF8string:       "utf8_string",
	universalTypeSequence:         "sequence",
	universalTypeSet:              "set",
	universalTypeNumericString:    "numeric_string",
	universalTypePrintableString:  "printable_string",
	universalTypeTeletexString:    "teletex_string",
	universalTypeVideotexString:   "videotex_string",
	universalTypeIA5String:        "ia5_string",
	universalTypeUTCTime:          "utc_time",
	universalTypeGeneralizedtime:  "generalized_time",
	universalTypeGraphicString:    "graphic_string",
	universalTypeVisibleString:    "visible_string",
	universalTypeGeneralString:    "general_string",
	universalTypeUniversalString:  "universal_string",
}

const (
	lengthIndefinite = 0
	lengthEndMarker  = 0x00_00
)

const (
	decimalPlusInfinity  = 0b00_00_00
	decimalMinusInfinity = 0b00_00_01
	decimalNan           = 0b00_00_10
	decimalMinusZero     = 0b00_00_11
)

var lengthMap = scalar.UintMapSymStr{
	0: "indefinite",
}

func decodeLength(d *decode.D) uint64 {
	n := d.U8()
	if n&0b1000_0000 != 0 {
		n = n & 0b0111_1111
		if n == 0 {
			return lengthIndefinite
		}
		if n == 127 {
			d.Errorf("length 127 reserved")
		}
		// TODO: bigint
		return d.U(int(n) * 8)
	}
	return n & 0b0111_1111
}

// TODO: bigint?
func decodeTagNumber(d *decode.D) uint64 {
	v := d.U5()
	moreBytes := v == 0b11111
	for moreBytes {
		moreBytes = d.Bool()
		v = v<<7 | d.U7()
	}
	return v
}

func decodeASN1BERValue(d *decode.D, bib *bitio.Buffer, sb *strings.Builder, parentForm uint64, parentTag uint64) {
	class := d.FieldU2("class", tagClassMap)
	form := d.FieldU1("form", constructedPrimitiveMap)

	// TODO: verify
	// TODO: constructed types verify
	_ = parentTag
	_ = parentForm

	var tag uint64
	switch class {
	case classUniversal:
		tag = d.FieldUintFn("tag", decodeTagNumber, universalTypeMap, scalar.UintHex)
	default:
		tag = d.FieldUintFn("tag", decodeTagNumber)
	}

	length := d.FieldUintFn("length", decodeLength, lengthMap)
	var l int64
	switch length {
	case lengthIndefinite:
		// null has zero length byte
		if !(class == classUniversal && tag == universalTypeNull) && form == formPrimitive {
			d.Fatalf("primitive with indefinite length")
		}
		l = d.BitsLeft()
	default:
		l = int64(length) * 8
	}

	d.LimitedFn(l, func(d *decode.D) {
		switch {
		case form == formConstructed || tag == universalTypeSequence || tag == universalTypeSet:
			d.FieldArray("constructed", func(d *decode.D) {
				for !d.End() {
					if length == lengthIndefinite && d.PeekUintBits(16) == lengthEndMarker {
						break
					}

					if form == formConstructed && bib == nil && sb == nil {
						switch tag {
						case universalTypeBitString:
							bib = &bitio.Buffer{}
						case universalTypeOctetString:
							bib = &bitio.Buffer{}
						case universalTypeUTF8string,
							universalTypeNumericString,
							universalTypePrintableString,
							universalTypeTeletexString,
							universalTypeVideotexString,
							universalTypeIA5String,
							universalTypeUTCTime,
							universalTypeVisibleString, // not encoded?
							universalTypeGeneralString: // not encoded?
							sb = &strings.Builder{}
						}
					}

					d.FieldStruct("object", func(d *decode.D) { decodeASN1BERValue(d, bib, sb, form, tag) })
				}
			})

			if length == lengthIndefinite {
				d.FieldU16("end_marker")
			}
			if form == formConstructed {
				switch tag {
				case universalTypeBitString:
					if bib != nil {
						buf, bufLen := bib.Bits()
						d.FieldRootBitBuf("value", bitio.NewBitReader(buf, bufLen))
					}
				case universalTypeOctetString:
					if bib != nil {
						buf, bufLen := bib.Bits()
						d.FieldRootBitBuf("value", bitio.NewBitReader(buf, bufLen))
					}
				case universalTypeUTF8string,
					universalTypeNumericString,
					universalTypePrintableString,
					universalTypeTeletexString,
					universalTypeVideotexString,
					universalTypeIA5String,
					universalTypeUTCTime,
					universalTypeVisibleString, // not encoded?
					universalTypeGeneralString: // not encoded?
					if sb != nil {
						d.FieldValueStr("value", sb.String())
					}
				}
			}
		case class == classUniversal && tag == universalTypeEndOfContent:
			// nop
		case class == classUniversal && tag == universalTypeBoolean:
			d.FieldU8("value", scalar.UintRangeToScalar{
				{Range: [2]uint64{0, 0}, S: scalar.Uint{Sym: false}},
				{Range: [2]uint64{0x01, 0xff1}, S: scalar.Uint{Sym: true}},
			})
		case class == classUniversal && tag == universalTypeInteger:
			if length > 8 {
				d.FieldSBigInt("value", int(length)*8)
			} else {
				d.FieldS("value", int(length)*8)
			}
		case class == classUniversal && tag == universalTypeBitString:
			unusedBitsCount := d.FieldU8("unused_bits_count")
			if unusedBitsCount > 7 {
				d.Fatalf("unusedBitsCount %d > 7", unusedBitsCount)
			}
			br := d.FieldRawLen("value", int64(length-1)*8-int64(unusedBitsCount))
			if bib != nil {
				// TODO: helper?
				if _, err := bitio.Copy(bib, br); err != nil {
					d.IOPanic(err, "value", "bitio.Copy")
				}
			}
			if unusedBitsCount > 0 {
				d.FieldRawLen("unused_bits", int64(unusedBitsCount))
			}
		case class == classUniversal && tag == universalTypeOctetString:
			br := d.FieldRawLen("value", int64(length)*8)
			if bib != nil {
				// TODO: helper?
				if _, err := bitio.Copy(bib, br); err != nil {
					d.IOPanic(err, "value", "bitio.Copy")
				}
			}
		case class == classUniversal && tag == universalTypeNull:
			d.FieldValueAny("value", nil)
		case class == classUniversal && tag == universalTypeObjectIdentifier:
			d.FieldArray("value", func(d *decode.D) {
				// first byte is = oid0*40 + oid1
				d.FieldUintFn("oid", func(d *decode.D) uint64 { return d.U8() / 40 })
				d.SeekRel(-8)
				d.FieldUintFn("oid", func(d *decode.D) uint64 { return d.U8() % 40 })
				for !d.End() {
					d.FieldUintFn("oid", func(d *decode.D) uint64 {
						more := true
						var n uint64
						for more {
							b := d.U8()
							n = n<<7 | b&0b0111_1111
							more = b&0b1000_0000 != 0
						}
						return n
					})
				}
			})
		case class == classUniversal && tag == universalTypeObjectDescriptor: // not encoded, just documentation?
			// nop
		case class == classUniversal && tag == universalTypeExternal:
			d.FieldRawLen("value", int64(length)*8)
		case class == classUniversal && tag == universalTypeReal:
			switch {
			case length == 0:
				d.FieldValueUint("value", 0)
			default:
				switch d.FieldBool("binary_encoding") {
				case true:
					s := d.FieldScalarBool("sign", scalar.BoolMapSymSint{
						true:  -1,
						false: 1,
					}).SymSint()
					base := d.FieldScalarU2("base", scalar.UintMapSymUint{
						0b00: 2,
						0b01: 8,
						0b10: 16,
						0b11: 0,
					}).SymUint()
					scale := d.FieldU2("scale")
					format := d.FieldU2("format")

					var exp int64
					switch format {
					case 0b00:
						exp = d.FieldS8("exp")
					case 0b01:
						exp = d.FieldS16("exp")
					case 0b10:
						exp = d.FieldS24("exp")
					default:
						n := d.FieldU8("exp_bytes")
						// TODO: bigint?
						exp = d.FieldS("exp", int(n)*8)
					}

					n := d.FieldU("n", int(d.BitsLeft()))
					m := float64(s) * float64(n) * math.Pow(float64(base), float64(exp)) * float64(int(1)<<scale)
					d.FieldValueFlt("value", m)

				case false:
					switch d.FieldBool("decimal_encoding") {
					case true:
						n := d.FieldU6("special", scalar.UintMapSymStr{
							decimalPlusInfinity:  "plus_infinity",
							decimalMinusInfinity: "minus_infinity",
							decimalNan:           "nan",
							decimalMinusZero:     "minus_zero",
						})

						switch n {
						case decimalPlusInfinity:
							d.FieldValueFlt("value", math.Inf(1))
						case decimalMinusInfinity:
							d.FieldValueFlt("value", math.Inf(-1))
						case decimalNan:
							d.FieldValueFlt("value", math.NaN())
						case decimalMinusZero:
							d.FieldValueFlt("value", -0)
						}
					case false:
						d.FieldU6("representation", scalar.UintMapSymStr{
							0b00_00_01: "nr1",
							0b00_00_10: "nr2",
							0b00_00_11: "nr3",
						})
						d.FieldFltFn("value", func(d *decode.D) float64 {
							// TODO: can ParseFloat do all ISO-6093 nr?
							n, _ := strconv.ParseFloat(d.UTF8(int(d.BitsLeft()/8)), 64)
							return n
						})
					}
				}
			}
		case class == classUniversal && tag == universalTypeUTF8string,
			class == classUniversal && tag == universalTypeNumericString,
			class == classUniversal && tag == universalTypePrintableString,
			class == classUniversal && tag == universalTypeTeletexString,
			class == classUniversal && tag == universalTypeVideotexString,
			class == classUniversal && tag == universalTypeIA5String,
			class == classUniversal && tag == universalTypeUTCTime,
			class == classUniversal && tag == universalTypeVisibleString, // not encoded?
			class == classUniversal && tag == universalTypeGeneralString: // not encoded?
			// TODO: restrict?
			s := d.FieldUTF8("value", int(length))
			if sb != nil {
				sb.WriteString(s)
			}
		case class == classUniversal && tag == universalTypeGeneralizedtime:
			d.FieldRawLen("value", int64(length)*8)
		default:
			d.FieldRawLen("value", l)
		}
	})
}

func decodeASN1BER(d *decode.D) any {
	decodeASN1BERValue(d, nil, nil, formConstructed, universalTypeSequence)
	return nil
}
