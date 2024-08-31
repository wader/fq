package midi

import (
	"github.com/wader/fq/pkg/scalar"
)

var notes = scalar.UintMapSymStr{
	uint64(127): "G9",
	uint64(126): "F♯9/G♭9",
	uint64(125): "F9",
	uint64(124): "E9",
	uint64(123): "D♯9/E♭9",
	uint64(122): "D9",
	uint64(121): "C♯9/D♭9",
	uint64(120): "C9",
	uint64(119): "B8",
	uint64(118): "A♯8/B♭8",
	uint64(117): "A8",
	uint64(116): "G♯8/A♭8",
	uint64(115): "G8",
	uint64(114): "F♯8/G♭8",
	uint64(113): "F8",
	uint64(112): "E8",
	uint64(111): "D♯8/E♭8",
	uint64(110): "D8",
	uint64(109): "C♯8/D♭8",
	uint64(108): "C8",
	uint64(107): "B7",
	uint64(106): "A♯7/B♭7",
	uint64(105): "A7",
	uint64(104): "G♯7/A♭7",
	uint64(103): "G7",
	uint64(102): "F♯7/G♭7",
	uint64(101): "F7",
	uint64(100): "E7",
	uint64(99):  "D♯7/E♭7",
	uint64(98):  "D7",
	uint64(97):  "C♯7/D♭7",
	uint64(96):  "C7",
	uint64(95):  "B6",
	uint64(94):  "A♯6/B♭6",
	uint64(93):  "A6",
	uint64(92):  "G♯6/A♭6",
	uint64(91):  "G6",
	uint64(90):  "F♯6/G♭6",
	uint64(89):  "F6",
	uint64(88):  "E6",
	uint64(87):  "D♯6/E♭6",
	uint64(86):  "D6",
	uint64(85):  "C♯6/D♭6",
	uint64(84):  "C6",
	uint64(83):  "B5",
	uint64(82):  "A♯5/B♭5",
	uint64(81):  "A5",
	uint64(80):  "G♯5/A♭5",
	uint64(79):  "G5",
	uint64(78):  "F♯5/G♭5",
	uint64(77):  "F5",
	uint64(76):  "E5",
	uint64(75):  "D♯5/E♭5",
	uint64(74):  "D5",
	uint64(73):  "C♯5/D♭5",
	uint64(72):  "C5",
	uint64(71):  "B4",
	uint64(70):  "A♯4/B♭4",
	uint64(69):  "A4",
	uint64(68):  "G♯4/A♭4",
	uint64(67):  "G4",
	uint64(66):  "F♯4/G♭4",
	uint64(65):  "F4",
	uint64(64):  "E4",
	uint64(63):  "D♯4/E♭4",
	uint64(62):  "D4",
	uint64(61):  "C♯4/D♭4",
	uint64(60):  "C4",
	uint64(59):  "B3",
	uint64(58):  "A♯3/B♭3",
	uint64(57):  "A3",
	uint64(56):  "G♯3/A♭3",
	uint64(55):  "G3",
	uint64(54):  "F♯3/G♭3",
	uint64(53):  "F3",
	uint64(52):  "E3",
	uint64(51):  "D♯3/E♭3",
	uint64(50):  "D3",
	uint64(49):  "C♯3/D♭3",
	uint64(48):  "C3",
	uint64(47):  "B2",
	uint64(46):  "A♯2/B♭2",
	uint64(45):  "A2",
	uint64(44):  "G♯2/A♭2",
	uint64(43):  "G2",
	uint64(42):  "F♯2/G♭2",
	uint64(41):  "F2",
	uint64(40):  "E2",
	uint64(39):  "D♯2/E♭2",
	uint64(38):  "D2",
	uint64(37):  "C♯2/D♭2",
	uint64(36):  "C2",
	uint64(35):  "B1",
	uint64(34):  "A♯1/B♭1",
	uint64(33):  "A1	A1",
	uint64(32):  "G♯1/A♭1",
	uint64(31):  "G1	G1",
	uint64(30):  "F♯1/G♭1",
	uint64(29):  "F1",
	uint64(28):  "E1",
	uint64(27):  "D♯1/E♭1",
	uint64(26):  "D1",
	uint64(25):  "C♯1/D♭1",
	uint64(24):  "C1",
	uint64(23):  "B0",
	uint64(22):  "A♯0/B♭0",
	uint64(21):  "A0",
}

