package nes

import (
	"embed"
	"fmt"
	"math"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed nes.jq
//go:embed nes.md
var nesFS embed.FS

func init() {
	interp.RegisterFormat(
		format.NES,
		&decode.Format{
			Description: "iNES/NES 2.0 cartridge ROM format",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeNES,
		})
	interp.RegisterFS(nesFS)
}

type nesContext struct {
	nes20       bool
	trainerSize uint64
	prgROMSize  uint64
	chrROMSize  uint64
	miscROMs    bool
}

func iNESCHRRAMSize(mapper uint64, chrROMSize uint64) uint64 {
	switch mapper {
	case 74, 191, 194:
		return 2 * 1024
	case 192, 195:
		return 4 * 1024
	case 119, 176:
		return 8 * 1024
	}
	if chrROMSize == 0 {
		return 8 * 1024
	} else {
		return 0
	}
}

func romSize(lowerByte uint64, upperNibble uint64, kbs uint64) uint64 {
	if upperNibble < 0xf {
		return (lowerByte + (upperNibble << 8)) * kbs * 1024
	}

	mm := (lowerByte << 6) >> 6
	exp := lowerByte >> 2
	return uint64(math.Pow(2, float64(exp))) * ((mm * 2) + 1)
}

func getROMSizeMapper(kbs uint64) scalar.UintFn {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		s.Sym = romSize(s.Actual, 0, kbs)
		return s, nil
	})
}

func getFlagMapper(flag uint64, descMappers ...scalar.UintMap) scalar.UintFn {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		s.Sym = s.Actual & flag

		var ds scalar.Uint = s
		for _, dm := range descMappers {
			ds = dm[ds.SymUint()]
		}
		s.Description = ds.Description

		return s, nil
	})
}

var multiplyRAMSizeMapper = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	if s.Actual == 0 {
		s.Sym = 8 * 1024
	} else {
		s.Sym = s.Actual * 8 * 1024
	}

	return s, nil
})

var shiftRAMSizeMapper = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	if s.Actual > 0 {
		s.Sym = 64 << s.Actual
	}
	return s, nil
})

var consoleTypeMapper = scalar.UintMap{
	0: {Sym: "nes", Description: "Nintendo Entertainment System/Family Computer"},
	1: {Sym: "vs", Description: "Nintendo Vs. System"},
	2: {Sym: "pc10", Description: "Nintendo Playchoice 10"},
	3: {Sym: "ext", Description: "Extended Console Type"},
}

var timingModeMapper = scalar.UintMap{
	0: {Sym: "ntsc", Description: "RP2C02 (NTSC NES)"},
	1: {Sym: "pal", Description: "RP2C07 (Licensed PAL NES)"},
	2: {Sym: "multi", Description: "Multiple-region"},
	3: {Sym: "dendy", Description: "UA6538 (Dendy)"},
}

var vsPPUMapper = scalar.UintMap{
	0:  {Description: "RP2C03B"},
	1:  {Description: "RP2C03G"},
	2:  {Description: "RP2C04-0001"},
	3:  {Description: "RP2C04-0002"},
	4:  {Description: "RP2C04-0003"},
	5:  {Description: "RP2C04-0004"},
	6:  {Description: "RC2C03B"},
	7:  {Description: "RC2C03C"},
	8:  {Description: "RC2C05-01"},
	9:  {Description: "RC2C05-02"},
	10: {Description: "RC2C05-03"},
	11: {Description: "RC2C05-04"},
	12: {Description: "RC2C05-05"},
}

var vsHardwareMapper = scalar.UintMap{
	0: {Description: "Vs. Unisystem (normal)"},
	1: {Description: "Vs. Unisystem (RBI Baseball protection)"},
	2: {Description: "Vs. Unisystem (TKO Boxing protection)"},
	3: {Description: "Vs. Unisystem (Super Xevious protection)"},
	4: {Description: "Vs. Unisystem (Vs. Ice Climber Japan protection)"},
	5: {Description: "Vs. Dual System (normal)"},
	6: {Description: "Vs. Dual System (Raid on Bungeling Bay protection)"},
}

