package pdu

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	// PDU Types.
	TYPE_FAST_PATH_INPUT            = 1  // Ex: scan codes, mouse, etc.
	TYPE_FAST_PATH_OUTPUT           = 2  // Ex: image
	TYPE_CLIENT_INFO                = 3  // Creds on connection
	TYPE_SLOW_PATH_PDU              = 4  // For slow-path PDUs
	TYPE_CONNECTION_CLOSE           = 5  // To advertise the end of the connection
	TYPE_CLIPBOARD_DATA             = 6  // To collect clipboard data
	TYPE_CLIENT_DATA                = 7  // Contains the clientName
	TYPE_MOUSE_MOVE                 = 8  // Mouse move event from the player
	TYPE_MOUSE_BUTTON               = 9  // Mouse button event from the player
	TYPE_MOUSE_WHEEL                = 10 // Mouse wheel event from the player
	TYPE_KEYBOARD                   = 11 // Keyboard event from the player
	TYPE_TEXT                       = 12 // Text event from the player
	TYPE_FORWARDING_STATE           = 13 // Event from the player to change the state of I/O forwarding
	TYPE_BITMAP                     = 14 // Bitmap event from the player
	TYPE_DEVICE_MAPPING             = 15 // Device mapping event notification
	TYPE_DIRECTORY_LISTING_REQUEST  = 16 // Directory listing request from the player
	TYPE_DIRECTORY_LISTING_RESPONSE = 17 // Directory listing response to the player
	TYPE_FILE_DOWNLOAD_REQUEST      = 18 // File download request from the player
	TYPE_FILE_DOWNLOAD_RESPONSE     = 19 // File download response to the player
	TYPE_FILE_DOWNLOAD_COMPLETE     = 20 // File download completion notification to the player
)

var TypesMap = scalar.UintMapSymStr{
	TYPE_FAST_PATH_INPUT:            "fastpath_input",
	TYPE_FAST_PATH_OUTPUT:           "fastpath_output",
	TYPE_CLIENT_INFO:                "client_info",
	TYPE_SLOW_PATH_PDU:              "slow_path_pdu",
	TYPE_CONNECTION_CLOSE:           "connection_close",
	TYPE_CLIPBOARD_DATA:             "clipboard_data",
	TYPE_CLIENT_DATA:                "client_data",
	TYPE_MOUSE_MOVE:                 "mouse_move",
	TYPE_MOUSE_BUTTON:               "mouse_button",
	TYPE_MOUSE_WHEEL:                "mouse_wheel",
	TYPE_KEYBOARD:                   "keyboard",
	TYPE_TEXT:                       "text",
	TYPE_FORWARDING_STATE:           "forwarding_state",
	TYPE_BITMAP:                     "bitmap",
	TYPE_DEVICE_MAPPING:             "device_mapping",
	TYPE_DIRECTORY_LISTING_REQUEST:  "directory_listing_request",
	TYPE_DIRECTORY_LISTING_RESPONSE: "directory_listing_response",
	TYPE_FILE_DOWNLOAD_REQUEST:      "file_download_request",
	TYPE_FILE_DOWNLOAD_RESPONSE:     "file_download_response",
	TYPE_FILE_DOWNLOAD_COMPLETE:     "file_download_complete",
}

func noParse(d *decode.D, length int64) {}

var ParsersMap = map[uint16]interface{}{
	TYPE_FAST_PATH_INPUT: parseFastPathInput,
	// TYPE_FAST_PATH_OUTPUT:           parseFastPathOut,
	TYPE_CLIENT_INFO: parseClientInfo,
	// TYPE_SLOW_PATH_PDU:              parseSlowPathPDU,
	TYPE_CONNECTION_CLOSE: noParse,
	TYPE_CLIPBOARD_DATA:   parseClipboardData,
	TYPE_CLIENT_DATA:      parseClientData,
	// TYPE_MOUSE_MOVE:                 parseMouseMove,
	// TYPE_MOUSE_BUTTON:               parseMouseButton,
	// TYPE_MOUSE_WHEEL:                parseMouseWheel,
	// TYPE_KEYBOARD:                   parseKeyboard,
	// TYPE_TEXT:                       parseText,
	// TYPE_FORWARDING_STATE:           parseForwardingState,
	// TYPE_BITMAP:                     parseBitmap,
	// TYPE_DEVICE_MAPPING:             parseDeviceMapping,
	// TYPE_DIRECTORY_LISTING_REQUEST:  parseDirectoryListingRequest,
	// TYPE_DIRECTORY_LISTING_RESPONSE: parseDirectoryListingResponse,
	// TYPE_FILE_DOWNLOAD_REQUEST:      parseFileDownloadRequest,
	// TYPE_FILE_DOWNLOAD_RESPONSE:     parseFileDownloadResponse,
	// TYPE_FILE_DOWNLOAD_COMPLETE:     parseFileDownloadComplete,
}