const (
	keyCMajor      = 0x0000
	keyGMajor      = 0x0100
	keyDMajor      = 0x0200
	keyAMajor      = 0x0300
	keyEMajor      = 0x0400
	keyBMajor      = 0x0500
	keyFSharpMajor = 0x0600
	keyCSharpMajor = 0x0700
	keyFMajor      = 0xff00
	keyBFlatMajor  = 0xfe00
	keyEFlatMajor  = 0xfd00
	keyAFlatMajor  = 0xfc00
	keyDFlatMajor  = 0xfb00
	keyGFlatMajor  = 0xfa00
	keyCFlatMajor  = 0xf900

	keyAMinor      = 0x0001
	keyEMinor      = 0x0101
	keyBMinor      = 0x0201
	keyFSharpMinor = 0x0301
	keyCSharpMinor = 0x0401
	keyGSharpMinor = 0x0501
	keyDSharpMinor = 0x0601
	keyASharpMinor = 0x0701
	keyDMinor      = 0xff01
	keyGMinor      = 0xfe01
	keyCMinor      = 0xfd01
	keyFMinor      = 0xfc01
	keyBFlatMinor  = 0xfb01
	keyEFlatMinor  = 0xfa01
	keyAFlatMinor  = 0xf901
)

var keys = scalar.UintMapSymStr{
	keyCMajor:      "C major",
	keyGMajor:      "G major",
	keyDMajor:      "D major",
	keyAMajor:      "A major",
	keyEMajor:      "E major",
	keyBMajor:      "B major",
	keyFSharpMajor: "F♯ major",
	keyCSharpMajor: "C♯ major",
	keyFMajor:      "F major",
	keyBFlatMajor:  "B♭ major",
	keyEFlatMajor:  "E♭ major",
	keyAFlatMajor:  "A♭ major",
	keyDFlatMajor:  "D♭ major",
	keyGFlatMajor:  "G♭ major",
	keyCFlatMajor:  "C♭ major",

	keyAMinor:      "A minor",
	keyEMinor:      "E minor",
	keyBMinor:      "B minor",
	keyFSharpMinor: "F♯ minor",
	keyCSharpMinor: "C♯ minor",
	keyGSharpMinor: "G♯ minor",
	keyDSharpMinor: "D♯ minor",
	keyASharpMinor: "A♯ minor",
	keyDMinor:      "D minor",
	keyGMinor:      "G minor",
	keyCMinor:      "C minor",
	keyFMinor:      "F minor",
	keyBFlatMinor:  "B♭ minor",
	keyEFlatMinor:  "E♭ minor",
	keyAFlatMinor:  "A♭ minor",
}

var controllers = scalar.UintMapSymStr{
	// High resolution continuous controllers (MSB)
	0:  "Bank Select (MSB)",
	1:  "Modulation Wheel (MSB)",
	2:  "Breath Controller (MSB)",
	4:  "Foot Controller (MSB)",
	5:  "Portamento Time (MSB)",
	6:  "Data Entry (MSB)",
	7:  "Channel Volume (MSB)",
	8:  "Balance (MSB)",
	10: "Pan (MSB)",
	11: "Expression Controller (MSB)",
	12: "Effect Control 1 (MSB)",
	13: "Effect Control 2 (MSB)",
	16: "General Purpose Controller 1 (MSB)",
	17: "General Purpose Controller 2 (MSB)",
	18: "General Purpose Controller 3 (MSB)",
	19: "General Purpose Controller 4 (MSB)",

	// High resolution continuous controllers (LSB)
	32: "Bank Select (LSB)",
	33: "Modulation Wheel (LSB)",
	34: "Breath Controller (LSB)",
	36: "Foot Controller (LSB)",
	37: "Portamento Time (LSB)",
	38: "Data Entry (LSB)",
	39: "Channel Volume (LSB)",
	40: "Balance (LSB)",
	42: "Pan (LSB)",
	43: "Expression Controller (LSB)",
	44: "Effect Control 1 (LSB)",
	45: "Effect Control 2 (LSB)",
	48: "General Purpose Controller 1 (LSB)",
	49: "General Purpose Controller 2 (LSB)",
	50: "General Purpose Controller 3 (LSB)",
	51: "General Purpose Controller 4 (LSB)",

	// Switches
	64: "Sustain On/Off",
	65: "Portamento On/Off",
	66: "Sostenuto On/Off",
	67: "Soft Pedal On/Off",
	68: "Legato On/Off",
	69: "Hold 2 On/Off",

	// Low resolution continuous controllers
	70: "Sound Controller 1  (TG: Sound Variation;  FX: Exciter On/Off)",
	71: "Sound Controller 2  (TG: Harmonic Content; FX: Compressor On/Off)",
	72: "Sound Controller 3  (TG: Release Time;     FX: Distortion On/Off)",
	73: "Sound Controller 4  (TG: Attack Time;      FX: EQ On/Off)",
	74: "Sound Controller 5  (TG: Brightness;       FX: Expander On/Off)",
	75: "Sound Controller 6  (TG: Decay Time;       FX: Reverb On/Off)",
	76: "Sound Controller 7  (TG: Vibrato Rate;     FX: Delay On/Off)",
	77: "Sound Controller 8  (TG: Vibrato Depth;    FX: Pitch Transpose On/Off)",
	78: "Sound Controller 9  (TG: Vibrato Delay;    FX: Flange/Chorus On/Off)",
	79: "Sound Controller 10 (TG: Undefined;        FX: Special Effects On/Off)",
	80: "General Purpose Controller 5",
	81: "General Purpose Controller 6",
	82: "General Purpose Controller 7",
	83: "General Purpose Controller 8",
	84: "Portamento Control",
	88: "High Resolution Velocity Prefix",
	91: "Effects 1 Depth (Reverb Send Level)",
	92: "Effects 2 Depth (Tremolo Depth)",
	93: "Effects 3 Depth (Chorus Send Level)",
	94: "Effects 4 Depth (Celeste Depth)",
	95: "Effects 5 Depth (Phaser Depth)",

	// RPNs / NRPNs
	96:  "Data Increment",
	97:  "Data Decrement",
	98:  "Non-Registered Parameter Number (LSB)",
	99:  "Non-Registered Parameter Number (MSB)",
	100: "Registered Parameter Number (LSB)",
	101: "Registered Parameter Number (MSB)",

	// Channel Mode messages
	120: "All Sound Off",
	121: "Reset All Controllers",
	122: "Local Control On/Off",
	123: "All Notes Off",
	124: "Omni Mode Off",
	125: "Omni Mode On ",
	126: "Mono Mode On",
	127: "Poly Mode On",
}

