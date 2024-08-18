package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
)

func decodeSysExEvent(d *decode.D, status uint8, ctx *context) {
	ctx.running = 0x00

	switch {
	case status == 0xf0 && ctx.casio:
		d.Errorf("SysExMessage F0 start byte without terminating F7")

	case status == 0xf0:
		d.FieldStruct("SysExMessage", func(d *decode.D) {
			ctx.casio = decodeSysExMessage(d)
		})

	case status == 0xf7 && ctx.casio:
		d.FieldStruct("SysExContinuation", func(d *decode.D) {
			decodeSysExContinuation(d, ctx)
		})

	case status == 0xf7:
		d.FieldStruct("SysExEscape", func(d *decode.D) {
			decodeSysExEscape(d, ctx)
		})

	default:
		flush(d, "unknown SysEx event (%02x)", status)
	}
}

func decodeSysExMessage(d *decode.D) bool {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	data := vlf(d)

	if len(data) > 0 {
		d.FieldValueStr("manufacturer", fmt.Sprintf("%02x", data[0]))

		if len(data) > 1 && data[len(data)-1] == 0xf7 {
			d.FieldValueStr("data", fmt.Sprintf("%v", data[1:len(data)-1]))
			return false
		} else {
			d.FieldValueStr("data", fmt.Sprintf("%v", data[1:]))
			d.FieldValueBool("more", true)
			return true
		}
	}

	return true
}

func decodeSysExContinuation(d *decode.D, ctx *context) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	data := vlf(d)

	if len(data) > 0 && data[len(data)-1] == 0xf7 {
		d.FieldValueStr("data", fmt.Sprintf("%v", data[:len(data)-1]))
		ctx.casio = false

	} else {
		d.FieldValueStr("data", fmt.Sprintf("%v", data))
		d.FieldValueBool("more", true)
		ctx.casio = true
	}
}

func decodeSysExEscape(d *decode.D, ctx *context) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	data := vlf(d)

	if len(data) > 0 {
		d.FieldValueStr("data", fmt.Sprintf("%v", data))
	}

	ctx.casio = true
}
