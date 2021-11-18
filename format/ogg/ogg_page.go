package ogg

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/crc"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.OGG_PAGE,
		Description: "OGG page",
		DecodeFn:    pageDecode,
	})
}

func pageDecode(d *decode.D, in interface{}) interface{} {
	p := format.OggPageOut{}
	startPos := d.Pos()

	d.FieldUTF8("capture_pattern", 4, d.AssertStr("OggS"))
	d.FieldU8("stream_structure_version", d.AssertU(0))
	d.FieldU5("unused_flags")
	p.IsLastPage = d.FieldBool("last_page")
	p.IsFirstPage = d.FieldBool("first_page")
	p.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64LE("absolute_granule_position")
	p.StreamSerialNumber = uint32(d.FieldU32LE("stream_serial_number"))
	p.SequenceNo = uint32(d.FieldU32LE("page_sequence_no"))
	d.FieldRawLen("page_checksum", 32, d.RawHexReverse)
	pageSegments := d.FieldU8("page_segments")
	var segmentTable []uint64
	d.FieldArray("segment_table", func(d *decode.D) {
		for i := uint64(0); i < pageSegments; i++ {
			segmentTable = append(segmentTable, d.FieldU8("segment_size"))
		}
	})
	d.FieldArray("segments", func(d *decode.D) {
		for _, ss := range segmentTable {
			p.Segments = append(p.Segments, d.FieldRawLen("segment", int64(ss)*8))
		}
	})
	endPos := d.Pos()

	pageChecksumValue := d.FieldGet("page_checksum")
	pageCRC := &crc.CRC{Bits: 32, Table: crc.Poly04c11db7Table}
	d.MustCopy(pageCRC, d.BitBufRange(startPos, pageChecksumValue.Range.Start-startPos))                      // header before checksum
	d.MustCopy(pageCRC, bytes.NewReader([]byte{0, 0, 0, 0}))                                                  // zero checksum bits
	d.MustCopy(pageCRC, d.BitBufRange(pageChecksumValue.Range.Stop(), endPos-pageChecksumValue.Range.Stop())) // rest of page
	_ = pageChecksumValue.ScalarFn(d.ValidateBitBuf(bitio.ReverseBytes(pageCRC.Sum(nil))))

	return p
}
