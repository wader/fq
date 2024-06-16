// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
package pdu

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func parseClientInfo(d *decode.D, length int64) {
	d.FieldStruct("client_info", func(d *decode.D) {
		pos := d.Pos()
		var (
			isUnicode bool
			hasNull   bool
			nullN     uint64 = 0
			unicodeN  uint64 = 0
		)
		codePage := d.FieldU32("code_page")
		flags := d.U32()
		d.SeekRel(-4 * 8)
		d.FieldStruct("flags", decodeFlagsFn)

		isUnicode = ((flags & INFO_UNICODE) != 0)
		hasNull = (codePage == 1252 || isUnicode)

		if hasNull {
			nullN = 1
		}
		if isUnicode {
			unicodeN = 2
		}

		domainLength := int(d.FieldU16("domain_length") + nullN*unicodeN)
		usernameLength := int(d.FieldU16("username_length") + nullN*unicodeN)
		passwordLength := int(d.FieldU16("password_length") + nullN*unicodeN)
		alternateShellLength := int(d.FieldU16("alternate_shell_length") + nullN*unicodeN)
		workingDirLength := int(d.FieldU16("working_dir_length") + nullN*unicodeN)

		d.FieldUTF16LE("domain", domainLength, scalar.StrActualTrim("\x00"))
		d.FieldUTF16LE("username", usernameLength, scalar.StrActualTrim("\x00"))
		d.FieldUTF16LE("password", passwordLength, scalar.StrActualTrim("\x00"))
		d.FieldUTF16LE("alternate_shell", alternateShellLength, scalar.StrActualTrim("\x00"))
		d.FieldUTF16LE("working_dir", workingDirLength, scalar.StrActualTrim("\x00"))

		extraLength := length - ((d.Pos() - pos) / 8)
		if extraLength > 0 {
			d.FieldStruct("extra_info", func(d *decode.D) {
				d.FieldU16("address_family", scalar.UintHex)
				addressLength := int(d.FieldU16("address_length"))
				d.FieldUTF16LE("address", addressLength, scalar.StrActualTrim("\x00"))
				clientDirLength := int(d.FieldU16("client_dir_length"))
				d.FieldUTF16LE("client_dir", clientDirLength, scalar.StrActualTrim("\x00"))
				// TS_TIME_ZONE_INFORMATION structure
				// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/526ed635-d7a9-4d3c-bbe1-4e3fb17585f4
				d.FieldU32("timezone_bias")
				d.FieldUTF16LE("timezone_standardname", 64, scalar.StrActualTrim("\x00"))
			})

			// XXX: there's more extra info but here's everything we need from the
			//			client (other than UTC info)
		}
	})
}

const (
	// flags
	INFO_MOUSE                  = 0x00000001
	INFO_DISABLECTRLALTDEL      = 0x00000002
	INFO_AUTOLOGON              = 0x00000008
	INFO_UNICODE                = 0x00000010
	INFO_MAXIMIZESHELL          = 0x00000020
	INFO_LOGONNOTIFY            = 0x00000040
	INFO_COMPRESSION            = 0x00000080
	INFO_ENABLEWINDOWSKEY       = 0x00000100
	INFO_REMOTECONSOLEAUDIO     = 0x00002000
	INFO_FORCE_ENCRYPTED_CS_PDU = 0x00004000
	INFO_RAIL                   = 0x00008000
	INFO_LOGONERRORS            = 0x00010000
	INFO_MOUSE_HAS_WHEEL        = 0x00020000
	INFO_PASSWORD_IS_SC_PIN     = 0x00040000
	INFO_NOAUDIOPLAYBACK        = 0x00080000
	INFO_USING_SAVED_CREDS      = 0x00100000
	INFO_AUDIOCAPTURE           = 0x00200000
	INFO_VIDEO_DISABLE          = 0x00400000
	INFO_RESERVED1              = 0x00800000
	INFO_RESERVED2              = 0x01000000
	INFO_HIDEF_RAIL_SUPPORTED   = 0x02000000
)

func decodeFlagsFn(d *decode.D) {
	d.FieldBool("mouse")
	d.FieldBool("disabledctrlaltdel")
	d.FieldRawLen("unused0", 1)
	d.FieldBool("autologon")
	d.FieldBool("unicode")
	d.FieldBool("maximizeshell")
	d.FieldBool("logonnotify")
	d.FieldBool("compression")
	d.FieldBool("enablewindowskey")
	d.FieldRawLen("unused1", 4)
	d.FieldBool("remoteconsoleaudio")
	d.FieldBool("force_encrypted_cs_pdu")
	d.FieldBool("rail")
	d.FieldBool("logonerrors")
	d.FieldBool("mouse_has_wheel")
	d.FieldBool("password_is_sc_pin")
	d.FieldBool("noaudioplayback")
	d.FieldBool("using_saved_creds")
	d.FieldBool("audiocapture")
	d.FieldBool("video_disable")
	d.FieldBool("reserved1")
	d.FieldBool("reserved2")
	d.FieldBool("hidef_rail_supported")
	d.FieldRawLen("unused2", 6)
}