var manufacturers = scalar.UintMapSymStr{
	// special purpose

	0x7D: "Non-Commercial",
	0x7E: "Non-RealTime Extensions",
	0x7F: "RealTime Extensions",

	// American
	0x01: "Sequential Circuits",
	0x04: "Moog",
	0x05: "Passport Designs",
	0x06: "Lexicon",
	0x07: "Kurzweil",
	0x08: "Fender",
	0x0A: "AKG Acoustics",
	0x0F: "Ensoniq",
	0x10: "Oberheim",
	0x11: "Apple",
	0x13: "Digidesign",
	0x18: "Emu",
	0x1A: "ART",
	0x1C: "Eventide",

	// European
	0x22: "Synthaxe",
	0x24: "Hohner",
	0x29: "PPG",
	0x2B: "SSL",
	0x2D: "Hinton Instruments",
	0x2F: "Elka / General Music",
	0x30: "Dynacord",
	0x33: "Clavia (Nord)",
	0x36: "Cheetah",
	0x3E: "Waldorf Electronics Gmbh",

	// Japanese
	0x40: "Kawai",
	0x41: "Roland",
	0x42: "Korg",
	0x43: "Yamaha",
	0x44: "Casio",
	0x47: "Akai",
	0x48: "Japan Victor (JVC)",
	0x4C: "Sony",
	0x4E: "Teac Corporation",
	0x51: "Fostex",
	0x52: "Zoom",
}

var manufacturers_extended = scalar.UintMapSymStr{
	0x0007: "Digital Music Corporation",
	0x0009: "New England Digital",
	0x000E: "Alesis",
	0x0015: "KAT",
	0x0016: "Opcode",
	0x001A: "Allen & Heath Brenell",
	0x001B: "Peavey Electronics",
	0x001C: "360 Systems",
	0x001F: "Zeta Systems",
	0x0020: "Axxes",
	0x003B: "Mark Of The Unicorn (MOTU)",
	0x004D: "Studio Electronics",
	0x0050: "MIDI Solutions Inc",
	0x0137: "Roger Linn Design",
	0x0172: "Kilpatrick Audio",
	0x0173: "iConnectivity",
	0x0214: "Intellijel Designs Inc",

	// European
	0x2011: "Forefront Technology",
	0x2013: "Kenton Electronics",
	0x201F: "TC Electronic",
	0x2020: "Doepfer",
	0x2027: "Acorn Computer",
	0x2029: "Focusrite / Novation",
	0x2032: "Behringer",
	0x2033: "Access Music Electronics",
	0x203C: "Elektron",
	0x204D: "Vermona",
	0x2052: "Analogue Systems",
	0x2069: "Elby Designs",
	0x206B: "Arturia",
	0x2076: "Teenage Engineering",
	0x2102: "Mutable Instruments",
	0x2109: "Native Instruments",
	0x2110: "ROLI Ltd",
	0x211A: "IK Multimedia",
	0x211C: "Modor Music",
	0x211D: "Ableton",
	0x2127: "Expert Sleepers",
}

var framerates = scalar.UintMapSymUint{
	0: 24,
	1: 25,
	2: 29,
	3: 30,
}
