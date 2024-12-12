package mp4

// Tries to decode ISOBMFF quicktime mov
// Uses naming from ISOBMFF when possible
// ISO/IEC 14496-12
// Quicktime file format https://developer.apple.com/standards/qtff-2001.pdf
// FLAC in ISOBMFF https://github.com/xiph/flac/blob/master/doc/isoflac.txt
// vp9 in ISOBMFF https://www.webmproject.org/vp9/mp4/
// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/Metadata/Metadata.html#//apple_ref/doc/uid/TP40000939-CH1-SW43

// TODO: validate structure better? trak/stco etc
// TODO: keep track of structure somehow to detect errors
// TODO: ISO-14496 says mp4 mdat can begin and end with original header/trailer (no used?)
// TODO: split into mov and mp4 decoder?
// TODO: split into mp4_box decoder? needs complex in/out args?
// TODO: better probe, find first 2 boxes, should be free,ftyp or mdat?

import (
	"cmp"
	"embed"
	"fmt"
	"slices"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed mp4.jq
//go:embed mp4.md
var mp4FS embed.FS

var aacFrameGroup decode.Group
var av1CCRGroup decode.Group
var av1FrameGroup decode.Group
var avcAUGroup decode.Group
var avcDCRGroup decode.Group
var flacFrameGroup decode.Group
var flacMetadatablocksGroup decode.Group
var hevcAUGroup decode.Group
var hevcCDCRGroup decode.Group
var iccProfileGroup decode.Group
var id3v2Group decode.Group
var imageGroup decode.Group
var jp2cGroup decode.Group
var jpegGroup decode.Group
var mp3FrameGroup decode.Group
var mpegESGroup decode.Group
var mpegPESPacketSampleGroup decode.Group
var opusPacketFrameGroup decode.Group
var pngGroup decode.Group
var proResFrameGroup decode.Group
var protoBufWidevineGroup decode.Group
var psshPlayreadyGroup decode.Group
var vorbisPacketGroup decode.Group
var vp9FrameGroup decode.Group
var vpxCCRGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.MP4,
		&decode.Format{
			Description: "ISOBMFF, QuickTime and similar",
			Groups: []*decode.Group{
				format.Probe,
				format.Image, // avif
			},
			DecodeFn: mp4Decode,
			DefaultInArg: format.MP4_In{
				DecodeSamples:  true,
				AllowTruncated: false,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AAC_Frame}, Out: &aacFrameGroup},
				{Groups: []*decode.Group{format.AV1_CCR}, Out: &av1CCRGroup},
				{Groups: []*decode.Group{format.AV1_Frame}, Out: &av1FrameGroup},
				{Groups: []*decode.Group{format.AVC_AU}, Out: &avcAUGroup},
				{Groups: []*decode.Group{format.AVC_DCR}, Out: &avcDCRGroup},
				{Groups: []*decode.Group{format.FLAC_Frame}, Out: &flacFrameGroup},
				{Groups: []*decode.Group{format.FLAC_Metadatablocks}, Out: &flacMetadatablocksGroup},
				{Groups: []*decode.Group{format.HEVC_AU}, Out: &hevcAUGroup},
				{Groups: []*decode.Group{format.HEVC_DCR}, Out: &hevcCDCRGroup},
				{Groups: []*decode.Group{format.ICC_Profile}, Out: &iccProfileGroup},
				{Groups: []*decode.Group{format.ID3v2}, Out: &id3v2Group},
				{Groups: []*decode.Group{format.Image}, Out: &imageGroup},
				{Groups: []*decode.Group{format.JP2C}, Out: &jp2cGroup},
				{Groups: []*decode.Group{format.JPEG}, Out: &jpegGroup},
				{Groups: []*decode.Group{format.MP3_Frame}, Out: &mp3FrameGroup},
				{Groups: []*decode.Group{format.MPEG_ES}, Out: &mpegESGroup},
				{Groups: []*decode.Group{format.MPEG_PES_Packet}, Out: &mpegPESPacketSampleGroup},
				{Groups: []*decode.Group{format.Opus_Packet}, Out: &opusPacketFrameGroup},
				{Groups: []*decode.Group{format.PNG}, Out: &pngGroup},
				{Groups: []*decode.Group{format.Prores_Frame}, Out: &proResFrameGroup},
				{Groups: []*decode.Group{format.ProtobufWidevine}, Out: &protoBufWidevineGroup},
				{Groups: []*decode.Group{format.PSSH_Playready}, Out: &psshPlayreadyGroup},
				{Groups: []*decode.Group{format.Vorbis_Packet}, Out: &vorbisPacketGroup},
				{Groups: []*decode.Group{format.VP9_Frame}, Out: &vp9FrameGroup},
				{Groups: []*decode.Group{format.VPX_CCR}, Out: &vpxCCRGroup},
			},
		})
	interp.RegisterFS(mp4FS)
}

