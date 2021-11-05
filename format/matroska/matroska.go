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
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/ranges"
)

//go:embed *.jq
var matroskaFS embed.FS

var aacFrameFormat []*decode.Format
var av1CCRFormat []*decode.Format
var av1FrameFormat []*decode.Format
var flacFrameFormat []*decode.Format
var flacMetadatablocksFormat []*decode.Format
var imageFormat []*decode.Format
var mp3FrameFormat []*decode.Format
var mpegASCFrameFormat []*decode.Format
var mpegAVCAUFormat []*decode.Format
var mpegAVCDCRFormat []*decode.Format
var mpegHEVCDCRFormat []*decode.Format
var mpegHEVCSampleFormat []*decode.Format
var mpegPESPacketSampleFormat []*decode.Format
var mpegSPUFrameFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp8FrameFormat []*decode.Format
var vp9CFMFormat []*decode.Format
var vp9FrameFormat []*decode.Format

var codecToFormat map[string]*[]*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MATROSKA,
		Description: "Matroska file",
		Groups:      []string{format.PROBE},
		DecodeFn:    matroskaDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.AV1_CCR}, Formats: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Formats: &av1FrameFormat},
			{Names: []string{format.AVC_AU}, Formats: &mpegAVCAUFormat},
			{Names: []string{format.AVC_DCR}, Formats: &mpegAVCDCRFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCKS}, Formats: &flacMetadatablocksFormat},
			{Names: []string{format.HEVC_AU}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.HEVC_DCR}, Formats: &mpegHEVCDCRFormat},
			{Names: []string{format.IMAGE}, Formats: &imageFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3FrameFormat},
			{Names: []string{format.MPEG_ASC}, Formats: &mpegASCFrameFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Formats: &mpegPESPacketSampleFormat},
			{Names: []string{format.MPEG_SPU}, Formats: &mpegSPUFrameFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP8_FRAME}, Formats: &vp8FrameFormat},
			{Names: []string{format.VP9_CFM}, Formats: &vp9CFMFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
		},
		Files: matroskaFS,
	})

	codecToFormat = map[string]*[]*decode.Format{
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
	formatInArg         interface{}
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
				var a ebml.Attribute

				tagID := d.FieldUFn("id", decodeRawVint, func(s decode.Scalar) (decode.Scalar, error) {
					n := s.ActualU()
					var ok bool
					a, ok = tag[n]
					if !ok {
						a, ok = ebml.Global[n]
						if !ok {
							d.Invalid(fmt.Sprintf("unknown id %d", n))
						}
					}
					return decode.Scalar{Actual: n, DisplayFormat: decode.NumberHex, Sym: a.Name, Description: a.Definition}, nil
				})
				d.FieldValueU("type", uint64(a.Type), d.Sym(ebml.TypeNames[a.Type]))

				if tagID == ebml_matroska.TrackEntryID {
					dc.currentTrack = &track{}
					dc.tracks = append(dc.tracks, dc.currentTrack)
				}

				// tagSize could be 0xffffffffffffff which means "unknown" size, then we will read until eof
				// TODO: should read until unknown id:
				//    The end of a Master-element with unknown size is determined by the beginning of the next
				//    element that is not a valid sub-element of that Master-element
				// TODO: should also handle garbage between
				tagSize := d.FieldUFn("size", decodeVint)

				if tagSize > 8 &&
					(a.Type == ebml.Integer ||
						a.Type == ebml.Uinteger ||
						a.Type == ebml.Float) {
					d.Invalid(fmt.Sprintf("invalid tagSize %d for non-master type", tagSize))
				}

				switch a.Type {
				case ebml.Integer:
					d.FieldU("value", int(tagSize)*8, func(s decode.Scalar) (decode.Scalar, error) {
						if a.IntegerEnums != nil {
							if e, ok := a.IntegerEnums[s.ActualS()]; ok {
								s.Sym = e.Label
								s.Description = e.Definition
							}
						}
						return s, nil
					})

				case ebml.Uinteger:
					v := d.FieldU("value", int(tagSize)*8, func(s decode.Scalar) (decode.Scalar, error) {
						if a.UintegerEnums != nil {
							if e, ok := a.UintegerEnums[s.ActualU()]; ok {
								s.Sym = e.Label
								s.Description = e.Definition
							}
						}
						return s, nil
					})

					if dc.currentTrack != nil && tagID == ebml_matroska.TrackNumberID {
						dc.currentTrack.number = int(v)
					}
				case ebml.Float:
					d.FieldF("value", int(tagSize)*8)
				case ebml.String:
					v := d.FieldUTF8("value", int(tagSize), func(s decode.Scalar) (decode.Scalar, error) {
						if a.StringEnums != nil {
							if e, ok := a.StringEnums[s.ActualStr()]; ok {
								s.Sym = e.Label
								s.Description = e.Definition
							}
						}
						return s, nil
					})

					if dc.currentTrack != nil && tagID == ebml_matroska.CodecIDID {
						dc.currentTrack.codec = v
					}
				case ebml.UTF8:
					d.FieldUTF8NullTerminatedLen("value", int(tagSize))
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
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegASCFrameFormat, nil)
		case "A_OPUS":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, opusPacketFrameFormat, nil)
		case "A_FLAC":
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldStruct("value", func(d *decode.D) {
					d.FieldUTF8("magic", 4, d.AssertStr("fLaC"))
					_, v := d.FieldFormat("metadatablocks", flacMetadatablocksFormat, nil)
					flacMetadatablockOut, ok := v.(format.FlacMetadatablocksOut)
					if !ok {
						d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
					}
					if flacMetadatablockOut.HasStreamInfo {
						t.formatInArg = format.FlacFrameIn{StreamInfo: flacMetadatablockOut.StreamInfo}
					}
				})
			})
		case "V_MPEG4/ISO/AVC":
			_, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegAVCDCRFormat, nil)
			avcDcrOut, ok := v.(format.AvcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			t.formatInArg = format.AvcIn{LengthSize: avcDcrOut.LengthSize} //nolint:gosimple
		case "V_MPEGH/ISO/HEVC":
			_, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, mpegHEVCDCRFormat, nil)
			hevcDcrOut, ok := v.(format.HevcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
			}
			t.formatInArg = format.HevcIn{LengthSize: hevcDcrOut.LengthSize} //nolint:gosimple
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
