package mp4

// Tries to decode both ISOBMFF and quicktime
// Uses naming from ISOBMFF when possible
// ISO/IEC 14496-12
// Quicktime file format https://developer.apple.com/standards/qtff-2001.pdf
// FLAC in ISOBMFF https://github.com/xiph/flac/blob/master/doc/isoflac.txt
// https://www.webmproject.org/vp9/mp4/

// TODO: validate structure better? trak/stco etc
// TODO: fmp4, default samples sizes etc
// TODO: keep track of structure somehow to detect errors
// TODO: ISO-14496 says mp4 mdat can begin and end with original header/trailer (no used i guess?)
// TODO: more metadata
// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/Metadata/Metadata.html#//apple_ref/doc/uid/TP40000939-CH1-SW43
// TODO: split into mov and mp4 decoder?
// TODO: split into mp4_box decoder? needs complex in/out args?
// TODO: fragmented: tracks per fragment? fragment_index in samples?
// TODO: better probe, find first 2 boxes, should be free,ftyp or mdat?
// TODO: mime

import (
	"embed"
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

//go:embed *.jq
var mp4FS embed.FS

var aacFrameFormat []*decode.Format
var av1CCRFormat []*decode.Format
var av1FrameFormat []*decode.Format
var flacFrameFormat []*decode.Format
var flacMetadatablockFormat []*decode.Format
var id3v2Format []*decode.Format
var imageFormat []*decode.Format
var jpegFormat []*decode.Format
var mp3FrameFormat []*decode.Format
var mpegAVCAUFormat []*decode.Format
var mpegAVCDCRFormat []*decode.Format
var mpegESFormat []*decode.Format
var mpegHEVCDCRFrameFormat []*decode.Format
var mpegHEVCSampleFormat []*decode.Format
var mpegPESPacketSampleFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var protoBufWidevineFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp9FrameFormat []*decode.Format
var vpxCCRFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MP4,
		Description: "MPEG-4 file and similar",
		Groups:      []string{format.PROBE},
		DecodeFn:    mp4Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.AV1_CCR}, Formats: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Formats: &av1FrameFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.ID3V2}, Formats: &id3v2Format},
			{Names: []string{format.IMAGE}, Formats: &imageFormat},
			{Names: []string{format.JPEG}, Formats: &jpegFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3FrameFormat},
			{Names: []string{format.AVC_AU}, Formats: &mpegAVCAUFormat},
			{Names: []string{format.AVC_DCR}, Formats: &mpegAVCDCRFormat},
			{Names: []string{format.MPEG_ES}, Formats: &mpegESFormat},
			{Names: []string{format.HEVC_AU}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.HEVC_DCR}, Formats: &mpegHEVCDCRFrameFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Formats: &mpegPESPacketSampleFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.PROTOBUF_WIDEVINE}, Formats: &protoBufWidevineFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
			{Names: []string{format.VPX_CCR}, Formats: &vpxCCRFormat},
		},
		FS: mp4FS,
	})
}

type stsc struct {
	firstChunk      uint32
	samplesPerChunk uint32
}

type moof struct {
	offset                        int64
	defaultSampleSize             uint32
	defaultSampleDescriptionIndex uint32
	dataOffset                    uint32
	samplesSizes                  []uint32
}

type sampleDescription struct {
	dataFormat string
}

type track struct {
	id                 uint32
	sampleDescriptions []sampleDescription
	subType            string
	stco               []uint64 //
	stsc               []stsc
	stsz               []uint32
	decodeOpts         []decode.Options
	objectType         int // if data format is "mp4a"

	moofs       []*moof // for fmp4
	currentMoof *moof
}

type decodeContext struct {
	path              []string
	tracks            map[uint32]*track
	currentTrack      *track
	currentMoofOffset int64
}

func isParent(ctx *decodeContext, typ string) bool {
	return len(ctx.path) >= 2 && ctx.path[len(ctx.path)-2] == typ
}

