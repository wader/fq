package mkv

// https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml
//go:generate sh -c "go run ebml_gen.go ebml_matroska.xml mkv | gofmt > ebml_matroska.go"

// https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
// https://matroska.org/technical/specs/index.html
// https://www.matroska.org/technical/basics.html
// https://www.matroska.org/technical/codec_specs.html
// https://wiki.xiph.org/MatroskaOpus

// TODO: refactor simepleblock/block to just defer decode etc?
// TODO: CRC
// TODO: value to names (TrackType etc)
// TODO: lacing
// TODO: handle garbage (see tcl and example files)

import (
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/ranges"
)

var aacFrameFormat []*decode.Format
var av1CCRFormat []*decode.Format
var av1FrameFormat []*decode.Format
var flacFrameFormat []*decode.Format
var flacMetadatablockFormat []*decode.Format
var mp3FrameFormat []*decode.Format
var mpegASCFrameFormat []*decode.Format
var mpegAVCDCRFormat []*decode.Format
var mpegAVCSampleFormat []*decode.Format
var mpegHEVCDCRFormat []*decode.Format
var mpegHEVCSampleFormat []*decode.Format
var mpegSPUFrameFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp9FrameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MKV,
		Description: "Matroska",
		Groups:      []string{format.PROBE},
		DecodeFn:    mkvDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AV1_CCR}, Formats: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Formats: &av1FrameFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3FrameFormat},
			{Names: []string{format.MPEG_AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.MPEG_ASC}, Formats: &mpegASCFrameFormat},
			{Names: []string{format.MPEG_AVC_DCR}, Formats: &mpegAVCDCRFormat},
			{Names: []string{format.MPEG_AVC}, Formats: &mpegAVCSampleFormat},
			{Names: []string{format.MPEG_HEVC_DCR}, Formats: &mpegHEVCDCRFormat},
			{Names: []string{format.MPEG_HEVC}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.MPEG_SPU}, Formats: &mpegSPUFrameFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
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
	decodeOpts          []decode.Options
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

func decodeMaster(d *decode.D, bitsLimit int64, tag ebmlTag, dc *decodeContext) {
	tagEndBit := d.Pos() + bitsLimit

	d.FieldArrayFn("elements", func(d *decode.D) {
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
						dc.simpleBlocks = append(dc.simpleBlocks, simpleBlock{
							d: d,
							r: ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
						})
						d.SeekRel(int64(tagSize) * 8)
					case Block:
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
					decodeMaster(d, int64(tagSize)*8, a.tag, dc)
				}
			})
		}

		// if crcD != nil {
		// 	crcValue := crcD.FieldMustRemove("value")
		// 	elementCRC := &crc.CRC{Bits: 32, Current: 0xffff_ffff, Table: crc.IEEELETable}
		// 	//log.Printf("crc: %x-%x %d\n", crcStart/8, d.Pos()/8, (d.Pos()-crcStart)/8)
		// 	ioextra.MustCopy(elementCRC, d.BitBufRange(crcStart, d.Pos()-crcStart))
		// 	crcD.FieldChecksumRange("value", crcValue.Range.Start, crcValue.Range.Len, elementCRC.Sum(nil), decode.LittleEndian)
		// }
	})

}

