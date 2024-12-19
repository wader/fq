package elf

// https://refspecs.linuxbase.org/elf/gabi4+/contents.html
// https://man7.org/linux/man-pages/man5/elf.5.html
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/elf.h
// https://sourceware.org/git/?p=binutils-gdb.git;a=blob;f=include/elf/external.h;hb=HEAD

// TODO: dwarf

import (
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.ELF,
		&decode.Format{
			Description: "Executable and Linkable Format",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    elfDecode,
		})
}

const (
	LITTLE_ENDIAN = 1
	BIG_ENDIAN    = 2
)

var endianNames = scalar.UintMapSymStr{
	LITTLE_ENDIAN: "little_endian",
	BIG_ENDIAN:    "big_endian",
}

var classBits = scalar.UintMapSymUint{
	1: 32,
	2: 64,
}

const (
	CLASS_32 = 1
	CLASS_64 = 2
)

var osABINames = scalar.UintMapSymStr{
	0:   "sysv",
	1:   "hpux",
	2:   "netbsd",
	3:   "linux",
	4:   "hurd",
	5:   "86open",
	6:   "solaris",
	7:   "monterey",
	8:   "irix",
	9:   "freebsd",
	10:  "tru64",
	11:  "modesto",
	12:  "openbsd",
	97:  "arm",
	255: "standalone",
}

const (
	ET_NONE = 0
	ET_REL  = 1
	ET_EXEC = 2
	ET_DYN  = 3
	ET_CORE = 4
)

var typeNames = scalar.UintRangeToScalar{
	{Range: [2]uint64{ET_NONE, ET_NONE}, S: scalar.Uint{Sym: "none"}},
	{Range: [2]uint64{ET_REL, ET_REL}, S: scalar.Uint{Sym: "rel"}},
	{Range: [2]uint64{ET_EXEC, ET_EXEC}, S: scalar.Uint{Sym: "exec"}},
	{Range: [2]uint64{ET_DYN, ET_DYN}, S: scalar.Uint{Sym: "dyn"}},
	{Range: [2]uint64{ET_CORE, ET_CORE}, S: scalar.Uint{Sym: "core"}},
	{Range: [2]uint64{0xfe00, 0xfeff}, S: scalar.Uint{Sym: "os"}},
	{Range: [2]uint64{0xff00, 0xffff}, S: scalar.Uint{Sym: "proc"}},
}

const (
	EM_X86_64 = 0x3e
	EM_ARM64  = 0xb7
)

var machineNames = scalar.UintMap{
	0x00:      {Description: "No specific instruction set"},
	0x01:      {Sym: "we_32100", Description: "AT&T WE 32100"},
	0x02:      {Sym: "sparc", Description: "SPARC"},
	0x03:      {Sym: "x86", Description: "x86"},
	0x04:      {Sym: "m68k", Description: "Motorola 68000 (M68k)"},
	0x05:      {Sym: "m88k", Description: "Motorola 88000 (M88k)"},
	0x06:      {Sym: "intel_mcu", Description: "Intel MCU"},
	0x07:      {Sym: "intel_80860", Description: "Intel 80860"},
	0x08:      {Sym: "mips", Description: "MIPS"},
	0x09:      {Sym: "s370", Description: "IBM_System/370"},
	0x0a:      {Sym: "mips_rs3000le", Description: "MIPS RS3000 Little-endian"},
	0x0e:      {Sym: "pa_risc", Description: "Hewlett-Packard PA-RISC"},
	0x0f:      {Description: "Reserved for future use"},
	0x13:      {Sym: "80960", Description: "Intel 80960"},
	0x14:      {Sym: "powerpc", Description: "PowerPC"},
	0x15:      {Sym: "powerpc64", Description: "PowerPC (64-bit)"},
	0x16:      {Sym: "s390", Description: "S390, including S390x"},
	0x17:      {Sym: "ibm_spu_spc", Description: "IBM SPU/SPC"},
	0x24:      {Sym: "nec_v800", Description: "NEC V800"},
	0x25:      {Sym: "fr20", Description: "Fujitsu FR20"},
	0x26:      {Sym: "trw_rh_32", Description: "TRW RH-32"},
	0x27:      {Sym: "motorola_rce", Description: "Motorola RCE"},
	0x28:      {Sym: "arm", Description: "ARM (up to ARMv7/Aarch32)"},
	0x29:      {Sym: "alpha", Description: "Digital Alpha"},
	0x2a:      {Sym: "superh", Description: "SuperH"},
	0x2b:      {Sym: "sparc_v9", Description: "SPARC Version 9"},
	0x2c:      {Sym: "siemens_tricore", Description: "Siemens TriCore embedded processor"},
	0x2d:      {Sym: "argonaut_risc", Description: "Argonaut RISC Core"},
	0x2e:      {Sym: "h8_300", Description: "Hitachi H8/300"},
	0x2f:      {Sym: "h8_300h", Description: "Hitachi H8/300H"},
	0x30:      {Sym: "h8s", Description: "Hitachi H8S"},
	0x31:      {Sym: "h8/500", Description: "Hitachi H8/500"},
	0x32:      {Sym: "ia_64", Description: "IA-64"},
	0x33:      {Sym: "mips_x", Description: "Stanford MIPS-X"},
	0x34:      {Sym: "coldfire", Description: "Motorola ColdFire"},
	0x35:      {Sym: "m68hc12", Description: "Motorola M68HC12"},
	0x36:      {Sym: "fujitsu_mma", Description: "Fujitsu MMA Multimedia Accelerator"},
	0x37:      {Sym: "siemens_pcp", Description: "Siemens PCP"},
	0x38:      {Sym: "sony_ncpu_risc", Description: "Sony nCPU embedded RISC processor"},
	0x39:      {Sym: "denso_ndr1", Description: "Denso NDR1 microprocessor"},
	0x3a:      {Sym: "motorola_star", Description: "Motorola Star*Core processor"},
	0x3b:      {Sym: "toyota_me16", Description: "Toyota ME16 processor"},
	0x3c:      {Sym: "st100", Description: "STMicroelectronics ST100 processor"},
	0x3d:      {Sym: "tinyj", Description: "Advanced Logic Corp. TinyJ embedded processor family"},
	EM_X86_64: {Sym: "x86_64", Description: "AMD x86-64"},
	0x8c:      {Sym: "tms320C6000", Description: "TMS320C6000 Family"},
	EM_ARM64:  {Sym: "arm64", Description: "ARM 64-bits (ARMv8/Aarch64)"},
	0xf3:      {Sym: "risc_v", Description: "RISC-V"},
	0xf7:      {Sym: "bpf", Description: "Berkeley Packet Filter"},
	0x101:     {Sym: "wdc_65C816", Description: "WDC 65C816"},
}

