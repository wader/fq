package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_TS_PMT,
		&decode.Format{
			Description: "MPEG TS Program Map Table",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    mpegTsPmtDecode,
		})
}

func mpegTsPmtDecode(d *decode.D) any {
	mtpo := format.MpegTsPmtOut{
		Streams: map[int]format.MpegTsStream{},
	}

	d.FieldU8("table_id")
	d.FieldU1("syntax_indicator")
	d.FieldU3("reserved0", scalar.UintHex)
	length := d.FieldU12("section_length")
	d.FramedFn(int64(length-4)*8, func(d *decode.D) {
		d.FieldU16("program_number")
		d.FieldU2("reserved1", scalar.UintHex)
		d.FieldU5("version_number")
		d.FieldU1("current_next_indicator")
		d.FieldU8("section_number")
		d.FieldU8("last_section_number")
		d.FieldU3("reserved2", scalar.UintHex)
		d.FieldU13("pcr_pid")
		d.FieldU4("reserved3", scalar.UintHex)
		programInfoLength := d.FieldU12("program_info_length")
		d.FramedFn(int64(programInfoLength)*8, func(d *decode.D) {
			d.FieldArray("decriptors", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("decriptor", func(d *decode.D) {
						d.FieldU8("tag", tsStreamTagMap, scalar.UintHex)
						length := d.FieldU8("length")
						// TODO:
						d.FieldRawLen("data", int64(length)*8)
					})
				}
			})
		})
		d.FieldArray("streams", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("stream", func(d *decode.D) {
					streamType := d.FieldU8("stream_type", scalar.UintHex, tsStreamTypeMap)
					d.FieldU3("reserved0", scalar.UintHex)
					streamPid := d.FieldU13("elementary_pid", scalar.UintHex)
					d.FieldU4("reserved1", scalar.UintHex)
					length := d.FieldU12("es_info_length")
					d.FramedFn(int64(length)*8, func(d *decode.D) {
						// TODO:
						d.FieldRawLen("data", d.BitsLeft())
					})

					mtpo.Streams[int(streamPid)] = format.MpegTsStream{Type: int(streamType)}
				})
			}
		})
	})
	d.FieldU32("crc", scalar.UintHex)
	if d.BitsLeft() > 0 {
		d.FieldRawLen("stuffing", d.BitsLeft())
	}

	return mtpo
}
