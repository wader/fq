package pe

// https://osandamalith.com/2020/07/19/exploring-the-ms-dos-stub/

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: probe?

func init() {
	interp.RegisterFormat(
		format.MSDOS_Stub,
		&decode.Format{
			Description: "MS-DOS Stub",
			DecodeFn:    msDosStubDecode,
		})
}

func msDosStubDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldU16("e_magic", scalar.UintDescription("Magic number"), d.UintAssert(0x5a4d), scalar.UintHex)
	d.FieldU16("e_cblp", scalar.UintDescription("Bytes on last page of file"))
	d.FieldU16("e_cp", scalar.UintDescription("Pages in file"))
	d.FieldU16("e_crlc", scalar.UintDescription("Relocations"))
	d.FieldU16("e_cparhdr", scalar.UintDescription("Size of header in paragraphs"))
	d.FieldU16("e_minalloc", scalar.UintDescription("Minimum extra paragraphs needed"))
	d.FieldU16("e_maxalloc", scalar.UintDescription("Maximum extra paragraphs needed"))
	d.FieldU16("e_ss", scalar.UintDescription("Initial (relative) SS value"))
	d.FieldU16("e_sp", scalar.UintDescription("Initial SP value"))
	d.FieldU16("e_csum", scalar.UintDescription("Checksum"))
	d.FieldU16("e_ip", scalar.UintDescription("Initial IP value"))
	d.FieldU16("e_cs", scalar.UintDescription("Initial (relative) CS value"))
	d.FieldU16("e_lfarlc", scalar.UintDescription("File address of relocation table"))
	d.FieldU16("e_ovno", scalar.UintDescription("Overlay number"))
	d.FieldRawLen("e_res", 4*16, scalar.BitBufDescription("Reserved words"))
	d.FieldU16("e_oemid", scalar.UintDescription("OEM identifier (for e_oeminfo)"))
	d.FieldU16("e_oeminfo", scalar.UintDescription("OEM information; e_oemid specific"))
	d.FieldRawLen("e_res2", 10*16, scalar.BitBufDescription("Reserved words"))
	lfanew := d.FieldU32("e_lfanew", scalar.UintDescription("File address of new exe header"))

	// TODO: how to detect UEFI?

	subEndPos := mathex.Min(d.Pos()+64*8, int64(lfanew)*8)

	// TODO: x86 format in the future
	d.FieldRawLen("stub", subEndPos-d.Pos(), scalar.BitBufDescription("Sub program"))

	// TODO: is not padding i guess?
	padding := lfanew*8 - uint64(subEndPos)
	if padding > 0 {
		d.FieldRawLen("padding", int64(padding))
	}

	return format.MS_DOS_Out{
		LFANew: int(lfanew),
	}
}
