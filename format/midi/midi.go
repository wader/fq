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

func init() {
	interp.RegisterFormat(
		format.MIDI,
		&decode.Format{
			Description: "Standard MIDI file",
			DecodeFn:    decodeMIDI,
		})

	interp.RegisterFS(midiFS)
}

func decodeMIDI(d *decode.D) any {
	d.Endian = decode.BigEndian

	decodeMIDIFile(d)

	return nil
}

func decodeMIDIFile(d *decode.D) {
	d.FieldStruct("header", decodeMThd)

	d.FieldArray("tracks", func(d *decode.D) {
		for d.BitsLeft() > 0 {
			d.FieldStruct("track", decodeMTrk)
		}
	})
}

func decodeMThd(d *decode.D) {
	d.AssertLeastBytesLeft(8)

	if !bytes.Equal(d.PeekBytes(4), []byte("MThd")) {
		d.Errorf("no MThd marker")
	}

	d.FieldArray("header", func(d *decode.D) {
		d.FieldUTF8NullFixedLen("tag", 4)
		length := d.FieldU32("length")

		d.FramedFn(int64(length)*8, func(d *decode.D) {
			format := d.FieldU16("format")
			if format != 0 && format != 1 && format != 2 {
				d.Errorf("invalid MThd format %v (expected 0,1 or 2)", format)
			}

			tracks := d.FieldU16("tracks")
			if format == 0 && tracks > 1 {
				d.Errorf("MIDI format 0 expects 1 track (got %v)", tracks)
			}

			division := d.FieldU16("division")
			if division&0x8000 == 0x8000 {
				SMPTE := (division & 0xff00) >> 8
				if SMPTE != 0xe8 && SMPTE != SMPTE && SMPTE != 0xe6 && SMPTE != 0xe5 {
					d.Errorf("invalid MThd division SMPTE timecode type %02X (expected E8,E7, E6 or E5)", SMPTE)
				}
			}
		})
	})

	return
}

func decodeMTrk(d *decode.D) {
	d.AssertLeastBytesLeft(8)

	if !bytes.Equal(d.PeekBytes(4), []byte("MTrk")) {
		d.Errorf("no MTrk marker")
	}

	d.FieldUTF8NullFixedLen("tag", 4)
	length := d.FieldU32("length")

	d.FieldArray("events", func(d *decode.D) {
		d.FramedFn(int64(length)*8, func(d *decode.D) {
			for d.BitsLeft() > 0 {
				// d.FieldStruct("event", decodeEvent)
				decodeEvent(d)
			}
		})
	})
}

func decodeEvent(d *decode.D) {
	_, status, event := peekEvent(d)

	fmt.Printf(">> status:%02x event:%02x\n", status, event)

	// ... meta event?
	if status == 0xff {
		decodeMetaEvent(d, event)
		return
	}

	// ... sysex event

	// ... midi event?
	decodeMIDIEvent(d, status)
}

func peekEvent(d *decode.D) (uint64, uint8, uint8) {
	N := 3

	for {
		bytes := d.PeekBytes(N)

		// ... peek at delta value
		delta := uint64(0)

		for i, b := range bytes[:N-2] {
			delta <<= 7
			delta += uint64(b & 0x7f)

			if b&0x80 == 0 {
				status := bytes[i+1]
				event := bytes[i+2]
				return delta, status, event
			}
		}

		N++
	}
}

func vlq(d *decode.D) uint64 {
	vlq := uint64(0)

	for {
		b := d.BytesLen(1)

		vlq <<= 7
		vlq += uint64(b[0] & 0x7f)

		if b[0]&0x80 == 0 {
			break
		}
	}

	return vlq
}

func vlf(d *decode.D) []uint8 {
	N := int(vlq(d))

	return d.BytesLen(N)
}
