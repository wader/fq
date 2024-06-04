// Processes PyRDP replay files
// https://github.com/GoSecure/pyrdp
//
// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
//
// Maintainer: Olivier Bilodeau <olivier.bilodeau@flare.io>
// Author: Lisandro Ubiedo
package pyrdp

import (
	"time"

	"github.com/wader/fq/format"
	pyrdp "github.com/wader/fq/format/pyrdp/pdu"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

const (
	READ_EXTRA = true

	// PDU Types.
	PDU_FAST_PATH_INPUT            = 1  // Ex: scan codes, mouse, etc.
	PDU_FAST_PATH_OUTPUT           = 2  // Ex: image
	PDU_CLIENT_INFO                = 3  // Creds on connection
	PDU_SLOW_PATH_PDU              = 4  // For slow-path PDUs
	PDU_CONNECTION_CLOSE           = 5  // To advertise the end of the connection
	PDU_CLIPBOARD_DATA             = 6  // To collect clipboard data
	PDU_CLIENT_DATA                = 7  // Contains the clientName
	PDU_MOUSE_MOVE                 = 8  // Mouse move event from the player
	PDU_MOUSE_BUTTON               = 9  // Mouse button event from the player
	PDU_MOUSE_WHEEL                = 10 // Mouse wheel event from the player
	PDU_KEYBOARD                   = 11 // Keyboard event from the player
	PDU_TEXT                       = 12 // Text event from the player
	PDU_FORWARDING_STATE           = 13 // Event from the player to change the state of I/O forwarding
	PDU_BITMAP                     = 14 // Bitmap event from the player
	PDU_DEVICE_MAPPING             = 15 // Device mapping event notification
	PDU_DIRECTORY_LISTING_REQUEST  = 16 // Directory listing request from the player
	PDU_DIRECTORY_LISTING_RESPONSE = 17 // Directory listing response to the player
	PDU_FILE_DOWNLOAD_REQUEST      = 18 // File download request from the player
	PDU_FILE_DOWNLOAD_RESPONSE     = 19 // File download response to the player
	PDU_FILE_DOWNLOAD_COMPLETE     = 20 // File download completion notification to the player
)

// TODO: Fill all descriptions.
var pduTypesMap = scalar.UintMap{
	PDU_FAST_PATH_INPUT:            {Sym: "pdu_fastpath_input", Description: ""},
	PDU_FAST_PATH_OUTPUT:           {Sym: "pdu_fastpath_output", Description: ""},
	PDU_CLIENT_INFO:                {Sym: "pdu_client_info", Description: ""},
	PDU_SLOW_PATH_PDU:              {Sym: "pdu_slow_path_pdu", Description: ""},
	PDU_CONNECTION_CLOSE:           {Sym: "pdu_connection_close", Description: ""},
	PDU_CLIPBOARD_DATA:             {Sym: "pdu_clipboard_data", Description: ""},
	PDU_CLIENT_DATA:                {Sym: "pdu_client_data", Description: ""},
	PDU_MOUSE_MOVE:                 {Sym: "pdu_mouse_move", Description: ""},
	PDU_MOUSE_BUTTON:               {Sym: "pdu_mouse_button", Description: ""},
	PDU_MOUSE_WHEEL:                {Sym: "pdu_mouse_wheel", Description: ""},
	PDU_KEYBOARD:                   {Sym: "pdu_keyboard", Description: ""},
	PDU_TEXT:                       {Sym: "pdu_text", Description: ""},
	PDU_FORWARDING_STATE:           {Sym: "pdu_forwarding_state", Description: ""},
	PDU_BITMAP:                     {Sym: "pdu_bitmap", Description: ""},
	PDU_DEVICE_MAPPING:             {Sym: "pdu_device_mapping", Description: ""},
	PDU_DIRECTORY_LISTING_REQUEST:  {Sym: "pdu_directory_listing_request", Description: ""},
	PDU_DIRECTORY_LISTING_RESPONSE: {Sym: "pdu_directory_listing_response", Description: ""},
	PDU_FILE_DOWNLOAD_REQUEST:      {Sym: "pdu_file_download_request", Description: ""},
	PDU_FILE_DOWNLOAD_RESPONSE:     {Sym: "pdu_file_download_response", Description: ""},
	PDU_FILE_DOWNLOAD_COMPLETE:     {Sym: "pdu_file_download_complete", Description: ""},
}

var pduParsersMap = map[uint16]interface{}{
	PDU_FAST_PATH_INPUT: pyrdp.ParseFastPathInput,
	// PDU_FAST_PATH_OUTPUT:           pyrdp.ParseFastPathOut,
	PDU_CLIENT_INFO: pyrdp.ParseClientInfo,
	// PDU_SLOW_PATH_PDU:              pyrdp.ParseSlowPathPDU,
	PDU_CONNECTION_CLOSE: noParse,
	PDU_CLIPBOARD_DATA:   pyrdp.ParseClipboardData,
	PDU_CLIENT_DATA:      pyrdp.ParseClientData,
	// PDU_MOUSE_MOVE:                 pyrdp.ParseMouseMove,
	// PDU_MOUSE_BUTTON:               pyrdp.ParseMouseButton,
	// PDU_MOUSE_WHEEL:                pyrdp.ParseMouseWheel,
	// PDU_KEYBOARD:                   pyrdp.ParseKeyboard,
	// PDU_TEXT:                       pyrdp.ParseText,
	// PDU_FORWARDING_STATE:           pyrdp.ParseForwardingState,
	// PDU_BITMAP:                     pyrdp.ParseBitmap,
	// PDU_DEVICE_MAPPING:             pyrdp.ParseDeviceMapping,
	// PDU_DIRECTORY_LISTING_REQUEST:  pyrdp.ParseDirectoryListingRequest,
	// PDU_DIRECTORY_LISTING_RESPONSE: pyrdp.ParseDirectoryListingResponse,
	// PDU_FILE_DOWNLOAD_REQUEST:      pyrdp.ParseFileDownloadRequest,
	// PDU_FILE_DOWNLOAD_RESPONSE:     pyrdp.ParseFileDownloadResponse,
	// PDU_FILE_DOWNLOAD_COMPLETE:     pyrdp.ParseFileDownloadComplete,
}

func init() {
	interp.RegisterFormat(
		format.PYRDP,
		&decode.Format{
			Description: "PyRDP Replay Files",
			DecodeFn:    decodePYRDP,
		})
}

func decodePYRDP(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldArray("events", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("event", func(d *decode.D) {
				pos := d.Pos()

				size := d.FieldU64("size") // minus the length
				pdu_type := uint16(d.FieldU16("pdu_type", pduTypesMap))
				d.FieldU64("timestamp", timestampMapper)
				pdu_size := int64(size - 18)

				if _, ok := pduParsersMap[pdu_type]; !ok { // catch undeclared parsers
					if pdu_size > 0 {
						d.FieldRawLen("data", int64(pdu_size*8))
					}
					return
				}
				pduParsersMap[uint16(pdu_type)].(func(d *decode.D, length int64))(
					d, pdu_size)

				curr := d.Pos() - pos
				if READ_EXTRA {
					d.FieldRawLen("extra", (int64(size)*8)-curr) // seek whatever is left
				} else {
					d.SeekRel((int64(size) * 8) - curr) // read whatever is left
				}
			})
		}
	})
	return nil
}

func noParse(d *decode.D, length int64) {
	return
}

var timestampMapper = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	s.Sym = time.UnixMilli(int64(s.Actual)).UTC().String()
	return s, nil
})
