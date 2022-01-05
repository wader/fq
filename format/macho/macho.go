package macho

// https://github.com/aidansteele/osx-abi-macho-file-format-reference

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MACHO,
		Description: "Mach-O is the native executable format of binaries in OS X",
		Groups:      []string{format.PROBE},
		DecodeFn:    machoDecode,
	})
}

const (
	MH_MAGIC    = 0xfeedface
	MH_CIGAM    = 0xcefaedfe
	MH_MAGIC_64 = 0xfeedfacf
	MH_CIGAM_64 = 0xcffaedfe
)

var classBits = scalar.UToSymU{
	MH_MAGIC:    32,
	MH_CIGAM:    32,
	MH_MAGIC_64: 64,
	MH_CIGAM_64: 64,
}

var endianNames = scalar.UToSymStr{
	MH_MAGIC:    "little-endian",
	MH_CIGAM:    "big-endian",
	MH_MAGIC_64: "little-endian",
	MH_CIGAM_64: "big-endian",
}

var cpuTypes = scalar.SToSymStr{
	-1:        "CPU_TYPE_ANY",
	1:         "CPU_TYPE_VAX",
	2:         "CPU_TYPE_ROMP",
	4:         "CPU_TYPE_NS32032",
	5:         "CPU_TYPE_NS32332",
	6:         "CPU_TYPE_MC680x0",
	7:         "CPU_TYPE_X86",
	8:         "CPU_TYPE_MIPS",
	9:         "CPU_TYPE_NS32532",
	10:        "CPU_TYPE_MC98000",
	11:        "CPU_TYPE_HPPA",
	12:        "CPU_TYPE_ARM",
	13:        "CPU_TYPE_MC88000",
	14:        "CPU_TYPE_SPARC",
	15:        "CPU_TYPE_I860",
	16:        "CPU_TYPE_I860_LITTLE",
	17:        "CPU_TYPE_RS6000",
	18:        "CPU_TYPE_POWERPC",
	0x1000007: "CPU_TYPE_X86_64",
	0x1000012: "CPU_TYPE_ARM64",
	0x1000018: "CPU_TYPE_POWERPC64",
	255:       "CPU_TYPE_VEO",
}

// TODO subtypes stosymstr depends on cputype and a signed integer
var cpuSubTypes = scalar.SToSymStr{
	-1: "CPU_SUBTYPE_MULTIPLE",
	0:  "CPU_SUBTYPE_VAX_ALL",
	1:  "CPU_SUBTYPE_VAX780",
	2:  "CPU_SUBTYPE_VAX785",
	3:  "CPU_SUBTYPE_VAX750",
	4:  "CPU_SUBTYPE_VAX730",
	5:  "CPU_SUBTYPE_UVAXI",
	6:  "CPU_SUBTYPE_UVAXII",
	7:  "CPU_SUBTYPE_VAX8200",
	8:  "CPU_SUBTYPE_VAX8500",
	9:  "CPU_SUBTYPE_VAX8600",
	10: "CPU_SUBTYPE_VAX8650",
	11: "CPU_SUBTYPE_VAX8800",
	12: "CPU_SUBTYPE_UVAXIII",
}

