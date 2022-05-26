package zip

// https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
// https://opensource.apple.com/source/zip/zip-6/unzip/unzip/proginfo/extra.fld

import (
	"bytes"
	"compress/flate"
	"embed"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed zip.jq
var zipFS embed.FS

var probeFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ZIP,
		Description: "ZIP archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    zipDecode,
		DecodeInArg: format.ZipIn{
			Uncompress: true,
		},
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
		Files:     zipFS,
		Functions: []string{"_help"},
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

var compressionMethodMap = scalar.UToSymStr{
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
	headerIDZip64ExtendedInformation = 0x001
)

var headerIDMap = scalar.UToDescription{
	headerIDZip64ExtendedInformation: "ZIP64 extended information extra field",
	0x0007:                           "AV Info",
	0x0009:                           "OS/2 extended attributes",
	0x000a:                           "NTFS (Win9x/WinNT FileTimes)",
	0x000c:                           "OpenVMS",
	0x000d:                           "Unix",
	0x000f:                           "Patch Descriptor",
	0x0014:                           "PKCS#7 Store for X.509 Certificates",
	0x0015:                           "X.509 Certificate ID and Signature for individual file",
	0x0016:                           "X.509 Certificate ID for Central Directory",
	0x0065:                           "IBM S/390 attributes - uncompressed",
	0x0066:                           "IBM S/390 attributes - compressed",
	0x07c8:                           "Info-ZIP Macintosh (old, J. Lee)",
	0x2605:                           "ZipIt Macintosh (first version)",
	0x2705:                           "ZipIt Macintosh v 1.3.5 and newer (w/o full filename)",
	0x334d:                           "Info-ZIP Macintosh (new, D. Haase's 'Mac3' field )",
	0x4154:                           "Tandem NSK",
	0x4341:                           "Acorn/SparkFS (David Pilling)",
	0x4453:                           "Windows NT security descriptor (binary ACL)",
	0x4704:                           "VM/CMS",
	0x470f:                           "MVS",
	// "inofficial" in original table
	//nolint:misspell
	0x4854: "Theos, old inofficial port",
	0x4b46: "FWKCS MD5 (see below)",
	0x4c41: "OS/2 access control list (text ACL)",
	0x4d49: "Info-ZIP OpenVMS (obsolete)",
	0x4d63: "Macintosh SmartZIP, by Macro Bambini",
	0x4f4c: "Xceed original location extra field",
	0x5356: "AOS/VS (binary ACL)",
	0x5455: "extended timestamp",
	0x5855: "Info-ZIP Unix (original; also OS/2, NT, etc.)",
	0x554e: "Xceed unicode extra field",
	0x6542: "BeOS (BeBox, PowerMac, etc.)",
	0x6854: "Theos",
	0x756e: "ASi Unix",
	0x7855: "Info-ZIP Unix (new)",
	0x7875: "UNIX UID/GID",
	0xfb4a: "SMS/QDOS",
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

func zipDecode(d *decode.D, in any) any {
	zi, _ := in.(format.ZipIn)

	// TODO: just decode instead?
	if !bytes.Equal(d.PeekBytes(4), []byte("PK\x03\x04")) {
		d.Errorf("expected PK header")
	}

	d.Endian = decode.LittleEndian

	d.SeekAbs(d.Len())

	// TODO: better EOCD probe
	p, _, err := d.TryPeekFind(32, -8, -10000, func(v uint64) bool {
		return v == uint64(endOfCentralDirectoryRecordSignatureN)
	})
	if err != nil {
		d.Fatalf("can't find end of central directory")
	}
	d.SeekAbs(d.Len() + p)

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

	// there is a end of central directory locator, is zip64
	if offsetCD == 0xff_ff_ff_ff {
		p, _, err := d.TryPeekFind(32, -8, -10000, func(v uint64) bool {
			return v == uint64(endOfCentralDirectoryLocatorSignatureN)
		})
		if err != nil {
			d.Fatalf("can't find zip64 end of central directory")
		}
		d.SeekAbs(d.Len() + p)

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
				for !d.End() {
					d.FieldStruct("extra_field", func(d *decode.D) {
						d.FieldU16("header_id", headerIDMap, scalar.ActualHex)
						dataSize := d.FieldU32("data_size")
						d.FieldRawLen("data", int64(dataSize)*8)
					})
				}
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
					d.FieldStruct("last_modification_date", fieldMSDOSTime)
					d.FieldStruct("last_modification_time", fieldMSDOSDate)
					d.FieldU32("crc32_uncompressed", scalar.ActualHex)
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
							for !d.End() {
								d.FieldStruct("extra_field", func(d *decode.D) {
									headerID := d.FieldU16("header_id", headerIDMap, scalar.ActualHex)
									dataSize := d.FieldU16("data_size")
									d.FramedFn(int64(dataSize)*8, func(d *decode.D) {
										switch headerID {
										case headerIDZip64ExtendedInformation:
											d.FieldU64("uncompressed_size")
											// TODO: spec says these should be here but real zip64 seems to not have them? optional?
											if !d.End() {
												d.FieldU64("compressed_size")
											}
											if !d.End() {
												localFileOffset = d.FieldU64("relative_offset_of_local_file_header")
											}
											if !d.End() {
												d.FieldU32("disk_number_where_file_starts")
											}
										default:
											d.FieldRawLen("data", int64(dataSize)*8)
										}
									})
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
				d.FieldStruct("last_modification_date", fieldMSDOSTime)
				d.FieldStruct("last_modification_time", fieldMSDOSDate)
				d.FieldU32("crc32_uncompressed", scalar.ActualHex)
				compressedSizeBytes := d.FieldU32("compressed_size")
				d.FieldU32("uncompressed_size")
				fileNameLength := d.FieldU16("file_name_length")
				extraFieldLength := d.FieldU16("extra_field_length")
				d.FieldUTF8("file_name", int(fileNameLength))
				d.FieldArray("extra_fields", func(d *decode.D) {
					d.FramedFn(int64(extraFieldLength)*8, func(d *decode.D) {
						for !d.End() {
							d.FieldStruct("extra_field", func(d *decode.D) {
								headerID := d.FieldU16("header_id", headerIDMap, scalar.ActualHex)
								dataSize := d.FieldU16("data_size")
								d.FramedFn(int64(dataSize)*8, func(d *decode.D) {
									switch headerID {
									case headerIDZip64ExtendedInformation:
										d.FieldU64("uncompressed_size")
										// TODO: spec says these should be here but real zip64 seems to not have them? optional?
										if !d.End() {
											compressedSizeBytes = d.FieldU64("compressed_size")
										}
									default:
										d.FieldRawLen("data", int64(dataSize)*8)
									}
								})
							})
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
					d.FieldFormatOrRawLen("uncompressed", compressedSize, probeFormat, nil)
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
						readCompressedSize, uncompressedBR, dv, _, _ := d.TryFieldReaderRangeFormat("uncompressed", d.Pos(), compressedLimit, rFn, probeFormat, nil)
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
						d.FieldU32("crc32_uncompressed", scalar.ActualHex)
						d.FieldU32("compressed_size")
						d.FieldU32("uncompressed_size")
					})
				}
			})
		}
	})

	return nil
}