const (
	PT_NULL    = 0
	PT_LOAD    = 1
	PT_DYNAMIC = 2
	PT_INTERP  = 3
	PT_NOTE    = 4
	PT_SHLIB   = 5
	PT_PHDR    = 6
	PT_TLS     = 7
)

var phTypeNames = scalar.UintRangeToScalar{
	{Range: [2]uint64{PT_NULL, PT_NULL}, S: scalar.Uint{Sym: "null", Description: "Unused element"}},
	{Range: [2]uint64{PT_LOAD, PT_LOAD}, S: scalar.Uint{Sym: "load", Description: "Loadable segment"}},
	{Range: [2]uint64{PT_DYNAMIC, PT_DYNAMIC}, S: scalar.Uint{Sym: "dynamic", Description: "Dynamic linking information"}},
	{Range: [2]uint64{PT_INTERP, PT_INTERP}, S: scalar.Uint{Sym: "interp", Description: "Interpreter to invoke"}},
	{Range: [2]uint64{PT_NOTE, PT_NOTE}, S: scalar.Uint{Sym: "note", Description: "Auxiliary information"}},
	{Range: [2]uint64{PT_SHLIB, PT_SHLIB}, S: scalar.Uint{Sym: "shlib", Description: "Reserved but has unspecified"}},
	{Range: [2]uint64{PT_PHDR, PT_PHDR}, S: scalar.Uint{Sym: "phdr", Description: "Program header location and size"}},
	{Range: [2]uint64{PT_TLS, PT_TLS}, S: scalar.Uint{Sym: "tls", Description: "Thread-Local Storage template"}},
	{Range: [2]uint64{0x6474e550, 0x6474e550}, S: scalar.Uint{Sym: "gnu_eh_frame", Description: "GNU frame unwind information"}},
	{Range: [2]uint64{0x6474e551, 0x6474e551}, S: scalar.Uint{Sym: "gnu_stack", Description: "GNU stack permission"}},
	{Range: [2]uint64{0x6474e552, 0x6474e552}, S: scalar.Uint{Sym: "gnu_relro", Description: "GNU read-only after relocation"}},
	{Range: [2]uint64{0x60000000, 0x6fffffff}, S: scalar.Uint{Sym: "os", Description: "Operating system-specific"}},
	{Range: [2]uint64{0x70000000, 0x7fffffff}, S: scalar.Uint{Sym: "proc", Description: "Processor-specific"}},
}

const (
	NT_PRSTATUS             = 1
	NT_PRFPREG              = 2
	NT_PRPSINFO             = 3
	NT_TASKSTRUCT           = 4
	NT_AUXV                 = 6
	NT_SIGINFO              = 0x53494749 // "SIGI"
	NT_FILE                 = 0x46494c45 // "FILE"
	NT_PRXFPREG             = 0x46e62b7f
	NT_PPC_VMX              = 0x100
	NT_PPC_SPE              = 0x101
	NT_PPC_VSX              = 0x102
	NT_PPC_TAR              = 0x103
	NT_PPC_PPR              = 0x104
	NT_PPC_DSCR             = 0x105
	NT_PPC_EBB              = 0x106
	NT_PPC_PMU              = 0x107
	NT_PPC_TM_CGPR          = 0x108
	NT_PPC_TM_CFPR          = 0x109
	NT_PPC_TM_CVMX          = 0x10a
	NT_PPC_TM_CVSX          = 0x10b
	NT_PPC_TM_SPR           = 0x10c
	NT_PPC_TM_CTAR          = 0x10d
	NT_PPC_TM_CPPR          = 0x10e
	NT_PPC_TM_CDSCR         = 0x10f
	NT_PPC_PKEY             = 0x110
	NT_386_TLS              = 0x200
	NT_386_IOPERM           = 0x201
	NT_X86_XSTATE           = 0x202
	NT_S390_HIGH_GPRS       = 0x300
	NT_S390_TIMER           = 0x301
	NT_S390_TODCMP          = 0x302
	NT_S390_TODPREG         = 0x303
	NT_S390_CTRS            = 0x304
	NT_S390_PREFIX          = 0x305
	NT_S390_LAST_BREAK      = 0x306
	NT_S390_SYSTEM_CALL     = 0x307
	NT_S390_TDB             = 0x308
	NT_S390_VXRS_LOW        = 0x309
	NT_S390_VXRS_HIGH       = 0x30a
	NT_S390_GS_CB           = 0x30b
	NT_S390_GS_BC           = 0x30c
	NT_S390_RI_CB           = 0x30d
	NT_S390_PV_CPU_DATA     = 0x30e
	NT_ARM_VFP              = 0x400
	NT_ARM_TLS              = 0x401
	NT_ARM_HW_BREAK         = 0x402
	NT_ARM_HW_WATCH         = 0x403
	NT_ARM_SYSTEM_CALL      = 0x404
	NT_ARM_SVE              = 0x405
	NT_ARM_PAC_MASK         = 0x406
	NT_ARM_PACA_KEYS        = 0x407
	NT_ARM_PACG_KEYS        = 0x408
	NT_ARM_TAGGED_ADDR_CTRL = 0x409
	NT_ARM_PAC_ENABLED_KEYS = 0x40a
	NT_ARM_SSVE             = 0x40b
	NT_ARM_ZA               = 0x40c
	NT_ARC_V2               = 0x600
	NT_VMCOREDD             = 0x700
	NT_MIPS_DSP             = 0x800
	NT_MIPS_FP_MODE         = 0x801
	NT_MIPS_MSA             = 0x802
	NT_LOONGARCH_CPUCFG     = 0xa00
	NT_LOONGARCH_CSR        = 0xa01
	NT_LOONGARCH_LSX        = 0xa02
	NT_LOONGARCH_LASX       = 0xa03
	NT_LOONGARCH_LBT        = 0xa04
)

