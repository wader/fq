package mkv

// https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
// https://matroska.org/technical/specs/index.html
// https://www.matroska.org/technical/basics.html
// https://www.matroska.org/technical/codec_specs.html

// TODO: rename simepleblock/block to just defer decode etc?
// TODO: CRC

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/ranges"
)

var vorbisPacketFormat []*decode.Format
var vp9FrameFormat []*decode.Format
var aacFrameFormat []*decode.Format
var mpegASCFrameFormat []*decode.Format
var mpegSPUFrameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MKV,
		Description: "Matroska (EBML)",
		Groups:      []string{format.PROBE},
		DecodeFn:    mkvDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.MPEG_ASC}, Formats: &mpegASCFrameFormat},
			{Names: []string{format.MPEG_SPU}, Formats: &mpegSPUFrameFormat},
		},
	})
}

// TODO: smarter?
func decodeRawVintWidth(d *decode.D) (uint64, int) {
	n := d.U8()
	w := 1
	for i := 0; (n & (1 << (7 - i))) == 0; i++ {
		w++
	}
	for i := 1; i < w; i++ {
		n = n<<8 | d.U8()
	}
	return n, w
}

func decodeRawVint(d *decode.D) uint64 {
	n, _ := decodeRawVintWidth(d)
	return n
}

func decodeVint(d *decode.D) uint64 {
	n, w := decodeRawVintWidth(d)
	m := (uint64(1<<((w-1)*8+(8-w))) - 1)
	return n & m
}

func fieldDecodeRawVint(d *decode.D, name string, displayFormat decode.DisplayFormat) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return decodeRawVint(d), displayFormat, ""
	})
}
func fieldDecodeVint(d *decode.D, name string, displayFormat decode.DisplayFormat) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return decodeVint(d), displayFormat, ""
	})
}

type ebmlType int

const (
	ebmlInteger ebmlType = iota
	ebmlUinteger
	ebmlFloat
	ebmlString
	ebmlUTF8
	ebmlDate
	ebmlBinary
	ebmlMaster
)

var ebmlTypeNames = map[ebmlType]string{
	ebmlInteger:  "integer",
	ebmlUinteger: "uinteger",
	ebmlFloat:    "float",
	ebmlString:   "string",
	ebmlUTF8:     "UTF8",
	ebmlDate:     "data",
	ebmlBinary:   "binary",
	ebmlMaster:   "master",
}

type ebmlAttribute struct {
	name string
	typ  ebmlType
	tag  ebmlTag
}

type ebmlTag map[uint64]ebmlAttribute

var ebmlGlobal = ebmlTag{
	0xbf: {name: "CRC-32", typ: ebmlBinary},
	0xec: {name: "Void", typ: ebmlBinary},
}

var ebmlHeader = ebmlTag{
	0x4286: {name: "EBMLVersion", typ: ebmlUinteger},
	0x42f7: {name: "EBMLReadVersion", typ: ebmlUinteger},
	0x42f2: {name: "EBMLMaxIDLength", typ: ebmlUinteger},
	0x42f3: {name: "EBMLMaxSizeLength", typ: ebmlUinteger},
	0x4282: {name: "DocType", typ: ebmlString},
	0x4287: {name: "DocTypeVersion", typ: ebmlUinteger},
	0x4285: {name: "DocTypeReadVersion", typ: ebmlUinteger},
}

var ebmlRoot = ebmlTag{
	0x1a45dfa3: {name: "EBML", typ: ebmlMaster, tag: ebmlHeader},
	0x18538067: {name: "Segment", typ: ebmlMaster, tag: mkvSegment},
}

type track struct {
	parentD             *decode.D
	number              int
	codec               string
	codecPrivatePos     int64
	codecPrivateTagSize int64
}

type simpleBlock struct {
	d *decode.D
	r ranges.Range
}

type decodeContext struct {
	currentTrack *track
	tracks       []*track
	simpleBlocks []simpleBlock
	blocks       []simpleBlock
}

/*
proc type_master {size _label extra} {
    upvar #0 "ebml_$extra" tags
    global ebml_Global
    set garbage_size 0

    # TODO: unknown-size might not be correct handled
    while {![end] && ($size > 0 || $size == -1)} {
        lassign [vint] tag_id_width tag_idnr
        set tag_id [format "%x" $tag_idnr]
        lassign [vint] tag_size_width tag_size_raw tag_size

        set tag_name "Unknown"
        set tag_type "binary"
        set tag_extra {}
        set tag_desc ""
        if {[dict exists $tags $tag_id]} {
            lassign [dict get $tags $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {[dict exists $ebml_Global $tag_id]} {
            lassign [dict get $ebml_Global $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {$size == -1} {
            incr garbage_size
            move [expr -($tag_id_width+$tag_size_width-1)]
            continue
        }

        if {$garbage_size != 0} {
            entry "Garbage" {} $garbage_size [expr [pos]-$garbage_size-$tag_id_width-$tag_size_width]
            set garbage_size 0
        }

        set type_fn "type_$tag_type"

        section "$tag_name ($tag_type)" {
            entry "ID" $tag_id $tag_id_width [expr [pos]-$tag_id_width-$tag_size_width]
            set tag_size_str $tag_size
            if {$tag_size_raw == 0xff} {
                append tag_size_str " (unknown)"
                set tag_size -1
            }
            entry "Size" "$tag_size_str" $tag_size_width [expr [pos]-$tag_size_width]
            $type_fn $tag_size $tag_name $tag_extra
        }

        if {$size == -1} {
            continue
        }
        incr size [expr -($tag_id_width+$tag_size_width+$tag_size)]
    }
}
*/

