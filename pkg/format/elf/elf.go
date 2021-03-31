package elf

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
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

func elfDecode(d *decode.D, in interface{}) interface{} {
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
		}, "Unknown", d.U8, decode.NumberDecimal)
		d.FieldU8("abi_version")
		d.FieldValidateZeroPadding("pad", 7*8)
	})

	d.Endian = endian

	// TODO: hex functions?

	d.FieldStringMapFn("type", map[uint64]string{
		0x00:   "None",
		0x01:   "Rel",
		0x02:   "Exec",
		0x03:   "Dyn",
		0x04:   "Core",
		0xfe00: "Loos",
		0xfeff: "Hios",
		0xff00: "Loproc",
		0xffff: "Hiproc",
	}, "Unknown", d.U16, decode.NumberHex)

	d.FieldStringMapFn("machine", map[uint64]string{
		0x00:  "No specific instruction set",
		0x01:  "AT&T WE 32100",
		0x02:  "SPARC",
		0x03:  "x86",
		0x04:  "Motorola 68000 (M68k)",
		0x05:  "Motorola 88000 (M88k)",
		0x06:  "Intel MCU",
		0x07:  "Intel 80860",
		0x08:  "MIPS",
		0x09:  "IBM_System/370",
		0x0A:  "MIPS RS3000 Little-endian",
		0x0E:  "Hewlett-Packard PA-RISC",
		0x0F:  "Reserved for future use",
		0x13:  "Intel 80960",
		0x14:  "PowerPC",
		0x15:  "PowerPC (64-bit)",
		0x16:  "S390, including S390x",
		0x17:  "IBM SPU/SPC",
		0x24:  "NEC V800",
		0x25:  "Fujitsu FR20",
		0x26:  "TRW RH-32",
		0x27:  "Motorola RCE",
		0x28:  "ARM (up to ARMv7/Aarch32)",
		0x29:  "Digital Alpha",
		0x2A:  "SuperH",
		0x2B:  "SPARC Version 9",
		0x2C:  "Siemens TriCore embedded processor",
		0x2D:  "Argonaut RISC Core",
		0x2E:  "Hitachi H8/300",
		0x2F:  "Hitachi H8/300H",
		0x30:  "Hitachi H8S",
		0x31:  "Hitachi H8/500",
		0x32:  "IA-64",
		0x33:  "Stanford MIPS-X",
		0x34:  "Motorola ColdFire",
		0x35:  "Motorola M68HC12",
		0x36:  "Fujitsu MMA Multimedia Accelerator",
		0x37:  "Siemens PCP",
		0x38:  "Sony nCPU embedded RISC processor",
		0x39:  "Denso NDR1 microprocessor",
		0x3A:  "Motorola Star*Core processor",
		0x3B:  "Toyota ME16 processor",
		0x3C:  "STMicroelectronics ST100 processor",
		0x3D:  "Advanced Logic Corp. TinyJ embedded processor family",
		0x3E:  "AMD x86-64",
		0x8C:  "TMS320C6000 Family",
		0xB7:  "ARM 64-bits (ARMv8/Aarch64)",
		0xF3:  "RISC-V",
		0xF7:  "Berkeley Packet Filter",
		0x101: "WDC 65C816",
	}, "Unknown", d.U16, decode.NumberHex)

	d.FieldU32("version")
	d.FieldU("entry", archBits)
	phoff := d.FieldU("phoff", archBits)
	shoff := d.FieldU("shoff", archBits)
	d.FieldU32("flags")
	d.FieldU16("ehsize")
	phsize := d.FieldU16("phentsize")
	phnum := d.FieldU16("phnum")
	shentsize := d.FieldU16("shentsize")
	shnum := d.FieldU16("shnum")
	shstrndx := d.FieldU16("shstrndx")

	// TODO: make this nicer, API to update fields?
	// TODO: is wrong: string table is one large string to index into
	// TODO: and string can overlap
	strTable := map[uint64]string{}
	if shstrndx != 0 {
		var strTableOffset uint64
		var strTableSize uint64
		d.DecodeRangeFn(int64((shoff+shstrndx*shentsize)*8), int64(shentsize*8), func(d *decode.D) {
			d.SeekRel(32)
			d.SeekRel(32)
			d.SeekRel(int64(archBits))
			d.SeekRel(int64(archBits))
			strTableOffset = d.U(archBits)
			strTableSize = d.U(archBits)
			_ = strTable
		})
		d.DecodeRangeFn(int64(strTableOffset*8), int64(strTableSize*8), func(d *decode.D) {
			var i uint64
			for d.NotEnd() {
				s := d.StrZeroTerminated()
				strTable[i] = s
				i += uint64(len(s)) + 1
			}
		})
	}

	log.Printf("strTable: %#+v\n", strTable)

	d.DecodeRangeFn(int64(phoff)*8, int64(phnum*phsize*8), func(d *decode.D) {
		d.FieldArrayFn("program_headers", func(d *decode.D) {
			for i := uint64(0); i < phnum; i++ {

				pTypeNames := map[uint64]string{
					0x00000000: "PT_NULL",
					0x00000001: "PT_LOAD",
					0x00000002: "PT_DYNAMIC",
					0x00000003: "PT_INTERP",
					0x00000004: "PT_NOTE",
					0x00000005: "PT_SHLIB",
					0x00000006: "PT_PHDR",
					0x00000007: "PT_TLS",
					0x60000000: "PT_LOOS",
					0x6FFFFFFF: "PT_HIOS",
					0x70000000: "PT_LOPROC",
					0x7FFFFFFF: "PT_HIPROC",
				}

				d.FieldStructFn("program_header", func(d *decode.D) {
					switch archBits {
					case 32:
						d.FieldStringMapFn("p_type", pTypeNames, "Unknown", d.U32, decode.NumberDecimal)
						d.FieldU("p_offset", archBits)
						d.FieldU("p_vaddr", archBits)
						d.FieldU("p_paddr", archBits)
						d.FieldU32("p_filesz")
						d.FieldU32("p_memsz")
						d.FieldU32("p_flags")
						d.FieldU32("p_align")
					case 64:
						d.FieldStringMapFn("p_type", pTypeNames, "Unknown", d.U32, decode.NumberDecimal)
						d.FieldU32("p_flags")
						d.FieldU("p_offset", archBits)
						d.FieldU("p_vaddr", archBits)
						d.FieldU("p_paddr", archBits)
						d.FieldU64("p_filesz")
						d.FieldU64("p_memsz")
						d.FieldU64("p_align")
					}
				})
			}
		})
	})

	d.DecodeRangeFn(int64(shoff)*8, int64(shnum*shentsize*8), func(d *decode.D) {
		d.FieldArrayFn("section_headers", func(d *decode.D) {
			for i := uint64(0); i < shnum; i++ {
				d.FieldStructFn("section_header", func(d *decode.D) {

					shTypeNames := map[uint64]string{
						0x0:        "SHT_NULL",
						0x1:        "SHT_PROGBITS",
						0x2:        "SHT_SYMTAB",
						0x3:        "SHT_STRTAB",
						0x4:        "SHT_RELA",
						0x5:        "SHT_HASH",
						0x6:        "SHT_DYNAMIC",
						0x7:        "SHT_NOTE",
						0x8:        "SHT_NOBITS",
						0x9:        "SHT_REL",
						0x0a:       "SHT_SHLIB",
						0x0b:       "SHT_DYNSYM",
						0x0e:       "SHT_INIT_ARRAY",
						0x0f:       "SHT_FINI_ARRAY",
						0x10:       "SHT_PREINIT_ARRAY",
						0x11:       "SHT_GROUP",
						0x12:       "SHT_SYMTAB_SHNDX",
						0x13:       "SHT_NUM",
						0x60000000: "SHT_LOOS",
					}

					switch archBits {
					case 32:
						d.FieldStringMapFn("sh_name", strTable, "Unknown", d.U32, decode.NumberDecimal)
						d.FieldStringMapFn("sh_type", shTypeNames, "Unknown", d.U32, decode.NumberHex)
						d.FieldU32("sh_flags")
						d.FieldU("sh_addr", archBits)
						d.FieldU("sh_offset", archBits)
						d.FieldU32("sh_size")
						d.FieldU32("sh_link")
						d.FieldU32("sh_info")
						d.FieldU32("sh_addralign")
						d.FieldU32("sh_entsize")
					case 64:
						d.FieldStringMapFn("sh_name", strTable, "Unknown", d.U32, decode.NumberDecimal)
						d.FieldStringMapFn("sh_type", shTypeNames, "Unknown", d.U32, decode.NumberHex)
						d.FieldU64("sh_flags")
						d.FieldU("sh_addr", archBits)
						d.FieldU("sh_offset", archBits)
						d.FieldU64("sh_size")
						d.FieldU32("sh_link")
						d.FieldU32("sh_info")
						d.FieldU64("sh_addralign")
						d.FieldU64("sh_entsize")
					}
				})
			}
		})
	})

	return nil
}
