package zip

// https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
// https://opensource.apple.com/source/zip/zip-6/unzip/unzip/proginfo/extra.fld

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

var probeFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ZIP,
		Description: "ZIP archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    zipDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
	})
}

const (
	compressionMethodNone                      = 0
	compressionMethodShrunk                    = 1
	compressionMethodReducedCompressionFactor1 = 2
	compressionMethodReducedCompressionFactor2 = 3
	compressionMethodReducedCompressionFactor3 = 4
	compressionMethodReducedCompressionFactor4 = 5
	compressionMethodImploded                  = 6
	compressionMethodDeflated                  = 8
	compressionMethodEnhancedDeflated          = 9
	compressionMethodPKWareDCLImploded         = 10
	compressionMethodBzip2                     = 12
	compressionMethodLZMA                      = 14
	compressionMethodIBMTERSE                  = 18
	compressionMethodIBMLZ77z                  = 19
	compressionMethodPPMd                      = 98
)

var compressionMethodMap = decode.UToStr{
	compressionMethodNone:                      "None",
	compressionMethodShrunk:                    "Shrunk",
	compressionMethodReducedCompressionFactor1: "ReducedCompressionFactor1",
	compressionMethodReducedCompressionFactor2: "ReducedCompressionFactor2",
	compressionMethodReducedCompressionFactor3: "ReducedCompressionFactor3",
	compressionMethodReducedCompressionFactor4: "ReducedCompressionFactor4",
	compressionMethodImploded:                  "Imploded",
	compressionMethodDeflated:                  "Deflated",
	compressionMethodEnhancedDeflated:          "EnhancedDeflated",
	compressionMethodPKWareDCLImploded:         "PKWareDCLImploded",
	compressionMethodBzip2:                     "Bzip2",
	compressionMethodLZMA:                      "LZMA",
	compressionMethodIBMTERSE:                  "IBMTERSE",
	compressionMethodIBMLZ77z:                  "IBMLZ77z",
	compressionMethodPPMd:                      "PPMd",
}

var (
	centralDirectorySignature       = []byte("PK\x01\x02")
	endOfCentralDirectorySignature  = []byte("PK\x05\x06")
	endOfCentralDirectorySignatureN = 0x06054b50
	localFileSignature              = []byte("PK\x03\x04")
	dataIndicatorSignature          = []byte("PK\x07\x08")
)

var headerIDMap = decode.UToScalar{
	0x0001: {Description: "ZIP64 extended information extra field"},
	0x0007: {Description: "AV Info"},
	0x0009: {Description: "OS/2 extended attributes"},
	0x000a: {Description: "NTFS (Win9x/WinNT FileTimes)"},
	0x000c: {Description: "OpenVMS"},
	0x000d: {Description: "Unix"},
	0x000f: {Description: "Patch Descriptor"},
	0x0014: {Description: "PKCS#7 Store for X.509 Certificates"},
	0x0015: {Description: "X.509 Certificate ID and Signature for individual file"},
	0x0016: {Description: "X.509 Certificate ID for Central Directory"},
	0x0065: {Description: "IBM S/390 attributes - uncompressed"},
	0x0066: {Description: "IBM S/390 attributes - compressed"},
	0x07c8: {Description: "Info-ZIP Macintosh (old, J. Lee)"},
	0x2605: {Description: "ZipIt Macintosh (first version)"},
	0x2705: {Description: "ZipIt Macintosh v 1.3.5 and newer (w/o full filename)"},
	0x334d: {Description: "Info-ZIP Macintosh (new, D. Haase's 'Mac3' field )"},
	0x4154: {Description: "Tandem NSK"},
	0x4341: {Description: "Acorn/SparkFS (David Pilling)"},
	0x4453: {Description: "Windows NT security descriptor (binary ACL)"},
	0x4704: {Description: "VM/CMS"},
	0x470f: {Description: "MVS"},
	// "inofficial" in original table
	//nolint:misspell
	0x4854: {Description: "Theos, old inofficial port"},
	0x4b46: {Description: "FWKCS MD5 (see below)"},
	0x4c41: {Description: "OS/2 access control list (text ACL)"},
	0x4d49: {Description: "Info-ZIP OpenVMS (obsolete)"},
	0x4d63: {Description: "Macintosh SmartZIP, by Macro Bambini"},
	0x4f4c: {Description: "Xceed original location extra field"},
	0x5356: {Description: "AOS/VS (binary ACL)"},
	0x5455: {Description: "extended timestamp"},
	0x5855: {Description: "Info-ZIP Unix (original; also OS/2, NT, etc.)"},
	0x554e: {Description: "Xceed unicode extra field"},
	0x6542: {Description: "BeOS (BeBox, PowerMac, etc.)"},
	0x6854: {Description: "Theos"},
	0x756e: {Description: "ASi Unix"},
	0x7855: {Description: "Info-ZIP Unix (new)"},
	0x7875: {Description: "UNIX UID/GID"},
	0xfb4a: {Description: "SMS/QDOS"},
}

