package matroska

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
	"fq/pkg/format/matroska/ebml"
	"fq/pkg/format/matroska/ebml_matroska"
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
var mpegAVCAUFormat []*decode.Format
var mpegHEVCDCRFormat []*decode.Format
var mpegHEVCSampleFormat []*decode.Format
var mpegSPUFrameFormat []*decode.Format
var mpegPESPacketSampleFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp8FrameFormat []*decode.Format
var vp9FrameFormat []*decode.Format
var vp9CFMFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MATROSKA,
		Description: "Matroska file",
		Groups:      []string{format.PROBE},
		DecodeFn:    matroskaDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AV1_CCR}, Formats: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Formats: &av1FrameFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3FrameFormat},
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.MPEG_ASC}, Formats: &mpegASCFrameFormat},
			{Names: []string{format.MPEG_AVC_DCR}, Formats: &mpegAVCDCRFormat},
			{Names: []string{format.MPEG_AVC_AU}, Formats: &mpegAVCAUFormat},
			{Names: []string{format.MPEG_HEVC_DCR}, Formats: &mpegHEVCDCRFormat},
			{Names: []string{format.MPEG_HEVC_AU}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.MPEG_SPU}, Formats: &mpegSPUFrameFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Formats: &mpegPESPacketSampleFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP8_FRAME}, Formats: &vp8FrameFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
			{Names: []string{format.VP9_CFM}, Formats: &vp9CFMFormat},
		},
	})
}

// TODO: smarter?
func decodeRawVintWidth(d *decode.D) (uint64, int) {
	n := d.U8()
	w := 1
	for i := 0; i <= 7 && (n&(1<<(7-i))) == 0; i++ {
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

func fieldDecodeVint(d *decode.D, name string, displayFormat decode.DisplayFormat) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return decodeVint(d), displayFormat, ""
	})
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

