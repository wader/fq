package pe

// string table:
// .coff.pointer_to_symbol_table as $off | .coff.number_of_symbols as $n | ($off+($n*18)) as $o | (tobytes[$o:$o+4] | explode | reverse |tobytes |  tonumber) as $s | tobytes[$o:$o+$s] | dd

// https://osandamalith.com/2020/07/19/exploring-the-ms-dos-stub/
// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format
// https://upload.wikimedia.org/wikipedia/commons/1/1b/Portable_Executable_32_bit_Structure_in_SVG_fixed.svg

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.COFF,
		&decode.Format{
			Description: "Common Object File Format",
			DecodeFn:    peCoffStubDecode,
			DefaultInArg: format.COFF_In{
				FilePointerOffset: 0,
			},
		})
}

const (
	peFormat32     = 0x10b
	peFormat32Plus = 0x20b
)

var peFormatNames = scalar.UintMapSymStr{
	peFormat32:     "pe32",
	peFormat32Plus: "pe32+",
}

const (
	MachineTypeUNKNOWN = 0x0
	MachineTypeALPHA   = 0x184
	MachineTypeALPHA64 = 0x284
	MachineTypeAM33    = 0x1d3
	MachineTypeAMD64   = 0x8664
	MachineTypeARM     = 0x1c0
	MachineTypeARM64   = 0xaa64
	MachineTypeARMNT   = 0x1c4
	// duplicate
	// MachineTypeAXP64       = 0x284
	MachineTypeEBC         = 0xebc
	MachineTypeI386        = 0x14c
	MachineTypeIA64        = 0x200
	MachineTypeLOONGARCH32 = 0x6232
	MachineTypeLOONGARCH64 = 0x6264
	MachineTypeM32R        = 0x9041
	MachineTypeMIPS16      = 0x266
	MachineTypeMIPSFPU     = 0x366
	MachineTypeMIPSFPU16   = 0x466
	MachineTypePOWERPC     = 0x1f0
	MachineTypePOWERPCFP   = 0x1f1
	MachineTypeR4000       = 0x166
	MachineTypeRISCV32     = 0x5032
	MachineTypeRISCV64     = 0x5064
	MachineTypeRISCV128    = 0x5128
	MachineTypeSH3         = 0x1a2
	MachineTypeSH3DSP      = 0x1a3
	MachineTypeSH4         = 0x1a6
	MachineTypeSH5         = 0x1a8
	MachineTypeTHUMB       = 0x1c2
	MachineTypeWCEMIPSV2   = 0x169
)

var MachineTypeNames = scalar.UintMapSymStr{
	MachineTypeUNKNOWN: "unknown",
	MachineTypeALPHA:   "alpha",
	MachineTypeALPHA64: "alpha64",
	MachineTypeAM33:    "am33",
	MachineTypeAMD64:   "amd64",
	MachineTypeARM:     "arm",
	MachineTypeARM64:   "arm64",
	MachineTypeARMNT:   "armnt",
	//MachineTypeAXP64:        "axp64",
	MachineTypeEBC:         "ebc",
	MachineTypeI386:        "i386",
	MachineTypeIA64:        "ia64",
	MachineTypeLOONGARCH32: "loongarch32",
	MachineTypeLOONGARCH64: "loongarch64",
	MachineTypeM32R:        "m32r",
	MachineTypeMIPS16:      "mips16",
	MachineTypeMIPSFPU:     "mipsfpu",
	MachineTypeMIPSFPU16:   "mipsfpu16",
	MachineTypePOWERPC:     "powerpc",
	MachineTypePOWERPCFP:   "powerpcfp",
	MachineTypeR4000:       "r4000",
	MachineTypeRISCV32:     "riscv32",
	MachineTypeRISCV64:     "riscv64",
	MachineTypeRISCV128:    "riscv128",
	MachineTypeSH3:         "sh3",
	MachineTypeSH3DSP:      "sh3dsp",
	MachineTypeSH4:         "sh4",
	MachineTypeSH5:         "sh5",
	MachineTypeTHUMB:       "thumb",
	MachineTypeWCEMIPSV2:   "wcemipsv2",
}

