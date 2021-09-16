package tiff

// https://www.adobe.io/content/dam/udp/en/open/standards/tiff/TIFF6.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var tiffIccProfile []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.TIFF,
		Description: "Tag Image File Format",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    tiffDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ICC_PROFILE}, Formats: &tiffIccProfile},
		},
	})
}

const littleEndian = 0x49492a00 // "II*\0"
const bigEndian = 0x4d4d002a    // "MM\0*"

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

var typeNames = map[uint64]string{
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
	d.FieldStructFn(name, func(d *decode.D) {
		numerator := d.FieldU32("numerator")
		denominator := d.FieldU32("denominator")
		v := float64(numerator) / float64(denominator)
		d.FieldFloatFn("float", func() (float64, string) {
			return v, ""
		})
	})
	return v
}

func fieldSRational(d *decode.D, name string) float64 {
	var v float64
	d.FieldStructFn(name, func(d *decode.D) {
		numerator := d.FieldS32("numerator")
		denominator := d.FieldS32("denominator")
		v := float64(numerator) / float64(denominator)
		d.FieldFloatFn("float", func() (float64, string) {
			return v, ""
		})
	})
	return v
}

type strips struct {
	offsets    []int64
	byteCounts []int64
}

func decodeIfd(d *decode.D, s *strips, tagNames map[uint64]string) int64 {
	var nextIfdOffset int64

	d.FieldStructFn("ifd", func(d *decode.D) {
		numberOfFields := d.FieldU16("number_of_field")
		d.FieldArrayFn("entries", func(d *decode.D) {
			for i := uint64(0); i < numberOfFields; i++ {
				d.FieldStructFn("entry", func(d *decode.D) {
					tag, _ := d.FieldStringMapFn("tag", tagNames, "unknown", d.U16, decode.NumberHex)
					typ, typOk := d.FieldStringMapFn("type", typeNames, "unknown", d.U16, decode.NumberDecimal)
					count := d.FieldU32("count")
					// TODO: short values stored in valueOffset directly?
					valueOrByteOffset := d.FieldU32("value_offset")

					if !typOk {
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

						d.FieldArrayFn("values", func(d *decode.D) {
							switch {
							case typ == UNDEFINED:
								switch tag {
								case InterColorProfile:
									d.FieldFormatRange("icc", int64(valueByteOffset)*8, int64(valueByteSize)*8, tiffIccProfile, nil)
								default:
									// log.Printf("tag: %#+v\n", tag)
									// log.Printf("valueByteSize: %#+v\n", valueByteSize)
									d.FieldBitBufRange("value", int64(valueByteOffset)*8, int64(valueByteSize)*8)
								}
							case typ == ASCII:
								d.DecodeRangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
									d.FieldUTF8("value", int(valueByteSize))
								})
							case typ == BYTE:
								d.FieldBitBufRange("value", int64(valueByteOffset*8), int64(valueByteSize*8))
							default:
								// log.Printf("valueOffset: %d\n", valueByteOffset)
								// log.Printf("valueSize: %d\n", valueByteSize)
								d.DecodeRangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func(d *decode.D) {
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
											panic("unknown type")
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

func tiffDecode(d *decode.D, in interface{}) interface{} {
	switch d.PeekBits(32) {
	case littleEndian, bigEndian:
		d.Endian = decode.BigEndian
	default:
		d.Invalid("unknown endian")
	}

	endian := d.FieldUFn("endian", func() (uint64, decode.DisplayFormat, string) {
		endian := d.U32()
		d.SeekRel(-4 * 8)
		d.FieldUTF8("order", 2)
		// TODO: validate?
		d.FieldU16("integer_42")
		switch endian {
		case littleEndian:
			return endian, decode.NumberHex, "little-endian"
		case bigEndian:
			return endian, decode.NumberHex, "big-endian"
		}
		return endian, decode.NumberDecimal, "unknown"
	})

	switch endian {
	case littleEndian:
		d.Endian = decode.LittleEndian
	case bigEndian:
		d.Endian = decode.BigEndian
	}

	ifdOffset := int64(d.FieldU32("first_ifd"))
	s := &strips{}

	d.FieldArrayFn("ifds", func(d *decode.D) {
		// TODO: inf loop?
		for ifdOffset != 0 {
			d.SeekAbs(ifdOffset * 8)
			ifdOffset = decodeIfd(d, s, tiffTagNames)
		}
	})

	if len(s.offsets) != len(s.byteCounts) {
		// TODO: warning
	} else {
		d.FieldArrayFn("strips", func(d *decode.D) {
			for i := 0; i < len(s.offsets); i++ {
				d.FieldBitBufRange("strip", s.offsets[i], s.byteCounts[i])
			}
		})
	}

	return nil
}