func mkvDecode(d *decode.D, in interface{}) interface{} {
	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekBits(32) != ebmlHeaderID {
		d.Invalid("no EBML header found")
	}
	dc := &decodeContext{tracks: []*track{}}
	decodeMaster(d, d.BitsLeft(), ebmlRoot, dc)

	trackNumberToTrack := map[int]*track{}
	for _, t := range dc.tracks {
		trackNumberToTrack[t.number] = t
	}

	for _, t := range dc.tracks {
		// no CodecPrivate found
		if t.parentD == nil {
			continue
		}

		// TODO: refactor, one DecodeRangeFn or refactor to add FieldDecodeRangeFn?

		switch t.codec {
		case "A_VORBIS":
			t.parentD.DecodeRangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				numPackets := d.FieldU8("num_packets")
				// TODO: lacing
				packetLengths := []int64{}
				// Xiph-style lacing (similar to ogg) of n-1 packets, last is reset of block
				d.FieldArrayFn("laces", func(d *decode.D) {
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
				d.FieldArrayFn("packets", func(d *decode.D) {
					for _, l := range packetLengths {
						d.FieldDecodeLen("packet", l*8, vorbisPacketFormat)
					}
					d.FieldDecodeLen("packet", d.BitsLeft(), vorbisPacketFormat)
				})
			})
		case "A_AAC":
			t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegASCFrameFormat)
		case "A_OPUS":
			t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, opusPacketFrameFormat)
		case "A_FLAC":
			t.parentD.DecodeRangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldStructFn("value", func(d *decode.D) {
					d.FieldValidateUTF8("magic", "fLaC")
					d.FieldArrayFn("metadatablocks", func(d *decode.D) {
						for {
							_, dv := d.FieldDecode("metadatablock", flacMetadatablockFormat)
							flacMetadatablockOut, ok := dv.(format.FlacMetadatablockOut)
							if !ok {
								d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", dv))
							}
							if flacMetadatablockOut.HasStreamInfo {
								t.decodeOpts = append(t.decodeOpts,
									decode.FormatOptions{InArg: format.FlacFrameIn{StreamInfo: flacMetadatablockOut.StreamInfo}})
							}
							if flacMetadatablockOut.IsLastBlock {
								return
							}
						}
					})
				})
			})
		case "V_MPEG4/ISO/AVC":
			_, dv := t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegAVCDCRFormat)
			avcDcrOut, ok := dv.(format.AvcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected AvcDcrOut got %#+v", dv))
			}
			t.decodeOpts = append(t.decodeOpts,
				decode.FormatOptions{InArg: format.AvcIn{LengthSize: avcDcrOut.LengthSize}})
		case "V_MPEGH/ISO/HEVC":
			_, dv := t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegHEVCDCRFormat)
			hevcDcrOut, ok := dv.(format.HevcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected HevcDcrOut got %#+v", dv))
			}
			t.decodeOpts = append(t.decodeOpts,
				decode.FormatOptions{InArg: format.HevcIn{LengthSize: hevcDcrOut.LengthSize}})
		case "V_AV1":
			t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, av1CCRFormat)
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

			codec := ""
			var decodeOpts []decode.Options
			t := trackNumberToTrack[int(trackNumber)]
			if t != nil {
				codec = t.codec
				decodeOpts = t.decodeOpts
			}

			switch codec {
			case "A_VORBIS":
				d.FieldDecodeLen("packet", d.BitsLeft(), vorbisPacketFormat)
			case "A_MPEG/L3":
				d.FieldDecodeLen("packet", d.BitsLeft(), mp3FrameFormat)
			case "A_FLAC":
				d.FieldDecodeLen("packet", d.BitsLeft(), flacFrameFormat, decodeOpts...)
				// TODO: could to md5 here somehow, see flac.go
			case "V_VP9":
				d.FieldDecodeLen("packet", d.BitsLeft(), vp9FrameFormat)
			case "V_AV1":
				d.FieldDecodeLen("packet", d.BitsLeft(), av1FrameFormat)
			case "V_VOBSUB":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegSPUFrameFormat)
			case "V_MPEG4/ISO/AVC":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegAVCSampleFormat, decodeOpts...)
			case "V_MPEGH/ISO/HEVC":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegHEVCSampleFormat, decodeOpts...)
			case "A_AAC":
				d.FieldDecodeLen("packet", d.BitsLeft(), aacFrameFormat)
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

			codec := ""
			t := trackNumberToTrack[int(trackNumber)]
			if t != nil {
				codec = t.codec
			}

			switch codec {
			case "S_VOBSUB":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegSPUFrameFormat)
			case "A_OPUS":
				d.FieldDecodeLen("packet", d.BitsLeft(), opusPacketFrameFormat)
			case "A_MPEG/L3":
				d.FieldDecodeLen("packet", d.BitsLeft(), mp3FrameFormat)
			case "V_AV1":
				d.FieldDecodeLen("packet", d.BitsLeft(), av1FrameFormat)
			default:
				d.FieldBitBufLen("data", d.BitsLeft())
			}

		})
	}

	return nil
}
