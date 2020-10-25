package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:      "ogg_page",
		DecodeFn:  oggDecode,
		SkipProbe: true,
	})
}

type page struct {
	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           []*bitbuf.Buffer
}

func oggDecode(d *decode.D) interface{} {
	p := &page{}

	// TODO: validate bits left
	d.FieldValidateString("capture_pattern", "OggS")
	d.FieldValidateUFn("stream_structure_version", 0, d.U8)
	d.FieldU5("unused_flags")
	p.IsLastPage = d.FieldBool("last_page")
	p.IsFirstPage = d.FieldBool("first_page")
	p.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64LE("absolute_granule_position")
	p.StreamSerialNumber = uint32(d.FieldU32LE("stream_serial_number"))
	p.SequenceNo = uint32(d.FieldU32LE("page_sequence_no"))
	d.FieldU32("page_checksum")
	pageSegments := d.FieldU8("page_segments")
	segmentTable := d.FieldBytesLen("segment_table", int64(pageSegments))

	d.FieldArrayFn("segment", func(d *decode.D) {
		for _, ss := range segmentTable {
			p.Segments = append(p.Segments, d.FieldBitBufLen("segment", int64(ss)*8))
		}
	})

	return p
}
