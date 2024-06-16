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
	"embed"
	"time"

	"github.com/wader/fq/format"
	pyrdp_pdu "github.com/wader/fq/format/pyrdp/pdu"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed pyrdp.md
var pyrdpFS embed.FS

func init() {
	interp.RegisterFormat(
		format.PYRDP,
		&decode.Format{
			Description: "PyRDP Replay Files",
			DecodeFn:    decodePYRDP,
		})
	interp.RegisterFS(pyrdpFS)
}

const (
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

var pduTypesMap = scalar.UintMapSymStr{
	PDU_FAST_PATH_INPUT:            "fastpath_input",
	PDU_FAST_PATH_OUTPUT:           "fastpath_output",
	PDU_CLIENT_INFO:                "client_info",
	PDU_SLOW_PATH_PDU:              "slow_path_pdu",
	PDU_CONNECTION_CLOSE:           "connection_close",
	PDU_CLIPBOARD_DATA:             "clipboard_data",
	PDU_CLIENT_DATA:                "client_data",
	PDU_MOUSE_MOVE:                 "mouse_move",
	PDU_MOUSE_BUTTON:               "mouse_button",
	PDU_MOUSE_WHEEL:                "mouse_wheel",
	PDU_KEYBOARD:                   "keyboard",
	PDU_TEXT:                       "text",
	PDU_FORWARDING_STATE:           "forwarding_state",
	PDU_BITMAP:                     "bitmap",
	PDU_DEVICE_MAPPING:             "device_mapping",
	PDU_DIRECTORY_LISTING_REQUEST:  "directory_listing_request",
	PDU_DIRECTORY_LISTING_RESPONSE: "directory_listing_response",
	PDU_FILE_DOWNLOAD_REQUEST:      "file_download_request",
	PDU_FILE_DOWNLOAD_RESPONSE:     "file_download_response",
	PDU_FILE_DOWNLOAD_COMPLETE:     "file_download_complete",
}

var pduParsersMap = map[uint16]interface{}{
	PDU_FAST_PATH_INPUT: pyrdp_pdu.ParseFastPathInput,
	// PDU_FAST_PATH_OUTPUT:           pyrdp_pdu.ParseFastPathOut,
	PDU_CLIENT_INFO: pyrdp_pdu.ParseClientInfo,
	// PDU_SLOW_PATH_PDU:              pyrdp_pdu.ParseSlowPathPDU,
	PDU_CONNECTION_CLOSE: noParse,
	PDU_CLIPBOARD_DATA:   pyrdp_pdu.ParseClipboardData,
	PDU_CLIENT_DATA:      pyrdp_pdu.ParseClientData,
	// PDU_MOUSE_MOVE:                 pyrdp_pdu.ParseMouseMove,
	// PDU_MOUSE_BUTTON:               pyrdp_pdu.ParseMouseButton,
	// PDU_MOUSE_WHEEL:                pyrdp_pdu.ParseMouseWheel,
	// PDU_KEYBOARD:                   pyrdp_pdu.ParseKeyboard,
	// PDU_TEXT:                       pyrdp_pdu.ParseText,
	// PDU_FORWARDING_STATE:           pyrdp_pdu.ParseForwardingState,
	// PDU_BITMAP:                     pyrdp_pdu.ParseBitmap,
	// PDU_DEVICE_MAPPING:             pyrdp_pdu.ParseDeviceMapping,
	// PDU_DIRECTORY_LISTING_REQUEST:  pyrdp_pdu.ParseDirectoryListingRequest,
	// PDU_DIRECTORY_LISTING_RESPONSE: pyrdp_pdu.ParseDirectoryListingResponse,
	// PDU_FILE_DOWNLOAD_REQUEST:      pyrdp_pdu.ParseFileDownloadRequest,
	// PDU_FILE_DOWNLOAD_RESPONSE:     pyrdp_pdu.ParseFileDownloadResponse,
	// PDU_FILE_DOWNLOAD_COMPLETE:     pyrdp_pdu.ParseFileDownloadComplete,
}

func decodePYRDP(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldArray("events", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("event", func(d *decode.D) {
				pos := d.Pos()

				size := d.FieldU64("size") // minus the length
				pduType := uint16(d.FieldU16("pdu_type", pduTypesMap))
				d.FieldU64("timestamp", scalar.UintActualUnixTimeDescription(time.Millisecond, time.RFC3339Nano))
				pduSize := int64(size - 18)

				pduParser, ok := pduParsersMap[pduType]
				if !ok { // catch undeclared parsers
					if pduSize > 0 {
						d.FieldRawLen("data", pduSize*8)
					}
					return
				}
				parseFn, ok := pduParser.(func(d *decode.D, length int64))
				if !ok {
					return
				}
				parseFn(d, pduSize)

				curr := d.Pos() - pos
				d.FieldRawLen("extra", (int64(size)*8)-curr) // seek whatever is left
			})
		}
	})
	return nil
}

func noParse(d *decode.D, length int64) {}
