package matroska

// https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
// https://matroska.org/technical/specs/index.html
// https://www.matroska.org/technical/basics.html
// https://www.matroska.org/technical/codec_specs.html
// https://wiki.xiph.org/MatroskaOpus

// TODO: refactor simepleblock/block to just defer decode etc?
// TODO: CRC
// TODO: handle garbage (see tcl and example files)
// TODO: could use md5 here somehow, see flac.go

import (
	"embed"
	"fmt"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/matroska/ebml"
	"github.com/wader/fq/format/matroska/ebml_matroska"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed matroska.jq
//go:embed matroska.md
var matroskaFS embed.FS

var aacFrameGroup decode.Group
var av1CCRGroup decode.Group
var av1FrameGroup decode.Group
var flacFrameGroup decode.Group
var flacMetadatablocksGroup decode.Group
var imageGroup decode.Group
var mp3FrameGroup decode.Group
var mpegASCFrameGroup decode.Group
var mpegAVCAUGroup decode.Group
var mpegAVCDCRGroup decode.Group
var mpegHEVCDCRGroup decode.Group
var mpegHEVCSampleGroup decode.Group
var mpegPESPacketSampleGroup decode.Group
var mpegSPUFrameGroup decode.Group
var opusPacketFrameGroup decode.Group
var vorbisPacketGroup decode.Group
var vp8FrameGroup decode.Group
var vp9CFMGroup decode.Group
var vp9FrameGroup decode.Group

var codecToGroup map[string]*decode.Group

func init() {
	interp.RegisterFormat(
		format.Matroska,
		&decode.Format{
			Description: "Matroska file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    matroskaDecode,
			DefaultInArg: format.Matroska_In{
				DecodeSamples: true,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AAC_Frame}, Out: &aacFrameGroup},
				{Groups: []*decode.Group{format.AV1_CCR}, Out: &av1CCRGroup},
				{Groups: []*decode.Group{format.AV1_Frame}, Out: &av1FrameGroup},
				{Groups: []*decode.Group{format.AVC_AU}, Out: &mpegAVCAUGroup},
				{Groups: []*decode.Group{format.AVC_DCR}, Out: &mpegAVCDCRGroup},
				{Groups: []*decode.Group{format.FLAC_Frame}, Out: &flacFrameGroup},
				{Groups: []*decode.Group{format.FLAC_Metadatablocks}, Out: &flacMetadatablocksGroup},
				{Groups: []*decode.Group{format.HEVC_AU}, Out: &mpegHEVCSampleGroup},
				{Groups: []*decode.Group{format.HEVC_DCR}, Out: &mpegHEVCDCRGroup},
				{Groups: []*decode.Group{format.Image}, Out: &imageGroup},
				{Groups: []*decode.Group{format.MP3_Frame}, Out: &mp3FrameGroup},
				{Groups: []*decode.Group{format.MPEG_ASC}, Out: &mpegASCFrameGroup},
				{Groups: []*decode.Group{format.MPEG_PES_Packet}, Out: &mpegPESPacketSampleGroup},
				{Groups: []*decode.Group{format.MPEG_SPU}, Out: &mpegSPUFrameGroup},
				{Groups: []*decode.Group{format.Opus_Packet}, Out: &opusPacketFrameGroup},
				{Groups: []*decode.Group{format.Vorbis_Packet}, Out: &vorbisPacketGroup},
				{Groups: []*decode.Group{format.VP8_Frame}, Out: &vp8FrameGroup},
				{Groups: []*decode.Group{format.VP9_CFM}, Out: &vp9CFMGroup},
				{Groups: []*decode.Group{format.VP9_Frame}, Out: &vp9FrameGroup},
			},
		})
	interp.RegisterFS(matroskaFS)

	codecToGroup = map[string]*decode.Group{
		"A_VORBIS":         &vorbisPacketGroup,
		"A_MPEG/L3":        &mp3FrameGroup,
		"A_FLAC":           &flacFrameGroup,
		"A_AAC":            &aacFrameGroup,
		"A_OPUS":           &opusPacketFrameGroup,
		"V_VP8":            &vp8FrameGroup,
		"V_VP9":            &vp9FrameGroup,
		"V_AV1":            &av1FrameGroup,
		"V_VOBSUB":         &mpegSPUFrameGroup,
		"V_MPEG4/ISO/AVC":  &mpegAVCAUGroup,
		"V_MPEGH/ISO/HEVC": &mpegHEVCSampleGroup,
		"V_MPEG2":          &mpegPESPacketSampleGroup,
		"S_VOBSUB":         &mpegSPUFrameGroup,
	}
}

const (
	lacingTypeNone  = 0b00
	lacingTypeXiph  = 0b01
	lacingTypeFixed = 0b10
	lacingTypeEBML  = 0b11
)