var coreNoteNames = scalar.UintMap{
	NT_PRSTATUS:             {Sym: "prstatus"},
	NT_PRFPREG:              {Sym: "prfpreg"},
	NT_PRPSINFO:             {Sym: "prpsinfo"},
	NT_TASKSTRUCT:           {Sym: "taskstruct"},
	NT_AUXV:                 {Sym: "auxv"},
	NT_SIGINFO:              {Sym: "siginfo", Description: "Signal info"},
	NT_FILE:                 {Sym: "file", Description: "File info"},
	NT_PRXFPREG:             {Sym: "prxfpreg"},
	NT_PPC_SPE:              {Sym: "ppc_spe", Description: "PowerPC SPE/EVR registers"},
	NT_PPC_VSX:              {Sym: "ppc_vsx", Description: "PowerPC VSX registers"},
	NT_PPC_TAR:              {Sym: "ppc_tar", Description: "Target Address Register"},
	NT_PPC_PPR:              {Sym: "ppc_ppr", Description: "Program Priority Register"},
	NT_PPC_DSCR:             {Sym: "ppc_dscr", Description: "Data Stream Control Register"},
	NT_PPC_EBB:              {Sym: "ppc_ebb", Description: "Event Based Branch Registers"},
	NT_PPC_PMU:              {Sym: "ppc_pmu", Description: "Performance Monitor Registers"},
	NT_PPC_TM_CGPR:          {Sym: "ppc_tm_cgpr", Description: "TM checkpointed GPR Registers"},
	NT_PPC_TM_CFPR:          {Sym: "ppc_tm_cfpr", Description: "TM checkpointed FPR Registers"},
	NT_PPC_TM_CVMX:          {Sym: "ppc_tm_cvmx", Description: "TM checkpointed VMX Registers"},
	NT_PPC_TM_CVSX:          {Sym: "ppc_tm_cvsx", Description: "TM checkpointed VSX Registers"},
	NT_PPC_TM_SPR:           {Sym: "ppc_tm_spr", Description: "TM Special Purpose Registers"},
	NT_PPC_TM_CTAR:          {Sym: "ppc_tm_ctar", Description: "TM checkpointed Target Address Register"},
	NT_PPC_TM_CPPR:          {Sym: "ppc_tm_cppr", Description: "TM checkpointed Program Priority Register"},
	NT_PPC_TM_CDSCR:         {Sym: "ppc_tm_cdscr", Description: "TM checkpointed Data Stream Control Register"},
	NT_PPC_PKEY:             {Sym: "ppc_pkey", Description: "Memory Protection Keys registers"},
	NT_386_TLS:              {Sym: "386_tls", Description: "i386 TLS slots (struct user_desc)"},
	NT_386_IOPERM:           {Sym: "386_ioperm", Description: "x86 io permission bitmap (1=deny)"},
	NT_X86_XSTATE:           {Sym: "x86_xstate", Description: "x86 extended state using xsave"},
	NT_S390_HIGH_GPRS:       {Sym: "s390_high_gprs", Description: "s390 upper register halves"},
	NT_S390_TIMER:           {Sym: "s390_timer", Description: "s390 timer register"},
	NT_S390_TODCMP:          {Sym: "s390_todcmp", Description: "s390 TOD clock comparator register"},
	NT_S390_TODPREG:         {Sym: "s390_todpreg", Description: "s390 TOD programmable register"},
	NT_S390_CTRS:            {Sym: "s390_ctrs", Description: "s390 control registers"},
	NT_S390_PREFIX:          {Sym: "s390_prefix", Description: "s390 prefix register"},
	NT_S390_LAST_BREAK:      {Sym: "s390_last_break", Description: "s390 breaking event address"},
	NT_S390_SYSTEM_CALL:     {Sym: "s390_system_call", Description: "s390 system call restart data"},
	NT_S390_TDB:             {Sym: "s390_tdb", Description: "s390 transaction diagnostic block"},
	NT_S390_VXRS_LOW:        {Sym: "s390_vxrs_low", Description: "s390 vector registers 0-15 upper half"},
	NT_S390_VXRS_HIGH:       {Sym: "s390_vxrs_high", Description: "s390 vector registers 16-31"},
	NT_S390_GS_CB:           {Sym: "s390_gs_cb", Description: "s390 guarded storage registers"},
	NT_S390_GS_BC:           {Sym: "s390_gs_bc", Description: "s390 guarded storage broadcast control block"},
	NT_S390_RI_CB:           {Sym: "s390_ri_cb", Description: "s390 runtime instrumentation"},
	NT_S390_PV_CPU_DATA:     {Sym: "s390_pv_cpu_data", Description: "s390 protvirt cpu dump data"},
	NT_ARM_VFP:              {Sym: "arm_vfp", Description: "ARM VFP/NEON registers"},
	NT_ARM_TLS:              {Sym: "arm_tls", Description: "ARM TLS register"},
	NT_ARM_HW_BREAK:         {Sym: "arm_hw_break", Description: "ARM hardware breakpoint registers"},
	NT_ARM_HW_WATCH:         {Sym: "arm_hw_watch", Description: "ARM hardware watchpoint registers"},
	NT_ARM_SYSTEM_CALL:      {Sym: "arm_system_call", Description: "ARM system call number"},
	NT_ARM_SVE:              {Sym: "arm_sve", Description: "ARM Scalable Vector Extension registers"},
	NT_ARM_PAC_MASK:         {Sym: "arm_pac_mask", Description: "ARM pointer authentication code masks"},
	NT_ARM_PACA_KEYS:        {Sym: "arm_paca_keys", Description: "ARM pointer authentication address keys"},
	NT_ARM_PACG_KEYS:        {Sym: "arm_pacg_keys", Description: "ARM pointer authentication generic key"},
	NT_ARM_TAGGED_ADDR_CTRL: {Sym: "arm_tagged_addr_ctrl", Description: "arm64 tagged address control (prctl())"},
	NT_ARM_PAC_ENABLED_KEYS: {Sym: "arm_pac_enabled_keys", Description: "arm64 ptr auth enabled keys (prctl())"},
	NT_ARM_SSVE:             {Sym: "arm_ssve", Description: "ARM Streaming SVE registers"},
	NT_ARM_ZA:               {Sym: "arm_za", Description: "ARM SME ZA registers"},
	NT_ARC_V2:               {Sym: "arc_v2", Description: "ARCv2 accumulator/extra registers"},
	NT_VMCOREDD:             {Sym: "vmcoredd", Description: "Vmcore Device Dump Note"},
	NT_MIPS_DSP:             {Sym: "mips_dsp", Description: "MIPS DSP ASE registers"},
	NT_MIPS_FP_MODE:         {Sym: "mips_fp_mode", Description: "MIPS floating-point mode"},
	NT_MIPS_MSA:             {Sym: "mips_msa", Description: "MIPS SIMD registers"},
	NT_LOONGARCH_CPUCFG:     {Sym: "loongarch_cpucfg", Description: "LoongArch CPU config registers"},
	NT_LOONGARCH_CSR:        {Sym: "loongarch_csr", Description: "LoongArch control and status registers"},
	NT_LOONGARCH_LSX:        {Sym: "loongarch_lsx", Description: "LoongArch Loongson SIMD Extension registers"},
	NT_LOONGARCH_LASX:       {Sym: "loongarch_lasx", Description: "LoongArch Loongson Advanced SIMD Extension registers"},
	NT_LOONGARCH_LBT:        {Sym: "loongarch_lbt", Description: "LoongArch Loongson Binary Translation registers"},
}

