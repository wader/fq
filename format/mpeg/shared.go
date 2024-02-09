package mpeg

import (
	"bytes"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFunc0("nal_unescape", makeBinaryTransformFn(func(r io.Reader) (io.Reader, error) {
		return &nalUnescapeReader{Reader: r}, nil
	}))
}

// transform to binary using fn
func makeBinaryTransformFn(fn func(r io.Reader) (io.Reader, error)) func(_ *interp.Interp, c any) any {
	return func(_ *interp.Interp, c any) any {
		inBR, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}

		r, err := fn(bitio.NewIOReader(inBR))
		if err != nil {
			return err
		}

		outBuf := &bytes.Buffer{}
		if _, err := io.Copy(outBuf, r); err != nil {
			return err
		}

		outBR := bitio.NewBitReader(outBuf.Bytes(), -1)

		bb, err := interp.NewBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}
		return bb
	}
}

func decodeEscapeValueFn(add int, b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return func(d *decode.D) uint64 {
		n1 := d.U(b1)
		n := n1
		if n1 == (1<<b1)-1 {
			n2 := d.U(b2)
			if add != -1 {
				n += n2 + uint64(add)
			} else {
				n = n2
			}
			if n2 == (1<<b2)-1 {
				n3 := d.U(b3)
				if add != -1 {
					n += n3 + uint64(add)
				} else {
					n = n3
				}
			}
		}
		return n
	}
}

// use last non-escaped value
func decodeEscapeValueAbsFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(-1, b1, b2, b3)
}

// add values and escaped values
//
//nolint:deadcode,unused
func decodeEscapeValueAddFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(0, b1, b2, b3)
}

// add values and escaped values+1
func decodeEscapeValueCarryFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(1, b1, b2, b3)
}

// TODO: move?
// TODO: make generic replace reader? share with id3v2 unsync?
type nalUnescapeReader struct {
	io.Reader
	lastTwoZeros [2]bool
}

func (r nalUnescapeReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	ni := 0
	for i, b := range p[0:n] {
		if r.lastTwoZeros[0] && r.lastTwoZeros[1] && b == 0x03 {
			n--
			r.lastTwoZeros[0] = false
			r.lastTwoZeros[1] = false
			continue
		} else {
			r.lastTwoZeros[1] = r.lastTwoZeros[0]
			r.lastTwoZeros[0] = b == 0
		}
		p[ni] = p[i]
		ni++
	}

	return n, err
}

const tsPacketLength = 188 * 8

const (
	pidPAT = 0
)

func tsPidIsTable(pid int, pmt map[int]format.MpegTsProgram) bool {
	// pid 0x0-0x1f seems to all be tables
	if pid >= 0 && pid <= 0x1f {
		return true
	}
	_, isPMT := pmt[pid]
	return isPMT
}

