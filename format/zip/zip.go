package zip

// https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
// https://opensource.apple.com/source/zip/zip-6/unzip/unzip/proginfo/extra.fld

import (
	"bytes"
	"compress/flate"
	"embed"
	"io"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed zip.md
var zipFS embed.FS

var probeGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Zip,
		&decode.Format{
			Description: "ZIP archive",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    zipDecode,
			DefaultInArg: format.Zip_In{
				Uncompress: true,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Probe}, Out: &probeGroup},
			},
		})
	interp.RegisterFS(zipFS)
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

var compressionMethodMap = scalar.UintMapSymStr{
	compressionMethodNone:                      "none",
	compressionMethodShrunk:                    "shrunk",
	compressionMethodReducedCompressionFactor1: "reduced_compression_factor1",
	compressionMethodReducedCompressionFactor2: "reduced_compression_factor2",
	compressionMethodReducedCompressionFactor3: "reduced_compression_factor3",
	compressionMethodReducedCompressionFactor4: "reduced_compression_factor4",
	compressionMethodImploded:                  "imploded",
	compressionMethodDeflated:                  "deflated",
	compressionMethodEnhancedDeflated:          "enhanced_deflated",
	compressionMethodPKWareDCLImploded:         "pk_ware_dcl_imploded",
	compressionMethodBzip2:                     "bzip2",
	compressionMethodLZMA:                      "lzma",
	compressionMethodIBMTERSE:                  "ibmterse",
	compressionMethodIBMLZ77z:                  "ibmlz77z",
	compressionMethodPPMd:                      "pp_md",
}

var (
	centralDirectorySignature              = []byte("PK\x01\x02")
	endOfCentralDirectoryRecordSignature   = []byte("PK\x05\x06")
	endOfCentralDirectoryRecordSignatureN  = 0x06054b50
	endOfCentralDirectoryRecord64Signature = []byte("PK\x06\x06")
	endOfCentralDirectoryLocatorSignature  = []byte("PK\x06\x07")
	endOfCentralDirectoryLocatorSignatureN = 0x07064b50
	localFileSignature                     = []byte("PK\x03\x04")
	dataIndicatorSignature                 = []byte("PK\x07\x08")
)

const (
	headerTagZip64ExtendedInformation = 0x001
	headerTagExtendedTimestamp        = 0x5455
)

var headerTagMap = scalar.UintMapDescription{
	headerTagZip64ExtendedInformation: "ZIP64 extended information extra field",
	0x0007:                            "AV Info",
	0x0009:                            "OS/2 extended attributes",
	0x000a:                            "NTFS (Win9x/WinNT FileTimes)",
	0x000c:                            "OpenVMS",
	0x000d:                            "Unix",
	0x000f:                            "Patch Descriptor",
	0x0014:                            "PKCS#7 Store for X.509 Certificates",
	0x0015:                            "X.509 Certificate ID and Signature for individual file",
	0x0016:                            "X.509 Certificate ID for Central Directory",
	0x0065:                            "IBM S/390 attributes - uncompressed",
	0x0066:                            "IBM S/390 attributes - compressed",
	0x07c8:                            "Info-ZIP Macintosh (old, J. Lee)",
	0x2605:                            "ZipIt Macintosh (first version)",
	0x2705:                            "ZipIt Macintosh v 1.3.5 and newer (w/o full filename)",
	0x334d:                            "Info-ZIP Macintosh (new, D. Haase's 'Mac3' field )",
	0x4154:                            "Tandem NSK",
	0x4341:                            "Acorn/SparkFS (David Pilling)",
	0x4453:                            "Windows NT security descriptor (binary ACL)",
	0x4704:                            "VM/CMS",
	0x470f:                            "MVS",
	// "inofficial" in original table
	//nolint:misspell
	0x4854:                     "Theos, old inofficial port",
	0x4b46:                     "FWKCS MD5 (see below)",
	0x4c41:                     "OS/2 access control list (text ACL)",
	0x4d49:                     "Info-ZIP OpenVMS (obsolete)",
	0x4d63:                     "Macintosh SmartZIP, by Macro Bambini",
	0x4f4c:                     "Xceed original location extra field",
	0x5356:                     "AOS/VS (binary ACL)",
	headerTagExtendedTimestamp: "extended timestamp",
	0x5855:                     "Info-ZIP Unix (original; also OS/2, NT, etc.)",
	0x554e:                     "Xceed unicode extra field",
	0x6542:                     "BeOS (BeBox, PowerMac, etc.)",
	0x6854:                     "Theos",
	0x756e:                     "ASi Unix",
	0x7855:                     "Info-ZIP Unix (new)",
	0x7875:                     "UNIX UID/GID",
	0xfb4a:                     "SMS/QDOS",
}

