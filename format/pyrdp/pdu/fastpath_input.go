// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
package pdu

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	// Security Flags.
	FASTPATH_INPUT_SECURE_CHECKSUM = 1
	FASTPATH_INPUT_ENCRYPTED       = 2

	// Event codes.
	FASTPATH_INPUT_EVENT_SCANCODE      = 0
	FASTPATH_INPUT_EVENT_MOUSE         = 1
	FASTPATH_INPUT_EVENT_MOUSEX        = 2
	FASTPATH_INPUT_EVENT_SYNC          = 3
	FASTPATH_INPUT_EVENT_UNICODE       = 4
	FASTPATH_INPUT_EVENT_QOE_TIMESTAMP = 6
)

// commented because unused but we should use one-day
//var eventCodesMap = scalar.UintMapSymStr{
//	FASTPATH_INPUT_EVENT_SCANCODE:      "fastpath_input_event_scancode",
//	FASTPATH_INPUT_EVENT_MOUSE:         "fastpath_input_event_mouse",
//	FASTPATH_INPUT_EVENT_MOUSEX:        "fastpath_input_event_mousex",
//	FASTPATH_INPUT_EVENT_SYNC:          "fastpath_input_event_sync",
//	FASTPATH_INPUT_EVENT_UNICODE:       "fastpath_input_event_unicode",
//	FASTPATH_INPUT_EVENT_QOE_TIMESTAMP: "fastpath_input_event_qoe_timestamp",
//}

//var eventFnMap = map[int]interface{}{
//	FASTPATH_INPUT_EVENT_SCANCODE:      parseFastpathInputEventScancode,
//	FASTPATH_INPUT_EVENT_MOUSE:         parseFastpathInputEventMouse,
//	FASTPATH_INPUT_EVENT_MOUSEX:        parseFastpathInputEventMousex,
//	FASTPATH_INPUT_EVENT_SYNC:          parseFastpathInputEventSync,
//	FASTPATH_INPUT_EVENT_UNICODE:       parseFastpathInputEventUnicode,
//	FASTPATH_INPUT_EVENT_QOE_TIMESTAMP: parseFastpathInputEventQoeTimestamp,
//}

//var fastPathInputEventLengthsMap = map[int]int{
//	FASTPATH_INPUT_EVENT_SCANCODE:      2,
//	FASTPATH_INPUT_EVENT_MOUSE:         7,
//	FASTPATH_INPUT_EVENT_MOUSEX:        7,
//	FASTPATH_INPUT_EVENT_SYNC:          1,
//	FASTPATH_INPUT_EVENT_UNICODE:       3,
//	FASTPATH_INPUT_EVENT_QOE_TIMESTAMP: 5,
//}

func parseFastPathInput(d *decode.D, length int64) {
	d.FieldStruct("fastpath_input", func(d *decode.D) {
		// var (
		// 	events uint8 = 1
		// )
		pos := d.Pos()

		d.FieldStruct("input_header", func(d *decode.D) {
			d.FieldU2("action", scalar.UintHex)
			// events = uint8(d.FieldU4("events") & 0xf)
			d.FieldU4("events", scalar.UintHex)
			flags := d.FieldU2("flags", scalar.UintHex)
			if flags&FASTPATH_INPUT_ENCRYPTED != 0 {
				panic("Encrypted fast-path not implemented.")
			}
		})

		inputLength := d.FieldU8("input_length1", scalar.UintHex)
		if inputLength&0x80 != 0 {
			inputLength = ((inputLength & 0x7f) << 8) | d.FieldU8("input_length2", scalar.UintHex)
		}

		// d.FieldU64("data_signature", scalar.Hex)
		// fmt.Fprintf(os.Stderr, "events:%d\n", events)
		// d.FieldArray("events", func(d *decode.D) {
		// 	for ; events > 0; events-- {
		// 		var event_type int
		//
		// 		d.FieldStruct("event_header", func(d *decode.D) {
		// 			event_type = int(d.FieldU3("type", eventCodesMap))
		// 			d.FieldU5("flags")
		// 		})
		//
		// 		if _, ok := eventFnMap[event_type]; !ok {
		// 			// panic("fastpath_input: Unknow event code.\n")
		// 			fmt.Fprint(os.Stderr, "fastpath_input: Unknow event code.\n")
		// 		} else {
		// 			eventFnMap[event_type].(func(d *decode.D))(d)
		// 		}
		// 	}
		// })

		inputLength -= uint64(d.Pos()-pos) / 8
		if inputLength > 0 {
			d.FieldRawLen("data", int64(inputLength*8))
		}
	})
}

//commented because unused but we should use one-day
//func parseFastpathInputEventScancode(d *decode.D) {
//	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/089d362b-31eb-4a1a-b6fa-92fe61bb5dbf
//	d.FieldU8("key_code", CharMapper)
//}

//func parseFastpathInputEventMouse(d *decode.D) {
//	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/16a96ded-b3d3-4468-b993-9c7a51297510
//	d.FieldU16("pointer_flags", scalar.UintHex)
//	d.FieldU16("x")
//	d.FieldU16("y")
//}
//func parseFastpathInputEventMousex(d *decode.D) {
//	// https: //docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/2ef7632f-2f2a-4de7-ab58-2585cedcdf48
//	d.FieldU16("pointer_flags", scalar.UintHex)
//	d.FieldU16("x")
//	d.FieldU16("y")
//}
//func parseFastpathInputEventSync(d *decode.D) {
//	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/6c5d0ef9-4653-4d69-9ba9-09ba3acd660f
//	d.FieldU16("padding")
//	d.FieldU32("toggle_flags")
//}
//func parseFastpathInputEventUnicode(d *decode.D) {
//	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/e7b13e98-d800-42bb-9a1d-6948537d2317
//	d.FieldU16("unicode_code", scalar.UintHex)
//}
//func parseFastpathInputEventQoeTimestamp(d *decode.D) {}