func decodeMaster(d *decode.D, bitsLimit int64, tag ebmlTag, dc *decodeContext) {
	tagEndBit := d.Pos() + bitsLimit

	d.FieldArrayFn("element", func(d *decode.D) {
		// var crcD *decode.D
		// var crcStart int64

		for d.Pos() < tagEndBit && d.NotEnd() {
			startPos := d.Pos()
			tagID := decodeRawVint(d)
			d.SeekAbs(startPos)

			a, ok := tag[tagID]
			if !ok {
				a, ok = ebmlGlobal[tagID]
				if !ok {
					panic("asdsad")
				}
			}

			const CRC = 0xbf
			const SimpleBlock = 0xa3
			const Block = 0xa1
			const CodecPrivate = 0x63a2
			const CodecID = 0x86
			const TrackNumber = 0xd7
			const TrackEntry = 0xae

			d.FieldStructFn("element", func(d *decode.D) {
				if tagID == TrackEntry {
					dc.currentTrack = &track{}
					dc.tracks = append(dc.tracks, dc.currentTrack)
				}

				d.FieldUFn("id", func() (uint64, decode.DisplayFormat, string) {
					n := decodeRawVint(d)
					return n, decode.NumberHex, a.name
				})
				// tagSize could be 0xffffffffffffff which means "unknown" size, then we will read until eof
				// TODO: should read until unknown id:
				//    The end of a Master-element with unknown size is determined by the beginning of the next
				//    element that is not a valid sub-element of that Master-element
				// TODO: should also handle garbage between
				tagSize := fieldDecodeVint(d, "size", decode.NumberDecimal)

				switch a.typ {
				case ebmlInteger:
					d.FieldS("value", int(tagSize)*8)
				case ebmlUinteger:
					v := d.FieldU("value", int(tagSize)*8)
					if dc.currentTrack != nil && tagID == TrackNumber {
						dc.currentTrack.number = int(v)
					}
				case ebmlFloat:
					d.FieldF("value", int(tagSize)*8)
				case ebmlString:
					v := d.FieldUTF8("value", int(tagSize))
					if dc.currentTrack != nil && tagID == CodecID {
						dc.currentTrack.codec = v
					}
				case ebmlUTF8:
					d.FieldUTF8("value", int(tagSize))
				case ebmlDate:
					// TODO:
					/*
						proc type_date {size label _extra} {
						    set s [clock scan {2001-01-01 00:00:00}]
						    set frac 0
						    switch $size {
						        0 {}
						        8 {
						            set nano [int64]
						            set s [clock add $s [expr $nano/1000000000] seconds]
						            set frac [expr ($nano%1000000000)/1000000000.0]
						        }
						        default {
						            bytes $size $label
						            return
						        }
						    }

						    entry $label "[clock format $s] ${frac}s" $size [expr [pos]-$size]
						}
					*/
					d.FieldBitBufLen("value", int64(tagSize)*8)
				case ebmlBinary:

					switch tagID {
					case SimpleBlock:

						// TODO: CodecPrivate
						// TODO: collect decode later when we know track codec?

						// d.DecodeLenFn(int64(tagSize)*8, func(d *decode.D) {
						// 	fieldDecodeVint(d, "track_number", decode.NumberDecimal)
						// 	d.FieldU16("timestamp")
						// 	d.FieldStructFn("flags", func(d *decode.D) {
						// 		d.FieldBool("key_frame")
						// 		d.FieldU3("reserved")
						// 		d.FieldBool("invisible")
						// 		d.FieldU2("lacing")
						// 		d.FieldBool("discardable")
						// 	})
						// 	// TODO: lacing
						// 	d.FieldBitBufLen("data", d.BitsLeft())

						// })

						dc.simpleBlocks = append(dc.simpleBlocks, simpleBlock{
							d: d,
							r: ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
						})

						d.SeekRel(int64(tagSize) * 8)
					case Block:

						// TODO: CodecPrivate
						// TODO: collect decode later when we know track codec?

						// d.DecodeLenFn(int64(tagSize)*8, func(d *decode.D) {
						// 	fieldDecodeVint(d, "track_number", decode.NumberDecimal)
						// 	d.FieldU16("timestamp")
						// 	d.FieldStructFn("flags", func(d *decode.D) {
						// 		d.FieldBool("key_frame")
						// 		d.FieldU3("reserved")
						// 		d.FieldBool("invisible")
						// 		d.FieldU2("lacing")
						// 		d.FieldBool("discardable")
						// 	})
						// 	// TODO: lacing
						// 	d.FieldBitBufLen("data", d.BitsLeft())

						// })

						dc.blocks = append(dc.blocks, simpleBlock{
							d: d,
							r: ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
						})

						d.SeekRel(int64(tagSize) * 8)
					case CodecPrivate:
						if dc.currentTrack != nil {
							dc.currentTrack.parentD = d
							dc.currentTrack.codecPrivatePos = d.Pos()
							dc.currentTrack.codecPrivateTagSize = int64(tagSize) * 8
						}
						d.SeekRel(int64(tagSize) * 8)
					default:
						d.FieldBitBufLen("value", int64(tagSize)*8)
						// if tagID == CRC {
						// 	crcD = d
						// 	crcStart = d.Pos()
						// }

					}

				case ebmlMaster:
					//d.DecodeLenFn(int64(tagSize)*8, func(d *decode.D) {
					decodeMaster(d, int64(tagSize)*8, a.tag, dc)
					//})
				}
			})
		}

		// if crcD != nil {
		// 	crcValue := crcD.FieldMustRemove("value")
		// 	elementCRC := &crc.CRC{Bits: 32, Current: 0xffff_ffff, Table: crc.IEEELETable}
		// 	//log.Printf("crc: %x-%x %d\n", crcStart/8, d.Pos()/8, (d.Pos()-crcStart)/8)
		// 	decode.MustCopy(elementCRC, d.BitBufRange(crcStart, d.Pos()-crcStart))
		// 	crcD.FieldChecksumRange("value", crcValue.Range.Start, crcValue.Range.Len, elementCRC.Sum(nil), decode.LittleEndian)
		// }
	})

}

