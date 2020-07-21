package elf

import "fq/pkg/decode"

/*

ElfN_Addr       Unsigned program address, uintN_t
ElfN_Off        Unsigned file offset, uintN_t
ElfN_Section    Unsigned section index, uint16_t
ElfN_Versym     Unsigned version symbol information, uint16_t
Elf_Byte        unsigned char
ElfN_Half       uint16_t
ElfN_Sword      int32_t
ElfN_Word       uint32_t
ElfN_Sxword     int64_t
ElfN_Xword      uint64_t

typedef struct {
	unsigned char e_ident[EI_NIDENT];
	uint16_t      e_type;
	uint16_t      e_machine;
	uint32_t      e_version;
	ElfN_Addr     e_entry;
	ElfN_Off      e_phoff;
	ElfN_Off      e_shoff;
	uint32_t      e_flags;
	uint16_t      e_ehsize;
	uint16_t      e_phentsize;
	uint16_t      e_phnum;
	uint16_t      e_shentsize;
	uint16_t      e_shnum;
	uint16_t      e_shstrndx;
} ElfN_Ehdr;

*/

var File = &decode.Format{
	Name: "elf",
	New:  func() decode.Decoder { return &FileDecoder{} },
}

// FileDecoder is ELF file decoder
type FileDecoder struct{ decode.Common }

// Decode a ELF file
func (d *FileDecoder) Decode() {
	d.ValidateAtLeastBitsLeft(128 * 8)

	var field16 func(name string) uint64
	var field32 func(name string) uint64
	var fieldN func(name string) uint64
	var fieldNX func(name string) uint64
	var dN func() uint64

	d.FieldNoneFn("ident", func() {
		d.FieldValidateString("magic", "\x7fELF")

		archBits := d.FieldUFn("class", func() (uint64, decode.NumberFormat, string) {
			switch d.U8() {
			case 1:
				return 32, decode.NumberDecimal, ""
			case 2:
				return 64, decode.NumberDecimal, ""
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
		_ = archBits
		isBigEndian := true
		d.FieldUFn("data", func() (uint64, decode.NumberFormat, string) {
			switch d.U8() {
			case 1:
				isBigEndian = false
				field16 = d.FieldU16LE
				field32 = d.FieldU32LE
				switch archBits {
				case 32:
					fieldN = d.FieldU32LE
					dN = d.U32LE
				case 64:
					fieldN = d.FieldU64LE
					dN = d.U64LE
				}
				return 1, decode.NumberDecimal, "Little-endian"
			case 2:
				field16 = d.FieldU16BE
				field32 = d.FieldU32BE
				switch archBits {
				case 32:
					fieldN = d.FieldU32BE
					dN = d.U32BE
				case 64:
					fieldN = d.FieldU64BE
					dN = d.U64BE
				}
				return 2, decode.NumberDecimal, "Big-endian"
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
		fieldNX = func(name string) uint64 {
			return d.FieldUFn(name, func() (uint64, decode.NumberFormat, string) {
				return dN(), decode.NumberHex, ""
			})
		}
		_ = isBigEndian
		d.FieldU8("version")
		d.FieldStringMapFn("os_abi", map[uint64]string{
			0:   "Sysv",
			1:   "HPUX",
			2:   "NetBSD",
			3:   "Linux",
			4:   "Hurd",
			5:   "86open",
			6:   "Solaris",
			7:   "Monterey",
			8:   "Irix",
			9:   "FreeBSD",
			10:  "Tru64",
			11:  "Modesto",
			12:  "OpenBSD",
			97:  "Arm",
			255: "Standalone",
		}, "Unknown", d.U8)
		d.FieldU8("abi_version")
		d.FieldValidateZeroPadding("pad", 7*8)
	})

	field16("type")
	field16("machine")
	field32("version")
	fieldNX("entry")
	fieldNX("phoff")
	fieldNX("shoff")
	field32("flags")
	ehsize := field16("ehsize")
	phentsize := field16("phentsize")
	phnum := field16("phnum")
	shentsize := field16("shentsize")
	shnum := field16("shnum")
	field16("shstrndx")

	for i := 0 < 
}