const (
	SubSystemUNKNOWN                  = 0
	SubSystemNATIVE                   = 1
	SubSystemWINDOWS_GUI              = 2
	SubSystemWINDOWS_CUI              = 3
	SubSystemOS2_CUI                  = 5
	SubSystemPOSIX_CUI                = 7
	SubSystemNATIVE_WINDOWS           = 8
	SubSystemWINDOWS_CE_GUI           = 9
	SubSystemEFI_APPLICATION          = 10
	SubSystemEFI_BOOT_SERVICE_DRIVER  = 11
	SubSystemEFI_RUNTIME_DRIVER       = 12
	SubSystemEFI_ROM                  = 13
	SubSystemXBOX                     = 14
	SubSystemWINDOWS_BOOT_APPLICATION = 16
)

var subSystemNames = scalar.UintMapSymStr{
	SubSystemUNKNOWN:                  "unknown",
	SubSystemNATIVE:                   "native",
	SubSystemWINDOWS_GUI:              "windows_gui",
	SubSystemWINDOWS_CUI:              "windows_cui",
	SubSystemOS2_CUI:                  "os2_cui",
	SubSystemPOSIX_CUI:                "posix_cui",
	SubSystemNATIVE_WINDOWS:           "native_windows",
	SubSystemWINDOWS_CE_GUI:           "windows_ce_gui",
	SubSystemEFI_APPLICATION:          "efi_application",
	SubSystemEFI_BOOT_SERVICE_DRIVER:  "efi_boot_service_driver",
	SubSystemEFI_RUNTIME_DRIVER:       "efi_runtime_driver",
	SubSystemEFI_ROM:                  "efi_rom",
	SubSystemXBOX:                     "xbox",
	SubSystemWINDOWS_BOOT_APPLICATION: "windows_boot_application",
}

const (
	symClassEndOfFunction   = 0xff // A special symbol that represents the end of function, for debugging purposes.
	symClassNull            = 0    // No assigned storage class.
	symClassAutomatic       = 1    // The automatic (stack) variable. The Value field specifies the stack frame offset.
	symClassExternal        = 2    // A value that Microsoft tools use for external symbols. The Value field indicates the size if the section number is IMAGE_SYM_UNDEFINED (0). If the section number is not zero, then the Value field specifies the offset within the section.
	symClassStatic          = 3    // The offset of the symbol within the section. If the Value field is zero, then the symbol represents a section name.
	symClassRegister        = 4    // A register variable. The Value field specifies the register number.
	symClassExternalDef     = 5    // A symbol that is defined externally.
	symClassLabel           = 6    // A code label that is defined within the module. The Value field specifies the offset of the symbol within the section.
	symClassUndefinedLabel  = 7    // A reference to a code label that is not defined.
	symClassMemberOfStruct  = 8    // The structure member. The Value field specifies the n th member.
	symClassArgument        = 9    // A formal argument (parameter) of a function. The Value field specifies the n th argument.
	symClassStructTag       = 10   // The structure tag-name entry.
	symClassMemberOfUnion   = 11   // A union member. The Value field specifies the n th member.
	symClassUnionTag        = 12   // The Union tag-name entry.
	symClassTypeDefinition  = 13   // A Typedef entry.
	symClassUndefinedStatic = 14   // A static data declaration.
	symClassEnumTag         = 15   // An enumerated type tagname entry.
	symClassMemberOfEnum    = 16   // A member of an enumeration. The Value field specifies the n th member.
	symClassRegisterParam   = 17   // A register parameter.
	symClassBitField        = 18   // A bit-field reference. The Value field specifies the n th bit in the bit field.
	symClassBlock           = 100  // A .bb (beginning of block) or .eb (end of block) record. The Value field is the relocatable address of the code location.
	symClassFunction        = 101  // A value that Microsoft tools use for symbol records that define the extent of a function: begin function (.bf ), end function ( .ef ), and lines in function ( .lf ). For .lf records, the Value field gives the number of source lines in the function. For .ef records, the Value field gives the size of the function code.
	symClassEndOfStruct     = 102  // An end-of-structure entry.
	symClassFile            = 103  // A value that Microsoft tools, as well as traditional COFF format, use for the source-file symbol record. The symbol is followed by auxiliary records that name the file.
	symClassSection         = 104  // A definition of a section (Microsoft tools use STATIC storage class instead).
	symClassWeakExternal    = 105  // A weak external. For more information, see Auxiliary Format 3: Weak Externals.
	symClassClrToken        = 107  // A CLR token symbol. The name is an ASCII string that consists of the hexadecimal value of the token. For more information, see CLR Token Definition (Object Only).
)

