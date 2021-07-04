package mpeg

// TODO: move aac things to mpeg?

// TODO:
// create ASC decoder
// https://github.com/mstorsjo/fdk-aac/blob/f285813ec15e7c6f8e4839c9eb4f6b0cd2da1990/libMpegTPEnc/src/tpenc_asc.cpp
// https://www.iis.fraunhofer.de/content/dam/iis/de/doc/ame/wp/FraunhoferIIS_Application-Bulletin_AAC-Transport-Formats.pdf
// https://github.com/FFmpeg/FFmpeg/blob/master/libavcodec/aac_adtstoasc_bsf.c

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

var aacFrameFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.ADTS_FRAME,
		Description: "Audio Data Transport Stream frame",
		DecodeFn:    adtsFrameDecoder,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
		},
	})
}

func adtsFrameDecoder(d *decode.D, in interface{}) interface{} {
	/*
	   adts_frame() {
	   	adts_fixed_header();
	   	adts_variable_header();
	   	if (number_of_raw_data_blocks_in_frame == 0) {
	   		adts_error_check();
	   		raw_data_block();
	   	} else {
	   		adts_header_error_check();
	   		for( i = 0; i <= number_of_raw_data_blocks_in_frame; i++ ) {
	   			raw_data_block();
	   			adts_raw_data_block_error_check();
	   		}
	   	}
	   }
	*/

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

	d.FieldValidateUFn("syncword", 0b1111_1111_1111, d.U12)
	d.FieldStringMapFn("mpeg_version", map[uint64]string{0: "MPEG-4", 1: "MPEG2- AAC"}, "Unknown", d.U1, decode.NumberDecimal)
	d.FieldValidateUFn("layer", 0, d.U2)
	protectionAbsent := d.FieldBoolFn("protection_absent", func() (bool, string) { return d.Bool(), "" })
	objectType, _ := d.FieldStringMapFn("profile", format.MPEGAudioObjectTypeNames, "Unknown", func() uint64 {
		return d.U2() + 1
	}, decode.NumberDecimal)
	d.FieldUFn("sampling_frequency_index", func() (uint64, decode.DisplayFormat, string) {
		v := d.U4()
		if v == 15 {
			return d.U24(), decode.NumberDecimal, ""
		}
		if f, ok := frequencyIndexHz[v]; ok {
			return uint64(f), decode.NumberDecimal, ""
		}
		return 0, decode.NumberDecimal, "Invalid"
	})
	d.FieldU1("private_bit")
	d.FieldStringMapFn("channel_configuration", channelConfigurationNames, "Reserved", d.U3, decode.NumberDecimal)
	d.FieldU1("originality")
	d.FieldU1("home")
	d.FieldU1("copyrighted")
	d.FieldU1("copyright")
	frameLength := d.FieldU13("frame_length")
	dataLength := frameLength - 7
	if !protectionAbsent {
		// TODO: multuple RDBs CRCs
		dataLength -= 2
	}
	d.FieldU11("buffer_fullness")
	numberOfRDBs := d.FieldUFn("number_of_rdbs", func() (uint64, decode.DisplayFormat, string) { return d.U2() + 1, decode.NumberDecimal, "" })
	if !protectionAbsent {
		d.FieldU16("crc")
	}

	d.FieldArrayFn("raw_data_blocks", func(d *decode.D) {
		for i := uint64(0); i < numberOfRDBs; i++ {
			d.FieldDecodeLen("raw_data_block", int64(dataLength)*8, aacFrameFormat, decode.FormatOptions{
				InArg: format.AACFrameIn{ObjectType: int(objectType)}},
			)
		}
	})

	return nil
}
