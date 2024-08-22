package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var sysex = scalar.UintMapSymStr{
	0xf0: "F0",
	0xf7: "F7",
}

func decodeSysExEvent(d *decode.D, status uint8, ctx *context) {
	ctx.running = 0x00

	delta := func(d *decode.D) {
		dt := d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)

		ctx.tick += dt
	}

	switch {
	case status == 0xf0 && ctx.casio:
		d.Errorf("SysExMessage F0 start byte without terminating F7")

	case status == 0xf0:
		d.FieldStruct("SysExMessage", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("sysex", sysex)
			decodeSysExMessage(d, ctx)
		})

	case status == 0xf7 && ctx.casio:
		d.FieldStruct("SysExContinuation", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("sysex", sysex)
			decodeSysExContinuation(d, ctx)
		})

	case status == 0xf7:
		d.FieldStruct("SysExEscape", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("sysex", sysex)
			decodeSysExEscape(d, ctx)
		})

	default:
		flush(d, "unknown SysEx event (%02x)", status)
	}
}

func decodeSysExMessage(d *decode.D, ctx *context) {
	var bytes []uint8

	d.FieldStrFn("bytes", func(d *decode.D) string {
		bytes = vlf(d)
		return fmt.Sprintf("%v", bytes)
	})

	if len(bytes) < 1 {
		ctx.casio = true
	} else {
		id := fmt.Sprintf("%02X", bytes[0])

		d.FieldValueStr("manufacturer", id, manufacturers)

		if len(bytes) > 1 {
			if bytes[len(bytes)-1] == 0xf7 {
				ctx.casio = false
			} else {
				ctx.casio = true
			}

			if bytes[len(bytes)-1] == 0xf7 {
				d.FieldValueStr("data", fmt.Sprintf("%v", bytes[1:len(bytes)-1]))
			} else {
				d.FieldValueStr("data", fmt.Sprintf("%v", bytes[1:]))
			}
		}
	}

	if ctx.casio {
		d.FieldValueBool("continued", true)
	}
}

func decodeSysExContinuation(d *decode.D, ctx *context) {
	d.FieldStrFn("data", func(d *decode.D) string {
		data := vlf(d)

		if len(data) > 0 && data[len(data)-1] == 0xf7 {
			ctx.casio = false
		} else {
			ctx.casio = true
		}

		if len(data) > 0 && data[len(data)-1] == 0xf7 {
			return fmt.Sprintf("%v", data[:len(data)-1])
		} else {
			return fmt.Sprintf("%v", data)
		}
	})
}

func decodeSysExEscape(d *decode.D, ctx *context) {
	d.FieldStrFn("data", func(d *decode.D) string {
		data := vlf(d)

		return fmt.Sprintf("%v", data)
	})

	ctx.casio = false
}
