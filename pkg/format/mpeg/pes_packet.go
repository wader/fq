package mpeg

// http://dvdnav.mplayerhq.hu/dvdinfo/mpeghdrs.html
// http://stnsoft.com/DVD/sys_hdr.html))

import (
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_PES_PACKET,
		Description: "MPEG Packetized elementary stream packet",
		DecodeFn:    pesPacketDecode,
	})
}

const (
	packHeader     = 0xba
	systemHeader   = 0xbb
	privateStream1 = 0xbd
)

type subStreamPacket struct {
	number int
	bb     *bitio.Buffer
}

// for ((i=0x01; i <= 0xaf; i++)); do printf '0x%x Slice\n' $i ; done | pbcopy
var startAndStreamNames = map[uint64]string{
	0x00: "Picture",
	0x01: "Slice",
	0x02: "Slice",
	0x03: "Slice",
	0x04: "Slice",
	0x05: "Slice",
	0x06: "Slice",
	0x07: "Slice",
	0x08: "Slice",
	0x09: "Slice",
	0x0a: "Slice",
	0x0b: "Slice",
	0x0c: "Slice",
	0x0d: "Slice",
	0x0e: "Slice",
	0x0f: "Slice",
	0x10: "Slice",
	0x11: "Slice",
	0x12: "Slice",
	0x13: "Slice",
	0x14: "Slice",
	0x15: "Slice",
	0x16: "Slice",
	0x17: "Slice",
	0x18: "Slice",
	0x19: "Slice",
	0x1a: "Slice",
	0x1b: "Slice",
	0x1c: "Slice",
	0x1d: "Slice",
	0x1e: "Slice",
	0x1f: "Slice",
	0x20: "Slice",
	0x21: "Slice",
	0x22: "Slice",
	0x23: "Slice",
	0x24: "Slice",
	0x25: "Slice",
	0x26: "Slice",
	0x27: "Slice",
	0x28: "Slice",
	0x29: "Slice",
	0x2a: "Slice",
	0x2b: "Slice",
	0x2c: "Slice",
	0x2d: "Slice",
	0x2e: "Slice",
	0x2f: "Slice",
	0x30: "Slice",
	0x31: "Slice",
	0x32: "Slice",
	0x33: "Slice",
	0x34: "Slice",
	0x35: "Slice",
	0x36: "Slice",
	0x37: "Slice",
	0x38: "Slice",
	0x39: "Slice",
	0x3a: "Slice",
	0x3b: "Slice",
	0x3c: "Slice",
	0x3d: "Slice",
	0x3e: "Slice",
	0x3f: "Slice",
	0x40: "Slice",
	0x41: "Slice",
	0x42: "Slice",
	0x43: "Slice",
	0x44: "Slice",
	0x45: "Slice",
	0x46: "Slice",
	0x47: "Slice",
	0x48: "Slice",
	0x49: "Slice",
	0x4a: "Slice",
	0x4b: "Slice",
	0x4c: "Slice",
	0x4d: "Slice",
	0x4e: "Slice",
	0x4f: "Slice",
	0x50: "Slice",
	0x51: "Slice",
	0x52: "Slice",
	0x53: "Slice",
	0x54: "Slice",
	0x55: "Slice",
	0x56: "Slice",
	0x57: "Slice",
	0x58: "Slice",
	0x59: "Slice",
	0x5a: "Slice",
	0x5b: "Slice",
	0x5c: "Slice",
	0x5d: "Slice",
	0x5e: "Slice",
	0x5f: "Slice",
	0x60: "Slice",
	0x61: "Slice",
	0x62: "Slice",
	0x63: "Slice",
	0x64: "Slice",
	0x65: "Slice",
	0x66: "Slice",
	0x67: "Slice",
	0x68: "Slice",
	0x69: "Slice",
	0x6a: "Slice",
	0x6b: "Slice",
	0x6c: "Slice",
	0x6d: "Slice",
	0x6e: "Slice",
	0x6f: "Slice",
	0x70: "Slice",
	0x71: "Slice",
	0x72: "Slice",
	0x73: "Slice",
	0x74: "Slice",
	0x75: "Slice",
	0x76: "Slice",
	0x77: "Slice",
	0x78: "Slice",
	0x79: "Slice",
	0x7a: "Slice",
	0x7b: "Slice",
	0x7c: "Slice",
	0x7d: "Slice",
	0x7e: "Slice",
	0x7f: "Slice",
	0x80: "Slice",
	0x81: "Slice",
	0x82: "Slice",
	0x83: "Slice",
	0x84: "Slice",
	0x85: "Slice",
	0x86: "Slice",
	0x87: "Slice",
	0x88: "Slice",
	0x89: "Slice",
	0x8a: "Slice",
	0x8b: "Slice",
	0x8c: "Slice",
	0x8d: "Slice",
	0x8e: "Slice",
	0x8f: "Slice",
	0x90: "Slice",
	0x91: "Slice",
	0x92: "Slice",
	0x93: "Slice",
	0x94: "Slice",
	0x95: "Slice",
	0x96: "Slice",
	0x97: "Slice",
	0x98: "Slice",
	0x99: "Slice",
	0x9a: "Slice",
	0x9b: "Slice",
	0x9c: "Slice",
	0x9d: "Slice",
	0x9e: "Slice",
	0x9f: "Slice",
	0xa0: "Slice",
	0xa1: "Slice",
	0xa2: "Slice",
	0xa3: "Slice",
	0xa4: "Slice",
	0xa5: "Slice",
	0xa6: "Slice",
	0xa7: "Slice",
	0xa8: "Slice",
	0xa9: "Slice",
	0xaa: "Slice",
	0xab: "Slice",
	0xac: "Slice",
	0xad: "Slice",
	0xae: "Slice",
	0xaf: "Slice",
	0xb0: "Reserved",
	0xb1: "Reserved",
	0xb2: "User data",
	0xb3: "SequenceHeader",
	0xb4: "SequenceError",
	0xb5: "Extension",
	0xb6: "Reserved",
	0xb7: "SequenceEnd",
	0xb8: "GroupOfPictures",
	0xb9: "ProgramEnd",
	0xba: "PackHeader",
	0xbb: "SystemHeader",
	0xbc: "ProgramStreamMap",
	0xbd: "PrivateStream1",
	0xbe: "PaddingStream",
	0xbf: "PrivateStream2",
	0xc0: "MPEG1OrMPEG2AudioStream",
	0xc1: "MPEG1OrMPEG2AudioStream",
	0xc2: "MPEG1OrMPEG2AudioStream",
	0xc3: "MPEG1OrMPEG2AudioStream",
	0xc4: "MPEG1OrMPEG2AudioStream",
	0xc5: "MPEG1OrMPEG2AudioStream",
	0xc6: "MPEG1OrMPEG2AudioStream",
	0xc7: "MPEG1OrMPEG2AudioStream",
	0xc8: "MPEG1OrMPEG2AudioStream",
	0xc9: "MPEG1OrMPEG2AudioStream",
	0xca: "MPEG1OrMPEG2AudioStream",
	0xcb: "MPEG1OrMPEG2AudioStream",
	0xcc: "MPEG1OrMPEG2AudioStream",
	0xcd: "MPEG1OrMPEG2AudioStream",
	0xce: "MPEG1OrMPEG2AudioStream",
	0xcf: "MPEG1OrMPEG2AudioStream",
	0xd0: "MPEG1OrMPEG2AudioStream",
	0xd1: "MPEG1OrMPEG2AudioStream",
	0xd2: "MPEG1OrMPEG2AudioStream",
	0xd3: "MPEG1OrMPEG2AudioStream",
	0xd4: "MPEG1OrMPEG2AudioStream",
	0xd5: "MPEG1OrMPEG2AudioStream",
	0xd6: "MPEG1OrMPEG2AudioStream",
	0xd7: "MPEG1OrMPEG2AudioStream",
	0xd8: "MPEG1OrMPEG2AudioStream",
	0xd9: "MPEG1OrMPEG2AudioStream",
	0xda: "MPEG1OrMPEG2AudioStream",
	0xdb: "MPEG1OrMPEG2AudioStream",
	0xdc: "MPEG1OrMPEG2AudioStream",
	0xdd: "MPEG1OrMPEG2AudioStream",
	0xde: "MPEG1OrMPEG2AudioStream",
	0xdf: "MPEG1OrMPEG2AudioStream",
	0xe0: "MPEG1OrMPEG2VideoStream",
	0xe1: "MPEG1OrMPEG2VideoStream",
	0xe2: "MPEG1OrMPEG2VideoStream",
	0xe3: "MPEG1OrMPEG2VideoStream",
	0xe4: "MPEG1OrMPEG2VideoStream",
	0xe5: "MPEG1OrMPEG2VideoStream",
	0xe6: "MPEG1OrMPEG2VideoStream",
	0xe7: "MPEG1OrMPEG2VideoStream",
	0xe8: "MPEG1OrMPEG2VideoStream",
	0xe9: "MPEG1OrMPEG2VideoStream",
	0xea: "MPEG1OrMPEG2VideoStream",
	0xeb: "MPEG1OrMPEG2VideoStream",
	0xec: "MPEG1OrMPEG2VideoStream",
	0xed: "MPEG1OrMPEG2VideoStream",
	0xee: "MPEG1OrMPEG2VideoStream",
	0xef: "MPEG1OrMPEG2VideoStream",
	0xf0: "ECMStream",
	0xf1: "EMMStream",
	0xf2: "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Annex A or ISO/IEC 13818-6_DSMCC_stream",
	0xf3: "ISO/IEC_13522_stream",
	0xf4: "ITU-T Rec. H.222.1 type A",
	0xf5: "ITU-T Rec. H.222.1 type B",
	0xf6: "ITU-T Rec. H.222.1 type C",
	0xf7: "ITU-T Rec. H.222.1 type D",
	0xf8: "ITU-T Rec. H.222.1 type E",
	0xf9: "Ancillary_stream",
	0xfa: "Reserved",
	0xfb: "Reserved",
	0xfc: "Reserved",
	0xfd: "Reserved",
	0xfe: "Reserved",
	0xff: "Program Stream Directory",
}

