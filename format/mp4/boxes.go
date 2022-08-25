package mp4

// TODO: flags?

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: keep track of list of sampleSize/entries instead and change sample read code
const maxSampleEntryCount = 10_000_000

var boxAliases = map[string]string{
	"styp": "ftyp",
}

const (
	boxSizeRestOfFile   = 0
	boxSizeUse64bitSize = 1
)

var boxSizeNames = scalar.UToDescription{
	boxSizeRestOfFile:   "Rest of file",
	boxSizeUse64bitSize: "Use 64 bit size",
}

var mediaTimeNames = scalar.SToDescription{
	-1: "empty",
}

var subTypeNames = scalar.StrToDescription{
	"alis": "Alias Data",
	"camm": "Camera Metadata",
	"crsm": "Clock Reference",
	"data": "Data",
	"hint": "Hint Track",
	"ipsm": "IPMP",
	"m7sm": "MPEG-7 Stream",
	"mdir": "Metadata",
	"mdta": "Metadata Tags",
	"meta": "NRT Metadata",
	"mjsm": "MPEG-J",
	"nrtm": "Non-Real Time Metadata",
	"ocsm": "Object Content",
	"odsm": "Object Descriptor",
	"pict": "Picture",
	"priv": "Private",
	"psmd": "Panasonic Static Metadata",
	"sbtl": "Subtitle",
	"sdsm": "Scene Description",
	"soun": "Audio Track",
	"subp": "Subpicture",
	"text": "Text",
	"tmcd": "Time Code",
	"url ": "URL",
	"vide": "Video Track",
}

var dataFormatNames = scalar.StrToDescription{
	"apch": "Apple ProRes 422 High Quality",
	"apcn": "Apple ProRes 422 Standard Definition",
	"apcs": "Apple ProRes 422 LT",
	"apco": "Apple ProRes 422 Proxy",
	"ap4h": "Apple ProRes 4444",
	"fLaC": "Fres Lossless Audio Codec",
	"Opus": "Xiph Opus",
	"vp09": "VP9",
	"avc1": "Advanced Video Coding / H.264 / MPEG-4 Part 10",
	"hev1": "High Efficiency Video Coding / H.265 / MPEG-H Part 2",
	"hvc1": "High Efficiency Video Coding / H.265 / MPEG-H Part 2",
	"av01": "AV1",
	"mp4a": "MPEG Audio",
	"mp4v": "MPEG Video",
	"jpeg": "JPEG Image",
}

var (
	uuidIsmlManifestBytes = [16]byte{0xa5, 0xd4, 0x0b, 0x30, 0xe8, 0x14, 0x11, 0xdd, 0xba, 0x2f, 0x08, 0x00, 0x20, 0x0c, 0x9a, 0x66}
	uuidXmpBytes          = [16]byte{0xbe, 0x7a, 0xcf, 0xcb, 0x97, 0xa9, 0x42, 0xe8, 0x9c, 0x71, 0x99, 0x94, 0x91, 0xe3, 0xaf, 0xac}
	uuidSphericalBytes    = [16]byte{0xff, 0xcc, 0x82, 0x63, 0xf8, 0x55, 0x4a, 0x93, 0x88, 0x14, 0x58, 0x7a, 0x02, 0x52, 0x1f, 0xdd}
	uuidPspUsmtBytes      = [16]byte{0x55, 0x53, 0x4d, 0x54, 0x21, 0xd2, 0x4f, 0xce, 0xbb, 0x88, 0x69, 0x5c, 0xfa, 0xc9, 0xc7, 0x40}
	uuidTfxdBytes         = [16]byte{0x6d, 0x1d, 0x9b, 0x05, 0x42, 0xd5, 0x44, 0xe6, 0x80, 0xe2, 0x14, 0x1d, 0xaf, 0xf7, 0x57, 0xb2}
	uuidTfrfBytes         = [16]byte{0xd4, 0x80, 0x7e, 0xf2, 0xca, 0x39, 0x46, 0x95, 0x8e, 0x54, 0x26, 0xcb, 0x9e, 0x46, 0xa7, 0x9f}
	uuidProfBytes         = [16]byte{0x50, 0x52, 0x4f, 0x46, 0x21, 0xd2, 0x4f, 0xce, 0xbb, 0x88, 0x69, 0x5c, 0xfa, 0xc9, 0xc7, 0x40}
	uuidIpodBytes         = [16]byte{0x6b, 0x68, 0x40, 0xf2, 0x5f, 0x24, 0x4f, 0xc5, 0xba, 0x39, 0xa5, 0x1b, 0xcf, 0x03, 0x23, 0xf3}
)

var uuidNames = scalar.BytesToScalar{
	{Bytes: uuidIsmlManifestBytes[:], Scalar: scalar.S{Sym: "isml_manifest"}},
	{Bytes: uuidXmpBytes[:], Scalar: scalar.S{Sym: "xmp"}},
	{Bytes: uuidSphericalBytes[:], Scalar: scalar.S{Sym: "spherical"}},
	{Bytes: uuidPspUsmtBytes[:], Scalar: scalar.S{Sym: "psp_usmt"}},
	{Bytes: uuidTfxdBytes[:], Scalar: scalar.S{Sym: "tfxd"}},
	{Bytes: uuidTfrfBytes[:], Scalar: scalar.S{Sym: "tfrf"}},
	{Bytes: uuidProfBytes[:], Scalar: scalar.S{Sym: "prof"}},
	{Bytes: uuidIpodBytes[:], Scalar: scalar.S{Sym: "ipod"}},
}

// ISO 639-2/T language code 3 * 5bit packed uint + 1 zero bit
func decodeLang(d *decode.D) string {
	d.U1()
	return string([]byte{
		byte(d.U5()) + 0x60,
		byte(d.U5()) + 0x60,
		byte(d.U5()) + 0x60},
	)
}

// Quicktime time seconds in January 1, 1904 UTC
var quicktimeEpochDate = time.Date(1904, time.January, 4, 0, 0, 0, 0, time.UTC)
var quicktimeEpoch = scalar.DescriptionActualUTime(quicktimeEpochDate, time.RFC3339)