var lacingTypeNames = scalar.UintMapSymStr{
	lacingTypeNone:  "none",
	lacingTypeXiph:  "xiph",
	lacingTypeFixed: "fixed",
	lacingTypeEBML:  "ebml",
}

const tagSizeUnknown = 0xffffffffffffff

var sintActualMatroskaEpochDescription = scalar.SintActualDateDescription(ebml.EpochDate, time.Nanosecond, time.RFC3339)

func decodeLacingFn(d *decode.D, lacingType int, fn func(d *decode.D)) {
	if lacingType == lacingTypeNone {
		fn(d)
		return
	}

	// -1 size means rest of buffer, used last sometimes
	var laceSizes []int64

	switch lacingType {
	case lacingTypeXiph:
		numLaces := int(d.FieldU8("num_laces"))
		d.FieldArray("lace_sizes", func(d *decode.D) {
			for i := 0; i < numLaces; i++ {
				s := int64(d.FieldUintFn("lace_size", decodeXiphLaceSize))
				laceSizes = append(laceSizes, s)
			}
			laceSizes = append(laceSizes, -1)
		})
	case lacingTypeEBML:
		numLaces := int(d.FieldU8("num_laces"))
		d.FieldArray("lace_sizes", func(d *decode.D) {
			s := int64(d.FieldUintFn("lace_size", decodeVint)) // first is unsigned, not ranged shifted
			laceSizes = append(laceSizes, s)
			for i := 0; i < numLaces-1; i++ {
				d := int64(d.FieldUintFn("lace_size_delta", decodeRawVint))
				// range shifting
				switch {
				case d&0b1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1000_0000 == 0b0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_1000_0000:
					// value -(2^6^-1) to 2^6^-1 (ie 0_to 2^7^-2 minus 2^6^-1, half of the range)
					d -= 0b0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_1011_1111
				case d&0b1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1100_0000_0000_0000 == 0b0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0100_0000_0000_0000:
					// value -(2^13^-1) to 2^13^-1
					d -= 0b0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0101_1111_1111_1111
				case d&0b1111_1111_1111_1111_1111_1111_1111_1111_1110_0000_0000_0000_0000_0000 == 0b0000_0000_0000_0000_0000_0000_0000_0000_0010_0000_0000_0000_0000_0000:
					// value -(2^20^-1) to 2^20^-1
					d -= 0b0000_0000_0000_0000_0000_0000_0000_0000_0010_1111_1111_1111_1111_1111
				case d&0b1111_1111_1111_1111_1111_1111_1111_0000_0000_0000_0000_0000_0000_0000 == 0b0000_0000_0000_0000_0000_0000_0001_0000_0000_0000_0000_0000_0000_0000:
					// value -(2^27^-1) to 2^27^-1
					d -= 0b0000_0000_0000_0000_0000_0000_0001_0111_1111_1111_1111_1111_1111_1111
				case d&0b1111_1111_1111_1111_1111_1000_0000_0000_0000_0000_0000_0000_0000_0000 == 0b0000_0000_0000_0000_0000_1000_0000_0000_0000_0000_0000_0000_0000_0000:
					// value -(2^34^-1) to 2^34^-1
					d -= 0b0000_0000_0000_0000_0000_1011_1111_1111_1111_1111_1111_1111_1111_1111
				case d&0b1111_1111_1111_1100_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000 == 0b0000_0000_0000_0100_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000:
					// value -(2^41^-1) to 2^41^-1
					d -= 0b0000_0000_0000_0101_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111
				case d&0b1111_1110_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000 == 0b0000_0010_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000:
					// value -(2^48^-1) to 2^48^-1
					d -= 0b0000_0010_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111
				}
				s += d
				laceSizes = append(laceSizes, s)
			}
			laceSizes = append(laceSizes, -1)
		})
	case lacingTypeFixed:
		numLaces := int(d.FieldU8("num_laces"))
		fixedSize := (d.BitsLeft() / 8) / int64(numLaces+1)
		for i := 0; i < numLaces+1; i++ {
			laceSizes = append(laceSizes, fixedSize)
		}
	default:
		panic("unreachable")
	}

	d.FieldArray("laces", func(d *decode.D) {
		for _, laceSize := range laceSizes {
			s := laceSize * 8
			if laceSize == -1 {
				s = d.BitsLeft()
			}
			d.FramedFn(s, fn)
		}
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

func peekRawVint(d *decode.D) uint64 {
	n, w := decodeRawVintWidth(d)
	d.SeekRel(int64(-w) * 8)
	return n
}

func decodeVint(d *decode.D) uint64 {
	n, w := decodeRawVintWidth(d)
	m := (uint64(1<<((w-1)*8+(8-w))) - 1)
	return n & m
}

func decodeXiphLaceSize(d *decode.D) uint64 {
	var s uint64
	for {
		n := d.U8()
		s += n
		if n < 255 {
			return s
		}
	}
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

func decodeMaster(d *decode.D, bitsLimit int64, elm *ebml.Master, unknownSize bool, dc *decodeContext) {
	tagEndBit := d.Pos() + bitsLimit

	d.FieldArray("elements", func(d *decode.D) {
		for d.Pos() < tagEndBit && !d.End() {
			// current we assume master with unknown size has ended if a valid parent is found
			// TODO:
			// https://github.com/ietf-wg-cellar/ebml-specification/blob/master/specification.markdown#unknown-data-size
			// > Any valid EBML Element according to the EBML Schema, Global Elements excluded, that
			// > is not a Descendant Element of the Unknown-Sized Element but shares a common direct
			// > parent, such as a Top-Level Element.
			// TODO: What to do if peeked is unknown?
			// TODO: Handle garbage between element
			if unknownSize {
				peekTagID := peekRawVint(d)
				_, validParent := ebml.FindParentID(ebml_matroska.IDToElement, elm.GetID(), ebml.ID(peekTagID))
				if validParent {
					break
				}
			}

			d.FieldStruct("element", func(d *decode.D) {
				var childElm ebml.Element
				childElm = &ebml.Unknown{}

				tagID := d.FieldUintFn("id", decodeRawVint, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
					n := s.Actual
					var ok bool
					childElm, ok = elm.Master[ebml.ID(n)]
					if !ok {
						childElm, ok = ebml.Global.Master[ebml.ID(n)]
						if !ok {
							childElm = &ebml.Unknown{}
							return scalar.Uint{Actual: n, DisplayFormat: scalar.NumberHex, Description: "Unknown"}, nil
						}
					}
					return scalar.Uint{
						Actual:        n,
						DisplayFormat: scalar.NumberHex,
						Sym:           childElm.GetName(),
						Description:   childElm.GetDefinition(),
					}, nil
				}))
				d.FieldValueStr("type", childElm.GetType())

				if tagID == ebml_matroska.TrackEntryID {
					dc.currentTrack = &track{}
					dc.tracks = append(dc.tracks, dc.currentTrack)
				}

				const maxStringTagSize = 100 * 1024 * 1024
				tagSize := d.FieldUintFn("size", decodeVint, scalar.UintMapDescription{
					0xffffffffffffff: "Unknown size",
				})
				unknownSize := tagSize == tagSizeUnknown
				if unknownSize {
					tagSize = uint64(d.BitsLeft() / 8)
				}

				// assert sane tag size
				// TODO: strings are limited for now because they are read into memory
				switch childElm.(type) {
				case *ebml.Integer,
					*ebml.Uinteger,
					*ebml.Float:
					if tagSize > 8 {
						d.Fatalf("invalid tagSize %d for number type", tagSize)
					}
				case *ebml.String,
					*ebml.UTF8:
					if tagSize > maxStringTagSize {
						d.Errorf("tagSize %d > maxStringTagSize %d", tagSize, maxStringTagSize)
					}
				case *ebml.Unknown,
					*ebml.Binary,
					*ebml.Date,
					*ebml.Master:
					// nop
				}

				switch childElm := childElm.(type) {
				case *ebml.Unknown:
					d.FieldRawLen("data", int64(tagSize)*8)
				case *ebml.Integer:
					var sm []scalar.SintMapper
					if childElm.Enums != nil {
						sm = append(sm, scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
							if e, ok := childElm.Enums[s.Actual]; ok {
								s.Sym = e.Name
								s.Description = e.Description
							}
							return s, nil
						}))
					}
					d.FieldS("value", int(tagSize)*8, sm...)
				case *ebml.Uinteger:
					var sm []scalar.UintMapper
					if childElm.Enums != nil {
						sm = append(sm, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
							if e, ok := childElm.Enums[s.Actual]; ok {
								s.Sym = e.Name
								s.Description = e.Description
							}
							return s, nil
						}))
					}
					v := d.FieldU("value", int(tagSize)*8, sm...)
					if dc.currentTrack != nil && tagID == ebml_matroska.TrackNumberID {
						dc.currentTrack.number = int(v)
					}
				case *ebml.Float:
					d.FieldF("value", int(tagSize)*8)
				case *ebml.String:
					var sm []scalar.StrMapper
					sm = append(sm, scalar.StrFn(func(s scalar.Str) (scalar.Str, error) {
						if e, ok := childElm.Enums[s.Actual]; ok {
							s.Sym = e.Name
							s.Description = e.Description
						}
						return s, nil
					}))
					v := d.FieldUTF8("value", int(tagSize), sm...)
					if dc.currentTrack != nil && tagID == ebml_matroska.CodecIDID {
						dc.currentTrack.codec = v
					}
				case *ebml.UTF8:
					d.FieldUTF8NullFixedLen("value", int(tagSize))
				case *ebml.Date:
					switch tagSize {
					case 0:
						d.FieldValueSint("value", 0, sintActualMatroskaEpochDescription)
					case 8:
						d.FieldS("value", int(tagSize)*8, sintActualMatroskaEpochDescription)
					default:
						d.FieldRawLen("value", int64(tagSize)*8)
					}
				case *ebml.Binary:
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
						d.FieldFormatOrRawLen("value", int64(tagSize)*8, &imageGroup, nil)
					default:
						d.FieldRawLen("value", int64(tagSize)*8)
					}

				case *ebml.Master:
					decodeMaster(d, int64(tagSize)*8, childElm, unknownSize, dc)
				}
			})
		}
	})
}

