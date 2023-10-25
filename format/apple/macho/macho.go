package macho

// https://github.com/aidansteele/osx-abi-macho-file-format-reference

import (
	"embed"
	"strings"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed macho.md
var machoFS embed.FS

func init() {
	interp.RegisterFormat(
		format.MachO,
		&decode.Format{
			Description: "Mach-O macOS executable",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    machoDecode,
		})
	interp.RegisterFS(machoFS)
}

func strIndexNull(idx int, s string) string {
	if idx > len(s) {
		return ""
	}
	i := strings.IndexByte(s[idx:], 0)
	if i == -1 {
		return ""
	}
	return s[idx : idx+i]
}

type strTable string

func (m strTable) MapUint(s scalar.Uint) (scalar.Uint, error) {
	s.Sym = strIndexNull(int(s.Actual), string(m))
	return s, nil
}

const (
	MH_MAGIC    = 0xfeed_face
	MH_CIGAM    = 0xcefa_edfe
	MH_MAGIC_64 = 0xfeed_facf
	MH_CIGAM_64 = 0xcffa_edfe
)

var magicSymMapper = scalar.UintMap{
	MH_MAGIC:    scalar.Uint{Sym: "32le", Description: "32-bit little endian"},
	MH_CIGAM:    scalar.Uint{Sym: "32be", Description: "32-bit big endian"},
	MH_MAGIC_64: scalar.Uint{Sym: "64le", Description: "64-bit little endian"},
	MH_CIGAM_64: scalar.Uint{Sym: "64be", Description: "64-bit big endian"},
}

var cpuTypes = scalar.UintMapSymStr{
	0xff_ff_ff_ff: "any",
	1:             "vax",
	2:             "romp",
	4:             "ns32032",
	5:             "ns32332",
	6:             "mc680x0",
	7:             "x86",
	8:             "mips",
	9:             "ns32532",
	10:            "mc98000",
	11:            "hppa",
	12:            "arm",
	13:            "mc88000",
	14:            "sparc",
	15:            "i860",
	16:            "i860_little",
	17:            "rs6000",
	18:            "powerpc",
	0x1000007:     "x86_64",
	0x100000c:     "arm64",
	0x1000013:     "powerpc64",
	255:           "veo",
}

func intelSubTypeHelper(f, m uint64) uint64 {
	return f + (m << 4)
}

var cpuSubTypes = map[uint64]scalar.UintMapSymStr{
	0xff_ff_ff_ff: {
		0xff_ff_ff_ff: "multiple",
	},
	1: {
		0xff_ff_ff_ff: "multiple",
		0:             "vax_all",
		1:             "vax780",
		2:             "vax785",
		3:             "vax750",
		4:             "vax730",
		5:             "uvaxi",
		6:             "uvaxii",
		7:             "vax8200",
		8:             "vax8500",
		9:             "vax8600",
		10:            "vax8650",
		11:            "vax8800",
		12:            "uvaxiii",
	},
	6: {
		0xff_ff_ff_ff: "multiple",
		1:             "mc680x0_all", // 1: mc68030
		2:             "mc68040",
		3:             "mc68030_only",
	},
	7: {
		0xff_ff_ff_ff:             "multiple",
		intelSubTypeHelper(3, 0):  "i386_all", // i386
		intelSubTypeHelper(4, 0):  "i486",
		intelSubTypeHelper(4, 8):  "486sx",
		intelSubTypeHelper(5, 0):  "pent",
		intelSubTypeHelper(6, 1):  "pentpro",
		intelSubTypeHelper(6, 3):  "pentii_m3",
		intelSubTypeHelper(6, 5):  "pentii_m5",
		intelSubTypeHelper(7, 6):  "celeron",
		intelSubTypeHelper(7, 7):  "celeron_mobile",
		intelSubTypeHelper(8, 0):  "pentium_3",
		intelSubTypeHelper(8, 1):  "pentium_3_m",
		intelSubTypeHelper(8, 2):  "pentium_3_xeon",
		intelSubTypeHelper(9, 0):  "pentium_m",
		intelSubTypeHelper(10, 0): "pentium_4",
		intelSubTypeHelper(10, 1): "pentium_4_m",
		intelSubTypeHelper(11, 0): "itanium",
		intelSubTypeHelper(11, 1): "itanium_2",
		intelSubTypeHelper(12, 0): "xeon",
		intelSubTypeHelper(12, 1): "xeon_2",
	},
	8: {
		0xff_ff_ff_ff: "multiple",
		0:             "mips_all",
		1:             "mips_r2300",
		2:             "mips_r2600",
		3:             "mips_r2800",
		4:             "mips_r2000a",
		5:             "mips_r2000",
		6:             "mips_r3000a",
		7:             "mips_r3000",
	},
	10: {
		0xff_ff_ff_ff: "multiple",
		0:             "mc98000_all",
		1:             "mc98001",
	},
	11: {
		0xff_ff_ff_ff: "multiple",
		0:             "hppa_all",
		1:             "hppa_7100",
		2:             "hppa_7100_lc",
	},
	12: {
		0xff_ff_ff_ff: "multiple",
		0:             "arm_all",
		5:             "arm_v4t",
		6:             "arm_v6",
		7:             "arm_v5tej",
		8:             "arm_xscale",
		9:             "arm_v7",
		10:            "arm_v7f",
		11:            "arm_v7s",
		12:            "arm_v7k",
		13:            "arm_v8",
		14:            "arm_v6m",
		15:            "arm_v7m",
		16:            "arm_v7em",
	},
	13: {
		0xff_ff_ff_ff: "multiple",
		0:             "mc88000_all",
		1:             "mc88100",
		2:             "mc88110",
	},
	14: {
		0xff_ff_ff_ff: "multiple",
		0:             "sparc_all",
	},
	15: {
		0xff_ff_ff_ff: "multiple",
		0:             "i860_all",
		1:             "i860_a860",
	},
	18: {
		0xff_ff_ff_ff: "multiple",
		0:             "powerpc_all",
		1:             "powerpc_601",
		2:             "powerpc_602",
		3:             "powerpc_603",
		4:             "powerpc_603e",
		5:             "powerpc_603ev",
		6:             "powerpc_604",
		7:             "powerpc_604e",
		8:             "powerpc_620",
		9:             "powerpc_750",
		10:            "powerpc_7400",
		11:            "powerpc_7450",
		100:           "powerpc_970",
	},
	0x1000012: {
		0xff_ff_ff_ff: "multiple",
		0:             "arm64_all",
		1:             "arm64_v8",
		2:             "arm64_e",
	},
}

var fileTypes = scalar.UintMapSymStr{
	0x1: "object",
	0x2: "execute",
	0x3: "fvmlib",
	0x4: "core",
	0x5: "preload",
	0x6: "dylib",
	0x7: "dylinker",
	0x8: "bundle",
	0x9: "dylib_stub",
	0xa: "dsym",
	0xb: "kext_bundle",
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
	LC_SOURCE_VERSION           = 0x2a
	LC_DYLIB_CODE_SIGN_DRS      = 0x2b
	LC_ENCRYPTION_INFO_64       = 0x2c
	LC_LINKER_OPTION            = 0x2d
	LC_LINKER_OPTIMIZATION_HINT = 0x2e
	LC_VERSION_MIN_TVOS         = 0x2f
	LC_VERSION_MIN_WATCHOS      = 0x30
	LC_NOTE                     = 0x31 // not implemented
	LC_BUILD_VERSION            = 0x32
)

var loadCommands = scalar.UintMapSymStr{
	LC_REQ_DYLD:                 "req_dyld",
	LC_SEGMENT:                  "segment",
	LC_SYMTAB:                   "symtab",
	LC_SYMSEG:                   "symseg",
	LC_THREAD:                   "thread",
	LC_UNIXTHREAD:               "unixthread",
	LC_LOADFVMLIB:               "loadfvmlib",
	LC_IDFVMLIB:                 "idfvmlib",
	LC_IDENT:                    "ident",
	LC_FVMFILE:                  "fvmfile",
	LC_PREPAGE:                  "prepage",
	LC_DYSYMTAB:                 "dysymtab",
	LC_LOAD_DYLIB:               "load_dylib",
	LC_ID_DYLIB:                 "id_dylib",
	LC_LOAD_DYLINKER:            "load_dylinker",
	LC_ID_DYLINKER:              "id_dylinker",
	LC_PREBOUND_DYLIB:           "prebound_dylib",
	LC_ROUTINES:                 "routines",
	LC_SUB_FRAMEWORK:            "sub_framework",
	LC_SUB_UMBRELLA:             "sub_umbrella",
	LC_SUB_CLIENT:               "sub_client",
	LC_SUB_LIBRARY:              "sub_library",
	LC_TWOLEVEL_HINTS:           "twolevel_hints",
	LC_PREBIND_CKSUM:            "prebind_cksum",
	LC_LOAD_WEAK_DYLIB:          "load_weak_dylib",
	LC_SEGMENT_64:               "segment_64",
	LC_ROUTINES_64:              "routines_64",
	LC_UUID:                     "uuid",
	LC_RPATH:                    "rpath",
	LC_CODE_SIGNATURE:           "code_signature",
	LC_SEGMENT_SPLIT_INFO:       "segment_split_info",
	LC_REEXPORT_DYLIB:           "reexport_dylib",
	LC_LAZY_LOAD_DYLIB:          "lazy_load_dylib",
	LC_ENCRYPTION_INFO:          "encryption_info",
	LC_DYLD_INFO:                "dyld_info",
	LC_DYLD_INFO_ONLY:           "dyld_info_only",
	LC_LOAD_UPWARD_DYLIB:        "load_upward_dylib",
	LC_VERSION_MIN_MACOSX:       "version_min_macosx",
	LC_VERSION_MIN_IPHONEOS:     "version_min_iphoneos",
	LC_FUNCTION_STARTS:          "function_starts",
	LC_DYLD_ENVIRONMENT:         "dyld_environment",
	LC_MAIN:                     "main",
	LC_DATA_IN_CODE:             "data_in_code",
	LC_SOURCE_VERSION:           "source_version",
	LC_DYLIB_CODE_SIGN_DRS:      "dylib_code_sign_drs",
	LC_ENCRYPTION_INFO_64:       "encryption_info_64",
	LC_LINKER_OPTION:            "linker_option",
	LC_LINKER_OPTIMIZATION_HINT: "linker_optimization_hint",
	LC_VERSION_MIN_TVOS:         "version_min_tvos",
	LC_VERSION_MIN_WATCHOS:      "version_min_watchos",
	LC_NOTE:                     "note",
	LC_BUILD_VERSION:            "build_version",
}

var sectionTypes = scalar.UintMapSymStr{
	0x0:  "regular",
	0x1:  "zerofill",
	0x2:  "cstring_literals",
	0x3:  "4byte_literals",
	0x4:  "8byte_literals",
	0x5:  "literal_pointers",
	0x6:  "non_lazy_symbol_pointers",
	0x7:  "lazy_symbol_pointers",
	0x8:  "symbol_stubs",
	0x9:  "mod_init_func_pointers",
	0xa:  "mod_term_func_pointers",
	0xb:  "coalesced",
	0xc:  "gb_zerofill",
	0xd:  "interposing",
	0xe:  "16byte_literals",
	0xf:  "dtrace_dof",
	0x10: "lazy_dylib_symbol_pointers",
	0x11: "thread_local_regular",
	0x12: "thread_local_zerofill",
	0x13: "thread_local_variables",
	0x14: "thread_local_variable_pointers",
	0x15: "thread_local_init_function_pointers",
}

func machoDecode(d *decode.D) any {
	var archBits int
	var cpuType uint64
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
		d.Fatalf("invalid magic")
	}

	d.SeekRel(-4 * 8)
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldValueSint("arch_bits", int64(archBits))
		d.FieldU32("magic", magicSymMapper, scalar.UintHex)
		d.FieldValueUint("bits", uint64(archBits))
		cpuType = d.FieldU32("cputype", cpuTypes, scalar.UintHex)
		d.FieldU32("cpusubtype", cpuSubTypes[cpuType], scalar.UintHex)
		d.FieldU32("filetype", fileTypes)
		ncmds = d.FieldU32("ncdms")
		d.FieldU32("sizeofncdms")
		d.FieldStruct("flags", parseMachHeaderFlags)
		if archBits == 64 {
			d.FieldRawLen("reserved", 4*8, d.BitBufIsZero())
		}
	})
	loadCommandsNext := d.Pos()
	d.FieldArray("load_commands", func(d *decode.D) {
		for i := uint64(0); i < ncmds; i++ {
			d.FieldStruct("load_command", func(d *decode.D) {
				d.SeekAbs(loadCommandsNext)

				cmd := d.FieldU32("cmd", loadCommands, scalar.UintHex)
				cmdSize := d.FieldU32("cmdsize")
				if cmdSize == 0 {
					d.Fatalf("cmdSize is zero")
				}

				loadCommandsNext += int64(cmdSize) * 8

				switch cmd {
				case LC_UUID:
					d.FieldStruct("uuid_command", func(d *decode.D) {
						d.FieldRawLen("uuid", 16*8)
					})
				case LC_SEGMENT,
					LC_SEGMENT_64:
					// nsect := (cmdsize - uint64(archBits)) / uint64(archBits)

					var vmaddr int64
					var fileoff int64

					var nsects uint64
					d.FieldStruct("segment_command", func(d *decode.D) {
						d.FieldValueSint("arch_bits", int64(archBits))
						d.FieldUTF8NullFixedLen("segname", 16) // OPCODE_DECODER segname==__TEXT
						if archBits == 32 {
							vmaddr = int64(d.FieldU32("vmaddr", scalar.UintHex))
							d.FieldU32("vmsize")
							fileoff = int64(d.FieldU32("fileoff", scalar.UintHex))
							d.FieldU32("tfilesize")
						} else {
							vmaddr = int64(d.FieldU64("vmaddr", scalar.UintHex))
							d.FieldU64("vmsize")
							fileoff = int64(d.FieldU64("fileoff", scalar.UintHex))
							d.FieldU64("tfilesize")
						}
						d.FieldS32("initprot")
						d.FieldS32("maxprot")
						nsects = d.FieldU32("nsects")
						d.FieldStruct("flags", parseSegmentFlags)
					})
					d.FieldArray("sections", func(d *decode.D) {
						for i := uint64(0); i < nsects; i++ {
							d.FieldStruct("section", func(d *decode.D) {
								// OPCODE_DECODER sectname==__text
								sectName := d.FieldUTF8NullFixedLen("sectname", 16)
								d.FieldUTF8NullFixedLen("segname", 16)
								var size uint64
								if archBits == 32 {
									d.FieldU32("address", scalar.UintHex)
									size = d.FieldU32("size")
								} else {
									d.FieldU64("address", scalar.UintHex)
									size = d.FieldU64("size")
								}
								offset := d.FieldU32("offset", scalar.UintHex)
								d.FieldU32("align")
								d.FieldU32("reloff")
								d.FieldU32("nreloc")
								// get section type
								d.FieldStruct("flags", parseSectionFlags)
								d.FieldU32("reserved1")
								d.FieldU32("reserved2")
								if archBits == 64 {
									d.FieldU32("reserved3")
								}

								switch sectName {
								case "__bss", // uninitialized data
									"__common": // allocated by linker
									// skip, no data from file
									// TODO: more?
								default:
									d.RangeFn(int64(offset)*8, int64(size)*8, func(d *decode.D) {
										switch sectName {
										case "__cstring":
											d.FieldArray("cstrings", func(d *decode.D) {
												for !d.End() {
													d.FieldUTF8Null("cstring")
												}
											})
										case "__ustring":
											d.FieldArray("ustrings", func(d *decode.D) {
												for !d.End() {
													// TODO: always LE?
													d.FieldUTF16LENull("ustring")
												}
											})
										case "__cfstring":
											d.FieldArray("cfstrings", func(d *decode.D) {
												for !d.End() {
													d.FieldStruct("cfstring", func(d *decode.D) {
														// https://github.com/llvm-mirror/clang/blob/aa231e4be75ac4759c236b755c57876f76e3cf05/lib/CodeGen/CodeGenModule.cpp#L4708
														const flagUTF8 = 0x07c8
														const flagUTF16 = 0x07d0

														d.FieldU("isa_vmaddr", archBits)
														flag := d.FieldU("flags", archBits, scalar.UintHex, scalar.UintMapSymStr{
															flagUTF8:  "utf8",
															flagUTF16: "utf16",
														})
														dataPtr := int64(d.FieldU("data_ptr", archBits, scalar.UintHex))
														length := int64(d.FieldU("length", archBits))

														offset := ((dataPtr - vmaddr) + fileoff) * 8
														switch flag {
														case flagUTF8:
															d.RangeFn(offset, length*8, func(d *decode.D) { d.FieldUTF8("string", int(length)) })
														case flagUTF16:
															// TODO: endian?
															d.RangeFn(offset, length*8*2, func(d *decode.D) { d.FieldUTF16("string", int(length*2)) })
														}
													})
												}
											})
										default:
											d.FieldRawLen("data", d.BitsLeft())
										}
									})
								}
							})
						}
					})
				case LC_TWOLEVEL_HINTS:
					d.FieldU32("offset", scalar.UintHex)
					d.FieldU32("nhints")
				case LC_LOAD_DYLIB,
					LC_ID_DYLIB,
					LC_LOAD_UPWARD_DYLIB,
					LC_LOAD_WEAK_DYLIB,
					LC_LAZY_LOAD_DYLIB,
					LC_REEXPORT_DYLIB:
					d.FieldStruct("dylib_command", func(d *decode.D) {
						offset := d.FieldU32("offset", scalar.UintHex)
						d.FieldU32("timestamp", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
						d.FieldU32("current_version")
						d.FieldU32("compatibility_version")
						d.FieldUTF8NullFixedLen("name", int(cmdSize)-int(offset))
					})
				case LC_LOAD_DYLINKER,
					LC_ID_DYLINKER,
					LC_DYLD_ENVIRONMENT:
					offset := d.FieldU32("offset", scalar.UintHex)
					d.FieldUTF8NullFixedLen("name", int(cmdSize)-int(offset))
				case LC_RPATH:
					offset := d.FieldU32("offset", scalar.UintHex)
					d.FieldUTF8NullFixedLen("name", int(cmdSize)-int(offset))
				case LC_PREBOUND_DYLIB:
					// https://github.com/aidansteele/osx-abi-macho-file-format-reference#prebound_dylib_command
					d.U32() // name_offset
					nmodules := d.FieldU32("nmodules")
					d.U32() // linked_modules_offset
					d.FieldUTF8Null("name")
					d.FieldBitBufFn("linked_modules", func(d *decode.D) bitio.ReaderAtSeeker {
						return d.RawLen(int64((nmodules / 8) + (nmodules % 8)))
					})
				case LC_THREAD,
					LC_UNIXTHREAD:
					d.FieldU32("flavor")
					count := d.FieldU32("count")
					switch cpuType {
					case 0x7:
						d.FieldStruct("state", threadStateI386Decode)
					case 0xC:
						d.FieldStruct("state", threadStateARM32Decode)
					case 0x13:
						d.FieldStruct("state", threadStatePPC32Decode)
					case 0x1000007:
						d.FieldStruct("state", threadStateX8664Decode)
					case 0x100000C:
						d.FieldStruct("state", threadStateARM64Decode)
					case 0x1000013:
						d.FieldStruct("state", threadStatePPC64Decode)
					default:
						d.FieldRawLen("state", int64(count*32))
					}
				case LC_ROUTINES,
					LC_ROUTINES_64:
					if archBits == 32 {
						d.FieldU32("init_address", scalar.UintHex)
						d.FieldU32("init_module")
						d.FieldU32("reserved1")
						d.FieldU32("reserved2")
						d.FieldU32("reserved3")
						d.FieldU32("reserved4")
						d.FieldU32("reserved5")
						d.FieldU32("reserved6")
					} else {
						d.FieldU64("init_address", scalar.UintHex)
						d.FieldU64("init_module")
						d.FieldU64("reserved1")
						d.FieldU64("reserved2")
						d.FieldU64("reserved3")
						d.FieldU64("reserved4")
						d.FieldU64("reserved5")
						d.FieldU64("reserved6")
					}
				case LC_SUB_UMBRELLA,
					LC_SUB_LIBRARY,
					LC_SUB_CLIENT,
					LC_SUB_FRAMEWORK:
					offset := d.FieldU32("offset", scalar.UintHex)
					d.FieldUTF8NullFixedLen("name", int(cmdSize)-int(offset))
				case LC_SYMTAB:
					symOff := d.FieldU32("symoff")
					nSyms := d.FieldU32("nsyms")
					strOff := d.FieldU32("stroff")
					strSize := d.FieldU32("strsize")

					d.RangeFn(int64(strOff)*8, int64(strSize)*8, func(d *decode.D) {
						d.FieldRawLen("str_table", d.BitsLeft())
					})
					symTabTable := strTable(string(d.BytesRange(int64(strOff)*8, int(strSize))))

					d.SeekAbs(int64(symOff) * 8)
					d.FieldArray("symbols", func(d *decode.D) {
						for i := 0; i < int(nSyms); i++ {
							symbolTypeMap := scalar.UintMapSymStr{
								0x0: "undef",
								0x1: "abs",
								0x5: "indr",
								0x6: "pbud",
								0x7: "sect",
							}

							d.FieldStruct("symbol", func(d *decode.D) {
								d.FieldU32("strx", symTabTable)
								d.FieldStruct("type", func(d *decode.D) {
									d.FieldU3("stab")
									d.FieldU1("pext")
									d.FieldU3("type", symbolTypeMap)
									d.FieldU1("ext")
								})
								d.FieldU8("sect")
								d.FieldU16("desc")
								d.FieldU("value", archBits, scalar.UintHex)
							})
						}
					})
				case LC_DYSYMTAB:
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
					d.FieldStructArrayLoop("tools", "tool", func() bool {
						return ntoolsIdx < ntools
					}, func(d *decode.D) {
						d.FieldU32("tool")
						d.FieldU32("version")
						ntoolsIdx++
					})
				case LC_CODE_SIGNATURE,
					LC_SEGMENT_SPLIT_INFO,
					LC_FUNCTION_STARTS,
					LC_DATA_IN_CODE,
					LC_DYLIB_CODE_SIGN_DRS,
					LC_LINKER_OPTIMIZATION_HINT:
					d.FieldStruct("linkedit_data", func(d *decode.D) {
						d.FieldU32("off")
						d.FieldU32("size")
					})
				case LC_VERSION_MIN_IPHONEOS,
					LC_VERSION_MIN_MACOSX,
					LC_VERSION_MIN_TVOS,
					LC_VERSION_MIN_WATCHOS:
					d.FieldU32("version")
					d.FieldU32("sdk")
				case LC_DYLD_INFO,
					LC_DYLD_INFO_ONLY:
					d.FieldStruct("dyld_info", func(d *decode.D) {
						d.FieldU32("rebase_off", scalar.UintHex)
						d.FieldU32("rebase_size")
						d.FieldU32("bind_off", scalar.UintHex)
						d.FieldU32("bind_size")
						d.FieldU32("weak_bind_off", scalar.UintHex)
						d.FieldU32("weak_bind_size")
						d.FieldU32("lazy_bind_off", scalar.UintHex)
						d.FieldU32("lazy_bind_size")
						d.FieldU32("export_off", scalar.UintHex)
						d.FieldU32("export_size")
					})
				case LC_MAIN:
					d.FieldStruct("entrypoint", func(d *decode.D) {
						d.FieldU64("entryoff", scalar.UintHex)
						d.FieldU64("stacksize")
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
				case LC_ENCRYPTION_INFO,
					LC_ENCRYPTION_INFO_64:
					d.FieldStruct("encryption_info", func(d *decode.D) {
						offset := d.FieldU32("offset", scalar.UintHex)
						size := d.FieldU32("size")
						d.FieldU32("id")
						d.RangeFn(int64(offset)*8, int64(size)*8, func(d *decode.D) {
							d.FieldRawLen("data", d.BitsLeft())
						})
					})
					if cmd == LC_ENCRYPTION_INFO_64 {
						// 64 bit align
						d.FieldU32("pad")
					}
				case LC_IDFVMLIB,
					LC_LOADFVMLIB:
					d.FieldStruct("fvmlib", func(d *decode.D) {
						offset := d.FieldU32("offset", scalar.UintHex)
						d.FieldU32("minor_version")
						d.FieldU32("header_addr", scalar.UintHex)
						d.FieldUTF8NullFixedLen("name", int(cmdSize)-int(offset))
					})
				default:
					// ignore
				}
			})
		}
	})

	return nil
}

// TODO: some kind of flags-endian helper?
func parseMachHeaderFlags(d *decode.D) {
	if d.Endian == decode.BigEndian {
		d.FieldRawLen("reserved", 6)
		d.FieldBool("app_extension_safe")
		d.FieldBool("no_heap_execution")

		d.FieldBool("has_tlv_descriptors")
		d.FieldBool("dead_strippable_dylib")
		d.FieldBool("pie")
		d.FieldBool("no_reexported_dylibs")
		d.FieldBool("setuid_safe")
		d.FieldBool("root_safe")
		d.FieldBool("allow_stack_execution")
		d.FieldBool("binds_to_weak")

		d.FieldBool("weak_defines")
		d.FieldBool("canonical")
		d.FieldBool("subsections_via_symbols")
		d.FieldBool("allmodsbound")
		d.FieldBool("prebindable")
		d.FieldBool("nofixprebinding")
		d.FieldBool("nomultidefs")
		d.FieldBool("force_flat")

		d.FieldBool("twolevel")
		d.FieldBool("lazy_init")
		d.FieldBool("split_segs")
		d.FieldBool("prebound")
		d.FieldBool("bindatload")
		d.FieldBool("dyldlink")
		d.FieldBool("incrlink")
		d.FieldBool("noundefs")
	} else {
		d.FieldBool("twolevel")
		d.FieldBool("lazy_init")
		d.FieldBool("split_segs")
		d.FieldBool("prebound")
		d.FieldBool("bindatload")
		d.FieldBool("dyldlink")
		d.FieldBool("incrlink")
		d.FieldBool("noundefs")

		d.FieldBool("weak_defines")
		d.FieldBool("canonical")
		d.FieldBool("subsections_via_symbols")
		d.FieldBool("allmodsbound")
		d.FieldBool("prebindable")
		d.FieldBool("nofixprebinding")
		d.FieldBool("nomultidefs")
		d.FieldBool("force_flat")

		d.FieldBool("has_tlv_descriptors")
		d.FieldBool("dead_strippable_dylib")
		d.FieldBool("pie")
		d.FieldBool("no_reexported_dylibs")
		d.FieldBool("setuid_safe")
		d.FieldBool("root_safe")
		d.FieldBool("allow_stack_execution")
		d.FieldBool("binds_to_weak")

		d.FieldRawLen("reserved", 6)
		d.FieldBool("app_extension_safe")
		d.FieldBool("no_heap_execution")
	}
}

func parseSegmentFlags(d *decode.D) {
	if d.Endian == decode.BigEndian {
		d.FieldRawLen("reserved0", 24)

		d.FieldRawLen("reserved1", 4)
		d.FieldBool("protected_version_1")
		d.FieldBool("noreloc")
		d.FieldBool("fvmlib")
		d.FieldBool("highvm")
	} else {
		d.FieldRawLen("reserved0", 4)
		d.FieldBool("protected_version_1")
		d.FieldBool("noreloc")
		d.FieldBool("fvmlib")
		d.FieldBool("highvm")

		d.FieldRawLen("reserved1", 24)
	}
}

func parseSectionFlags(d *decode.D) {
	if d.Endian == decode.BigEndian {
		d.FieldBool("attr_pure_instructions")
		d.FieldBool("attr_no_toc")
		d.FieldBool("attr_strip_static_syms")
		d.FieldBool("attr_no_dead_strip")
		d.FieldBool("attr_live_support")
		d.FieldBool("attr_self_modifying_code")
		d.FieldBool("attr_debug")
		d.FieldRawLen("reserved0", 1)

		d.FieldRawLen("reserved1", 8)

		d.FieldRawLen("reserved2", 5)
		d.FieldBool("attr_some_instructions")
		d.FieldBool("attr_ext_reloc")
		d.FieldBool("attr_loc_reloc")

		d.FieldU8("type", sectionTypes)
	} else {
		d.FieldU8("type", sectionTypes)

		d.FieldRawLen("reserved2", 5)
		d.FieldBool("attr_some_instructions")
		d.FieldBool("attr_ext_reloc")
		d.FieldBool("attr_loc_reloc")

		d.FieldRawLen("reserved1", 8)

		d.FieldBool("attr_pure_instructions")
		d.FieldBool("attr_no_toc")
		d.FieldBool("attr_strip_static_syms")
		d.FieldBool("attr_no_dead_strip")
		d.FieldBool("attr_live_support")
		d.FieldBool("attr_self_modifying_code")
		d.FieldBool("attr_debug")
		d.FieldRawLen("reserved0", 1)
	}
}

func threadStateI386Decode(d *decode.D) {
	d.FieldU32("eax")
	d.FieldU32("ebx")
	d.FieldU32("ecx")
	d.FieldU32("edx")
	d.FieldU32("edi")
	d.FieldU32("esi")
	d.FieldU32("ebp")
	d.FieldU32("esp")
	d.FieldU32("ss")
	d.FieldU32("eflags")
	d.FieldU32("eip")
	d.FieldU32("cs")
	d.FieldU32("ds")
	d.FieldU32("es")
	d.FieldU32("fs")
	d.FieldU32("gs")
}

func threadStateX8664Decode(d *decode.D) {
	d.FieldU64("rax")
	d.FieldU64("rbx")
	d.FieldU64("rcx")
	d.FieldU64("rdx")
	d.FieldU64("rdi")
	d.FieldU64("rsi")
	d.FieldU64("rbp")
	d.FieldU64("rsp")
	d.FieldU64("r8")
	d.FieldU64("r9")
	d.FieldU64("r10")
	d.FieldU64("r11")
	d.FieldU64("r12")
	d.FieldU64("r13")
	d.FieldU64("r14")
	d.FieldU64("r15")
	d.FieldU64("rip")
	d.FieldU64("rflags")
	d.FieldU64("cs")
	d.FieldU64("fs")
	d.FieldU64("gs")
}

func threadStateARM32Decode(d *decode.D) {
	rIdx := 0
	d.FieldStructArrayLoop("r", "r", func() bool {
		return rIdx < 13
	}, func(d *decode.D) {
		d.FieldU32("value")
		rIdx++
	})
	d.FieldU32("sp")
	d.FieldU32("lr")
	d.FieldU32("pc")
	d.FieldU32("cpsr")
}

func threadStateARM64Decode(d *decode.D) {
	rIdx := 0
	d.FieldStructArrayLoop("r", "r", func() bool {
		return rIdx < 29
	}, func(d *decode.D) {
		d.FieldU64("value")
		rIdx++
	})
	d.FieldU64("fp")
	d.FieldU64("lr")
	d.FieldU64("sp")
	d.FieldU64("pc")
	d.FieldU32("cpsr")
	d.FieldU32("pad")
}

func threadStatePPC32Decode(d *decode.D) {
	srrIdx := 0
	d.FieldStructArrayLoop("srr", "srr", func() bool {
		return srrIdx < 2
	}, func(d *decode.D) {
		d.FieldU32("value")
		srrIdx++
	})
	rIdx := 0
	d.FieldStructArrayLoop("r", "r", func() bool {
		return rIdx < 32
	}, func(d *decode.D) {
		d.FieldU32("value")
		rIdx++
	})
	d.FieldU32("ct")
	d.FieldU32("xer")
	d.FieldU32("lr")
	d.FieldU32("ctr")
	d.FieldU32("mq")
	d.FieldU32("vrsave")
}

func threadStatePPC64Decode(d *decode.D) {
	srrIdx := 0
	d.FieldStructArrayLoop("srr", "srr", func() bool {
		return srrIdx < 2
	}, func(d *decode.D) {
		d.FieldU64("value")
		srrIdx++
	})
	rIdx := 0
	d.FieldStructArrayLoop("r", "r", func() bool {
		return rIdx < 32
	}, func(d *decode.D) {
		d.FieldU64("value")
		rIdx++
	})
	d.FieldU32("ct")
	d.FieldU64("xer")
	d.FieldU64("lr")
	d.FieldU64("ctr")
	d.FieldU32("vrsave")
}