var tsStreamTagMap = scalar.UintRangeToScalar{
	{Range: [2]uint64{0, 0}, S: scalar.Uint{Description: "Reserved"}},
	{Range: [2]uint64{1, 1}, S: scalar.Uint{Description: "Reserved"}},
	{Range: [2]uint64{2, 2}, S: scalar.Uint{Description: "video_stream_descriptor"}},
	{Range: [2]uint64{3, 3}, S: scalar.Uint{Description: "audio_stream_descriptor"}},
	{Range: [2]uint64{4, 4}, S: scalar.Uint{Description: "hierarchy_descriptor"}},
	{Range: [2]uint64{5, 5}, S: scalar.Uint{Description: "registration_descriptor"}},
	{Range: [2]uint64{6, 6}, S: scalar.Uint{Description: "data_stream_alignment_descriptor"}},
	{Range: [2]uint64{7, 7}, S: scalar.Uint{Description: "target_background_grid_descriptor"}},
	{Range: [2]uint64{8, 8}, S: scalar.Uint{Description: "video_window_descriptor"}},
	{Range: [2]uint64{9, 9}, S: scalar.Uint{Description: "CA_descriptor"}},
	{Range: [2]uint64{10, 10}, S: scalar.Uint{Description: "ISO_639_language_descriptor"}},
	{Range: [2]uint64{11, 11}, S: scalar.Uint{Description: "system_clock_descriptor"}},
	{Range: [2]uint64{12, 12}, S: scalar.Uint{Description: "multiplex_buffer_utilization_descriptor"}},
	{Range: [2]uint64{13, 13}, S: scalar.Uint{Description: "copyright_descriptor"}},
	{Range: [2]uint64{14, 14}, S: scalar.Uint{Description: "maximum_bitrate_descriptor"}},
	{Range: [2]uint64{15, 15}, S: scalar.Uint{Description: "private_data_indicator_descriptor"}},
	{Range: [2]uint64{16, 16}, S: scalar.Uint{Description: "smoothing_buffer_descriptor"}},
	{Range: [2]uint64{17, 17}, S: scalar.Uint{Description: "STD_descriptor"}},
	{Range: [2]uint64{18, 18}, S: scalar.Uint{Description: "IBP_descriptor"}},
	{Range: [2]uint64{19, 26}, S: scalar.Uint{Description: "Defined in ISO/IEC 13818-6"}},
	{Range: [2]uint64{27, 27}, S: scalar.Uint{Description: "MPEG-4_video_descriptor"}},
	{Range: [2]uint64{28, 28}, S: scalar.Uint{Description: "MPEG-4_audio_descriptor"}},
	{Range: [2]uint64{29, 29}, S: scalar.Uint{Description: "IOD_descriptor"}},
	{Range: [2]uint64{30, 30}, S: scalar.Uint{Description: "SL_descriptor"}},
	{Range: [2]uint64{31, 31}, S: scalar.Uint{Description: "FMC_descriptor"}},
	{Range: [2]uint64{32, 32}, S: scalar.Uint{Description: "external_ES_ID_descriptor"}},
	{Range: [2]uint64{33, 33}, S: scalar.Uint{Description: "MuxCode_descriptor"}},
	{Range: [2]uint64{34, 34}, S: scalar.Uint{Description: "FmxBufferSize_descriptor"}},
	{Range: [2]uint64{35, 35}, S: scalar.Uint{Description: "multiplexbuffer_descriptor"}},
	{Range: [2]uint64{36, 36}, S: scalar.Uint{Description: "content_labeling_descriptor"}},
	{Range: [2]uint64{37, 37}, S: scalar.Uint{Description: "metadata_pointer_descriptor"}},
	{Range: [2]uint64{38, 38}, S: scalar.Uint{Description: "metadata_descriptor"}},
	{Range: [2]uint64{39, 39}, S: scalar.Uint{Description: "metadata_STD_descriptor"}},
	{Range: [2]uint64{40, 40}, S: scalar.Uint{Description: "AVC video descriptor"}},
	{Range: [2]uint64{41, 41}, S: scalar.Uint{Description: "IPMP_descriptor (defined in ISO/IEC 13818-11, MPEG-2 IPMP)"}},
	{Range: [2]uint64{42, 42}, S: scalar.Uint{Description: "AVC timing and HRD descriptor"}},
	{Range: [2]uint64{43, 43}, S: scalar.Uint{Description: "MPEG-2_AAC_audio_descriptor"}},
	{Range: [2]uint64{44, 44}, S: scalar.Uint{Description: "FlexMuxTiming_descriptor"}},
	{Range: [2]uint64{45, 63}, S: scalar.Uint{Description: "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Reserved"}},
	{Range: [2]uint64{64, 255}, S: scalar.Uint{Description: "User Private"}},
}

