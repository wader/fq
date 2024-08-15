package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type MetaEventType uint8

const (
	TypeSequenceNumber         MetaEventType = 0x00
	TypeText                   MetaEventType = 0x01
	TypeCopyright              MetaEventType = 0x02
	TypeTrackName              MetaEventType = 0x03
	TypeInstrumentName         MetaEventType = 0x04
	TypeLyric                  MetaEventType = 0x05
	TypeMarker                 MetaEventType = 0x06
	TypeCuePoint               MetaEventType = 0x07
	TypeProgramName            MetaEventType = 0x08
	TypeDeviceName             MetaEventType = 0x09
	TypeMIDIChannelPrefix      MetaEventType = 0x20
	TypeMIDIPort               MetaEventType = 0x21
	TypeTempo                  MetaEventType = 0x51
	TypeSMPTEOffset            MetaEventType = 0x54
	TypeTimeSignature          MetaEventType = 0x58
	TypeKeySignature           MetaEventType = 0x59
	TypeEndOfTrack             MetaEventType = 0x2f
	TypeSequencerSpecificEvent MetaEventType = 0x7f
)

const (
	keyCMajor      = 0x0001
	keyGMajor      = 0x0101
	keyDMajor      = 0x0201
	keyAMajor      = 0x0301
	keyEMajor      = 0x0401
	keyBMajor      = 0x0501
	keyFSharpMajor = 0x0601
	keyCSharpMajor = 0x0701
	keyFMajor      = 0x0801
	keyBFlatMajor  = 0x0901
	keyEFlatMajor  = 0x0a01
	keyAFlatMajor  = 0x0b01
	keyDFlatMajor  = 0x0c01
	keyGFlatMajor  = 0x0d01
	keyCFlatMajor  = 0x0e01

	keyAMinor      = 0x0000
	keyEMinor      = 0x0100
	keyBMinor      = 0x0200
	keyASharpMinor = 0x0300
	keyDSharpMinor = 0x0400
	keyGSharpMinor = 0x0500
	keyCSharpMinor = 0x0600
	keyFSharpMinor = 0x0700
	keyDMinor      = 0x0800
	keyGMinor      = 0x0900
	keyCMinor      = 0x0a00
	keyFMinor      = 0x0b00
	keyBFlatMinor  = 0x0c00
	keyEFlatMinor  = 0x0d00
	keyAFlatMinor  = 0x0e00
)

var keys = scalar.UintMapSymStr{
	keyCMajor:      "C major",
	keyGMajor:      "G major",
	keyDMajor:      "D major",
	keyAMajor:      "A major",
	keyEMajor:      "E major",
	keyBMajor:      "B major",
	keyFSharpMajor: "F♯ major",
	keyCSharpMajor: "C♯ major",
	keyFMajor:      "F major",
	keyBFlatMajor:  "B♭ major",
	keyEFlatMajor:  "E♭ major",
	keyAFlatMajor:  "A♭ major",
	keyDFlatMajor:  "D♭ major",
	keyGFlatMajor:  "G♭ major",
	keyCFlatMajor:  "C♭ major",

	keyAMinor:      "A minor",
	keyEMinor:      "E minor",
	keyBMinor:      "B minor",
	keyASharpMinor: "A♯ minor",
	keyDSharpMinor: "D♯ minor",
	keyGSharpMinor: "G♯ minor",
	keyCSharpMinor: "C♯ minor",
	keyFSharpMinor: "F♯ minor",
	keyDMinor:      "D minor",
	keyGMinor:      "G minor",
	keyCMinor:      "C minor",
	keyFMinor:      "F minor",
	keyBFlatMinor:  "B♭ minor",
	keyEFlatMinor:  "E♭ minor",
	keyAFlatMinor:  "A♭ minor",
}

func decodeMetaEvent(d *decode.D, event uint8) {
	switch MetaEventType(event) {
	case TypeTrackName:
		d.FieldStruct("TrackName", decodeTrackName)
		return

	case TypeTempo:
		d.FieldStruct("Tempo", decodeTempo)
		return

	case TypeTimeSignature:
		d.FieldStruct("TimeSignature", decodeTimeSignature)
		return

	case TypeKeySignature:
		d.FieldStruct("KeySignature", decodeKeySignature)
		return

	case TypeEndOfTrack:
		d.FieldStruct("EndOfTrack", decodeEndOfTrack)
		return

		// TypeSequenceNumber         MetaEventType = 0x00
		// TypeText                   MetaEventType = 0x01
		// TypeCopyright              MetaEventType = 0x02
		// TypeInstrumentName         MetaEventType = 0x04
		// TypeLyric                  MetaEventType = 0x05
		// TypeMarker                 MetaEventType = 0x06
		// TypeCuePoint               MetaEventType = 0x07
		// TypeProgramName            MetaEventType = 0x08
		// TypeDeviceName             MetaEventType = 0x09
		// TypeMIDIChannelPrefix      MetaEventType = 0x20
		// TypeMIDIPort               MetaEventType = 0x21
		// TypeSMPTEOffset            MetaEventType = 0x54
		// TypeSequencerSpecificEvent MetaEventType = 0x7f
	}

	// ... unknown event - flush remaining data
	fmt.Printf("UNKNOWN META EVENT:%02x\n", event)

	var N int = int(d.BitsLeft())

	d.Bits(N)
}

func decodeTrackName(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("name", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeTempo(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	d.FieldUintFn("tempo", func(d *decode.D) uint64 {
		tempo := uint64(0)
		bytes := vlf(d)

		for _, b := range bytes {
			tempo <<= 8
			tempo |= uint64(b & 0x00ff)
		}

		return tempo
	})
}

func decodeTimeSignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStruct("signature", func(d *decode.D) {
		bytes := vlf(d)

		if len(bytes) > 0 {
			d.FieldValueUint("numerator", uint64(bytes[0]))
		}

		if len(bytes) > 1 {
			denominator := uint16(1)
			for i := uint8(0); i < bytes[1]; i++ {
				denominator *= 2
			}

			d.FieldValueUint("denominator", uint64(denominator))
		}

		if len(bytes) > 2 {
			d.FieldValueUint("ticksPerClick", uint64(bytes[2]))
		}

		if len(bytes) > 3 {
			d.FieldValueUint("thirtySecondsPerQuarter", uint64(bytes[3]))
		}
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	bytes := vlf(d)
	if len(bytes) > 1 {
		key := (uint64(bytes[0]) << 8) & 0xff00
		key |= (uint64(bytes[1]) << 0) & 0x00ff

		d.FieldValueUint("key", key, keys)
	}
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	vlf(d)
}