const (
	SHT_NULL          = 0x0
	SHT_PROGBITS      = 0x1
	SHT_SYMTAB        = 0x2
	SHT_STRTAB        = 0x3
	SHT_RELA          = 0x4
	SHT_HASH          = 0x5
	SHT_DYNAMIC       = 0x6
	SHT_NOTE          = 0x7
	SHT_NOBITS        = 0x8
	SHT_REL           = 0x9
	SHT_SHLIB         = 0x0a
	SHT_DYNSYM        = 0x0b
	SHT_INIT_ARRAY    = 0x0e
	SHT_FINI_ARRAY    = 0x0f
	SHT_PREINIT_ARRAY = 0x10
	SHT_GROUP         = 0x11
	SHT_SYMTAB_SHNDX  = 0x12
	SHT_GNU_HASH      = 0x6ffffff6
)

var sectionHeaderTypeMap = scalar.UintMap{
	SHT_NULL:          {Sym: "null", Description: "Header inactive"},
	SHT_PROGBITS:      {Sym: "progbits", Description: "Information defined by the program"},
	SHT_SYMTAB:        {Sym: "symtab", Description: "Symbol table"},
	SHT_STRTAB:        {Sym: "strtab", Description: "String table"},
	SHT_RELA:          {Sym: "rela", Description: "Relocation entries with explicit addends"},
	SHT_HASH:          {Sym: "hash", Description: "Symbol hash table"},
	SHT_DYNAMIC:       {Sym: "dynamic", Description: "Information for dynamic linking"},
	SHT_NOTE:          {Sym: "note", Description: "Information that marks the file in some way"},
	SHT_NOBITS:        {Sym: "nobits", Description: "No space in the file"},
	SHT_REL:           {Sym: "rel", Description: "Relocation entries without explicit addends"},
	SHT_SHLIB:         {Sym: "shlib", Description: "Reserved but has unspecified semantics"},
	SHT_DYNSYM:        {Sym: "dynsym", Description: "Dynamic linking symbol table"},
	SHT_INIT_ARRAY:    {Sym: "init_array", Description: "Initialization functions"},
	SHT_FINI_ARRAY:    {Sym: "fini_array", Description: "Termination functions"},
	SHT_PREINIT_ARRAY: {Sym: "preinit_array", Description: "Pre initialization functions"},
	SHT_GROUP:         {Sym: "group", Description: "Section group"},
	SHT_SYMTAB_SHNDX:  {Sym: "symtab_shndx", Description: ""},
	SHT_GNU_HASH:      {Sym: "gnu_hash", Description: "GNU symbol hash table"},
}

const (
	STRTAB_DYNSTR   = ".dynstr"
	STRTAB_SHSTRTAB = ".shstrtab"
	STRTAB_STRTAB   = ".strtab"
)

const (
	DT_NULL            = 0
	DT_NEEDED          = 1
	DT_PLTRELSZ        = 2
	DT_PLTGOT          = 3
	DT_HASH            = 4
	DT_STRTAB          = 5
	DT_SYMTAB          = 6
	DT_RELA            = 7
	DT_RELASZ          = 8
	DT_RELAENT         = 9
	DT_STRSZ           = 10
	DT_SYMENT          = 11
	DT_INIT            = 12
	DT_FINI            = 13
	DT_SONAME          = 14
	DT_RPATH           = 15
	DT_SYMBOLIC        = 16
	DT_REL             = 17
	DT_RELSZ           = 18
	DT_RELENT          = 19
	DT_PLTREL          = 20
	DT_DEBUG           = 21
	DT_TEXTREL         = 22
	DT_JMPREL          = 23
	DT_BIND_NOW        = 24
	DT_INIT_ARRAY      = 25
	DT_FINI_ARRAY      = 26
	DT_INIT_ARRAYSZ    = 27
	DT_FINI_ARRAYSZ    = 28
	DT_RUNPATH         = 29
	DT_FLAGS           = 30 // TODO: flag map
	DT_ENCODING        = 32 // or DT_PREINIT_ARRAY ?
	DT_PREINIT_ARRAYSZ = 33
	DT_LOOS            = 0x6000000D
	DT_HIOS            = 0x6ffff000
	DT_LOPROC          = 0x70000000
	DT_HIPROC          = 0x7fffffff
)

const (
	dUnIgnored = iota
	dUnVal
	dUnPtr
	dUnUnspecified
)

type dtEntry struct {
	r   [2]uint64
	dUn int
	s   scalar.Uint
}

type dynamicTableEntries []dtEntry

func (d dynamicTableEntries) lookup(u uint64) (dtEntry, bool) {
	for _, de := range d {
		if de.r[0] >= u && de.r[1] <= u {
			return de, true
		}
	}
	return dtEntry{}, false
}

func (d dynamicTableEntries) MapUint(s scalar.Uint) (scalar.Uint, error) {
	u := s.Actual
	if de, ok := d.lookup(u); ok {
		s = de.s
		s.Actual = u
	}
	return s, nil
}

