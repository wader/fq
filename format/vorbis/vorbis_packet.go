package vorbis

// https://xiph.org/vorbis/doc/Vorbis_I_spec.html
// TODO: setup? more audio?
// TODO: end padding? byte align?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var vorbisComment decode.Group

func init() {
	interp.RegisterFormat(
		format.Vorbis_Packet,
		&decode.Format{
			Description: "Vorbis packet",
			DecodeFn:    vorbisDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Vorbis_Comment}, Out: &vorbisComment},
			},
		})
}

const (
	packetTypeAudio          = 0
	packetTypeIdentification = 1
	packetTypeComment        = 3
	packetTypeSetup          = 5
)

var packetTypeNames = map[uint]string{
	packetTypeAudio:          "audio",
	packetTypeIdentification: "identification",
	packetTypeComment:        "comment",
	packetTypeSetup:          "setup",
}

func vorbisDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	packetType := d.FieldScalarUintFn("packet_type", func(d *decode.D) scalar.Uint {
		packetTypeName := "unknown"
		t := d.U8()
		// 4.2.1. Common header decode
		// "these types are all odd as a packet with a leading single bit of ’0’ is an audio packet"
		if t&1 == 0 {
			t = packetTypeAudio
		}
		if n, ok := packetTypeNames[uint(t)]; ok {
			packetTypeName = n
		}
		return scalar.Uint{Actual: t, Sym: packetTypeName}
	})

	switch packetType.Actual {
	case packetTypeIdentification, packetTypeSetup, packetTypeComment:
		d.FieldUTF8("magic", 6, d.StrAssert("vorbis"))
	case packetTypeAudio:
	default:
		d.Fatalf("unknown packet type %d", packetType.Actual)
	}

	switch packetType.Actual {
	case packetTypeAudio:
	case packetTypeIdentification:
		// 1   1) [vorbis_version] = read 32 bits as unsigned integer
		// 2   2) [audio_channels] = read 8 bit integer as unsigned
		// 3   3) [audio_sample_rate] = read 32 bits as unsigned integer
		// 4   4) [bitrate_maximum] = read 32 bits as signed integer
		// 5   5) [bitrate_nominal] = read 32 bits as signed integer
		// 6   6) [bitrate_minimum] = read 32 bits as signed integer
		// 7   7) [blocksize_0] = 2 exponent (read 4 bits as unsigned integer)
		// 8   8) [blocksize_1] = 2 exponent (read 4 bits as unsigned integer)
		// 9   9) [framing_flag] = read one bit
		d.FieldU32("vorbis_version", d.UintValidate(0))
		d.FieldU8("audio_channels")
		d.FieldU32("audio_sample_rate")
		d.FieldU32("bitrate_maximum")
		d.FieldU32("bitrate_nominal")
		d.FieldU32("bitrate_minimum")
		// TODO: code/comment about 2.1.4. coding bits into byte sequences
		d.FieldUintFn("blocksize_1", func(d *decode.D) uint64 { return 1 << d.U4() })
		d.FieldUintFn("blocksize_0", func(d *decode.D) uint64 { return 1 << d.U4() })
		// TODO: warning if blocksize0 > blocksize1
		// TODO: warning if not 64-8192
		d.FieldRawLen("padding0", 7, d.BitBufIsZero())
		d.FieldU1("framing_flag", d.UintValidate(1))
	case packetTypeSetup:
		d.FieldUintFn("vorbis_codebook_count", func(d *decode.D) uint64 { return d.U8() + 1 })
		d.FieldU24("codecooke_sync", d.UintValidate(0x564342), scalar.UintHex)
		d.FieldU16("codebook_dimensions")
		d.FieldU24("codebook_entries")

		// d.SeekRel(7)
		// ordered := d.FieldBool("ordered")

		// if ordered {

		// } else {
		// 	d.SeekRel(-2)
		// 	sparse := d.FieldBool("sparse")
		// 	d.SeekRel(1)

		// 	if sparse {

		// 	} else {
		// 		d.SeekRel(-7)
		// 		d.FieldU5("length")

		// 	}
		// }

	case packetTypeComment:
		d.FieldFormat("comment", &vorbisComment, nil)

		// note this uses vorbis bitpacking convention, bits are added LSB first per byte
		d.FieldRawLen("padding0", 7, d.BitBufIsZero())
		d.FieldU1("frame_bit", d.UintValidate(1))
	}

	if d.BitsLeft() > 0 {
		d.FieldRawLen("data", d.BitsLeft())
	}

	return nil
}