var tsStreamTypeMap = scalar.UintRangeToScalar{
	{Range: [2]uint64{0x00, 0x00}, S: scalar.Uint{Description: "Reserved"}},
	{Range: [2]uint64{0x01, 0x01}, S: scalar.Uint{Sym: "video", Description: "ISO/IEC 11172-2 Video"}}, // TODO: video_mpeg? codec?
	{Range: [2]uint64{0x02, 0x02}, S: scalar.Uint{Description: "ISO/IEC 13818-2 or ISO/IEC 11172-2"}},
	{Range: [2]uint64{0x03, 0x03}, S: scalar.Uint{Sym: "audio_mpeg1", Description: "ISO/IEC 11172-3 Audio"}},
	{Range: [2]uint64{0x04, 0x04}, S: scalar.Uint{Sym: "audio_mpeg2", Description: "ISO/IEC 13818-3 Audio"}},
	{Range: [2]uint64{0x05, 0x05}, S: scalar.Uint{Description: "ISO/IEC 13818-1 private_sections"}},
	{Range: [2]uint64{0x06, 0x06}, S: scalar.Uint{Description: "ISO/IEC 13818-1 PES packets containing private data"}},
	{Range: [2]uint64{0x07, 0x07}, S: scalar.Uint{Description: "ISO/IEC 13522 MHEG"}},
	{Range: [2]uint64{0x08, 0x08}, S: scalar.Uint{Description: "ISO/IEC 13818-1 Annex A DSM-CC"}},
	{Range: [2]uint64{0x09, 0x09}, S: scalar.Uint{Description: "ITU-T Rec. H.222.1"}},
	{Range: [2]uint64{0x0a, 0x0a}, S: scalar.Uint{Description: "ISO/IEC 13818-6 type A"}},
	{Range: [2]uint64{0x0b, 0x0b}, S: scalar.Uint{Description: "ISO/IEC 13818-6 type B"}},
	{Range: [2]uint64{0x0c, 0x0c}, S: scalar.Uint{Description: "ISO/IEC 13818-6 type C"}},
	{Range: [2]uint64{0x0d, 0x0d}, S: scalar.Uint{Description: "ISO/IEC 13818-6 type D"}},
	{Range: [2]uint64{0x0e, 0x0e}, S: scalar.Uint{Description: "ISO/IEC 13818-1 auxiliary"}},
	{Range: [2]uint64{0x0f, 0x0f}, S: scalar.Uint{Sym: "audio_adts", Description: "ISO/IEC 13818-7 Audio with ADTS transport syntax"}},
	{Range: [2]uint64{0x10, 0x10}, S: scalar.Uint{Description: "ISO/IEC 14496-2 Visual"}},
	{Range: [2]uint64{0x11, 0x11}, S: scalar.Uint{Sym: "audio_latm", Description: "ISO/IEC 14496-3 Audio with the LATM"}},
	{Range: [2]uint64{0x12, 0x12}, S: scalar.Uint{Description: "ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in PES packets"}},
	{Range: [2]uint64{0x13, 0x13}, S: scalar.Uint{Description: "ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in ISO/IEC 14496_sections"}},
	{Range: [2]uint64{0x14, 0x14}, S: scalar.Uint{Description: "ISO/IEC 13818-6 Synchronized Download Protocol"}},
	{Range: [2]uint64{0x15, 0x15}, S: scalar.Uint{Description: "Metadata carried in PES packets"}},
	{Range: [2]uint64{0x16, 0x16}, S: scalar.Uint{Description: "Metadata carried in metadata_sections"}},
	{Range: [2]uint64{0x17, 0x17}, S: scalar.Uint{Description: "Metadata carried in ISO/IEC 13818-6 Data Carousel"}},
	{Range: [2]uint64{0x18, 0x18}, S: scalar.Uint{Description: "Metadata carried in ISO/IEC 13818-6 Object Carousel"}},
	{Range: [2]uint64{0x19, 0x19}, S: scalar.Uint{Description: "Metadata carried in ISO/IEC 13818-6 Synchronized Download Protocol"}},
	{Range: [2]uint64{0x1a, 0x1a}, S: scalar.Uint{Description: "IPMP stream (defined in ISO/IEC 13818-11, MPEG-2 IPMP)"}},
	{Range: [2]uint64{0x1b, 0x1b}, S: scalar.Uint{Sym: "video_avc", Description: "AVC video stream as defined in ITU-T Rec. H.264 | ISO/IEC 14496-10 Video"}},
	{Range: [2]uint64{0x1c, 0x7e}, S: scalar.Uint{Description: "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Reserved"}},
	{Range: [2]uint64{0x7f, 0x7f}, S: scalar.Uint{Description: "IPMP stream"}},
	{Range: [2]uint64{0x80, 0xff}, S: scalar.Uint{Description: "User Private"}},
}

var tsPidMap = scalar.UintRangeToScalar{
	{Range: [2]uint64{pidPAT, pidPAT}, S: scalar.Uint{Sym: "pat", Description: "Program association table"}},
	{Range: [2]uint64{0x0001, 0x0001}, S: scalar.Uint{Sym: "cat", Description: "Conditional access table"}},
	{Range: [2]uint64{0x0002, 0x0002}, S: scalar.Uint{Description: "Transport stream description table"}},
	{Range: [2]uint64{0x0003, 0x0003}, S: scalar.Uint{Description: "IPMP control information table"}},
	{Range: [2]uint64{0x0004, 0x000f}, S: scalar.Uint{Description: "Reserved for future use"}},
	{Range: [2]uint64{0x0010, 0x001f}, S: scalar.Uint{Description: "DVB metadata"}},
	{Range: [2]uint64{0x0010, 0x0010}, S: scalar.Uint{Sym: "nit", Description: "NIT, ST"}},
	{Range: [2]uint64{0x0011, 0x0011}, S: scalar.Uint{Sym: "sdt", Description: "SDT, BAT, ST"}},
	{Range: [2]uint64{0x0012, 0x0012}, S: scalar.Uint{Sym: "eit", Description: "EIT, ST, CIT"}},
	{Range: [2]uint64{0x0013, 0x0013}, S: scalar.Uint{Sym: "rst", Description: "RST, ST"}},
	{Range: [2]uint64{0x0014, 0x0014}, S: scalar.Uint{Sym: "tdt", Description: "TDT, TOT, ST"}},
	{Range: [2]uint64{0x0015, 0x0015}, S: scalar.Uint{Description: "Network synchronization"}},
	{Range: [2]uint64{0x0016, 0x0016}, S: scalar.Uint{Sym: "rnt", Description: "RNT"}},
	{Range: [2]uint64{0x0017, 0x001b}, S: scalar.Uint{Description: "Reserved for future use"}},
	{Range: [2]uint64{0x001c, 0x001c}, S: scalar.Uint{Description: "Inband signalling"}},
	{Range: [2]uint64{0x001d, 0x001d}, S: scalar.Uint{Description: "Measurement"}},
	{Range: [2]uint64{0x001e, 0x001e}, S: scalar.Uint{Sym: "dit", Description: "DIT"}},
	{Range: [2]uint64{0x001f, 0x001f}, S: scalar.Uint{Sym: "sit", Description: "SIT"}},
	{Range: [2]uint64{0x0020, 0x1ffa}, S: scalar.Uint{Description: "Program maps, elementary streams and data"}},
	{Range: [2]uint64{0x1ffb, 0x1ffb}, S: scalar.Uint{Description: "DigiCipher 2/ATSC MGT metadata"}},
	{Range: [2]uint64{0x1ffc, 0x1ffe}, S: scalar.Uint{Description: "Program association table assigned"}},
	{Range: [2]uint64{0x1fff, 0x1fff}, S: scalar.Uint{Description: "Null packet (padding)"}},
}

