package tzx

// https://worldofspectrum.net/zx-modules/fileformats/tapformat.html

import (
	"embed"

	"golang.org/x/text/encoding/charmap"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed tap.md
var tapFS embed.FS

func init() {
	interp.RegisterFormat(
		format.TAP,
		&decode.Format{
			Description: "TAP tape format for ZX Spectrum computers",
			DecodeFn:    tapDecoder,
		})
	interp.RegisterFS(tapFS)
}

// The TAP- (and BLK-) format is nearly a direct copy of the data that is stored
// in real tapes, as it is written by the ROM save routine of the ZX-Spectrum.
// A TAP file is simply one data block or a group of 2 or more data blocks, one
// followed after the other. The TAP file may be empty.
func tapDecoder(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var ti format.TAP_In
	d.ArgAs(&ti)

	if ti.ReadOneBlock {
		decodeTapBlock(d)
		return nil
	}

	d.FieldArray("blocks", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("block", func(d *decode.D) {
				decodeTapBlock(d)
			})
		}
	})
	return nil
}

func decodeTapBlock(d *decode.D) {
	// Length of the following data.
	length := d.FieldU16("length")

	// read header, fragment, or data block
	switch length {
	case 0:
		// fragment with no data
	case 1:
		d.FieldRawLen("data", 8)
	case 19:
		d.FieldStruct("header", func(d *decode.D) {
			decodeHeader(d)
		})
	default:
		d.FieldStruct("data", func(d *decode.D) {
			decodeDataBlock(d, length)
		})
	}
}

// decodes the different types of 19-byte header blocks.
func decodeHeader(d *decode.D) {
	// Always 0: byte indicating a standard ROM loading header
	d.FieldU8("flag", scalar.UintMapSymStr{0: "standard_speed_data"})
	// Header type
	dataType := d.FieldU8("data_type", scalar.UintMapSymStr{
		0x00: "program",
		0x01: "numeric",
		0x02: "alphanumeric",
		0x03: "data",
	})
	// Loading name of the program. Filled with spaces (0x20) to 10 characters.
	d.FieldStr("program_name", 10, charmap.ISO8859_1)

	switch dataType {
	case 0:
		// Length of data following the header = length of BASIC program + variables.
		d.FieldU16("data_length")
		// LINE parameter of SAVE command. Value 32768 means "no auto-loading".
		// 0..9999 are valid line numbers.
		d.FieldU16("auto_start_line")
		// Length of BASIC program;
		// remaining bytes ([data length] - [program length]) = offset of variables.
		d.FieldU16("program_length")
	case 1:
		// Length of data following the header = length of number array * 5 + 3.
		d.FieldU16("data_length")
		// Unused byte.
		d.FieldU8("unused0")
		// (1..26 meaning A..Z) + 128.
		d.FieldU8("variable_name", scalar.UintHex)
		// UnusedWord: 32768.
		d.FieldU16("unused1")
	case 2:
		// Length of data following the header = length of string array + 3.
		d.FieldU16("data_length")
		// Unused byte.
		d.FieldU8("unused0")
		// (1..26 meaning A$..Z$) + 192.
		d.FieldU8("variable_name", scalar.UintHex)
		// UnusedWord: 32768.
		d.FieldU16("unused1")
	case 3:
		// Length of data following the header, in case of a SCREEN$ header = 6912.
		d.FieldU16("data_length")
		// In case of a SCREEN$ header = 16384.
		d.FieldU16("start_address", scalar.UintHex)
		//	UnusedWord: 32768.
		d.FieldU16("unused")
	default:
		d.Fatalf("invalid TAP header type, got: %d", dataType)
	}

	// Simply all bytes XORed (including flag byte).
	d.FieldU8("checksum", scalar.UintHex)
}

func decodeDataBlock(d *decode.D, length uint64) {
	// flag indicating the type of data block, usually 255 (standard speed data)
	d.FieldU8("flag", scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		if s.Actual == 0xFF {
			s.Sym = "standard_speed_data"
		} else {
			s.Sym = "custom_data_block"
		}
		return s, nil
	}))
	// The essential data: length minus the flag/checksum bytes (may be empty)
	d.FieldRawLen("data", int64(length-2)*8)
	// Simply all bytes (including flag byte) XORed
	d.FieldU8("checksum", scalar.UintHex)
}
