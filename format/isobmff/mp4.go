package isobmff

// Tries to decode mp4 and quicktime mov
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

type sampleDescription struct {
	dataFormat     string
	originalFormat string
}

type stscEntry struct {
	firstChunk      int
	samplesPerChunk int
}

type stszEntry struct {
	size  int64
	count int
}

type tkhdBox struct {
	trackID int
}

type hdlrBox struct {
	subType string
}

type stsdBox struct {
	sampleDescriptions []sampleDescription
	numAudioChannels   uint64
}

type stscBox struct {
	entries []stscEntry
}

type stszBox struct {
	entries []stszEntry
}

type stcoBox struct {
	entries []int64
}

type drefBox struct {
	url string
}

type tencBox struct {
	defaultIVSize int
}

// Fragmented MP4 boxes
type tfhdBox struct {
	trackID                       int
	baseDataOffset                int64
	defaultSampleSize             int64
	defaultSampleDescriptionIndex int
}

type trunBox struct {
	dataOffset  int64
	sampleSizes []int64
}

type sencBox struct {
	entries []struct{}
}

type moofBox struct {
	offset int64
}

type ftypBox struct {
	majorBrand   string
	minorVersion uint32
	minorBrands  []string
}

type avcCBox struct {
	formatInArg format.AVC_AU_In
}

type hvcCBox struct {
	formatInArg format.HEVC_AU_In
}

type dfLaBox struct {
	formatInArg format.FLAC_Frame_In
}

type esdsBox struct {
	formatInArg format.AAC_Frame_In
	objectType  int
}

type frmaBox struct {
	originalFormat string
}

type trackCollected struct {
	trackID   int
	trakNode  *box
	order     int
	trafNodes []*box
}

func stblNodeFindFormatInArg(n *box) any {
	if v := findData[*avcCBox](n, ">>avcC"); v != nil {
		return v.formatInArg
	}
	if v := findData[*hvcCBox](n, ">>hvcC"); v != nil {
		return v.formatInArg
	}
	if v := findData[*dfLaBox](n, ">>dfLa"); v != nil {
		return v.formatInArg
	}
	if v := findData[*esdsBox](n, ">>esds"); v != nil {
		return v.formatInArg
	}
	return nil
}

func stblNodeFindObjectType(n *box) int {
	if v := findData[*esdsBox](n, ">>esds"); v != nil {
		return v.objectType
	}
	return 0
}

func findMajorBrand(root *box) string {
	if fb := findData[*ftypBox](root, "ftyp"); fb != nil {
		return fb.majorBrand
	}
	if fb := findData[*ftypBox](root, "styp"); fb != nil {
		return fb.majorBrand
	}
	return ""
}

type trafNodeData struct {
	moofOffset int64
	tfhd       *tfhdBox
	truns      []*trunBox
	sencs      []*sencBox
}

func trafNodeSampleData(tn *box) *trafNodeData {
	td := &trafNodeData{}
	if mb := findData[*moofBox](tn, "<moof"); mb != nil {
		td.moofOffset = mb.offset
	}
	if tb := findData[*tfhdBox](tn, "tfhd"); tb != nil {
		td.tfhd = tb
	}
	td.truns = slices.Collect(findAllData[*trunBox](tn, "trun"))
	td.sencs = slices.Collect(findAllData[*sencBox](tn, "senc"))
	return td
}