var dynamicTableMap = dynamicTableEntries{
	{r: [2]uint64{DT_NULL, DT_NULL}, dUn: dUnIgnored, s: scalar.Uint{Sym: "null", Description: "Marks end of dynamic section"}},
	{r: [2]uint64{DT_NEEDED, DT_NEEDED}, dUn: dUnVal, s: scalar.Uint{Sym: "needed", Description: "String table offset to name of a needed library"}},
	{r: [2]uint64{DT_PLTRELSZ, DT_PLTRELSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "pltrelsz", Description: "Size in bytes of PLT relocation entries"}},
	{r: [2]uint64{DT_PLTGOT, DT_PLTGOT}, dUn: dUnPtr, s: scalar.Uint{Sym: "pltgot", Description: "Address of PLT and/or GOT"}},
	{r: [2]uint64{DT_HASH, DT_HASH}, dUn: dUnPtr, s: scalar.Uint{Sym: "hash", Description: "Address of symbol hash table"}},
	{r: [2]uint64{DT_STRTAB, DT_STRTAB}, dUn: dUnPtr, s: scalar.Uint{Sym: "strtab", Description: "Address of string table"}},
	{r: [2]uint64{DT_SYMTAB, DT_SYMTAB}, dUn: dUnPtr, s: scalar.Uint{Sym: "symtab", Description: "Address of symbol table"}},
	{r: [2]uint64{DT_RELA, DT_RELA}, dUn: dUnPtr, s: scalar.Uint{Sym: "rela", Description: "Address of Rela relocation table"}},
	{r: [2]uint64{DT_RELASZ, DT_RELASZ}, dUn: dUnVal, s: scalar.Uint{Sym: "relasz", Description: "Size in bytes of the Rela relocation table"}},
	{r: [2]uint64{DT_RELAENT, DT_RELAENT}, dUn: dUnVal, s: scalar.Uint{Sym: "relaent", Description: "Size in bytes of a Rela relocation table entry"}},
	{r: [2]uint64{DT_STRSZ, DT_STRSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "strsz", Description: "Size in bytes of string table"}},
	{r: [2]uint64{DT_SYMENT, DT_SYMENT}, dUn: dUnVal, s: scalar.Uint{Sym: "syment", Description: "Size in bytes of a symbol table entry"}},
	{r: [2]uint64{DT_INIT, DT_INIT}, dUn: dUnPtr, s: scalar.Uint{Sym: "init", Description: "Address of the initialization function"}},
	{r: [2]uint64{DT_FINI, DT_FINI}, dUn: dUnPtr, s: scalar.Uint{Sym: "fini", Description: "Address of the termination function"}},
	{r: [2]uint64{DT_SONAME, DT_SONAME}, dUn: dUnVal, s: scalar.Uint{Sym: "soname", Description: "String table offset to name of shared object"}},
	{r: [2]uint64{DT_RPATH, DT_RPATH}, dUn: dUnVal, s: scalar.Uint{Sym: "rpath", Description: "String table offset to library search path (deprecated)"}},
	{r: [2]uint64{DT_SYMBOLIC, DT_SYMBOLIC}, dUn: dUnIgnored, s: scalar.Uint{Sym: "symbolic", Description: "Alert linker to search this shared object before the executable for symbols DT_REL Address of Rel relocation table"}},
	{r: [2]uint64{DT_REL, DT_REL}, dUn: dUnPtr, s: scalar.Uint{Sym: "rel", Description: ""}},
	{r: [2]uint64{DT_RELSZ, DT_RELSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "relsz", Description: "Size in bytes of Rel relocation table"}},
	{r: [2]uint64{DT_RELENT, DT_RELENT}, dUn: dUnVal, s: scalar.Uint{Sym: "relent", Description: "Size in bytes of a Rel table entry"}},
	{r: [2]uint64{DT_PLTREL, DT_PLTREL}, dUn: dUnVal, s: scalar.Uint{Sym: "pltrel", Description: "Type of relocation entry to which the PLT refers (Rela or Rel)"}},
	{r: [2]uint64{DT_DEBUG, DT_DEBUG}, dUn: dUnPtr, s: scalar.Uint{Sym: "debug", Description: "Undefined use for debugging"}},
	{r: [2]uint64{DT_TEXTREL, DT_TEXTREL}, dUn: dUnIgnored, s: scalar.Uint{Sym: "textrel", Description: "Absence of this entry indicates that no relocation entries should apply to a nonwritable segment"}},
	{r: [2]uint64{DT_JMPREL, DT_JMPREL}, dUn: dUnPtr, s: scalar.Uint{Sym: "jmprel", Description: "Address of relocation entries associated solely with the PLT"}},
	{r: [2]uint64{DT_BIND_NOW, DT_BIND_NOW}, dUn: dUnIgnored, s: scalar.Uint{Sym: "bind_now", Description: "Instruct dynamic linker to process all relocations before transferring control to the executable"}},
	{r: [2]uint64{DT_INIT_ARRAY, DT_INIT_ARRAY}, dUn: dUnPtr, s: scalar.Uint{Sym: "init_array", Description: "Address of the array of pointers to initialization functions"}},
	{r: [2]uint64{DT_FINI_ARRAY, DT_FINI_ARRAY}, dUn: dUnPtr, s: scalar.Uint{Sym: "fini_array", Description: "Address of the array of pointers to termination functions"}},
	{r: [2]uint64{DT_INIT_ARRAYSZ, DT_INIT_ARRAYSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "init_arraysz", Description: "Size in bytes of the array of initialization functions"}},
	{r: [2]uint64{DT_FINI_ARRAYSZ, DT_FINI_ARRAYSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "fini_arraysz", Description: "Size in bytes of the array of termination functions "}},
	{r: [2]uint64{DT_RUNPATH, DT_RUNPATH}, dUn: dUnVal, s: scalar.Uint{Sym: "runpath", Description: "String table offset to library search path"}},
	{r: [2]uint64{DT_FLAGS, DT_FLAGS}, dUn: dUnVal, s: scalar.Uint{Sym: "flags", Description: "Flag values specific to the object being loaded"}}, // TODO: flag ma}},
	{r: [2]uint64{DT_ENCODING, DT_ENCODING}, dUn: dUnUnspecified, s: scalar.Uint{Sym: "encoding", Description: ""}},                               // or DT_PREINIT_ARRAY }},
	{r: [2]uint64{DT_PREINIT_ARRAYSZ, DT_PREINIT_ARRAYSZ}, dUn: dUnVal, s: scalar.Uint{Sym: "preinit_arraysz", Description: "Address of the array of pointers to pre-initialization functions"}},
	{r: [2]uint64{DT_LOOS, DT_HIOS}, dUn: dUnUnspecified, s: scalar.Uint{Sym: "lo", Description: "Operating system-specific semantics"}},
	{r: [2]uint64{DT_LOPROC, DT_HIPROC}, dUn: dUnUnspecified, s: scalar.Uint{Sym: "proc", Description: "Processor-specific semantics"}},
}

var symbolTableBindingMap = scalar.UintMapSymStr{
	0:  "local",
	1:  "global",
	2:  "weak",
	10: "loos",
	12: "hios",
	13: "proc",
	14: "proc",
	15: "proc",
}

var symbolTableTypeMap = scalar.UintMapSymStr{
	0:  "notype",
	1:  "object",
	2:  "func",
	3:  "section",
	4:  "file",
	5:  "common",
	6:  "tls",
	10: "loos",
	12: "hios",
	13: "proc",
	14: "proc",
	15: "proc",
}

