package mpeg

// http://dvdnav.mplayerhq.hu/dvdinfo/mpeghdrs.html
// http://stnsoft.com/DVD/sys_hdr.html))

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_PES_Packet,
		&decode.Format{
			Description: "MPEG Packetized elementary stream packet",
			DecodeFn:    pesPacketDecode,
		})
}

const (
	sequenceHeader = 0xb3
	packHeader     = 0xba
	systemHeader   = 0xbb
	privateStream1 = 0xbd
)

type subStreamPacket struct {
	number int
	buf    []byte
}

var startAndStreamNames = scalar.UintRangeToScalar{
	{Range: [2]uint64{0x00, 0x00}, S: scalar.Uint{Sym: "picture"}},
	{Range: [2]uint64{0x01, 0xaf}, S: scalar.Uint{Sym: "slice"}},
	{Range: [2]uint64{0xb0, 0xb1}, S: scalar.Uint{Sym: "reserved"}},
	{Range: [2]uint64{0xb2, 0xb2}, S: scalar.Uint{Sym: "user_data"}},
	{Range: [2]uint64{0xb3, 0xb3}, S: scalar.Uint{Sym: "sequence_header"}},
	{Range: [2]uint64{0xb4, 0xb4}, S: scalar.Uint{Sym: "sequence_error"}},
	{Range: [2]uint64{0xb5, 0xb5}, S: scalar.Uint{Sym: "extension"}},
	{Range: [2]uint64{0xb6, 0xb6}, S: scalar.Uint{Sym: "reserved"}},
	{Range: [2]uint64{0xb7, 0xb7}, S: scalar.Uint{Sym: "sequence_end"}},
	{Range: [2]uint64{0xb8, 0xb8}, S: scalar.Uint{Sym: "group_of_pictures"}},
	{Range: [2]uint64{0xb9, 0xb9}, S: scalar.Uint{Sym: "program_end"}},
	{Range: [2]uint64{0xba, 0xba}, S: scalar.Uint{Sym: "pack_header"}},
	{Range: [2]uint64{0xbb, 0xbb}, S: scalar.Uint{Sym: "system_header"}},
	{Range: [2]uint64{0xbc, 0xbc}, S: scalar.Uint{Sym: "program_stream_map"}},
	{Range: [2]uint64{0xbd, 0xbd}, S: scalar.Uint{Sym: "private_stream1"}},
	{Range: [2]uint64{0xbe, 0xbe}, S: scalar.Uint{Sym: "padding_stream"}},
	{Range: [2]uint64{0xbf, 0xbf}, S: scalar.Uint{Sym: "private_stream2"}},
	{Range: [2]uint64{0xc0, 0xdf}, S: scalar.Uint{Sym: "audio_stream"}},
	{Range: [2]uint64{0xe0, 0xef}, S: scalar.Uint{Sym: "video_stream"}},
	{Range: [2]uint64{0xf0, 0xf0}, S: scalar.Uint{Sym: "ecm_stream"}},
	{Range: [2]uint64{0xf1, 0xf1}, S: scalar.Uint{Sym: "emm_stream"}},
	{Range: [2]uint64{0xf2, 0xf2}, S: scalar.Uint{Sym: "itu_t_rec_h_222_0"}},
	{Range: [2]uint64{0xf3, 0xf3}, S: scalar.Uint{Sym: "iso_iec_13522_stream"}},
	{Range: [2]uint64{0xf4, 0xf4}, S: scalar.Uint{Sym: "itu_t_rec_h_222_1_type_a"}},
	{Range: [2]uint64{0xf5, 0xf5}, S: scalar.Uint{Sym: "itu_t_rec_h_222_1_type_b"}},
	{Range: [2]uint64{0xf6, 0xf6}, S: scalar.Uint{Sym: "itu_t_rec_h_222_1_type_c"}},
	{Range: [2]uint64{0xf7, 0xf7}, S: scalar.Uint{Sym: "itu_t_rec_h_222_1_type_d"}},
	{Range: [2]uint64{0xf8, 0xf8}, S: scalar.Uint{Sym: "itu_t_rec_h_222_1_type_e"}},
	{Range: [2]uint64{0xf9, 0xf9}, S: scalar.Uint{Sym: "ancillary_stream"}},
	{Range: [2]uint64{0xfa, 0xfe}, S: scalar.Uint{Sym: "reserved"}},
	{Range: [2]uint64{0xff, 0xff}, S: scalar.Uint{Sym: "program_stream_directory"}},
}

