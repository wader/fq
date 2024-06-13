// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
package pdu

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	RDP4     = 0x80001
	RDP5     = 0x80004
	RDP10    = 0x80005
	RDP10_1  = 0x80006
	RDP10_2  = 0x80007
	RDP10_3  = 0x80008
	RDP10_4  = 0x80009
	RDP10_5  = 0x8000A
	RDP10_6  = 0x8000B
	RDP10_7  = 0x8000C
	RDP10_8  = 0x8000d
	RDP10_9  = 0x8000e
	RDP10_10 = 0x8000f
)

var RDPVersionMap = scalar.UintMapSymStr{
	RDP4:     "rdp4",
	RDP5:     "rdp5",
	RDP10:    "rdp10",
	RDP10_1:  "rdp10_1",
	RDP10_2:  "rdp10_2",
	RDP10_3:  "rdp10_3",
	RDP10_4:  "rdp10_4",
	RDP10_5:  "rdp10_5",
	RDP10_6:  "rdp10_6",
	RDP10_7:  "rdp10_7",
	RDP10_8:  "rdp10_8",
	RDP10_9:  "rdp10_9",
	RDP10_10: "rdp10_10",
}

const (
	CLIENT_CORE     = 0xC001
	CLIENT_SECURITY = 0xC002
	CLIENT_NETWORK  = 0xC003
	CLIENT_CLUSTER  = 0xC004
)

var clientDataMap = scalar.UintMapSymStr{
	CLIENT_CORE:     "client_core",
	CLIENT_SECURITY: "client_security",
	CLIENT_NETWORK:  "client_network",
	CLIENT_CLUSTER:  "client_cluster",
}

func ParseClientData(d *decode.D, length int64) {
	d.FieldStruct("client_data", func(d *decode.D) {
		header := d.FieldU16("header", clientDataMap)
		data_len := int64(d.FieldU16("length") - 4)

		switch header {
		case CLIENT_CORE:
			ParseClientDataCore(d, data_len)
		case CLIENT_SECURITY:
			ParseClientDataSecurity(d, data_len)
		case CLIENT_NETWORK:
			ParseClientDataNetwork(d, data_len)
		case CLIENT_CLUSTER:
			ParseClientDataCluster(d, data_len)
		default:
			// Assert() once all functions are implemented and tested.
			d.FieldRawLen("data", data_len*8)
			return
		}
	})
}

func ParseClientDataCore(d *decode.D, length int64) {
	d.FieldU32("version", RDPVersionMap)
	d.FieldU16("desktop_width")
	d.FieldU16("desktop_height")
	d.FieldU16("color_depth")
	d.FieldU16("sas_sequence")
	d.FieldU32("keyboard_layout")
	d.FieldU32("client_build")
	d.FieldUTF16LE("client_name", 32, scalar.StrActualTrim("\x00"))
	d.FieldU32("keyboard_type")
	d.FieldU32("keyboard_sub_type")
	d.FieldU32("keyboard_function_key")
	d.FieldRawLen("ime_file_name", 64*8)
	d.FieldRawLen("code_data", 98*8)
}

func ParseClientDataSecurity(d *decode.D, length int64) {
	d.FieldU32("encryption_methods")
	d.FieldU32("ext_encryption_methods")
}

func ParseClientDataNetwork(d *decode.D, length int64) {
	d.FieldU32("channel_count")
	length -= 4
	d.FieldRawLen("channel_def_array", length*8)
}

func ParseClientDataCluster(d *decode.D, length int64) {
	d.FieldU32("flags")
	d.FieldU32("redirected_session_id")
}
