package id3

// https://id3.org/id3v2.3.0
// https://id3.org/id3v2.4.0-structure
// https://id3.org/id3v2.4.0-frames
// https://id3.org/id3v2-chapters-1.0

import (
	"fmt"
	"io"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var imageFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.ID3V2,
		Description: "ID3v2 metadata",
		DecodeFn:    id3v2Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Formats: &imageFormat},
		},
	})
}

var idDescriptions = map[string]string{
	"CHAP": "Chapter",
	"CTOC": "Table of contents",

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

// id3v2 MPEG/AAC unsynchronisation reader
// Replace 0xff 0x00 0xab with 0xff 0xab in byte stream
type unsyncReader struct {
	io.Reader
	lastFF bool
}

func (r unsyncReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	ni := 0
	for i, b := range p[0:n] {
		if r.lastFF && b == 0x00 {
			n--
			r.lastFF = false
			continue
		} else {
			r.lastFF = b == 0xff
		}
		p[ni] = p[i]
		ni++
	}

	return n, err
}

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

var encodingLen = map[uint64]int64{
	encodingISO8859_1: 1,
	encodingUTF16:     2,
	encodingUTF16BE:   2,
	encodingUTF8:      1,
}

func decodeToString(e int, b []byte) string {
	var enc encoding.Encoding

	switch e {
	case encodingISO8859_1:
		enc = charmap.ISO8859_1
	case encodingUTF16:
		enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	case encodingUTF16BE:
		enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	default:
		enc = unicode.UTF8
	}

	// TODO: try decode?
	s, _ := enc.NewDecoder().String(string(b))
	return s
}

func syncSafeU32(d *decode.D) uint64 {
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

func text(d *decode.D, encoding int, nBytes int) string {
	return strings.TrimRight(decodeToString(encoding, d.BytesLen(nBytes)), "\x00")
}

func textNull(d *decode.D, encoding int) string {
	nullLen := encodingLen[encodingUTF8]
	if n, ok := encodingLen[uint64(encoding)]; ok {
		nullLen = n
	}

	offset, _ := d.PeekFind(
		int(nullLen)*8,
		nullLen*8,
		func(v uint64) bool { return v == 0 },
		-1,
	)
	offsetBytes := offset / 8
	text := text(d, encoding, int(offsetBytes))

	d.SeekRel(nullLen * 8)
	// seems sometimes utf16 etc has en exta null byte
	if nullLen > 1 && d.PeekBits(8) == 0 {
		d.SeekRel(8)
	}

	return text
}

func fieldSyncSafeU32(d *decode.D, name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return syncSafeU32(d), decode.NumberDecimal, ""
	})
}

func fieldTextNull(d *decode.D, name string, encoding int) string { //nolint:unparam
	return d.FieldStrFn(name, func() (string, string) {
		return textNull(d, encoding), ""
	})
}

func fieldText(d *decode.D, name string, encoding int, nBytes int) string { //nolint:unparam
	return d.FieldStrFn(name, func() (string, string) {
		return text(d, encoding, nBytes), ""
	})
}

