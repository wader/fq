package midi

import (
	"github.com/wader/fq/pkg/scalar"
)

var manufacturers = scalar.StrMapSymStr{
	// special purpose

	"7D": "Non-Commercial",
	"7E": "Non-RealTime Extensions",
	"7F": "RealTime Extensions",

	// American
	"01": "Sequential Circuits",
	"04": "Moog",
	"05": "Passport Designs",
	"06": "Lexicon",
	"07": "Kurzweil",
	"08": "Fender",
	"0A": "AKG Acoustics",
	"0F": "Ensoniq",
	"10": "Oberheim",
	"11": "Apple",
	"13": "Digidesign",
	"18": "Emu",
	"1A": "ART",
	"1C": "Eventide",

	// European
	"22": "Synthaxe",
	"24": "Hohner",
	"29": "PPG",
	"2B": "SSL",
	"2D": "Hinton Instruments",
	"2F": "Elka / General Music",
	"30": "Dynacord",
	"33": "Clavia (Nord)",
	"36": "Cheetah",
	"3E": "Waldorf Electronics Gmbh",

	// Japanese
	"40": "Kawai",
	"41": "Roland",
	"42": "Korg",
	"43": "Yamaha",
	"44": "Casio",
	"47": "Akai",
	"48": "Japan Victor (JVC)",
	"4C": "Sony",
	"4E": "Teac Corporation",
	"51": "Fostex",
	"52": "Zoom",

	// American
	"0007": "Digital Music Corporation",
	"0009": "New England Digital",
	"000E": "Alesis",
	"0015": "KAT",
	"0016": "Opcode",
	"001A": "Allen & Heath Brenell",
	"001B": "Peavey Electronics",
	"001C": "360 Systems",
	"001F": "Zeta Systems",
	"0020": "Axxes",
	"003B": "Mark Of The Unicorn (MOTU)",
	"004D": "Studio Electronics",
	"0050": "MIDI Solutions Inc",
	"0137": "Roger Linn Design",
	"0172": "Kilpatrick Audio",
	"0173": "iConnectivity",
	"0214": "Intellijel Designs Inc",

	// European
	"2011": "Forefront Technology",
	"2013": "Kenton Electronics",
	"201F": "TC Electronic",
	"2020": "Doepfer",
	"2027": "Acorn Computer",
	"2029": "Focusrite / Novation",
	"2032": "Behringer",
	"2033": "Access Music Electronics",
	"203C": "Elektron",
	"204D": "Vermona",
	"2052": "Analogue Systems",
	"2069": "Elby Designs",
	"206B": "Arturia",
	"2076": "Teenage Engineering",
	"2102": "Mutable Instruments",
	"2109": "Native Instruments",
	"2110": "ROLI Ltd",
	"211A": "IK Multimedia",
	"211C": "Modor Music",
	"211D": "Ableton",
	"2127": "Expert Sleepers",
}
