package elf

import "fq/internal/decode"

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

var Register = &decode.Register{
	Name: "elf",
	MIME: "",
	New:  func(common decode.Common) decode.Decoder { return &Decoder{Common: common} },
}

// Decoder is ELF decoder
type Decoder struct{ decode.Common }

// Decode ID3v1
func (d *Decoder) Decode(opts decode.Options) {
	d.ValidateAtLeastBitsLeft(128 * 8)

	d.FieldNoneFn("ident", func() {
		d.FieldValidateString("magic", "\x7fELF")

		archBits := d.FieldUFn("class", func() (uint64, decode.Format, string) {
			switch d.U8() {
			case 1:
				return 32, decode.FormatDecimal, ""
			case 2:
				return 64, decode.FormatDecimal, ""
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
		_ = archBits
		isBigEndian := true
		d.FieldUFn("data", func() (uint64, decode.Format, string) {
			switch d.U8() {
			case 1:
				isBigEndian = false
				return 1, decode.FormatDecimal, "Little-endian"
			case 2:
				return 2, decode.FormatDecimal, "Big-endian"
			default:
				//d.Invalid()
			}
			panic("unreachable")
		})
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

	})
}
