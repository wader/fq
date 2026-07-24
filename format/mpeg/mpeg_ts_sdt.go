package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_TS_SDT,
		&decode.Format{
			Description: "MPEG TS Service Description Table",
			DecodeFn:    mpegTsSdtDecode,
		})
}

func mpegTsSdtDecode(d *decode.D) any {
	d.FieldU8("table_id", tsTableMap, scalar.UintHex)
	d.FieldU1("syntax_indicator")
	d.FieldU3("reserved0", scalar.UintHex)

	length := d.FieldU12("section_length")
	d.FramedFn(int64(length-4)*8, func(d *decode.D) {
		d.FieldU16("transport_stream_id")
		d.FieldU2("reserved1", scalar.UintHex)
		d.FieldU5("version_number")
		d.FieldU1("current_next_indicator")
		d.FieldU8("section_number")
		d.FieldU8("last_section_number")
		d.FieldU16("original_network_id", scalar.UintHex)
		d.FieldU8("reserved3", scalar.UintHex)
		d.FieldArray("services", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("stream", func(d *decode.D) {
					d.FieldU16("service_id")
					d.FieldU6("reserved0", scalar.UintHex)
					d.FieldBool("eit_schedule_flag")
					d.FieldBool("present_following_flag")
					d.FieldU3("running_status")
					d.FieldBool("free_ca_mode")
					descriptorsLoopLength := d.FieldU12("descriptors_loop_length")
					d.FramedFn(int64(descriptorsLoopLength)*8, func(d *decode.D) {
						// TODO:
						d.FieldRawLen("descriptor", d.BitsLeft())
					})
				})
			}
		})
	})
	d.FieldU32("crc", scalar.UintHex)
	if d.BitsLeft() > 0 {
		d.FieldRawLen("stuffing", d.BitsLeft())
	}
	return nil
}
