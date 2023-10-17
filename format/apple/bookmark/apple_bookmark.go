package bookmarkdata

import (
	"embed"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/apple"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed apple_bookmark.jq apple_bookmark.md
var bookmarkFS embed.FS

func init() {
	interp.RegisterFormat(format.Apple_Bookmark,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinUnique,
			Description: "Apple BookmarkData",
			Groups:      []*decode.Group{format.Probe},
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

var dataTypeMap = scalar.UintMap{
	dataTypeString:       {Sym: "string", Description: "UTF-8 String"},
	dataTypeData:         {Sym: "data", Description: "Raw bytes"},
	dataTypeNumber8:      {Sym: "byte", Description: "(signed 8-bit) 1-byte number"},
	dataTypeNumber16:     {Sym: "short", Description: "(signed 16-bit) 2-byte number"},
	dataTypeNumber32:     {Sym: "int", Description: "(signed 32-bit) 4-byte number"},
	dataTypeNumber64:     {Sym: "long", Description: "(signed 64-bit) 8-byte number"},
	dataTypeNumber32F:    {Sym: "float", Description: "(32-bit float) IEEE single precision"},
	dataTypeNumber64F:    {Sym: "double", Description: "(64-bit float) IEEE double precision"},
	dataTypeDate:         {Sym: "date", Description: "Big-endian IEEE double precision seconds since 2001-01-01 00:00:00 UTC"},
	dataTypeBooleanFalse: {Sym: "boolean_false", Description: "False"},
	dataTypeBooleanTrue:  {Sym: "boolean_true", Description: "True"},
	dataTypeArray:        {Sym: "array", Description: "Array of 4-byte offsets to data items"},
	dataTypeDictionary:   {Sym: "dictionary", Description: "Array of pairs of 4-byte (key, value) data item offsets"},
	dataTypeUUID:         {Sym: "uuid", Description: "Raw bytes"},
	dataTypeURL:          {Sym: "url", Description: "UTF-8 string"},
	dataTypeRelativeURL:  {Sym: "relative_url", Description: "4-byte offset to base URL, 4-byte offset to UTF-8 string"},
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

var elementTypeMap = scalar.UintMap{
	elementTypeTargetURL:             {Sym: "target_url", Description: "A URL"},
	elementTypeTargetPath:            {Sym: "target_path", Description: "Array of individual path components"},
	elementTypeTargetCNIDPath:        {Sym: "target_cnid_path", Description: "Array of CNIDs"},
	elementTypeTargetFlags:           {Sym: "target_flags", Description: "flag bitfield"},
	elementTypeTargetFilename:        {Sym: "target_filename", Description: "String"},
	elementTypeCNID:                  {Sym: "target_cnid", Description: "4-byte integer"},
	elementTypeTargetCreationDate:    {Sym: "target_creation_date", Description: "Date"},
	elementTypeUnknown1:              {Sym: "unknown1", Description: "Unknown"},
	elementTypeUnknown2:              {Sym: "unknown2", Description: "Unknown"},
	elementTypeUnknown3:              {Sym: "unknown3", Description: "Unknown"},
	elementTypeUnknown4:              {Sym: "unknown4", Description: "Unknown"},
	elementTypeUnknown5:              {Sym: "unknown5", Description: "Unknown"},
	elementTypeTOCPath:               {Sym: "toc_path", Description: "Array - see below"},
	elementTypeVolumePath:            {Sym: "volume_path", Description: "Array of individual path components"},
	elementTypeVolumeURL:             {Sym: "volume_url", Description: "URL of volume root"},
	elementTypeVolumeName:            {Sym: "volume_name", Description: "String"},
	elementTypeVolumeUUID:            {Sym: "volume_uuid", Description: "String UUID"},
	elementTypeVolumeSize:            {Sym: "volume_size", Description: "8-byte integer"},
	elementTypeVolumeCreationDate:    {Sym: "volume_creation_date", Description: "Date"},
	elementTypeVolumeFlags:           {Sym: "volume_flags", Description: "flag bitfield"},
	elementTypeVolumeIsRoot:          {Sym: "volume_is_root", Description: "True if the volume was the filesystem root"},
	elementTypeVolumeBookmark:        {Sym: "volume_bookmark", Description: "TOC identifier for disk image"},
	elementTypeVolumeMountPointURL:   {Sym: "volume_mount_point", Description: "URL"},
	elementTypeUnknown6:              {Sym: "unknown6", Description: "Unknown"},
	elementTypeContainingFolderIndex: {Sym: "containing_folder_index", Description: "Integer index of containing folder in target path array"},
	elementTypeCreatorUsername:       {Sym: "creator_username", Description: "Name of user that created bookmark"},
	elementTypeCreatorUID:            {Sym: "creator_uid", Description: "UID of user that created bookmark"},
	elementTypeFileReferenceFlag:     {Sym: "file_reference_flag", Description: "True if creating URL was a file reference URL"},
	elementTypeCreationOptions:       {Sym: "creation_options", Description: "Integer containing flags passed to CFURLCreateBookmarkData"},
	elementTypeURLLengthArray:        {Sym: "url_length_array", Description: "Array of integers"},
	elementTypeDisplayName:           {Sym: "display_name", Description: "String"},
	elementTypeIconData:              {Sym: "icon_data", Description: "icns format data"},
	elementTypeIconImageData:         {Sym: "icon_image", Description: "Data"},
	elementTypeTypeBindingInfo:       {Sym: "type_binding_info", Description: "dnib byte array"},
	elementTypeBookmarkCreationTime:  {Sym: "bookmark_creation_time", Description: "64-bit float seconds since January 1st 2001"},
	elementTypeSandboxRWExtension:    {Sym: "sandbox_rw_extension", Description: "Looks like a hash with data and an access right"},
	elementTypeSandboxROExtension:    {Sym: "sandbox_ro_extension", Description: "Looks like a hash with data and an access right"},
}

const dataObjectLen = 24

func decodeFlagDataObject(d *decode.D, flagFn func(d *decode.D)) {
	d.FieldStruct("record", func(d *decode.D) {
		d.FieldU32("length", d.UintAssert(dataObjectLen))
		d.FieldU32("raw_type", dataTypeMap, d.UintAssert(dataTypeData))
		d.FieldValueStr("type", "flag_data")
		d.FieldStruct("property_flags", flagFn)
		d.FieldStruct("enabled_property_flags", flagFn)
		d.FieldRawLen("reserved", 64)
	})
}

func decodeTgtPropertyFlagBits(d *decode.D) {
	start := d.Pos()
	d.FieldBool("is_hidden")
	d.FieldBool("is_user_immutable")
	d.FieldBool("is_system_immutable")
	d.FieldBool("is_package")
	d.FieldBool("is_volume")
	d.FieldBool("is_symbolic_link")
	d.FieldBool("is_directory")
	d.FieldBool("is_regular_file")

	d.FieldBool("is_alias_file")
	d.FieldBool("is_executable")
	d.FieldBool("is_writeable")
	d.FieldBool("is_readable")
	d.FieldBool("can_set_hidden_extension")
	d.FieldBool("is_compressed")
	d.FieldBool("is_application")
	d.FieldBool("has_hidden_extension")

	d.FieldRawLen("reserved_bits_0", 7)
	d.FieldBool("is_mount_trigger")

	d.FieldRawLen("reserved", 64-(d.Pos()-start))
}

func decodeVolPropertyFlagBits(d *decode.D) {
	d.FieldBool("is_internal")
	d.FieldBool("is_removable")
	d.FieldBool("is_ejectable")
	d.FieldBool("is_quarantined")
	d.FieldBool("is_read_only")
	d.FieldBool("dont_browse")
	d.FieldBool("is_automount")
	d.FieldBool("is_local")

	d.FieldBool("is_dvd")
	d.FieldBool("is_cd")
	d.FieldBool("is_idisk")
	d.FieldBool("is_ipod")
	d.FieldBool("is_local_idisk_mirror")
	d.FieldBool("is_file_vault")
	d.FieldBool("is_disk_image")
	d.FieldBool("is_external")

	d.FieldRawLen("reserved_0", 7)
	d.FieldBool("is_device_file_system")

	d.FieldRawLen("reserved_1", 8)

	d.FieldBool("supports_read_dir_attr")
	d.FieldBool("supports_copy_file")
	d.FieldBool("supports_deny_modes")
	d.FieldBool("supports_symbolic_links")
	d.FieldBool("reserved_2")
	d.FieldBool("supports_exchange")
	d.FieldBool("supports_search_fs")
	d.FieldBool("supports_persistent_ids")

	d.FieldBool("supports_extended_security")
	d.FieldBool("has_no_root_directory_times")
	d.FieldBool("supports_flock")
	d.FieldBool("supports_case_preserved_names")
	d.FieldBool("supports_case_sensitive_names")
	d.FieldBool("supports_fast_stat_fs")
	d.FieldBool("supports_rename")
	d.FieldBool("supports_journaling")

	d.FieldBool("supports_zero_runs")
	d.FieldBool("supports_sparse_files")
	d.FieldBool("is_journaling")
	d.FieldBool("reserved_3")
	d.FieldBool("supports_path_from_id")
	d.FieldBool("supports_mandatory_byte_range_locks")
	d.FieldBool("supports_hard_links")
	d.FieldBool("supports_2_tb_file_size")

	d.FieldRawLen("reserved_4", 3)
	d.FieldBool("has64_bit_object_ids")
	d.FieldBool("supports_decmp_fs_compression")
	d.FieldBool("supports_hidden_files")
	d.FieldBool("supports_remote_events")
	d.FieldBool("supports_volume_sizes")
}

var cocoaTimeEpochDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

type tocHeader struct {
	nextTOCOffset    uint64
	numEntries       uint64
	entryArrayOffset int64
}

func (hdr *tocHeader) decodeEntries(d *decode.D) {
	for k := uint64(0); k < hdr.numEntries; k++ {
		d.FieldStruct("entry", func(d *decode.D) {
			key := d.FieldU32("key", elementTypeMap)

			// if the key has the top bit set, then (key & 0x7fffffff)
			// gives the offset of a string record.
			if key&0x80000000 != 0 {
				d.FieldStruct("key_string", func(d *decode.D) {
					d.SeekAbs(calcOffset(key&0x7fffffff), makeDecodeRecord())
				})
			}

			recordOffset := calcOffset(d.FieldU32("offset_to_record"))

			d.FieldU32("unused")

			switch key {
			case elementTypeTargetFlags:
				d.SeekAbs(recordOffset, func(d *decode.D) { decodeFlagDataObject(d, decodeTgtPropertyFlagBits) })
			case elementTypeVolumeFlags:
				d.SeekAbs(recordOffset, func(d *decode.D) { decodeFlagDataObject(d, decodeVolPropertyFlagBits) })
			default:
				d.SeekAbs(recordOffset, makeDecodeRecord())
			}

		})
	}
}

func decodeTOCHeader(d *decode.D) *tocHeader {
	hdr := new(tocHeader)

	d.FieldStruct("toc_header", func(d *decode.D) {
		d.FieldU32("toc_size")
		d.FieldU32("magic", d.UintAssert(0xfffffffe))
		d.FieldU32("identifier")
		hdr.nextTOCOffset = d.FieldU32("next_toc_offset")
		hdr.numEntries = d.FieldU32("num_entries_in_toc")
		hdr.entryArrayOffset = d.Pos()
	})

	return hdr
}

const (
	arrayEntrySize = 4
	dictEntrySize  = 4
)

func makeDecodeRecord() func(d *decode.D) {
	var pld apple.PosLoopDetector[int64]

	var decodeRecord func(d *decode.D)
	decodeRecord = func(d *decode.D) {
		defer pld.PushAndPop(
			d.Pos(),
			func() { d.Fatalf("infinite recursion detected in record decode function") },
		)()

		d.FieldStruct("record", func(d *decode.D) {
			n := int(d.FieldU32("length"))
			typ := d.FieldU32("type", dataTypeMap)
			switch typ {
			case dataTypeString:
				d.FieldUTF8("data", n)
				d.FieldRawLen("alignment_bytes", 32-((d.Pos()+32)%32))
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
				d.FieldF64BE("data", scalar.FltActualDateDescription(cocoaTimeEpochDate, time.Second, time.RFC3339))
			case dataTypeBooleanFalse:
			case dataTypeBooleanTrue:
			case dataTypeArray:
				d.FieldStructNArray("data", "element", int64(n/arrayEntrySize), func(d *decode.D) {
					offset := calcOffset(d.FieldU32("offset"))
					d.SeekAbs(offset, decodeRecord)
				})
			case dataTypeDictionary:
				d.FieldStructNArray("data", "element", int64(n/dictEntrySize), func(d *decode.D) {
					keyOffset := calcOffset(d.FieldU32("key_offset"))
					d.FieldStruct("key", func(d *decode.D) {
						d.SeekAbs(keyOffset, decodeRecord)
					})

					valueOffset := calcOffset(d.FieldU32("value_offset"))
					d.FieldStruct("value", func(d *decode.D) {
						d.SeekAbs(valueOffset, decodeRecord)
					})
				})
			case dataTypeUUID:
				d.FieldRawLen("data", int64(n*8))
			case dataTypeURL:
				d.FieldUTF8("data", n)
			case dataTypeRelativeURL:
				baseOffset := d.FieldU32("base_url_offset")
				d.FieldStruct("base_url", func(d *decode.D) {
					d.SeekAbs(int64(baseOffset), decodeRecord)
				})

				suffixOffset := d.FieldU32("suffix_offset")
				d.FieldStruct("suffix", func(d *decode.D) {
					d.SeekAbs(int64(suffixOffset), decodeRecord)
				})
			}
		})
	}
	return decodeRecord
}

const reservedSize = 32
const headerEnd = 48
const headerEndBitPos = headerEnd * 8

// all offsets are calculated relative to the end of the bookmark header
func calcOffset(i uint64) int64 { return int64(8 * (i + headerEnd)) }

func bookmarkDecode(d *decode.D) any {
	// all fields are little-endian with the exception of the Date datatype.
	d.Endian = decode.LittleEndian

	// decode bookmarkdata header, one at the top of each "file",
	// although these may be nested inside of binary plists
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("magic", 4, d.StrAssert("book", "alis"))
		d.FieldU32("total_size")
		d.FieldU32("unknown")
		d.FieldU32("header_size", d.UintAssert(48))
		d.FieldRawLen("reserved", reservedSize*8)
	})

	tocOffset := calcOffset(d.FieldU32("first_toc_offset"))

	var currentHdr *tocHeader
	var tocHeaders []*tocHeader

	hdrCount := 0
	d.FieldArrayLoop("toc_headers", func() bool {
		return tocOffset != headerEndBitPos || hdrCount > 100
	}, func(d *decode.D) {
		d.SeekAbs(tocOffset, func(d *decode.D) {
			currentHdr = decodeTOCHeader(d)
			tocOffset = calcOffset(currentHdr.nextTOCOffset)
			tocHeaders = append(tocHeaders, currentHdr)
		})
		hdrCount++
	})

	// now that we've collected all toc headers, iterate through each one's
	// entries and decode associated records.
	d.FieldArray("bookmark_entries",
		func(d *decode.D) {
			for _, hdr := range tocHeaders {
				d.SeekAbs(hdr.entryArrayOffset, hdr.decodeEntries)
			}
		})
	return nil
}
