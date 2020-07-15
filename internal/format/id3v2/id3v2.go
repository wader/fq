package id3v2

// https://id3.org/id3v2.3.0
// https://id3.org/id3v2.4.0-structure
// https://id3.org/id3v2.4.0-frames

import (
	"bytes"
	"fmt"
	"fq/internal/decode"
	"fq/internal/format/group"
	"strings"
)

var Tag = &decode.Format{
	Name: "id3v2",
	New:  func() decode.Decoder { return &TagDecoder{} },
}

var idDesriptions = map[string]string{
	"AENC": "Audio encryption",
	"APIC": "Attached picture",
	"ASPI": "Audio seek point index",
	"COMM": "Comments",
	"COMR": "Commercial frame",
	"ENCR": "Encryption method registration",
	"EQU2": "Equalisation (2)",
	"EQUA": "Equalization",
	"ETCO": "Event timing codes",
	"GEOB": "General encapsulated object",
	"GRID": "Group identification registration",
	"IPLS": "Involved people list",
	"LINK": "Linked information",
	"MCDI": "Music CD identifier",
	"MLLT": "MPEG location lookup table",
	"OWNE": "Ownership frame",
	"PCNT": "Play counter",
	"POPM": "Popularimeter",
	"POSS": "Position synchronisation frame",
	"PRIV": "Private frame",
	"RBUF": "Recommended buffer size",
	"RVA2": "Relative volume adjustment (2)",
	"RVAD": "Relative volume adjustment",
	"RVRB": "Reverb",
	"SEEK": "Seek frame",
	"SIGN": "Signature frame",
	"SYLT": "Synchronized lyric/text",
	"SYTC": "Synchronized tempo codes",
	"TALB": "Album/Movie/Show title",
	"TBPM": "BPM (beats per minute)",
	"TCOM": "Composer",
	"TCON": "Content type",
	"TCOP": "Copyright message",
	"TDAT": "Date",
	"TDEN": "Encoding time",
	"TDLY": "Playlist delay",
	"TDOR": "Original release time",
	"TDRC": "Recording time",
	"TDRL": "Release time",
	"TDTG": "Tagging time",
	"TENC": "Encoded by",
	"TEXT": "Lyricist/Text writer",
	"TFLT": "File type",
	"TIME": "Time",
	"TIPL": "Involved people list",
	"TIT1": "Content group description",
	"TIT2": "Title/songname/content description",
	"TIT3": "Subtitle/Description refinement",
	"TKEY": "Initial key",
	"TLAN": "Language(s)",
	"TLEN": "Length",
	"TMCL": "Musician credits list",
	"TMED": "Media type",
	"TMOO": "Mood",
	"TOAL": "Original album/movie/show title",
	"TOFN": "Original filename",
	"TOLY": "Original lyricist(s)/text writer(s)",
	"TOPE": "Original artist(s)/performer(s)",
	"TORY": "Original release year",
	"TOWN": "File owner/licensee",
	"TPE1": "Lead performer(s)/Soloist(s)",
	"TPE2": "Band/orchestra/accompaniment",
	"TPE3": "Conductor/performer refinement",
	"TPE4": "Interpreted, remixed, or otherwise modified by",
	"TPOS": "Part of a set",
	"TPRO": "Produced notice",
	"TPUB": "Publisher",
	"TRCK": "Track number/Position in set",
	"TRDA": "Recording dates",
	"TRSN": "Internet radio station name",
	"TRSO": "Internet radio station owner",
	"TSIZ": "Size",
	"TSOA": "Album sort order",
	"TSOP": "Performer sort order",
	"TSOT": "Title sort order",
	"TSRC": "ISRC (international standard recording code)",
	"TSSE": "Software/Hardware and settings used for encoding",
	"TSST": "Set subtitle",
	"TXXX": "User defined text information frame",
	"TYER": "Year",
	"UFID": "Unique file identifier",
	"USER": "Terms of use",
	"USLT": "Unsychronized lyric/text transcription",
	"WCOM": "Commercial information",
	"WCOP": "Copyright/Legal information",
	"WOAF": "Official audio file webpage",
	"WOAR": "Official artist/performer webpage",
	"WOAS": "Official audio source webpage",
	"WORS": "Official Internet radio station homepage",
	"WPAY": "Payment",
	"WPUB": "Publishers official webpage",
	"WXXX": "User defined URL link frame",
	"BUF":  "Recommended buffer size",
	"CNT":  "Play counter",
	"COM":  "Comments",
	"CRA":  "Audio encryption",
	"CRM":  "Encrypted meta frame",
	"ETC":  "Event timing codes",
	"EQU":  "Equalization",
	"GEO":  "General encapsulated object",
	"IPL":  "Involved people list",
	"LNK":  "Linked information",
	"MCI":  "Music CD Identifier",
	"MLL":  "MPEG location lookup table",
	"PIC":  "Attached picture",
	"POP":  "Popularimeter",
	"REV":  "Reverb",
	"RVA":  "Relative volume adjustment",
	"SLT":  "Synchronized lyric/text",
	"STC":  "Synced tempo codes",
	"TAL":  "Album/Movie/Show title",
	"TBP":  "BPM (Beats Per Minute)",
	"TCM":  "Composer",
	"TCO":  "Content type",
	"TCR":  "Copyright message",
	"TDA":  "Date",
	"TDY":  "Playlist delay",
	"TEN":  "Encoded by",
	"TFT":  "File type",
	"TIM":  "Time",
	"TKE":  "Initial key",
	"TLA":  "Language(s)",
	"TLE":  "Length",
	"TMT":  "Media type",
	"TOA":  "Original artist(s)/performer(s)",
	"TOF":  "Original filename",
	"TOL":  "Original Lyricist(s)/text writer(s)",
	"TOR":  "Original release year",
	"TOT":  "Original album/Movie/Show title",
	"TP1":  "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group",
	"TP2":  "Band/Orchestra/Accompaniment",
	"TP3":  "Conductor/Performer refinement",
	"TP4":  "Interpreted, remixed, or otherwise modified by",
	"TPA":  "Part of a set",
	"TPB":  "Publisher",
	"TRC":  "ISRC (International Standard Recording Code)",
	"TRD":  "Recording dates",
	"TRK":  "Track number/Position in set",
	"TSI":  "Size",
	"TSS":  "Software/hardware and settings used for encoding",
	"TT1":  "Content group description",
	"TT2":  "Title/Songname/Content description",
	"TT3":  "Subtitle/Description refinement",
	"TXT":  "Lyricist/text writer",
	"TXX":  "User defined text information frame",
	"TYE":  "Year",
	"UFI":  "Unique file identifier",
	"ULT":  "Unsychronized lyric/text transcription",
	"WAF":  "Official audio file webpage",
	"WAR":  "Official artist/performer webpage",
	"WAS":  "Official audio source webpage",
	"WCM":  "Commercial information",
	"WCP":  "Copyright/Legal information",
	"WPB":  "Publishers official webpage",
	"WXX":  "User defined URL link frame",
}

