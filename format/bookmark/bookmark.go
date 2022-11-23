package bplist

import (
	"embed"
	"fmt"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bookmark.jq bplist.md
var bookmarkFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.BOOKMARK,
		ProbeOrder:  format.ProbeOrderBinUnique,
		Description: "Apple BookmarkData",
		Groups:      []string{format.PROBE},
		DecodeFn:    bookmarkDecode,
		Functions:   []string{"torepr"},
	})
	interp.RegisterFS(bookmarkFS)
}

const (
	dataTypeString       = 0x0101
	dataTypeData         = 0x0201
	dataTypeNumber8      = 0x0301
	dataTypeNumber16     = 0x0302
	dataTypeNumber32     = 0x0303
	dataTypeNumber64     = 0x0304
	dataTypeNumber32F    = 0x0305
	dataTypeNumber64F    = 0x0306
	dataTypeDate         = 0x0400
	dataTypeBooleanFalse = 0x0500
	dataTypeBooleanTrue  = 0x0501
	dataTypeArray        = 0x0601
	dataTypeDictionary   = 0x0701
	dataTypeUUID         = 0x0801
	dataTypeURL          = 0x0901
	dataTypeRelativeURL  = 0x0902
)

var dataTypeMap = scalar.UToScalar{
	dataTypeString:       {Sym: "String", Description: "UTF-8 String"},
	dataTypeData:         {Sym: "Data", Description: "Raw bytes"},
	dataTypeNumber8:      {Sym: "Byte", Description: "(signed 8-bit) 1-byte number"},
	dataTypeNumber16:     {Sym: "Short", Description: "(signed 16-bit) 2-byte number"},
	dataTypeNumber32:     {Sym: "Int", Description: "(signed 32-bit) 4-byte number"},
	dataTypeNumber64:     {Sym: "Long", Description: "(signed 64-bit) 8-byte number"},
	dataTypeNumber32F:    {Sym: "Float", Description: "(32-bit float) IEEE single precision"},
	dataTypeNumber64F:    {Sym: "Double", Description: "(64-bit float) IEEE double precision"},
	dataTypeDate:         {Sym: "Date", Description: "Big-endian IEEE double precision seconds since 2001-01-01 00:00:00 UTC"},
	dataTypeBooleanFalse: {Sym: "BooleanFalse", Description: "(false)"},
	dataTypeBooleanTrue:  {Sym: "BooleanTrue", Description: "(true)"},
	dataTypeArray:        {Sym: "Array", Description: "Array of 4-byte offsets to data items"},
	dataTypeDictionary:   {Sym: "Dictionary", Description: "Array of pairs of 4-byte (key, value) data item offsets"},
	dataTypeUUID:         {Sym: "UUID", Description: "Raw bytes"},
	dataTypeURL:          {Sym: "URL", Description: "UTF-8 string"},
	dataTypeRelativeURL:  {Sym: "RelativeURL", Description: "4-byte offset to base URL, 4-byte offset to UTF-8 string"},
}

const (
	elementTypeTargetURL             = 0x1003
	elementTypeTargetPath            = 0x1004
	elementTypeTargetCNIDPath        = 0x1005
	elementTypeTargetFlags           = 0x1010
	elementTypeTargetFilename        = 0x1020
	elementTypeCNID                  = 0x1030
	elementTypeTargetCreationDate    = 0x1040
	elementTypeUnknown1              = 0x1054
	elementTypeUnknown2              = 0x1055
	elementTypeUnknown3              = 0x1056
	elementTypeUnknown4              = 0x1101
	elementTypeUnknown5              = 0x1102
	elementTypeTOCPath               = 0x2000
	elementTypeVolumePath            = 0x2002
	elementTypeVolumeURL             = 0x2005
	elementTypeVolumeName            = 0x2010
	elementTypeVolumeUUID            = 0x2011
	elementTypeVolumeSize            = 0x2012
	elementTypeVolumeCreationDate    = 0x2013
	elementTypeVolumeFlags           = 0x2020
	elementTypeVolumeIsRoot          = 0x2030
	elementTypeVolumeBookmark        = 0x2040
	elementTypeVolumeMountPointURL   = 0x2050
	elementTypeUnknown6              = 0x2070
	elementTypeContainingFolderIndex = 0xc001
	elementTypeCreatorUsername       = 0xc011
	elementTypeCreatorUID            = 0xc012
	elementTypeFileReferenceFlag     = 0xd001
	elementTypeCreationOptions       = 0xd010
	elementTypeURLLengthArray        = 0xe003
	elementTypeDisplayName           = 0xf017
	elementTypeIconData              = 0xf020
	elementTypeIconImageData         = 0xf021
	elementTypeTypeBindingInfo       = 0xf022
	elementTypeBookmarkCreationTime  = 0xf030
	elementTypeSandboxRWExtension    = 0xf080
	elementTypeSandboxROExtension    = 0xf081
)