var extConsoleMapper = scalar.UintMap{
	0:  {Description: "[Regular NES/Famicom/Dendy]"},
	1:  {Description: "[Nintendo Vs. System]"},
	2:  {Description: "[Playchoice 10]"},
	3:  {Description: "Regular Famiclone, but with CPU that supports Decimal Mode"},
	4:  {Description: "Regular NES/Famicom with EPSM module or plug-through cartridge"},
	5:  {Description: "V.R. Technology VT01 with red/cyan STN palette"},
	6:  {Description: "V.R. Technology VT02"},
	7:  {Description: "V.R. Technology VT03"},
	8:  {Description: "V.R. Technology VT09"},
	9:  {Description: "V.R. Technology VT32"},
	10: {Description: "V.R. Technology VT369"},
	11: {Description: "UMC UM6578"},
	12: {Description: "Famicom Network System"},
}

var expDeviceMapper = scalar.UintMap{
	0:  {Description: "Unspecified"},
	1:  {Description: "Standard NES/Famicom controllers"},
	2:  {Description: "NES Four Score/Satellite with two additional standard controllers"},
	3:  {Description: "Famicom Four Players Adapter with two additional standard controllers using the 'simple' protocol"},
	4:  {Description: "Vs. System (1P via $4016)"},
	5:  {Description: "Vs. System (1P via $4017)"},
	6:  {Description: "Reserved"},
	7:  {Description: "Vs. Zapper"},
	8:  {Description: "Zapper ($4017)"},
	9:  {Description: "Two Zappers"},
	10: {Description: "Bandai Hyper Shot Lightgun"},
	11: {Description: "Power Pad Side A"},
	12: {Description: "Power Pad Side B"},
	13: {Description: "Family Trainer Side A"},
	14: {Description: "Family Trainer Side B"},
	15: {Description: "Arkanoid Vaus Controller (NES)"},
	16: {Description: "Arkanoid Vaus Controller (Famicom)"},
	17: {Description: "Two Vaus Controllers plus Famicom Data Recorder"},
	18: {Description: "Konami Hyper Shot Controller"},
	19: {Description: "Coconuts Pachinko Controller"},
	20: {Description: "Exciting Boxing Punching Bag (Blowup Doll)"},
	21: {Description: "Jissen Mahjong Controller"},
	22: {Description: "Party Tap "},
	23: {Description: "Oeka Kids Tablet"},
	24: {Description: "Sunsoft Barcode Battler"},
	25: {Description: "Miracle Piano Keyboard"},
	26: {Description: "Pokkun Moguraa (Whack-a-Mole Mat and Mallet)"},
	27: {Description: "Top Rider (Inflatable Bicycle)"},
	28: {Description: "Double-Fisted (Requires or allows use of two controllers by one player)"},
	29: {Description: "Famicom 3D System"},
	30: {Description: "Doremikko Keyboard"},
	31: {Description: "R.O.B. Gyro Set"},
	32: {Description: "Famicom Data Recorder ('silent' keyboard)"},
	33: {Description: "ASCII Turbo File"},
	34: {Description: "IGS Storage Battle Box"},
	35: {Description: "Family BASIC Keyboard plus Famicom Data Recorder"},
	36: {Description: "Dongda PEC-586 Keyboard"},
	37: {Description: "Bit Corp. Bit-79 Keyboard"},
	38: {Description: "Subor Keyboard"},
	39: {Description: "Subor Keyboard plus mouse (3x8-bit protocol)"},
	40: {Description: "Subor Keyboard plus mouse (24-bit protocol via $4016)"},
	41: {Description: "SNES Mouse ($4017.d0)"},
	42: {Description: "Multicart"},
	43: {Description: "Two SNES controllers replacing the two standard NES controllers"},
	44: {Description: "RacerMate Bicycle"},
	45: {Description: "U-Force"},
	46: {Description: "R.O.B. Stack-Up"},
	47: {Description: "City Patrolman Lightgun"},
	48: {Description: "Sharp C1 Cassette Interface"},
	49: {Description: "Standard Controller with swapped Left-Right/Up-Down/B-A"},
	50: {Description: "Excalibur Sudoku Pad"},
	51: {Description: "ABL Pinball"},
	52: {Description: "Golden Nugget Casino extra buttons"},
	53: {Description: "Unknown famiclone keyboard used by the 'Golden Key' educational cartridge"},
	54: {Description: "Subor Keyboard plus mouse (24-bit protocol via $4017)"},
	55: {Description: "Port test controller"},
	56: {Description: "Bandai Multi Game Player Gamepad buttons"},
	57: {Description: "Venom TV Dance Mat"},
	58: {Description: "LG TV Remote Control"},
	59: {Description: "Famicom Network Controller"},
	60: {Description: "King Fishing Controller"},
}

