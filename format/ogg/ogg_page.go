package ogg

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/checksum"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.Ogg_Page,
		&decode.Format{
			Description: "OGG page",
			DecodeFn:    pageDecode,
		})
}

func pageDecode(d *decode.D) any {
	p := format.Ogg_Page_Out{}
	startPos := d.Pos()

	d.Endian = decode.LittleEndian

	d.FieldUTF8("capture_pattern", 4, d.StrAssert("OggS"))
	d.FieldU8("version", d.UintAssert(0))
	d.FieldU5("unused_flags")
	p.IsLastPage = d.FieldBool("last_page")
	p.IsFirstPage = d.FieldBool("first_page")
	p.IsContinuedPacket = d.FieldBool("continued_packet")
	d.FieldU64("granule_position")
	p.StreamSerialNumber = uint32(d.FieldU32("bitstream_serial_number"))
	p.SequenceNo = uint32(d.FieldU32("page_sequence_no"))
	d.FieldU32("crc", scalar.UintHex)
	pageSegments := d.FieldU8("page_segments")
	var segmentTable []uint64
	d.FieldArray("segment_table", func(d *decode.D) {
		for i := uint64(0); i < pageSegments; i++ {
			segmentTable = append(segmentTable, d.FieldU8("segment_size"))
		}
	})
	d.FieldArray("segments", func(d *decode.D) {
		for _, ss := range segmentTable {
			bs := d.ReadAllBits(d.FieldRawLen("segment", int64(ss)*8))
			p.Segments = append(p.Segments, bs)
		}
	})
	endPos := d.Pos()

	pageChecksumValue := d.FieldGet("crc")
	pageCRC := &checksum.CRC{Bits: 32, Table: checksum.Poly04c11db7Table}
	d.Copy(pageCRC, bitio.NewIOReader(d.BitBufRange(startPos, pageChecksumValue.Range.Start-startPos)))                      // header before checksum
	d.Copy(pageCRC, bytes.NewReader([]byte{0, 0, 0, 0}))                                                                     // zero checksum bits
	d.Copy(pageCRC, bitio.NewIOReader(d.BitBufRange(pageChecksumValue.Range.Stop(), endPos-pageChecksumValue.Range.Stop()))) // rest of page
	_ = pageChecksumValue.TryUintScalarFn(d.UintValidateBytes(pageCRC.Sum(nil)))

	return p
}