func mkvDecode(d *decode.D) interface{} {
	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekBits(32) != ebmlHeaderID {
		d.Invalid("no EBML header found")
	}
	dc := &decodeContext{tracks: []*track{}}
	decodeMaster(d, d.BitsLeft(), ebmlRoot, dc)

	trackCodec := map[int]string{}

	for _, t := range dc.tracks {
		if t.codec != "" {
			trackCodec[t.number] = t.codec
		}
		// no CodecPrivate found
		if t.parentD == nil {
			continue
		}

		switch t.codec {
		case "A_VORBIS":
			t.parentD.DecodeRangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				numPackets := d.FieldU8("num_packets")
				// TODO: lacing
				packetLengths := []int64{}
				// Xiph-style lacing (similar to ogg) of n-1 packets, last is reset of block
				d.FieldArrayFn("lace", func(d *decode.D) {
					for i := uint64(0); i < numPackets; i++ {
						l := d.FieldUFn("lace", func() (uint64, decode.DisplayFormat, string) {
							var l uint64
							for {
								n := d.U8()
								l += n
								if n < 255 {
									return l, decode.NumberDecimal, ""
								}
							}
						})
						packetLengths = append(packetLengths, int64(l))
					}
				})
				d.FieldArrayFn("packet", func(d *decode.D) {
					for _, l := range packetLengths {
						d.FieldDecodeLen("packet", l*8, vorbisPacketFormat)
					}
					d.FieldDecodeLen("packet", d.BitsLeft(), vorbisPacketFormat)
				})
			})
		case "A_AAC":
			t.parentD.FieldDecodeRange("asc", t.codecPrivatePos, t.codecPrivateTagSize, mpegASCFrameFormat)
		default:
			t.parentD.FieldBitBufRange("value", t.codecPrivatePos, t.codecPrivateTagSize)
		}
	}

	for _, s := range dc.simpleBlocks {
		s.d.DecodeRangeFn(s.r.Start, s.r.Len, func(d *decode.D) {
			trackNumber := fieldDecodeVint(d, "track_number", decode.NumberDecimal)
			d.FieldU16("timestamp")
			d.FieldStructFn("flags", func(d *decode.D) {
				d.FieldBool("key_frame")
				d.FieldU3("reserved")
				d.FieldBool("invisible")
				d.FieldU2("lacing")
				d.FieldBool("discardable")
			})
			// TODO: lacing etc

			switch trackCodec[int(trackNumber)] {
			case "A_VORBIS":
				d.FieldDecodeLen("packet", d.BitsLeft(), vorbisPacketFormat)
			case "V_VP9":
				d.FieldDecodeLen("packet", d.BitsLeft(), vp9FrameFormat)
			case "V_VOBSUB":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegSPUFrameFormat)
			// case "A_AAC":
			// 	log.Println("bla")
			// 	d.FieldDecodeLen("packet", d.BitsLeft(), aacFrameFormat)
			default:
				d.FieldBitBufLen("data", d.BitsLeft())
			}

		})
	}

	for _, s := range dc.blocks {
		s.d.DecodeRangeFn(s.r.Start, s.r.Len, func(d *decode.D) {
			trackNumber := fieldDecodeVint(d, "track_number", decode.NumberDecimal)
			d.FieldU16("timestamp")
			d.FieldStructFn("flags", func(d *decode.D) {
				d.FieldU4("reserved")
				d.FieldBool("invisible")
				d.FieldU2("lacing")
				d.FieldBool("not_used")
			})

			// d.FieldStructFn("flags", func(d *decode.D) {
			// 	d.FieldBool("key_frame")
			// 	d.FieldU3("reserved")
			// 	d.FieldBool("invisible")
			// 	d.FieldU2("lacing")
			// 	d.FieldBool("discardable")
			// })
			// // TODO: lacing etc

			// switch trackCodec[int(trackNumber)] {
			// case "A_VORBIS":
			// 	d.FieldDecodeLen("packet", d.BitsLeft(), vorbisPacketFormat)
			// case "V_VP9":
			// 	d.FieldDecodeLen("packet", d.BitsLeft(), vp9FrameFormat)
			// // case "A_AAC":
			// // 	log.Println("bla")
			// // 	d.FieldDecodeLen("packet", d.BitsLeft(), aacFrameFormat)
			// default:
			// 	d.FieldBitBufLen("data", d.BitsLeft())
			// }

			switch trackCodec[int(trackNumber)] {
			case "S_VOBSUB":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegSPUFrameFormat)
			// case "A_AAC":
			// 	log.Println("bla")
			// 	d.FieldDecodeLen("packet", d.BitsLeft(), aacFrameFormat)
			default:
				d.FieldBitBufLen("data", d.BitsLeft())
			}

		})
	}

	return nil
}