const (
	LC_REQ_DYLD                 = 0x80000000
	LC_SEGMENT                  = 0x1
	LC_SYMTAB                   = 0x2
	LC_SYMSEG                   = 0x3
	LC_THREAD                   = 0x4
	LC_UNIXTHREAD               = 0x5
	LC_LOADFVMLIB               = 0x6
	LC_IDFVMLIB                 = 0x7
	LC_IDENT                    = 0x8 // not implemented
	LC_FVMFILE                  = 0x9 // not implemented
	LC_PREPAGE                  = 0xa // not implemented
	LC_DYSYMTAB                 = 0xb
	LC_LOAD_DYLIB               = 0xc
	LC_ID_DYLIB                 = 0xd
	LC_LOAD_DYLINKER            = 0xe
	LC_ID_DYLINKER              = 0xf
	LC_PREBOUND_DYLIB           = 0x10
	LC_ROUTINES                 = 0x11
	LC_SUB_FRAMEWORK            = 0x12
	LC_SUB_UMBRELLA             = 0x13
	LC_SUB_CLIENT               = 0x14
	LC_SUB_LIBRARY              = 0x15
	LC_TWOLEVEL_HINTS           = 0x16
	LC_PREBIND_CKSUM            = 0x17 // not implemented
	LC_LOAD_WEAK_DYLIB          = 0x80000018
	LC_SEGMENT_64               = 0x19
	LC_ROUTINES_64              = 0x1a
	LC_UUID                     = 0x1b
	LC_RPATH                    = 0x8000001c
	LC_CODE_SIGNATURE           = 0x1d
	LC_SEGMENT_SPLIT_INFO       = 0x1e
	LC_REEXPORT_DYLIB           = 0x8000001f
	LC_LAZY_LOAD_DYLIB          = 0x20
	LC_ENCRYPTION_INFO          = 0x21
	LC_DYLD_INFO                = 0x22
	LC_DYLD_INFO_ONLY           = 0x80000022
	LC_LOAD_UPWARD_DYLIB        = 0x80000023
	LC_VERSION_MIN_MACOSX       = 0x24
	LC_VERSION_MIN_IPHONEOS     = 0x25
	LC_FUNCTION_STARTS          = 0x26
	LC_DYLD_ENVIRONMENT         = 0x27
	LC_MAIN                     = 0x80000028
	LC_DATA_IN_CODE             = 0x29
	LC_SOURCE_VERSION           = 0x2A
	LC_DYLIB_CODE_SIGN_DRS      = 0x2B
	LC_ENCRYPTION_INFO_64       = 0x2C
	LC_LINKER_OPTION            = 0x2D
	LC_LINKER_OPTIMIZATION_HINT = 0x2E
	LC_VERSION_MIN_TVOS         = 0x2F
	LC_VERSION_MIN_WATCHOS      = 0x30
	LC_NOTE                     = 0x31 // not implemented
	LC_BUILD_VERSION            = 0x32
)

var loadCommands = scalar.UToSymStr{
	LC_REQ_DYLD:                 "LC_REQ_DYLD",
	LC_SEGMENT:                  "LC_SEGMENT",
	LC_SYMTAB:                   "LC_SYMTAB",
	LC_SYMSEG:                   "LC_SYMSEG",
	LC_THREAD:                   "LC_THREAD",
	LC_UNIXTHREAD:               "LC_UNIXTHREAD",
	LC_LOADFVMLIB:               "LC_LOADFVMLIB",
	LC_IDFVMLIB:                 "LC_IDFVMLIB",
	LC_IDENT:                    "LC_IDENT",
	LC_FVMFILE:                  "LC_FVMFILE",
	LC_PREPAGE:                  "LC_PREPAGE",
	LC_DYSYMTAB:                 "LC_DYSYMTAB",
	LC_LOAD_DYLIB:               "LC_LOAD_DYLIB",
	LC_ID_DYLIB:                 "LC_ID_DYLIB",
	LC_LOAD_DYLINKER:            "LC_LOAD_DYLINKER",
	LC_ID_DYLINKER:              "LC_ID_DYLINKER",
	LC_PREBOUND_DYLIB:           "LC_PREBOUND_DYLIB",
	LC_ROUTINES:                 "LC_ROUTINES",
	LC_SUB_FRAMEWORK:            "LC_SUB_FRAMEWORK",
	LC_SUB_UMBRELLA:             "LC_SUB_UMBRELLA",
	LC_SUB_CLIENT:               "LC_SUB_CLIENT",
	LC_SUB_LIBRARY:              "LC_SUB_LIBRARY",
	LC_TWOLEVEL_HINTS:           "LC_TWOLEVEL_HINTS",
	LC_PREBIND_CKSUM:            "LC_PREBIND_CKSUM",
	LC_LOAD_WEAK_DYLIB:          "LC_LOAD_WEAK_DYLIB",
	LC_SEGMENT_64:               "LC_SEGMENT_64",
	LC_ROUTINES_64:              "LC_ROUTINES_64",
	LC_UUID:                     "LC_UUID",
	LC_RPATH:                    "LC_RPATH",
	LC_CODE_SIGNATURE:           "LC_CODE_SIGNATURE",
	LC_SEGMENT_SPLIT_INFO:       "LC_SEGMENT_SPLIT_INFO",
	LC_REEXPORT_DYLIB:           "LC_REEXPORT_DYLIB",
	LC_LAZY_LOAD_DYLIB:          "LC_LAZY_LOAD_DYLIB",
	LC_ENCRYPTION_INFO:          "LC_ENCRYPTION_INFO",
	LC_DYLD_INFO:                "LC_DYLD_INFO",
	LC_DYLD_INFO_ONLY:           "LC_DYLD_INFO_ONLY",
	LC_LOAD_UPWARD_DYLIB:        "LC_LOAD_UPWARD_DYLIB",
	LC_VERSION_MIN_MACOSX:       "LC_VERSION_MIN_MACOSX",
	LC_VERSION_MIN_IPHONEOS:     "LC_VERSION_MIN_IPHONEOS",
	LC_FUNCTION_STARTS:          "LC_FUNCTION_STARTS",
	LC_DYLD_ENVIRONMENT:         "LC_DYLD_ENVIRONMENT",
	LC_MAIN:                     "LC_MAIN",
	LC_DATA_IN_CODE:             "LC_DATA_IN_CODE",
	LC_SOURCE_VERSION:           "LC_SOURCE_VERSION",
	LC_DYLIB_CODE_SIGN_DRS:      "LC_DYLIB_CODE_SIGN_DRS",
	LC_ENCRYPTION_INFO_64:       "LC_ENCRYPTION_INFO_64",
	LC_LINKER_OPTION:            "LC_LINKER_OPTION",
	LC_LINKER_OPTIMIZATION_HINT: "LC_LINKER_OPTIMIZATION_HINT",
	LC_VERSION_MIN_TVOS:         "LC_VERSION_MIN_TVOS",
	LC_VERSION_MIN_WATCHOS:      "LC_VERSION_MIN_WATCHOS",
	LC_NOTE:                     "LC_NOTE",
	LC_BUILD_VERSION:            "LC_BUILD_VERSION",
}