type stsc struct {
	firstChunk      int
	samplesPerChunk int
}

type moof struct {
	offset                        int64
	defaultSampleSize             int64
	defaultSampleDescriptionIndex int
	truns                         []trun
	sencs                         []senc
}

// TODO: nothing for now
type senc struct {
	entries []struct{}
}

type trun struct {
	dataOffset   int64
	samplesSizes []int64
}

type sampleDescription struct {
	dataFormat     string
	originalFormat string
}

type stsz struct {
	size  int64
	count int
}

type track struct {
	seenHdlr           bool
	fragment           bool
	id                 int
	sampleDescriptions []sampleDescription
	subType            string
	stco               []int64
	stsc               []stsc
	stsz               []stsz
	formatInArg        any
	objectType         int // if data format is "mp4a"
	defaultIVSize      int
	moofs              []*moof // for fmp4
	dref               bool
	drefURL            string
}

type pathEntry struct {
	typ  string
	data any
}

type decodeContext struct {
	opts   format.MP4_In
	path   []pathEntry
	tracks []*track
}

func (ctx *decodeContext) isParent(typ string) bool {
	return ctx.parent().typ == typ
}

func (ctx *decodeContext) parent() pathEntry {
	return ctx.path[len(ctx.path)-2]
}

func (ctx *decodeContext) findParent(typ string) any {
	for i := len(ctx.path) - 1; i >= 0; i-- {
		p := ctx.path[i]
		if p.typ == typ {
			return p.data
		}
	}
	return nil
}

func (ctx *decodeContext) rootBox() *rootBox {
	t, _ := ctx.findParent("").(*rootBox)
	return t
}

func (ctx *decodeContext) currentTrakBox() *trakBox {
	t, _ := ctx.findParent("trak").(*trakBox)
	return t
}

func (ctx *decodeContext) currentTrafBox() *trafBox {
	t, _ := ctx.findParent("traf").(*trafBox)
	return t
}

func (ctx *decodeContext) currentMoofBox() *moofBox {
	t, _ := ctx.findParent("moof").(*moofBox)
	return t
}

func (ctx *decodeContext) currentMetaBox() *metaBox {
	t, _ := ctx.findParent("meta").(*metaBox)
	return t
}

func (ctx *decodeContext) currentTrack() *track {
	if t := ctx.currentTrakBox(); t != nil {
		return t.track
	}
	if t := ctx.currentTrafBox(); t != nil {
		return t.track
	}
	return nil
}