// code below generated with ebml_gen.go
// https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml

var mkvSegment = ebmlTag{
	0x114d9b74: {name: "SeekHead", typ: ebmlMaster, tag: mkvSeekHead},
	0x1549a966: {name: "Info", typ: ebmlMaster, tag: mkvInfo},
	0x1f43b675: {name: "Cluster", typ: ebmlMaster, tag: mkvCluster},
	0x1654ae6b: {name: "Tracks", typ: ebmlMaster, tag: mkvTracks},
	0x1c53bb6b: {name: "Cues", typ: ebmlMaster, tag: mkvCues},
	0x1941a469: {name: "Attachments", typ: ebmlMaster, tag: mkvAttachments},
	0x1043a770: {name: "Chapters", typ: ebmlMaster, tag: mkvChapters},
	0x1254c367: {name: "Tags", typ: ebmlMaster, tag: mkvTags},
}

var mkvSeekHead = ebmlTag{
	0x4dbb: {name: "Seek", typ: ebmlMaster, tag: mkvSeek},
}

var mkvSeek = ebmlTag{
	0x53ab: {name: "SeekID", typ: ebmlBinary},
	0x53ac: {name: "SeekPosition", typ: ebmlUinteger},
}

var mkvInfo = ebmlTag{
	0x73a4:   {name: "SegmentUID", typ: ebmlBinary},
	0x7384:   {name: "SegmentFilename", typ: ebmlUTF8},
	0x3cb923: {name: "PrevUID", typ: ebmlBinary},
	0x3c83ab: {name: "PrevFilename", typ: ebmlUTF8},
	0x3eb923: {name: "NextUID", typ: ebmlBinary},
	0x3e83bb: {name: "NextFilename", typ: ebmlUTF8},
	0x4444:   {name: "SegmentFamily", typ: ebmlBinary},
	0x6924:   {name: "ChapterTranslate", typ: ebmlMaster, tag: mkvChapterTranslate},
	0x2ad7b1: {name: "TimestampScale", typ: ebmlUinteger},
	0x4489:   {name: "Duration", typ: ebmlFloat},
	0x4461:   {name: "DateUTC", typ: ebmlDate},
	0x7ba9:   {name: "Title", typ: ebmlUTF8},
	0x4d80:   {name: "MuxingApp", typ: ebmlUTF8},
	0x5741:   {name: "WritingApp", typ: ebmlUTF8},
}

var mkvChapterTranslate = ebmlTag{
	0x69fc: {name: "ChapterTranslateEditionUID", typ: ebmlUinteger},
	0x69bf: {name: "ChapterTranslateCodec", typ: ebmlUinteger},
	0x69a5: {name: "ChapterTranslateID", typ: ebmlBinary},
}

var mkvCluster = ebmlTag{
	0xe7:   {name: "Timestamp", typ: ebmlUinteger},
	0x5854: {name: "SilentTracks", typ: ebmlMaster, tag: mkvSilentTracks},
	0xa7:   {name: "Position", typ: ebmlUinteger},
	0xab:   {name: "PrevSize", typ: ebmlUinteger},
	0xa3:   {name: "SimpleBlock", typ: ebmlBinary},
	0xa0:   {name: "BlockGroup", typ: ebmlMaster, tag: mkvBlockGroup},
	0xaf:   {name: "EncryptedBlock", typ: ebmlBinary},
}

var mkvSilentTracks = ebmlTag{
	0x58d7: {name: "SilentTrackNumber", typ: ebmlUinteger},
}

