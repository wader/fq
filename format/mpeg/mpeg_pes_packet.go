package mpeg

// http://dvdnav.mplayerhq.hu/dvdinfo/mpeghdrs.html
// http://stnsoft.com/DVD/sys_hdr.html))

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MPEG_PES_PACKET,
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
	bb     *bitio.Buffer
}

var startAndStreamNames = scalar.URangeToScalar{
	{0x00, 0x00}: {Sym: "Picture"},
	{0x01, 0xaf}: {Sym: "Slice"},
	{0xb0, 0xb1}: {Sym: "Reserved"},
	{0xb2, 0xb2}: {Sym: "User data"},
	{0xb3, 0xb3}: {Sym: "SequenceHeader"},
	{0xb4, 0xb4}: {Sym: "SequenceError"},
	{0xb5, 0xb5}: {Sym: "Extension"},
	{0xb6, 0xb6}: {Sym: "Reserved"},
	{0xb7, 0xb7}: {Sym: "SequenceEnd"},
	{0xb8, 0xb8}: {Sym: "GroupOfPictures"},
	{0xb9, 0xb9}: {Sym: "ProgramEnd"},
	{0xba, 0xba}: {Sym: "PackHeader"},
	{0xbb, 0xbb}: {Sym: "SystemHeader"},
	{0xbc, 0xbc}: {Sym: "ProgramStreamMap"},
	{0xbd, 0xbd}: {Sym: "PrivateStream1"},
	{0xbe, 0xbe}: {Sym: "PaddingStream"},
	{0xbf, 0xbf}: {Sym: "PrivateStream2"},
	{0xc0, 0xdf}: {Sym: "MPEG1OrMPEG2AudioStream"},
	{0xe0, 0xef}: {Sym: "MPEG1OrMPEG2VideoStream"},
	{0xf0, 0xf0}: {Sym: "ECMStream"},
	{0xf1, 0xf1}: {Sym: "EMMStream"},
	{0xf2, 0xf2}: {Sym: "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Annex A or ISO/IEC 13818-6_DSMCC_stream"},
	{0xf3, 0xf3}: {Sym: "ISO/IEC_13522_stream"},
	{0xf4, 0xf4}: {Sym: "ITU-T Rec. H.222.1 type A"},
	{0xf5, 0xf5}: {Sym: "ITU-T Rec. H.222.1 type B"},
	{0xf6, 0xf6}: {Sym: "ITU-T Rec. H.222.1 type C"},
	{0xf7, 0xf7}: {Sym: "ITU-T Rec. H.222.1 type D"},
	{0xf8, 0xf8}: {Sym: "ITU-T Rec. H.222.1 type E"},
	{0xf9, 0xf9}: {Sym: "Ancillary_stream"},
	{0xfa, 0xfe}: {Sym: "Reserved"},
	{0xff, 0xff}: {Sym: "Program Stream Directory"},
}

func pesPacketDecode(d *decode.D, in interface{}) interface{} {
	var v interface{}

	d.FieldU24("prefix", d.AssertU(0b0000_0000_0000_0000_0000_0001), scalar.Bin)
	startCode := d.FieldU8("start_code", startAndStreamNames, scalar.Hex)

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
		d.FieldStruct("scr", func(d *decode.D) {
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
			d.FieldValueU("scr", scr)
		})
		d.FieldU22("mux_rate")
		d.FieldU2("skip0")
		d.FieldU5("reserved")
		packStuffingLength := d.FieldU3("pack_stuffing_length")
		if packStuffingLength > 0 {
			d.FieldRawLen("stuffing", int64(packStuffingLength*8))
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
			for d.PeekBits(1) == 1 {
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
				d.LenFn(dataLen, func(d *decode.D) {
					substreamNumber := d.FieldU8("substream")
					substreamBB := d.FieldRawLen("data", dataLen-8)

					v = subStreamPacket{
						number: int(substreamNumber),
						bb:     substreamBB,
					}
				})
			})
		default:
			d.FieldRawLen("data", dataLen)
		}
	default:
		d.FieldRawLen("data", d.BitsLeft())
	}

	return v
}