var symClassNames = scalar.UintMapSymStr{
	symClassEndOfFunction:   "end_of_function",
	symClassNull:            "null",
	symClassAutomatic:       "automatic",
	symClassExternal:        "external",
	symClassStatic:          "static",
	symClassRegister:        "register",
	symClassExternalDef:     "external_def",
	symClassLabel:           "label",
	symClassUndefinedLabel:  "undefined_label",
	symClassMemberOfStruct:  "member_of_struct",
	symClassArgument:        "argument",
	symClassStructTag:       "struct_tag",
	symClassMemberOfUnion:   "member_of_union",
	symClassUnionTag:        "union_tag",
	symClassTypeDefinition:  "type_definition",
	symClassUndefinedStatic: "undefined_static",
	symClassEnumTag:         "enum_tag",
	symClassMemberOfEnum:    "member_of_enum",
	symClassRegisterParam:   "register_param",
	symClassBitField:        "bit_field",
	symClassBlock:           "block",
	symClassFunction:        "function",
	symClassEndOfStruct:     "end_of_struct",
	symClassFile:            "file",
	symClassSection:         "section",
	symClassWeakExternal:    "weak_external",
	symClassClrToken:        "clr_token",
}

const (
	symTypeNull   = 0
	symTypeVoid   = 1
	symTypeChar   = 2
	symTypeShort  = 3
	symTypeInt    = 4
	symTypeLong   = 5
	symTypeFloat  = 6
	symTypeDouble = 7
	symTypeStruct = 8
	symTypeUnion  = 9
	symTypeEnum   = 10
	symTypeMoe    = 11
	symTypeByte   = 12
	symTypeWord   = 13
	symTypeUint   = 14
	symTypeDword  = 15
)

var symBaseTypeNames = scalar.UintMapSymStr{
	symTypeNull:   "null",
	symTypeVoid:   "void",
	symTypeChar:   "char",
	symTypeShort:  "short",
	symTypeInt:    "int",
	symTypeLong:   "long",
	symTypeFloat:  "float",
	symTypeDouble: "double",
	symTypeStruct: "struct",
	symTypeUnion:  "union",
	symTypeEnum:   "enum",
	symTypeMoe:    "moe",
	symTypeByte:   "byte",
	symTypeWord:   "word",
	symTypeUint:   "uint",
	symTypeDword:  "dword",
}

const (
	symDtypeNull     = 0
	symDtypePointer  = 1
	symDtypeFunction = 2
	symDtypeArray    = 3
)

var symBaseDTypeNames = scalar.UintMapSymStr{
	symDtypeNull:     "null",
	symDtypePointer:  "pointer",
	symDtypeFunction: "function",
	symDtypeArray:    "array",
}

// type stringTable []string

// func (m stringTable) MapStr(s scalar.Str) (scalar.Str, error) {
// 	if s.Actual == "" || s.Actual[0] != '/' {
// 		return s, nil
// 	}
// 	un, err := strconv.ParseUint(s.Actual[1:], 10, 64)
// 	if err != nil {
// 		// ignore error
// 		//nolint: nilerr
// 		return s, nil
// 	}
// 	n := int(un)
// 	if n >= len(m) {
// 		return s, nil
// 	}

// 	s.Sym = m[n]

// 	return s, nil
// }

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

type stringTable string

func (m stringTable) MapStr(s scalar.Str) (scalar.Str, error) {
	if s.Actual[0] == '/' {
		// /### section name

		s.Actual = strings.TrimRight(s.Actual, "\x00")

		un, err := strconv.ParseUint(s.Actual[1:], 10, 64)
		if err != nil {
			// ignore error
			//nolint: nilerr
			return s, nil
		}
		n := int(un) - 4

		s.Sym = strIndexNull(n, string(m))

		return s, nil
	} else if s.Actual[0:4] == "\x00\x00\x00\x00" {
		// \0\0\0\0LE32 symbol name
		n := binary.LittleEndian.Uint32([]byte(s.Actual)[4:8]) - 4
		s.Sym = strIndexNull(int(n), string(m))
	} else {
		// right null padded
		s.Actual = strings.TrimRight(s.Actual, "\x00")
	}

	return s, nil
}

