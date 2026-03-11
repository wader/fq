package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_TS_PAT,
		&decode.Format{
			Description: "MPEG TS Program Association Table",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    mpegTsPatDecode,
		})
}

func mpegTsPatDecode(d *decode.D) any {
	mtpo := format.MpegTsPatOut{
		PidMap: map[int]int{},
	}

	d.FieldU8("table_id")
	d.FieldU1("syntax_indicator")
	d.FieldU3("reserved0", scalar.UintHex)
	length := d.FieldU12("section_length")
	d.FramedFn(int64(length-4)*8, func(d *decode.D) {
		d.FieldU16("transport_stream_id")
		d.FieldU2("reserved1", scalar.UintHex)
		d.FieldU5("version_number") // TODO: output?
		d.FieldU1("current_next_indicator")
		d.FieldU8("section_number")
		d.FieldU8("last_section_number")
		d.FieldArray("programs", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("program", func(d *decode.D) {
					programNumber := d.FieldU16("program_number")
					d.FieldU3("reserved", scalar.UintHex)
					switch programNumber {
					case 0:
						d.FieldU13("network_pid", scalar.UintHex)
					default:
						programPid := d.FieldU13("program_map_pid", scalar.UintHex)
						mtpo.PidMap[int(programPid)] = int(programNumber)
					}
				})
			}
		})
	})
	// TODO: move
	d.FieldU32("crc", scalar.UintHex)
	if d.BitsLeft() > 0 {
		d.FieldRawLen("stuffing", d.BitsLeft())
	}

	return mtpo
}