var elementTypeMap = scalar.UToScalar{
	elementTypeTargetURL:             {Sym: "Target URL", Description: "A URL"},
	elementTypeTargetPath:            {Sym: "Target path", Description: "Array of individual path components"},
	elementTypeTargetCNIDPath:        {Sym: "Target CNID path", Description: "Array of CNIDs"},
	elementTypeTargetFlags:           {Sym: "Target flags", Description: "Data - see below"},
	elementTypeTargetFilename:        {Sym: "Target filename", Description: "String"},
	elementTypeCNID:                  {Sym: "Target CNID", Description: "4-byte integer"},
	elementTypeTargetCreationDate:    {Sym: "Target creation date", Description: "Date"},
	elementTypeUnknown1:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeUnknown2:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeUnknown3:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeUnknown4:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeUnknown5:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeTOCPath:               {Sym: "TOC path", Description: "Array - see below"},
	elementTypeVolumePath:            {Sym: "Volume path", Description: "Array of individual path components"},
	elementTypeVolumeURL:             {Sym: "Volume URL", Description: "URL of volume root"},
	elementTypeVolumeName:            {Sym: "Volume name", Description: "String"},
	elementTypeVolumeUUID:            {Sym: "Volume UUID", Description: "String (not a UUID!)"},
	elementTypeVolumeSize:            {Sym: "Volume size", Description: "8-byte integer"},
	elementTypeVolumeCreationDate:    {Sym: "Volume creation date", Description: "Date"},
	elementTypeVolumeFlags:           {Sym: "Volume flags", Description: "Data - see below"},
	elementTypeVolumeIsRoot:          {Sym: "Volume is root", Description: "True if the volume was the filesystem root"},
	elementTypeVolumeBookmark:        {Sym: "Volume bookmark", Description: "TOC identifier for disk image"},
	elementTypeVolumeMountPointURL:   {Sym: "Volume mount point", Description: "URL"},
	elementTypeUnknown6:              {Sym: "Unknown", Description: "Unknown"},
	elementTypeContainingFolderIndex: {Sym: "Containing folder index", Description: "Integer index of containing folder in target path array"},
	elementTypeCreatorUsername:       {Sym: "Creator username", Description: "Name of user that created bookmark"},
	elementTypeCreatorUID:            {Sym: "Creator UID", Description: "UID of user that created bookmark"},
	elementTypeFileReferenceFlag:     {Sym: "File reference flag", Description: "True if creating URL was a file reference URL"},
	elementTypeCreationOptions:       {Sym: "Creation options", Description: "Integer containing flags passed to CFURLCreateBookmarkData"},
	elementTypeURLLengthArray:        {Sym: "URL length array", Description: "Array of integers - see below"},
	elementTypeDisplayName:           {Sym: "Display name", Description: "String"},
	elementTypeIconData:              {Sym: "Icon data", Description: "icns format data"},
	elementTypeIconImageData:         {Sym: "Icon image", Description: "Data"},
	elementTypeTypeBindingInfo:       {Sym: "Type binding info", Description: "dnib byte array"},
	elementTypeBookmarkCreationTime:  {Sym: "Bookmark creation time", Description: "64-bit float seconds since January 1st 2001"},
	elementTypeSandboxRWExtension:    {Sym: "Sandbox RW extension", Description: "Looks like a hash with data and an access right"},
	elementTypeSandboxROExtension:    {Sym: "Sandbox RO extension", Description: "As above"},
}

var cocoaTimeEpochDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

type tocHeader struct {
	tocSize          uint64
	nextTOCOffset    uint64
	numEntries       uint64
	entryArrayOffset uint64
}

func decodeTOCHeader(d *decode.D, idx int) *tocHeader {
	hdr := new(tocHeader)

	d.FieldStruct(fmt.Sprintf("toc_header_%d", idx), func(d *decode.D) {
		hdr.tocSize = d.FieldU32("toc_size")
		d.FieldU32("magic", d.AssertU(0xfffffffe))
		d.FieldU32("identifier")
		hdr.nextTOCOffset = d.FieldU32("next_toc_offset")
		hdr.numEntries = d.FieldU32("num_entries_in_toc")
		hdr.entryArrayOffset = uint64(d.Pos())
	})

	return hdr
}

type tocEntry struct {
	key          uint64
	recordOffset uint64
}

func decodeTOCEntry(d *decode.D) *tocEntry {
	var entry *tocEntry

	entry.key = d.FieldU32("key")
	entry.recordOffset = d.FieldU32("offset_to_record")
	d.FieldU32("reserved")

	return entry
}

