package mp4

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/text/encoding/charmap"
)

// TODO: keep track of list of sampleSize/entries instead and change sample read code
const maxSampleEntryCount = 10_000_000

const (
	boxSizeRestOfFile   = 0
	boxSizeUse64bitSize = 1
)

var boxSizeNames = scalar.UintMapDescription{
	boxSizeRestOfFile:   "Rest of file",
	boxSizeUse64bitSize: "Use 64 bit size",
}

var mediaTimeNames = scalar.SintMapDescription{
	-1: "empty",
}

var subTypeNames = scalar.StrMapDescription{
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
	"url":  "URL",
	"vide": "Video Track",
}

var dataFormatNames = scalar.StrMapDescription{
	// additional codecs
	"apch": "Apple ProRes 422 High Quality",
	"apcn": "Apple ProRes 422 Standard Definition",
	"apcs": "Apple ProRes 422 LT",
	"apco": "Apple ProRes 422 Proxy",
	"ap4h": "Apple ProRes 4444",
	"jpeg": "JPEG Image",

	// codecs from https://mp4ra.org/
	"3gvo": "3GPP Video Orientation",
	"a3d1": "Multiview Video Coding",
	"a3d2": "Multiview Video Coding",
	"a3d3": "Multiview Video Coding",
	"a3d4": "Multiview Video Coding",
	"a3ds": "Auro-Cx 3D audio",
	"ac-3": "AC-3 audio",
	"ac-4": "AC-4 audio",
	"alac": "Apple lossless audio codec",
	"alaw": "a-Law",
	"av01": "AV1 video",
	"avc1": "Advanced Video Coding",
	"avc2": "Advanced Video Coding",
	"avc3": "Advanced Video Coding",
	"avc4": "Advanced Video Coding",
	"avcp": "Advanced Video Coding Parameters",
	"dra1": "DRA Audio",
	"drac": "Dirac Video Coder",
	"dts-": "Dependent base layer for DTS layered audio",
	"dts+": "Enhancement layer for DTS layered audio",
	"dtsc": "DTS Coherent Acoustics audio",
	"dtse": "DTS Express low bit rate audio, also known as DTS LBR",
	"dtsh": "DTS-HD High Resolution Audio",
	"dtsl": "DTS-HD Master Audio",
	"dtsx": "DTS:X",
	"dvav": "AVC-based “Dolby Vision”",
	"dvhe": "HEVC-based “Dolby Vision”",
	"ec-3": "Enhanced AC-3 audio",
	"enca": "Encrypted/Protected audio",
	"encf": "Encrypted/Protected font",
	"encm": "Encrypted/Protected metadata stream",
	"encs": "Encrypted Systems stream",
	"enct": "Encrypted Text",
	"encv": "Encrypted/protected video",
	"fdp":  "File delivery hints",
	"fLaC": "Fres Lossless Audio Codec",
	"g719": "ITU-T Recommendation G.719 (2008)",
	"g726": "ITU-T Recommendation G.726 (1990)",
	"hev1": "High Efficiency Video Coding",
	"hvc1": "High Efficiency Video Coding",
	"hvt1": "High Efficiency Video Coding",
	"ixse": "DVB Track Level Index Track",
	"lhe1": "Layered High Efficiency Video Coding",
	"lht1": "Layered High Efficiency Video Coding",
	"lhv1": "Layered High Efficiency Video Coding",
	"m2ts": "MPEG-2 transport stream for DMB",
	"m4ae": "MPEG-4 Audio Enhancement MP4v1/2",
	"mett": "Text timed metadata that is not XML",
	"metx": "XML timed metadata",
	"mha1": "MPEG-H Audio (single stream, uncapsulated)",
	"mha2": "MPEG-H Audio (multi-stream, unencapsulated)",
	"mhm1": "MPEG-H Audio (single stream, MHAS encapsulated)",
	"mhm2": "MPEG-H Audio (multi-stream, MHAS encapsulated)",
	"mjp2": "Motion JPEG 2000",
	"mlix": "DVB Movie level index track",
	"mlpa": "MLP Audio",
	"mp4a": "MPEG-4 Audio",
	"mp4s": "MPEG-4 Systems",
	"mp4v": "MPEG-4 Visual",
	"mvc1": "Multiview coding",
	"mvc2": "Multiview coding",
	"mvc3": "Multiview coding",
	"mvc4": "Multiview coding",
	"mvd1": "Multiview coding",
	"mvd2": "Multiview coding",
	"mvd3": "Multiview coding",
	"mvd4": "Multiview coding",
	"oksd": "OMA Keys",
	"Opus": "Opus audio coding",
	"pm2t": "Protected MPEG-2 Transport",
	"prtp": "Protected RTP Reception",
	"raw":  "Uncompressed audio",
	"resv": "Restricted Video",
	"rm2t": "MPEG-2 Transport Reception",
	"rrtp": "RTP reception",
	"rsrp": "SRTP Reception",
	"rtmd": "Real Time Metadata Sample Entry(XAVC Format)",
	"rtp":  "RTP Hints",
	"s263": "ITU H.263 video (3GPP format)",
	"samr": "Narrowband AMR voice",
	"sawb": "Wideband AMR voice",
	"sawp": "Extended AMR-WB (AMR-WB+)",
	"sevc": "EVRC Voice",
	"sm2t": "MPEG-2 Transport Server",
	"sqcp": "13K Voice",
	"srtp": "SRTP Hints",
	"ssmv": "SMV Voice",
	"STGS": "Subtitle Sample Entry (HMMP)",
	"stpp": "Subtitles (Timed Text)",
	"svc1": "Scalable Video Coding",
	"svc2": "Scalable Video Coding",
	"svcM": "SVC Metadata",
	"tc64": "64 bit timecode samples",
	"tmcd": "32 bit timecode samples",
	"twos": "Uncompressed 16-bit audio",
	"tx3g": "Timed Text stream",
	"ulaw": "Samples have been compressed using uLaw 2:1.",
	"unid": "Dynamic Range Control (DRC) data",
	"urim": "Binary timed metadata identified by URI",
	"vc-1": "SMPTE VC-1",
	"vp08": "VP8 video",
	"vp09": "VP9 video",
	"wvtt": "WebVTT",
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

var uuidNames = scalar.RawBytesMap{
	{Bytes: uuidIsmlManifestBytes[:], Scalar: scalar.BitBuf{Sym: "isml_manifest"}},
	{Bytes: uuidXmpBytes[:], Scalar: scalar.BitBuf{Sym: "xmp"}},
	{Bytes: uuidSphericalBytes[:], Scalar: scalar.BitBuf{Sym: "spherical"}},
	{Bytes: uuidPspUsmtBytes[:], Scalar: scalar.BitBuf{Sym: "psp_usmt"}},
	{Bytes: uuidTfxdBytes[:], Scalar: scalar.BitBuf{Sym: "tfxd"}},
	{Bytes: uuidTfrfBytes[:], Scalar: scalar.BitBuf{Sym: "tfrf"}},
	{Bytes: uuidProfBytes[:], Scalar: scalar.BitBuf{Sym: "prof"}},
	{Bytes: uuidIpodBytes[:], Scalar: scalar.BitBuf{Sym: "ipod"}},
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
var quicktimeEpochDate = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)

var uintActualQuicktimeEpochDescription = scalar.UintActualDateDescription(quicktimeEpochDate, time.Second, time.RFC3339)

func decodeMvhdFieldMatrix(d *decode.D, name string) {
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

// ISO 14496-12 8.40.2.3 Sample dependency box semantics
func decodeSampleFlags(d *decode.D) {
	d.FieldU4("reserved0")
	d.FieldU2("is_leading")
	d.FieldU2("sample_depends_on", scalar.UintMap{
		0: scalar.Uint{Sym: "unknown"},
		1: scalar.Uint{Sym: "other", Description: "Not I-picture"},
		2: scalar.Uint{Sym: "none", Description: "Is I-picture"},
	})
	d.FieldU2("sample_is_depended_on", scalar.UintMap{
		0: scalar.Uint{Sym: "unknown"},
		1: scalar.Uint{Sym: "other", Description: "Not disposable"},
		2: scalar.Uint{Sym: "none", Description: "Is disposable"},
	})
	d.FieldU2("sample_has_redundancy", scalar.UintMap{
		0: scalar.Uint{Sym: "unknown"},
		2: scalar.Uint{Sym: "none", Description: "No redundant coding"},
	})
	d.FieldU3("sample_padding_value")
	d.FieldU1("sample_is_non_sync_sample")
	d.FieldU16("sample_degradation_priority")
}

func decodeBoxWithParentData(ctx *decodeContext, d *decode.D, parentData any, extraTypeMappers ...scalar.StrMapper) {
	var dataSize uint64
	typeMappers := []scalar.StrMapper{scalar.ActualTrimSpace, boxDescriptions}
	if len(extraTypeMappers) > 0 {
		typeMappers = append(typeMappers, extraTypeMappers...)
	}

	boxSize := d.FieldU32("size", boxSizeNames)
	typ := d.FieldStr("type", 4, charmap.ISO8859_1, typeMappers...)

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

	if parentData != nil {
		ctx.path[len(ctx.path)-1].data = parentData
	}
	ctx.path = append(ctx.path, pathEntry{typ: typ, data: parentData})
	d.FramedFn(int64(dataSize*8), func(d *decode.D) {
		decodeBox(ctx, d, typ)
	})
	ctx.path = ctx.path[0 : len(ctx.path)-1]
}

func decodeBoxes(ctx *decodeContext, d *decode.D, extraTypeMappers ...scalar.StrMapper) {
	decodeBoxesWithParentData(ctx, d, nil, extraTypeMappers...)
}

func decodeBoxesWithParentData(ctx *decodeContext, d *decode.D, parentData any, extraTypeMappers ...scalar.StrMapper) {
	d.FieldStructArrayLoop("boxes", "box",
		func() bool { return d.BitsLeft() >= 8*8 },
		func(d *decode.D) {
			decodeBoxWithParentData(ctx, d, parentData, extraTypeMappers...)
		})

	if d.BitsLeft() > 0 {
		// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
		if d.BitsLeft() >= 32 && d.PeekUintBits(32) == 0 {
			d.FieldU32("zero_terminator")
		}
		if d.BitsLeft() > 0 {
			d.FieldRawLen("padding", d.BitsLeft())
		}
	}
}

type rootBox struct {
	ftypMajorBrand string
}

type irefBox struct {
	version int
}

type trakBox struct {
	track *track
}

type moofBox struct {
	offset int64
}

type trafBox struct {
	track          *track
	baseDataOffset int64
	moof           *moof
}

type metaBox struct {
	subType string
	keys    *keysBox
}

type keysBoxKey struct {
	namespace string
	name      string
}

type keysBox struct {
	keys []keysBoxKey
}

func decodeBoxIrefEntry(irefBox *irefBox, d *decode.D) {
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

func decodeBoxFtyp(ctx *decodeContext, d *decode.D) {
	root := ctx.rootBox()

	brand := d.FieldUTF8("major_brand", 4, scalar.ActualTrimSpace)
	root.ftypMajorBrand = brand

	d.FieldU32("minor_version", scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		switch brand {
		case "qt":
			// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap1/qtff1.html#//apple_ref/doc/uid/TP40000939-CH203-BBCGDDDF
			// "For QuickTime movie files, this takes the form of four binary-coded decimal values, indicating the century,
			//  year, and month of the QuickTime File Format Specification, followed by a binary coded decimal zero. For example,
			//  for the June 2004 minor version, this field is set to the BCD values 20 04 06 00."
			s.Description = fmt.Sprintf("%.4d.%.2d", (s.Actual>>24)&0xff_ff, (s.Actual>>8)&0xff)
		}
		return s, nil
	}))
	numBrands := d.BitsLeft() / 8 / 4
	var i int64
	d.FieldArrayLoop("brands", func() bool { return i < numBrands }, func(d *decode.D) {
		d.FieldUTF8("brand", 4, brandDescriptions, scalar.ActualTrimSpace)
		i++
	})
}

func decodeBox(ctx *decodeContext, d *decode.D, typ string) {
	switch typ {
	case "ftyp":
		decodeBoxFtyp(ctx, d)
	case "styp":
		decodeBoxFtyp(ctx, d)
	case "mvhd":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		switch version {
		case 0:
			d.FieldU32("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("modification_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("time_scale")
			d.FieldU32("duration")
		case 1:
			d.FieldU64("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU64("modification_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("time_scale")
			d.FieldU64("duration")
		default:
			return
		}
		d.FieldFP32("preferred_rate")
		d.FieldFP16("preferred_volume")
		d.FieldUTF8("reserved", 10)
		decodeMvhdFieldMatrix(d, "matrix_structure")
		d.FieldU32("preview_time")
		d.FieldU32("preview_duration")
		d.FieldU32("poster_time")
		d.FieldU32("selection_time")
		d.FieldU32("selection_duration")
		d.FieldU32("current_time")
		d.FieldU32("next_track_id")
	case "trak":
		t := &track{}
		ctx.tracks = append(ctx.tracks, t)
		decodeBoxesWithParentData(ctx, d, &trakBox{
			track: t,
		})
	case "edts":
		decodeBoxes(ctx, d)
	case "elst":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		var i uint64
		d.FieldStructArrayLoop("entries", "entry", func() bool { return i < entryCount }, func(d *decode.D) {
			switch version {
			case 0:
				d.FieldS32("segment_duration")
				d.FieldS32("media_time", mediaTimeNames)
			case 1:
				d.FieldS64("segment_duration")
				d.FieldS64("media_time", mediaTimeNames)
			default:
				return
			}
			d.FieldFP32("media_rate")
			i++
		})
	case "tref":
		decodeBoxes(ctx, d)
	case "tkhd":
		var trackID int
		version := d.FieldU8("version")
		d.FieldStruct("flags", func(d *decode.D) {
			d.FieldU20("unused")
			d.FieldBool("size_is_aspect_ratio")
			d.FieldBool("in_preview")
			d.FieldBool("in_movie")
			d.FieldBool("enabled")
		})
		switch version {
		case 0:
			d.FieldU32("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("modification_time", uintActualQuicktimeEpochDescription)
			trackID = int(d.FieldU32("track_id"))
			d.FieldU32("reserved1")
			d.FieldU32("duration")
		case 1:
			d.FieldU64("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU64("modification_time", uintActualQuicktimeEpochDescription)
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
		decodeMvhdFieldMatrix(d, "matrix_structure")
		d.FieldFP32("track_width")
		d.FieldFP32("track_height")

		if t := ctx.currentTrakBox(); t != nil {
			t.track.id = trackID
		}
	case "mdia":
		decodeBoxes(ctx, d)
	case "mdhd":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		// TODO: timestamps
		switch version {
		case 0:
			d.FieldU32("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("modification_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("time_scale")
			d.FieldU32("duration")
		case 1:
			d.FieldU64("creation_time", uintActualQuicktimeEpochDescription)
			d.FieldU64("modification_time", uintActualQuicktimeEpochDescription)
			d.FieldU32("time_scale")
			d.FieldU64("duration")
		default:
			return
		}
		d.FieldStrFn("language", decodeLang)
		d.FieldU16("quality")
	case "vmhd":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU16("graphicsmode")
		d.FieldArray("opcolor", func(d *decode.D) {
			// TODO: is FP16?
			d.FieldU16("value")
			d.FieldU16("value")
			d.FieldU16("value")
		})
	case "hdlr":
		majorBrand := ctx.rootBox().ftypMajorBrand

		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldUTF8NullFixedLen("component_type", 4)
		subType := d.FieldUTF8("component_subtype", 4, scalar.ActualTrimSpace, subTypeNames)
		d.FieldUTF8NullFixedLen("component_manufacturer", 4)
		d.FieldU32("component_flags")
		d.FieldU32("component_flags_mask")

		switch majorBrand {
		case "qt":
			// qt brand seems to use length prefixed strings
			// From QuickTime File Format specification:
			// > A (counted) string that specifies the name of the component—that is, the media handler used
			// > when this media was created. This field may contain a zero-length (empty) string.
			if d.BitsLeft() > 0 {
				d.FieldUTF8ShortStringFixedLen("component_name", int(d.BitsLeft()/8))
			}
		default:
			d.FieldUTF8NullFixedLen("component_name", int(d.BitsLeft()/8))
		}

		if t := ctx.currentTrack(); t != nil {
			t.seenHdlr = true
			// component_type seems to be all zero sometimes so can't look for "mhlr"
			switch subType {
			case "vide", "soun":
				t.subType = subType
			}
		} else if m := ctx.currentMetaBox(); m != nil {
			m.subType = subType
		}
	case "minf":
		decodeBoxes(ctx, d)
	case "dinf":
		decodeBoxes(ctx, d)
	case "dref":
		d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		var i uint64
		var drefURL string
		d.FieldStructArrayLoop("boxes", "box", func() bool { return i < entryCount }, func(d *decode.D) {
			size := d.FieldU32("size")
			typ := d.FieldUTF8("type", 4, scalar.ActualTrimSpace)
			d.FieldU8("version")
			d.FieldU24("flags")
			dataSize := size - 12
			switch typ {
			case "url":
				drefURL = d.FieldUTF8("data", int(dataSize))
			default:
				d.FieldRawLen("data", int64(dataSize*8))
			}
			i++
		})

		if t := ctx.currentTrack(); t != nil {
			t.dref = true
			t.drefURL = drefURL
		}

	case "stbl":
		decodeBoxes(ctx, d)
	case "stsd":
		d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		var i uint64
		// note called "boxes" here instead of "sample_descriptions" and data format is named "type".
		// this is to make it easier to threat them as normal boxes
		d.FieldArrayLoop("boxes", func() bool { return i < entryCount }, func(d *decode.D) {
			d.FieldStruct("box", func(d *decode.D) {
				size := d.FieldU32("size")
				dataFormat := d.FieldUTF8("type", 4, dataFormatNames, scalar.ActualTrimSpace)
				subType := ""
				if t := ctx.currentTrack(); t != nil {
					t.sampleDescriptions = append(t.sampleDescriptions, sampleDescription{
						dataFormat: dataFormat,
					})

					if t.seenHdlr {
						subType = t.subType
					} else {
						// TODO: seems to be ffmpeg mov.c, where is this documented in specs?
						// no hdlr box found, guess using dataFormat
						// ex PNG samples but there is no hdlr box saying it's video, but the esds says MPEGObjectTypePNG
						switch dataFormat {
						case "mp4v":
							subType = "vide"
						case "mp4a":
							subType = "soun"
						}
					}
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
	case "avcC":
		_, v := d.FieldFormat("descriptor", &avcDCRGroup, nil)
		avcDcrOut, ok := v.(format.AVC_DCR_Out)
		if !ok {
			panic(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
		}
		if t := ctx.currentTrack(); t != nil {
			t.formatInArg = format.AVC_AU_In(avcDcrOut)
		}
	case "hvcC":
		_, v := d.FieldFormat("descriptor", &hevcCDCRGroup, nil)
		hevcDcrOut, ok := v.(format.HEVC_DCR_Out)
		if !ok {
			panic(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
		}
		if t := ctx.currentTrack(); t != nil {
			t.formatInArg = format.HEVC_AU_In(hevcDcrOut)
		}
	case "dfLa":
		d.FieldU8("version")
		d.FieldU24("flags")
		_, v := d.FieldFormat("descriptor", &flacMetadatablocksGroup, nil)
		flacMetadatablockOut, ok := v.(format.FLAC_Metadatablocks_Out)
		if !ok {
			panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
		}
		if flacMetadatablockOut.HasStreamInfo {
			if t := ctx.currentTrack(); t != nil {
				t.formatInArg = format.FLAC_Frame_In{BitsPerSample: int(flacMetadatablockOut.StreamInfo.BitsPerSample)}
			}
		}
	case "dOps":
		d.FieldFormat("descriptor", &opusPacketFrameGroup, nil)
	case "av1C":
		d.FieldFormat("descriptor", &av1CCRGroup, nil)
	case "vpcC":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldFormat("descriptor", &vpxCCRGroup, nil)
	case "iods":
		d.FieldU32("version")
		d.FieldFormat("descriptor", &mpegESGroup, nil)
	case "esds":
		d.FieldU32("version")
		_, v := d.FieldFormat("descriptor", &mpegESGroup, nil)
		mpegEsOut, ok := v.(format.MPEG_ES_Out)
		if !ok {
			panic(fmt.Sprintf("expected mpegEsOut got %#+v", v))
		}

		if t := ctx.currentTrack(); t != nil && len(mpegEsOut.DecoderConfigs) > 0 {
			dc := mpegEsOut.DecoderConfigs[0]
			t.objectType = dc.ObjectType
			t.formatInArg = format.AAC_Frame_In{ObjectType: dc.ASCObjectType}
		}
	case "stts":
		d.FieldU8("version")
		d.FieldU24("flags")
		numEntries := d.FieldU32("entry_count")
		var i uint64
		d.FieldStructArrayLoop("entries", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
			d.FieldU32("count")
			d.FieldU32("delta")
			i++
		})
	case "stsc":
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
	case "stsz":
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
	case "stz2":
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
	case "stco":
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
	case "stss":
		d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		d.FieldArray("entries", func(d *decode.D) {
			for i := uint64(0); i < entryCount; i++ {
				d.FieldU32("sample_number")
			}
		})
	case "sdtp":
		d.FieldU8("version")
		d.FieldU24("flags")
		// TODO: should be count from stsz
		// TODO: can we know count here or do we need to defer decoding somehow?
		d.FieldArray("entries", func(d *decode.D) {
			for d.NotEnd() {
				d.FieldStruct("entry", func(d *decode.D) {
					d.FieldU2("reserved")
					values := scalar.UintMapSymStr{
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
	case "ctts":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		var i uint64
		d.FieldStructArrayLoop("entries", "entry", func() bool { return i < entryCount }, func(d *decode.D) {
			d.FieldS32("sample_count")
			// ISO/IEC14496-12 says version 0 is unsigned and version 1 is signed
			// in preactice it seems muxers write it as signed for both version
			switch version {
			case 0, 1:
				d.FieldS32("sample_offset")
			}
			i++
		})
		// TODO: refactor: merge with stco?
	case "co64":
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
	case "sidx":
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
	case "udta":
		decodeBoxes(ctx, d)
	case "meta":
		// TODO: meta box sometimes has a 4 byte unknown field? (flag/version?)
		maybeFlags := d.PeekUintBits(32)
		if maybeFlags == 0 {
			// TODO: rename?
			d.FieldU32("maybe_flags")
		}
		decodeBoxesWithParentData(ctx, d, &metaBox{})
	case "ilst":
		if mb := ctx.currentMetaBox(); mb != nil && mb.keys != nil && len(mb.keys.keys) > 0 {
			// meta box had a keys box
			var b [4]byte
			typeSymMapper := scalar.StrMapSymStr{}
			for k, v := range mb.keys.keys {
				// type will be a uint32 be integer
				// +1 as they seem to be counted from 1
				binary.BigEndian.PutUint32(b[:], uint32(k+1))
				typeSymMapper[string(b[:])] = v.namespace + "." + v.name
			}

			decodeBoxes(ctx, d, typeSymMapper)
		} else {
			decodeBoxes(ctx, d)
		}
	case "data":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU32("reserved")
		if ctx.isParent("covr") {
			d.FieldFormatOrRawLen("data", d.BitsLeft(), &imageGroup, nil)
		} else {
			d.FieldUTF8("data", int(d.BitsLeft()/8))
		}
	case "moov":
		decodeBoxes(ctx, d)
	case "moof":
		offset := (d.Pos() / 8) - 8
		decodeBoxesWithParentData(ctx, d, &moofBox{offset: offset})
	case "traf": // Track Fragment
		t := &track{fragment: true}
		ctx.tracks = append(ctx.tracks, t)
		decodeBoxesWithParentData(ctx, d, &trafBox{
			track: t,
		})
	case "mfhd": // Movie Fragment Header
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU32("sequence_number")
	case "tfhd": // Track Fragment Header
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
			d.FieldStruct("default_sample_flags", decodeSampleFlags)
		}

		if t := ctx.currentTrafBox(); t != nil {
			t.track.id = trackID
			t.moof = m
			t.baseDataOffset = baseDataOffset
		}
		if t := ctx.currentTrack(); t != nil {
			t.moofs = append(t.moofs, m)
		}
	case "trun": // Track Fragment Run
		m := &moof{}
		if t := ctx.currentTrafBox(); t != nil {
			m = t.moof
		}

		version := d.FieldU8("version")
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
			d.FieldStruct("first_sample_flags", decodeSampleFlags)
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
						d.FieldStruct("sample_flags", decodeSampleFlags)
					}
					if sampleCompositionTimeOffsetsPresent {
						if version == 0 {
							d.FieldU32("sample_composition_time_offset")
						} else {
							d.FieldS32("sample_composition_time_offset")
						}
					}
				})

				t.samplesSizes = append(t.samplesSizes, sampleSize)
			}
		})

		m.truns = append(m.truns, t)
	case "tfdt":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		if version == 1 {
			d.FieldU64("start_time")
		} else {
			d.FieldU32("start_time")
		}
	case "mvex":
		decodeBoxes(ctx, d)
	case "trex":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU32("track_id")
		d.FieldU32("default_sample_description_index")
		d.FieldU32("default_sample_duration")
		d.FieldU32("default_sample_size")
		d.FieldStruct("default_sample_flags", decodeSampleFlags)
	case "mfra":
		decodeBoxes(ctx, d)
	case "tfra":
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
	case "mfro":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU32("mfra_size")

	case "iloc": // HEIC image
		// TODO: item location
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
	case "infe":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		if version == 0 || version == 1 {
			d.FieldU16("item_id")
			d.FieldU16("item_protection_index")
			d.FieldUTF8Null("item_name")
			d.FieldUTF8Null("content_type")
			if !d.End() {
				d.FieldUTF8Null("content_encoding")
			}
		}
		if version == 1 {
			if !d.End() {
				d.FieldU32("extension_type")
			}
			if !d.End() {
				d.FieldU32("extension_type_extra")
			}
		}
		if version >= 2 {
			switch version {
			case 2:
				d.FieldU16("item_id")
			case 3:
				d.FieldU32("item_id")
			}
			d.FieldU16("item_protection_index")
			itemType := d.FieldUTF8("item_type", 4)
			d.FieldUTF8Null("item_name")
			switch itemType {
			case "mime":
				d.FieldUTF8Null("content_type")
				if !d.End() {
					d.FieldUTF8Null("content_encoding")
				}
			case "uri":
				d.FieldUTF8Null("item_uri_type")
			}
		}
	case "iinf":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		if version == 0 {
			d.FieldU16("entry_count")
		} else {
			d.FieldU32("entry_count")
		}
		decodeBoxes(ctx, d)
	case "iprp":
		decodeBoxes(ctx, d)
	case "ipco":
		decodeBoxes(ctx, d)
	case "ID32":
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
		d.FieldFormat("data", &id3v2Group, nil)
	case "mehd":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		switch version {
		case 0:
			d.FieldU32("fragment_duration")
		case 1:
			d.FieldU64("fragment_duration")
		}
	case "pssh":
		var (
			systemIDCommon    = [16]byte{0x10, 0x77, 0xef, 0xec, 0xc0, 0xb2, 0x4d, 0x02, 0xac, 0xe3, 0x3c, 0x1e, 0x52, 0xe2, 0xfb, 0x4b}
			systemIDWidevine  = [16]byte{0xed, 0xef, 0x8b, 0xa9, 0x79, 0xd6, 0x4a, 0xce, 0xa3, 0xc8, 0x27, 0xdc, 0xd5, 0x1d, 0x21, 0xed}
			systemIDPlayReady = [16]byte{0x9a, 0x04, 0xf0, 0x79, 0x98, 0x40, 0x42, 0x86, 0xab, 0x92, 0xe6, 0x5b, 0xe0, 0x88, 0x5f, 0x95}
			systemIDFairPlay  = [16]byte{0x94, 0xce, 0x86, 0xfb, 0x07, 0xff, 0x4f, 0x43, 0xad, 0xb8, 0x93, 0xd2, 0xfa, 0x96, 0x8c, 0xa2}
		)
		systemIDNames := scalar.RawBytesMap{
			{Bytes: systemIDCommon[:], Scalar: scalar.BitBuf{Sym: "common"}},
			{Bytes: systemIDWidevine[:], Scalar: scalar.BitBuf{Sym: "widevine"}},
			{Bytes: systemIDPlayReady[:], Scalar: scalar.BitBuf{Sym: "playready"}},
			{Bytes: systemIDFairPlay[:], Scalar: scalar.BitBuf{Sym: "fairplay"}},
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
			d.FieldFormatLen("data", int64(dataLen)*8, &protoBufWidevineGroup, nil)
		case bytes.Equal(systemID, systemIDPlayReady[:]):
			d.FieldFormatLen("data", int64(dataLen)*8, &psshPlayreadyGroup, nil)
		case systemID == nil:
			fallthrough
		default:
			d.FieldRawLen("data", int64(dataLen)*8)
		}
	case "sinf":
		decodeBoxes(ctx, d)
	case "frma":
		format := d.FieldUTF8("format", 4)

		// set to original data format
		// TODO: how to handle multiple descriptors? track current?
		if t := ctx.currentTrack(); t != nil && len(t.sampleDescriptions) > 0 {
			t.sampleDescriptions[0].originalFormat = format
		}
	case "schm":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldUTF8("encryption_type", 4)
		d.FieldU16("encryption_version")
		if d.BitsLeft() > 0 {
			d.FieldUTF8("uri", int(d.BitsLeft())/8)
		}
	case "schi":
		decodeBoxes(ctx, d)
	case "btrt":
		d.FieldU32("decoding_buffer_size")
		d.FieldU32("max_bitrate")
		d.FieldU32("avg_bitrate")
	case "pasp":
		d.FieldU32("h_spacing")
		d.FieldU32("v_spacing")
	case "uuid":
		d.FieldRawLen("uuid", 16*8, scalar.RawUUID, uuidNames)
		d.FieldRawLen("data", d.BitsLeft())
	case "keys":
		mb := ctx.currentMetaBox()
		var kb *keysBox
		if mb != nil {
			kb = &keysBox{}
			mb.keys = kb
		}

		d.FieldU8("version")
		d.FieldU24("flags")
		entryCount := d.FieldU32("entry_count")
		d.FieldArray("entries", func(d *decode.D) {
			for i := uint64(0); i < entryCount; i++ {
				d.FieldStruct("entry", func(d *decode.D) {
					keySize := d.FieldU32("key_size")
					namespace := d.FieldUTF8("key_namespace", 4)
					name := d.FieldUTF8("key_name", int(keySize)-8)
					if kb != nil {
						kb.keys = append(kb.keys, keysBoxKey{
							namespace: namespace,
							name:      name,
						})
					}
				})
			}
		})
	case "wave":
		decodeBoxes(ctx, d)
	case "saiz":
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
	case "sgpd":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		var groupingType = d.FieldUTF8("grouping_type", 4)
		var defaultLength uint64
		if version == 1 {
			defaultLength = d.FieldU32("default_length")
		} else if version >= 2 {
			d.FieldU32("default_sample_description_index")
		}
		entryCount := d.FieldU32("entry_count")

		d.FieldStructNArray("entries", "entry", int64(entryCount), func(d *decode.D) {
			entryLen := defaultLength
			if version == 1 {
				if defaultLength == 0 {
					entryLen = d.FieldU32("description_length")
				} else if entryLen == 0 {
					d.Fatalf("sgpd groups entry len == 0")
				}
			} else if entryLen == 0 {
				// TODO: this is likely a mistake: here version is != 1, so defaultLength is default (i.e. 0),
				// TODO: so entryLen is also 0. So for version != 1 this fatal should always throw
				d.Fatalf("sgpd groups entry len == 0")
			}
			// CENC Sample Encryption Info Entries
			switch groupingType {
			case "seig":
				d.FieldU8("reserved")
				d.FieldU4("crypto_bytes")
				d.FieldU4("skip_bytes")
				isEncrypted := d.FieldU8("is_encrypted")
				perSampleIVSize := d.FieldU8("per_sample_iv_size")
				d.FieldRawLen("kid", 8*16)
				if isEncrypted != 0 {
					// If perSampleIVSize > 0 then
					if perSampleIVSize == 0 {
						// This means whole fragment is encrypted with a constant IV
						iVSize := d.FieldU8("constant_iv_size")
						d.FieldRawLen("constant_iv", 8*int64(iVSize))
					}
				}
			case "roll":
				d.FieldU16("roll_distance")
			default:
				d.FieldRawLen("data", int64(entryLen)*8)
			}
		})
	case "sbgp":
		version := d.FieldU8("version")
		d.FieldU24("flags")

		d.FieldUTF8("grouping_type", 4)
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
	case "saio":
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
	case "senc":
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
	case "tenc":
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
	case "covr":
		decodeBoxes(ctx, d)
	case "dec3":
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
	case "dac4":
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
		// 				presBytes := d.FieldUintFn("pres_bytes", func() (uint64, decode.DisplayFormat, string) {
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
	case "tapt":
		decodeBoxes(ctx, d)
	case "clef":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldFP32("width")
		d.FieldFP32("height")
	case "prof":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldFP32("width")
		d.FieldFP32("height")
	case "enof":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldFP32("width")
		d.FieldFP32("height")
	case "clap":
		d.FieldU32("aperture_width_n")
		d.FieldU32("aperture_width_d")
		d.FieldU32("aperture_height_n")
		d.FieldU32("aperture_height_d")
		d.FieldU32("horiz_off_n")
		d.FieldU32("horiz_off_d")
		d.FieldU32("vert_off_n")
		d.FieldU32("vert_off_d")
	case "smhd":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldFP16("balance")
		d.FieldU16("reserved")
	case "colr":
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
			d.FieldFormat("profile", &iccProfileGroup, nil)
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	case "ispe":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldU32("image_width")
		d.FieldU32("image_height")
	case "ipma":
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
	case "pitm":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		if version == 0 {
			d.FieldU16("item_id")
		} else {
			d.FieldU32("item_id")
		}
	case "iref":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		decodeBoxesWithParentData(ctx, d, &irefBox{version: int(version)})
	case "dimg":
		if irefBox, ok := ctx.parent().data.(*irefBox); ok {
			decodeBoxIrefEntry(irefBox, d)
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	case "thmb":
		if irefBox, ok := ctx.parent().data.(*irefBox); ok {
			decodeBoxIrefEntry(irefBox, d)
		} else {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8("format", 4)
			d.FieldFormatOrRawLen("image", d.BitsLeft(), &imageGroup, nil)
		}
	case "cdsc":
		if irefBox, ok := ctx.parent().data.(*irefBox); ok {
			decodeBoxIrefEntry(irefBox, d)
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	case "irot":
		d.FieldU8("rotation", scalar.UintMapSymUint{
			0: 0,
			1: 90,
			2: 180,
			3: 270,
		})
	case "hnti":
		decodeBoxes(ctx, d)
	case "hint":
		decodeBoxes(ctx, d)
	case "pdin":
		d.FieldU8("version")
		d.FieldU24("flags")
		d.FieldArray("entries", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("entry", func(d *decode.D) {
					d.FieldU32("rate")
					d.FieldU32("initial_delay")
				})
			}
		})
	case "cslg":
		version := d.FieldU8("version")
		d.FieldU24("flags")
		switch version {
		case 0:
			d.FieldS32("composition_to_dts_shift")
			d.FieldS32("least_decode_to_display_delta")
			d.FieldS32("greatest_decode_to_display_delta")
			d.FieldS32("composition_start_time")
			d.FieldS32("composition_end_time")
		case 1:
			d.FieldS64("composition_to_dts_shift")
			d.FieldS64("least_decode_to_display_delta")
			d.FieldS64("greatest_decode_to_display_delta")
			d.FieldS64("composition_start_time")
			d.FieldS64("composition_end_time")
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	case "emsg":
		// https://aomediacodec.github.io/id3-emsg/
		version := d.FieldU8("version")
		d.FieldU24("flags")
		var schemeIdUri string
		switch version {
		case 0:
			schemeIdUri = d.FieldUTF8Null("scheme_id_uri")
			d.FieldUTF8Null("value")
			d.FieldU32("timescale")
			d.FieldU32("presentation_time_delta")
			d.FieldU32("event_duration")
			d.FieldU32("id")
		case 1:
			d.FieldU32("timescale")
			d.FieldU64("presentation_time")
			d.FieldU32("event_duration")
			d.FieldU32("id")
			schemeIdUri = d.FieldUTF8Null("scheme_id_uri")
			d.FieldUTF8Null("value")
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
		switch schemeIdUri {
		case "https://aomedia.org/emsg/ID3":
			d.FieldFormat("message_data", &id3v2Group, nil)
		default:
			d.FieldRawLen("message_data", d.BitsLeft())
		}
	case "jp2h":
		decodeBoxes(ctx, d)
	case "ihdr":
		d.FieldU32("width")
		d.FieldU32("height")
		d.FieldU16("nc")
		d.FieldU8("bpc", scalar.UintActualAdd(1))
		d.FieldU8("c", scalar.UintMapSymStr{
			0: "uncompressed",
			1: "modified_huffman",
			2: "modified_read",
			3: "modified_modified_read",
			4: "jbig_bi_level",
			5: "jpeg",
			6: "jpeg_ls",
			7: "jpeg_2000",
			8: "jbig2",
			9: "jbig",
		})
		d.FieldU8("unk_c")
		d.FieldU8("ipr")
	case "jP":
		d.FieldRawLen("signature", 4*8, d.AssertBitBuf([]byte{0x0d, 0x0a, 0x87, 0x0a}))
	case "jp2c":
		d.FieldFormat("segments", &jp2cGroup, nil)
	case "uinf":
		decodeBoxes(ctx, d)
	case "ulst":
		nu := d.FieldU16("nu")
		d.FieldArray("uids", func(d *decode.D) {
			for i := 0; i < int(nu); i++ {
				d.FieldRawLen("uid", 128)
			}
		})

	default:
		// there are at least 4 ways to encode udta metadata in mov/mp4 files.
		//
		// mdta subtype:
		//
		// udta:
		//   meta
		//     hdlr with subtype "mdta"
		//     keys with 1-based <index> to key namespace.name table
		//     ilst
		//       <index>-box (box type is 32bit BE 1-based number into table above)
		//         data box with value
		//
		// mdir subtype:
		//
		// udta
		//   meta
		//     hdlr with subtype "mdir"
		//     ilst
		//       ©<abc> or similar
		//         data with value
		//
		// no-meta-box with length and language:
		//
		// udta
		//   ©<abc> or similar
		//     value length and language
		//
		// no-meta-box value rest of box:
		//
		// udta
		//   <name>
		//     value rest of box
		if mb := ctx.currentMetaBox(); mb != nil && ctx.parent().typ == "ilst" {
			// unknown type under a meta box with ilst as parent, decode as boxes
			// is probably one or more data boxes
			decodeBoxes(ctx, d)
		} else if ctx.parent().typ == "udta" {
			// TODO: better probe? ffmpeg uses box name heuristics?
			// if 16 length field seems to match assume box with length, language and value
			// otherwise decode as box with value rest of box

			// only probe if there is something
			probeLength := int64(0)
			if d.BitsLeft() >= 16 {
				probeLength = int64(d.PeekUintBits(16))
			}
			// +2 for length field, +2 for language field
			if (probeLength+2+2)*8 == d.BitsLeft() {
				length := d.FieldU16("length")
				d.FieldStrFn("language", decodeLang)
				d.FieldUTF8("value", int(length))
			} else {
				d.FieldUTF8("value", int(d.BitsLeft()/8))
			}
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	}
}