var mkvBlockGroup = ebmlTag{
	0xa1:   {name: "Block", typ: ebmlBinary},
	0xa2:   {name: "BlockVirtual", typ: ebmlBinary},
	0x75a1: {name: "BlockAdditions", typ: ebmlMaster, tag: mkvBlockAdditions},
	0x9b:   {name: "BlockDuration", typ: ebmlUinteger},
	0xfa:   {name: "ReferencePriority", typ: ebmlUinteger},
	0xfb:   {name: "ReferenceBlock", typ: ebmlInteger},
	0xfd:   {name: "ReferenceVirtual", typ: ebmlInteger},
	0xa4:   {name: "CodecState", typ: ebmlBinary},
	0x75a2: {name: "DiscardPadding", typ: ebmlInteger},
	0x8e:   {name: "Slices", typ: ebmlMaster, tag: mkvSlices},
	0xc8:   {name: "ReferenceFrame", typ: ebmlMaster, tag: mkvReferenceFrame},
}

var mkvBlockAdditions = ebmlTag{
	0xa6: {name: "BlockMore", typ: ebmlMaster, tag: mkvBlockMore},
}

var mkvBlockMore = ebmlTag{
	0xee: {name: "BlockAddID", typ: ebmlUinteger},
	0xa5: {name: "BlockAdditional", typ: ebmlBinary},
}

var mkvSlices = ebmlTag{
	0xe8: {name: "TimeSlice", typ: ebmlMaster, tag: mkvTimeSlice},
}

var mkvTimeSlice = ebmlTag{
	0xcc: {name: "LaceNumber", typ: ebmlUinteger},
	0xcd: {name: "FrameNumber", typ: ebmlUinteger},
	0xcb: {name: "BlockAdditionID", typ: ebmlUinteger},
	0xce: {name: "Delay", typ: ebmlUinteger},
	0xcf: {name: "SliceDuration", typ: ebmlUinteger},
}

var mkvReferenceFrame = ebmlTag{
	0xc9: {name: "ReferenceOffset", typ: ebmlUinteger},
	0xca: {name: "ReferenceTimestamp", typ: ebmlUinteger},
}

var mkvTracks = ebmlTag{
	0xae: {name: "TrackEntry", typ: ebmlMaster, tag: mkvTrackEntry},
}

var mkvTrackEntry = ebmlTag{
	0xd7:     {name: "TrackNumber", typ: ebmlUinteger},
	0x73c5:   {name: "TrackUID", typ: ebmlUinteger},
	0x83:     {name: "TrackType", typ: ebmlUinteger},
	0xb9:     {name: "FlagEnabled", typ: ebmlUinteger},
	0x88:     {name: "FlagDefault", typ: ebmlUinteger},
	0x55aa:   {name: "FlagForced", typ: ebmlUinteger},
	0x9c:     {name: "FlagLacing", typ: ebmlUinteger},
	0x6de7:   {name: "MinCache", typ: ebmlUinteger},
	0x6df8:   {name: "MaxCache", typ: ebmlUinteger},
	0x23e383: {name: "DefaultDuration", typ: ebmlUinteger},
	0x234e7a: {name: "DefaultDecodedFieldDuration", typ: ebmlUinteger},
	0x23314f: {name: "TrackTimestampScale", typ: ebmlFloat},
	0x537f:   {name: "TrackOffset", typ: ebmlInteger},
	0x55ee:   {name: "MaxBlockAdditionID", typ: ebmlUinteger},
	0x41e4:   {name: "BlockAdditionMapping", typ: ebmlMaster, tag: mkvBlockAdditionMapping},
	0x536e:   {name: "Name", typ: ebmlUTF8},
	0x22b59c: {name: "Language", typ: ebmlString},
	0x22b59d: {name: "LanguageIETF", typ: ebmlString},
	0x86:     {name: "CodecID", typ: ebmlString},
	0x63a2:   {name: "CodecPrivate", typ: ebmlBinary},
	0x258688: {name: "CodecName", typ: ebmlUTF8},
	0x7446:   {name: "AttachmentLink", typ: ebmlUinteger},
	0x3a9697: {name: "CodecSettings", typ: ebmlUTF8},
	0x3b4040: {name: "CodecInfoURL", typ: ebmlString},
	0x26b240: {name: "CodecDownloadURL", typ: ebmlString},
	0xaa:     {name: "CodecDecodeAll", typ: ebmlUinteger},
	0x6fab:   {name: "TrackOverlay", typ: ebmlUinteger},
	0x56aa:   {name: "CodecDelay", typ: ebmlUinteger},
	0x56bb:   {name: "SeekPreRoll", typ: ebmlUinteger},
	0x6624:   {name: "TrackTranslate", typ: ebmlMaster, tag: mkvTrackTranslate},
	0xe0:     {name: "Video", typ: ebmlMaster, tag: mkvVideo},
	0xe1:     {name: "Audio", typ: ebmlMaster, tag: mkvAudio},
	0xe2:     {name: "TrackOperation", typ: ebmlMaster, tag: mkvTrackOperation},
	0xc0:     {name: "TrickTrackUID", typ: ebmlUinteger},
	0xc1:     {name: "TrickTrackSegmentUID", typ: ebmlBinary},
	0xc6:     {name: "TrickTrackFlag", typ: ebmlUinteger},
	0xc7:     {name: "TrickMasterTrackUID", typ: ebmlUinteger},
	0xc4:     {name: "TrickMasterTrackSegmentUID", typ: ebmlBinary},
	0x6d80:   {name: "ContentEncodings", typ: ebmlMaster, tag: mkvContentEncodings},
}