func mp4Tracks(d *decode.D, ctx *decodeContext, trakNodes, moofNodes []*box) {
	var tracksCollected []*trackCollected
	tracksCollectedSeen := map[int]*trackCollected{}

	// TODO: error on dup id or mixing fragmend and non-fragmented?

	for i, tn := range trakNodes {
		trackID := 0
		if tb := findData[*tkhdBox](tn, "tkhd"); tb != nil {
			trackID = tb.trackID
		}
		tc := &trackCollected{
			trackID:  trackID,
			order:    i,
			trakNode: tn,
		}
		tracksCollectedSeen[trackID] = tc
		tracksCollected = append(tracksCollected, tc)
	}

	for _, moofNode := range moofNodes {
		for trafNode := range moofNode.findAll("traf") {
			trackID := 0
			if tb := findData[*tfhdBox](trafNode, "tfhd"); tb != nil {
				trackID = tb.trackID
			}
			tc, ok := tracksCollectedSeen[trackID]
			if !ok {
				tc = &trackCollected{
					trackID: trackID,
					order:   len(tracksCollected),
				}
				tracksCollectedSeen[trackID] = tc
				tracksCollected = append(tracksCollected, tc)
			}
			tc.trafNodes = append(tc.trafNodes, trafNode)
		}
	}

	slices.SortStableFunc(tracksCollected, func(a, b *trackCollected) int {
		if r := cmp.Compare(a.trackID, b.trackID); r != 0 {
			return r
		}
		return cmp.Compare(a.order, b.order)
	})

	d.FieldArray("tracks", func(d *decode.D) {
		for _, tc := range tracksCollected {
			decodeSampleRange := func(d *decode.D, objectType int, decodeSample bool, dataFormat string, name string, firstBit int64, nBits int64, inArg any) {
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
					case dataFormat == "mp4a" && objectType == format.MPEGObjectTypeMP3:
						d.FieldFormatLen(name, nBits, &mp3FrameGroup, inArg)
					case dataFormat == "mp4a" && objectType == format.MPEGObjectTypeAAC:
						d.FieldFormatLen(name, nBits, &aacFrameGroup, inArg)
					case dataFormat == "mp4a" && objectType == format.MPEGObjectTypeVORBIS:
						d.FieldFormatLen(name, nBits, &vorbisPacketGroup, inArg)
					case dataFormat == "mp4v" && objectType == format.MPEGObjectTypeMPEG2VideoMain:
						d.FieldFormatLen(name, nBits, &mpegPESPacketSampleGroup, inArg)
					case dataFormat == "mp4v" && objectType == format.MPEGObjectTypeMJPEG:
						d.FieldFormatLen(name, nBits, &jpegGroup, inArg)
					case dataFormat == "mp4v" && objectType == format.MPEGObjectTypePNG:
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
				tn := tc.trakNode

				d.FieldValueUint("id", uint64(tc.trackID))

				var stblNode *box
				if tn != nil {
					stblNode = tn.find("mdia/minf/stbl")
				}

				trackSDDataFormat := "unknown"
				sampleDescriptions := []sampleDescription(nil)
				var formatInArg any
				objectType := 0
				if stblNode != nil {
					if sb := findData[*stsdBox](stblNode, "stsd"); sb != nil {
						sampleDescriptions = sb.sampleDescriptions
					}
					formatInArg = stblNodeFindFormatInArg(stblNode)
					objectType = stblNodeFindObjectType(stblNode)
				}

				if len(sampleDescriptions) > 0 {
					sd := sampleDescriptions[0]
					trackSDDataFormat = sd.dataFormat
					if sd.originalFormat != "" {
						trackSDDataFormat = sd.originalFormat
					}
				}
				d.FieldValueStr("data_format", trackSDDataFormat, dataFormatNames)

				if db := findData[*drefBox](tn, "minf/dinf/dref"); db != nil && db.url != "" {
					d.FieldValueStr("data_reference_url", db.url)
					return
				}

				switch trackSDDataFormat {
				case "lpcm",
					"ipcm",
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

				if ctx.opts.SkipSamples {
					return
				}

				d.FieldArray("samples", func(d *decode.D) {
					stsz := findData[*stszBox](stblNode, "stsz")
					if stsz == nil {
						stsz = findData[*stszBox](stblNode, "stz2")
					}
					stsc := findData[*stscBox](stblNode, "stsc")
					stco := findData[*stcoBox](stblNode, "stco")
					if stco == nil {
						stco = findData[*stcoBox](stblNode, "co64")
					}

					if stblNode != nil && stsz != nil && stsc != nil && stco != nil {
						stszEntries := stsz.entries
						stscEntries := stsc.entries
						stcoEntries := stco.entries

						if len(stszEntries) > 0 && len(stscEntries) > 0 && len(stcoEntries) > 0 {
							stszIndex := 0
							stszEntryNr := 0
							sampleNr := 0
							stscIndex := 0
							stscEntryNr := 0
							stcoIndex := 0

							stszEntry := stszEntries[stszIndex]
							stscEntry := stscEntries[stscIndex]
							sampleOffset := stcoEntries[stcoIndex]

							logStrFn := func() string {
								return fmt.Sprintf("%d: %s: nr=%d: stsz[%d/%d] nr=%d %#v stsc[%d/%d] nr=%d %#v stco[%d/%d]=%d \n",
									tc.trackID,
									trackSDDataFormat,
									sampleNr,
									stszIndex, len(stszEntries), stszEntryNr, stszEntry,
									stscIndex, len(stscEntries), stscEntryNr, stscEntry,
									stcoIndex, len(stcoEntries), sampleOffset,
								)
							}

							for stszIndex < len(stszEntries) {
								if stszEntryNr >= stszEntry.count {
									stszIndex++
									if stszIndex >= len(stszEntries) {
										break
									}

									stszEntry = stszEntries[stszIndex]
									stszEntryNr = 0
								}

								if stscEntryNr >= stscEntry.samplesPerChunk {
									stscEntryNr = 0
									stcoIndex++
									if stcoIndex >= len(stcoEntries) {
										d.Fatalf("outside stco: %s", logStrFn())
									}
									sampleOffset = stcoEntries[stcoIndex]

									if stscIndex < len(stscEntries)-1 && stcoIndex >= stscEntries[stscIndex+1].firstChunk-1 {
										stscIndex++
										if stscIndex >= len(stscEntries) {
											d.Fatalf("outside stsc: %s", logStrFn())
										}
										stscEntry = stscEntries[stscIndex]
									}
								}

								decodeSampleRange(d, objectType, ctx.opts.DecodeSamples, trackSDDataFormat, "sample", sampleOffset*8, stszEntry.size*8, formatInArg)

								sampleOffset += stszEntry.size
								stscEntryNr++
								stszEntryNr++
								sampleNr++
							}
						}
					}

					sampleNr := 0

					for _, trafNode := range tc.trafNodes {
						td := trafNodeSampleData(trafNode)
						for trunNr, tr := range td.truns {
							var sencEntries []struct{}
							if trunNr < len(td.sencs) {
								sencEntries = td.sencs[trunNr].entries
							}
							sampleOffset := td.moofOffset + tr.dataOffset

							for trunSampleNr, sz := range tr.sampleSizes {
								dataFormat := trackSDDataFormat
								if td.tfhd != nil && td.tfhd.defaultSampleDescriptionIndex != 0 && td.tfhd.defaultSampleDescriptionIndex-1 < len(sampleDescriptions) {
									sd := sampleDescriptions[td.tfhd.defaultSampleDescriptionIndex-1]
									dataFormat = sd.dataFormat
									if sd.originalFormat != "" {
										dataFormat = sd.originalFormat
									}
								}

								decodeSample := ctx.opts.DecodeSamples
								if trunSampleNr < len(sencEntries) {
									// TODO: encrypted
									decodeSample = false
								}

								decodeSampleRange(d, objectType, decodeSample, dataFormat, "sample", sampleOffset*8, sz*8, formatInArg)

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

func mp4Decode(d *decode.D) any {
	isobmffDecode(d, func(firstType string, ftyp ftypBox) {
		switch firstType {
		case "free", // seems to happen
			"moov", // seems to happen
			"pnot", // video preview file
			"jP":   // JPEG 2000
		case "ftyp",
			"styp": // segment
			switch ftyp.majorBrand {
			case "isom", // iso media (mp4ish)
				"isml",
				"qt",  // quicktime mov
				"jp2": // JPEG 2000
			default:
				d.Errorf("major_brand not isom, qt or jp2")
			}
		default:
			d.Errorf("type not ftyp, styp or moov")
		}
	})

	return nil
}
