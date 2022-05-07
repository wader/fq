package ogg

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/checksum"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
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

	d.Endian = decode.LittleEndian

	d.FieldUTF8("capture_pattern", 4, d.AssertStr("OggS"))
	d.FieldU8("version", d.AssertU(0))
	d.FieldU5("unused_flags")
	p.IsLastPage = d.FieldBool("last_page")
	p.IsFirstPage = d.FieldBool("first_page")
	p.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64("granule_position")
	p.StreamSerialNumber = uint32(d.FieldU32("bitstream_serial_number"))
	p.SequenceNo = uint32(d.FieldU32("page_sequence_no"))
	d.FieldU32("crc", scalar.ActualHex)
	pageSegments := d.FieldU8("page_segments")
	var segmentTable []uint64
	d.FieldArray("segment_table", func(d *decode.D) {
		for i := uint64(0); i < pageSegments; i++ {
			segmentTable = append(segmentTable, d.FieldU8("segment_size"))
		}
	})
	d.FieldArray("segments", func(d *decode.D) {
		for _, ss := range segmentTable {
			bs := d.MustReadAllBits(d.FieldRawLen("segment", int64(ss)*8))
			p.Segments = append(p.Segments, bs)
		}
	})
	endPos := d.Pos()

	pageChecksumValue := d.FieldGet("crc")
	pageCRC := &checksum.CRC{Bits: 32, Table: checksum.Poly04c11db7Table}
	d.MustCopy(pageCRC, bitio.NewIOReader(d.BitBufRange(startPos, pageChecksumValue.Range.Start-startPos)))                      // header before checksum
	d.MustCopy(pageCRC, bytes.NewReader([]byte{0, 0, 0, 0}))                                                                     // zero checksum bits
	d.MustCopy(pageCRC, bitio.NewIOReader(d.BitBufRange(pageChecksumValue.Range.Stop(), endPos-pageChecksumValue.Range.Stop()))) // rest of page
	_ = pageChecksumValue.TryScalarFn(d.ValidateUBytes(pageCRC.Sum(nil)))

	return p
}
