package id3

// https://id3.org/id3v2.3.0
// https://id3.org/id3v2.4.0-structure
// https://id3.org/id3v2.4.0-frames
// https://id3.org/id3v2-chapters-1.0

import (
	"io"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var imageFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ID3V2,
		Description: "ID3v2 metadata",
		DecodeFn:    id3v2Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Group: &imageFormat},
		},
	})
}

var idDescriptions = decode.StrToScalar{
	"BUF":  {Description: "Recommended buffer size"},
	"CNT":  {Description: "Play counter"},
	"COM":  {Description: "Comments"},
	"CRA":  {Description: "Audio encryption"},
	"CRM":  {Description: "Encrypted meta frame"},
	"EQU":  {Description: "Equalization"},
	"ETC":  {Description: "Event timing codes"},
	"GEO":  {Description: "General encapsulated object"},
	"IPL":  {Description: "Involved people list"},
	"LNK":  {Description: "Linked information"},
	"MCI":  {Description: "Music CD Identifier"},
	"MLL":  {Description: "MPEG location lookup table"},
	"PIC":  {Description: "Attached picture"},
	"POP":  {Description: "Popularimeter"},
	"REV":  {Description: "Reverb"},
	"RVA":  {Description: "Relative volume adjustment"},
	"SLT":  {Description: "Synchronized lyric/text"},
	"STC":  {Description: "Synced tempo codes"},
	"TAL":  {Description: "Album/Movie/Show title"},
	"TBP":  {Description: "BPM (Beats Per Minute)"},
	"TCM":  {Description: "Composer"},
	"TCO":  {Description: "Content type"},
	"TCR":  {Description: "Copyright message"},
	"TDA":  {Description: "Date"},
	"TDY":  {Description: "Playlist delay"},
	"TEN":  {Description: "Encoded by"},
	"TFT":  {Description: "File type"},
	"TIM":  {Description: "Time"},
	"TKE":  {Description: "Initial key"},
	"TLA":  {Description: "Language(s)"},
	"TLE":  {Description: "Length"},
	"TMT":  {Description: "Media type"},
	"TOA":  {Description: "Original artist(s)/performer(s)"},
	"TOF":  {Description: "Original filename"},
	"TOL":  {Description: "Original Lyricist(s)/text writer(s)"},
	"TOR":  {Description: "Original release year"},
	"TOT":  {Description: "Original album/Movie/Show title"},
	"TP1":  {Description: "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group"},
	"TP2":  {Description: "Band/Orchestra/Accompaniment"},
	"TP3":  {Description: "Conductor/Performer refinement"},
	"TP4":  {Description: "Interpreted, remixed, or otherwise modified by"},
	"TPA":  {Description: "Part of a set"},
	"TPB":  {Description: "Publisher"},
	"TRC":  {Description: "ISRC (International Standard Recording Code)"},
	"TRD":  {Description: "Recording dates"},
	"TRK":  {Description: "Track number/Position in set"},
	"TSI":  {Description: "Size"},
	"TSS":  {Description: "Software/hardware and settings used for encoding"},
	"TT1":  {Description: "Content group description"},
	"TT2":  {Description: "Title/Songname/Content description"},
	"TT3":  {Description: "Subtitle/Description refinement"},
	"TXT":  {Description: "Lyricist/text writer"},
	"TXX":  {Description: "User defined text information frame"},
	"TYE":  {Description: "Year"},
	"UFI":  {Description: "Unique file identifier"},
	"ULT":  {Description: "Unsychronized lyric/text transcription"},
	"WAF":  {Description: "Official audio file webpage"},
	"WAR":  {Description: "Official artist/performer webpage"},
	"WAS":  {Description: "Official audio source webpage"},
	"WCM":  {Description: "Commercial information"},
	"WCP":  {Description: "Copyright/Legal information"},
	"WPB":  {Description: "Publishers official webpage"},
	"WXX":  {Description: "User defined URL link frame"},
	"AENC": {Description: "Audio encryption"},
	"APIC": {Description: "Attached picture"},
	"ASPI": {Description: "Audio seek point index"},
	"CHAP": {Description: "Chapter"},
	"COMM": {Description: "Comments"},
	"COMR": {Description: "Commercial frame"},
	"CTOC": {Description: "Table of contents"},
	"ENCR": {Description: "Encryption method registration"},
	"EQU2": {Description: "Equalisation (2)"},
	"EQUA": {Description: "Equalization"},
	"ETCO": {Description: "Event timing codes"},
	"GEOB": {Description: "General encapsulated object"},
	"GRID": {Description: "Group identification registration"},
	"IPLS": {Description: "Involved people list"},
	"LINK": {Description: "Linked information"},
	"MCDI": {Description: "Music CD identifier"},
	"MLLT": {Description: "MPEG location lookup table"},
	"OWNE": {Description: "Ownership frame"},
	"PCNT": {Description: "Play counter"},
	"POPM": {Description: "Popularimeter"},
	"POSS": {Description: "Position synchronisation frame"},
	"PRIV": {Description: "Private frame"},
	"RBUF": {Description: "Recommended buffer size"},
	"RVA2": {Description: "Relative volume adjustment (2)"},
	"RVAD": {Description: "Relative volume adjustment"},
	"RVRB": {Description: "Reverb"},
	"SEEK": {Description: "Seek frame"},
	"SIGN": {Description: "Signature frame"},
	"SYLT": {Description: "Synchronized lyric/text"},
	"SYTC": {Description: "Synchronized tempo codes"},
	"TALB": {Description: "Album/Movie/Show title"},
	"TBPM": {Description: "BPM (beats per minute)"},
	"TCOM": {Description: "Composer"},
	"TCON": {Description: "Content type"},
	"TCOP": {Description: "Copyright message"},
	"TDAT": {Description: "Date"},
	"TDEN": {Description: "Encoding time"},
	"TDLY": {Description: "Playlist delay"},
	"TDOR": {Description: "Original release time"},
	"TDRC": {Description: "Recording time"},
	"TDRL": {Description: "Release time"},
	"TDTG": {Description: "Tagging time"},
	"TENC": {Description: "Encoded by"},
	"TEXT": {Description: "Lyricist/Text writer"},
	"TFLT": {Description: "File type"},
	"TIME": {Description: "Time"},
	"TIPL": {Description: "Involved people list"},
	"TIT1": {Description: "Content group description"},
	"TIT2": {Description: "Title/songname/content description"},
	"TIT3": {Description: "Subtitle/Description refinement"},
	"TKEY": {Description: "Initial key"},
	"TLAN": {Description: "Language(s)"},
	"TLEN": {Description: "Length"},
	"TMCL": {Description: "Musician credits list"},
	"TMED": {Description: "Media type"},
	"TMOO": {Description: "Mood"},
	"TOAL": {Description: "Original album/movie/show title"},
	"TOFN": {Description: "Original filename"},
	"TOLY": {Description: "Original lyricist(s)/text writer(s)"},
	"TOPE": {Description: "Original artist(s)/performer(s)"},
	"TORY": {Description: "Original release year"},
	"TOWN": {Description: "File owner/licensee"},
	"TPE1": {Description: "Lead performer(s)/Soloist(s)"},
	"TPE2": {Description: "Band/orchestra/accompaniment"},
	"TPE3": {Description: "Conductor/performer refinement"},
	"TPE4": {Description: "Interpreted, remixed, or otherwise modified by"},
	"TPOS": {Description: "Part of a set"},
	"TPRO": {Description: "Produced notice"},
	"TPUB": {Description: "Publisher"},
	"TRCK": {Description: "Track number/Position in set"},
	"TRDA": {Description: "Recording dates"},
	"TRSN": {Description: "Internet radio station name"},
	"TRSO": {Description: "Internet radio station owner"},
	"TSIZ": {Description: "Size"},
	"TSOA": {Description: "Album sort order"},
	"TSOP": {Description: "Performer sort order"},
	"TSOT": {Description: "Title sort order"},
	"TSRC": {Description: "ISRC (international standard recording code)"},
	"TSSE": {Description: "Software/Hardware and settings used for encoding"},
	"TSST": {Description: "Set subtitle"},
	"TXXX": {Description: "User defined text information frame"},
	"TYER": {Description: "Year"},
	"UFID": {Description: "Unique file identifier"},
	"USER": {Description: "Terms of use"},
	"USLT": {Description: "Unsychronized lyric/text transcription"},
	"WCOM": {Description: "Commercial information"},
	"WCOP": {Description: "Copyright/Legal information"},
	"WOAF": {Description: "Official audio file webpage"},
	"WOAR": {Description: "Official artist/performer webpage"},
	"WOAS": {Description: "Official audio source webpage"},
	"WORS": {Description: "Official Internet radio station homepage"},
	"WPAY": {Description: "Payment"},
	"WPUB": {Description: "Publishers official webpage"},
	"WXXX": {Description: "User defined URL link frame"},
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
var encodingNames = decode.UToStr{
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

func decodeSyncSafeU32(d *decode.D) uint64 {
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

func textFn(encoding int, nBytes int) func(d *decode.D) string {
	return func(d *decode.D) string {
		return strings.TrimRight(decodeToString(encoding, d.BytesLen(nBytes)), "\x00")
	}
}

func textNullFn(encoding int) func(d *decode.D) string {
	return func(d *decode.D) string {
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
		text := textFn(encoding, int(offsetBytes))(d)

		d.SeekRel(nullLen * 8)
		// seems sometimes utf16 etc has en exta null byte
		if nullLen > 1 && d.PeekBits(8) == 0 {
			d.SeekRel(8)
		}

		return text
	}
}

func decodeFrame(d *decode.D, version int) uint64 {
	var id string
	var size uint64
	var dataSize uint64
	// TODO: global tag unsync?
	unsyncFlag := false

	switch version {
	case 2:
		// Frame ID   "XXX"
		// Frame size $xx xx xx
		id = d.FieldUTF8("id", 3, d.MapStrToScalar(idDescriptions))
		dataSize = d.FieldU24("size")
		size = dataSize + 6
	case 3:
		// Frame ID   $xx xx xx xx  (four characters)
		// Size       $xx xx xx xx
		// Flags      $xx xx
		id = d.FieldUTF8("id", 4, d.MapStrToScalar(idDescriptions))
		dataSize = d.FieldU32("size")

		d.FieldStruct("flags", func(d *decode.D) {
			// %abc00000 %ijk00000
			d.FieldBool("tag_alter_preservation")
			d.FieldBool("file_alter_preservation")
			d.FieldBool("read_only")

			d.FieldU5("unused0")

			d.FieldBool("compression")
			// TODO: read encryption byte, skip decode of frame data?
			d.FieldBool("encryption")
			d.FieldBool("grouping_identity")

			d.FieldU5("unused1")
		})

		size = dataSize + 10
	case 4:
		// Frame ID      $xx xx xx xx  (four characters)
		// Size      4 * %0xxxxxxx  (synchsafe integer)
		// Flags         $xx xx
		id = d.FieldUTF8("id", 4, d.MapStrToScalar(idDescriptions))
		dataSize = d.FieldUFn("size", decodeSyncSafeU32)
		var headerLen uint64 = 10

		dataLenFlag := false
		d.FieldStruct("flags", func(d *decode.D) {
			// %0abc0000 %0h00kmnp
			d.FieldU1("unused0")
			d.FieldBool("tag_alter_preservation")
			d.FieldBool("file_alter_preservation")
			d.FieldBool("read_only")

			d.FieldU5("unused1")

			d.FieldBool("grouping_identity")

			d.FieldU2("unused2")

			d.FieldBool("compression")
			// TODO: read encryption byte, skip decode of frame data?
			d.FieldBool("encryption")
			unsyncFlag = d.FieldBool("unsync")
			dataLenFlag = d.FieldBool("data_length_indicator")
		})

		if dataLenFlag {
			d.FieldUFn("data_length_indicator", decodeSyncSafeU32)
			dataSize -= 4
			headerLen += 4
		}

		size = dataSize + headerLen
	default:
		// can't know size
		d.Fatalf("unknown version")
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
			d.FieldStrFn("element_id", textNullFn(encodingUTF8))
			d.FieldU32("start_time")
			d.FieldU32("end_time")
			d.FieldU32("start_offset")
			d.FieldU32("end_offset")
			decodeFrames(d, version, uint64(d.BitsLeft()/8))
		},
		"CTOC": func(d *decode.D) {
			d.FieldStrFn("element_id", textNullFn(encodingUTF8))
			d.FieldU8("ctoc_flags")
			entryCount := d.FieldU8("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStrFn("entry", textNullFn(encodingUTF8))
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
			encoding := d.FieldU8("text_encoding", d.MapUToStrSym(encodingNames))
			d.FieldStrFn("mime_type", textNullFn(encodingUTF8))
			d.FieldU8("picture_type") // TODO: table
			d.FieldStrFn("description", textNullFn(int(encoding)))
			dv, _, _ := d.TryFieldFormatLen("picture", d.BitsLeft(), imageFormat, nil)
			if dv == nil {
				d.FieldRawLen("picture", d.BitsLeft())
			}
		},

		// <Header for 'General encapsulated object', ID: "GEOB">
		// Text encoding          $xx
		// MIME type              <text string> $00
		// Filename               <text string according to encoding> $00 (00)
		// Content description    <text string according to encoding> $00 (00)
		// Encapsulated object    <binary data>
		"GEOB": func(d *decode.D) {
			encoding := d.FieldU8("text_encoding", d.MapUToStrSym(encodingNames))
			d.FieldStrFn("mime_type", textNullFn(encodingUTF8))
			d.FieldStrFn("filename", textNullFn(int(encoding)))
			d.FieldStrFn("description", textNullFn(int(encoding)))
			dv, _, _ := d.TryFieldFormatLen("data", d.BitsLeft(), imageFormat, nil)
			if dv == nil {
				d.FieldRawLen("data", d.BitsLeft())
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
			encoding := d.FieldU8("text_encoding", d.MapUToStrSym(encodingNames))
			d.FieldUTF8("language", 3)
			d.FieldStrFn("description", textNullFn(int(encoding)))
			d.FieldStrFn("value", textFn(int(encoding), int(d.BitsLeft()/8)))
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
			encoding := d.FieldU8("text_encoding", d.MapUToStrSym(encodingNames))
			d.FieldStrFn("text", textFn(int(encoding), int(d.BitsLeft()/8)))
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
			encoding := d.FieldU8("text_encoding", d.MapUToStrSym(encodingNames))
			d.FieldStrFn("description", textNullFn(int(encoding)))
			d.FieldStrFn("value", textFn(int(encoding), int(d.BitsLeft()/8)))
		},
		// <Header for 'Private frame', ID: "PRIV">
		// Owner identifier      <text string> $00
		// The private data      <binary data>
		"PRIV": func(d *decode.D) {
			// TODO: is default ISO8859-1?
			d.FieldStrFn("owner", textNullFn(int(encodingISO8859_1)))
			d.FieldRawLen("data", d.BitsLeft())
		},
	}

	idNormalized := id
	switch {
	case id == "COMM", id == "COM", id == "USLT", id == "ULT":
		idNormalized = "COMM"
	case id == "TXX", id == "TXXX":
		idNormalized = "TXXX"
	case len(id) > 0 && id[0] == 'T':
		idNormalized = "T000"
	}

	if unsyncFlag {
		// TODO: DecodeFn
		// TODO: unknown after frame decode
		unsyncedBb := d.MustNewBitBufFromReader(unsyncReader{Reader: d.BitBufRange(d.Pos(), int64(dataSize)*8)})
		d.FieldFormatBitBuf("unsync", unsyncedBb, decode.FormatFn(func(d *decode.D, in interface{}) interface{} {
			if fn, ok := frames[idNormalized]; ok {
				fn(d)
			} else {
				d.FieldRawLen("data", d.BitsLeft())
			}

			return nil
		}), nil)
		d.FieldRawLen("data", int64(dataSize*8))
	} else {
		if fn, ok := frames[idNormalized]; ok {
			d.LenFn(int64(dataSize)*8, func(d *decode.D) {
				fn(d)
			})
		} else {
			d.FieldRawLen("data", int64(dataSize*8))
		}
	}

	return size
}

func decodeFrames(d *decode.D, version int, size uint64) {
	d.FieldArray("frames", func(d *decode.D) {
		for size > 0 {
			if d.PeekBits(8) == 0 {
				return
			}

			d.FieldStruct("frame", func(d *decode.D) {
				size -= decodeFrame(d, version)
			})
		}
	})

	if size > 0 {
		d.FieldRawLen("padding", int64(size*8), d.BitBufIsZero)
	}
}

func id3v2Decode(d *decode.D, in interface{}) interface{} {
	d.AssertAtLeastBitsLeft(4 * 8)
	d.FieldUTF8("magic", 3, d.ValidateStr("ID3"))
	version := int(d.FieldU8("version"))
	versionValid := version == 2 || version == 3 || version == 4
	if !versionValid {
		d.Fatalf("unsupported version %d", version)
	}

	d.FieldU8("revision")
	var extendedHeader bool
	d.FieldStruct("flags", func(d *decode.D) {
		d.FieldBool("unsynchronisation")
		extendedHeader = d.FieldBool("extended_header")
		d.FieldBool("experimental_indicator")
		d.FieldU5("unused")
	})
	size := d.FieldUFn("size", decodeSyncSafeU32)

	var extHeaderSize uint64
	if extendedHeader {
		d.FieldStruct("extended_header", func(d *decode.D) {
			switch version {
			case 3:
				extHeaderSize = d.FieldU32("size")
				d.FieldRawLen("data", int64(extHeaderSize)*8)
			case 4:
				extHeaderSize = d.FieldUFn("size", decodeSyncSafeU32)
				// in v4 synchsafe integer includes itself
				d.FieldRawLen("data", (int64(extHeaderSize)-4)*8)
			}
		})
	}

	decodeFrames(d, version, size)

	return nil
}