func decodeFileHeader(d *decode.D, nc *nesContext) {
	// byte 0-3
	d.FieldRawLen("identifier", 4*8, d.AssertBitBuf([]byte{0x4e, 0x45, 0x53, 0x1a})) // NES<EOF>

	// Peek nes20_identifier
	d.SeekRel(28)
	nes20 := d.PeekUintBits(2)
	d.SeekRel(-28)
	nc.nes20 = nes20 == 0b10

	// byte 4-5
	var prgROMSize0, chrROMSize0 uint64
	if nes20 == 0b10 {
		prgROMSize0 = d.FieldU8("prg_rom_size0")
		chrROMSize0 = d.FieldU8("chr_rom_size0")
	} else {
		prgROMSize := d.FieldU8("prg_rom_size", getROMSizeMapper(16))
		chrROMSize := d.FieldU8("chr_rom_size", getROMSizeMapper(8))
		nc.prgROMSize = romSize(prgROMSize, 0, 16)
		nc.chrROMSize = romSize(chrROMSize, 0, 8)
	}

	// byte 6
	mapper0 := d.FieldU4("mapper0")
	d.FieldU1("alternative_nametables", scalar.UintMapSymStr{0: "no", 1: "yes"})
	trainer := d.FieldU1("trainer", scalar.UintMapSymStr{0: "no", 1: "yes"})
	d.FieldU1("battery", scalar.UintMapSymStr{0: "no", 1: "yes"})
	d.FieldU1("nametable_layout", scalar.UintMapSymStr{0: "vertical", 1: "horizontal"})

	// byte 7
	mapper1 := d.FieldU4("mapper1")
	d.FieldU2("nes20_identifier", scalar.UintMap{0b00: {Sym: "no", Description: "iNES"}, 0b10: {Sym: "yes", Description: "NES 2.0"}})
	consoleType := d.FieldU2("console_type", consoleTypeMapper)

	if nc.nes20 { // NES 2.0
		// byte 8
		d.FieldU4("submapper")
		mapper2 := d.FieldU4("mapper2")
		d.FieldValueUint("mapper", mapper0+(mapper1<<4)+(mapper2<<8))

		// byte 9
		chrROMSize1 := d.FieldU4("chr_rom_size1")
		prgROMSize1 := d.FieldU4("prg_rom_size1")
		nc.chrROMSize = romSize(chrROMSize0, chrROMSize1, 8)
		d.FieldValueUint("chr_rom_size", nc.chrROMSize)
		nc.prgROMSize = romSize(prgROMSize0, prgROMSize1, 16)
		d.FieldValueUint("prg_rom_size", nc.prgROMSize)

		// byte 10
		d.FieldU4("prg_nvram_size", shiftRAMSizeMapper)
		d.FieldU4("prg_ram_size", shiftRAMSizeMapper)

		// byte 11
		d.FieldU4("chr_nvram_size", shiftRAMSizeMapper)
		d.FieldU4("chr_ram_size", shiftRAMSizeMapper)

		// byte 12
		d.FieldU8("cpu_ppu_timing_mode", getFlagMapper(0x3), timingModeMapper)

		// byte 13
		switch consoleType {
		case 1:
			d.FieldU4("vs_hardware_type", vsHardwareMapper)
			d.FieldU4("vs_ppu_type", vsPPUMapper)
		case 3:
			d.FieldU8("ext_console_type", getFlagMapper(0xf, extConsoleMapper))
		default:
			d.FieldU8("byte_13")
		}

		// byte 14
		miscROMs := d.FieldU8("misc_roms", getFlagMapper(0x3))
		nc.miscROMs = (miscROMs & 0x3) > 0

		// byte 15
		d.FieldU8("default_exp_device", getFlagMapper(0x3f, expDeviceMapper))

	} else { // iNES
		mapper := mapper0 + (mapper1 << 4)
		d.FieldValueUint("mapper", mapper)
		d.FieldValueUint("chr_ram_size", iNESCHRRAMSize(mapper, nc.chrROMSize))

		// byte 8
		d.FieldU8("prg_ram_size", multiplyRAMSizeMapper)

		// byte 9
		d.FieldU8("byte_9")

		// byte 10
		d.FieldU8("byte_10")

		// byte 11-15
		d.FieldRawLen("unused", 5*8, d.AssertBitBuf([]byte{0x00, 0x00, 0x00, 0x00, 0x00}))
	}

	if trainer == 1 {
		nc.trainerSize = 512
	} else {
		nc.trainerSize = 0
	}
}

