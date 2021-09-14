package ogg

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/crc"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.OGG_PAGE,
		Description: "OGG page",
		DecodeFn:    pageDecode,
	})
}

func pageDecode(d *decode.D, in interface{}) interface{} {
	p := format.OggPageOut{}
	startPos := d.Pos()

	// TODO: validate bits left
	d.FieldValidateUTF8("capture_pattern", "OggS")
	d.FieldValidateUFn("stream_structure_version", 0, d.U8)
	d.FieldU5("unused_flags")
	p.IsLastPage = d.FieldBool("last_page")
	p.IsFirstPage = d.FieldBool("first_page")
	p.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64LE("absolute_granule_position")
	p.StreamSerialNumber = uint32(d.FieldU32LE("stream_serial_number"))
	p.SequenceNo = uint32(d.FieldU32LE("page_sequence_no"))
	d.FieldU32LE("page_checksum")
	pageSegments := d.FieldU8("page_segments")
	var segmentTable []uint64
	d.FieldArrayFn("segment_table", func(d *decode.D) {
		for i := uint64(0); i < pageSegments; i++ {
			segmentTable = append(segmentTable, d.FieldU8("segment_size"))
		}
	})
	d.FieldArrayFn("segments", func(d *decode.D) {
		for _, ss := range segmentTable {
			p.Segments = append(p.Segments, d.FieldBitBufLen("segment", int64(ss)*8))
		}
	})
	endPos := d.Pos()

	pageChecksum := d.FieldMustRemove("page_checksum")
	pageCRC := &crc.CRC{Bits: 32, Table: crc.Poly04c11db7Table}
	decode.MustCopy(pageCRC, d.BitBufRange(startPos, pageChecksum.Range.Start-startPos))                 // header before checksum
	decode.MustCopy(pageCRC, bytes.NewReader([]byte{0, 0, 0, 0}))                                        // zero checksum bits
	decode.MustCopy(pageCRC, d.BitBufRange(pageChecksum.Range.Stop(), endPos-pageChecksum.Range.Stop())) // rest of page
	d.FieldChecksumRange("page_checksum", pageChecksum.Range.Start, pageChecksum.Range.Len, pageCRC.Sum(nil), decode.LittleEndian)

	return p
}