func matroskaDecode(d *decode.D) any {
	var mi format.Matroska_In
	d.ArgAs(&mi)

	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekUintBits(32) != ebmlHeaderID {
		d.Errorf("no EBML header found")
	}
	dc := &decodeContext{tracks: []*track{}}
	decodeMaster(d, d.BitsLeft(), ebml_matroska.RootElement, false, dc)

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
				decodeLacingFn(d, lacingTypeXiph, func(d *decode.D) {
					if mi.DecodeSamples {
						d.FieldFormat("packet", &vorbisPacketGroup, nil)
					} else {
						d.FieldRawLen("packet", d.BitsLeft())
					}
				})
			})
		case "A_AAC":
			_, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &mpegASCFrameGroup, nil)
			mpegASCOut, ok := v.(format.MPEG_ASC_Out)
			if !ok {
				panic(fmt.Sprintf("expected mpegASCOut got %#+v", v))
			}
			t.formatInArg = format.AAC_Frame_In(mpegASCOut)
		case "A_OPUS":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &opusPacketFrameGroup, nil)
		case "A_FLAC":
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldStruct("value", func(d *decode.D) {
					d.FieldUTF8("magic", 4, d.StrAssert("fLaC"))
					_, v := d.FieldFormat("metadatablocks", &flacMetadatablocksGroup, nil)
					flacMetadatablockOut, ok := v.(format.FLAC_Metadatablocks_Out)
					if !ok {
						panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
					}
					if flacMetadatablockOut.HasStreamInfo {
						t.formatInArg = format.FLAC_Frame_In{BitsPerSample: int(flacMetadatablockOut.StreamInfo.BitsPerSample)}
					}
				})
			})
		case "V_MPEG4/ISO/AVC":
			_, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &mpegAVCDCRGroup, nil)
			avcDcrOut, ok := v.(format.AVC_DCR_Out)
			if !ok {
				panic(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			t.formatInArg = format.AVC_AU_In(avcDcrOut)
		case "V_MPEGH/ISO/HEVC":
			_, v := t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &mpegHEVCDCRGroup, nil)
			hevcDcrOut, ok := v.(format.HEVC_DCR_Out)
			if !ok {
				panic(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
			}
			t.formatInArg = format.HEVC_AU_In(hevcDcrOut)
		case "V_AV1":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &av1CCRGroup, nil)
		case "V_VP9":
			t.parentD.FieldFormatRange("value", t.codecPrivatePos, t.codecPrivateTagSize, &vp9CFMGroup, nil)
		default:
			t.parentD.RangeFn(t.codecPrivatePos, t.codecPrivateTagSize, func(d *decode.D) {
				d.FieldRawLen("value", d.BitsLeft())
			})
		}
	}

	for _, b := range dc.blocks {
		b.d.RangeFn(b.r.Start, b.r.Len, func(d *decode.D) {
			var lacing uint64
			trackNumber := d.FieldUintFn("track_number", decodeVint)
			d.FieldU16("timestamp")
			if b.simple {
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldBool("key_frame")
					d.FieldU3("reserved")
					d.FieldBool("invisible")
					lacing = d.FieldU2("lacing", lacingTypeNames)
					d.FieldBool("discardable")
				})
			} else {
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldU4("reserved")
					d.FieldBool("invisible")
					lacing = d.FieldU2("lacing", lacingTypeNames)
					d.FieldBool("not_used")
				})
			}

			var f *decode.Group
			var track *track
			track, trackOk := trackNumberToTrack[int(trackNumber)]
			if trackOk {
				f = codecToGroup[track.codec]
			}

			decodeLacingFn(d, int(lacing), func(d *decode.D) {
				if mi.DecodeSamples && f != nil {
					d.FieldFormat("packet", f, track.formatInArg)
				} else {
					d.FieldRawLen("packet", d.BitsLeft())
				}
			})
		})
	}

	return nil
}