var symbolTableVisibilityMap = scalar.UintMapSymStr{
	0: "default",
	1: "internal",
	2: "hidden",
	3: "protected",
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

type strTable string

func (m strTable) MapUint(s scalar.Uint) (scalar.Uint, error) {
	s.Sym = strIndexNull(int(s.Actual), string(m))
	return s, nil
}

func elfDecodeSymbolHashTable(d *decode.D) {
	nBucket := d.FieldU32("nbucket")
	nChain := d.FieldU32("nchain")

	repeatFn := func(r int, fn func(d *decode.D)) func(d *decode.D) {
		return func(d *decode.D) {
			for i := 0; i < r; i++ {
				fn(d)
			}
		}
	}

	d.FieldArray("buckets", repeatFn(int(nBucket), func(d *decode.D) { d.FieldU32("bucket") }))
	d.FieldArray("chains", repeatFn(int(nChain), func(d *decode.D) { d.FieldU32("chain") }))
}

func elfDecodeSymbolTable(d *decode.D, ec elfContext, nEntries int, strTab string) {
	for i := 0; i < nEntries; i++ {
		d.FieldStruct("symbol", func(d *decode.D) {
			switch ec.archBits {
			case 32:
				d.FieldU32("name", strTable(strTab))
				d.FieldU32("value")
				d.FieldU32("size")
				d.FieldU4("bind", symbolTableBindingMap)
				d.FieldU4("type", symbolTableTypeMap)
				d.FieldU6("other_unused")
				d.FieldU2("visibility", symbolTableVisibilityMap)
				d.FieldU16("shndx")
			case 64:
				d.FieldU32("name", strTable(strTab))
				d.FieldU4("bind", symbolTableBindingMap)
				d.FieldU4("type", symbolTableTypeMap)
				d.FieldU6("other_unused")
				d.FieldU2("visibility", symbolTableVisibilityMap)
				d.FieldU16("shndx")
				d.FieldU64("value")
				d.FieldU64("size")
			}
		})
	}
}

func elfDecodeGNUHash(d *decode.D, ec elfContext, size int64, strTab string) {
	d.FramedFn(size, func(d *decode.D) {
		nBuckets := d.FieldU32("nbuckets")
		d.FieldU32("symndx")
		maskwords := d.FieldU32("maskwords")
		d.FieldU32("shift2")

		repeatFn := func(r int, fn func(d *decode.D)) func(d *decode.D) {
			return func(d *decode.D) {
				for i := 0; i < r; i++ {
					fn(d)
				}
			}
		}
		// TODO: possible to map to symbols?
		_ = strTab
		d.FieldArray("bloom_filter", repeatFn(int(maskwords), func(d *decode.D) { d.FieldU("maskword", ec.archBits) }))
		d.FieldArray("buckets", repeatFn(int(nBuckets), func(d *decode.D) { d.FieldU32("bucket") }))
		d.FieldArray("values", func(d *decode.D) {
			for !d.End() {
				d.FieldU32("value")
			}
		})
	})
}

type dynamicContext struct {
	entries   int
	strTabPtr int64
	strSzVal  int64
	strTab    string
	symEnt    int64
}

func elfReadDynamicTags(d *decode.D, ec *elfContext) dynamicContext {
	var strTabPtr int64
	var strSzVal int64
	var symEnt int64
	var entries int

	seenNull := false
	for !seenNull {
		entries++
		tag := d.U(ec.archBits)
		valPtr := d.U(ec.archBits)

		switch tag {
		case DT_STRTAB:
			strTabPtr = int64(valPtr) * 8
		case DT_STRSZ:
			strSzVal = int64(valPtr) * 8
		case DT_SYMENT:
			symEnt = int64(valPtr) * 8
		case DT_NULL:
			seenNull = true
		}
	}

	return dynamicContext{
		entries:   entries,
		strTabPtr: strTabPtr,
		strSzVal:  strSzVal,
		symEnt:    symEnt,
	}
}

type symbol struct {
	name  uint64
	value uint64
}

func elfReadSymbolTable(d *decode.D, ec *elfContext, sh sectionHeader) []symbol {
	var ss []symbol

	for i := 0; i < int(sh.size/sh.entSize); i++ {
		var name uint64
		var value uint64
		switch ec.archBits {
		case 32:
			name = d.U32()  // name
			value = d.U32() // value
			d.U32()         // size
			d.U4()          // bind
			d.U4()          // type
			d.U6()          // other_unused
			d.U2()          // visibility
			d.U16()         // shndx
		case 64:
			name = d.U32()  // name
			d.U4()          // bind
			d.U4()          // type
			d.U6()          // other_unused
			d.U2()          // visibility
			d.U16()         // shndx
			value = d.U64() // value
			d.U64()         // size
		}
		ss = append(ss, symbol{name: name, value: value})
	}

	return ss
}

type sectionHeader struct {
	addr    int64
	offset  int64
	size    int64
	entSize int64
	name    int
	typ     int
	dc      dynamicContext // if SHT_DYNAMIC
	symbols []symbol
}

const maxStrTabSize = 100_000_000

func readStrTab(d *decode.D, firstBit int64, nBytes int64) string {
	if nBytes > maxStrTabSize {
		d.Errorf("string table too large %d > %d", nBytes, maxStrTabSize)
	}
	return string(d.BytesRange(firstBit, int(nBytes)))
}

func elfReadSectionHeaders(d *decode.D, ec *elfContext) {
	for i := 0; i < ec.shNum; i++ {
		d.SeekAbs(ec.shOff + int64(i)*ec.shEntSize)
		var sh sectionHeader

		switch ec.archBits {
		case 32:
			sh.name = int(d.U32())
			sh.typ = int(d.U32())
			d.U32()                      // flags
			sh.addr = int64(d.U32() * 8) // addr
			sh.offset = int64(d.U32()) * 8
			sh.size = int64(d.U32()) * 8
			d.U32() // link
			d.U32() // info
			d.U32() // addralign
			sh.entSize = int64(d.U32()) * 8
		case 64:
			sh.name = int(d.U32())
			sh.typ = int(d.U32())
			d.U64()                      // addr
			sh.addr = int64(d.U64() * 8) // flags
			sh.offset = int64(d.U64()) * 8
			sh.size = int64(d.U64()) * 8
			d.U32() // link
			d.U32() // info
			d.U64() // addralign
			sh.entSize = int64(d.U64()) * 8
		default:
			panic("unreachable")
		}

		switch sh.typ {
		case SHT_DYNAMIC:
			d.SeekAbs(sh.offset)
			sh.dc = elfReadDynamicTags(d, ec)
		case SHT_SYMTAB:
			d.SeekAbs(sh.offset)
			sh.symbols = elfReadSymbolTable(d, ec, sh)
		}

		ec.sections = append(ec.sections, sh)
	}

	// for dynamic linking sections find offset to string table by looking up
	// section by address using string stable address
	for i := range ec.sections {
		sh := &ec.sections[i]
		if sh.typ != SHT_DYNAMIC {
			continue
		}
		if i, ok := ec.sectionIndexByAddr(sh.dc.strTabPtr); ok {
			strTabSh := ec.sections[i]
			sh.dc.strTab = readStrTab(d, strTabSh.offset, sh.dc.strSzVal/8)
		}
	}

	// provide default empty string tables to be more robust
	ec.strTabMap = map[string]string{
		STRTAB_DYNSTR:   "",
		STRTAB_SHSTRTAB: "",
		STRTAB_STRTAB:   "",
	}
	var shStrTab string
	if ec.shStrNdx < len(ec.sections) {
		shStr := ec.sections[ec.shStrNdx]
		shStrTab = readStrTab(d, shStr.offset, shStr.size/8)

		for _, sh := range ec.sections {
			if sh.typ != SHT_STRTAB {
				continue
			}
			ec.strTabMap[strIndexNull(sh.name, shStrTab)] = readStrTab(d, sh.offset, sh.size/8)
		}
	}
}

type elfContext struct {
	archBits int
	typ      int
	machine  int
	endian   decode.Endian

	phOff  int64
	phNum  int
	phSize int64

	shOff     int64
	shNum     int
	shEntSize int64

	shStrNdx int

	sections  []sectionHeader
	strTabMap map[string]string
}

func (ec *elfContext) sectionIndexByAddr(addr int64) (int, bool) {
	for i, s := range ec.sections {
		if s.addr == addr {
			return i, true
		}
	}
	return 0, false
}

func elfDecodeHeader(d *decode.D, ec *elfContext) {
	var class uint64
	var archBits int
	var endian uint64

	d.FieldStruct("ident", func(d *decode.D) {
		d.FieldRawLen("magic", 4*8, d.AssertBitBuf([]byte("\x7fELF")))
		class = d.FieldU8("class", classBits)
		endian = d.FieldU8("data", endianNames)
		d.FieldU8("version")
		d.FieldU8("os_abi", osABINames)
		d.FieldU8("abi_version")
		d.FieldRawLen("pad", 7*8, d.BitBufIsZero())
	})

	switch class {
	case CLASS_32:
		archBits = 32
	case CLASS_64:
		archBits = 64
	default:
		d.Fatalf("unknown class %d", class)
	}

	switch endian {
	case LITTLE_ENDIAN:
		d.Endian = decode.LittleEndian
	case BIG_ENDIAN:
		d.Endian = decode.BigEndian
	default:
		d.Fatalf("unknown endian %d", endian)
	}

	typ := d.FieldU16("type", typeNames, scalar.UintHex)
	machine := d.FieldU16("machine", machineNames, scalar.UintHex)
	d.FieldU32("version")
	d.FieldU("entry", archBits)
	phOff := d.FieldU("phoff", archBits)
	shOff := d.FieldU("shoff", archBits)
	d.FieldU32("flags")
	d.FieldU16("ehsize")
	phSize := d.FieldU16("phentsize")
	phNum := d.FieldU16("phnum")
	shEntSize := d.FieldU16("shentsize")
	shNum := d.FieldU16("shnum")
	shStrNdx := d.FieldU16("shstrndx")

	ec.archBits = archBits
	ec.endian = d.Endian
	ec.typ = int(typ)
	ec.machine = int(machine)
	ec.phOff = int64(phOff) * 8
	ec.phNum = int(phNum)
	ec.phSize = int64(phSize) * 8
	ec.shOff = int64(shOff) * 8
	ec.shNum = int(shNum)
	ec.shEntSize = int64(shEntSize) * 8
	ec.shStrNdx = int(shStrNdx)
}

func elfDecodeProgramHeader(d *decode.D, ec elfContext) {
	pFlags := func(d *decode.D) {
		d.FieldStruct("flags", func(d *decode.D) {
			if d.Endian == decode.LittleEndian {
				d.FieldU5("unused0")
				d.FieldBool("r")
				d.FieldBool("w")
				d.FieldBool("x")
				d.FieldU24("unused1")
			} else {
				d.FieldU29("unused0")
				d.FieldBool("r")
				d.FieldBool("w")
				d.FieldBool("x")
			}
		})
	}

	var typ uint64
	var offset uint64
	var size uint64

	switch ec.archBits {
	case 32:
		typ = d.FieldU32("type", phTypeNames)
		offset = d.FieldU("offset", ec.archBits, scalar.UintHex)
		d.FieldU("vaddr", ec.archBits, scalar.UintHex)
		d.FieldU("paddr", ec.archBits, scalar.UintHex)
		size = d.FieldU32("filesz")
		d.FieldU32("memsz")
		pFlags(d)
		d.FieldU32("align")
	case 64:
		typ = d.FieldU32("type", phTypeNames)
		pFlags(d)
		offset = d.FieldU("offset", ec.archBits, scalar.UintHex)
		d.FieldU("vaddr", ec.archBits, scalar.UintHex)
		d.FieldU("paddr", ec.archBits, scalar.UintHex)
		size = d.FieldU64("filesz")
		d.FieldU64("memsz")
		d.FieldU64("align")
	}

	d.RangeFn(int64(offset*8), int64(size*8), func(d *decode.D) {
		switch {
		case typ == PT_NOTE:
			d.FieldArray("notes", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("note", func(d *decode.D) {
						// elf manpage says this is 32 or 64 bit but it seems it is always 32
						// and that is also what readelf external.h says
						nameSz := d.FieldU32("n_namesz")
						descSz := d.FieldU32("n_descsz")
						if ec.typ == ET_CORE {
							d.FieldU32("n_type", coreNoteNames, scalar.UintHex)
						} else {
							d.FieldU32("n_type", scalar.UintHex)
						}
						d.FieldUTF8NullFixedLen("name", int(nameSz))
						nameAlign := d.AlignBits(4 * 8)
						if nameAlign != 0 {
							d.FieldRawLen("name_align", int64(nameAlign))
						}
						d.FieldRawLen("desc", int64(descSz)*8)
						descAlign := d.AlignBits(4 * 8)
						if descAlign != 0 {
							d.FieldRawLen("decs_align", int64(descAlign))
						}
					})
				}
			})
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	})
}