// "MS-DOS uses year values relative to 1980 and 2 second precision."
// https://learn.microsoft.com/en-gb/windows/win32/api/winbase/nf-winbase-dosdatetimetofiletime?redirectedfrom=MSDN
// https://formats.kaitai.io/dos_datetime/
// Note all of this is a mess because time/date is stored in bit ranges inside 16 LE numbers
// TODO: maybe can be cleaned up if bit-endian decoding is added?
func fieldMSDOSTime(d *decode.D) (int, int, int) {
	fatTime := d.FieldU16("fat_time", scalar.UintHex)

	// second/2 b5
	// minute b6
	// hour b5
	second := (fatTime >> 0) & 0b1_1111
	minute := (fatTime >> 5) & 0b11_1111
	hour := (fatTime >> (5 + 6)) & 0b1_1111
	d.FieldValueUint("second", second, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		s.Sym = s.Actual * 2
		return s, nil
	}))
	d.FieldValueUint("minute", minute)
	d.FieldValueUint("hour", hour)

	return int(second), int(minute), int(hour)
}

func fieldMSDOSDate(d *decode.D) (int, int, int) {
	fatDate := d.FieldU16("fat_date", scalar.UintHex)

	// day b5
	// month b4
	// day b7
	day := (fatDate >> 0) & 0b1_1111
	month := (fatDate >> 5) & 0b1111
	year := (fatDate >> (5 + 4)) & 0b111_1111
	d.FieldValueUint("day", day)
	d.FieldValueUint("month", month)
	d.FieldValueUint("year", year, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		s.Sym = s.Actual + 1980
		return s, nil
	}))

	return int(day), int(month), int(year)
}

// time.RFC3339 but no timezone
const rfc3339Local = "2006-01-02T15:04:05"

func fieldTimeDate(d *decode.D) {
	var second, minute, hour int
	var day, month, year int
	second, minute, hour = fieldMSDOSTime(d)
	day, month, year = fieldMSDOSDate(d)
	t := time.Date(1980+year, time.Month(month), day, hour, minute, second*2, 0, time.UTC)
	d.FieldValueUint("unix_guess", uint64(t.Unix()),
		scalar.UintActualUnixTimeDescription(time.Second, rfc3339Local))
}