func decodeFrame(d *decode.D, version int) uint64 {
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
	// TODO: global tag unsync?
	unsyncFlag := false

	idDescription := ""
	if d, ok := idDescriptions[id]; ok {
		idDescription = d
	}

	switch version {
	case 2:
		// Frame ID   "XXX"
		// Frame size $xx xx xx
		d.FieldStrFn("id", func() (string, string) { return d.UTF8(3), idDescription })
		dataSize = d.FieldU24("size")
		size = dataSize + 6
	case 3:
		// Frame ID   $xx xx xx xx  (four characters)
		// Size       $xx xx xx xx
		// Flags      $xx xx
		d.FieldStrFn("id", func() (string, string) { return d.UTF8(4), idDescription })
		dataSize = d.FieldU32("size")

		d.FieldStructFn("flags", func(d *decode.D) {
			// %abc00000 %ijk00000
			d.FieldBool("tag_alter_preservation")
			d.FieldBool("file_alter_preservation")
			d.FieldBool("read_only")

			d.FieldU5("unused0")

			d.FieldBool("compression")
			// TODO: read encruption byte, skip decode of frame data?
			d.FieldBool("encryption")
			d.FieldBool("grouping_identity")

			d.FieldU5("unused1")
		})

		size = dataSize + 10
	case 4:
		// Frame ID      $xx xx xx xx  (four characters)
		// Size      4 * %0xxxxxxx  (synchsafe integer)
		// Flags         $xx xx
		d.FieldStrFn("id", func() (string, string) { return d.UTF8(4), idDescription })
		dataSize = fieldSyncSafeU32(d, "size")
		var headerLen uint64 = 10

		dataLenFlag := false
		d.FieldStructFn("flags", func(d *decode.D) {
			// %0abc0000 %0h00kmnp
			d.FieldU1("unused0")
			d.FieldBool("tag_alter_preservation")
			d.FieldBool("file_alter_preservation")
			d.FieldBool("read_only")

			d.FieldU5("unused1")

			d.FieldBool("grouping_identity")

			d.FieldU2("unused2")

			d.FieldBool("compression")
			// TODO: read encruption byte, skip decode of frame data?
			d.FieldBool("encryption")
			unsyncFlag = d.FieldBool("unsync")
			dataLenFlag = d.FieldBool("data_length_indicator")
		})

		if dataLenFlag {
			fieldSyncSafeU32(d, "data_length_indicator")
			dataSize -= 4
			headerLen += 4
		}

		size = dataSize + headerLen
	}

	// note frame function run inside a SubLenFn so they can use BitLefts and
	// can't accidentally read too far
	frames := map[string]func(d *decode.D){
		// <ID3v2.3 or ID3v2.4 frame header, ID: "CHAP">           (10 bytes)
		// Element ID      <text string> $00
		// Start time      $xx xx xx xx
		// End time        $xx xx xx xx
		// Start offset    $xx xx xx xx
		// End offset      $xx xx xx xx
		// <Optional embedded sub-frames>
		"CHAP": func(d *decode.D) {
			fieldTextNull(d, "element_id", encodingUTF8)
			d.FieldU32("start_time")
			d.FieldU32("end_time")
			d.FieldU32("start_offset")
			d.FieldU32("end_offset")
			decodeFrames(d, version, uint64(d.BitsLeft()/8))
		},
		"CTOC": func(d *decode.D) {
			fieldTextNull(d, "element_id", encodingUTF8)
			d.FieldU8("ctoc_flags")
			entryCount := d.FieldU8("entry_count")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					fieldTextNull(d, "entry", encodingUTF8)
				}
			})
		},

		// <Header for 'Attached picture', ID: "APIC">
		// Text encoding      $xx
		// MIME type          <text string> $00
		// Picture type       $xx
		// Description        <text string according to encoding> $00 (00)
		// Picture data       <binary data>
		"APIC": func(d *decode.D) {
			encoding, _ := d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8, decode.NumberDecimal)
			fieldTextNull(d, "mime_type", encodingUTF8)
			d.FieldU8("picture_type") // TODO: table
			fieldTextNull(d, "description", int(encoding))
			dv, _, _ := d.FieldTryFormatLen("picture", d.BitsLeft(), imageFormat)
			if dv == nil {
				d.FieldBitBufLen("picture", d.BitsLeft())
			}
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
		"COMM": func(d *decode.D) {
			encoding, _ := d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8, decode.NumberDecimal)
			d.FieldUTF8("language", 3)
			fieldTextNull(d, "description", int(encoding))
			fieldText(d, "value", int(encoding), int(d.BitsLeft()/8))
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
		"T000": func(d *decode.D) {
			encoding, _ := d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8, decode.NumberDecimal)
			fieldText(d, "text", int(encoding), int(d.BitsLeft()/8))
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
		"TXXX": func(d *decode.D) {
			encoding, _ := d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8, decode.NumberDecimal)
			fieldTextNull(d, "description", int(encoding))
			fieldText(d, "value", int(encoding), int(d.BitsLeft()/8))
		},
		// <Header for 'Private frame', ID: "PRIV">
		// Owner identifier      <text string> $00
		// The private data      <binary data>
		"PRIV": func(d *decode.D) {
			// TODO: is default ISO8859-1?
			fieldTextNull(d, "owner", int(encodingISO8859_1))
			d.FieldBitBufLen("data", d.BitsLeft())
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

	if unsyncFlag {
		// TODO: DecodeFn
		// TODO: unknown after frame decode
		unsyncedBb := decode.MustNewBitBufFromReader(unsyncReader{Reader: d.BitBufRange(d.Pos(), int64(dataSize)*8)})
		d.FieldFormatBitBuf("unsync", unsyncedBb, decode.FormatFn(func(d *decode.D, in interface{}) interface{} {
			if fn, ok := frames[idNormalized]; ok {
				fn(d)
			} else {
				d.FieldBitBufLen("data", d.BitsLeft())
			}

			return nil
		}))
		d.FieldBitBufLen("data", int64(dataSize*8))
	} else {
		if fn, ok := frames[idNormalized]; ok {
			d.DecodeLenFn(int64(dataSize)*8, func(d *decode.D) {
				fn(d)
			})
		} else {
			d.FieldBitBufLen("data", int64(dataSize*8))
		}
	}

	return size
}

func decodeFrames(d *decode.D, version int, size uint64) {
	d.FieldArrayFn("frames", func(d *decode.D) {
		for size > 0 {
			if d.PeekBits(8) == 0 {
				return
			}

			d.FieldStructFn("frame", func(d *decode.D) {
				size -= decodeFrame(d, version)
			})
		}

	})

	if size > 0 {
		d.FieldValidateZeroPadding("padding", int(size)*8)
	}
}

func id3v2Decode(d *decode.D, in interface{}) interface{} {
	d.ValidateAtLeastBitsLeft(4 * 8)
	d.FieldValidateUTF8("magic", "ID3")
	version := int(d.FieldU8("version"))
	versionValid := version == 2 || version == 3 || version == 4
	if !versionValid {
		d.Invalid(fmt.Sprintf("unsupported version %d", version))
	}

	d.FieldU8("revision")
	var extendedHeader bool
	d.FieldStructFn("flags", func(d *decode.D) {
		d.FieldBool("unsynchronisation")
		extendedHeader = d.FieldBool("extended_header")
		d.FieldBool("experimental_indicator")
		d.FieldU5("unused")
	})
	size := d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) {
		return syncSafeU32(d), decode.NumberDecimal, ""
	})

	var extHeaderSize uint64
	if extendedHeader {
		d.FieldStructFn("extended_header", func(d *decode.D) {
			switch version {
			case 3:
				extHeaderSize = d.FieldU32("size")
				d.FieldBitBufLen("data", int64(extHeaderSize)*8)
			case 4:
				extHeaderSize = d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) {
					return syncSafeU32(d), decode.NumberDecimal, ""
				})
				// in v4 synchsafe integer includes itself
				d.FieldBitBufLen("data", (int64(extHeaderSize)-4)*8)
			}
		})
	}

	decodeFrames(d, version, size)

	return nil
}