func elfDecodeProgramHeaders(d *decode.D, ec elfContext) {
	for i := 0; i < ec.phNum; i++ {
		d.FieldStruct("program_header", func(d *decode.D) {
			d.SeekAbs(ec.phOff + int64(i)*ec.phSize)
			elfDecodeProgramHeader(d, ec)
		})
	}
}

func elfDecodeDynamicTag(d *decode.D, ec elfContext, dc dynamicContext) {
	dtTag := d.FieldU("tag", ec.archBits, dynamicTableMap)
	name := "unspecified"
	dfMapper := scalar.UintHex
	if de, ok := dynamicTableMap.lookup(dtTag); ok {
		switch de.dUn {
		case dUnIgnored:
			name = "ignored"
		case dUnVal:
			name = "val"
			dfMapper = scalar.UintDec
		case dUnPtr:
			name = "ptr"
		}
	}

	switch dtTag {
	case DT_NEEDED:
		d.FieldU(name, ec.archBits, dfMapper, strTable(dc.strTab))
	case DT_HASH:
		v := d.FieldU(name, ec.archBits, dfMapper)
		if i, ok := ec.sectionIndexByAddr(int64(v) * 8); ok {
			d.FieldValueUint("section_index", uint64(i))
		}
	case DT_SYMTAB,
		DT_STRTAB,
		DT_PLTGOT,
		DT_JMPREL,
		DT_INIT,
		DT_FINI:
		v := d.FieldU(name, ec.archBits, dfMapper)
		if i, ok := ec.sectionIndexByAddr(int64(v) * 8); ok {
			d.FieldValueUint("section_index", uint64(i))
		}
	default:
		d.FieldU(name, ec.archBits, dfMapper)
	}
}

