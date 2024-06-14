// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
package pdu

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	// Message types.
	CB_MONITOR_READY         = 0x0001
	CB_FORMAT_LIST           = 0x0002
	CB_FORMAT_LIST_RESPONSE  = 0x0003
	CB_FORMAT_DATA_REQUEST   = 0x0004
	CB_FORMAT_DATA_RESPONSE  = 0x0005
	CB_TEMP_DIRECTORY        = 0x0006
	CB_CLIP_CAPS             = 0x0007
	CB_FILECONTENTS_REQUEST  = 0x0008
	CB_FILECONTENTS_RESPONSE = 0x0009
	CB_LOCK_CLIPDATA         = 0x000a
	CB_UNLOCK_CLIPDATA       = 0x000b

	// Message flags.
	NONE             = 0
	CB_RESPONSE_OK   = 0x0001
	CB_RESPONSE_FAIL = 0x0002
	CB_ASCII_NAMES   = 0x0004
)

var cbTypesMap = scalar.UintMapSymStr{
	CB_MONITOR_READY:         "cb_monitor_ready",
	CB_FORMAT_LIST:           "cb_format_list",
	CB_FORMAT_LIST_RESPONSE:  "cb_format_list_response",
	CB_FORMAT_DATA_REQUEST:   "cb_format_data_request",
	CB_FORMAT_DATA_RESPONSE:  "cb_format_data_response",
	CB_TEMP_DIRECTORY:        "cb_temp_directory",
	CB_CLIP_CAPS:             "cb_clip_caps",
	CB_FILECONTENTS_REQUEST:  "cb_filecontents_request",
	CB_FILECONTENTS_RESPONSE: "cb_filecontents_response",
	CB_LOCK_CLIPDATA:         "cb_lock_clipdata",
	CB_UNLOCK_CLIPDATA:       "cb_unlock_clipdata",
}

var cbFlagsMap = scalar.UintMapSymStr{
	NONE:             "none",
	CB_RESPONSE_OK:   "cb_response_ok",
	CB_RESPONSE_FAIL: "cb_response_fail",
	CB_ASCII_NAMES:   "cb_ascii_names",
}

var cbParseFnMap = map[uint16]interface{}{
	CB_FORMAT_DATA_RESPONSE: parseCbFormatDataResponse,
}

func ParseClipboardData(d *decode.D, length int64) {
	d.FieldStruct("clipboard_data", func(d *decode.D) {
		msgType := uint16(d.FieldU16("msg_type", cbTypesMap))
		d.FieldU16("msg_flags", cbFlagsMap)
		dataLength := d.FieldU32("data_len")

		cbParser, ok := cbParseFnMap[msgType]
		if ok {
			parseFn, ok := cbParser.(func(d *decode.D, length uint64))
			if ok {
				parseFn(d, dataLength)
				return
			}
		}
		// Assert() once all functions are implemented.
		d.FieldRawLen("data", int64(dataLength*8))
	})
}

func parseCbFormatDataResponse(d *decode.D, length uint64) {
	d.FieldRawLen("data", int64(length*8))
}
