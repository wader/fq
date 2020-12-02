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
		Name:        format.ELF,
		Description: "Executable and Linkable Format",
		Groups:      []string{format.PROBE},
		DecodeFn:    elfDecode,
	})
}

func elfDecode(d *decode.D) interface{} {
	d.ValidateAtLeastBitsLeft(128 * 8)

	var archBits int
	var endian decode.Endian

	d.FieldStructFn("ident", func(d *decode.D) {
		d.FieldValidateUTF8("magic", "\x7fELF")

		archBits = int(d.FieldUFn("class", func() (uint64, decode.DisplayFormat, string) {
			switch d.U8() {
			case 1:
				return 32, decode.NumberDecimal, ""
			case 2:
				return 64, decode.NumberDecimal, ""
			default:
				//d.Invalid()
			}
			panic("unreachable")
		}))
		d.FieldUFn("data", func() (uint64, decode.DisplayFormat, string) {
			switch d.U8() {
			case 1:
				endian = decode.LittleEndian
				return 1, decode.NumberDecimal, "Little-endian"
			case 2:
				endian = decode.BigEndian
				return 2, decode.NumberDecimal, "Big-endian"
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
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

	d.Endian = endian

	// TODO: hex functions?

	d.FieldU16LE("type")
	d.FieldU16LE("machine")
	d.FieldU32LE("version")
	d.FieldU("entry", archBits)
	d.FieldU("phoff", archBits)
	d.FieldU("shoff", archBits)
	d.FieldU32LE("flags")
	d.FieldU16LE("ehsize")
	d.FieldU16LE("phentsize")
	phnum := d.FieldU16LE("phnum")
	d.FieldU16LE("shentsize")
	shnum := d.FieldU16LE("shnum")
	d.FieldU16LE("shstrndx")

	d.FieldArrayFn("program_header", func(d *decode.D) {
		for i := uint64(0); i < phnum; i++ {
			d.FieldStructFn("program_header", func(d *decode.D) {
				switch archBits {
				case 32:
					d.FieldU32LE("p_type")
					d.FieldU("p_offset", archBits)
					d.FieldU("p_vaddr", archBits)
					d.FieldU("p_paddr", archBits)
					d.FieldU32LE("p_filesz")
					d.FieldU32LE("p_memsz")
					d.FieldU32LE("p_flags")
					d.FieldU32LE("p_align")
				case 64:
					d.FieldU32LE("p_type")
					d.FieldU32LE("p_flags")
					d.FieldU("p_offset", archBits)
					d.FieldU("p_vaddr", archBits)
					d.FieldU("p_paddr", archBits)
					d.FieldU64LE("p_filesz")
					d.FieldU64LE("p_memsz")
					d.FieldU64LE("p_align")
				}
			})
		}
	})

	d.FieldArrayFn("section_header", func(d *decode.D) {
		for i := uint64(0); i < shnum; i++ {
			d.FieldStructFn("section_header", func(d *decode.D) {
				switch archBits {
				case 32:
					d.FieldU32LE("sh_name")
					d.FieldU32LE("sh_type")
					d.FieldU32LE("sh_flags")
					d.FieldU("sh_addr", archBits)
					d.FieldU("sh_offset", archBits)
					d.FieldU32LE("sh_size")
					d.FieldU32LE("sh_link")
					d.FieldU32LE("sh_info")
					d.FieldU32LE("sh_addralign")
					d.FieldU32LE("sh_entsize")
				case 64:
					d.FieldU32LE("sh_name")
					d.FieldU32LE("sh_type")
					d.FieldU64LE("sh_flags")
					d.FieldU("sh_addr", archBits)
					d.FieldU("sh_offset", archBits)
					d.FieldU64LE("sh_size")
					d.FieldU32LE("sh_link")
					d.FieldU32LE("sh_info")
					d.FieldU64LE("sh_addralign")
					d.FieldU64LE("sh_entsize")
				}
			})
		}
	})

	return nil
}