func elfDecodeDynamicTags(d *decode.D, ec elfContext, dc dynamicContext) {
	for i := 0; i < dc.entries; i++ {
		d.FieldStruct("dynamic_tags", func(d *decode.D) {
			elfDecodeDynamicTag(d, ec, dc)
		})
	}
}

func elfDecodeSectionHeader(d *decode.D, ec elfContext, sh sectionHeader) {
	shFlags := func(d *decode.D, archBits int) {
		d.FieldStruct("flags", func(d *decode.D) {
			if d.Endian == decode.LittleEndian {
				d.FieldBool("link_order")
				d.FieldBool("info_link")
				d.FieldBool("strings")
				d.FieldBool("merge")
				d.FieldU1("unused0")
				d.FieldBool("execinstr")
				d.FieldBool("alloc")
				d.FieldBool("write")
				d.FieldBool("tls")
				d.FieldBool("group")
				d.FieldBool("os_nonconforming")

				d.FieldU9("unused1")

				d.FieldU8("os_specific")
				d.FieldU4("processor_specific")
				if archBits == 64 {
					d.FieldU32("unused2")
				}
			} else {
				// TODO: add.FieldUnused that is per decoder?
				if archBits == 64 {
					d.FieldU32("unused0")
				}
				d.FieldU4("processor_specific")
				d.FieldU8("os_specific")
				d.FieldU9("unused1")
				d.FieldBool("tls")
				d.FieldBool("group")
				d.FieldBool("os_nonconforming")
				d.FieldBool("link_order")
				d.FieldBool("info_link")
				d.FieldBool("strings")
				d.FieldBool("merge")
				d.FieldU1("unused2")
				d.FieldBool("execinstr")
				d.FieldBool("alloc")
				d.FieldBool("write")
			}
		})
	}

	var offset int64
	var size int64
	var entSize int64
	var typ uint64

	switch ec.archBits {
	case 32:
		d.FieldU32("name", strTable(ec.strTabMap[STRTAB_SHSTRTAB]))
		typ = d.FieldU32("type", sectionHeaderTypeMap, scalar.UintHex)
		shFlags(d, ec.archBits)
		d.FieldU("addr", ec.archBits, scalar.UintHex)
		offset = int64(d.FieldU("offset", ec.archBits)) * 8
		size = int64(d.FieldU32("size", scalar.UintHex) * 8)
		d.FieldU32("link")
		d.FieldU32("info")
		d.FieldU32("addralign")
		entSize = int64(d.FieldU32("entsize") * 8)
	case 64:
		d.FieldU32("name", strTable(ec.strTabMap[STRTAB_SHSTRTAB]))
		typ = d.FieldU32("type", sectionHeaderTypeMap, scalar.UintHex)
		shFlags(d, ec.archBits)
		d.FieldU("addr", ec.archBits, scalar.UintHex)
		offset = int64(d.FieldU("offset", ec.archBits, scalar.UintHex) * 8)
		size = int64(d.FieldU64("size") * 8)
		d.FieldU32("link")
		d.FieldU32("info")
		d.FieldU64("addralign")
		entSize = int64(d.FieldU64("entsize") * 8)
	}

	// SHT_NOBITS:
	// A section of this type occupies no space in the file but otherwise resembles
	// SHT_PROGBITS. Although this section contains no bytes, the sh_offset member
	// contains the conceptual file offset.
	// SHT_NULL:
	// This value marks the section header as inactive; it does not have an
	// associated section. Other members of the section header have undefined
	// values.
	if typ == SHT_NOBITS || typ == SHT_NULL {
		// section occupies no space in file
		return
	}

	d.SeekAbs(offset)
	switch typ {
	case SHT_STRTAB:
		d.FieldUTF8("string", int(size/8))
	case SHT_DYNAMIC:
		d.FieldArray("dynamic_tags", func(d *decode.D) {
			elfDecodeDynamicTags(d, ec, sh.dc)
		})
	case SHT_HASH:
		d.FieldStruct("symbol_hash_table", elfDecodeSymbolHashTable)
	case SHT_SYMTAB:
		d.FieldArray("symbol_table", func(d *decode.D) {
			elfDecodeSymbolTable(d, ec, int(size/entSize), ec.strTabMap[STRTAB_STRTAB])
		})
	case SHT_DYNSYM:
		d.FieldArray("symbol_table", func(d *decode.D) {
			elfDecodeSymbolTable(d, ec, int(size/entSize), ec.strTabMap[STRTAB_DYNSTR])
		})
	case SHT_PROGBITS:
		// TODO: name progbits?
		// TODO: decode opcodes
		d.FieldRawLen("data", size)
	case SHT_GNU_HASH:
		d.FieldStruct("gnu_hash", func(d *decode.D) {
			elfDecodeGNUHash(d, ec, size, ec.strTabMap[STRTAB_DYNSTR])
		})
	default:
		d.FieldRawLen("data", size)
	}
}

func elfDecodeSectionHeaders(d *decode.D, ec elfContext) {
	for i := 0; i < ec.shNum; i++ {
		d.SeekAbs(ec.shOff + int64(i)*ec.shEntSize)
		d.FieldStruct("section_header", func(d *decode.D) {
			elfDecodeSectionHeader(d, ec, ec.sections[i])
		})
	}
}

func elfDecode(d *decode.D) any {
	var ec elfContext

	d.FieldStruct("header", func(d *decode.D) { elfDecodeHeader(d, &ec) })
	d.Endian = ec.endian
	// a first pass to find all sections and string table information etc
	elfReadSectionHeaders(d, &ec)
	d.FieldArray("program_headers", func(d *decode.D) {
		elfDecodeProgramHeaders(d, ec)
	})
	d.FieldArray("section_headers", func(d *decode.D) {
		elfDecodeSectionHeaders(d, ec)
	})

	return nil
}
