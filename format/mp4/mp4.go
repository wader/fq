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
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

//go:embed *.jq
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

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MP4,
		Description: "MPEG-4 file and similar",
		Groups: []string{
			format.PROBE,
			format.IMAGE, // avif
		},
		DecodeFn: mp4Decode,
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
		},
		Files: mp4FS,
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
	dataFormat     string
	originalFormat string
}

type track struct {
	id                 uint32
	sampleDescriptions []sampleDescription
	subType            string
	stco               []uint64 //
	stsc               []stsc
	stsz               []uint32
	formatInArg        interface{}
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
	d.AssertLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Fatalf("first box size too small < 8")
	}
	firstType := d.UTF8(4)
	switch firstType {
	case "styp", "ftyp", "free", "moov":
	default:
		d.Errorf("no styp, ftyp, free or moov box found")
	}

	d.SeekRel(-8 * 8)

	decodeBoxes(ctx, d)

	// keep track order stable
	var sortedTracks []*track
	for _, t := range ctx.tracks {
		sortedTracks = append(sortedTracks, t)
	}
	sort.Slice(sortedTracks, func(i, j int) bool { return sortedTracks[i].id < sortedTracks[j].id })

	d.FieldArray("tracks", func(d *decode.D) {
		for _, t := range sortedTracks {
			decodeSampleRange := func(d *decode.D, t *track, dataFormat string, name string, firstBit int64, nBits int64, inArg interface{}) {
				d.RangeFn(firstBit, nBits, func(d *decode.D) {
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
				// TODO: handle progressive/fragmented mp4 differently somehow?

				trackSdDataFormat := "unknown"
				if len(t.sampleDescriptions) > 0 {
					sd := t.sampleDescriptions[0]
					trackSdDataFormat = sd.dataFormat
					if sd.originalFormat != "" {
						trackSdDataFormat = sd.originalFormat
					}
				}

				d.FieldArray("samples", func(d *decode.D) {
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
							decodeSampleRange(d, t, trackSdDataFormat, "sample", int64(sampleOffset)*8, int64(sampleSize)*8, t.formatInArg)

							// log.Printf("%s %d/%d %d/%d sample=%d/%d chunk=%d size=%d %d-%d\n",
							// 	trackSdDataFormat, stscIndex, len(t.stsc),
							// 	i, stscEntry.samplesPerChunk,
							// 	sampleNr, len(t.stsz),
							// 	chunkNr,
							// 	sampleSize,
							// 	sampleOffset,
							// 	sampleOffset+uint64(sampleSize))

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

							dataFormat := trackSdDataFormat
							if m.defaultSampleDescriptionIndex != 0 && int(m.defaultSampleDescriptionIndex-1) < len(t.sampleDescriptions) {
								sd := t.sampleDescriptions[m.defaultSampleDescriptionIndex-1]
								dataFormat = sd.dataFormat
								if sd.originalFormat != "" {
									dataFormat = sd.originalFormat
								}
							}

							// log.Printf("moof %#+v dataFormat: %#+v\n", m, dataFormat)

							decodeSampleRange(d, t, dataFormat, "sample", sampleOffset*8, int64(sz)*8, t.formatInArg)
							sampleOffset += int64(sz)
						}
					}
				})
			})
		}
	})

	return nil

}