func peCoffStubDecode(d *decode.D) any {
	var pci format.COFF_In
	d.ArgAs(&pci)

	addrSize := 64

	d.Endian = decode.LittleEndian

	d.FieldRawLen("signature", 4*8, d.AssertBitBuf([]byte("PE\x00\x00")))
	d.FieldU16("machine", MachineTypeNames, scalar.UintHex)
	numberOfSections := d.FieldU16("number_of_sections")
	d.FieldU32("time_date_stamp", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
	pointerToSymbolTable := d.FieldU32("pointer_to_symbol_table", scalar.UintHex)
	numberOfSymbols := d.FieldU32("number_of_symbols")
	sizeOfOptionalHeader := d.FieldU16("size_of_optional_header")
	d.FieldStruct("characteristics", func(d *decode.D) {
		// bytes swapped as it's a 16bit LE
		d.FieldBool("bytes_reversed_lo")   //  0x0080 // Little endian: the least significant bit (LSB) precedes the most significant bit (MSB) in memory. This flag is deprecated and should be zero.
		d.FieldBool("reserved")            //  0x0040 // This flag is reserved for future use.
		d.FieldBool("large_address_aware") //  0x0020 // Application can handle > 2-GB addresses.
		d.FieldBool("aggressive_ws_trim")  //  0x0010 // Obsolete. Aggressively trim working set. This flag is deprecated for Windows 2000 and later and must be zero.
		d.FieldBool("local_syms_stripped") //  0x0008 // COFF symbol table entries for local symbols have been removed. This flag is deprecated and should be zero.
		d.FieldBool("line_nums_stripped")  //  0x0004 // COFF line numbers have been removed. This flag is deprecated and should be zero.
		d.FieldBool("executable_image")    //  0x0002 // Image only. This indicates that the image file is valid and can be run. If this flag is not set, it indicates a linker error.
		d.FieldBool("relocs_stripped")     //  0x0001 // Image only, Windows CE, and Microsoft Windows NT and later. This indicates that the file does not contain base relocations and must therefore be loaded at its preferred base address. If the base address is not available, the loader reports an error. The default behavior of the linker is to strip base relocations from executable (EXE) files.

		d.FieldBool("bytes_reversed_hi")       //  0x8000 // Big endian: the MSB precedes the LSB in memory. This flag is deprecated and should be zero.
		d.FieldBool("up_system_only")          //  0x4000 // The file should be run only on a uniprocessor machine.
		d.FieldBool("dll")                     //  0x2000 // The image file is a dynamic-link library (DLL). Such files are considered executable files for almost all purposes, although they cannot be directly run.
		d.FieldBool("system")                  //  0x1000 // The image file is a system file, not a user program.
		d.FieldBool("net_run_from_swap")       //  0x0800 // If the image is on network media, fully load it and copy it to the swap file.
		d.FieldBool("removable_run_from_swap") //  0x0400 // If the image is on removable media, fully load it and copy it to the swap file.
		d.FieldBool("debug_stripped")          //  0x0200 // Debugging information is removed from the image file.
		d.FieldBool("32bit_machine")           //  0x0100 // Machine is based on a 32-bit-word architecture.
	})

	if pointerToSymbolTable != 0 {
		pointerToSymbolTable -= uint64(pci.FilePointerOffset)
	}
	stringTablePos := (int64(pointerToSymbolTable) + int64(numberOfSymbols)*18) * 8

	var stringTableMapper stringTable
	if stringTablePos < d.Len()+4*8 {
		d.SeekAbs(stringTablePos, func(d *decode.D) {
			stringTableSize := d.U32() - 4
			if stringTableSize*8 > uint64(d.BitsLeft()) {
				return
			}
			stringTableMapper = stringTable(d.UTF8(int(stringTableSize)))
			// d.FramedFn(int64(stringTableSize)*8, func(d *decode.D) {
			// 	for !d.End() {
			// 		stringTable = append(stringTable, d.UTF8Null())
			// 	}
			// })
		})
	}

	// how to know if image only? windows specific?
	if sizeOfOptionalHeader > 0 {
		d.FieldStruct("optional_header", func(d *decode.D) {
			d.FramedFn(int64(sizeOfOptionalHeader)*8, func(d *decode.D) {
				peFormat := d.FieldU16("format", peFormatNames, scalar.UintHex)
				d.FieldU8("major_linker_version")
				d.FieldU8("minor_linker_version")
				d.FieldU32("size_of_code")
				d.FieldU32("size_of_initialized_data")
				d.FieldU32("size_of_uninitialized_data")
				d.FieldU32("address_of_entry_point", scalar.UintHex)
				d.FieldU32("base_of_code", scalar.UintHex)
				if peFormat == peFormat32 {
					d.FieldU32("base_of_data", scalar.UintHex)
					addrSize = 32
				}

				d.FieldU("image_base", addrSize, scalar.UintHex)
				d.FieldU32("section_alignment")
				d.FieldU32("file_alignment")
				d.FieldU16("major_os_version")
				d.FieldU16("minor_os_version")
				d.FieldU16("major_image_version")
				d.FieldU16("minor_image_version")
				d.FieldU16("major_subsystem_version")
				d.FieldU16("minor_subsystem_version")
				d.FieldU32("win32_version")
				d.FieldU32("size_of_image")
				d.FieldU32("size_of_headers")
				d.FieldU32("chunk_sum", scalar.UintHex)
				d.FieldU16("subsystem", subSystemNames)
				d.FieldStruct("dll_characteristics", func(d *decode.D) {
					d.FieldBool("force_integrity") // Code Integrity checks are enforced.
					d.FieldBool("dynamic_base")    // DLL can be relocated at load time.
					d.FieldBool("high_entropy_va") // Image can handle a high entropy 64-bit virtual address space.
					d.FieldBool("reserved0")       // ??
					d.FieldBool("reserved1")
					d.FieldBool("reserved2")
					d.FieldBool("reserved3")
					d.FieldBool("reserved4")

					d.FieldBool("terminal_server_aware") // Terminal Server aware.
					d.FieldBool("guard_cf")              // Image supports Control Flow Guard.
					d.FieldBool("wdm_driver")            // A WDM driver.
					d.FieldBool("appcontainer")          // Image must execute in an AppContainer.
					d.FieldBool("no_bind")               // Do not bind the image.
					d.FieldBool("no_seh")                // Does not use structured exception (SE) handling. No SE handler may be called in this image.
					d.FieldBool("no_isolation")          // Isolation aware, but do not isolate the image.
					d.FieldBool("nx_compat")             // Image is NX compatible.
				})
				d.FieldU("size_of_track_reserve", addrSize)
				d.FieldU("size_of_stack_commit", addrSize)
				d.FieldU("size_of_heap_reserve", addrSize)
				d.FieldU("size_of_heap_commit", addrSize)
				d.FieldU32("loader_flags")
				d.FieldU32("number_of_rva_and_sizes")

				d.FieldU32("export_table_address", scalar.UintHex) //The export table address and size. For more information see .edata Section (Image Only).
				d.FieldU32("export_table_size")
				d.FieldU32("import_table_address", scalar.UintHex) //The import table address and size. For more information, see The .idata Section.
				d.FieldU32("import_table_size")
				d.FieldU32("resource_table_address", scalar.UintHex) //The resource table address and size. For more information, see The .rsrc Section.
				d.FieldU32("resource_table_size")
				d.FieldU32("exception_table_address", scalar.UintHex) //The exception table address and size. For more information, see The .pdata Section.
				d.FieldU32("exception_table_size")
				d.FieldU32("certificate_table_address", scalar.UintHex) //The attribute certificate table address and size. For more information, see The Attribute Certificate Table (Image Only).
				d.FieldU32("certificate_table_size")
				d.FieldU32("base_relocation_table_address", scalar.UintHex) //The base relocation table address and size. For more information, see The .reloc Section (Image Only).
				d.FieldU32("base_relocation_table_size")
				d.FieldU32("debug_address", scalar.UintHex) //The debug data starting address and size. For more information, see The .debug Section.
				d.FieldU32("debug_size")
				d.FieldU64("architecture")                      //Reserved, must be 0
				d.FieldU64("global_ptr", scalar.UintHex)        //The RVA of the value to be stored in the global pointer register. The size member of this structure must be set to zero.
				d.FieldU32("tls_table_address", scalar.UintHex) //The thread local storage (TLS) table address and size. For more information, see The .tls Section.
				d.FieldU32("tls_table_size")
				d.FieldU32("load_config_table_address", scalar.UintHex) //The load configuration table address and size. For more information, see The Load Configuration Structure (Image Only).
				d.FieldU32("load_config_table_size")
				d.FieldU32("bound_import_address", scalar.UintHex) //The bound import table address and size.
				d.FieldU32("bound_import_size")
				d.FieldU32("iat_address", scalar.UintHex) //The import address table address and size. For more information, see Import Address Table.
				d.FieldU32("iat_size")
				d.FieldU32("delay_import_descriptor_address", scalar.UintHex) //The delay import descriptor address and size. For more information, see Delay-Load Import Tables (Image Only).
				d.FieldU32("delay_import_descriptor_size")
				d.FieldU32("clr_runtime_header_address", scalar.UintHex) //The CLR runtime header address and size. For more information, see The .cormeta Section (Object Only).
				d.FieldU32("clr_runtime_header_size")
				d.FieldU64("reserved") //must be zero

				// TODO: where?
				/*numberOfRvaAndSizes :=*/
				/*
					d.FieldArray("data_directories", func(d *decode.D) {
						for i := 0; i < int(numberOfRvaAndSizes); i++ {
							d.FieldStruct("data_directory", func(d *decode.D) {
								d.FieldU32("virtual_address", scalar.UintHex)
								d.FieldU32("size")
							})
						}
					})
				*/

				d.FieldRawLen("unknown", d.BitsLeft())
			})
		})
	}

	// TODO: section_alignment?

	d.FieldArray("sections", func(d *decode.D) {
		for i := uint64(0); i < numberOfSections; i++ {
			d.FieldStruct("section", func(d *decode.D) {
				name := d.FieldUTF8("name", 8, stringTableMapper)                     // An 8-byte, null-padded UTF-8 encoded string. If the string is exactly 8 characters long, there is no terminating null. For longer names, this field contains a slash (/) that is followed by an ASCII representation of a decimal number that is an offset into the string table. Executable images do not use a string table and do not support section names longer than 8 characters. Long names in object files are truncated if they are emitted to an executable file.
				d.FieldU32("virtual_size")                                            // The total size of the section when loaded into memory. If this value is greater than SizeOfRawData, the section is zero-padded. This field is valid only for executable images and should be set to zero for object files.
				virtualAddress := d.FieldU32("virtual_address", scalar.UintHex)       // For executable images, the address of the first byte of the section relative to the image base when the section is loaded into memory. For object files, this field is the address of the first byte before relocation is applied; for simplicity, compilers should set this to zero. Otherwise, it is an arbitrary value that is subtracted from offsets during relocation.
				sizeOfRawData := d.FieldU32("size_of_raw_data")                       // The size of the section (for object files) or the size of the initialized data on disk (for image files). For executable images, this must be a multiple of FileAlignment from the optional header. If this is less than VirtualSize, the remainder of the section is zero-filled. Because the SizeOfRawData field is rounded but the VirtualSize field is not, it is possible for SizeOfRawData to be greater than VirtualSize as well. When a section contains only uninitialized data, this field should be zero.
				pointerToRawData := d.FieldU32("pointer_to_raw_data", scalar.UintHex) // The file pointer to the first page of the section within the COFF file. For executable images, this must be a multiple of FileAlignment from the optional header. For object files, the value should be aligned on a 4-byte boundary for best performance. When a section contains only uninitialized data, this field should be zero.
				d.FieldU32("pointer_to_relocations", scalar.UintHex)                  // The file pointer to the beginning of relocation entries for the section. This is set to zero for executable images or if there are no relocations.
				d.FieldU32("pointer_to_line_numbers", scalar.UintHex)                 // The file pointer to the beginning of line-number entries for the section. This is set to zero if there are no COFF line numbers. This value should be zero for an image because COFF debugging information is deprecated.
				d.FieldU16("number_of_relocations")                                   // The number of relocation entries for the section. This is set to zero for executable images.
				d.FieldU16("number_of_line_numbers")                                  // The number of line-number entries for the section. This value should be zero for an image because COFF debugging information is deprecated.

				d.FieldStruct("characteristics", func(d *decode.D) {

					// 32 bit LE flags

					d.FieldBool("cnt_uninitialized_data") // The section contains uninitialized data.
					d.FieldBool("cnt_initialized_data")   // The section contains initialized data.
					d.FieldBool("cnt_code")               // The section contains executable code.
					d.FieldBool("reserved")               // Reserved for future use.
					d.FieldBool("type_no_pad")            // The section should not be padded to the next boundary. This flag is obsolete and is replaced by IMAGE_SCN_ALIGN_1BYTES. This is valid only for object files.
					d.FieldBool("reserved0")              // Reserved for future use.
					d.FieldBool("reserved1")              // Reserved for future use.
					d.FieldBool("reserved2")              // Reserved for future use.

					d.FieldBool("gprel")      // The section contains data referenced through the global pointer (GP).
					d.FieldBool("unknown0")   // ??
					d.FieldBool("unknown1")   // ??
					d.FieldBool("lnk_comdat") // The section contains COMDAT data. For more information, see COMDAT Sections (Object Only). This is valid only for object files.
					d.FieldBool("lnk_remove") // The section will not become part of the image. This is valid only for object files.
					d.FieldBool("reserved3")  // Reserved for future use.
					d.FieldBool("lnk_info")   // The section contains comments or other information. The .drectve section has this type. This is valid for object files only.
					d.FieldBool("lnk_other")  // Reserved for future use.

					d.FieldBool("align_128bytes") // Align data on a 128-byte boundary. Valid only for object files.
					d.FieldBool("align_8bytes")   // Align data on an 8-byte boundary. Valid only for object files.
					d.FieldBool("align_2bytes")   // Align data on a 2-byte boundary. Valid only for object files.
					d.FieldBool("align_1bytes")   // Align data on a 1-byte boundary. Valid only for object files.
					d.FieldBool("mem_preload")    // Reserved for future use.
					d.FieldBool("mem_locked")     // Reserved for future use.
					d.FieldBool("mem_16bit")      // Reserved for future use.
					d.FieldBool("mem_purgeable")  // Reserved for future use. TODO was 0x00020000 in docnumberOfSymbols

					d.FieldBool("mem_write")       // The section can be written to.
					d.FieldBool("mem_read")        // The section can be read.
					d.FieldBool("mem_execute")     // The section can be executed as code.
					d.FieldBool("mem_shared")      // The section can be shared in memory.
					d.FieldBool("mem_not_paged")   // The section is not pageable.
					d.FieldBool("mem_not_cached")  // The section cannot be cached.
					d.FieldBool("mem_discardable") // The section can be discarded as needed.
					d.FieldBool("lnk_nreloc_ovfl") // The section contains extended relocations.

					// IMAGE_SCN_ALIGN_4BYTES 0x00300000 Align data on a 4-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_16BYTES 0x00500000 Align data on a 16-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_32BYTES 0x00600000 Align data on a 32-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_64BYTES 0x00700000 Align data on a 64-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_256BYTES 0x00900000 Align data on a 256-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_512BYTES 0x00A00000 Align data on a 512-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_1024BYTES 0x00B00000 Align data on a 1024-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_2048BYTES 0x00C00000 Align data on a 2048-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_4096BYTES 0x00D00000 Align data on a 4096-byte boundary. Valid only for object files.
					// IMAGE_SCN_ALIGN_8192BYTES 0x00E00000 Align data on an 8192-byte boundary. Valid only for object files.

					// d.FieldBool("reserved")               // Reserved for future use.
					// d.FieldBool("reserved")               // Reserved for future use.
					// d.FieldBool("reserved")               // Reserved for future use.
					// d.FieldBool("type_no_pad")            // The section should not be padded to the next boundary. This flag is obsolete and is replaced by IMAGE_SCN_ALIGN_1BYTES. This is valid only for object files.
					// d.FieldBool("reserved")               // Reserved for future use.
					// d.FieldBool("cnt_code")               // The section contains executable code.
					// d.FieldBool("cnt_initialized_data")   // The section contains initialized data.
					// d.FieldBool("cnt_uninitialized_data") // The section contains uninitialized data.

					// d.FieldBool("lnk_other")              // Reserved for future use.
					// d.FieldBool("lnk_info")               // The section contains comments or other information. The .drectve section has this type. This is valid for object files only.
					// d.FieldBool("reserved")               // Reserved for future use.
					// d.FieldBool("lnk_remove")             // The section will not become part of the image. This is valid only for object files.
					// d.FieldBool("lnk_comdat")             // The section contains COMDAT data. For more information, see COMDAT Sections (Object Only). This is valid only for object files.
					// d.FieldBool("unknown")                // The section contains data referenced through the global pointer (GP).
					// d.FieldBool("unknown")                // The section contains data referenced through the global pointer (GP).
					// d.FieldBool("gprel")                  // The section contains data referenced through the global pointer (GP).

					// d.FieldBool("mem_purgeable")          // Reserved for future use. TODO was 0x00020000 in docnumberOfSymbols
					// d.FieldBool("mem_16bit")              // Reserved for future use.
					// d.FieldBool("mem_locked")             // Reserved for future use.
					// d.FieldBool("mem_preload")            // Reserved for future use.
					// d.FieldBool("align_1bytes")           // Align data on a 1-byte boundary. Valid only for object files.
					// d.FieldBool("align_2bytes")           // Align data on a 2-byte boundary. Valid only for object files.
					// d.FieldBool("align_8bytes")           // Align data on an 8-byte boundary. Valid only for object files.
					// d.FieldBool("align_128bytes")         // Align data on a 128-byte boundary. Valid only for object files.

					// d.FieldBool("lnk_nreloc_ovfl")        // The section contains extended relocations.
					// d.FieldBool("mem_discardable")        // The section can be discarded as needed.
					// d.FieldBool("mem_not_cached")         // The section cannot be cached.
					// d.FieldBool("mem_not_paged")          // The section is not pageable.
					// d.FieldBool("mem_shared")             // The section can be shared in memory.
					// d.FieldBool("mem_execute")            // The section can be executed as code.
					// d.FieldBool("mem_read")               // The section can be read.
					// d.FieldBool("mem_write")              // The section can be written to.

				})

				_ = name
				_ = virtualAddress
				_ = sizeOfRawData
				_ = pointerToRawData

				/*
					if pointerToRawData != 0 {
						pointerToRawData -= uint64(pci.FilePointerOffset)
						d.SeekAbs(int64(pointerToRawData)*8, func(d *decode.D) {
							switch name {
							case ".idata":
								d.FieldStruct("data", func(d *decode.D) {

									count := 0
									d.FieldArray("directory_table", func(d *decode.D) {
										for {
											var n uint64
											d.FieldStruct("directory_entry", func(d *decode.D) {
												n = d.FieldU32("lookup_table_rva", scalar.UintHex)
												d.FieldU32("timestamp")
												d.FieldU32("forward_chain")
												d.FieldU32("name_rva", scalar.UintHex)
												d.FieldU32("address_table_rva", scalar.UintHex)
											})

											if n == 0 {
												break
											}

											count++
										}
									})

									_ = virtualAddress

									d.FieldArray("dlls", func(d *decode.D) {

										for i := 0; i < count; i++ {

											d.FieldStruct("dll", func(d *decode.D) {

												d.FieldArray("lookup_table1", func(d *decode.D) {
													for {
														if d.PeekUintBits(addrSize) == 0 {
															break
														}
														d.FieldU("bla", addrSize)
													}
												})
												d.FieldU("lookup_table1_end", addrSize)
											})
										}

									})

									d.FieldArray("hint_names", func(d *decode.D) {

										for i := 0; true; i++ {

											log.Printf("i: %#+v\n", i)

											d.FieldStruct("hint_name", func(d *decode.D) {

												d.FieldU16("hint")
												d.FieldUTF8Null("name")

												// if d.AlignBits(16) != 0 {
												// 	d.FieldU8("align_pad")
												// }

											})

											if i > 800 {
												break
											}
										}

									})

								})

							default:
								d.FieldRawLen("data", int64(sizeOfRawData)*8)
							}
						})
					}
				*/
			})
		}
	})

	// var stringTableMapperPos int64

	// TODO: if pointerToSymbolTable != 0?

	if pointerToSymbolTable != 0 {
		d.FieldArray("symbol_table", func(d *decode.D) {
			d.SeekAbs(int64(pointerToSymbolTable*8), func(d *decode.D) {
				for i := uint64(0); i < numberOfSymbols; i++ {
					d.FieldStruct("symbol", func(d *decode.D) {
						// TODO: name
						d.FieldUTF8("name", 8, stringTableMapper) // The name of the symbol, represented by a union of three structures. An array of 8 bytes is used if the name is not more than 8 bytes long. For more information, see Symbol Name Representation.
						d.FieldU32("value")                       // The value that is associated with the symbol. The interpretation of this field depends on SectionNumber and StorageClass. A typical meaning is the relocatable address.
						d.FieldU16("section_number")              // The signed integer that identifies the section, using a one-based index into the section table. Some values have special meaning, as defined in section 5.4.2, "Section Number Values."
						d.FieldU8("base_type", symBaseTypeNames)
						d.FieldU8("complex_type", symBaseDTypeNames)
						d.FieldU8("storage_class", symClassNames) // An enumerated value that represents storage class. For more information, see Storage Class.
						d.FieldU8("number_of_aux_symbols")        // The number of auxiliary symbol table entries that follow this record.
					})
				}
				// stringTablePos = d.Pos()
			})
		})

		d.SeekAbs(stringTablePos, func(d *decode.D) {
			// TODO: if pos != 0?
			d.FieldStruct("string_table", func(d *decode.D) {
				stringTableSize := d.FieldU32("size") - 4
				d.FramedFn(int64(stringTableSize*8), func(d *decode.D) {
					d.FieldArray("entries", func(d *decode.D) {
						for !d.End() {
							d.FieldUTF8Null("entry")
						}
					})
				})
			})
		})
	}

	return nil
}
