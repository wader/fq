package tzx

// https://worldofspectrum.net/TZXformat.html

import (
	"embed"

	"golang.org/x/text/encoding/charmap"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed tzx.md
var tzxFS embed.FS

var tapFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.TZX,
		&decode.Format{
			Description: "TZX tape format for ZX Spectrum computers",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    tzxDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.TAP}, Out: &tapFormat},
			},
		})
	interp.RegisterFS(tzxFS)
}

func tzxDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldRawLen("signature", 8*8, d.AssertBitBuf([]byte("ZXTape!\x1A")))
	d.FieldU8("major_version")
	d.FieldU8("minor_version")
	decodeBlocks(d)

	return nil
}

func decodeBlocks(d *decode.D) {
	d.FieldArray("blocks", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("block", func(d *decode.D) {
				decodeBlock(d)
			})
		}
	})
}

func decodeBlock(d *decode.D) {
	blocks := map[uint64]func(d *decode.D){
		// ID: 10h (16d) | Standard Speed Data
		// This block is replayed with the standard Spectrum ROM timing values
		// (the values in curly brackets in block ID 11). The pilot tone
		// consists of 8063 pulses if the first data byte (the flag byte)
		// is < 128, 3223 otherwise.
		0x10: func(d *decode.D) {
			// Pause after this block (ms.) {1000}
			d.FieldU16("pause")

			// A single TAP Data Block
			peekBytes := d.PeekBytes(2)                              // get the TAP data block length
			length := uint16(peekBytes[1])<<8 | uint16(peekBytes[0]) // bytes are stored in LittleEndian
			length += 2                                              // include the two bytes for this value
			d.FieldFormatLen("tap", int64(length)*8, &tapFormat, nil)
		},

		// ID: 11h (17d) | Turbo Speed Data
		// This block is very similar to the normal TAP block but with some
		// additional info on the timings and other important differences. The
		// same tape encoding is used as for the standard speed data block. If
		// a block should use some non-standard sync or pilot tones (i.e. all
		// sorts of protection schemes) then the next three blocks describe it.
		0x11: func(d *decode.D) {
			d.FieldU16("pilot_pulse")  // Length of PILOT pulse {2168}
			d.FieldU16("sync_pulse_1") // Length of SYNC first pulse {667}
			d.FieldU16("sync_pulse_2") // Length of SYNC second pulse {735}
			d.FieldU16("bit0_pulse")   // Length of ZERO bit pulse {855}
			d.FieldU16("bit1_pulse")   // Length of ONE bit pulse {1710}

			// Length of PILOT tone (number of pulses)
			// {8063 header (flag<128), 3223 data (flag>=128)}
			d.FieldU16("pilot_tone")

			// Used bits in the last byte (other bits should be 0) {8}
			// e.g. if this is 6, then the bits used (x) in the last byte are: xxxxxx00,
			// where MSb is the leftmost bit, LSb is the rightmost bit
			d.FieldU8("used_bits")

			d.FieldU16("pause")            // Pause after this block (ms.) {1000}
			length := d.FieldU24("length") // Length of data that follows

			// Data as in .TAP files
			d.FieldRawLen("data", int64(length)*8)
		},

		// ID: 12h (18d) | Pure Tone
		// This will produce a tone which is basically the same as the pilot
		// tone in 10h and 11h blocks.
		0x12: func(d *decode.D) {
			d.FieldU16("pulse_length") // Length of one pulse in T-states
			d.FieldU16("pulse_count")  // Number of pulses
		},

		// ID: 13h (19d) | Sequence of Pulses
		// This will produce N pulses, each having its own timing. Up to 255
		// pulses can be stored in this block.
		0x13: func(d *decode.D) {
			count := d.FieldU8("pulse_count")
			d.FieldArray("pulses", func(d *decode.D) {
				for i := uint64(0); i < count; i++ {
					d.FieldU16("pulse")
				}
			})
		},

		// ID: 14h (20d) | Pure Data
		// This is the same as in the turbo loading data block, except that it
		// has no pilot or sync pulses.
		0x14: func(d *decode.D) {
			d.FieldU16("bit0_pulse")       // Length of ZERO bit pulse
			d.FieldU16("bit1_pulse")       // Length of ONE bit pulse
			d.FieldU8("used_bits")         // Used bits in last byte
			d.FieldU16("pause")            // Pause after this block (ms.)
			length := d.FieldU24("length") // Length of data that follows

			// Data as in .TAP files
			d.FieldRawLen("data", int64(length)*8)
		},

		// ID: 15h (21d) | Direct Recording
		// This block is used for tapes which have some parts in a format such
		// that the turbo loader block cannot be used. This is not like a VOC
		// file since the information is much more compact. Each sample value
		// is represented by one bit only (0 for low, 1 for high) which means
		// that the block will be at most 1/8 the size of the equivalent VOC.
		// The preferred sampling frequencies are 22050 or 44100 Hz
		// (158 or 79 T-states/sample).
		0x15: func(d *decode.D) {
			d.FieldU16("t_states")                 // Number of T-states per sample (bit of data)
			d.FieldU16("pause")                    // Pause after this block in milliseconds (ms.)
			d.FieldU8("used_bits")                 // Used bits (samples) in last byte of data (1-8)
			length := d.FieldU24("length")         // Length of data that follows
			d.FieldRawLen("data", int64(length)*8) // Samples data. Each bit represents a state on the EAR port
		},

		// ID: 18h (24d) | CSW Recording
		// This block contains a sequence of raw pulses encoded in CSW format
		// v2 (Compressed Square Wave).
		0x18: func(d *decode.D) {
			length := d.FieldU32("length") // Block length (without these four bytes)

			// NOTE: remove these next 4 fields from the length so
			// the data size is calculated correctly
			length -= 2 + 3 + 1 + 4

			// Pause after this block (in ms)
			d.FieldU16("pause")
			// Sampling rate
			d.FieldU24("sample_rate")
			// Compression type
			d.FieldU8("compression_type", scalar.UintMapSymStr{0x00: "unknown", 0x01: "rle", 0x02: "zrle"})
			// Number of stored pulses (after decompression)
			d.FieldU32("stored_pulse_count")

			// CSW data, encoded according to the CSW specification
			d.FieldRawLen("data", int64(length)*8)
		},

		// ID: 19h (25d) | Generalized Data
		// This block was developed to represent an extremely wide range of data
		// encoding techniques. Each loading component (pilot tone, sync pulses,
		// data) is associated to a specific sequence of pulses, where each
		// sequence (wave) can contain a different number of pulses from the
		// others. In this way it is possible to have a situation where bit 0 is
		// represented with 4 pulses and bit 1 with 8 pulses.
		0x19: func(d *decode.D) {
			length := d.FieldU32("length") // Block length (without these four bytes)
			// TBD:
			//	Pause        uint16     // Pause after this block (ms)
			//	TOTP         uint32     // Total number of symbols in pilot/sync block (can be 0)
			//	NPP          uint8      // Maximum number of pulses per pilot/sync symbol
			//	ASP          uint8      // Number of pilot/sync symbols in the alphabet table (0=256)
			//	TOTD         uint32     // Total number of symbols in data stream (can be 0)
			//	NPD          uint8      // Maximum number of pulses per data symbol
			//	ASD          uint8      // Number of data symbols in the alphabet table (0=256)
			//	PilotSymbols []Symbol   // 0x12  SYMDEF[ASP] Pilot and sync symbols definition table
			//	PilotStreams []PilotRLE // 0x12+ (2*NPP+1)*ASP - PRLE[TOTP]  Pilot and sync data stream
			//	DataSymbols  []Symbol   // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+TOTP*3  - SYMDEF[ASD] Data symbols definition table
			//	DataStreams  []uint8    // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+ TOTP*3+(2*NPD+1)*ASD - BYTE[DS]  Data stream
			d.FieldRawLen("data", int64(length)*8)
		},

		// ID: 20h (32d) | Pause Tape Command
		// This will make a silence (low amplitude level (0)) for a given time
		// in milliseconds. If the value is 0 then the emulator or utility should
		// (in effect) STOP THE TAPE, until the user or emulator requests it.
		0x20: func(d *decode.D) {
			d.FieldU16("pause") // Pause duration in ms.
		},

		// ID: 21h (33d) | Group Start
		// This block marks the start of a group of blocks which are to be
		// treated as one single (composite) block. For each group start block
		// there must be a group end block. Nesting of groups is not allowed.
		0x21: func(d *decode.D) {
			length := d.FieldU8("length")
			d.FieldStr("group_name", int(length), charmap.ISO8859_1)
		},

		// ID: 22h (34d) | Group End
		// This indicates the end of a group. This block has no body.
		0x22: func(d *decode.D) {},

		// JumpTo
		// ID: 23h (35d)
		// This block will allow for jumping from one block to another within
		// the file. All blocks are included in the block count!
		0x23: func(d *decode.D) {
			d.FieldS16("value", scalar.SintMapSymStr{
				0:  "loop_forever",
				1:  "next_block",
				2:  "skip_block",
				-1: "prev_block",
			})
		},

		// ID: 24h (36d) | Loop Start
		// Indicates a sequence of identical blocks, or of identical groups of
		// blocks. This block is the same as the FOR statement in BASIC.
		0x24: func(d *decode.D) {
			d.FieldU16("repetitions") // Number of repetitions (greater than 1)
		},

		// ID: 25h (37d) | Loop End
		// This is the same as BASIC's NEXT statement. It means that the utility
		// should jump back to the start of the loop if it hasn't been run for
		// the specified number of times. This block has no body.
		0x25: func(d *decode.D) {},

		// ID: 26h (38d) | Call Sequence
		// This block is an analogue of the CALL Subroutine statement. It
		// basically executes a sequence of blocks that are somewhere else and
		// then goes back to the next block. Because more than one call can be
		// normally used you can include a list of sequences to be called. CALL
		// blocks can be used in the LOOP sequences and vice versa. The value
		// is relative so that you can add some blocks in the beginning of the
		// file without disturbing the call values.
		// Look at 'Jump To Block' for reference on the values.
		0x26: func(d *decode.D) {
			count := d.FieldU16("count")
			d.FieldArray("call_blocks", func(d *decode.D) {
				for i := uint64(0); i < count; i++ {
					d.FieldS16("offset")
				}
			})
		},

		// ID: 27h (39d) | Return From Sequence
		// This block indicates the end of the Called Sequence. The next block
		// played will be the block after the last CALL block (or the next Call,
		// if the Call block had multiple calls). This block has no body.
		0x27: func(d *decode.D) {},

		// ID: 28h (40d) | Select
		// This block is useful when the tape consists of two or more separately
		// loadable parts. With this block it is possible to select one of the
		// parts and the utility/emulator will start loading from that block.
		// All offsets are relative signed words.
		0x28: func(d *decode.D) {
			// Length of the whole block (without these two bytes)
			d.FieldU16("length")

			count := d.FieldU8("count")
			d.FieldArray("selections", func(d *decode.D) {
				for i := 0; i < int(count); i++ {
					d.FieldStruct("selection", func(d *decode.D) {
						d.FieldS16("offset")          // Relative Offset as `signed` value
						length := d.FieldU8("length") // Length of description text (max 30 chars)
						d.FieldStr("description", int(length), charmap.ISO8859_1)
					})
				}
			})
		},

		// ID: 2Ah (42d) | Stop Tape When 48k Mode
		// When this block is encountered, the tape will stop ONLY if the machine
		// is an 48K Spectrum. This block is to be used for multi-loading games
		// that load one level at a time in 48K mode, but load the entire tape at
		// once if in 128K mode.
		// This block has no body of its own, but follows the extension rule.
		0x2A: func(d *decode.D) {
			d.FieldU32("length") // Length of the block without these four bytes (0)
		},

		// ID: 2Bh (43d) | Set Signal Level
		// This block sets the current signal level to the specified value
		// (high or low). It should be used whenever it is necessary to avoid
		// any ambiguities, e.g. with custom loaders which are level-sensitive.
		0x2B: func(d *decode.D) {
			d.FieldU32("length") // Block length (without these four bytes)
			d.FieldU8("signal_level", scalar.UintMapSymStr{0: "low", 1: "high"})
		},

		// ID: 30h (48d) | Text Description
		// This is meant to identify parts of the tape, such as where level 1
		// starts, where to rewind to when the game ends, etc. This description
		// is not guaranteed to be shown while the tape is playing, but can be
		// read while browsing the tape or changing the tape pointer.
		// The description can be up to 255 characters long.
		0x30: func(d *decode.D) {
			length := d.FieldU8("length")
			d.FieldStr("description", int(length), charmap.ISO8859_1)
		},

		// ID: 31h (49d) | Message
		// This will enable the emulators to display a message for a given time.
		// This should not stop the tape and it should not make silence. If the
		// time is 0 then the emulator should wait for the user to press a key.
		0x31: func(d *decode.D) {
			// Time (in seconds) for which the message should be displayed
			d.FieldU8("display_time")
			// Length of the text message
			length := d.FieldU8("length")
			// Message that should be displayed in ASCII format
			d.FieldStr("message", int(length), charmap.ISO8859_1)
		},

		// ID: 32h (50d) | Archive Info
		// This optional block is used at the beginning of the tape containing
		// various metadata about the tape.
		0x32: func(d *decode.D) {
			d.FieldU16("length")        // Length of the whole block without these two bytes
			count := d.FieldU8("count") // Number of entries in the archive info

			// the archive strings
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < count; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						d.FieldU8("id", scalar.UintMapSymStr{
							0x00: "title",
							0x01: "publisher",
							0x02: "author",
							0x03: "year",
							0x04: "language",
							0x05: "category",
							0x06: "price",
							0x07: "loader",
							0x08: "origin",
							0xFF: "comment",
						})
						length := d.FieldU8("length")
						d.FieldStr("value", int(length), charmap.ISO8859_1)
					})
				}
			})
		},

		// ID: 33h (51d) | Hardware Type
		// This blocks contains information about the hardware that the programs
		// on this tape use.
		0x33: func(d *decode.D) {
			// Number of machines and hardware types for which info is supplied
			count := d.FieldU8("count")
			d.FieldArray("hardware_info", func(d *decode.D) {
				for i := uint64(0); i < count; i++ {
					d.FieldStruct("info", func(d *decode.D) {
						// Hardware Type ID (computers, printers, mice, etc.)
						typeId := d.FieldU8("type_id", hwInfoTypes)
						// Hardware Device ID (ZX81, Kempston Joystick, etc.)
						d.FieldU8("device_id", hwInfoDevices[typeId])
						// Hardware compatibility information
						d.FieldU8("info_id", hwInfoCompatibilityInfo)
					})
				}
			})
		},

		// ID: 35h (53d) | Custom Info
		// This block contains various custom data. For example, it might contain
		// some information written by a utility, extra settings required by a
		// particular emulator, etc.
		0x35: func(d *decode.D) {
			d.FieldStr("identification", 16, charmap.ISO8859_1)
			length := d.FieldU32("length")
			d.FieldRawLen("info", int64(length)*8)
		},

		// ID: 5Ah (90d) | Glue Block
		// This block is generated when two ZX Tape files are merged together.
		// It is here so that you can easily copy the files together and use
		// them. Of course, this means that resulting file would be 10 bytes
		// longer than if this block was not used. All you have to do if you
		// encounter this block ID is to skip next 9 bytes. If you can avoid
		// using this block for this purpose, then do so; it is preferable to
		// use a utility to join the two files and ensure that they are both
		// of the higher version number.
		0x5A: func(d *decode.D) {
			// Value: { "XTape!",0x1A,MajR,MinR }
			// Just skip these 9 bytes and you will end up on the next ID.
			d.FieldRawLen("value", 9*8)
		},
	}

	blockType := d.PeekUintBits(8)
	// Deprecated block types: C64RomType, C64TurboData, EmulationInfo, Snapshot
	if blockType == 0x16 || blockType == 0x17 || blockType == 0x34 || blockType == 0x40 {
		d.Fatalf("deprecated block type encountered: %02x", blockType)
	}

	blockLabel := blockTypeMapper[blockType]
	d.FieldStruct(blockLabel, func(d *decode.D) {
		d.FieldU8("type", blockTypeMapper)

		if fn, ok := blocks[blockType]; ok {
			fn(d)
		} else {
			d.Fatalf("block type not valid, got: %02x", blockType)
		}
	})
}