func decodeFieldMatrix(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		d.FieldFP32("a")
		d.FieldFP32("b")
		d.FieldFP("u", 32, 30)
		d.FieldFP32("c")
		d.FieldFP32("d")
		d.FieldFP("v", 32, 30)
		d.FieldFP32("x")
		d.FieldFP32("y")
		d.FieldFP("w", 32, 30)
	})
}

func decodeBoxWithParentData(ctx *decodeContext, d *decode.D, parentData any) {
	var typ string
	var dataSize uint64

	boxSize := d.FieldU32("size", boxSizeNames)
	typ = d.FieldUTF8("type", 4, boxDescriptions)

	switch boxSize {
	case boxSizeRestOfFile:
		dataSize = uint64(d.Len()-d.Pos()) / 8
	case boxSizeUse64bitSize:
		boxSize = d.FieldU64("size64")
		dataSize = boxSize - 16
	default:
		dataSize = boxSize - 8
	}

	if ctx.opts.AllowTruncated && dataSize > uint64(d.BitsLeft()/8) {
		dataSize = uint64(d.BitsLeft() / 8)

	}

	// TODO: not sure about this
	switch {
	case typ == "�too":
		typ = "_apple_list"
	case typ[0] == 0xa9:
		typ = "_apple_entry"
	}

	if a, ok := boxAliases[typ]; ok {
		typ = a
	}

	if parentData != nil {
		ctx.path[len(ctx.path)-1].data = parentData
	}
	ctx.path = append(ctx.path, pathEntry{typ: typ, data: parentData})

	if decodeFn, ok := boxDecoders[typ]; ok {
		d.FramedFn(int64(dataSize*8), func(d *decode.D) {
			decodeFn(ctx, d)
		})
	} else {
		d.FieldRawLen("data", int64(dataSize*8))
	}

	ctx.path = ctx.path[0 : len(ctx.path)-1]
}

func decodeBoxes(ctx *decodeContext, d *decode.D) {
	decodeBoxesWithParentData(ctx, d, nil)
}

func decodeBoxesWithParentData(ctx *decodeContext, d *decode.D, parentData any) {
	d.FieldStructArrayLoop("boxes", "box", func() bool { return d.BitsLeft() >= 8*8 }, func(d *decode.D) {
		decodeBoxWithParentData(ctx, d, parentData)
	})

	if d.BitsLeft() > 0 {
		// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
		if d.BitsLeft() >= 32 && d.PeekBits(32) == 0 {
			d.FieldU32("zero_terminator")
		}
		if d.BitsLeft() > 0 {
			d.FieldRawLen("padding", d.BitsLeft())
		}
	}
}

var boxDecoders map[string]func(ctx *decodeContext, d *decode.D)

type irefBox struct {
	version int
}

type trakBox struct {
	trackID int
}

type moofBox struct {
	offset int64
}

type trafBox struct {
	trackID        int
	baseDataOffset int64
	moof           *moof
}

func irefEntryDecode(ctx *decodeContext, d *decode.D) {
	irefBox, ok := ctx.parent().data.(*irefBox)
	if !ok {
		d.FieldRawLen("data", d.BitsLeft())
		return
	}

	idSize := 16
	if irefBox.version != 0 {
		idSize = 32
	}

	d.FieldU("from_id", idSize)
	count := d.FieldU16("count")

	d.FieldArray("ids", func(d *decode.D) {
		for i := uint64(0); i < count; i++ {
			d.FieldU("id", idSize)
		}
	})
}

