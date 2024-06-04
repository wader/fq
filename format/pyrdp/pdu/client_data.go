// Copyright (c) 2022-2023 GoSecure Inc.
// Licensed under the MIT License
package pyrdp

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

var RDPVersionMap = scalar.UintMap{
	RDP4:     {Sym: "rdp4", Description: "RDP 4"},
	RDP5:     {Sym: "rdp5", Description: "RDP 5"},
	RDP10:    {Sym: "rdp10", Description: "RDP 10"},
	RDP10_1:  {Sym: "rdp10_1", Description: "RDP 10.1"},
	RDP10_2:  {Sym: "rdp10_2", Description: "RDP 10.2"},
	RDP10_3:  {Sym: "rdp10_3", Description: "RDP 10.3"},
	RDP10_4:  {Sym: "rdp10_4", Description: "RDP 10.4"},
	RDP10_5:  {Sym: "rdp10_5", Description: "RDP 10.5"},
	RDP10_6:  {Sym: "rdp10_6", Description: "RDP 10.6"},
	RDP10_7:  {Sym: "rdp10_7", Description: "RDP 10.7"},
	RDP10_8:  {Sym: "rdp10_8", Description: "RDP 10.8"},
	RDP10_9:  {Sym: "rdp10_9", Description: "RDP 10.9"},
	RDP10_10: {Sym: "rdp10_10", Description: "RDP 10.10"},
}

const (
	CLIENT_CORE     = 0xC001
	CLIENT_SECURITY = 0xC002
	CLIENT_NETWORK  = 0xC003
	CLIENT_CLUSTER  = 0xC004
)

// TODO: Fill descriptions.
var clientDataMap = scalar.UintMap{
	CLIENT_CORE:     {Sym: "client_core", Description: ""},
	CLIENT_SECURITY: {Sym: "client_security", Description: ""},
	CLIENT_NETWORK:  {Sym: "client_network", Description: ""},
	CLIENT_CLUSTER:  {Sym: "client_cluster", Description: ""},
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
	d.FieldStrFn("client_name", toTextUTF16Fn(32))
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
