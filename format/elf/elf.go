package elf

// https://en.wikipedia.org/wiki/Executable_and_Linkable_Format
// https://man7.org/linux/man-pages/man5/elf.5.html
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/elf.h

import (
	"fq/format"
	"fq/pkg/decode"
	"strings"
)

// TODO: p_type hi/lo

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.ELF,
		Description: "Executable and Linkable Format",
		Groups:      []string{format.PROBE},
		DecodeFn:    elfDecode,
	})
}

func strIndexNull(idx int, s string) string {
	if idx > len(s) {
		return ""
	}
	i := strings.IndexByte(s[idx:], 0)
	if i == -1 {
		return s
	}
	return s[idx : idx+i]
}

func fieldStringStrIndexFn(d *decode.D, name string, strTable string, fn func() uint64) string {
	return d.FieldFn(name, func() *decode.Value {
		idx := fn()
		return &decode.Value{V: idx, Symbol: strIndexNull(int(idx), strTable)}
	}).Symbol
}

func elfDecode(d *decode.D, in interface{}) interface{} {
	d.ValidateAtLeastBitsLeft(128 * 8)

	rootD := d
	var archBits int

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
				rootD.Endian = decode.LittleEndian
				return 1, decode.NumberDecimal, "Little-endian"
			case 2:
				rootD.Endian = decode.BigEndian
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
		0x0a:  "MIPS RS3000 Little-endian",
		0x0e:  "Hewlett-Packard PA-RISC",
		0x0f:  "Reserved for future use",
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
		0x2a:  "SuperH",
		0x2b:  "SPARC Version 9",
		0x2c:  "Siemens TriCore embedded processor",
		0x2d:  "Argonaut RISC Core",
		0x2e:  "Hitachi H8/300",
		0x2f:  "Hitachi H8/300H",
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
		0x3a:  "Motorola Star*Core processor",
		0x3b:  "Toyota ME16 processor",
		0x3c:  "STMicroelectronics ST100 processor",
		0x3d:  "Advanced Logic Corp. TinyJ embedded processor family",
		0x3e:  "AMD x86-64",
		0x8c:  "TMS320C6000 Family",
		0xb7:  "ARM 64-bits (ARMv8/Aarch64)",
		0xf3:  "RISC-V",
		0xf7:  "Berkeley Packet Filter",
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
	var strIndexTable string
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
			_ = strIndexTable
		})

		strIndexTable = string(d.BytesRange(int64(strTableOffset*8), int(strTableSize)*8))
	}

	// d.DecodeRangeFn(int64(phoff)*8, int64(phnum*phsize*8), func(d *decode.D) {
	d.FieldArrayFn("program_headers", func(d *decode.D) {
		for i := uint64(0); i < phnum; i++ {
			d.SeekAbs(int64(phoff*8) + int64(i*phsize*8))

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
				0x6fffffff: "PT_HIOS",
				0x70000000: "PT_LOPROC",
				0x7fffffff: "PT_HIPROC",
			}

			pFlags := func(d *decode.D) {
				d.FieldStructFn("p_flags", func(d *decode.D) {
					if d.Endian == decode.LittleEndian {
						d.FieldU5("unused0")
						d.FieldBool("PF_R")
						d.FieldBool("PF_W")
						d.FieldBool("PF_X")
						d.FieldU24("unused1")
					} else {
						d.FieldU29("unused0")
						d.FieldBool("PF_R")
						d.FieldBool("PF_W")
						d.FieldBool("PF_X")
					}
				})
			}

			d.FieldStructFn("program_header", func(d *decode.D) {
				var offset uint64
				var size uint64

				switch archBits {
				case 32:
					d.FieldStringMapFn("p_type", pTypeNames, "Unknown", func() uint64 { return d.U32() & 0xf }, decode.NumberDecimal)
					offset = d.FieldU("p_offset", archBits)
					d.FieldU("p_vaddr", archBits)
					d.FieldU("p_paddr", archBits)
					size = d.FieldU32("p_filesz")
					d.FieldU32("p_memsz")
					pFlags(d)
					d.FieldU32("p_align")
				case 64:
					d.FieldStringMapFn("p_type", pTypeNames, "Unknown", func() uint64 { return d.U32() & 0xf }, decode.NumberDecimal)
					pFlags(d)
					offset = d.FieldU("p_offset", archBits)
					d.FieldU("p_vaddr", archBits)
					d.FieldU("p_paddr", archBits)
					size = d.FieldU64("p_filesz")
					d.FieldU64("p_memsz")
					d.FieldU64("p_align")
				}

				d.FieldBitBufRange("data", int64(offset*8), int64(size*8))
			})
		}
	})
	// })

	// d.DecodeRangeFn(int64(shoff)*8, int64(shnum*shentsize*8), func(d *decode.D) {
	d.FieldArrayFn("section_headers", func(d *decode.D) {
		for i := uint64(0); i < shnum; i++ {
			d.SeekAbs(int64(shoff*8) + int64(i*shentsize*8))

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

			shFlags := func(d *decode.D, archBits int) {
				d.FieldStructFn("sh_flags", func(d *decode.D) {
					if d.Endian == decode.LittleEndian {
						d.FieldBool("SHF_LINK_ORDER")
						d.FieldBool("SHF_INFO_LINK")
						d.FieldBool("SHF_STRINGS")
						d.FieldBool("SHF_MERGE")
						d.FieldU1("unused0")
						d.FieldBool("SHF_EXECINSTR")
						d.FieldBool("SHF_ALLOC")
						d.FieldBool("SHF_WRITE")
						d.FieldBool("SHF_TLS")
						d.FieldBool("SHF_GROUP")
						d.FieldBool("SHF_OS_NONCONFORMING")

						d.FieldU9("unused1")

						d.FieldU8("os_specific")
						d.FieldU4("processor_specific")
						if archBits == 64 {
							d.FieldU32("unused2")
						}
					} else {
						// TODO: add d.FieldUnused that is per decoder?
						if archBits == 64 {
							d.FieldU32("unused0")
						}
						d.FieldU4("processor_specific")
						d.FieldU8("os_specific")
						d.FieldU9("unused1")
						d.FieldBool("SHF_TLS")
						d.FieldBool("SHF_GROUP")
						d.FieldBool("SHF_OS_NONCONFORMING")
						d.FieldBool("SHF_LINK_ORDER")
						d.FieldBool("SHF_INFO_LINK")
						d.FieldBool("SHF_STRINGS")
						d.FieldBool("SHF_MERGE")
						d.FieldU1("unused2")
						d.FieldBool("SHF_EXECINSTR")
						d.FieldBool("SHF_ALLOC")
						d.FieldBool("SHF_WRITE")
						// 0x1	SHF_WRITE	Writable
						// 0x2	SHF_ALLOC	Occupies memory during execution
						// 0x4	SHF_EXECINSTR	Executable
						// 0x10	SHF_MERGE	Might be merged
						// 0x20	SHF_STRINGS	Contains null-terminated strings
						// 0x40	SHF_INFO_LINK	'sh_info' contains SHT index
						// 0x80	SHF_LINK_ORDER	Preserve order after combining
						// 0x100	SHF_OS_NONCONFORMING	Non-standard OS specific handling required
						// 0x200	SHF_GROUP	Section is member of a group
						// 0x400	SHF_TLS	Section hold thread-local data
						// 0x0ff00000	SHF_MASKOS	OS-specific
						// 0xf0000000	SHF_MASKPROC	Processor-specific
						// 0x4000000	SHF_ORDERED	Special ordering requirement (Solaris)
						// 0x8000000	SHF_EXCLUDE	Section is excluded unless referenced or allocated (Solaris)
					}
				})
			}

			d.FieldStructFn("section_header", func(d *decode.D) {
				var offset uint64
				var size uint64
				var shname string

				const (
					DT_NULL     = 0
					DT_NEEDED   = 1
					DT_PLTRELSZ = 2
					DT_PLTGOT   = 3
					DT_HASH     = 4
					DT_STRTAB   = 5
					DT_SYMTAB   = 6
					DT_RELA     = 7
					DT_RELASZ   = 8
					DT_RELAENT  = 9
					DT_STRSZ    = 10
					DT_SYMENT   = 11
					DT_INIT     = 12
					DT_FINI     = 13
					DT_SONAME   = 14
					DT_RPATH    = 15
					DT_SYMBOLIC = 16
					DT_REL      = 17
					DT_RELSZ    = 18
					DT_RELENT   = 19
					DT_PLTREL   = 20
					DT_DEBUG    = 21
					DT_TEXTREL  = 22
					DT_JMPREL   = 23
					DT_ENCODING = 32
				)
				var dtNames = map[uint64]string{
					DT_NULL:     "DT_NULL",
					DT_NEEDED:   "DT_NEEDED",
					DT_PLTRELSZ: "DT_PLTRELSZ",
					DT_PLTGOT:   "DT_PLTGOT",
					DT_HASH:     "DT_HASH",
					DT_STRTAB:   "DT_STRTAB",
					DT_SYMTAB:   "DT_SYMTAB",
					DT_RELA:     "DT_RELA",
					DT_RELASZ:   "DT_RELASZ",
					DT_RELAENT:  "DT_RELAENT",
					DT_STRSZ:    "DT_STRSZ",
					DT_SYMENT:   "DT_SYMENT",
					DT_INIT:     "DT_INIT",
					DT_FINI:     "DT_FINI",
					DT_SONAME:   "DT_SONAME",
					DT_RPATH:    "DT_RPATH",
					DT_SYMBOLIC: "DT_SYMBOLIC",
					DT_REL:      "DT_REL",
					DT_RELSZ:    "DT_RELSZ",
					DT_RELENT:   "DT_RELENT",
					DT_PLTREL:   "DT_PLTREL",
					DT_DEBUG:    "DT_DEBUG",
					DT_TEXTREL:  "DT_TEXTREL",
					DT_JMPREL:   "DT_JMPREL",
					DT_ENCODING: "DT_ENCODING",
				}

				switch archBits {
				case 32:
					shname = fieldStringStrIndexFn(d, "sh_name", strIndexTable, d.U32)
					d.FieldStringMapFn("sh_type", shTypeNames, "Unknown", d.U32, decode.NumberHex)
					shFlags(d, archBits)
					d.FieldU("sh_addr", archBits)
					offset = d.FieldU("sh_offset", archBits)
					size = d.FieldU32("sh_size")
					d.FieldU32("sh_link")
					d.FieldU32("sh_info")
					d.FieldU32("sh_addralign")
					d.FieldU32("sh_entsize")
				case 64:
					shname = fieldStringStrIndexFn(d, "sh_name", strIndexTable, d.U32)
					d.FieldStringMapFn("sh_type", shTypeNames, "Unknown", d.U32, decode.NumberHex)
					shFlags(d, archBits)
					d.FieldU("sh_addr", archBits)
					offset = d.FieldU("sh_offset", archBits)
					size = d.FieldU64("sh_size")
					d.FieldU32("sh_link")
					d.FieldU32("sh_info")
					d.FieldU64("sh_addralign")
					d.FieldU64("sh_entsize")
				}

				d.FieldBitBufRange("data", int64(offset*8), int64(size*8))

				d.DecodeRangeFn(int64(offset)*8, int64(size*8), func(d *decode.D) {
					switch shname {
					// TODO: PT_DYNAMIC?
					case ".dynamic":
						d.FieldArrayFn("dynamic_tags", func(d *decode.D) {
							for d.NotEnd() {
								d.FieldStructFn("tag", func(d *decode.D) {
									tag, _ := d.FieldStringMapFn("tag", dtNames, "Unknown", func() uint64 { return d.U(archBits) }, decode.NumberHex)
									switch tag {
									case DT_NEEDED:
										// TODO: DT_STRTAB
										fieldStringStrIndexFn(d, "val", strIndexTable, func() uint64 { return d.U(archBits) })
									default:
										d.FieldU("d_un", archBits)
									}
								})
							}
						})
					}

				})
			})
		}
	})
	// })

	return nil
}