type encoding int

const (
	encodingISO8859_1 = 0
	encodingUTF16     = 1
	encodingUTF16BE   = 2
	encodingUTF8      = 3
)

// $00 ISO-8859-1 [ISO-8859-1]. Terminated with $00.
// $01 UTF-16 [UTF-16] encoded Unicode [UNICODE] with BOM. All
//     strings in the same frame SHALL have the same byteorder.
//     Terminated with $00 00.
// $02 UTF-16BE [UTF-16] encoded Unicode [UNICODE] without BOM.
//     Terminated with $00 00.
// $03 UTF-8 [UTF-8] encoded Unicode [UNICODE]. Terminated with $00.
var encodingNames = map[uint64]string{
	encodingISO8859_1: "ISO-8859-1",
	encodingUTF16:     "UTF-16",
	encodingUTF16BE:   "UTF-16BE",
	encodingUTF8:      "UTF-8",
}

var encodingLen = map[uint64]int{
	encodingISO8859_1: 1,
	encodingUTF16:     2,
	encodingUTF16BE:   2,
	encodingUTF8:      1,
}

var encodingToUTF8 = map[int]func(b []byte) string{
	encodingISO8859_1: func(b []byte) string {
		rs := make([]rune, len(b))
		for i, r := range b {
			rs[i] = rune(r)
		}
		return string(rs)
	},
	encodingUTF16: func(b []byte) string {
		beBOM := []byte("\xfe\xff")
		leBOM := []byte("\xff\xfe")
		var rs []rune
		switch {
		case bytes.HasPrefix(b, leBOM):
			// strip BOM
			b = b[2:]
			rs = make([]rune, len(b)/2)
			for i := 0; i < len(b)/2; i++ {
				rs[i] = rune(uint(b[i*2]) | uint(b[i*2+1])<<8)
			}
		case bytes.HasPrefix(b, beBOM):
			b = b[2:]
			fallthrough
		default:
			rs = make([]rune, len(b)/2)
			for i := 0; i < len(b)/2; i++ {
				rs[i] = rune(uint(b[i*2])<<8 | uint(b[i*2+1]))
			}
		}
		return string(rs)
	},
	encodingUTF16BE: func(b []byte) string {
		rs := make([]rune, len(b)/2)
		for i := 0; i < len(b)/2; i++ {
			rs[i] = rune(uint(b[i*2])<<8 + uint(b[i*2+1]))
		}
		return string(rs)
	},
	encodingUTF8: func(b []byte) string {
		return string(b)
	},
}

