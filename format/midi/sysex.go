package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
)

func decodeSysExEvent(d *decode.D, status uint8, casio *bool) {
	switch {
	case status == 0xf0 && *casio:
		d.Errorf("SysExMessage F0 start byte without terminating F7")
		return

	case status == 0xf0:
		d.FieldStruct("SysExMessage", func(d *decode.D) {
			*casio = decodeSysExMessage(d)
		})

	case status == 0xf7 && *casio:
		d.FieldStruct("SysExContinuation", func(d *decode.D) {
			*casio = decodeSysExContinuation(d)
		})

	case status == 0xf7:
		d.FieldStruct("SysExEscape", func(d *decode.D) {
			*casio = decodeSysExEscape(d)
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

	return false
}

func decodeSysExContinuation(d *decode.D) bool {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	data := vlf(d)

	if len(data) > 0 && data[len(data)-1] == 0xf7 {
		d.FieldValueStr("data", fmt.Sprintf("%v", data[:len(data)-1]))
		return false

	} else {
		d.FieldValueStr("data", fmt.Sprintf("%v", data))
		d.FieldValueBool("more", true)
		return true
	}
}

func decodeSysExEscape(d *decode.D) bool {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")

	data := vlf(d)

	if len(data) > 0 {
		d.FieldValueStr("data", fmt.Sprintf("%v", data))
	}

	return false
}