func machoDecode(d *decode.D, in interface{}) interface{} {
	var archBits int
	var ncmds uint64
	magicBuffer := d.U32LE()

	if magicBuffer == MH_MAGIC || magicBuffer == MH_MAGIC_64 {
		d.Endian = decode.LittleEndian
		if magicBuffer == MH_MAGIC {
			archBits = 32
		} else {
			archBits = 64
		}
	} else if magicBuffer == MH_CIGAM || magicBuffer == MH_CIGAM_64 {
		d.Endian = decode.BigEndian
		if magicBuffer == MH_CIGAM {
			archBits = 32
		} else {
			archBits = 64
		}
	} else {
		d.Fatalf("Invalid magic field")
	}

	d.SeekAbs(0)
	d.FieldStruct(fmt.Sprintf("mach_header_%d", archBits), func(d *decode.D) {
		d.FieldU32("magic", scalar.Hex, classBits, endianNames)
		d.FieldS32("cputype", cpuTypes)
		d.FieldS32("cpusubtype")
		// TODO ask about how to symmap this as it depends on a pair of values
		d.FieldU32("filetype") // TODO expand this
		ncmds = d.FieldU32("ncdms")
		d.FieldU32("sizeofncdms")
		d.FieldU32("flags") // TODO expand flags
		if archBits == 64 {
			d.FieldRawLen("reserved", 4*8, d.BitBufIsZero())
		}
	})
	ncmdsIdx := 0
	d.FieldStructArrayLoop("load_command", "load_command", func() bool {
		return ncmdsIdx < int(ncmds)
	}, func(d *decode.D) {
		cmd := d.FieldU32("cmd", loadCommands)
		cmdsize := d.FieldU32("cmdsize")
		switch cmd {
		case LC_UUID:
			d.FieldStruct("uuid_command", func(d *decode.D) {
				d.FieldRawLen("uuid", 16*8)
			})
		case LC_SEGMENT, LC_SEGMENT_64:
			// nsect := (cmdsize - uint64(archBits)) / uint64(archBits)
			var nsects uint64
			if archBits == 32 {
				d.FieldStruct("segment_command", func(d *decode.D) {
					d.FieldRawLen("segname", 16*8)
					d.FieldU32("vmaddr")
					d.FieldU32("vmsize")
					d.FieldU32("fileoff")
					d.FieldU32("tfilesize")
					d.FieldS32("initprot")
					d.FieldS32("maxprot")
					nsects = d.FieldU32("nsects")
					d.FieldU32("flags") // TODO expand flags
				})
			} else {
				d.FieldStruct("segment_command_64", func(d *decode.D) {
					d.FieldStrFn("segname", func(d *decode.D) string {
						return string(d.BytesLen(16))
					})
					d.FieldU64("vmaddr")
					d.FieldU64("vmsize")
					d.FieldU64("fileoff")
					d.FieldU64("tfilesize")
					d.FieldS32("initprot")
					d.FieldS32("maxprot")
					nsects = d.FieldU32("nsects")
					d.FieldU32("flags") // TODO expand flags
				})
			}
			var nsectIdx uint64
			d.FieldStructArrayLoop("section", "section", func() bool {
				return nsectIdx < nsects
			},
				func(d *decode.D) {
					d.FieldStrFn("sectname", func(d *decode.D) string {
						return string(d.BytesLen(16))
					})
					d.FieldStrFn("segname", func(d *decode.D) string {
						return string(d.BytesLen(16))
					})
					if archBits == 32 {
						d.FieldU32("address")
						d.FieldU32("size")
					} else {
						d.FieldU64("address")
						d.FieldU64("size")
					}
					d.FieldU32("offset")
					d.FieldU32("align")
					d.FieldU32("reloff")
					d.FieldU32("nreloc")
					d.FieldU32("flags") // TODO expand flags
					d.FieldU32("reserved1")
					d.FieldU32("reserved2")
					if archBits == 64 {
						d.FieldU32("reserved3")
					}
					nsectIdx++
				})
		case LC_TWOLEVEL_HINTS:
			d.FieldU32("offset")
			d.FieldU32("nhints")
		case LC_LOAD_DYLIB, LC_ID_DYLIB, LC_LOAD_UPWARD_DYLIB, LC_LOAD_WEAK_DYLIB, LC_LAZY_LOAD_DYLIB, LC_REEXPORT_DYLIB:
			d.FieldStruct("dylib_command", func(d *decode.D) {
				offset := d.FieldU32("offset")
				d.FieldU32("timestamp") // TODO human readable
				d.FieldU32("current_version")
				d.FieldU32("compatibility_version")
				d.FieldStrFn("name", func(d *decode.D) string {
					return string(d.BytesLen(int(cmdsize) - int(offset)))
				})
			})
		case LC_LOAD_DYLINKER, LC_ID_DYLINKER, LC_DYLD_ENVIRONMENT:
			offset := d.FieldU32("offset")
			d.FieldStrFn("name", func(d *decode.D) string {
				return string(d.BytesLen(int(cmdsize) - int(offset)))
			})
		case LC_RPATH:
			offset := d.FieldU32("offset")
			d.FieldStrFn("name", func(d *decode.D) string {
				return string(d.BytesLen(int(cmdsize) - int(offset)))
			})
		case LC_PREBOUND_DYLIB:
			offset := d.FieldU32("offset")
			nmodules := d.FieldU32("nmodules")
			d.FieldBitBufFn("linked_modules", func(d *decode.D) *bitio.Buffer {
				return d.RawLen(int64((nmodules / 8) + (nmodules % 8)))
			}) // TODO this needs better representation
			d.FieldStrFn("name", func(d *decode.D) string {
				return string(d.BytesLen(int(cmdsize) - int(offset)))
			}) // TODO visualize this bitset
		case LC_THREAD, LC_UNIXTHREAD:
			d.FieldU32("flavor")
			count := d.FieldU32("count")
			d.FieldRawLen("state", int64(count*32))
			// better visualization needed for this specific for major architectures
		case LC_ROUTINES, LC_ROUTINES_64:
			if archBits == 32 {
				d.FieldU32("init_address")
				d.FieldU32("init_module")
				d.FieldU32("reserved1", d.BitBufIsZero())
				d.FieldU32("reserved2", d.BitBufIsZero())
				d.FieldU32("reserved3", d.BitBufIsZero())
				d.FieldU32("reserved4", d.BitBufIsZero())
				d.FieldU32("reserved5", d.BitBufIsZero())
				d.FieldU32("reserved6", d.BitBufIsZero())
			} else {
				d.FieldU64("init_address")
				d.FieldU64("init_module")
				d.FieldU64("reserved1", d.BitBufIsZero())
				d.FieldU64("reserved2", d.BitBufIsZero())
				d.FieldU64("reserved3", d.BitBufIsZero())
				d.FieldU64("reserved4", d.BitBufIsZero())
				d.FieldU64("reserved5", d.BitBufIsZero())
				d.FieldU64("reserved6", d.BitBufIsZero())
			}
		case LC_SUB_UMBRELLA, LC_SUB_LIBRARY, LC_SUB_CLIENT, LC_SUB_FRAMEWORK:
			offset := d.FieldU32("offset")
			d.FieldStrFn("name", func(d *decode.D) string {
				return string(d.BytesLen(int(cmdsize) - int(offset)))
			})
		case LC_SYMTAB:
			d.FieldU32("symoff")
			d.FieldU32("nsyms")
			d.FieldU32("stroff")
			d.FieldU32("strsize")
		case LC_DYSYMTAB:
			// TODO
			// if archBits == 32 {
			// 	d.FieldStruct("nlist", func(d *decode.D) {
			// 		d.FieldStruct("n_un", func(d *decode.D) {
			// 			d.FieldS32("n_strx")
			// 		})
			// 		d.FieldU8("n_type")
			// 		d.FieldU8("n_sect")
			// 		d.FieldU16("n_desc")
			// 		d.FieldU32("n_value")
			// 	})
			// } else {
			// 	d.FieldStruct("nlist", func(d *decode.D) {
			// 		d.FieldStruct("n_un", func(d *decode.D) {
			// 			d.FieldS32("n_strx")
			// 		})
			// 		d.FieldU8("n_type")
			// 		d.FieldU8("n_sect")
			// 		d.FieldU16("n_desc")
			// 		d.FieldU64("n_value")
			// 	})
			// }
			d.FieldU32("ilocalsym")
			d.FieldU32("nlocalsym")
			d.FieldU32("iextdefsym")
			d.FieldU32("nextdefsym")
			d.FieldU32("iundefsym")
			d.FieldU32("nundefsym")
			d.FieldU32("tocoff")
			d.FieldU32("ntoc")
			d.FieldU32("modtaboff")
			d.FieldU32("nmodtab")
			d.FieldU32("extrefsymoff")
			d.FieldU32("nextrefsyms")
			d.FieldU32("indirectsymoff")
			d.FieldU32("nindirectsyms")

			d.FieldU32("extreloff")
			d.FieldU32("nextrel")
			d.FieldU32("locreloff")
			d.FieldU32("nlocrel")
		case LC_BUILD_VERSION:
			d.FieldU32("platform")
			d.FieldU32("minos")
			d.FieldU32("sdk")
			ntools := d.FieldU32("ntools")
			var ntoolsIdx uint64
			d.FieldStructArrayLoop("tools", "tools", func() bool {
				return ntoolsIdx < ntools
			}, func(d *decode.D) {
				d.FieldU32("tool")
				d.FieldU32("version")
			})
		case LC_CODE_SIGNATURE, LC_SEGMENT_SPLIT_INFO, LC_FUNCTION_STARTS, LC_DATA_IN_CODE, LC_DYLIB_CODE_SIGN_DRS, LC_LINKER_OPTIMIZATION_HINT:
			d.FieldStruct("linkedit_data", func(d *decode.D) {
				d.FieldU32("off")
				d.FieldU32("size")
			})
		case LC_VERSION_MIN_IPHONEOS, LC_VERSION_MIN_MACOSX, LC_VERSION_MIN_TVOS, LC_VERSION_MIN_WATCHOS:
			d.FieldU32("version")
			d.FieldU32("sdk")
		case LC_DYLD_INFO, LC_DYLD_INFO_ONLY:
			d.FieldStruct("dyld_info", func(d *decode.D) {
				d.FieldU32("rebase_off")
				d.FieldU32("rebase_size")
				d.FieldU32("bind_off")
				d.FieldU32("bind_size")
				d.FieldU32("weak_bind_off")
				d.FieldU32("weak_bind_size")
				d.FieldU32("lazy_bind_off")
				d.FieldU32("lazy_bind_size")
				d.FieldU32("export_off")
				d.FieldU32("export_size")
			})
		case LC_MAIN:
			d.FieldStruct("entrypoint", func(d *decode.D) {
				d.FieldU32("entryoff")
				d.FieldU32("stacksize")
			})
		case LC_SOURCE_VERSION:
			d.FieldStruct("source_version_tag", func(d *decode.D) {
				d.FieldU64("tag")
			})
		case LC_LINKER_OPTION:
			d.FieldStruct("linker_option", func(d *decode.D) {
				count := d.FieldU32("count")
				d.FieldUTF8NullFixedLen("option", int(count))
			})
		case LC_ENCRYPTION_INFO, LC_ENCRYPTION_INFO_64:
			d.FieldStruct(fmt.Sprintf("encryption_info_%d", archBits), func(d *decode.D) {
				d.FieldU32("offset")
				d.FieldU32("size")
				d.FieldU32("id")
			})
		case LC_IDFVMLIB, LC_LOADFVMLIB:
			d.FieldStruct("fvmlib", func(d *decode.D) {
				offset := d.FieldU32("offset")
				d.FieldU32("minor_version")
				d.FieldU32("header_addr")
				d.FieldStrFn("name", func(d *decode.D) string {
					return string(d.BytesLen(int(cmdsize) - int(offset)))
				})
			})
		default:
			if _, ok := loadCommands[cmd]; !ok {
				d.FieldRawLen("unknown", int64((cmdsize-8)*8))
			}
		}
		ncmdsIdx++
	})
	return nil
}