func fieldExtendedTimestamp(d *decode.D) {
	modificationTimePresent := false
	accessTimePresent := false
	creationTimePresent := false
	d.FieldStruct("flags", func(d *decode.D) {
		d.FieldU5("unused")
		creationTimePresent = d.FieldBool("creation_time_present")
		accessTimePresent = d.FieldBool("access_time_present")
		modificationTimePresent = d.FieldBool("modification_time_present")
	})
	// Spec says this but seem like flags and size is not in sync sometimes?
	// ex: flags is 0x03 but size is 5
	// > TSize should equal (1 + 4*(number of set bits in Flags))
	if modificationTimePresent && !d.End() {
		d.FieldU32("modification_time", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
	}
	if accessTimePresent && !d.End() {
		d.FieldU32("access_time", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
	}
	if creationTimePresent && !d.End() {
		d.FieldU32("creation_time", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
	}
}

type zip64ExtendedInformation struct {
	uncompressedSize                 uint64
	uncompressedSizePresent          bool
	compressedSize                   uint64
	compressedSizePresent            bool
	localFileOffset                  uint64
	localFileOffsetPresent           bool
	diskNumberWhereFileStarts        uint64
	diskNumberWhereFileStartsPresent bool
}

func fieldTagZip64ExtendedInformation(d *decode.D) zip64ExtendedInformation {
	zi := zip64ExtendedInformation{}

	zi.uncompressedSize = d.FieldU64("uncompressed_size")
	zi.uncompressedSizePresent = true
	// TODO: spec says these should be here but real zip64 seems to not have them? optional?
	if !d.End() {
		zi.compressedSize = d.FieldU64("compressed_size")
		zi.compressedSizePresent = true
	}
	if !d.End() {
		zi.localFileOffset = d.FieldU64("relative_offset_of_local_file_header")
		zi.localFileOffsetPresent = true
	}
	if !d.End() {
		zi.diskNumberWhereFileStarts = d.FieldU32("disk_number_where_file_starts")
		zi.diskNumberWhereFileStartsPresent = true
	}
	return zi
}

type extraFields struct {
	zip64ExtendedInformation        zip64ExtendedInformation
	zip64ExtendedInformationPresent bool
}

func fieldsExtraFields(d *decode.D) extraFields {
	ef := extraFields{}

	for !d.End() {
		d.FieldStruct("extra_field", func(d *decode.D) {
			tag := d.FieldU16("tag", headerTagMap, scalar.UintHex)
			size := d.FieldU16("size")
			d.FramedFn(int64(size)*8, func(d *decode.D) {
				switch tag {
				case headerTagZip64ExtendedInformation:
					ef.zip64ExtendedInformation = fieldTagZip64ExtendedInformation(d)
					ef.zip64ExtendedInformationPresent = true
				case headerTagExtendedTimestamp:
					fieldExtendedTimestamp(d)
				default:
					d.FieldRawLen("data", int64(size)*8)
				}
			})
		})
	}
	return ef
}

func zipDecode(d *decode.D) any {
	var zi format.Zip_In
	d.ArgAs(&zi)

	d.Endian = decode.LittleEndian

	// zip files are parsed from end
	d.SeekAbs(d.Len())

	// TODO: better EOCD probe
	p, _, err := d.TryPeekFind(32, -8, 128*8, func(v uint64) bool {
		return v == uint64(endOfCentralDirectoryRecordSignatureN)
	})
	if err != nil {
		d.Fatalf("can't find end of central directory")
	}
	d.SeekRel(p)

	var offsetCD uint64
	var sizeCD uint64
	var diskNr uint64

	d.FieldStruct("end_of_central_directory_record", func(d *decode.D) {
		d.FieldRawLen("signature", 4*8, d.AssertBitBuf(endOfCentralDirectoryRecordSignature))
		diskNr = d.FieldU16("disk_nr")
		d.FieldU16("central_directory_start_disk_nr")
		d.FieldU16("nr_of_central_directory_records_on_disk")
		d.FieldU16("nr_of_central_directory_records")
		sizeCD = d.FieldU32("size_of_central_directory")
		offsetCD = d.FieldU32("offset_of_start_of_central_directory")
		commentLength := d.FieldU16("comment_length")
		d.FieldUTF8("comment", int(commentLength))
	})

	// is there a zip64 end of central directory locator?
	p, _, err = d.TryPeekFind(32, -8, 128*8, func(v uint64) bool {
		return v == uint64(endOfCentralDirectoryLocatorSignatureN)
	})
	if err == nil && p != -1 {
		d.SeekRel(p)

		var offsetEOCD uint64
		d.FieldStruct("end_of_central_directory_locator", func(d *decode.D) {
			d.FieldRawLen("signature", 4*8, d.AssertBitBuf(endOfCentralDirectoryLocatorSignature))
			diskNr = d.FieldU32("disk_nr")
			offsetEOCD = d.FieldU64("offset_of_end_of_central_directory_record")
			diskNr = d.FieldU32("total_disk_nr")
		})

		d.SeekAbs(int64(offsetEOCD) * 8)
		d.FieldStruct("end_of_central_directory_record_zip64", func(d *decode.D) {
			d.FieldRawLen("signature", 4*8, d.AssertBitBuf(endOfCentralDirectoryRecord64Signature))
			sizeEOCD := d.FieldU64("size_of_end_of_central_directory")
			d.FieldU16("version_made_by")
			d.FieldU16("version_needed_to_extract")
			diskNr = d.FieldU32("disk_nr")
			d.FieldU32("central_directory_start_disk_nr")
			d.FieldU64("nr_of_central_directory_records_on_disk")
			d.FieldU64("nr_of_central_directory_records")
			sizeCD = d.FieldU64("size_of_central_directory")
			offsetCD = d.FieldU64("offset_of_start_of_central_directory")
			const sizeOfFixedFields = 44
			d.FramedFn(int64(sizeEOCD-sizeOfFixedFields)*8, func(d *decode.D) {
				d.FieldArray("extensible_data", func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("extensible_data", func(d *decode.D) {
							d.FieldU16("tag", headerTagMap, scalar.UintHex)
							dataSize := d.FieldU32("size")
							d.FieldRawLen("data", int64(dataSize)*8)
						})
					}
				})
			})
		})
	}

	var localFileOffsets []uint64

	d.SeekAbs(int64(offsetCD) * 8)
	d.FieldArray("central_directories", func(d *decode.D) {
		d.FramedFn(int64(sizeCD)*8, func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("central_directory", func(d *decode.D) {
					d.FieldRawLen("signature", 4*8, d.AssertBitBuf(centralDirectorySignature))
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
					d.FieldU16("compression_method", compressionMethodMap)
					d.FieldStruct("last_modification", fieldTimeDate)
					d.FieldU32("crc32_uncompressed", scalar.UintHex)
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
						d.FramedFn(int64(extraFieldLength)*8, func(d *decode.D) {
							ef := fieldsExtraFields(d)
							if ef.zip64ExtendedInformationPresent &&
								ef.zip64ExtendedInformation.localFileOffsetPresent {
								localFileOffset = ef.zip64ExtendedInformation.localFileOffset
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
				d.FieldRawLen("signature", 4*8, d.AssertBitBuf(localFileSignature))
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
				compressionMethod := d.FieldU16("compression_method", compressionMethodMap)
				d.FieldStruct("last_modification", fieldTimeDate)
				d.FieldU32("crc32_uncompressed", scalar.UintHex)
				compressedSizeBytes := d.FieldU32("compressed_size")
				d.FieldU32("uncompressed_size")
				fileNameLength := d.FieldU16("file_name_length")
				extraFieldLength := d.FieldU16("extra_field_length")
				d.FieldUTF8("file_name", int(fileNameLength))
				d.FieldArray("extra_fields", func(d *decode.D) {
					d.FramedFn(int64(extraFieldLength)*8, func(d *decode.D) {
						ef := fieldsExtraFields(d)
						if ef.zip64ExtendedInformationPresent &&
							ef.zip64ExtendedInformation.compressedSizePresent {
							compressedSizeBytes = ef.zip64ExtendedInformation.compressedSize
						}
					})
				})
				compressedSize := int64(compressedSizeBytes) * 8
				compressedStart := d.Pos()

				compressedLimit := compressedSize
				if compressedLimit == 0 {
					compressedLimit = d.BitsLeft()
				}

				if compressionMethod == compressionMethodNone {
					d.FieldFormatOrRawLen("uncompressed", compressedSize, &probeGroup, format.Probe_In{})
				} else {
					var rFn func(r io.Reader) io.Reader
					if zi.Uncompress {
						switch compressionMethod {
						case compressionMethodDeflated:
							// bitio.NewIOReadSeeker implements io.ByteReader so that deflate don't do own
							// buffering and might read more than needed messing up knowing compressed size
							rFn = func(r io.Reader) io.Reader { return flate.NewReader(r) }
						}
					}

					if rFn != nil {
						readCompressedSize, uncompressedBR, dv, _, _ :=
							d.TryFieldReaderRangeFormat("uncompressed", d.Pos(), compressedLimit, rFn, &probeGroup, format.Probe_In{})
						if dv == nil && uncompressedBR != nil {
							d.FieldRootBitBuf("uncompressed", uncompressedBR)
						}
						if compressedSize == 0 {
							compressedSize = readCompressedSize
						}
						d.FieldRawLen("compressed", compressedSize)

					} else {
						if compressedSize != 0 {
							d.FieldRawLen("compressed", compressedSize)
						}
					}
				}

				d.SeekAbs(compressedStart + compressedSize)

				if hasDataDescriptor {
					d.FieldStruct("data_indicator", func(d *decode.D) {
						if bytes.Equal(d.PeekBytes(4), dataIndicatorSignature) {
							d.FieldRawLen("signature", 4*8, d.AssertBitBuf(dataIndicatorSignature))
						}
						d.FieldU32("crc32_uncompressed", scalar.UintHex)
						d.FieldU32("compressed_size")
						d.FieldU32("uncompressed_size")
					})
				}
			})
		}
	})

	return nil
}
