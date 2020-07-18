package tiff

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html
// TODO: gps

import (
	"fq/pkg/decode"
	"log"
)

var File = &decode.Format{
	Name:  "tiff",
	MIMEs: []string{"image/tiff"},
	New:   func() decode.Decoder { return &FileDecoder{} },
}

const littleEndian = 0x49492a00
const bigEndian = 0x4d4d002a

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

const (
	imageWidth                  = 256
	imageLength                 = 257
	bitsPerSample               = 258
	compression                 = 259
	photometricInterpretation   = 262
	imageDescription            = 270
	make                        = 271
	model                       = 272
	stripOffsets                = 273
	orientation                 = 274
	samplesPerPixel             = 277
	rowsPerStrip                = 278
	stripByteCounts             = 279
	planarConfiguration         = 284
	xResolution                 = 282
	yResolution                 = 283
	resolutionUnit              = 296
	transferFunction            = 301
	software                    = 305
	dateTime                    = 306
	artist                      = 315
	whitePoint                  = 318
	primaryChromaticities       = 319
	jpegInterchangeFormat       = 513
	jpegInterchangeFormatLength = 514
	yCbCrCoefficients           = 529
	yCbCrSubSampling            = 530
	yCbCrPositioning            = 531
	referenceBlackWhite         = 532
	copyright                   = 33432
	exifTag                     = 34665
	gpsTag                      = 34853
)

var tagNames = map[uint64]string{
	imageWidth:                  "ImageWidth",
	imageLength:                 "ImageLength",
	bitsPerSample:               "BitsPerSample",
	compression:                 "Compression",
	photometricInterpretation:   "PhotometricInterpretation",
	imageDescription:            "ImageDescription",
	make:                        "Make",
	model:                       "Model",
	stripOffsets:                "StripOffsets",
	orientation:                 "Orientation",
	samplesPerPixel:             "SamplesPerPixel",
	rowsPerStrip:                "RowsPerStrip",
	stripByteCounts:             "StripByteCounts",
	planarConfiguration:         "PlanarConfiguration",
	xResolution:                 "XResolution",
	yResolution:                 "YResolution",
	resolutionUnit:              "ResolutionUnit",
	transferFunction:            "TransferFunction",
	software:                    "Software",
	dateTime:                    "DateTime",
	artist:                      "Artist",
	whitePoint:                  "WhitePoint",
	primaryChromaticities:       "PrimaryChromaticities",
	jpegInterchangeFormat:       "JPEGInterchangeFormat",
	jpegInterchangeFormatLength: "JPEGInterchangeFormatLength",
	yCbCrCoefficients:           "YCbCrCoefficients",
	yCbCrSubSampling:            "YCbCrSubSampling",
	yCbCrPositioning:            "YCbCrPositioning",
	referenceBlackWhite:         "ReferenceBlackWhite",
	copyright:                   "Copyright",
	exifTag:                     "ExifTag",
	gpsTag:                      "GPSTag",
}

// FileDecoder is a TIFF decoder
type FileDecoder struct {
	decode.Common
}

// Decode TIFF file
func (d *FileDecoder) Decode() {
	switch d.PeekBits(32) {
	case littleEndian, bigEndian:
	default:
		d.Invalid("unknown endian")
	}
	var fu16 func(name string) uint64
	var u16 func() uint64
	var fu32 func(name string) uint64

	endian := d.FieldUFn("endian", func() (uint64, decode.NumberFormat, string) {
		endian := d.U32()
		d.SeekRel(-4 * 8)
		d.FieldUTF8("order", 2)
		d.FieldU16("integer_42")
		switch endian {
		case littleEndian:
			fu16 = d.FieldU16LE
			u16 = d.U16LE
			fu32 = d.FieldU32LE
			return endian, decode.NumberHex, "little-endian"
		case bigEndian:
			fu16 = d.FieldU16BE
			u16 = d.U16BE
			fu32 = d.FieldU32BE
			return endian, decode.NumberHex, "big-endian"
		}
		return endian, decode.NumberDecimal, "unknown"
	})

	ifdOffset := fu32("ifd_offset")

	_ = ifdOffset

	// TODO: inf loop?
	for ifdOffset != 0 {
		log.Printf("ifdOffset: %#+v\n", ifdOffset)
		d.SeekAbs(ifdOffset * 8)

		d.FieldNoneFn("ifd", func() {
			numberOfFields := fu16("number_of_field")
			for i := uint64(0); i < numberOfFields; i++ {
				d.FieldNoneFn("ifd", func() {
					d.FieldStringMapFn("tag", tagNames, "unknown", u16)
					typ := d.FieldStringMapFn("type", typeNames, "unknown", u16)
					count := fu32("count")
					// TODO: short values stored in valueOffset directly?
					valueOffset := fu32("value_offset")
					_ = valueOffset
					_ = count
					_ = typ
				})
			}

			ifdOffset = fu32("next_ifd")
		})
	}

	_ = endian

}
