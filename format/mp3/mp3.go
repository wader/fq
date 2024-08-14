package mp3

// TODO: vbri

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var headerGroup decode.Group
var footerGroup decode.Group
var mp3FrameGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.MP3,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinFuzzy, // after most others (silent samples and jpeg header can look like mp3 sync)
			Description: "MP3 file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    mp3Decode,
			DefaultInArg: format.MP3_In{
				MaxUniqueHeaderConfigs: 5,
				MaxUnknown:             50,
				MaxSyncSeek:            4 * 1024 * 8,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.ID3v2}, Out: &headerGroup},
				{
					Groups: []*decode.Group{
						format.ID3v1,
						format.ID3v11,
						format.Apev2,
					},
					Out: &footerGroup,
				},
				{Groups: []*decode.Group{format.MP3_Frame}, Out: &mp3FrameGroup},
			},
		})
}

func mp3Decode(d *decode.D) any {
	var mi format.MP3_In
	d.ArgAs(&mi)

	// things in a mp3 stream usually have few unique combinations of.
	// does not include bitrate on purpose
	type headerConfig struct {
		MPEGVersion      int
		ProtectionAbsent bool
		SampleRate       int
		ChannelsIndex    int
		ChannelModeIndex int
	}
	uniqueHeaderConfigs := map[headerConfig]struct{}{}
	knownSize := int64(0)

	// there are mp3s files in the wild with multiple headers, two id3v2 tags etc
	d.FieldArray("headers", func(d *decode.D) {
		for d.NotEnd() {
			headerStart := d.Pos()
			if dv, _, _ := d.TryFieldFormat("header", &headerGroup, nil); dv == nil {
				return
			}
			knownSize += d.Pos() - headerStart
		}
	})

	lastValidEnd := int64(0)
	validFrames := 0
	decodeFailures := 0
	d.FieldArray("frames", func(d *decode.D) {
		for d.NotEnd() {
			syncLen, _, err := d.TryPeekFind(16, 8, int64(mi.MaxSyncSeek), func(v uint64) bool {
				return (v&0b1111_1111_1110_0000 == 0b1111_1111_1110_0000 && // sync header
					v&0b0000_0000_0001_1000 != 0b0000_0000_0000_1000 && // not reserved mpeg version
					v&0b0000_0000_0000_0110 == 0b0000_0000_0000_0010) // layer 3
			})
			if err != nil || syncLen < 0 {
				break
			}
			if syncLen > 0 {
				d.SeekRel(syncLen)
			}

			frameStart := d.Pos()
			dv, v, _ := d.TryFieldFormat("frame", &mp3FrameGroup, nil)
			if dv == nil {
				decodeFailures++
				d.SeekRel(8)
				continue
			}
			mfo, ok := v.(format.MP3_Frame_Out)
			if !ok {
				panic(fmt.Sprintf("expected MP3FrameOut got %#+v", v))
			}
			knownSize += d.Pos() - frameStart
			uniqueHeaderConfigs[headerConfig{
				MPEGVersion:      mfo.MPEGVersion,
				ProtectionAbsent: mfo.ProtectionAbsent,
				SampleRate:       mfo.SampleRate,
				ChannelsIndex:    mfo.ChannelsIndex,
				ChannelModeIndex: mfo.ChannelModeIndex,
			}] = struct{}{}

			lastValidEnd = d.Pos()
			validFrames++

			if len(uniqueHeaderConfigs) >= mi.MaxUniqueHeaderConfigs {
				d.Errorf("too many unique header configurations")
			}
		}
	})

	if validFrames == 0 || (validFrames < 2 && decodeFailures > 0) {
		d.Errorf("no frames found")
	}

	d.SeekAbs(lastValidEnd)

	d.FieldArray("footers", func(d *decode.D) {
		for d.NotEnd() {
			footerStart := d.Pos()
			if dv, _, _ := d.TryFieldFormat("footer", &footerGroup, nil); dv == nil {
				return
			}
			knownSize += d.Pos() - footerStart
		}
	})

	unknownPercent := int(float64((d.Len() - knownSize)) / float64(d.Len()) * 100.0)
	if unknownPercent > mi.MaxUnknown {
		d.Errorf("exceeds max precent unknown bits, %d > %d", unknownPercent, mi.MaxUnknown)
	}

	return nil
}
