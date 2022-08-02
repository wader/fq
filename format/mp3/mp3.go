package mp3

// TODO: vbri

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var headerFormat decode.Group
var footerFormat decode.Group
var mp3Frame decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.MP3,
		ProbeOrder:  format.ProbeOrderBinFuzzy, // after most others (silent samples and jpeg header can look like mp3 sync)
		Description: "MP3 file",
		Groups:      []string{format.PROBE},
		DecodeFn:    mp3Decode,
		DecodeInArg: format.Mp3In{
			MaxUniqueHeaderConfigs: 5,
			MaxSyncSeek:            4 * 1024 * 8,
		},
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Group: &headerFormat},
			{
				Names: []string{
					format.ID3V1,
					format.ID3V11,
					format.APEV2,
				},
				Group: &footerFormat,
			},
			{Names: []string{format.MP3_FRAME}, Group: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D, in any) any {
	mi, _ := in.(format.Mp3In)

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

	// there are mp3s files in the wild with multiple headers, two id3v2 tags etc
	d.FieldArray("headers", func(d *decode.D) {
		for d.NotEnd() {
			if dv, _, _ := d.TryFieldFormat("header", headerFormat, nil); dv == nil {
				return
			}
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

			dv, v, _ := d.TryFieldFormat("frame", mp3Frame, nil)
			if dv == nil {
				decodeFailures++
				d.SeekRel(8)
				continue
			}
			mfo, ok := v.(format.MP3FrameOut)
			if !ok {
				panic(fmt.Sprintf("expected MP3FrameOut got %#+v", v))
			}
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
			if dv, _, _ := d.TryFieldFormat("footer", footerFormat, nil); dv == nil {
				return
			}
		}
	})

	return nil
}