func decodePRGROM(d *decode.D) {
	for !d.End() {
		peek := d.PeekUintBits(8)
		peekTyp := opMap[peek].Type
		peekArgL := ArgLength(peekTyp)

		if d.BitsLeft() < int64(1+peekArgL)*8 {
			d.FieldRawLen("padding", d.BitsLeft())
		} else {
			d.FieldStruct("instruction", func(d *decode.D) {
				op := d.FieldU8("op_code", scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
					opLookup := opMap[s.Actual]
					s.Sym = opLookup.Name

					return s, nil
				}))

				typ := opMap[op].Type
				argL := ArgLength(typ)
				switch argL {
				case 1:
					d.FieldU8("args", GetArgFormatter(typ))
				case 2:
					d.FieldU16("args", GetArgFormatter(typ))
				}
			})
		}
	}
}

func decodeTilePart(d *decode.D) scalar.Uint {
	a := d.U64()
	sym := fmt.Sprintf("%064b", a)

	return scalar.Uint{Actual: a, Sym: sym}
}

func decodeCHRROM(d *decode.D) {
	d.Endian = decode.BigEndian
	for !d.End() {
		d.FieldStruct("tile", func(d *decode.D) {
			lsbsScalar := d.FieldScalarUintFn("pixels_lsb", decodeTilePart)
			msbsScalar := d.FieldScalarUintFn("pixels_msb", decodeTilePart)

			final := make([]rune, 64*3)
			msbsStr := []rune(msbsScalar.SymStr())
			lsbsStr := []rune(lsbsScalar.SymStr())

			for r := 0; r < 8; r++ {
				msbRow := msbsStr[r*8 : r*8+8]
				lsbRow := lsbsStr[r*8 : r*8+8]
				for c := 0; c < 8; c++ {
					final[r*24+c*3] = msbRow[c]
					final[r*24+c*3+1] = lsbRow[c]
					final[r*24+c*3+2] = rune(' ')
				}
			}

			d.FieldValueStr("combined", string(final))
		})
	}
}

func decodeNES(d *decode.D) any {
	var nc nesContext

	d.Endian = decode.LittleEndian

	d.FramedFn(16*8, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) { decodeFileHeader(d, &nc) })
	})

	if nc.trainerSize > 0 {
		d.FramedFn(int64(nc.trainerSize)*8, func(d *decode.D) {
			d.FieldArray("trainer", decodePRGROM)
		})
	}

	d.FramedFn(int64(nc.prgROMSize)*8, func(d *decode.D) {
		d.FieldArray("prg_rom", decodePRGROM)
	})

	d.FramedFn(int64(nc.chrROMSize)*8, func(d *decode.D) {
		d.FieldArray("chr_rom", decodeCHRROM)
	})

	if nc.miscROMs {
		d.FieldRawLen("misc_roms", d.BitsLeft())
	}

	return nil
}
