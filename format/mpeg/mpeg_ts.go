package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_TS,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinFuzzy, // make sure to be after gif, both start with 0x47
			Description: "MPEG Transport Stream",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    tsDecode,
		})
}

// TODO: ts_packet

func tsDecode(d *decode.D) any {
	d.FieldU8("sync", d.UintAssert(0x47), scalar.UintHex)
	d.FieldBool("transport_error_indicator")
	d.FieldBool("payload_unit_start")
	d.FieldBool("transport_priority")
	d.FieldU13("pid")
	d.FieldU2("transport_scrambling_control")
	d.FieldU2("adaptation_field_control")
	d.FieldU4("continuity_counter")

	return nil
}