var mkvBlockAdditionMapping = ebmlTag{
	0x41f0: {name: "BlockAddIDValue", typ: ebmlUinteger},
	0x41a4: {name: "BlockAddIDName", typ: ebmlString},
	0x41e7: {name: "BlockAddIDType", typ: ebmlUinteger},
	0x41ed: {name: "BlockAddIDExtraData", typ: ebmlBinary},
}

var mkvTrackTranslate = ebmlTag{
	0x66fc: {name: "TrackTranslateEditionUID", typ: ebmlUinteger},
	0x66bf: {name: "TrackTranslateCodec", typ: ebmlUinteger},
	0x66a5: {name: "TrackTranslateTrackID", typ: ebmlBinary},
}

var mkvVideo = ebmlTag{
	0x9a:     {name: "FlagInterlaced", typ: ebmlUinteger},
	0x9d:     {name: "FieldOrder", typ: ebmlUinteger},
	0x53b8:   {name: "StereoMode", typ: ebmlUinteger},
	0x53c0:   {name: "AlphaMode", typ: ebmlUinteger},
	0x53b9:   {name: "OldStereoMode", typ: ebmlUinteger},
	0xb0:     {name: "PixelWidth", typ: ebmlUinteger},
	0xba:     {name: "PixelHeight", typ: ebmlUinteger},
	0x54aa:   {name: "PixelCropBottom", typ: ebmlUinteger},
	0x54bb:   {name: "PixelCropTop", typ: ebmlUinteger},
	0x54cc:   {name: "PixelCropLeft", typ: ebmlUinteger},
	0x54dd:   {name: "PixelCropRight", typ: ebmlUinteger},
	0x54b0:   {name: "DisplayWidth", typ: ebmlUinteger},
	0x54ba:   {name: "DisplayHeight", typ: ebmlUinteger},
	0x54b2:   {name: "DisplayUnit", typ: ebmlUinteger},
	0x54b3:   {name: "AspectRatioType", typ: ebmlUinteger},
	0x2eb524: {name: "ColourSpace", typ: ebmlBinary},
	0x2fb523: {name: "GammaValue", typ: ebmlFloat},
	0x2383e3: {name: "FrameRate", typ: ebmlFloat},
	0x55b0:   {name: "Colour", typ: ebmlMaster, tag: mkvColour},
	0x7670:   {name: "Projection", typ: ebmlMaster, tag: mkvProjection},
}

var mkvColour = ebmlTag{
	0x55b1: {name: "MatrixCoefficients", typ: ebmlUinteger},
	0x55b2: {name: "BitsPerChannel", typ: ebmlUinteger},
	0x55b3: {name: "ChromaSubsamplingHorz", typ: ebmlUinteger},
	0x55b4: {name: "ChromaSubsamplingVert", typ: ebmlUinteger},
	0x55b5: {name: "CbSubsamplingHorz", typ: ebmlUinteger},
	0x55b6: {name: "CbSubsamplingVert", typ: ebmlUinteger},
	0x55b7: {name: "ChromaSitingHorz", typ: ebmlUinteger},
	0x55b8: {name: "ChromaSitingVert", typ: ebmlUinteger},
	0x55b9: {name: "Range", typ: ebmlUinteger},
	0x55ba: {name: "TransferCharacteristics", typ: ebmlUinteger},
	0x55bb: {name: "Primaries", typ: ebmlUinteger},
	0x55bc: {name: "MaxCLL", typ: ebmlUinteger},
	0x55bd: {name: "MaxFALL", typ: ebmlUinteger},
	0x55d0: {name: "MasteringMetadata", typ: ebmlMaster, tag: mkvMasteringMetadata},
}

var mkvMasteringMetadata = ebmlTag{
	0x55d1: {name: "PrimaryRChromaticityX", typ: ebmlFloat},
	0x55d2: {name: "PrimaryRChromaticityY", typ: ebmlFloat},
	0x55d3: {name: "PrimaryGChromaticityX", typ: ebmlFloat},
	0x55d4: {name: "PrimaryGChromaticityY", typ: ebmlFloat},
	0x55d5: {name: "PrimaryBChromaticityX", typ: ebmlFloat},
	0x55d6: {name: "PrimaryBChromaticityY", typ: ebmlFloat},
	0x55d7: {name: "WhitePointChromaticityX", typ: ebmlFloat},
	0x55d8: {name: "WhitePointChromaticityY", typ: ebmlFloat},
	0x55d9: {name: "LuminanceMax", typ: ebmlFloat},
	0x55da: {name: "LuminanceMin", typ: ebmlFloat},
}

