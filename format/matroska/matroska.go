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
// TODO: could use md5 here somehow, see flac.go

import (
	"embed"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/matroska/ebml"
	"github.com/wader/fq/format/matroska/ebml_matroska"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed matroska.jq
var matroskaFS embed.FS

var aacFrameFormat decode.Group
var av1CCRFormat decode.Group
var av1FrameFormat decode.Group
var flacFrameFormat decode.Group
var flacMetadatablocksFormat decode.Group
var imageFormat decode.Group
var mp3FrameFormat decode.Group
var mpegASCFrameFormat decode.Group
var mpegAVCAUFormat decode.Group
var mpegAVCDCRFormat decode.Group
var mpegHEVCDCRFormat decode.Group
var mpegHEVCSampleFormat decode.Group
var mpegPESPacketSampleFormat decode.Group
var mpegSPUFrameFormat decode.Group
var opusPacketFrameFormat decode.Group
var vorbisPacketFormat decode.Group
var vp8FrameFormat decode.Group
var vp9CFMFormat decode.Group
var vp9FrameFormat decode.Group

var codecToFormat map[string]*decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.MATROSKA,
		Description: "Matroska file",
		Groups:      []string{format.PROBE},
		DecodeFn:    matroskaDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AAC_FRAME}, Group: &aacFrameFormat},
			{Names: []string{format.AV1_CCR}, Group: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Group: &av1FrameFormat},
			{Names: []string{format.AVC_AU}, Group: &mpegAVCAUFormat},
			{Names: []string{format.AVC_DCR}, Group: &mpegAVCDCRFormat},
			{Names: []string{format.FLAC_FRAME}, Group: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCKS}, Group: &flacMetadatablocksFormat},
			{Names: []string{format.HEVC_AU}, Group: &mpegHEVCSampleFormat},
			{Names: []string{format.HEVC_DCR}, Group: &mpegHEVCDCRFormat},
			{Names: []string{format.IMAGE}, Group: &imageFormat},
			{Names: []string{format.MP3_FRAME}, Group: &mp3FrameFormat},
			{Names: []string{format.MPEG_ASC}, Group: &mpegASCFrameFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Group: &mpegPESPacketSampleFormat},
			{Names: []string{format.MPEG_SPU}, Group: &mpegSPUFrameFormat},
			{Names: []string{format.OPUS_PACKET}, Group: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Group: &vorbisPacketFormat},
			{Names: []string{format.VP8_FRAME}, Group: &vp8FrameFormat},
			{Names: []string{format.VP9_CFM}, Group: &vp9CFMFormat},
			{Names: []string{format.VP9_FRAME}, Group: &vp9FrameFormat},
		},
		Functions: []string{"_help"},
	})
	interp.RegisterFS(matroskaFS)

	codecToFormat = map[string]*decode.Group{
		"A_VORBIS":         &vorbisPacketFormat,
		"A_MPEG/L3":        &mp3FrameFormat,
		"A_FLAC":           &flacFrameFormat,
		"A_AAC":            &aacFrameFormat,
		"A_OPUS":           &opusPacketFrameFormat,
		"V_VP8":            &vp8FrameFormat,
		"V_VP9":            &vp9FrameFormat,
		"V_AV1":            &av1FrameFormat,
		"V_VOBSUB":         &mpegSPUFrameFormat,
		"V_MPEG4/ISO/AVC":  &mpegAVCAUFormat,
		"V_MPEGH/ISO/HEVC": &mpegHEVCSampleFormat,
		"V_MPEG2":          &mpegPESPacketSampleFormat,
		"S_VOBSUB":         &mpegSPUFrameFormat,
	}
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

type track struct {
	parentD             *decode.D
	number              int
	codec               string
	codecPrivatePos     int64
	codecPrivateTagSize int64
	formatInArg         any
}

type block struct {
	d      *decode.D
	r      ranges.Range
	simple bool
}

type decodeContext struct {
	currentTrack *track
	tracks       []*track
	blocks       []block
}

