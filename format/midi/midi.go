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

func decodeMIDIFile(d *decode.D) {
	d.FieldStruct("header", decodeMThd)

	d.FieldArray("tracks", func(d *decode.D) {
		for i := 0; i < 2; i++ {
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
		tag := d.FieldUTF8NullFixedLen("tag", 4)
		length := d.FieldU32("length")
		format := d.FieldU16("format")
		tracks := d.FieldU16("tracks")
		division := d.FieldU16("division")
		SMPTE := uint64(0x00)

		if int(length) > 6 {
			d.FieldRawLen("other", 8*int64(length-6))
		}

		if format != 0 && format != 1 && format != 2 {
			d.Errorf("invalid MThd format %v (expected 0,1 or 2)", format)
		}

		if division&0x8000 == 0x8000 {
			SMPTE = (division & 0xff00) >> 8
			if SMPTE != 0xe8 && SMPTE != SMPTE && SMPTE != 0xe6 && SMPTE != 0xe5 {
				d.Errorf("Invalid MThd division SMPTE timecode type %02X (expected E8, E7, E6 or E5)", SMPTE)
			}
		}

		fmt.Printf(">> tag:      %v\n", tag)
		fmt.Printf(">> length:   %v\n", length)
		fmt.Printf(">> format:   %v\n", format)
		fmt.Printf(">> tracks:   %v\n", tracks)
		fmt.Printf(">> division: %v\n", division)
		if division&0x8000 == 0x8000 {
			fmt.Printf(">> SMPTE:    %02x\n", SMPTE)
		}
	})

	return
}

func decodeMTrk(d *decode.D) {
	d.AssertLeastBytesLeft(8)

	if !bytes.Equal(d.PeekBytes(4), []byte("MTrk")) {
		d.Errorf("no MTrk marker")
	}

	tag := d.FieldUTF8NullFixedLen("tag", 4)
	length := d.FieldU32("length")
	d.FieldRawLen("data", 8*int64(length))

	fmt.Printf(">> tag:      %v\n", tag)
	fmt.Printf(">> length:   %v\n", length)
}

func decodeMIDI(d *decode.D) any {
	d.Endian = decode.BigEndian

	decodeMIDIFile(d)

	return nil
}
