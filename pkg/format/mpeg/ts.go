package mpeg

import (
	"bytes"
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_TS,
		Description: "MPEG Transport Stream",
		Groups:      []string{format.PROBE},
		DecodeFn:    tsDecode,
	})
}

// TODO: ts_packet

func tsDecode(d *decode.D, in interface{}) interface{} {
	gifHeader := []byte("GIF89a")
	if d.BitsLeft()*8 >= int64(len(gifHeader)) && bytes.Equal(d.PeekBytes(len(gifHeader)), []byte("GIF89a")) {
		d.Invalid("looks like GIF")
	}

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
