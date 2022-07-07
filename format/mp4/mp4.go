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
	"embed"
	"fmt"
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

//go:embed mp4.jq
var mp4FS embed.FS

var aacFrameFormat decode.Group
var av1CCRFormat decode.Group
var av1FrameFormat decode.Group
var flacFrameFormat decode.Group
var flacMetadatablocksFormat decode.Group
var id3v2Format decode.Group
var imageFormat decode.Group
var jpegFormat decode.Group
var mp3FrameFormat decode.Group
var mpegAVCAUFormat decode.Group
var mpegAVCDCRFormat decode.Group
var mpegESFormat decode.Group
var mpegHEVCDCRFrameFormat decode.Group
var mpegHEVCSampleFormat decode.Group
var mpegPESPacketSampleFormat decode.Group
var opusPacketFrameFormat decode.Group
var protoBufWidevineFormat decode.Group
var psshPlayreadyFormat decode.Group
var vorbisPacketFormat decode.Group
var vp9FrameFormat decode.Group
var vpxCCRFormat decode.Group
var iccProfileFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MP4,
		Description: "ISOBMFF MPEG-4 part 12 and similar",
		Groups: []string{
			format.PROBE,
			format.IMAGE, // avif
		},
		DecodeFn: mp4Decode,
		DecodeInArg: format.Mp4In{
			DecodeSamples:  true,
			AllowTruncated: false,
		},
		Dependencies: []decode.Dependency{
			{Names: []string{format.AAC_FRAME}, Group: &aacFrameFormat},
			{Names: []string{format.AV1_CCR}, Group: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Group: &av1FrameFormat},
			{Names: []string{format.FLAC_FRAME}, Group: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCKS}, Group: &flacMetadatablocksFormat},
			{Names: []string{format.ID3V2}, Group: &id3v2Format},
			{Names: []string{format.IMAGE}, Group: &imageFormat},
			{Names: []string{format.JPEG}, Group: &jpegFormat},
			{Names: []string{format.MP3_FRAME}, Group: &mp3FrameFormat},
			{Names: []string{format.AVC_AU}, Group: &mpegAVCAUFormat},
			{Names: []string{format.AVC_DCR}, Group: &mpegAVCDCRFormat},
			{Names: []string{format.MPEG_ES}, Group: &mpegESFormat},
			{Names: []string{format.HEVC_AU}, Group: &mpegHEVCSampleFormat},
			{Names: []string{format.HEVC_DCR}, Group: &mpegHEVCDCRFrameFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Group: &mpegPESPacketSampleFormat},
			{Names: []string{format.OPUS_PACKET}, Group: &opusPacketFrameFormat},
			{Names: []string{format.PROTOBUF_WIDEVINE}, Group: &protoBufWidevineFormat},
			{Names: []string{format.PSSH_PLAYREADY}, Group: &psshPlayreadyFormat},
			{Names: []string{format.VORBIS_PACKET}, Group: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Group: &vp9FrameFormat},
			{Names: []string{format.VPX_CCR}, Group: &vpxCCRFormat},
			{Names: []string{format.ICC_PROFILE}, Group: &iccProfileFormat},
		},
		Files:     mp4FS,
		Functions: []string{"_help"},
	})
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
}

type pathEntry struct {
	typ  string
	data any
}

type decodeContext struct {
	opts   format.Mp4In
	path   []pathEntry
	tracks map[int]*track
}