var streamHasExtension = map[uint64]bool{
	0xbd: true, // Privatestream1
	0xc0: true, // MPEG1OrMPEG2AudioStream
	0xc1: true, // MPEG1OrMPEG2AudioStream
	0xc2: true, // MPEG1OrMPEG2AudioStream
	0xc3: true, // MPEG1OrMPEG2AudioStream
	0xc4: true, // MPEG1OrMPEG2AudioStream
	0xc5: true, // MPEG1OrMPEG2AudioStream
	0xc6: true, // MPEG1OrMPEG2AudioStream
	0xc7: true, // MPEG1OrMPEG2AudioStream
	0xc8: true, // MPEG1OrMPEG2AudioStream
	0xc9: true, // MPEG1OrMPEG2AudioStream
	0xca: true, // MPEG1OrMPEG2AudioStream
	0xcb: true, // MPEG1OrMPEG2AudioStream
	0xcc: true, // MPEG1OrMPEG2AudioStream
	0xcd: true, // MPEG1OrMPEG2AudioStream
	0xce: true, // MPEG1OrMPEG2AudioStream
	0xcf: true, // MPEG1OrMPEG2AudioStream
	0xd0: true, // MPEG1OrMPEG2AudioStream
	0xd1: true, // MPEG1OrMPEG2AudioStream
	0xd2: true, // MPEG1OrMPEG2AudioStream
	0xd3: true, // MPEG1OrMPEG2AudioStream
	0xd4: true, // MPEG1OrMPEG2AudioStream
	0xd5: true, // MPEG1OrMPEG2AudioStream
	0xd6: true, // MPEG1OrMPEG2AudioStream
	0xd7: true, // MPEG1OrMPEG2AudioStream
	0xd8: true, // MPEG1OrMPEG2AudioStream
	0xd9: true, // MPEG1OrMPEG2AudioStream
	0xda: true, // MPEG1OrMPEG2AudioStream
	0xdb: true, // MPEG1OrMPEG2AudioStream
	0xdc: true, // MPEG1OrMPEG2AudioStream
	0xdd: true, // MPEG1OrMPEG2AudioStream
	0xde: true, // MPEG1OrMPEG2AudioStream
	0xdf: true, // MPEG1OrMPEG2AudioStream
	0xe0: true, // MPEG1OrMPEG2VideoStream
	0xe1: true, // MPEG1OrMPEG2VideoStream
	0xe2: true, // MPEG1OrMPEG2VideoStream
	0xe3: true, // MPEG1OrMPEG2VideoStream
	0xe4: true, // MPEG1OrMPEG2VideoStream
	0xe5: true, // MPEG1OrMPEG2VideoStream
	0xe6: true, // MPEG1OrMPEG2VideoStream
	0xe7: true, // MPEG1OrMPEG2VideoStream
	0xe8: true, // MPEG1OrMPEG2VideoStream
	0xe9: true, // MPEG1OrMPEG2VideoStream
	0xea: true, // MPEG1OrMPEG2VideoStream
	0xeb: true, // MPEG1OrMPEG2VideoStream
	0xec: true, // MPEG1OrMPEG2VideoStream
	0xed: true, // MPEG1OrMPEG2VideoStream
	0xee: true, // MPEG1OrMPEG2VideoStream
	0xef: true, // MPEG1OrMPEG2VideoStream
}