var mkvProjection = ebmlTag{
	0x7671: {name: "ProjectionType", typ: ebmlUinteger},
	0x7672: {name: "ProjectionPrivate", typ: ebmlBinary},
	0x7673: {name: "ProjectionPoseYaw", typ: ebmlFloat},
	0x7674: {name: "ProjectionPosePitch", typ: ebmlFloat},
	0x7675: {name: "ProjectionPoseRoll", typ: ebmlFloat},
}

var mkvAudio = ebmlTag{
	0xb5:   {name: "SamplingFrequency", typ: ebmlFloat},
	0x78b5: {name: "OutputSamplingFrequency", typ: ebmlFloat},
	0x9f:   {name: "Channels", typ: ebmlUinteger},
	0x7d7b: {name: "ChannelPositions", typ: ebmlBinary},
	0x6264: {name: "BitDepth", typ: ebmlUinteger},
}

var mkvTrackOperation = ebmlTag{
	0xe3: {name: "TrackCombinePlanes", typ: ebmlMaster, tag: mkvTrackCombinePlanes},
	0xe9: {name: "TrackJoinBlocks", typ: ebmlMaster, tag: mkvTrackJoinBlocks},
}

var mkvTrackCombinePlanes = ebmlTag{
	0xe4: {name: "TrackPlane", typ: ebmlMaster, tag: mkvTrackPlane},
}

var mkvTrackPlane = ebmlTag{
	0xe5: {name: "TrackPlaneUID", typ: ebmlUinteger},
	0xe6: {name: "TrackPlaneType", typ: ebmlUinteger},
}

var mkvTrackJoinBlocks = ebmlTag{
	0xed: {name: "TrackJoinUID", typ: ebmlUinteger},
}

var mkvContentEncodings = ebmlTag{
	0x6240: {name: "ContentEncoding", typ: ebmlMaster, tag: mkvContentEncoding},
}

var mkvContentEncoding = ebmlTag{
	0x5031: {name: "ContentEncodingOrder", typ: ebmlUinteger},
	0x5032: {name: "ContentEncodingScope", typ: ebmlUinteger},
	0x5033: {name: "ContentEncodingType", typ: ebmlUinteger},
	0x5034: {name: "ContentCompression", typ: ebmlMaster, tag: mkvContentCompression},
	0x5035: {name: "ContentEncryption", typ: ebmlMaster, tag: mkvContentEncryption},
}

var mkvContentCompression = ebmlTag{
	0x4254: {name: "ContentCompAlgo", typ: ebmlUinteger},
	0x4255: {name: "ContentCompSettings", typ: ebmlBinary},
}

var mkvContentEncryption = ebmlTag{
	0x47e1: {name: "ContentEncAlgo", typ: ebmlUinteger},
	0x47e2: {name: "ContentEncKeyID", typ: ebmlBinary},
	0x47e7: {name: "ContentEncAESSettings", typ: ebmlMaster, tag: mkvContentEncAESSettings},
	0x47e3: {name: "ContentSignature", typ: ebmlBinary},
	0x47e4: {name: "ContentSigKeyID", typ: ebmlBinary},
	0x47e5: {name: "ContentSigAlgo", typ: ebmlUinteger},
	0x47e6: {name: "ContentSigHashAlgo", typ: ebmlUinteger},
}

var mkvContentEncAESSettings = ebmlTag{
	0x47e8: {name: "AESSettingsCipherMode", typ: ebmlUinteger},
}

var mkvCues = ebmlTag{
	0xbb: {name: "CuePoint", typ: ebmlMaster, tag: mkvCuePoint},
}

var mkvCuePoint = ebmlTag{
	0xb3: {name: "CueTime", typ: ebmlUinteger},
	0xb7: {name: "CueTrackPositions", typ: ebmlMaster, tag: mkvCueTrackPositions},
}

var mkvCueTrackPositions = ebmlTag{
	0xf7:   {name: "CueTrack", typ: ebmlUinteger},
	0xf1:   {name: "CueClusterPosition", typ: ebmlUinteger},
	0xf0:   {name: "CueRelativePosition", typ: ebmlUinteger},
	0xb2:   {name: "CueDuration", typ: ebmlUinteger},
	0x5378: {name: "CueBlockNumber", typ: ebmlUinteger},
	0xea:   {name: "CueCodecState", typ: ebmlUinteger},
	0xdb:   {name: "CueReference", typ: ebmlMaster, tag: mkvCueReference},
}

var mkvCueReference = ebmlTag{
	0x96:   {name: "CueRefTime", typ: ebmlUinteger},
	0x97:   {name: "CueRefCluster", typ: ebmlUinteger},
	0x535f: {name: "CueRefNumber", typ: ebmlUinteger},
	0xeb:   {name: "CueRefCodecState", typ: ebmlUinteger},
}

