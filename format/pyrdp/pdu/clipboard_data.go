// Copyright (c) 2022-2023 GoSecure Inc.
package pyrdp

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
	CB_LOCK_CLIPDATA         = 0x000A
	CB_UNLOCK_CLIPDATA       = 0x000B

	// Message flags.
	NONE             = 0
	CB_RESPONSE_OK   = 0x0001
	CB_RESPONSE_FAIL = 0x0002
	CB_ASCII_NAMES   = 0x0004
)

// TODO: Fill the descriptions.
var cbTypesMap = scalar.UintMap{
	CB_MONITOR_READY:         {Sym: "cb_monitor_ready", Description: ""},
	CB_FORMAT_LIST:           {Sym: "cb_format_list", Description: ""},
	CB_FORMAT_LIST_RESPONSE:  {Sym: "cb_format_list_response", Description: ""},
	CB_FORMAT_DATA_REQUEST:   {Sym: "cb_format_data_request", Description: ""},
	CB_FORMAT_DATA_RESPONSE:  {Sym: "cb_format_data_response", Description: ""},
	CB_TEMP_DIRECTORY:        {Sym: "cb_temp_directory", Description: ""},
	CB_CLIP_CAPS:             {Sym: "cb_clip_caps", Description: ""},
	CB_FILECONTENTS_REQUEST:  {Sym: "cb_filecontents_request", Description: ""},
	CB_FILECONTENTS_RESPONSE: {Sym: "cb_filecontents_response", Description: ""},
	CB_LOCK_CLIPDATA:         {Sym: "cb_lock_clipdata", Description: ""},
	CB_UNLOCK_CLIPDATA:       {Sym: "cb_unlock_clipdata", Description: ""},
}

var cbFlagsMap = scalar.UintMap{
	NONE:             {Sym: "none", Description: ""},
	CB_RESPONSE_OK:   {Sym: "cb_response_ok", Description: ""},
	CB_RESPONSE_FAIL: {Sym: "cb_response_fail", Description: ""},
	CB_ASCII_NAMES:   {Sym: "cb_ascii_names", Description: ""},
}

var cbParseFnMap = map[uint16]interface{}{
	CB_FORMAT_DATA_RESPONSE: parseCbFormatDataResponse,
}

func ParseClipboardData(d *decode.D, length int64) {
	d.FieldStruct("clipboard_data", func(d *decode.D) {
		msg_type := uint16(d.FieldU16("msg_type", cbTypesMap))
		d.FieldU16("msg_flags", cbFlagsMap)
		data_length := d.FieldU32("data_len")

		if _, ok := cbParseFnMap[msg_type]; ok {
			cbParseFnMap[msg_type].(func(d *decode.D, length uint64))(d, data_length)
			return
		}
		// Assert() once all functions are implemented.
		d.FieldRawLen("data", int64(data_length*8))
	})
}

func parseCbFormatDataResponse(d *decode.D, length uint64) {
	d.FieldRawLen("data", int64(length*8))
}