var mpegVersion = scalar.UintMapDescription{
	0b01: "MPEG2",
	0b10: "MPEG1",
}

func pesPacketDecode(d *decode.D) any {
	var v any

	d.FieldU24("prefix", d.UintAssert(0b0000_0000_0000_0000_0000_0001), scalar.UintBin)
	startCode := d.FieldU8("start_code", startAndStreamNames, scalar.UintHex)

	switch {
	case startCode == sequenceHeader:
		d.FieldU12("horizontal_size")
		d.FieldU12("vertical_size")
		d.FieldU4("aspect_ratio")
		d.FieldU4("frame_rate_code")
		// TODO: bit rate * 400, rounded upwards. Use 0x3FFFF for variable bit rate
		d.FieldU18("bit_rate")
		d.FieldU1("marker_bit")
		d.FieldU10("vbv_buf_size")
		d.FieldU1("constrained_parameters_flag")
		loadIntraQuantizerMatrix := d.FieldBool("load_intra_quantizer_matrix")
		if loadIntraQuantizerMatrix {
			d.FieldRawLen("intra_quantizer_matrix", 8*64)

		}
		loadNonIntraQuantizerMatrix := d.FieldBool("load_non_intra_quantizer_matrix")
		if loadNonIntraQuantizerMatrix {
			d.FieldRawLen("non_intra_quantizer_matrix", 8*64)

		}
	case startCode == packHeader:
		isMPEG2 := d.PeekUintBits(2) == 0b01
		if isMPEG2 {
			d.FieldU2("marker_bits0", mpegVersion)
		} else {
			d.FieldU4("marker_bits0", mpegVersion)
		}
		scr0 := d.FieldU3("system_clock0")
		d.FieldU1("marker_bits1")
		scr1 := d.FieldU15("system_clock1")
		d.FieldU1("marker_bits2")
		scr2 := d.FieldU15("system_clock2")
		d.FieldU1("marker_bits3")
		if isMPEG2 {
			d.FieldU9("scr_ext")
		}
		d.FieldU1("marker_bits4")
		scr := scr0<<30 | scr1<<15 | scr2
		d.FieldValueUint("scr", scr)
		d.FieldU22("mux_rate")
		d.FieldU1("marker_bits5")
		if isMPEG2 {
			d.FieldU1("marker_bits6")
			d.FieldU5("reserved")
			packStuffingLength := d.FieldU3("pack_stuffing_length")
			if packStuffingLength > 0 {
				d.FieldRawLen("stuffing", int64(packStuffingLength*8))
			}
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
		d.FieldArray("stream_bound_entries", func(d *decode.D) {
			for d.PeekUintBits(1) == 1 {
				d.FieldStruct("stream_bound_entry", func(d *decode.D) {
					d.FieldU8("stream_id")
					d.FieldU2("skip0")
					d.FieldU1("pstd_buffer_bound_scale")
					d.FieldU13("pstd_buffer_size_bound")
				})
			}
		})
	case startCode >= 0xbd:
		length := d.FieldU16("length")
		// 0xbd-0xbd // Privatestream1
		// 0xc0-0xdf // MPEG1OrMPEG2AudioStream
		// 0xe0-0xef // MPEG1OrMPEG2VideoStream
		hasExtension := startCode == 0xbd || (startCode >= 0xc0 && startCode <= 0xef)
		var headerDataLength uint64
		var extensionLength uint64
		if hasExtension {
			extensionLength = 3
			d.FieldStruct("extension", func(d *decode.D) {
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
			d.FieldRawLen("header_data", int64(headerDataLength)*8)
		}

		dataLen := int64(length-headerDataLength-extensionLength) * 8

		switch startCode {
		case privateStream1:
			d.FieldStruct("data", func(d *decode.D) {
				d.FramedFn(dataLen, func(d *decode.D) {
					substreamNumber := d.FieldU8("substream")
					substreamBR := d.FieldRawLen("data", dataLen-8)

					v = subStreamPacket{
						number: int(substreamNumber),
						buf:    d.ReadAllBits(substreamBR),
					}
				})
			})
		default:
			d.FieldRawLen("stream_data", dataLen)
		}
	default:
		// nop
	}

	if d.BitsLeft() > 0 {
		d.FieldRawLen("data", d.BitsLeft())
	}

	return v
}