func pesPacketDecode(d *decode.D, in interface{}) interface{} {
	var v interface{}

	d.FieldValidateUFn("prefix", 0b0000_0000_0000_0000_0000_0001, d.U24)
	startCode, _ := d.FieldStringMapFn("start_code", startAndStreamNames, "Unknown", d.U8)

	switch {
	case startCode == packHeader:
		d.FieldStructFn("scr", func(d *decode.D) {
			d.FieldU2("skip0")
			scr0 := d.FieldU3("scr0")
			d.FieldU1("skip1")
			scr1 := d.FieldU15("scr1")
			d.FieldU1("skip2")
			scr2 := d.FieldU15("scr2")
			d.FieldU1("skip3")
			d.FieldU9("scr_ext")
			d.FieldU1("skip4")
			scr := scr0<<30 | scr1<<15 | scr2
			d.FieldValueU("scr", scr, "")
		})
		d.FieldU22("mux_rate")
		d.FieldU2("skip0")
		d.FieldU5("reserved")
		packStuffingLength := d.FieldU3("pack_stuffing_length")
		if packStuffingLength > 0 {
			d.FieldBitBufLen("stuffing", int64(packStuffingLength*8))
		}
	case startCode == systemHeader:
		d.FieldU16("length")
		d.FieldU1("skip0")
		d.FieldU22("rate_bound")
		d.FieldU1("skip1")
		d.FieldU6("audio_bound")
		d.FieldU1("fixed_flag")
		d.FieldU1("csps_flag")
		d.FieldU1("system_audio_lock_flag")
		d.FieldU1("system_video_lock_flag")
		d.FieldU1("skip2")
		d.FieldU5("video_bound")
		d.FieldU1("packet_rate_restriction_flag")
		d.FieldU7("reserved")
		d.FieldArrayFn("stream_bound_entries", func(d *decode.D) {
			for d.PeekBits(1) == 1 {
				d.FieldStructFn("stream_bound_entry", func(d *decode.D) {
					d.FieldU8("stream_id")
					d.FieldU2("skip0")
					d.FieldU1("pstd_buffer_bound_scale")
					d.FieldU13("pstd_buffer_size_bound")
				})
			}
		})
	case startCode >= 0xbd:
		//log.Printf("startCode: %#+v\n", startCode)
		length := d.FieldU16("length")
		hasExtension := streamHasExtension[startCode]
		var headerDataLength uint64
		var extensionLength uint64
		if hasExtension {
			extensionLength = 3
			d.FieldStructFn("extension", func(d *decode.D) {
				d.FieldU2("skip0")
				d.FieldU2("scramble_control")
				d.FieldU1("priority")
				d.FieldU1("data_alignment_indicator")
				d.FieldU1("copyright")
				d.FieldU1("original")
				d.FieldU2("pts_dts_flags")
				d.FieldU1("escr_flag")
				d.FieldU1("es_rate_flag")
				d.FieldU1("dsm_trick_mode_flag")
				d.FieldU1("additional_copy_info_flag")
				d.FieldU1("pes_crc_flag")
				d.FieldU1("pes_ext_flag")
				headerDataLength = d.FieldU8("header_data_length")
			})
			// TODO:
			d.FieldBitBufLen("header_data", int64(headerDataLength)*8)
		}

		dataLen := int64(length-headerDataLength-extensionLength) * 8

		switch startCode {
		case privateStream1:
			d.FieldStructFn("data", func(d *decode.D) {
				d.DecodeLenFn(dataLen, func(d *decode.D) {
					substreamNumber := d.FieldU8("substream")
					substreamBB := d.FieldBitBufLen("data", dataLen-8)

					v = subStreamPacket{
						number: int(substreamNumber),
						bb:     substreamBB,
					}
				})
			})
		default:
			d.FieldBitBufLen("data", dataLen)
		}
	}

	return v
}
