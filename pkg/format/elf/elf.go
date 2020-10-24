package elf

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

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

func init() {
	format.MustRegister(&decode.Format{
		Name:     "elf",
		DecodeFn: elfDecode,
	})
}

func elfDecode(d *decode.Common) interface{} {
	d.ValidateAtLeastBitsLeft(128 * 8)

	// TODO: make endian switching nicer somehow?
	var archBits uint64
	var field16 func(d *decode.Common, name string) uint64
	var field32 func(d *decode.Common, name string) uint64
	var field64 func(d *decode.Common, name string) uint64
	var fieldNX func(d *decode.Common, name string) uint64
	var dN func(d *decode.Common) uint64

	d.FieldStructFn2("ident", func(d *decode.Common) {
		d.FieldValidateString("magic", "\x7fELF")

		archBits = d.FieldUFn("class", func() (uint64, decode.DisplayFormat, string) {
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
		isBigEndian := true
		d.FieldUFn("data", func() (uint64, decode.DisplayFormat, string) {
			switch d.U8() {
			case 1:
				isBigEndian = false
				field16 = func(d *decode.Common, name string) uint64 { return d.FieldU16LE(name) }
				field32 = func(d *decode.Common, name string) uint64 { return d.FieldU32LE(name) }
				field64 = func(d *decode.Common, name string) uint64 { return d.FieldU64LE(name) }
				switch archBits {
				case 32:
					dN = func(d *decode.Common) uint64 { return d.U32LE() }
				case 64:
					dN = func(d *decode.Common) uint64 { return d.U64LE() }
				}
				return 1, decode.NumberDecimal, "Little-endian"
			case 2:
				field16 = func(d *decode.Common, name string) uint64 { return d.FieldU16BE(name) }
				field32 = func(d *decode.Common, name string) uint64 { return d.FieldU32BE(name) }
				field64 = func(d *decode.Common, name string) uint64 { return d.FieldU64BE(name) }
				switch archBits {
				case 32:
					dN = func(d *decode.Common) uint64 { return d.U32BE() }
				case 64:
					dN = func(d *decode.Common) uint64 { return d.U64BE() }
				}
				return 2, decode.NumberDecimal, "Big-endian"
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
		fieldNX = func(d *decode.Common, name string) uint64 {
			return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
				return dN(d), decode.NumberHex, ""
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

	field16(d, "type")
	field16(d, "machine")
	field32(d, "version")
	fieldNX(d, "entry")
	fieldNX(d, "phoff")
	fieldNX(d, "shoff")
	field32(d, "flags")
	field16(d, "ehsize")
	field16(d, "phentsize")
	phnum := field16(d, "phnum")
	field16(d, "shentsize")
	shnum := field16(d, "shnum")
	field16(d, "shstrndx")

	d.FieldArrayFn2("program_header", func(d *decode.Common) {
		for i := uint64(0); i < phnum; i++ {
			d.FieldStructFn2("program_header", func(d *decode.Common) {
				switch archBits {
				case 32:
					field32(d, "p_type")
					fieldNX(d, "p_offset")
					fieldNX(d, "p_vaddr")
					fieldNX(d, "p_paddr")
					field32(d, "p_filesz")
					field32(d, "p_memsz")
					field32(d, "p_flags")
					field32(d, "p_align")
				case 64:
					field32(d, "p_type")
					field32(d, "p_flags")
					fieldNX(d, "p_offset")
					fieldNX(d, "p_vaddr")
					fieldNX(d, "p_paddr")
					field64(d, "p_filesz")
					field64(d, "p_memsz")
					field64(d, "p_align")
				}
			})
		}
	})

	d.FieldArrayFn2("section_header", func(d *decode.Common) {
		for i := uint64(0); i < shnum; i++ {
			d.FieldStructFn2("section_header", func(d *decode.Common) {
				switch archBits {
				case 32:
					field32(d, "sh_name")
					field32(d, "sh_type")
					field32(d, "sh_flags")
					fieldNX(d, "sh_addr")
					fieldNX(d, "sh_offset")
					field32(d, "sh_size")
					field32(d, "sh_link")
					field32(d, "sh_info")
					field32(d, "sh_addralign")
					field32(d, "sh_entsize")
				case 64:
					field32(d, "sh_name")
					field32(d, "sh_type")
					field64(d, "sh_flags")
					fieldNX(d, "sh_addr")
					fieldNX(d, "sh_offset")
					field64(d, "sh_size")
					field32(d, "sh_link")
					field32(d, "sh_info")
					field64(d, "sh_addralign")
					field64(d, "sh_entsize")
				}
			})
		}
	})

	return nil
}