// "MS-DOS uses year values relative to 1980 and 2 second precision."
func fieldMSDOSTime(d *decode.D) {
	d.FieldU5("hours")
	d.FieldU6("minutes")
	d.FieldU5("seconds")
}

func fieldMSDOSDate(d *decode.D) {
	d.FieldU7("year")
	d.FieldU4("month")
	d.FieldU5("day")
}

func zipDecode(d *decode.D, in interface{}) interface{} {
	// TODO: just decode instead?
	if !bytes.Equal(d.PeekBytes(4), []byte("PK\x03\x04")) {
		d.Errorf("expected PK header")
	}

	d.Endian = decode.LittleEndian

	d.SeekAbs(d.Len())

	// TODO: better EOCD probe
	p, _, err := d.TryPeekFind(32, -8, -10000, func(v uint64) bool {
		return v == uint64(endOfCentralDirectorySignatureN)
	})
	if err != nil {
		d.Fatalf("can't find end of central directory")
	}
	d.SeekAbs(d.Len() + p)

	var offsetCD uint64
	var sizeCD uint64
	var diskNr uint64

	d.FieldStruct("end_of_central_directory", func(d *decode.D) {
		d.FieldRawLen("signature", 4*8, d.ValidateBitBuf(endOfCentralDirectorySignature))
		diskNr = d.FieldU16("disk_nr")
		d.FieldU16("central_directory_start_disk_nr")
		d.FieldU16("nr_of_central_directory_records_on_disk")
		d.FieldU16("nr_of_central_directory_records")
		sizeCD = d.FieldU32("size_of_central directory")
		offsetCD = d.FieldU32("offset_of_start_of_central_directory")
		commentLength := d.FieldU16("comment_length")
		d.FieldUTF8("comment", int(commentLength))
	})

	var localFileOffsets []uint64

	d.SeekAbs(int64(offsetCD) * 8)
	d.FieldArray("central_directories", func(d *decode.D) {
		d.LenFn(int64(sizeCD)*8, func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("central_directory", func(d *decode.D) {
					d.FieldRawLen("signature", 4*8, d.ValidateBitBuf(centralDirectorySignature))
					d.FieldU16("version_made_by")
					d.FieldU16("version_needed")
					d.FieldStruct("flags", func(d *decode.D) {
						// TODO: 16LE, should have some kind of native endian flag reader helper?
						d.FieldU1("unused0")
						d.FieldBool("strong_encryption")
						d.FieldBool("compressed_patched_data")
						d.FieldBool("enhanced_deflation")
						d.FieldBool("data_descriptor")
						d.FieldBool("compression0")
						d.FieldBool("compression1")
						d.FieldBool("encrypted")

						d.FieldU2("reserved0")
						d.FieldBool("mask_header_values")
						d.FieldBool("reserved1")
						d.FieldBool("language_encoding")
						d.FieldU3("unused1")
					})
					d.FieldU16("compression_method", d.MapUToStrSym(compressionMethodMap))
					d.FieldStruct("last_modification_date", fieldMSDOSTime)
					d.FieldStruct("last_modification_time", fieldMSDOSDate)
					d.FieldU32("crc32_uncompressed", d.Hex)
					d.FieldU32("compressed_size")
					d.FieldU32("uncompressed_size")
					fileNameLength := d.FieldU16("file_name_length")
					extraFieldLength := d.FieldU16("extra_field_length")
					fileCommentLength := d.FieldU16("file_comment_length")
					diskNrStart := d.FieldU16("disk_number_where_file_starts")
					d.FieldU16("internal_file_attributes")
					d.FieldU32("external_file_attributes")
					localFileOffset := d.FieldU32("relative_offset_of_local_file_header")
					d.FieldUTF8("file_name", int(fileNameLength))
					d.FieldArray("extra_fields", func(d *decode.D) {
						d.LenFn(int64(extraFieldLength)*8, func(d *decode.D) {
							for !d.End() {
								d.FieldStruct("extra_field", func(d *decode.D) {
									d.FieldU16("header_id", d.MapUToScalar(headerIDMap), d.Hex)
									dataSize := d.FieldU16("data_size")
									d.FieldRawLen("data", int64(dataSize)*8)
								})
							}
						})
					})
					d.FieldUTF8("file_comment", int(fileCommentLength))

					if diskNrStart == diskNr {
						localFileOffsets = append(localFileOffsets, localFileOffset)
					}
				})
			}
		})
	})

	d.FieldArray("local_files", func(d *decode.D) {
		for _, o := range localFileOffsets {
			d.SeekAbs(int64(o) * 8)
			d.FieldStruct("local_file", func(d *decode.D) {
				var hasDataDescriptor bool
				d.FieldRawLen("signature", 4*8, d.ValidateBitBuf(localFileSignature))
				d.FieldU16("version_needed")
				d.FieldStruct("flags", func(d *decode.D) {
					// TODO: 16LE, should have some kind of native endian flag reader helper?
					d.FieldU1("unused0")
					d.FieldBool("strong_encryption")
					d.FieldBool("compressed_patched_data")
					d.FieldBool("enhanced_deflation")
					hasDataDescriptor = d.FieldBool("data_descriptor")
					d.FieldBool("compression0")
					d.FieldBool("compression1")
					d.FieldBool("encrypted")

					d.FieldU2("reserved0")
					d.FieldBool("mask_header_values")
					d.FieldBool("reserved1")
					d.FieldBool("language_encoding")
					d.FieldU3("unused1")
				})
				compressionMethod := d.FieldU16("compression_method", d.MapUToStrSym(compressionMethodMap))
				d.FieldStruct("last_modification_date", fieldMSDOSTime)
				d.FieldStruct("last_modification_time", fieldMSDOSDate)
				d.FieldU32("crc32_uncompressed", d.Hex)
				compressedSize := d.FieldU32("compressed_size")
				d.FieldU32("uncompressed_size")
				fileNameLength := d.FieldU16("file_name_length")
				extraFieldLength := d.FieldU16("extra_field_length")
				d.FieldUTF8("file_name", int(fileNameLength))
				d.FieldArray("extra_fields", func(d *decode.D) {
					d.LenFn(int64(extraFieldLength)*8, func(d *decode.D) {
						for !d.End() {
							d.FieldStruct("extra_field", func(d *decode.D) {
								d.FieldU16("header_id", d.MapUToScalar(headerIDMap), d.Hex)
								dataSize := d.FieldU16("data_size")
								d.FieldRawLen("data", int64(dataSize)*8)
							})
						}
					})
				})

				compressedLimit := int64(compressedSize) * 8
				if compressedLimit == 0 {
					compressedLimit = d.BitsLeft()
				}

				compressedStart := d.Pos()

				d.LenFn(compressedLimit, func(d *decode.D) {
					if compressionMethod == compressionMethodNone {
						d.FieldRawLen("uncompressed", int64(compressedSize)*8)
						return
					}

					var decompressR io.Reader
					compressedBB := d.BitBufRange(d.Pos(), d.BitsLeft())
					switch compressionMethod {
					case compressionMethodDeflated:
						// *bitio.Buffer implements io.ByteReader so hat deflate don't do own
						// buffering and might read more than needed messing up knowing compressed size
						decompressR = flate.NewReader(compressedBB)
					}

					if decompressR != nil {
						uncompressed := &bytes.Buffer{}
						if _, err := d.Copy(uncompressed, decompressR); err != nil {
							d.IOPanic(err)
						}
						uncompressedBB := bitio.NewBufferFromBytes(uncompressed.Bytes(), -1)
						dv, _, _ := d.FieldTryFormatBitBuf("uncompressed", uncompressedBB, probeFormat, nil)
						if dv == nil {
							d.FieldRootBitBuf("uncompressed", uncompressedBB)
						}

						// no compressed size, is a streaming zip, figure out size by checking what
						// position compressed buffer ended at
						if compressedSize == 0 {
							pos, err := compressedBB.Pos()
							if err != nil {
								d.IOPanic(err)
							}
							compressedSize = uint64(pos) / 8
						}
					}

					if compressedSize != 0 {
						d.FieldRawLen("compressed", int64(compressedSize)*8)
					}
				})

				d.SeekAbs(compressedStart + int64(compressedSize*8))

				if hasDataDescriptor {
					d.FieldStruct("data_indicator", func(d *decode.D) {
						if bytes.Equal(d.PeekBytes(4), dataIndicatorSignature) {
							d.FieldRawLen("signature", 4*8, d.ValidateBitBuf(dataIndicatorSignature))
						}
						d.FieldU32("crc32_uncompressed", d.Hex)
						d.FieldU32("compressed_size")
						d.FieldU32("uncompressed_size")
					})
				}
			})
		}
	})

	return nil
}