func init() {
	boxDecoders = map[string]func(ctx *decodeContext, d *decode.D){
		"ftyp": func(_ *decodeContext, d *decode.D) {
			d.FieldUTF8("major_brand", 4)
			d.FieldU32("minor_version")
			numBrands := d.BitsLeft() / 8 / 4
			var i int64
			d.FieldArrayLoop("brands", func() bool { return i < numBrands }, func(d *decode.D) {
				d.FieldUTF8("brand", 4, brandDescriptions, scalar.ActualTrimSpace)
				i++
			})
		},
		"mvhd": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			switch version {
			case 0:
				d.FieldU32("creation_time", quicktimeEpoch)
				d.FieldU32("modification_time", quicktimeEpoch)
				d.FieldU32("time_scale")
				d.FieldU32("duration")
			case 1:
				d.FieldU64("creation_time", quicktimeEpoch)
				d.FieldU64("modification_time", quicktimeEpoch)
				d.FieldU32("time_scale")
				d.FieldU64("duration")
			default:
				return
			}
			d.FieldFP32("preferred_rate")
			d.FieldFP16("preferred_volume")
			d.FieldUTF8("reserved", 10)
			decodeFieldMatrix(d, "matrix_structure")
			d.FieldU32("preview_time")
			d.FieldU32("preview_duration")
			d.FieldU32("poster_time")
			d.FieldU32("selection_time")
			d.FieldU32("selection_duration")
			d.FieldU32("current_time")
			d.FieldU32("next_track_id")
		},
		"trak": func(ctx *decodeContext, d *decode.D) {
			decodeBoxesWithParentData(ctx, d, &trakBox{})
		},
		"edts": decodeBoxes,
		"elst": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < entryCount }, func(d *decode.D) {
				d.FieldS32("segment_duration")
				d.FieldSFn("media_time", func(d *decode.D) int64 {
					var t int64
					if version == 0 {
						t = d.S32()
					} else {
						t = d.S64()
					}
					return t
				}, mediaTimeNames)
				d.FieldFP32("media_rate")
				i++
			})
		},
		"tref": decodeBoxes,
		"tkhd": func(ctx *decodeContext, d *decode.D) {
			var trackID int
			version := d.FieldU8("version")
			d.FieldU24("flags")
			switch version {
			case 0:
				d.FieldU32("creation_time", quicktimeEpoch)
				d.FieldU32("modification_time", quicktimeEpoch)
				trackID = int(d.FieldU32("track_id"))
				d.FieldU32("reserved1")
				d.FieldU32("duration")
			case 1:
				d.FieldU64("creation_time", quicktimeEpoch)
				d.FieldU64("modification_time", quicktimeEpoch)
				trackID = int(d.FieldU32("track_id"))
				d.FieldU32("reserved1")
				d.FieldU64("duration")
			default:
				return
			}
			d.FieldRawLen("reserved2", 8*8)
			d.FieldU16("layer")
			d.FieldU16("alternate_group")
			d.FieldFP16("volume")
			d.FieldU16("reserved3")
			decodeFieldMatrix(d, "matrix_structure")
			d.FieldFP32("track_width")
			d.FieldFP32("track_height")

			if t := ctx.currentTrakBox(); t != nil {
				t.trackID = trackID
				_ = ctx.currentTrack()
			}
		},
		"mdia": decodeBoxes,
		"mdhd": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			// TODO: timestamps
			switch version {
			case 0:
				d.FieldU32("creation_time", quicktimeEpoch)
				d.FieldU32("modification_time", quicktimeEpoch)
				d.FieldU32("time_scale")
				d.FieldU32("duration")
			case 1:
				d.FieldU64("creation_time", quicktimeEpoch)
				d.FieldU64("modification_time", quicktimeEpoch)
				d.FieldU32("time_scale")
				d.FieldU64("duration")
			default:
				return
			}
			d.FieldStrFn("language", decodeLang)
			d.FieldU16("quality")
		},
		"vmhd": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU16("graphicsmode")
			d.FieldArray("opcolor", func(d *decode.D) {
				// TODO: is FP16?
				d.FieldU16("value")
				d.FieldU16("value")
				d.FieldU16("value")
			})
		},
		"hdlr": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8NullFixedLen("component_type", 4)
			subType := d.FieldUTF8("component_subtype", 4, subTypeNames, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("component_manufacturer", 4)
			d.FieldU32("component_flags")
			d.FieldU32("component_flags_mask")
			// TODO: sometimes has a length prefix byte, how to know?
			d.FieldUTF8NullFixedLen("component_name", int(d.BitsLeft()/8))

			if t := ctx.currentTrack(); t != nil {
				// component_type seems to be all zero sometimes so can't look for "mhlr"
				switch subType {
				case "vide", "soun":
					t.subType = subType
				}
			}
		},
		"minf": decodeBoxes,
		"dinf": decodeBoxes,
		"dref": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			d.FieldStructArrayLoop("boxes", "box", func() bool { return i < entryCount }, func(d *decode.D) {
				size := d.FieldU32("size")
				d.FieldUTF8("type", 4)
				d.FieldU8("version")
				d.FieldU24("flags")
				dataSize := size - 12
				d.FieldRawLen("data", int64(dataSize*8))
				i++
			})
		},
		"stbl": decodeBoxes,
		"stsd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			// note called "boxes" here instead of "sample_descriptions" and data format is named "type".
			// this is to make it easier to threat them as normal boxes
			d.FieldArrayLoop("boxes", func() bool { return i < entryCount }, func(d *decode.D) {
				d.FieldStruct("box", func(d *decode.D) {
					size := d.FieldU32("size")
					dataFormat := d.FieldUTF8("type", 4, dataFormatNames)
					subType := ""
					if t := ctx.currentTrack(); t != nil {
						t.sampleDescriptions = append(t.sampleDescriptions, sampleDescription{
							dataFormat: dataFormat,
						})
						subType = t.subType
					}

					d.FramedFn(int64(size-8)*8, func(d *decode.D) {
						d.FieldRawLen("reserved", 6*8)
						d.FieldU16("data_reference_index")

						switch subType {
						case "soun", "vide":

							version := d.FieldU16("version")
							d.FieldU16("revision_level")
							d.FieldU32("max_packet_size") // TODO: vendor for some subtype?

							switch subType {
							case "soun":
								// AudioSampleEntry
								// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap3/qtff3.html#//apple_ref/doc/uid/TP40000939-CH205-SW1
								switch version {
								case 0:
									d.FieldU16("num_audio_channels")
									d.FieldU16("sample_size")
									d.FieldU16("compression_id")
									d.FieldU16("packet_size")
									d.FieldFP32("sample_rate")
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								case 1:
									d.FieldU16("num_audio_channels")
									d.FieldU16("sample_size")
									d.FieldU16("compression_id")
									d.FieldU16("packet_size")
									d.FieldFP32("sample_rate")
									d.FieldU32("samples_per_packet")
									d.FieldU32("bytes_per_packet")
									d.FieldU32("bytes_per_frame")
									d.FieldU32("bytes_per_sample")
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								case 2:
									d.FieldU16("always_3")
									d.FieldU16("always_16")
									d.FieldU16("always_minus_2") // TODO: as in const -2?
									d.FieldU16("always_0")
									d.FieldU32("always_65536")
									d.FieldU32("size_of_struct_only")
									d.FieldF64("audio_sample_rate")
									d.FieldU32("num_audio_channels")
									d.FieldU32("always_7f000000")
									d.FieldU32("const_bits_per_channel")
									d.FieldU32("format_specific_flags")
									d.FieldU32("const_bytes_per_audio_packet")
									d.FieldU32("const_lpcm_frames_per_audio_packet")
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								default:
									d.FieldRawLen("data", d.BitsLeft())
								}
							case "vide":
								// VideoSampleEntry
								// TODO: version 0 and 1 same?
								switch version {
								case 0, 1:
									d.FieldU32("temporal_quality")
									d.FieldU32("spatial_quality")
									d.FieldU16("width")
									d.FieldU16("height")
									d.FieldFP32("horizontal_resolution")
									d.FieldFP32("vertical_resolution")
									d.FieldU32("data_size")
									d.FieldU16("frame_count")
									d.FieldUTF8ShortStringFixedLen("compressor_name", 32)
									d.FieldU16("depth")
									d.FieldS16("color_table_id")
									// TODO: if 0 decode ctab
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								default:
									d.FieldRawLen("data", d.BitsLeft())
								}
								// case "hint": TODO: Hint entry
							default:
								d.FieldRawLen("data", d.BitsLeft())
							}
						default:
							d.FieldRawLen("data", d.BitsLeft())
						}

					})
				})
				i++
			})
		},
		"avcC": func(ctx *decodeContext, d *decode.D) {
			dv, v := d.FieldFormat("descriptor", mpegAVCDCRFormat, nil)
			avcDcrOut, ok := v.(format.AvcDcrOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			if t := ctx.currentTrack(); t != nil {
				t.formatInArg = format.AvcAuIn{LengthSize: avcDcrOut.LengthSize} //nolint:gosimple
			}
		},
		"hvcC": func(ctx *decodeContext, d *decode.D) {
			dv, v := d.FieldFormat("descriptor", mpegHEVCDCRFrameFormat, nil)
			hevcDcrOut, ok := v.(format.HevcDcrOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
			}
			if t := ctx.currentTrack(); t != nil {
				t.formatInArg = format.HevcAuIn{LengthSize: hevcDcrOut.LengthSize} //nolint:gosimple
			}
		},
		"dfLa": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			dv, v := d.FieldFormat("descriptor", flacMetadatablocksFormat, nil)
			flacMetadatablockOut, ok := v.(format.FlacMetadatablocksOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
			}
			if flacMetadatablockOut.HasStreamInfo {
				if t := ctx.currentTrack(); t != nil {
					t.formatInArg = format.FlacFrameIn{BitsPerSample: int(flacMetadatablockOut.StreamInfo.BitsPerSample)}
				}
			}
		},
		"dOps": func(_ *decodeContext, d *decode.D) {
			d.FieldFormat("descriptor", opusPacketFrameFormat, nil)
		},
		"av1C": func(_ *decodeContext, d *decode.D) {
			d.FieldFormat("descriptor", av1CCRFormat, nil)
		},
		"vpcC": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldFormat("descriptor", vpxCCRFormat, nil)
		},
		"esds": func(ctx *decodeContext, d *decode.D) {
			d.FieldU32("version")

			dv, v := d.FieldFormat("descriptor", mpegESFormat, nil)
			mpegEsOut, ok := v.(format.MpegEsOut)
			if dv != nil && !ok {
				panic(fmt.Sprintf("expected mpegEsOut got %#+v", v))
			}

			if t := ctx.currentTrack(); t != nil && len(mpegEsOut.DecoderConfigs) > 0 {
				dc := mpegEsOut.DecoderConfigs[0]
				t.objectType = dc.ObjectType
				t.formatInArg = format.AACFrameIn{ObjectType: dc.ASCObjectType}
			}
		},
		"stts": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			numEntries := d.FieldU32("entry_count")
			var i uint64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("count")
				d.FieldU32("delta")
				i++
			})
		},
		"stsc": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < entryCount }, func(d *decode.D) {
				firstChunk := uint32(d.FieldU32("first_chunk"))
				samplesPerChunk := uint32(d.FieldU32("samples_per_chunk"))
				d.FieldU32("sample_description_id")

				if t := ctx.currentTrack(); t != nil {
					t.stsc = append(t.stsc, stsc{
						firstChunk:      int(firstChunk),
						samplesPerChunk: int(samplesPerChunk),
					})
				}
				i++
			})
		},
		"stsz": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			// TODO: bytes_per_sample from audio stsd?
			sampleSize := d.FieldU32("sample_size")
			entryCount := d.FieldU32("entry_count")

			t := ctx.currentTrack()

			if t != nil && len(t.stsz) > 0 {
				d.Errorf("multiple stsz or stz2 boxes")
			}
			if sampleSize == 0 {
				var i uint64
				d.FieldArrayLoop("entries", func() bool { return i < entryCount }, func(d *decode.D) {
					size := uint32(d.FieldU32("size"))
					if t != nil {
						t.stsz = append(t.stsz, stsz{
							size:  int64(size),
							count: 1,
						})
					}
					i++
				})
			} else {
				if t != nil {
					t.stsz = append(t.stsz, stsz{
						size:  int64(sampleSize),
						count: int(entryCount),
					})
				}
			}
		},
		"stz2": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			fieldSize := d.FieldU32("field_size")
			if fieldSize > 16 {
				d.Errorf("field_size %d > 16", fieldSize)
			}
			entryCount := d.FieldU32("entry_count")
			var i uint64
			t := ctx.currentTrack()
			if t != nil && len(t.stsz) > 0 {
				d.Errorf("multiple stsz or stz2 boxes")
			}
			d.FieldArrayLoop("entries", func() bool { return i < entryCount }, func(d *decode.D) {
				size := uint32(d.FieldU("size", int(fieldSize)))
				if t != nil {
					t.stsz = append(t.stsz, stsz{
						size:  int64(size),
						count: 1,
					})
				}
				i++
			})
		},
		"stco": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			t := ctx.currentTrack()
			d.FieldArrayLoop("entries", func() bool { return i < entryCount }, func(d *decode.D) {
				chunkOffset := d.FieldU32("chunk_offset")
				if t != nil {
					t.stco = append(t.stco, int64(chunkOffset))
				}
				i++
			})
		},
		"stss": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldU32("sample_number")
				}
			})
		},
		"sdtp": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			// TODO: should be count from stsz
			// TODO: can we know count here or do we need to defer decoding somehow?
			d.FieldArray("entries", func(d *decode.D) {
				for d.NotEnd() {
					d.FieldStruct("entry", func(d *decode.D) {
						d.FieldU2("reserved")
						values := scalar.UToSymStr{
							0: "unknown",
							1: "yes",
							2: "no",
						}
						d.FieldU2("sample_depends_on", values)
						d.FieldU2("sample_is_depended_on", values)
						d.FieldU2("sample_has_redundancy", values)
					})
				}
			})
		},
		"ctts": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < entryCount }, func(d *decode.D) {
				d.FieldS32("sample_count")
				d.FieldS32("sample_offset")
				i++
			})
		},
		// TODO: refactor: merge with stco?
		"co64": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			var i uint64
			t := ctx.currentTrack()
			d.FieldArrayLoop("entries", func() bool { return i < entryCount }, func(d *decode.D) {
				offset := d.FieldU64("offset")
				if t != nil {
					t.stco = append(t.stco, int64(offset))
				}
				i++
			})
		},
		"sidx": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("reference_id")
			d.FieldU32("timescale")
			if version == 0 {
				d.FieldU32("pts")
				d.FieldU32("offset")
			} else {
				d.FieldU64("pts")
				d.FieldU64("offset")
			}
			d.FieldU16("reserved")
			numEntries := d.FieldU16("entry_count")
			var i uint64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU1("reference_type")
				d.FieldU31("size")
				d.FieldU32("duration")
				d.FieldU1("starts_with_sap")
				d.FieldU3("sap_type")
				d.FieldU28("sap_delta_time")
				i++
			})
		},
		"udta": decodeBoxes,
		"meta": func(ctx *decodeContext, d *decode.D) {
			// TODO: meta box sometimes has a 4 byte unknown field? (flag/version?)
			maybeFlags := d.PeekBits(32)
			if maybeFlags == 0 {
				// TODO: rename?
				d.FieldU32("maybe_flags")
			}
			decodeBoxes(ctx, d)
		},
		"ilst":        decodeBoxes,
		"_apple_list": decodeBoxes,
		"_apple_entry": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8("data", int(d.BitsLeft()/8))
		},
		"data": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("reserved")
			if ctx.isParent("covr") {
				d.FieldFormatOrRawLen("data", d.BitsLeft(), imageFormat, nil)
			} else {
				d.FieldUTF8("data", int(d.BitsLeft()/8))
			}
		},
		"moov": decodeBoxes,
		"moof": func(ctx *decodeContext, d *decode.D) {
			offset := (d.Pos() / 8) - 8
			decodeBoxesWithParentData(ctx, d, &moofBox{offset: offset})
		},
		// Track Fragment
		"traf": func(ctx *decodeContext, d *decode.D) {
			decodeBoxesWithParentData(ctx, d, &trafBox{})
		},
		// Movie Fragment Header
		"mfhd": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("sequence_number")
		},
		// Track Fragment Header
		"tfhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			baseDataOffsetPresent := false
			sampleDescriptionIndexPresent := false
			defaultSampleDurationPresent := false
			defaultSampleSizePresent := false
			defaultSampleFlagsPresent := false
			d.FieldStruct("flags", func(d *decode.D) {
				d.FieldU7("unused0")
				d.FieldBool("duration_is_empty")
				d.FieldU10("unused1")
				defaultSampleFlagsPresent = d.FieldBool("default_sample_flags_present")
				defaultSampleSizePresent = d.FieldBool("default_sample_size_present")
				defaultSampleDurationPresent = d.FieldBool("default_sample_duration_present")
				d.FieldU1("unused2")
				sampleDescriptionIndexPresent = d.FieldBool("sample_description_index_present")
				baseDataOffsetPresent = d.FieldBool("base_data_offset_present")
			})
			trackID := int(d.FieldU32("track_id"))

			m := &moof{}
			if mb := ctx.currentMoofBox(); mb != nil {
				m.offset = mb.offset
			}

			baseDataOffset := int64(0)
			if baseDataOffsetPresent {
				baseDataOffset = int64(d.FieldU64("base_data_offset"))
			}
			if sampleDescriptionIndexPresent {
				m.defaultSampleDescriptionIndex = int(d.FieldU32("sample_description_index"))
			}
			if defaultSampleDurationPresent {
				d.FieldU32("default_sample_duration")
			}
			if defaultSampleSizePresent {
				m.defaultSampleSize = int64(d.FieldU32("default_sample_size"))
			}
			if defaultSampleFlagsPresent {
				d.FieldU32("default_sample_flags")
			}

			if t := ctx.currentTrafBox(); t != nil {
				t.trackID = trackID
				t.moof = m
				t.baseDataOffset = baseDataOffset
			}
			if t := ctx.currentTrack(); t != nil {
				t.moofs = append(t.moofs, m)
			}
		},
		// Track Fragment Run
		"trun": func(ctx *decodeContext, d *decode.D) {
			m := &moof{}
			if t := ctx.currentTrafBox(); t != nil {
				m = t.moof
			}

			d.FieldU8("version")
			sampleCompositionTimeOffsetsPresent := false
			sampleFlagsPresent := false
			sampleSizePresent := false
			sampleDurationPresent := false
			firstSampleFlagsPresent := false
			dataOffsetPresent := false
			d.FieldStruct("flags", func(d *decode.D) {
				d.FieldU12("unused0")
				sampleCompositionTimeOffsetsPresent = d.FieldBool("sample_composition_time_offsets_present")
				sampleFlagsPresent = d.FieldBool("sample_flags_present")
				sampleSizePresent = d.FieldBool("sample_size_present")
				sampleDurationPresent = d.FieldBool("sample_duration_present")
				d.FieldU5("unused1")
				firstSampleFlagsPresent = d.FieldBool("first_sample_flags_present")
				d.FieldU1("unused2")
				dataOffsetPresent = d.FieldBool("data_offset_present")
			})
			sampleCount := d.FieldU32("sample_count")
			dataOffset := int64(0)
			if dataOffsetPresent {
				dataOffset = d.FieldS32("data_offset")
			}
			if firstSampleFlagsPresent {
				d.FieldU32("first_sample_flags")
			}

			if sampleCount > maxSampleEntryCount {
				d.Errorf("too many sample trun entries %d > %d", sampleCount, maxSampleEntryCount)
			}

			t := trun{
				dataOffset: dataOffset,
			}
			d.FieldArray("samples", func(d *decode.D) {
				for i := uint64(0); i < sampleCount; i++ {
					sampleSize := m.defaultSampleSize
					d.FieldStruct("sample", func(d *decode.D) {
						if sampleDurationPresent {
							d.FieldU32("sample_duration")
						}
						if sampleSizePresent {
							sampleSize = int64(d.FieldU32("sample_size"))
						}
						if sampleFlagsPresent {
							d.FieldU32("sample_flags")
						}
						if sampleCompositionTimeOffsetsPresent {
							d.FieldU32("sample_composition_time_offset")
						}
					})

					t.samplesSizes = append(t.samplesSizes, sampleSize)
				}
			})

			m.truns = append(m.truns, t)
		},
		"tfdt": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			if version == 1 {
				d.FieldU64("start_time")
			} else {
				d.FieldU32("start_time")
			}
		},
		"mvex": decodeBoxes,
		"trex": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("track_id")
			d.FieldU32("default_sample_description_index")
			d.FieldU32("default_sample_duration")
			d.FieldU32("default_sample_size")
			d.FieldU4("reserved0")
			d.FieldU2("is_leading")
			d.FieldU2("sample_depends_on")
			d.FieldU2("sample_is_depended_on")
			d.FieldU2("sample_has_redundancy")
			d.FieldU3("sample_padding_value")
			d.FieldU1("sample_is_non_sync_sample")
			d.FieldU16("sample_degradation_priority")
		},
		"mfra": decodeBoxes,
		"tfra": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("track_id")
			d.FieldU26("reserved")
			lengthSizeOfTrafNum := d.FieldU2("length_size_of_traf_num")
			sampleLengthSizeOfTrunNum := d.FieldU2("sample_length_size_of_trun_num")
			lengthSizeOfSampleNum := d.FieldU2("length_size_of_sample_num")
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						if version == 1 {
							d.FieldU64("time")
							d.FieldU64("moof_offset")
						} else {
							d.FieldU32("time")
							d.FieldU32("moof_offset")
						}
						d.FieldU("traf_number", int(lengthSizeOfTrafNum+1)*8)
						d.FieldU("trun_number", int(sampleLengthSizeOfTrunNum+1)*8)
						d.FieldU("sample_number", int(lengthSizeOfSampleNum+1)*8)
					})
				}
			})
		},
		"mfro": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("mfra_size")
		},
		// TODO: item location
		// HEIC image
		"iloc": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")

			offsetSize := d.FieldU4("offset_size")
			lengthSize := d.FieldU4("length_size")
			baseOffsetSize := d.FieldU4("base_offset_size")
			var indexSize uint64
			switch version {
			case 1, 2:
				indexSize = d.FieldU4("index_size")
			default:
				d.FieldU4("reserved")
			}
			var itemCount uint64
			if version < 2 {
				itemCount = d.FieldU16("item_count")
			} else {
				itemCount = d.FieldU32("item_count")
			}
			d.FieldArray("items", func(d *decode.D) {
				for i := uint64(0); i < itemCount; i++ {
					d.FieldStruct("item", func(d *decode.D) {
						switch version {
						case 0, 1:
							d.FieldU16("id")
						case 2:
							d.FieldU32("id")
						}
						switch version {
						case 1, 2:
							d.FieldU12("reserved")
							d.FieldU4("construction_method")
						}
						d.FieldU16("data_reference_index")
						d.FieldU("base_offset", int(baseOffsetSize)*8)
						extentCount := d.FieldU16("extent_count")
						d.FieldArray("extends", func(d *decode.D) {
							for i := uint64(0); i < extentCount; i++ {
								d.FieldStruct("extent", func(d *decode.D) {
									if (version == 1 || version == 2) && indexSize > 0 {
										d.FieldU("index", int(offsetSize)*8)
									}
									d.FieldU("offset", int(offsetSize)*8)
									d.FieldU("length", int(lengthSize)*8)
								})
							}
						})
					})
				}
			})
		},
		"infe": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU16("id")
			d.FieldU16("protection_index")
			d.FieldUTF8Null("item_name")
			// TODO: really optional? seems so
			if d.NotEnd() {
				d.FieldUTF8Null("content_type")
			}
			if d.NotEnd() {
				d.FieldUTF8Null("content_encoding")
			}
		},
		"iinf": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			if version == 0 {
				_ = d.FieldU16("entry_count")
				decodeBoxes(ctx, d)
			} else {
				d.FieldRawLen("data", d.BitsLeft())
			}
		},
		"iprp": decodeBoxes,
		"ipco": decodeBoxes,
		"ID32": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU1("pad")
			// ISO-639-2/T as 3*5 bit integers - 0x60
			d.FieldStrFn("language", func(d *decode.D) string {
				s := ""
				for i := 0; i < 3; i++ {
					s += fmt.Sprintf("%c", int(d.U5())+0x60)
				}
				return s
			})
			d.FieldFormat("data", id3v2Format, nil)
		},
		"mehd": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			switch version {
			case 0:
				d.FieldU32("fragment_duration")
			case 1:
				d.FieldU64("fragment_duration")
			}
		},
		"pssh": func(_ *decodeContext, d *decode.D) {
			var (
				systemIDCommon    = [16]byte{0x10, 0x77, 0xef, 0xec, 0xc0, 0xb2, 0x4d, 0x02, 0xac, 0xe3, 0x3c, 0x1e, 0x52, 0xe2, 0xfb, 0x4b}
				systemIDWidevine  = [16]byte{0xed, 0xef, 0x8b, 0xa9, 0x79, 0xd6, 0x4a, 0xce, 0xa3, 0xc8, 0x27, 0xdc, 0xd5, 0x1d, 0x21, 0xed}
				systemIDPlayReady = [16]byte{0x9a, 0x04, 0xf0, 0x79, 0x98, 0x40, 0x42, 0x86, 0xab, 0x92, 0xe6, 0x5b, 0xe0, 0x88, 0x5f, 0x95}
				systemIDFairPlay  = [16]byte{0x94, 0xce, 0x86, 0xfb, 0x07, 0xff, 0x4f, 0x43, 0xad, 0xb8, 0x93, 0xd2, 0xfa, 0x96, 0x8c, 0xa2}
			)
			systemIDNames := scalar.BytesToScalar{
				{Bytes: systemIDCommon[:], Scalar: scalar.S{Sym: "common"}},
				{Bytes: systemIDWidevine[:], Scalar: scalar.S{Sym: "widevine"}},
				{Bytes: systemIDPlayReady[:], Scalar: scalar.S{Sym: "playready"}},
				{Bytes: systemIDFairPlay[:], Scalar: scalar.S{Sym: "fairplay"}},
			}

			version := d.FieldU8("version")
			d.FieldU24("flags")
			systemIDBR := d.FieldRawLen("system_id", 16*8, systemIDNames)
			// TODO: make nicer
			systemID := d.ReadAllBits(systemIDBR)
			switch version {
			case 0:
			case 1:
				kidCount := d.FieldU32("kid_count")
				d.FieldArray("kids", func(d *decode.D) {
					for i := uint64(0); i < kidCount; i++ {
						d.FieldRawLen("kid", 16*8)
					}
				})
			}
			dataLen := d.FieldU32("data_size")

			switch {
			case bytes.Equal(systemID, systemIDWidevine[:]):
				d.FieldFormatLen("data", int64(dataLen)*8, protoBufWidevineFormat, nil)
			case bytes.Equal(systemID, systemIDPlayReady[:]):
				d.FieldFormatLen("data", int64(dataLen)*8, psshPlayreadyFormat, nil)
			case systemID == nil:
				fallthrough
			default:
				d.FieldRawLen("data", int64(dataLen)*8)
			}
		},
		"sinf": decodeBoxes,
		"frma": func(ctx *decodeContext, d *decode.D) {
			format := d.FieldUTF8("format", 4)

			// set to original data format
			// TODO: how to handle multiple descriptors? track current?
			if t := ctx.currentTrack(); t != nil && len(t.sampleDescriptions) > 0 {
				t.sampleDescriptions[0].originalFormat = format
			}
		},
		"schm": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8("encryption_type", 4)
			d.FieldU16("encryption_version")
			if d.BitsLeft() > 0 {
				d.FieldUTF8("uri", int(d.BitsLeft())/8)
			}
		},
		"schi": decodeBoxes,
		"btrt": func(_ *decodeContext, d *decode.D) {
			d.FieldU32("decoding_buffer_size")
			d.FieldU32("max_bitrate")
			d.FieldU32("avg_bitrate")
		},
		"pasp": func(_ *decodeContext, d *decode.D) {
			d.FieldU32("h_spacing")
			d.FieldU32("v_spacing")
		},
		"uuid": func(_ *decodeContext, d *decode.D) {
			d.FieldRawLen("uuid", 16*8, scalar.RawUUID, uuidNames)
			d.FieldRawLen("data", d.BitsLeft())
		},
		"keys": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						keySize := d.FieldU32("key_size")
						d.FieldUTF8("key_namespace", 4)
						d.FieldUTF8("key_name", int(keySize)-8)
					})
				}
			})
		},
		"wave": decodeBoxes,
		"saiz": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			flags := d.FieldU24("flags")
			if flags&0b1 != 0 {
				d.FieldU32("aux_info_type")
				d.FieldU32("aux_info_type_parameter")
			}
			defaultSampleInfoSize := d.FieldU8("default_sample_info_size")
			sampleCount := d.FieldU32("sample_count")
			if defaultSampleInfoSize == 0 {
				d.FieldArray("sample_size_info_table", func(d *decode.D) {
					for i := uint64(0); i < sampleCount; i++ {
						d.FieldU8("sample_size")
					}
				})
			}
		},
		"sgpd": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("grouping_type")
			var defaultLength uint64
			if version == 1 {
				defaultLength = d.FieldU32("default_length")
			}
			if version >= 2 {
				d.FieldU32("default_sample_description_index")
			}
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					entryLen := defaultLength
					if version == 1 {
						if defaultLength == 0 {
							entryLen = d.FieldU32("description_length")
						} else if entryLen == 0 {
							d.Fatalf("sgpd groups entry len <= 0 version 1")
						}
					} else if entryLen == 0 {
						d.Fatalf("sgpd groups entry len <= 0")
					}
					d.FieldRawLen("data", int64(entryLen)*8)
				}
			})
		},
		"sbgp": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")

			d.FieldU32("grouping_type")
			if version == 1 {
				d.FieldU32("grouping_type_parameter")
			}
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						d.FieldU32("sample_count")
						d.FieldU32("group_description_index")
					})
				}
			})
		},
		"saio": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			flags := d.FieldU24("flags")

			if flags&0b1 != 0 {
				d.FieldU32("aux_info_type")
				d.FieldU32("aux_info_type_parameter")
			}
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					if version == 0 {
						d.FieldU32("offset")
					} else {
						d.FieldU64("offset")
					}
				}
			})
		},
		"senc": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			flags := d.FieldU24("flags")

			t := ctx.currentTrack()
			if t == nil {
				// need to know iv size
				return
			}
			m := &moof{}
			if t := ctx.currentTrafBox(); t != nil {
				m = t.moof
			}

			s := senc{}
			sampleCount := d.FieldU32("sample_count")
			d.FieldArray("samples", func(d *decode.D) {
				for i := uint64(0); i < sampleCount; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						if t.defaultIVSize != 0 {
							d.FieldRawLen("iv", int64(t.defaultIVSize*8))
						}
						if flags&0b10 != 0 {
							subSampleCount := d.FieldU16("subsample_count")
							d.FieldArray("subsamples", func(d *decode.D) {
								for i := uint64(0); i < subSampleCount; i++ {
									d.FieldStruct("entry", func(d *decode.D) {
										d.FieldU16("bytes_of_clean_data")
										d.FieldU32("bytes_of_encrypted_data")
									})
								}
							})
						}
					})

					// TODO: add iv etc
					s.entries = append(s.entries, struct{}{})
				}
			})
			m.sencs = append(m.sencs, s)
		},
		"tenc": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")

			d.FieldU8("reserved0")
			switch version {
			case 0:
				d.FieldU8("reserved1")
			default:
				d.FieldU4("default_crypto_bytes")
				d.FieldU4("default_skip_bytes")
			}

			defaultIsEncrypted := d.FieldU8("default_is_encrypted")
			defaultIVSize := d.FieldU8("default_iv_size")
			d.FieldRawLen("default_kid", 8*16)

			if defaultIsEncrypted != 0 && defaultIVSize == 0 {
				defaultConstantIVSize := d.FieldU8("default_constant_iv_size")
				d.FieldRawLen("default_constant_iv", int64(defaultConstantIVSize)*8)
			}
			if t := ctx.currentTrack(); t != nil {
				t.defaultIVSize = int(defaultIVSize)
			}
		},
		"covr": decodeBoxes,
		"dec3": func(_ *decodeContext, d *decode.D) {
			d.FieldU13("data_rate")
			d.FieldU3("num_ind_sub")
			d.FieldU2("fscod")
			d.FieldU5("bsid")
			d.FieldU5("bsmod")
			d.FieldU3("acmod")
			d.FieldU1("lfeon")
			d.FieldU3("reserved0")
			numDepSub := d.FieldU4("num_dep_sub")
			if numDepSub > 0 {
				d.FieldU9("chan_loc")
			} else {
				d.FieldU1("reserved1")
			}

			if d.BitsLeft() >= 16 {
				d.FieldU7("reserved2")
				ec3JocFlag := d.FieldBool("ec3_job_flag")
				if ec3JocFlag {
					d.FieldU1("ec3_job_complexity")
				}
			}
		},
		"dac4": func(_ *decodeContext, d *decode.D) {
			d.FieldU3("ac4_dsi_version")
			bitstreamVersion := d.FieldU7("bitstream_version")
			d.FieldU1("fs_index")
			d.FieldU4("frame_rate_index")
			d.FieldU9("n_presentation")

			if bitstreamVersion > 1 {
				hasProgramID := d.FieldBool("has_program_id")
				if hasProgramID {
					d.FieldU16("short_program_id")
					hasUUID := d.FieldBool("has_uuid")
					if hasUUID {
						d.FieldRawLen("uuid", 16*8)
					}
				}
			}

			// if ac4DsiVersion == 1 {
			// 	d.FieldU2("bit_rate_mode")
			// 	d.FieldU32("bit_rate")
			// 	d.FieldU32("bit_rate_precision")
			// }

			// if ac4DsiVersion == 1 {

			// 	d.FieldArray("presentations", func(d *decode.D) {
			// 		for i := uint64(0); i < nPresentation; i++ {
			// 			d.FieldStruct("presentation", func(d *decode.D) {
			// 				d.FieldU8("presentation_version")
			// 				presBytes := d.FieldUFn("pres_bytes", func() (uint64, decode.DisplayFormat, string) {
			// 					n := d.U8()
			// 					if n == 0x0ff {
			// 						n += d.U16()
			// 					}
			// 					return n, decode.NumberDecimal, ""
			// 				})
			// 				d.FieldRawLen("data", int64(presBytes)*8)
			// 			})
			// 		}
			// 	})
			// }

			if d.BitsLeft() > 0 {
				d.FieldRawLen("data", d.BitsLeft())
			}
		},
		"tapt": decodeBoxes,
		"clef": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldFP32("width")
			d.FieldFP32("height")
		},
		"prof": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldFP32("width")
			d.FieldFP32("height")
		},
		"enof": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldFP32("width")
			d.FieldFP32("height")
		},
		"clap": func(_ *decodeContext, d *decode.D) {
			d.FieldU32("aperture_width_n")
			d.FieldU32("aperture_width_d")
			d.FieldU32("aperture_height_n")
			d.FieldU32("aperture_height_d")
			d.FieldU32("horiz_off_n")
			d.FieldU32("horiz_off_d")
			d.FieldU32("vert_off_n")
			d.FieldU32("vert_off_d")
		},
		"smhd": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldFP16("balance")
			d.FieldU16("reserved")
		},
		"colr": func(_ *decodeContext, d *decode.D) {
			parameterType := d.FieldUTF8("parameter_type", 4)

			switch parameterType {
			case "nclx", "nclc":
				d.FieldU16("primaries_index", format.ISO_23091_2_ColourPrimariesMap)
				d.FieldU16("transfer_function_index", format.ISO_23091_2_TransferCharacteristicMap)
				d.FieldU16("matrix_index", format.ISO_23091_2_MatrixCoefficients)
				switch parameterType {
				case "nclx":
					d.FieldU8("color_range")
				}
			case "prof":
				d.FieldFormat("profile", iccProfileFormat, nil)
			default:
				d.FieldRawLen("data", d.BitsLeft())
			}
		},
		"ispe": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("image_width")
			d.FieldU32("image_height")
		},
		"ipma": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			flags := d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			d.FieldArray("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStruct("entry", func(d *decode.D) {
						if version < 1 {
							d.FieldU16("item_id")
						} else {
							d.FieldU32("item_id")
						}
						associationCount := d.FieldU8("association_count")
						d.FieldArray("associations", func(d *decode.D) {
							for j := uint64(0); j < associationCount; j++ {
								d.FieldStruct("association", func(d *decode.D) {
									d.FieldBool("essential")
									if flags&0b1 != 0 {
										d.FieldU15("property_index")
									} else {
										d.FieldU7("item_id")
									}
								})
							}
						})
					})
				}
			})
		},
		"pitm": func(_ *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			if version < 1 {
				d.FieldU16("item_id")
			} else {
				d.FieldU32("item_id")
			}
		},
		"iref": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			decodeBoxesWithParentData(ctx, d, &irefBox{version: int(version)})
		},
		"dimg": irefEntryDecode,
		"thmb": irefEntryDecode,
		"cdsc": irefEntryDecode,
		"irot": func(_ *decodeContext, d *decode.D) {
			d.FieldU8("rotation", scalar.UToSymU{
				0: 0,
				1: 90,
				2: 180,
				3: 270,
			})
		},
	}
}
