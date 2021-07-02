package mpeg

import (
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.MPEG_TS,
		ProbeWeight: 10, // make sure to be after gif, both start with 0x47
		Description: "MPEG Transport Stream",
		Groups:      []string{format.PROBE},
		DecodeFn:    tsDecode,
	})
}

// TODO: ts_packet

func tsDecode(d *decode.D, in interface{}) interface{} {
	d.FieldValidateUFn("sync", 0x47, d.U8)
	d.FieldBool("transport_error_indicator")
	d.FieldBool("payload_unit_start")
	d.FieldBool("transport_priority")
	d.FieldU13("pid")
	d.FieldU2("transport_scrambling_control")
	d.FieldU2("adaptation_field_control")
	d.FieldU4("continuity_counter")

	return nil
}