const (
	arrayEntrySize = 4
	dictEntrySize  = 4
)

func decodeRecord(d *decode.D, offset uint64) {
	d.SeekAbs(int64(offset), func(d *decode.D) {
		d.FieldStruct("record", func(d *decode.D) {
			n := int(d.FieldU32("length"))
			typ := d.FieldU32("type", dataTypeMap)
			switch typ {
			case dataTypeString:
				d.FieldUTF8("data", n)
			case dataTypeData:
				d.FieldRawLen("data", int64(n*8))
			case dataTypeNumber8:
				d.FieldS8("data")
			case dataTypeNumber16:
				d.FieldS16("data")
			case dataTypeNumber32:
				d.FieldS32("data")
			case dataTypeNumber64:
				d.FieldS64("data")
			case dataTypeNumber32F:
				d.FieldF32("data")
			case dataTypeNumber64F:
				d.FieldF64("data")
			case dataTypeDate:
				d.FieldF64BE("data")
			case dataTypeBooleanFalse:
			case dataTypeBooleanTrue:
			case dataTypeArray:
				d.FieldStructNArray("data", "element", int64(n/arrayEntrySize), func(d *decode.D) {
					offset := calcOffset(d.FieldU32("offset"))
					decodeRecord(d, offset)
				})
			case dataTypeDictionary:
				d.FieldStructNArray("data", "element", int64(n/dictEntrySize), func(d *decode.D) {
					keyOffset := calcOffset(d.FieldU32("key_offset"))
					d.FieldStruct("key", func(d *decode.D) {
						decodeRecord(d, keyOffset)
					})

					valueOffset := calcOffset(d.FieldU32("value_offset"))
					d.FieldStruct("value", func(d *decode.D) {
						decodeRecord(d, valueOffset)
					})
				})
			case dataTypeUUID:
				d.FieldRawLen("data", int64(n*8))
			case dataTypeURL:
				d.FieldUTF8("data", n)
			case dataTypeRelativeURL:
				baseOffset := d.FieldU32("base_url_offset")
				d.FieldStruct("base_url", func(d *decode.D) {
					decodeRecord(d, baseOffset)
				})

				suffixOffset := d.FieldU32("suffix_offset")
				d.FieldStruct("suffix", func(d *decode.D) {
					decodeRecord(d, suffixOffset)
				})
			}
		})
	})
}

const reservedSize = 32
const headerEnd = 48
const headerEndBitPos = headerEnd * 8

// all offsets are calculated relative to the end of the bookmark header
func calcOffset(i uint64) uint64 { return 8 * (i + headerEnd) }

func bookmarkDecode(d *decode.D, _ any) any {

	// all fields are little-endian with the exception of the Date datatype.
	d.Endian = decode.LittleEndian

	// decode bookmarkdata header, one at the top of each "file",
	// although these may be nested inside of binary plists
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 4, d.AssertStr("book", "alis"))
		d.FieldU32LE("total_size")
		d.FieldU32("unknown")
		d.FieldU32("header_size", d.AssertU(48))
		d.FieldRawLen("reserved", reservedSize*8)
	})

	tocOffset := calcOffset(d.FieldU32("first_toc_offset"))

	var tocHeaders []*tocHeader

	for i := 0; tocOffset != headerEndBitPos; i++ {
		// seek to the TOC, and decode the header and entries
		// for this TOC instance. SeekAbs restores our offset each time.
		d.SeekAbs(int64(tocOffset), func(d *decode.D) {

			tocHdr := decodeTOCHeader(d, i)
			// store the toc header. we're going to decode the entries in one
			// big array once we have decoded all toc's
			tocHeaders = append(tocHeaders, tocHdr)
			// save the next toc_offset value. 0 indicates that we have reached
			// the last TOC instance.
			tocOffset = calcOffset(tocHdr.nextTOCOffset)

		})

		j := 0

		// now that we've collected all toc headers, iterate through each one's
		// entries and decode associated records.
		d.FieldArrayLoop("bookmark_entries",
			func() bool { return j < len(tocHeaders) },
			func(d *decode.D) {

				tocHdr := tocHeaders[j]
				j++

				d.SeekAbs(int64(tocHdr.entryArrayOffset), func(d *decode.D) {
					for k := uint64(0); k < tocHdr.numEntries; k++ {
						entry := new(tocEntry)

						d.FieldStruct("entry", func(d *decode.D) {
							entry.key = d.FieldU32("key", elementTypeMap)

							entry.recordOffset = calcOffset(d.FieldU32("offset_to_record"))

							d.FieldU32("reserved")

							decodeRecord(d, entry.recordOffset)
						})
					}
				})
			})
	}

	return nil
}
