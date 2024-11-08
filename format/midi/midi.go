/*
Package midi implements an fq plugin to decode [standard MIDI files].

The MIDI decoder is a member of the 'probe' group and fq should automatically invoke the
decoder when opening a MIDI file. The decoder can be explicitly specified with the '-d midi'
command line option.

The decoder currently only supports MIDI 1.0 files and does only basic validation on the
MIDI file structure.

[standard MIDI files]: https://midi.org/standard-midi-files.
*/
package midi

// https://www.midi.org/specifications/item/the-midi-1-0-specification

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed midi.md
var midiFS embed.FS

// context is a container struct for the running parse information required to
// decode a MIDI track.
type context struct {
	tick    uint64
	running uint8
	casio   bool
}

// init registers the MIDI format decoder and adds it to the 'probe' group.
func init() {
	interp.RegisterFormat(
		format.MIDI,
		&decode.Format{
			Description: "Standard MIDI file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeMIDI,
		})

	interp.RegisterFS(midiFS)
}

// decodeMIDI implements the MIDI file decoder.
//
// The decoder parses the file as a set of chunks, each comprising a 4 character tag
// followed by a uint32 length field. The decoder parses the MTHd and MTrk MIDI chunks.
func decodeMIDI(d *decode.D) any {
	d.Endian = decode.BigEndian

	// ... decode header
	format, _ := decodeMThd(d)

	// ... decode tracks (and other chunks)
	d.FieldArray("content", func(d *decode.D) {
		tracks := uint16(0)

		for d.BitsLeft() > 0 {
			if bytes.Equal(d.PeekBytes(4), []byte("MTrk")) {
				switch {
				case format == 0 && tracks > 0: // decode 'extra' format 0 tracks as data
					d.FieldStruct("other", decodeOther)

				default:
					d.FieldStruct("track", decodeMTrk)
				}

				tracks++
			} else {
				d.FieldStruct("other", decodeOther)
			}
		}
	})

	return nil
}

// decodeMThd decodes an MThd MIDI header chunk into a 'header' struct with the fields:
//   - tag       "MThd"
//   - length    Header chunk size
//   - format    MIDI format (0,1 or 2)
//   - tracks    Number of tracks
//   - division  Time division
func decodeMThd(d *decode.D) (uint16, uint16) {
	var format uint16
	var tracks uint16

	f := func(d *decode.D) {
		if !bytes.Equal(d.PeekBytes(4), []byte("MThd")) {
			d.Errorf("missing MThd tag")
		}

		d.FieldUTF8("tag", 4)
		length := d.FieldS32("length")

		d.FramedFn(length*8, func(d *decode.D) {
			format = uint16(d.FieldU16("format"))
			tracks = uint16(d.FieldU16("tracks"))

			d.FieldStruct("division", func(d *decode.D) {
				if division := d.PeekUintBits(16); division&0x8000 == 0x8000 {
					d.FieldS8("fps", fpsMap)
					d.FieldU8("resolution")
				} else {
					d.FieldU16("ppqn")
				}
			})
		})
	}

	d.FieldStruct("header", f)

	return format, tracks
}

// decodeMTrk decodes an MTrk MIDI track chunk into a struct with the header fields:
//   - tag      "MTrk"
//   - length   Track chunk size
//   - events   List of track events
func decodeMTrk(d *decode.D) {
	if !bytes.Equal(d.PeekBytes(4), []byte("MTrk")) {
		d.Errorf("missing MTrk tag")
	}

	d.FieldUTF8("tag", 4)
	length := d.FieldS32("length")

	d.FieldArray("events", func(d *decode.D) {
		d.FramedFn(length*8, func(d *decode.D) {
			ctx := context{
				tick:    0,
				running: 0x000,
				casio:   false,
			}

			for d.BitsLeft() > 0 {
				decodeEvent(d, &ctx)
			}
		})
	})
}

// decodeEvent decodes a single MIDI event as either:
//   - Meta event
//   - MIDI channel event
//   - SysEx system event
func decodeEvent(d *decode.D, ctx *context) {
	_, status, event := peekEvent(d)

	if status == 0xf0 || status == 0xf7 {
		decodeSysExEvent(d, status, ctx)
	} else if status == 0xff {
		decodeMetaEvent(d, event, ctx)
	} else if status < 0xf0 {
		decodeMIDIEvent(d, status, ctx)
	} else {
		d.Errorf("invalid status byte (%02x)", status)
	}
}

// peekEvent retrieves the type of the next event without moving the reader location.
func peekEvent(d *decode.D) (uint64, uint8, uint8) {
	var N int = 1

	for {
		bytes := d.PeekBytes(N)
		delta := uint64(0)
		ix := 0

		for ix < N {
			b := bytes[ix]
			ix++

			delta <<= 7
			delta += uint64(b & 0x7f)

			if b&0x80 == 0 {
				if ix < N {
					status := bytes[ix]
					ix++

					// ... MIDI event?
					if status < 0xf0 {
						return delta, status, 0x00
					}

					// ... sysex?
					if status == 0xf0 || status == 0xf7 {
						return delta, status, 0x00
					}

					// ... (invalid) real-time event
					if status != 0xff {
						return delta, status, 0x00
					}

					// ... meta-event
					if ix < N {
						event := bytes[ix]
						return delta, status, event
					}
				}
			}
		}

		N++
	}
}

// decodeOther decodes non-MIDI chunks as raw data.
func decodeOther(d *decode.D) {
	d.FieldUTF8("tag", 4)
	length := d.FieldS32("length")
	d.FieldRawLen("data", length*8)
}

// vlq decodes a MIDI big-endian varuint.
func vlq(d *decode.D) uint64 {
	vlq := uint64(0)

	for {
		b := d.U8()

		vlq <<= 7
		vlq += b & 0x7f

		if b&0x80 == 0 {
			break
		}
	}

	return vlq
}

// vlf decodes a MIDI byte array prefixed with a varuint length.
func vlf(d *decode.D) ([]uint8, error) {
	N := vlq(d)

	if N*8 > uint64(d.BitsLeft()) {
		d.Fatalf("invalid field length")
	}

	if int64(N*8) > d.BitsLeft() {
		return nil, fmt.Errorf("invalid field length")
	} else {
		return d.BytesLen(int(N)), nil
	}
}

// vlstring decodes a MIDI string prefixed with a varuint length.
func vlstring(d *decode.D) string {
	if data, err := vlf(d); err != nil {
		d.Fatalf("%v", err)
	} else {
		return string(data)
	}

	return ""
}

// flush reads and discards any remaining bits in a chunk after encountering an
// invalid event.
func flush(d *decode.D, format string, args ...any) {
	d.Errorf(format, args...)

	var N int = int(d.BitsLeft())

	d.Bits(N)
}