var blockTypeMapper = scalar.UintMapSymStr{
	0x10: "standard_speed_data",
	0x11: "turbo_speed_data",
	0x12: "pure_tone",
	0x13: "sequence_of_pulses",
	0x14: "pure_data",
	0x15: "direct_recording", // deprecated
	0x16: "c64_rom_type",     // deprecated
	0x17: "c64_turbo_data",
	0x18: "csw_recording",
	0x19: "generalized_data",
	0x20: "pause_tape_command",
	0x21: "group_start",
	0x22: "group_end",
	0x23: "jump_to",
	0x24: "loop_start",
	0x25: "loop_end",
	0x26: "call_sequence",
	0x27: "return_from_sequence",
	0x28: "select",
	0x2A: "stop_tape_when_48k_mode",
	0x2B: "set_signal_level",
	0x30: "text_description",
	0x31: "message",
	0x32: "archive_info",
	0x33: "hardware_type",
	0x34: "emulation_info", // deprecated
	0x35: "custom_info",
	0x40: "snapshot", // deprecated
	0x5A: "glue_block",
}

var hwInfoTypes = scalar.UintMapDescription{
	0x00: "Computers",
	0x01: "External storage",
	0x02: "ROM/RAM type add-ons",
	0x03: "Sound devices",
	0x04: "Joysticks",
	0x05: "Mice",
	0x06: "Other controllers",
	0x07: "Serial ports",
	0x08: "Parallel ports",
	0x09: "Printers",
	0x0a: "Modems",
	0x0b: "Digitizers",
	0x0c: "Network adapters",
	0x0d: "Keyboards & keypads",
	0x0e: "AD/DA converters",
	0x0f: "EPROM programmers",
	0x10: "Graphics",
}

