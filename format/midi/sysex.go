package midi

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var sysex = scalar.UintMapSymStr{
	0xf0: "sysex_message",
	0xf7: "sysex_escape",
}

var sysex_extensions = scalar.UintMapSymStr{
	0xf7: "sysex_continuation",
}

func decodeSysExEvent(d *decode.D, status uint8, ctx *context) {
	ctx.running = 0x00

	delta := func(d *decode.D) {
		ctx.tick += d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)
	}

	switch {

	case status == 0xf0:
		d.FieldStruct("sysex_event", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("event", sysex)
			decodeSysExMessage(d, ctx)
		})

	case status == 0xf7 && ctx.casio:
		d.FieldStruct("sysex_event", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("event", sysex_extensions)
			decodeSysExContinuation(d, ctx)
		})

	case status == 0xf7:
		d.FieldStruct("sysex_event", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("event", sysex)
			decodeSysExEscape(d, ctx)
		})

	default:
		flush(d, "unknown SysEx event (%02x)", status)
	}
}

func decodeSysExMessage(d *decode.D, ctx *context) {
	length := d.FieldUintFn("length", vlq)

	if length*8 > uint64(d.BitsLeft()) {
		d.Fatalf("invalid field length")
	}

	d.FieldStruct("message", func(d *decode.D) {
		d.FieldU8("manufacturer", manufacturersMap)

		if length < 1 {
			ctx.casio = true
			d.FieldValueBool("continued", true)
		} else {
			bytes := d.PeekBytes(int(length - 1))
			N := len(bytes)

			if N > 0 && bytes[N-1] == 0xf7 {
				ctx.casio = false
			} else {
				ctx.casio = true
			}

			if N > 0 && bytes[N-1] == 0xf7 {
				d.FieldRawLen("data", int64(8*(N-1)))
				d.FieldU8("end_of_message")
			} else {
				d.FieldRawLen("data", int64(8*N))
				d.FieldValueBool("continued", true)
			}
		}
	})
}

func decodeSysExContinuation(d *decode.D, ctx *context) {
	length := d.FieldUintFn("length", vlq)

	if length*8 > uint64(d.BitsLeft()) {
		d.Fatalf("invalid field length")
	}

	d.FieldStruct("continuation", func(d *decode.D) {
		if length > 0 {
			bytes := d.PeekBytes(int(length))
			N := len(bytes)

			if N > 0 && bytes[N-1] == 0xf7 {
				ctx.casio = false
			} else {
				ctx.casio = true
			}

			if N > 0 && bytes[N-1] == 0xf7 {
				d.FieldRawLen("data", int64(8*(N-1)))
				d.FieldU8("end_of_message")
			} else {
				d.FieldRawLen("data", int64(8*N))
				d.FieldValueBool("continued", true)
			}
		}
	})
}

func decodeSysExEscape(d *decode.D, ctx *context) {
	length := d.FieldUintFn("length", vlq)

	if length*8 > uint64(d.BitsLeft()) {
		d.Fatalf("invalid field length")
	}

	d.FieldStruct("escape", func(d *decode.D) {
		d.FieldRawLen("data", int64(8*length))
	})

	ctx.casio = false
}
