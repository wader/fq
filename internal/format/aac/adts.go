package aac

import (
	"fq/internal/decode"
)

// Audio Data Transport Stream (ADTS)

var ADTS = &decode.Format{
	Name: "adts",
	New:  func() decode.Decoder { return &ADTSDecoder{} },
}

// ADTSDecoder is a adts  decoder
type ADTSDecoder struct {
	decode.Common
}

// Decode adts
func (d *ADTSDecoder) Decode() {

	// A	12	syncword 0xFFF, all bits must be 1
	// B	1	MPEG Version: 0 for MPEG-4, 1 for MPEG-2
	// C	2	Layer: always 0
	// D	1	protection absent, Warning, set to 1 if there is no CRC and 0 if there is CRC
	// E	2	profile, the MPEG-4 Audio Object Type minus 1
	// F	4	MPEG-4 Sampling Frequency Index (15 is forbidden)
	// G	1	private bit, guaranteed never to be used by MPEG, set to 0 when encoding, ignore when decoding
	// H	3	MPEG-4 Channel Configuration (in the case of 0, the channel configuration is sent via an inband PCE)
	// I	1	originality, set to 0 when encoding, ignore when decoding
	// J	1	home, set to 0 when encoding, ignore when decoding
	// K	1	copyrighted id bit, the next bit of a centrally registered copyright identifier, set to 0 when encoding, ignore when decoding
	// L	1	copyright id start, signals that this frame's copyright id bit is the first bit of the copyright id, set to 0 when encoding, ignore when decoding
	// M	13	frame length, this value must include 7 or 9 bytes of header length: FrameLength = (ProtectionAbsent == 1 ? 7 : 9) + size(AACFrame)
	// O	11	Buffer fullness
	// P	2	Number of AAC frames (RDBs) in ADTS frame minus 1, for maximum compatibility always use 1 AAC frame per ADTS frame
	// Q	16	CRC if protection absent is 0

	d.FieldValidateUFn("syncword", 0b111111111111, d.U12)
	d.FieldU1("mpeg_version")
	d.FieldValidateUFn("layer", 0, d.U2)
	hasCRC := !d.FieldBoolFn("protection", func() (bool, string) { return d.Bool(), "" })
	d.FieldU2("profile")
	d.FieldU4("sampling_frequency_index")
	d.FieldU1("private_bit")
	d.FieldU3("channel_configuration")
	d.FieldU1("originality")
	d.FieldU1("home")
	d.FieldU1("copyrighted")
	d.FieldU1("copyright")
	frameLength := d.FieldU13("frame_length")
	dataLength := frameLength - 7
	if hasCRC {
		dataLength -= 2
	}
	d.FieldU11("buffer_fullness")
	d.FieldUFn("number_of_frames", func() (uint64, decode.NumberFormat, string) { return d.U2() + 1, decode.NumberDecimal, "" })
	if hasCRC {
		d.FieldU16("crc")
	}
	d.FieldDecodeLen("frame", dataLength*8, Frame)
	// d.FieldBytesLen("frame", dataLength)
}