var tsTableMap = scalar.UintRangeToScalar{
	{Range: [2]uint64{0x00, 0x00}, S: scalar.Uint{Description: "program_association_section"}},
	{Range: [2]uint64{0x01, 0x01}, S: scalar.Uint{Description: "conditional_access_section"}},
	{Range: [2]uint64{0x02, 0x02}, S: scalar.Uint{Description: "program_map_section"}},
	{Range: [2]uint64{0x03, 0x03}, S: scalar.Uint{Description: "transport_stream_description_section"}},
	{Range: [2]uint64{0x04, 0x3f}, S: scalar.Uint{Description: "reserved"}},
	{Range: [2]uint64{0x40, 0x40}, S: scalar.Uint{Description: "network_information_section - actual_network"}},
	{Range: [2]uint64{0x41, 0x41}, S: scalar.Uint{Description: "network_information_section - other_network"}},
	{Range: [2]uint64{0x42, 0x42}, S: scalar.Uint{Sym: "sdt", Description: "service_description_section - actual_transport_stream"}},
	{Range: [2]uint64{0x43, 0x45}, S: scalar.Uint{Description: "reserved for future use"}},
	{Range: [2]uint64{0x46, 0x46}, S: scalar.Uint{Description: "service_description_section - other_transport_stream"}},
	{Range: [2]uint64{0x47, 0x47}, S: scalar.Uint{Description: "to 0x49 reserved for future use"}},
	{Range: [2]uint64{0x4a, 0x4a}, S: scalar.Uint{Description: "bouquet_association_section"}},
	{Range: [2]uint64{0x4b, 0x4d}, S: scalar.Uint{Description: "reserved for future use"}},
	{Range: [2]uint64{0x4e, 0x4e}, S: scalar.Uint{Description: "event_information_section - actual_transport_stream, present/following"}},
	{Range: [2]uint64{0x4f, 0x4f}, S: scalar.Uint{Description: "event_information_section - other_transport_stream, present/following"}},
	{Range: [2]uint64{0x50, 0x5f}, S: scalar.Uint{Description: "event_information_section - actual_transport_stream, schedule"}},
	{Range: [2]uint64{0x60, 0x6f}, S: scalar.Uint{Description: "event_information_section - other_transport_stream, schedule"}},
	{Range: [2]uint64{0x70, 0x70}, S: scalar.Uint{Description: "time_date_section"}},
	{Range: [2]uint64{0x71, 0x71}, S: scalar.Uint{Description: "running_status_section"}},
	{Range: [2]uint64{0x72, 0x72}, S: scalar.Uint{Description: "stuffing_section"}},
	{Range: [2]uint64{0x73, 0x73}, S: scalar.Uint{Description: "time_offset_section"}},
	{Range: [2]uint64{0x74, 0x7d}, S: scalar.Uint{Description: "reserved for future use"}},
	{Range: [2]uint64{0x7e, 0x7e}, S: scalar.Uint{Description: "discontinuity_information_section"}},
	{Range: [2]uint64{0x7f, 0x7f}, S: scalar.Uint{Description: "selection_information_section"}},
	{Range: [2]uint64{0x80, 0xfe}, S: scalar.Uint{Description: "user defined"}},
	{Range: [2]uint64{0xff, 0xff}, S: scalar.Uint{Description: "reserved"}},
}
