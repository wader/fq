package tiff

// https://www.adobe.io/content/dam/udp/en/open/standards/tiff/TIFF6.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var tiffIccProfile decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.TIFF,
		Description: "Tag Image File Format",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    tiffDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ICC_PROFILE}, Group: &tiffIccProfile},
		},
	})
}

const littleEndian = 0x49492a00 // "II*\0"
const bigEndian = 0x4d4d002a    // "MM\0*"

var endianNames = scalar.UToSymStr{
	littleEndian: "little-endian",
	bigEndian:    "big-endian",
}

const (
	BYTE      = 1
	ASCII     = 2
	SHORT     = 3
	LONG      = 4
	RATIONAL  = 5
	UNDEFINED = 7
	SLONG     = 9
	SRATIONAL = 10
)

var typeNames = scalar.UToSymStr{
	BYTE:      "BYTE",
	ASCII:     "ASCII",
	SHORT:     "SHORT",
	LONG:      "LONG",
	RATIONAL:  "RATIONAL",
	UNDEFINED: "UNDEFINED",
	SLONG:     "SLONG",
	SRATIONAL: "SRATIONAL",
}

// TODO: tiff 6.0 types
var typeByteSize = map[uint64]uint64{
	BYTE:      1,
	ASCII:     1,
	SHORT:     2,
	LONG:      4,
	RATIONAL:  4 + 4,
	UNDEFINED: 1,
	SLONG:     4,
	SRATIONAL: 4 + 4,
}

func fieldRational(d *decode.D, name string) float64 {
	var v float64
	d.FieldStruct(name, func(d *decode.D) {
		numerator := d.FieldU32("numerator")
		denominator := d.FieldU32("denominator")
		v := float64(numerator) / float64(denominator)
		d.FieldValueFloat("float", v)
	})
	return v
}

func fieldSRational(d *decode.D, name string) float64 {
	var v float64
	d.FieldStruct(name, func(d *decode.D) {
		numerator := d.FieldS32("numerator")
		denominator := d.FieldS32("denominator")
		v := float64(numerator) / float64(denominator)
		d.FieldValueFloat("float", v)
	})
	return v
}

type strips struct {
	offsets    []int64
	byteCounts []int64
}

func decodeIfd(d *decode.D, s *strips, tagNames scalar.UToSymStr) int64 {
	var nextIfdOffset int64

	d.FieldStruct("ifd", func(d *decode.D) {
		numberOfFields := d.FieldU16("number_of_field")
		d.FieldArray("entries", func(d *decode.D) {
			for i := uint64(0); i < numberOfFields; i++ {
				d.FieldStruct("entry", func(d *decode.D) {
					tag := d.FieldU16("tag", tagNames, scalar.ActualHex)
					typ := d.FieldU16("type", typeNames)
					count := d.FieldU32("count")
					valueOrByteOffset := d.FieldU32("value_offset")

					if _, ok := typeNames[typ]; !ok {
						return
					}

					valueByteOffset := valueOrByteOffset
					valueByteSize := typeByteSize[typ] * count
					if valueByteSize <= 4 {
						// if value fits in offset itself use offset to value_offset
						valueByteOffset = uint64(d.Pos()/8) - 4
					}

					switch {
					case typ == LONG && (tag == ExifIFD || tag == GPSInfo):
						ifdPos := valueOrByteOffset
						pos := d.Pos()
						d.SeekAbs(int64(ifdPos * 8))

						switch tag {
						case ExifIFD:
							// TODO: exif tag names?
							decodeIfd(d, &strips{}, tiffTagNames)
						case GPSInfo:
							decodeIfd(d, &strips{}, gpsInfoTagNames)
						}

						d.SeekAbs(pos)
					default:

						d.FieldArray("values", func(d *decode.D) {
							switch {
							case typ == UNDEFINED:
								switch tag {
								case InterColorProfile:
									d.FieldFormatRange("icc", int64(valueByteOffset)*8, int64(valueByteSize)*8, tiffIccProfile, nil)
								default:
									d.RangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
										d.FieldRawLen("value", d.BitsLeft())
									})
								}
							case typ == ASCII:
								d.RangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
									d.FieldUTF8NullFixedLen("value", int(valueByteSize))
								})
							case typ == BYTE:
								d.RangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
									d.FieldRawLen("value", d.BitsLeft())
								})
							default:
								d.RangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
									for i := uint64(0); i < count; i++ {
										switch typ {
										// TODO: only some typ?
										case BYTE:
											d.FieldU8("value")
										case SHORT:
											v := d.FieldU16("value")
											_ = v
											switch tag {
											case StripOffsets:
												s.offsets = append(s.offsets, int64(v*8))
											case StripByteCounts:
												s.byteCounts = append(s.byteCounts, int64(v*8))
											}
										case LONG:
											v := d.FieldU32("value")
											_ = v
											switch tag {
											case StripOffsets:
												s.offsets = append(s.offsets, int64(v*8))
											case StripByteCounts:
												s.byteCounts = append(s.byteCounts, int64(v*8))
											}
										case RATIONAL:
											fieldRational(d, "value")
										case SLONG:
											d.FieldS32("value")
										case SRATIONAL:
											fieldSRational(d, "value")
										default:
											d.Errorf("unknown type")
										}
									}
								})
							}
						})
					}
				})
			}
		})

		nextIfdOffset = int64(d.FieldU32("next_ifd"))
	})

	return nextIfdOffset
}

func tiffDecode(d *decode.D, in any) any {
	endian := d.FieldU32("endian", endianNames, scalar.ActualHex)

	switch endian {
	case littleEndian:
		d.Endian = decode.LittleEndian
	case bigEndian:
		d.Endian = decode.BigEndian
	default:
		d.Fatalf("unknown endian")
	}

	d.SeekRel(-4 * 8)

	d.FieldUTF8("order", 2, d.AssertStr("II", "MM"))
	d.FieldU16("integer_42", d.AssertU(42))

	ifdOffset := int64(d.FieldU32("first_ifd"))
	s := &strips{}

	// to catch infinite loops
	ifdSeen := map[int64]struct{}{}

	d.FieldArray("ifds", func(d *decode.D) {
		for ifdOffset != 0 {
			if _, ok := ifdSeen[ifdOffset]; ok {
				d.Fatalf("ifd loop detected for %d", ifdOffset)
			}
			ifdSeen[ifdOffset] = struct{}{}
			d.SeekAbs(ifdOffset * 8)
			ifdOffset = decodeIfd(d, s, tiffTagNames)
		}
	})

	if len(s.offsets) != len(s.byteCounts) {
		// TODO: warning
	} else {
		d.FieldArray("strips", func(d *decode.D) {
			for i := 0; i < len(s.offsets); i++ {
				d.RangeFn(s.offsets[i], s.byteCounts[i], func(d *decode.D) {
					d.FieldRawLen("strip", d.BitsLeft())
				})
			}
		})
	}

	return nil
}