func mp4Decode(d *decode.D, in interface{}) interface{} {
	ctx := &decodeContext{
		tracks: map[uint32]*track{},
	}

	// TODO: nicer, validate functions without field?
	d.ValidateAtLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Invalid("first box size too small < 8")
	}
	firstType := d.UTF8(4)
	switch firstType {
	case "styp", "ftyp", "free", "moov":
	default:
		d.Invalid("no styp, ftyp, free or moov box found")
	}

	d.SeekRel(-8 * 8)

	decodeBoxes(ctx, d)

	// keep track order stable
	var sortedTracks []*track
	for _, t := range ctx.tracks {
		sortedTracks = append(sortedTracks, t)
	}
	sort.Slice(sortedTracks, func(i, j int) bool { return sortedTracks[i].id < sortedTracks[j].id })

	d.FieldArrayFn("tracks", func(d *decode.D) {
		for _, t := range sortedTracks {
			decodeSampleRange := func(d *decode.D, t *track, dataFormat string, name string, firstBit int64, nBits int64, opts ...decode.Options) {
				switch dataFormat {
				case "fLaC":
					d.FieldFormatRange(name, firstBit, nBits, flacFrameFormat, opts...)
				case "Opus":
					d.FieldFormatRange(name, firstBit, nBits, opusPacketFrameFormat, opts...)
				case "vp09":
					d.FieldFormatRange(name, firstBit, nBits, vp9FrameFormat, opts...)
				case "avc1":
					d.FieldFormatRange(name, firstBit, nBits, mpegAVCAUFormat, opts...)
				case "hev1", "hvc1":
					d.FieldFormatRange(name, firstBit, nBits, mpegHEVCSampleFormat, opts...)
				case "av01":
					d.FieldFormatRange(name, firstBit, nBits, av1FrameFormat, opts...)
				case "mp4a":
					switch t.objectType {
					case format.MPEGObjectTypeMP3:
						d.FieldFormatRange(name, firstBit, nBits, mp3FrameFormat, opts...)
					case format.MPEGObjectTypeAAC:
						// TODO: MPEGObjectTypeAACLow, Main etc?
						d.FieldFormatRange(name, firstBit, nBits, aacFrameFormat, opts...)
					case format.MPEGObjectTypeVORBIS:
						d.FieldFormatRange(name, firstBit, nBits, vorbisPacketFormat, opts...)
					default:
						d.FieldBitBufRange(name, firstBit, nBits)
					}
				case "mp4v":
					switch t.objectType {
					case format.MPEGObjectTypeMPEG2VideoMain:
						d.FieldFormatRange(name, firstBit, nBits, mpegPESPacketSampleFormat, opts...)
					case format.MPEGObjectTypeMJPEG:
						d.FieldFormatRange(name, firstBit, nBits, jpegFormat, opts...)
					default:
						d.FieldBitBufRange(name, firstBit, nBits)
					}
				case "jpeg":
					d.FieldFormatRange(name, firstBit, nBits, jpegFormat, opts...)
				default:
					d.FieldBitBufRange(name, firstBit, nBits)
				}
			}

			d.FieldStructFn("track", func(d *decode.D) {
				// TODO: handle progressive/fragmented mp4 differently somehow?
				if t.moofs == nil && len(t.sampleDescriptions) > 0 {
					d.FieldStrFn("data_format", func() (string, string) { return t.sampleDescriptions[0].dataFormat, "" })
				}

				d.FieldArrayFn("samples", func(d *decode.D) {
					stscIndex := 0
					chunkNr := uint32(0)
					sampleNr := uint64(0)

					for sampleNr < uint64(len(t.stsz)) {
						if stscIndex >= len(t.stsc) {
							// TODO: add warning
							break
						}
						stscEntry := t.stsc[stscIndex]
						if int(chunkNr) >= len(t.stco) {
							// TODO: add warning
							break
						}
						sampleOffset := t.stco[chunkNr]

						for i := uint32(0); i < stscEntry.samplesPerChunk; i++ {
							if int(sampleNr) >= len(t.stsz) {
								// TODO: add warning
								break
							}

							sampleSize := t.stsz[sampleNr]
							dataFormat := "unknown"
							if len(t.sampleDescriptions) > 0 {
								dataFormat = t.sampleDescriptions[0].dataFormat
							}

							decodeSampleRange(d, t, dataFormat, "sample", int64(sampleOffset)*8, int64(sampleSize)*8, t.decodeOpts...)

							// log.Printf("%s %d/%d %d/%d sample=%d/%d chunk=%d size=%d %d-%d\n", t.dataFormat, stscIndex, len(t.stsc), i, stscEntry.samplesPerChunk, sampleNr, len(t.stsz), chunkNr, sampleSize, sampleOffset, sampleOffset+uint64(sampleSize))

							sampleOffset += uint64(sampleSize)
							sampleNr++

						}

						chunkNr++
						if stscIndex < len(t.stsc)-1 && chunkNr >= t.stsc[stscIndex+1].firstChunk-1 {
							stscIndex++
						}
					}

					for _, m := range t.moofs {
						sampleOffset := m.offset + int64(m.dataOffset)
						for _, sz := range m.samplesSizes {
							// log.Printf("moof sample %s %d-%d\n", t.dataFormat, sampleOffset, int64(sz))

							dataFormat := "unknown"
							if len(t.sampleDescriptions) > 0 {
								dataFormat = t.sampleDescriptions[0].dataFormat
							}
							if m.defaultSampleDescriptionIndex != 0 && int(m.defaultSampleDescriptionIndex-1) < len(t.sampleDescriptions) {
								dataFormat = t.sampleDescriptions[m.defaultSampleDescriptionIndex-1].dataFormat
							}

							// log.Printf("moof %#+v dataFormat: %#+v\n", m, dataFormat)

							decodeSampleRange(d, t, dataFormat, "sample", sampleOffset*8, int64(sz)*8, t.decodeOpts...)
							sampleOffset += int64(sz)
						}
					}
				})
			})
		}
	})

	return nil

}