func decodeMaster(d *decode.D, bitsLimit int64, tag ebml.Tag, dc *decodeContext) {
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
				a, ok = ebml.Global[tagID]
				if !ok {
					d.Invalid(fmt.Sprintf("unknown id %d", tagID))
				}
			}

			d.FieldStructFn("element", func(d *decode.D) {
				if tagID == ebml_matroska.TrackEntryID {
					dc.currentTrack = &track{}
					dc.tracks = append(dc.tracks, dc.currentTrack)
				}

				d.FieldUFn("id", func() (uint64, decode.DisplayFormat, string) {
					n := decodeRawVint(d)
					return n, decode.NumberHex, a.Name
				})
				// tagSize could be 0xffffffffffffff which means "unknown" size, then we will read until eof
				// TODO: should read until unknown id:
				//    The end of a Master-element with unknown size is determined by the beginning of the next
				//    element that is not a valid sub-element of that Master-element
				// TODO: should also handle garbage between
				tagSize := fieldDecodeVint(d, "size", decode.NumberDecimal)

				if tagSize > 8 &&
					(a.Type == ebml.Integer ||
						a.Type == ebml.Uinteger ||
						a.Type == ebml.Float) {
					d.Invalid(fmt.Sprintf("invalid tagSize %d for non-master type", tagSize))
				}

				switch a.Type {
				case ebml.Integer:
					d.FieldSFn("value", func() (int64, decode.DisplayFormat, string) {
						n := d.S(int(tagSize) * 8)
						if len(a.UintegerEnums) > 0 {
							// TODO: use enum Defintion as description
							return n, decode.NumberDecimal, a.IntegerEnums[n].Label

						}
						return n, decode.NumberDecimal, ""
					})
				case ebml.Uinteger:
					v := d.FieldUFn("value", func() (uint64, decode.DisplayFormat, string) {
						n := d.U(int(tagSize) * 8)
						if len(a.UintegerEnums) > 0 {
							// TODO: use enum Defintion as description
							return n, decode.NumberDecimal, a.UintegerEnums[n].Label

						}
						return n, decode.NumberDecimal, ""
					})

					if dc.currentTrack != nil && tagID == ebml_matroska.TrackNumberID {
						dc.currentTrack.number = int(v)
					}
				case ebml.Float:
					d.FieldF("value", int(tagSize)*8)
				case ebml.String:
					v := d.FieldStrFn("value", func() (string, string) {
						s := d.UTF8(int(tagSize))
						if len(a.StringEnums) > 0 {
							// TODO: use enum Defintion as description
							return s, a.StringEnums[s].Label

						}
						return s, ""
					})

					if dc.currentTrack != nil && tagID == ebml_matroska.CodecIDID {
						dc.currentTrack.codec = v
					}
				case ebml.UTF8:
					d.FieldUTF8("value", int(tagSize))
				case ebml.Date:
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
				case ebml.Binary:
					switch tagID {
					case ebml_matroska.SimpleBlockID:
						dc.simpleBlocks = append(dc.simpleBlocks, simpleBlock{
							d: d,
							r: ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
						})
						d.SeekRel(int64(tagSize) * 8)
					case ebml_matroska.BlockID:
						dc.blocks = append(dc.blocks, simpleBlock{
							d: d,
							r: ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
						})
						d.SeekRel(int64(tagSize) * 8)
					case ebml_matroska.CodecPrivateID:
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

				case ebml.Master:
					decodeMaster(d, int64(tagSize)*8, a.Tag, dc)
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

func matroskaDecode(d *decode.D, in interface{}) interface{} {
	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekBits(32) != ebmlHeaderID {
		d.Invalid("no EBML header found")
	}
	dc := &decodeContext{tracks: []*track{}}
	decodeMaster(d, d.BitsLeft(), ebml_matroska.Root, dc)

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
							_, v := d.FieldDecode("metadatablock", flacMetadatablockFormat)
							flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
							if !ok {
								d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
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
			_, v := t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegAVCDCRFormat)
			avcDcrOut, ok := v.(format.AvcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			t.decodeOpts = append(t.decodeOpts,
				decode.FormatOptions{InArg: format.AvcIn{LengthSize: avcDcrOut.LengthSize}})
		case "V_MPEGH/ISO/HEVC":
			_, v := t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegHEVCDCRFormat)
			hevcDcrOut, ok := v.(format.HevcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
			}
			t.decodeOpts = append(t.decodeOpts,
				decode.FormatOptions{InArg: format.HevcIn{LengthSize: hevcDcrOut.LengthSize}})
		case "V_AV1":
			t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, av1CCRFormat)
		case "V_VP9":
			t.parentD.FieldDecodeRange("value", t.codecPrivatePos, t.codecPrivateTagSize, vp9CFMFormat)
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
				// TODO: could use md5 here somehow, see flac.go
			case "V_VP8":
				d.FieldDecodeLen("packet", d.BitsLeft(), vp8FrameFormat)
			case "V_VP9":
				d.FieldDecodeLen("packet", d.BitsLeft(), vp9FrameFormat)
			case "V_AV1":
				d.FieldDecodeLen("packet", d.BitsLeft(), av1FrameFormat)
			case "V_VOBSUB":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegSPUFrameFormat)
			case "V_MPEG4/ISO/AVC":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegAVCAUFormat, decodeOpts...)
			case "V_MPEGH/ISO/HEVC":
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegHEVCSampleFormat, decodeOpts...)
			case "V_MPEG2":
				// TODO: many samples? maybe should be in es decoder etc?
				d.FieldDecodeLen("packet", d.BitsLeft(), mpegPESPacketSampleFormat, decodeOpts...)
			case "A_AAC":
				d.FieldDecodeLen("packet", d.BitsLeft(), aacFrameFormat)
			case "A_OPUS":
				d.FieldDecodeLen("packet", d.BitsLeft(), opusPacketFrameFormat)
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