func (ctx *decodeContext) lookupTrack(id int) *track {
	t, ok := ctx.tracks[id]
	if !ok {
		t = &track{id: id}
		ctx.tracks[id] = t
	}
	return t
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

func (ctx *decodeContext) currentTrack() *track {
	if t := ctx.currentTrakBox(); t != nil {
		return ctx.lookupTrack(t.trackID)
	}
	if t := ctx.currentTrafBox(); t != nil {
		return ctx.lookupTrack(t.trackID)
	}
	return nil
}

func mp4Tracks(d *decode.D, ctx *decodeContext) {
	// keep track order stable
	var sortedTracks []*track
	for _, t := range ctx.tracks {
		sortedTracks = append(sortedTracks, t)
	}
	sort.Slice(sortedTracks, func(i, j int) bool { return sortedTracks[i].id < sortedTracks[j].id })

	d.FieldArray("tracks", func(d *decode.D) {
		d.RangeSorted = false

		for _, t := range sortedTracks {
			decodeSampleRange := func(d *decode.D, t *track, decodeSample bool, dataFormat string, name string, firstBit int64, nBits int64, inArg any) {
				d.RangeFn(firstBit, nBits, func(d *decode.D) {
					if !decodeSample {
						d.FieldRawLen(name, d.BitsLeft())
						return
					}

					switch {
					case dataFormat == "fLaC":
						d.FieldFormatLen(name, nBits, flacFrameFormat, inArg)
					case dataFormat == "Opus":
						d.FieldFormatLen(name, nBits, opusPacketFrameFormat, inArg)
					case dataFormat == "vp09":
						d.FieldFormatLen(name, nBits, vp9FrameFormat, inArg)
					case dataFormat == "avc1":
						d.FieldFormatLen(name, nBits, mpegAVCAUFormat, inArg)
					case dataFormat == "hev1",
						dataFormat == "hvc1":
						d.FieldFormatLen(name, nBits, mpegHEVCSampleFormat, inArg)
					case dataFormat == "av01":
						d.FieldFormatLen(name, nBits, av1FrameFormat, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeMP3:
						d.FieldFormatLen(name, nBits, mp3FrameFormat, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeAAC:
						d.FieldFormatLen(name, nBits, aacFrameFormat, inArg)
					case dataFormat == "mp4a" && t.objectType == format.MPEGObjectTypeVORBIS:
						d.FieldFormatLen(name, nBits, vorbisPacketFormat, inArg)
					case dataFormat == "mp4v" && t.objectType == format.MPEGObjectTypeMPEG2VideoMain:
						d.FieldFormatLen(name, nBits, mpegPESPacketSampleFormat, inArg)
					case dataFormat == "mp4v" && t.objectType == format.MPEGObjectTypeMJPEG:
						d.FieldFormatLen(name, nBits, jpegFormat, inArg)
					case dataFormat == "jpeg":
						d.FieldFormatLen(name, nBits, jpegFormat, inArg)
					default:
						d.FieldRawLen(name, d.BitsLeft())
					}
				})
			}

			d.FieldStruct("track", func(d *decode.D) {
				d.FieldValueU("id", uint64(t.id))

				trackSDDataFormat := "unknown"
				if len(t.sampleDescriptions) > 0 {
					sd := t.sampleDescriptions[0]
					trackSDDataFormat = sd.dataFormat
					if sd.originalFormat != "" {
						trackSDDataFormat = sd.originalFormat
					}
				}

				d.FieldValueStr("data_foramt", trackSDDataFormat)

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

							// log.Println(logStrFn())

							decodeSampleRange(d, t, ctx.opts.DecodeSamples, trackSDDataFormat, "sample", sampleOffset*8, stszEntry.size*8, t.formatInArg)

							sampleOffset += stszEntry.size
							stscEntryNr++
							stszEntryNr++
							sampleNr++
						}
					}

					sampleNr := 0
					for _, m := range t.moofs {
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
								// log.Println(logStrFn())

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
				})
			})
		}
	})
}

func mp4Decode(d *decode.D, in any) any {
	mi, _ := in.(format.Mp4In)

	ctx := &decodeContext{
		opts:   mi,
		path:   []pathEntry{{typ: "root"}},
		tracks: map[int]*track{},
	}

	// TODO: nicer, validate functions without field?
	d.AssertLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Fatalf("first box size too small < 8")
	}
	firstType := d.UTF8(4)
	switch firstType {
	case "styp", // mp4 segment
		"ftyp", // mp4 file
		"free", // seems to happen
		"moov", // seems to happen
		"pnot", // video preview file
		"jP  ": // JPEG 2000
	default:
		d.Errorf("no styp, ftyp, free or moov box found")
	}

	d.SeekRel(-8 * 8)

	decodeBoxes(ctx, d)
	if len(ctx.tracks) > 0 {
		mp4Tracks(d, ctx)
	}

	return nil
}