var mkvAttachments = ebmlTag{
	0x61a7: {name: "AttachedFile", typ: ebmlMaster, tag: mkvAttachedFile},
}

var mkvAttachedFile = ebmlTag{
	0x467e: {name: "FileDescription", typ: ebmlUTF8},
	0x466e: {name: "FileName", typ: ebmlUTF8},
	0x4660: {name: "FileMimeType", typ: ebmlString},
	0x465c: {name: "FileData", typ: ebmlBinary},
	0x46ae: {name: "FileUID", typ: ebmlUinteger},
	0x4675: {name: "FileReferral", typ: ebmlBinary},
	0x4661: {name: "FileUsedStartTime", typ: ebmlUinteger},
	0x4662: {name: "FileUsedEndTime", typ: ebmlUinteger},
}

var mkvChapters = ebmlTag{
	0x45b9: {name: "EditionEntry", typ: ebmlMaster, tag: mkvEditionEntry},
}

var mkvEditionEntry = ebmlTag{
	0x45bc: {name: "EditionUID", typ: ebmlUinteger},
	0x45bd: {name: "EditionFlagHidden", typ: ebmlUinteger},
	0x45db: {name: "EditionFlagDefault", typ: ebmlUinteger},
	0x45dd: {name: "EditionFlagOrdered", typ: ebmlUinteger},
	0xb6:   {name: "ChapterAtom", typ: ebmlMaster, tag: mkvChapterAtom},
}

var mkvChapterAtom = ebmlTag{
	0x73c4: {name: "ChapterUID", typ: ebmlUinteger},
	0x5654: {name: "ChapterStringUID", typ: ebmlUTF8},
	0x91:   {name: "ChapterTimeStart", typ: ebmlUinteger},
	0x92:   {name: "ChapterTimeEnd", typ: ebmlUinteger},
	0x98:   {name: "ChapterFlagHidden", typ: ebmlUinteger},
	0x4598: {name: "ChapterFlagEnabled", typ: ebmlUinteger},
	0x6e67: {name: "ChapterSegmentUID", typ: ebmlBinary},
	0x6ebc: {name: "ChapterSegmentEditionUID", typ: ebmlUinteger},
	0x63c3: {name: "ChapterPhysicalEquiv", typ: ebmlUinteger},
	0x8f:   {name: "ChapterTrack", typ: ebmlMaster, tag: mkvChapterTrack},
	0x80:   {name: "ChapterDisplay", typ: ebmlMaster, tag: mkvChapterDisplay},
	0x6944: {name: "ChapProcess", typ: ebmlMaster, tag: mkvChapProcess},
}

var mkvChapterTrack = ebmlTag{
	0x89: {name: "ChapterTrackUID", typ: ebmlUinteger},
}

var mkvChapterDisplay = ebmlTag{
	0x85:   {name: "ChapString", typ: ebmlUTF8},
	0x437c: {name: "ChapLanguage", typ: ebmlString},
	0x437d: {name: "ChapLanguageIETF", typ: ebmlString},
	0x437e: {name: "ChapCountry", typ: ebmlString},
}

var mkvChapProcess = ebmlTag{
	0x6955: {name: "ChapProcessCodecID", typ: ebmlUinteger},
	0x450d: {name: "ChapProcessPrivate", typ: ebmlBinary},
	0x6911: {name: "ChapProcessCommand", typ: ebmlMaster, tag: mkvChapProcessCommand},
}

var mkvChapProcessCommand = ebmlTag{
	0x6922: {name: "ChapProcessTime", typ: ebmlUinteger},
	0x6933: {name: "ChapProcessData", typ: ebmlBinary},
}

var mkvTags = ebmlTag{
	0x7373: {name: "Tag", typ: ebmlMaster, tag: mkvTag},
}

var mkvTag = ebmlTag{
	0x63c0: {name: "Targets", typ: ebmlMaster, tag: mkvTargets},
	0x67c8: {name: "SimpleTag", typ: ebmlMaster, tag: mkvSimpleTag},
}

var mkvTargets = ebmlTag{
	0x68ca: {name: "TargetTypeValue", typ: ebmlUinteger},
	0x63ca: {name: "TargetType", typ: ebmlString},
	0x63c5: {name: "TagTrackUID", typ: ebmlUinteger},
	0x63c9: {name: "TagEditionUID", typ: ebmlUinteger},
	0x63c4: {name: "TagChapterUID", typ: ebmlUinteger},
	0x63c6: {name: "TagAttachmentUID", typ: ebmlUinteger},
}

var mkvSimpleTag = ebmlTag{
	0x45a3: {name: "TagName", typ: ebmlUTF8},
	0x447a: {name: "TagLanguage", typ: ebmlString},
	0x447b: {name: "TagLanguageIETF", typ: ebmlString},
	0x4484: {name: "TagDefault", typ: ebmlUinteger},
	0x4487: {name: "TagString", typ: ebmlUTF8},
	0x4485: {name: "TagBinary", typ: ebmlBinary},
}
