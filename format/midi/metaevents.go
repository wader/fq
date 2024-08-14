package midi

import (
	"github.com/wader/fq/pkg/decode"
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

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	vlf(d)
}