var hwInfoDevices = map[uint64]scalar.UintMapDescription{
	0x00: { // Computers
		0x00: "ZX Spectrum 16k",
		0x01: "ZX Spectrum 48k, Plus",
		0x02: "ZX Spectrum 48k ISSUE 1",
		0x03: "ZX Spectrum 128k +(Sinclair)",
		0x04: "ZX Spectrum 128k +2 (grey case)",
		0x05: "ZX Spectrum 128k +2A, +3",
		0x06: "Timex Sinclair TC-2048",
		0x07: "Timex Sinclair TS-2068",
		0x08: "Pentagon 128",
		0x09: "Sam Coupe",
		0x0a: "Didaktik M",
		0x0b: "Didaktik Gama",
		0x0c: "ZX-80",
		0x0d: "ZX-81",
		0x0e: "ZX Spectrum 128k, Spanish version",
		0x0f: "ZX Spectrum, Arabic version",
		0x10: "Microdigital TK 90-X",
		0x11: "Microdigital TK 95",
		0x12: "Byte",
		0x13: "Elwro 800-3 ",
		0x14: "ZS Scorpion 256",
		0x15: "Amstrad CPC 464",
		0x16: "Amstrad CPC 664",
		0x17: "Amstrad CPC 6128",
		0x18: "Amstrad CPC 464+",
		0x19: "Amstrad CPC 6128+",
		0x1a: "Jupiter ACE",
		0x1b: "Enterprise",
		0x1c: "Commodore 64",
		0x1d: "Commodore 128",
		0x1e: "Inves Spectrum+",
		0x1f: "Profi",
		0x20: "GrandRomMax",
		0x21: "Kay 1024",
		0x22: "Ice Felix HC 91",
		0x23: "Ice Felix HC 2000",
		0x24: "Amaterske RADIO Mistrum",
		0x25: "Quorum 128",
		0x26: "MicroART ATM",
		0x27: "MicroART ATM Turbo 2",
		0x28: "Chrome",
		0x29: "ZX Badaloc",
		0x2a: "TS-1500",
		0x2b: "Lambda",
		0x2c: "TK-65",
		0x2d: "ZX-97",
	},
	0x01: { // External storage
		0x00: "ZX Microdrive",
		0x01: "Opus Discovery",
		0x02: "MGT Disciple",
		0x03: "MGT Plus-D",
		0x04: "Rotronics Wafadrive",
		0x05: "TR-DOS (BetaDisk)",
		0x06: "Byte Drive",
		0x07: "Watsford",
		0x08: "FIZ",
		0x09: "Radofin",
		0x0a: "Didaktik disk drives",
		0x0b: "BS-DOS (MB-02)",
		0x0c: "ZX Spectrum +3 disk drive",
		0x0d: "JLO (Oliger) disk interface",
		0x0e: "Timex FDD3000",
		0x0f: "Zebra disk drive",
		0x10: "Ramex Millennia",
		0x11: "Larken",
		0x12: "Kempston disk interface",
		0x13: "Sandy",
		0x14: "ZX Spectrum +3e hard disk",
		0x15: "ZXATASP",
		0x16: "DivIDE",
		0x17: "ZXCF",
	},
	0x02: { // ROM/RAM type add_ons
		0x00: "Sam Ram",
		0x01: "Multiface ONE",
		0x02: "Multiface 128k",
		0x03: "Multiface +3",
		0x04: "MultiPrint",
		0x05: "MB-02 ROM/RAM expansion",
		0x06: "SoftROM",
		0x07: "1k",
		0x08: "16k",
		0x09: "48k",
		0x0a: "Memory in 8-16k used",
	},
	0x03: { // Sound devices
		0x00: "Classic AY hardware (compatible with 128k ZXs)",
		0x01: "Fuller Box AY sound hardware",
		0x02: "Currah microSpeech",
		0x03: "SpecDrum",
		0x04: "AY ACB stereo (A+C=left, B+C=right); Melodik",
		0x05: "AY ABC stereo (A+B=left, B+C=right)",
		0x06: "RAM Music Machine",
		0x07: "Covox",
		0x08: "General Sound",
		0x09: "Intec Electronics Digital Interface B8001",
		0x0a: "Zon-X AY",
		0x0b: "QuickSilva AY",
		0x0c: "Jupiter ACE",
	},
	0x04: { // Joysticks
		0x00: "Kempston",
		0x01: "Cursor, Protek, AGF",
		0x02: "Sinclair 2 Left (12345)",
		0x03: "Sinclair 1 Right (67890)",
		0x04: "Fuller",
	},
	0x05: { // Mice
		0x00: "AMX mouse",
		0x01: "Kempston mouse",
	},
	0x06: { // Other controllers
		0x00: "Trickstick",
		0x01: "ZX Light Gun",
		0x02: "Zebra Graphics Tablet",
		0x03: "Defender Light Gun",
	},
	0x07: { // Serial ports
		0x00: "ZX Interface 1",
		0x01: "ZX Spectrum 128k",
	},
	0x08: { // Parallel ports
		0x00: "Kempston S",
		0x01: "Kempston E",
		0x02: "ZX Spectrum +3",
		0x03: "Tasman",
		0x04: "DK'Tronics",
		0x05: "Hilderbay",
		0x06: "INES Printerface",
		0x07: "ZX LPrint Interface 3",
		0x08: "MultiPrint",
		0x09: "Opus Discovery",
		0x0a: "Standard 8255 chip with ports 31,63,95",
	},
	0x09: { // Printers
		0x00: "ZX Printer, Alphacom 32 & compatibles",
		0x01: "Generic printer",
		0x02: "EPSON compatible",
	},
	0x0a: { // Modems
		0x00: "Prism VTX 5000",
		0x01: "T/S 2050 or Westridge 2050",
	},
	0x0b: { // Digitizers
		0x00: "RD Digital Tracer",
		0x01: "DK'Tronics Light Pen",
		0x02: "British MicroGraph Pad",
		0x03: "Romantic Robot Videoface",
	},
	0x0c: { // Network adapters
		0x00: "ZX Interface 1",
	},
	0x0d: { // Keyboards & keypads
		0x00: "Keypad for ZX Spectrum 128k",
	},
	0x0e: { // AD/DA converters
		0x00: "Harley Systems ADC 8.2",
		0x01: "Blackboard Electronics",
	},
	0x0f: { // EPROM programmers
		0x00: "Orme Electronics",
	},
	0x10: { // Graphics
		0x00: "WRX Hi-Res",
		0x01: "G007",
		0x02: "Memotech",
		0x03: "Lambda Colour",
	},
}

var hwInfoCompatibilityInfo = scalar.UintMapDescription{
	00: "RUNS on this machine or with this hardware, but may or may not use the hardware or special features of the machine.",
	01: "USES the hardware or special features of the machine, such as extra memory or a sound chip.",
	02: "RUNS but it DOESN'T use the hardware or special features of the machine.",
	03: "DOESN'T RUN on this machine or with this hardware.",
}