// Decoder is ID3v2 tag decoder
type TagDecoder struct {
	decode.Common
}

func (d *TagDecoder) SyncSafeU32() uint64 {
	u := d.U32()
	// syncsafe integer is a number encoded
	// with 8th bit in each byte set to zero
	// 0aaaaaaa0bbbbbbb0ccccccc0ddddddd ->
	// 0000aaaaaaabbbbbbbcccccccddddddd
	return (((u & 0x7f000000) >> 3) |
		((u & 0x007f0000) >> 2) |
		((u & 0x00007f00) >> 1) |
		((u & 0x0000007f) >> 0))
}

func (d *TagDecoder) Text(encoding int, nBytes uint64) string {
	encodingFn := encodingToUTF8[encodingUTF8]
	if fn, ok := encodingToUTF8[encoding]; ok {
		encodingFn = fn
	}
	return strings.TrimRight(encodingFn(d.BytesLen(nBytes)), "\x00")
}

func (d *TagDecoder) TextNull(encoding int) string {
	nullLen := encodingLen[encodingUTF8]
	if n, ok := encodingLen[uint64(encoding)]; ok {
		nullLen = n
	}

	textLen := d.PeekFind(uint64(nullLen*8), 0, -1)/8 - uint64(nullLen)
	text := d.Text(encoding, textLen)
	// TODO: field?
	d.SeekRel(int64(nullLen) * 8)

	return text
}

func (d *TagDecoder) FieldSyncSafeU32(name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.NumberFormat, string) {
		return d.SyncSafeU32(), decode.NumberDecimal, ""
	})
}

func (d *TagDecoder) FieldTextNull(name string, encoding int) string {
	return d.FieldStrFn(name, func() (string, string) {
		return d.TextNull(encoding), ""
	})
}

func (d *TagDecoder) FieldText(name string, encoding int, nBytes uint64) string {
	return d.FieldStrFn(name, func() (string, string) {
		return d.Text(encoding, nBytes), ""
	})
}