func decodeMaster(d *decode.D, bitsLimit int64, tag ebml.Tag, dc *decodeContext) {
	tagEndBit := d.Pos() + bitsLimit

	d.FieldArray("elements", func(d *decode.D) {
		// var crcD *decode.D
		// var crcStart int64

		for d.Pos() < tagEndBit && d.NotEnd() {
			d.FieldStruct("element", func(d *decode.D) {
				a := ebml.Attribute{
					Type: ebml.Unknown,
				}

				tagID := d.FieldUFn("id", decodeRawVint, scalar.Fn(func(s scalar.S) (scalar.S, error) {
					n := s.ActualU()
					var ok bool
					a, ok = tag[n]
					if !ok {
						a, ok = ebml.Global[n]
						if !ok {
							a = ebml.Attribute{
								Type: ebml.Unknown,
							}
							return scalar.S{Actual: n, ActualDisplay: scalar.NumberHex, Description: "Unknown"}, nil
						}
					}
					return scalar.S{Actual: n, ActualDisplay: scalar.NumberHex, Sym: a.Name, Description: a.Definition}, nil
				}))
				d.FieldValueStr("type", ebml.TypeNames[a.Type])

				if tagID == ebml_matroska.TrackEntryID {
					dc.currentTrack = &track{}
					dc.tracks = append(dc.tracks, dc.currentTrack)
				}

				// tagSize could be 0xffffffffffffff which means "unknown" size, then we will read until eof
				// TODO: should read until unknown id:
				//    The end of a Master-element with unknown size is determined by the beginning of the next
				//    element that is not a valid sub-element of that Master-element
				// TODO: should also handle garbage between
				const maxStringTagSize = 100 * 1024 * 1024
				tagSize := d.FieldUFn("size", decodeVint)

				// assert sane tag size
				// TODO: strings are limited for now because they are read into memory
				switch a.Type {
				case ebml.Integer,
					ebml.Uinteger,
					ebml.Float:
					if tagSize > 8 {
						d.Fatalf("invalid tagSize %d for number type", tagSize)
					}
				case ebml.String,
					ebml.UTF8:
					if tagSize > maxStringTagSize {
						d.Errorf("tagSize %d > maxStringTagSize %d", tagSize, maxStringTagSize)
					}
				case ebml.Unknown,
					ebml.Binary,
					ebml.Date,
					ebml.Master:
					// nop
				}

				optionalMap := func(sm scalar.Mapper) scalar.Mapper {
					return scalar.Fn(func(s scalar.S) (scalar.S, error) {
						if sm != nil {
							return sm.MapScalar(s)
						}
						return s, nil
					})
				}

				switch a.Type {
				case ebml.Unknown:
					d.FieldRawLen("data", int64(tagSize)*8)
				case ebml.Integer:
					d.FieldS("value", int(tagSize)*8, optionalMap(a.IntegerEnums))
				case ebml.Uinteger:
					v := d.FieldU("value", int(tagSize)*8, optionalMap(a.UintegerEnums))
					if dc.currentTrack != nil && tagID == ebml_matroska.TrackNumberID {
						dc.currentTrack.number = int(v)
					}
				case ebml.Float:
					d.FieldF("value", int(tagSize)*8)
				case ebml.String:
					v := d.FieldUTF8("value", int(tagSize), optionalMap(a.StringEnums))
					if dc.currentTrack != nil && tagID == ebml_matroska.CodecIDID {
						dc.currentTrack.codec = v
					}
				case ebml.UTF8:
					d.FieldUTF8NullFixedLen("value", int(tagSize))
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
					d.FieldRawLen("value", int64(tagSize)*8)
				case ebml.Binary:
					switch tagID {
					case ebml_matroska.SimpleBlockID:
						dc.blocks = append(dc.blocks, block{
							d:      d,
							r:      ranges.Range{Start: d.Pos(), Len: int64(tagSize) * 8},
							simple: true,
						})
						d.SeekRel(int64(tagSize) * 8)
					case ebml_matroska.BlockID:
						dc.blocks = append(dc.blocks, block{
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
					case ebml_matroska.FileDataID:
						d.FieldFormatLen("value", int64(tagSize)*8, imageFormat, nil)
					default:
						d.FieldRawLen("value", int64(tagSize)*8)
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
		// 	ioex.MustCopy(elementCRC, d.BitBufRange(crcStart, d.Pos()-crcStart))
		// 	crcD.FieldChecksumRange("value", crcValue.Range.Start, crcValue.Range.Len, elementCRC.Sum(nil), decode.LittleEndian)
		// }
	})

}

func matroskaDecode(d *decode.D, _ any) any {
	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekBits(32) != ebmlHeaderID {
		d.Errorf("no EBML header found")
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

		// TODO: refactor, one DecodeRangeFn or refactor to add FieldFormatRangeFn?

		switch t.codec {
		case "A_VORBIS":
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				numPackets := d.FieldU8("num_packets")
				// TODO: lacing
				packetLengths := []int64{}
				// Xiph-style lacing (similar to ogg) of n-1 packets, last is reset of block
				d.FieldArray("laces", func(d *decode.D) {
					for i := uint64(0); i < numPackets; i++ {
						l := d.FieldUFn("lace", func(d *decode.D) uint64 {
							var l uint64
							for {
								n := d.U8()
								l += n
								if n < 255 {
									return l
								}
							}
						})
						packetLengths = append(packetLengths, int64(l))
					}
				})
				d.FieldArray("packets", func(d *decode.D) {
					for _, l := range packetLengths {
						d.FieldFormatLen("packet", l*8, vorbisPacketFormat, nil)
					}
					d.FieldFormatLen("packet", d.BitsLeft(), vorbisPacketFormat, nil)
				})
			})
		case "A_AAC":
			dv, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegASCFrameFormat, nil)
			mpegASCOut, ok := v.(format.MPEGASCOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected mpegASCOut got %#+v", v))
			}
			//nolint:gosimple
			t.formatInArg = format.AACFrameIn{ObjectType: mpegASCOut.ObjectType}
		case "A_OPUS":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, opusPacketFrameFormat, nil)
		case "A_FLAC":
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldStruct("value", func(d *decode.D) {
					d.FieldUTF8("magic", 4, d.AssertStr("fLaC"))
					dv, v := d.FieldFormat("metadatablocks", flacMetadatablocksFormat, nil)
					flacMetadatablockOut, ok := v.(format.FlacMetadatablocksOut)
					if dv != nil && !ok {
						panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
					}
					if flacMetadatablockOut.HasStreamInfo {
						t.formatInArg = format.FlacFrameIn{BitsPerSample: int(flacMetadatablockOut.StreamInfo.BitsPerSample)}
					}
				})
			})
		case "V_MPEG4/ISO/AVC":
			dv, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegAVCDCRFormat, nil)
			avcDcrOut, ok := v.(format.AvcDcrOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			t.formatInArg = format.AvcAuIn{LengthSize: avcDcrOut.LengthSize} //nolint:gosimple
		case "V_MPEGH/ISO/HEVC":
			dv, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegHEVCDCRFormat, nil)
			hevcDcrOut, ok := v.(format.HevcDcrOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
			}
			t.formatInArg = format.HevcAuIn{LengthSize: hevcDcrOut.LengthSize} //nolint:gosimple
		case "V_AV1":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, av1CCRFormat, nil)
		case "V_VP9":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, vp9CFMFormat, nil)
		default:
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldRawLen("value", d.BitsLeft())
			})
		}
	}

	for _, b := range dc.blocks {
		b.d.RangeFn(b.r.Start, b.r.Len, func(d *decode.D) {
			trackNumber := d.FieldUFn("track_number", decodeVint)
			d.FieldU16("timestamp")
			if b.simple {
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldBool("key_frame")
					d.FieldU3("reserved")
					d.FieldBool("invisible")
					d.FieldU2("lacing")
					d.FieldBool("discardable")
				})
			} else {
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldU4("reserved")
					d.FieldBool("invisible")
					d.FieldU2("lacing")
					d.FieldBool("not_used")
				})
			}
			// TODO: lacing etc

			// TODO: fixed/unknown?
			if t, ok := trackNumberToTrack[int(trackNumber)]; ok {
				if f, ok := codecToFormat[t.codec]; ok {
					d.FieldFormat("packet", *f, t.formatInArg)
				}
			}

			if d.BitsLeft() > 0 {
				d.FieldRawLen("data", d.BitsLeft())
			}
		})
	}

	return nil
}