func mp4Tracks(d *decode.D, ctx *decodeContext) {
	type trackCollected struct {
		track  *track
		order  int
		moofss [][]*moof
	}

	var tracksCollected []*trackCollected
	tracksCollectedSeen := map[int]*trackCollected{}
	for i, t := range ctx.tracks {
		tc, ok := tracksCollectedSeen[t.id]
		if !ok {
			tc = &trackCollected{
				order: i,
				track: t,
			}
			tracksCollectedSeen[t.id] = tc
			tracksCollected = append(tracksCollected, tc)
		}

		// TODO: error if not fragmented and seen before?

		tc.moofss = append(tc.moofss, t.moofs)
	}

	// sort by id then order in file
	slices.SortStableFunc(tracksCollected, func(a, b *trackCollected) int {
		if r := cmp.Compare(a.track.id, b.track.id); r != 0 {
			return r
		}
		return cmp.Compare(a.order, b.order)
	})

	d.FieldArray("tracks", func(d *decode.D) {
		for _, tc := range tracksCollected {
			decodeSampleRange := func(d *decode.D, t *track, decodeSample bool, dataFormat string, name string, firstBit int64, nBits int64, inArg any) {
				d.RangeFn(firstBit, nBits, func(d *decode.D) {
					if !decodeSample {
						d.FieldRawLen(name, d.BitsLeft())
						return
					}

					switch {
					case dataFormat == "fLaC":
						d.FieldFormatLen(name, nBits, &flacFrameGroup, inArg)
					case dataFormat == "Opus":
						d.FieldFormatLen(name, nBits, &opusPacketFrameGroup, inArg)
					case dataFormat == "vp09":
						d.FieldFormatLen(name, nBits, &vp9FrameGroup, inArg)
					case dataFormat == "avc1":
						d.FieldFormatLen(name, nBits, &avcAUGroup, inArg)
					case dataFormat == "hev1",
						dataFormat == "hvc1":
						d.FieldFormatLen(name, nBits, &hevcAUGroup, inArg)
					case dataFormat == "av01":
						d.FieldFormatLen(name, nBits, &av1FrameGroup, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeMP3:
						d.FieldFormatLen(name, nBits, &mp3FrameGroup, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeAAC:
						d.FieldFormatLen(name, nBits, &aacFrameGroup, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeVORBIS:
						d.FieldFormatLen(name, nBits, &vorbisPacketGroup, inArg)
					case dataFormat == "mp4v" && t.objectType == format.MPEGObjectTypeMPEG2VideoMain:
						d.FieldFormatLen(name, nBits, &mpegPESPacketSampleGroup, inArg)
					case dataFormat == "mp4v" && t.objectType == format.MPEGObjectTypeMJPEG:
						d.FieldFormatLen(name, nBits, &jpegGroup, inArg)
					case dataFormat == "mp4v" && t.objectType == format.MPEGObjectTypePNG:
						d.FieldFormatLen(name, nBits, &pngGroup, inArg)
					case dataFormat == "jpeg":
						d.FieldFormatLen(name, nBits, &jpegGroup, inArg)
					case dataFormat == "apch",
						dataFormat == "apcn",
						dataFormat == "scpa",
						dataFormat == "apco",
						dataFormat == "ap4h":
						d.FieldFormatLen(name, nBits, &proResFrameGroup, inArg)
					default:
						d.FieldRawLen(name, d.BitsLeft())
					}
				})
			}

			d.FieldStruct("track", func(d *decode.D) {
				t := tc.track

				d.FieldValueUint("id", uint64(t.id))

				trackSDDataFormat := "unknown"
				if len(t.sampleDescriptions) > 0 {
					sd := t.sampleDescriptions[0]
					trackSDDataFormat = sd.dataFormat
					if sd.originalFormat != "" {
						trackSDDataFormat = sd.originalFormat
					}
				}
				d.FieldValueStr("data_format", trackSDDataFormat, dataFormatNames)

				if t.dref && t.drefURL != "" {
					d.FieldValueStr("data_reference_url", t.drefURL)
					return
				}

				switch trackSDDataFormat {
				case "lpcm",
					"raw ",
					"twos",
					"sowt",
					"in24",
					"in32",
					"fl23",
					"fl64",
					"alaw",
					"ulaw":
					// TODO: treat raw samples format differently, a bit too much to have one field per sample.
					// maybe in some future fq could have smart array fields
					return
				}

				d.FieldArray("samples", func(d *decode.D) {
					// TODO: warning? could also be init fragment etc

					if len(t.stsz) > 0 && len(t.stsc) > 0 && len(t.stco) > 0 {
						stszIndex := 0
						stszEntryNr := 0
						sampleNr := 0
						stscIndex := 0
						stscEntryNr := 0
						stcoIndex := 0

						stszEntry := t.stsz[stszIndex]
						stscEntry := t.stsc[stscIndex]
						sampleOffset := t.stco[stcoIndex]

						logStrFn := func() string {
							return fmt.Sprintf("%d: %s: nr=%d: stsz[%d/%d] nr=%d %#v stsc[%d/%d] nr=%d %#v stco[%d/%d]=%d \n",
								t.id,
								trackSDDataFormat,
								sampleNr,
								stszIndex, len(t.stsz), stszEntryNr, stszEntry,
								stscIndex, len(t.stsc), stscEntryNr, stscEntry,
								stcoIndex, len(t.stco), sampleOffset,
							)
						}

						for stszIndex < len(t.stsz) {
							if stszEntryNr >= stszEntry.count {
								stszIndex++
								if stszIndex >= len(t.stsz) {
									// TODO: warning if unused stsc/stco entries?
									break
								}

								stszEntry = t.stsz[stszIndex]
								stszEntryNr = 0
							}

							if stscEntryNr >= stscEntry.samplesPerChunk {
								stscEntryNr = 0
								stcoIndex++
								if stcoIndex >= len(t.stco) {
									d.Fatalf("outside stco: %s", logStrFn())
								}
								sampleOffset = t.stco[stcoIndex]

								if stscIndex < len(t.stsc)-1 && stcoIndex >= t.stsc[stscIndex+1].firstChunk-1 {
									stscIndex++
									if stscIndex >= len(t.stsc) {
										d.Fatalf("outside stsc: %s", logStrFn())
									}
									stscEntry = t.stsc[stscIndex]
								}
							}

							decodeSampleRange(d, t, ctx.opts.DecodeSamples, trackSDDataFormat, "sample", sampleOffset*8, stszEntry.size*8, t.formatInArg)

							sampleOffset += stszEntry.size
							stscEntryNr++
							stszEntryNr++
							sampleNr++
						}
					}

					sampleNr := 0

					for _, ms := range tc.moofss {
						for _, m := range ms {
							for trunNr, trun := range m.truns {
								var senc senc
								if trunNr < len(m.sencs) {
									senc = m.sencs[trunNr]
								}
								sampleOffset := m.offset + trun.dataOffset

								for trunSampleNr, sz := range trun.samplesSizes {
									dataFormat := trackSDDataFormat
									if m.defaultSampleDescriptionIndex != 0 && m.defaultSampleDescriptionIndex-1 < len(t.sampleDescriptions) {
										sd := t.sampleDescriptions[m.defaultSampleDescriptionIndex-1]
										dataFormat = sd.dataFormat
										if sd.originalFormat != "" {
											dataFormat = sd.originalFormat
										}
									}

									// logStrFn := func() string {
									// 	return fmt.Sprintf("%d: %s: %d: (%s): sz=%d %d+%d=%d",
									// 		t.id,
									// 		dataFormat,
									// 		sampleNr,
									// 		trackSDDataFormat,
									// 		sz,
									// 		m.offset,
									// 		m.dataOffset,
									// 		sampleOffset,
									// 	)
									// }

									decodeSample := ctx.opts.DecodeSamples
									if trunSampleNr < len(senc.entries) {
										// TODO: encrypted
										decodeSample = false
									}

									decodeSampleRange(d, t, decodeSample, dataFormat, "sample", sampleOffset*8, sz*8, t.formatInArg)

									sampleOffset += sz
									sampleNr++
								}
							}
						}
					}
				})
			})
		}
	})
}

func mp4Decode(d *decode.D) any {
	var mi format.MP4_In
	d.ArgAs(&mi)

	ctx := &decodeContext{
		opts:   mi,
		path:   []pathEntry{{typ: "root"}},
		tracks: []*track{},
	}

	// TODO: nicer, validate functions without field?
	d.AssertLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Fatalf("first box size too small < 8")
	}
	firstType := strings.TrimSpace(d.UTF8(4))
	switch firstType {
	case "styp", // mp4 segment
		"ftyp", // mp4 file
		"free", // seems to happen
		"moov", // seems to happen
		"pnot", // video preview file
		"jP":   // JPEG 2000
	default:
		d.Errorf("no styp, ftyp, free or moov box found")
	}

	d.SeekRel(-8 * 8)

	ctx.path = []pathEntry{{typ: "", data: &rootBox{}}}

	decodeBoxes(ctx, d)
	if len(ctx.tracks) > 0 {
		mp4Tracks(d, ctx)
	}

	return nil
}
