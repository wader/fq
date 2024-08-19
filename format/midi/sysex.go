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
			decodeSysExMessage(d, ctx)
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

func decodeSysExMessage(d *decode.D, ctx *context) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	N := vlq(d)
	if N < 1 {
		ctx.casio = true
	} else {
		d.FieldStrFn("manufacturer", func(d *decode.D) string {
			data := d.BytesLen(1)

			return fmt.Sprintf("%02X", data[0])
		}, manufacturers)

		if N > 1 {
			d.FieldStrFn("data", func(d *decode.D) string {
				data := d.BytesLen(int(N) - 1)

				if data[len(data)-1] == 0xf7 {
					ctx.casio = false
				} else {
					ctx.casio = true
				}

				if data[len(data)-1] == 0xf7 {
					return fmt.Sprintf("%v", data[1:len(data)-1])
				} else {
					return fmt.Sprintf("%v", data[1:])
				}
			})
		}
	}
}

func decodeSysExContinuation(d *decode.D, ctx *context) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldStrFn("data", func(d *decode.D) string {
		data := vlf(d)

		if len(data) > 0 && data[len(data)-1] == 0xf7 {
			ctx.casio = false
			return fmt.Sprintf("%v", data[:len(data)-1])
		} else {
			ctx.casio = true
			return fmt.Sprintf("%v", data)
		}
	})
}

func decodeSysExEscape(d *decode.D, ctx *context) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldStrFn("data", func(d *decode.D) string {
		return fmt.Sprintf("%v", vlf(d))
	})

	ctx.casio = true
}
