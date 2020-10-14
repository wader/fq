package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/format"
)

var Page = format.MustRegister(&decode.Format{
	Name: "ogg_page",
	New:  func() decode.Decoder { return &PageDecoder{} },
})

// Decoder is a ogg page decoder
type PageDecoder struct {
	decode.Common

	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           []*bitbuf.Buffer
}

// Decode ogg page
func (d *PageDecoder) Decode() {
	// TODO: validate bits left
	d.FieldValidateString("capture_pattern", "OggS")
	d.FieldValidateUFn("stream_structure_version", 0, d.U8)
	d.FieldU5("unused_flags")
	d.IsLastPage = d.FieldBool("last_page")
	d.IsFirstPage = d.FieldBool("first_page")
	d.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64LE("absolute_granule_position")
	d.StreamSerialNumber = uint32(d.FieldU32LE("stream_serial_number"))
	d.SequenceNo = uint32(d.FieldU32LE("page_sequence_no"))
	d.FieldU32("page_checksum")
	pageSegments := d.FieldU8("page_segments")
	segmentTable := d.FieldBytesLen("segment_table", int64(pageSegments))

	for _, ss := range segmentTable {
		d.Segments = append(d.Segments, d.FieldBitBufLen("segment", int64(ss)*8))
	}
}