func (d *TagDecoder) DecodeFrame(version int) uint64 {
	var id string

	switch version {
	case 2:
		id = d.UTF8(3)
		d.SeekRel(-3 * 8)
	case 3, 4:
		id = d.UTF8(4)
		d.SeekRel(-4 * 8)
	}

	var size uint64
	var dataSize uint64

	d.FieldStrFn(id, func() (string, string) {
		switch version {
		case 2:
			// Frame ID   "XXX"
			// Frame size $xx xx xx
			d.FieldUTF8("id", 3)
			dataSize = d.FieldU24("size")
			size = dataSize + 6
		case 3:
			// Frame ID   $xx xx xx xx  (four characters)
			// Size       $xx xx xx xx
			// Flags      $xx xx
			d.FieldUTF8("id", 4)
			dataSize = d.FieldU32("size")
			d.FieldU16("flags")
			size = dataSize + 10
		case 4:
			// Frame ID      $xx xx xx xx  (four characters)
			// Size      4 * %0xxxxxxx  (synchsafe integer)
			// Flags         $xx xx
			d.FieldUTF8("id", 4)
			dataSize = d.FieldSyncSafeU32("size")
			var headerLen uint64 = 10

			const flagUnsync = 0b10
			const flagDataLen = 0b1

			dataLenFlag := false
			d.FieldNoneFn("flags", func() {
				d.FieldU14("unused")
				d.FieldBool("unsync")
				dataLenFlag = d.FieldBool("data_length_indicator")
			})

			if dataLenFlag {
				d.FieldSyncSafeU32("data_length_indicator")
				dataSize -= 4
				headerLen = 4
			}

			size = dataSize + headerLen
		}

		// note frame function run inside a SubLenFn so they can use BitLefts and
		// can't accidentally read too far
		frames := map[string]func(){
			// <Header for 'Attached picture', ID: "APIC">
			// Text encoding      $xx
			// MIME type          <text string> $00
			// Picture type       $xx
			// Description        <text string according to encoding> $00 (00)
			// Picture data       <binary data>
			"APIC": func() {
				encoding := int(d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8))
				d.FieldTextNull("mime_type", encodingUTF8)
				d.FieldU8("picture_type") // TODO: table
				d.FieldTextNull("description", encoding)
				d.FieldDecodeLen("picture", d.BitsLeft(), group.Images...)
			},
			// Unsynced lyrics/text "ULT"
			// Frame size           $xx xx xx
			// Text encoding        $xx
			// Language             $xx xx xx
			// Content descriptor   <textstring> $00 (00)
			// Lyrics/text          <textstring>
			//
			// <Header for 'Unsynchronised lyrics/text transcription', ID: "USLT">
			// Text encoding        $xx
			// Language             $xx xx xx
			// Content descriptor   <text string according to encoding> $00 (00)
			// Lyrics/text          <full text string according to encoding>
			//
			// Comment                   "COM"
			// Frame size                $xx xx xx
			// Text encoding             $xx
			// Language                  $xx xx xx
			// Short content description <textstring> $00 (00)
			// The actual text           <textstring>
			//
			// <Header for 'Comment', ID: "COMM">
			// Text encoding          $xx
			// Language               $xx xx xx
			// Short content descrip. <text string according to encoding> $00 (00)
			// The actual text        <full text string according to encoding>
			"COMM": func() {
				encoding := int(d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8))
				d.FieldUTF8("language", 3)
				d.FieldTextNull("description", encoding)
				d.FieldText("value", encoding, d.BitsLeft()/8)
			},
			// Text information identifier  "T00" - "TZZ" , excluding "TXX",
			//                             described in 4.2.2.
			// Frame size                   $xx xx xx
			// Text encoding                $xx
			// Information                  <textstring>
			//
			// <Header for 'Text information frame', ID: "T000" - "TZZZ",
			// excluding "TXXX" described in 4.2.6.>
			// Text encoding                $xx
			// Information                  <text string(s) according to encoding>
			"T000": func() {
				encoding := int(d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8))
				d.FieldText("text", encoding, d.BitsLeft()/8)
			},
			// User defined...   "TXX"
			// Frame size        $xx xx xx
			// Text encoding     $xx
			// Description       <textstring> $00 (00)
			// Value             <textstring>
			//
			// <Header for 'User defined text information frame', ID: "TXXX">
			// Text encoding     $xx
			// Description       <text string according to encoding> $00 (00)
			// Value             <text string according to encoding>
			"TXXX": func() {
				encoding := int(d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8))
				d.FieldTextNull("description", encoding)
				d.FieldText("value", encoding, d.BitsLeft()/8)
			},
		}

		idNormalized := id
		switch {
		case id == "COMM", id == "COM", id == "USLT", id == "ULT":
			idNormalized = "COMM"
		case id == "TXX", id == "TXXX":
			idNormalized = "TXXX"
		case id[0] == 'T':
			idNormalized = "T000"
		}

		if fn, ok := frames[idNormalized]; ok {
			d.SubLenFn(dataSize*8, fn)
		} else {
			d.FieldBytesLen("data", dataSize)
		}

		idDescription := ""
		if d, ok := idDesriptions[id]; ok {
			idDescription = d
		}

		return id, idDescription
	})

	// TODO
	return size
}

func (d *TagDecoder) DecodeFrames(version int, size uint64) {
	for size > 0 {
		for d.PeekBits(8) == 0 {
			d.FieldValidateZeroPadding("padding", size*8)
			return
		}

		size -= d.DecodeFrame(version)
	}

	// TODO: padding?
}

// Decode ID3v2
func (d *TagDecoder) Decode() {
	d.ValidateAtLeastBitsLeft(4 * 8)
	d.FieldValidateString("magic", "ID3")
	version := int(d.FieldU8("version"))
	versionValid := version == 2 || version == 3 || version == 4
	if !versionValid {
		d.Invalid(fmt.Sprintf("unsupported version %d", version))
	}

	d.FieldU8("revision")
	var extendedHeader bool
	d.FieldNoneFn("flags", func() {
		d.FieldU1("unsynchronisation")
		extendedHeader = d.FieldBool("extended_header")
		d.FieldU1("experimental_indicator")
		d.FieldU5("unused")
	})
	size := d.FieldUFn("size", func() (uint64, decode.NumberFormat, string) {
		return d.SyncSafeU32(), decode.NumberDecimal, ""
	})

	var extHeaderSize uint64
	if extendedHeader {
		d.FieldNoneFn("extended_header", func() {
			switch version {
			case 3:
				extHeaderSize = d.FieldU32("size")
				d.FieldBytesLen("data", extHeaderSize)
			case 4:
				extHeaderSize = d.FieldUFn("size", func() (uint64, decode.NumberFormat, string) {
					return d.SyncSafeU32(), decode.NumberDecimal, ""
				})
				// in v4 synchsafe integer includes itself
				d.FieldBytesLen("data", extHeaderSize-4)
			}
		})
	}

	d.DecodeFrames(version, size)
}
