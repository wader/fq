package protobuf

// TODO: move? make collection on known protobuf messages?

import (
	"fq/format"
	"fq/pkg/decode"
)

var widevineProtoBufFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.PROTOBUF_WIDEVINE,
		Description: "Widevine protobuf",
		DecodeFn:    widevineDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROTOBUF}, Formats: &widevineProtoBufFormat},
		},
	})
}

func widevineDecode(d *decode.D, in interface{}) interface{} {
	widewinePb := format.ProtoBufMessage{
		1: {Type: format.ProtoBufTypeEnum, Name: "algorithm", Enums: map[uint64]string{
			0: "UNENCRYPTED",
			1: "AESCTR",
		}},
		2: {Type: format.ProtoBufTypeBytes, Name: "key_id"},
		3: {Type: format.ProtoBufTypeString, Name: "provider"},
		4: {Type: format.ProtoBufTypeBytes, Name: "content_id"},
		6: {Type: format.ProtoBufTypeString, Name: "policy"},
		7: {Type: format.ProtoBufTypeUInt32, Name: "crypto_period_index"},
		8: {Type: format.ProtoBufTypeBytes, Name: "grouped_license"},
		9: {Type: format.ProtoBufTypeUInt32, Name: "protection_scheme", Enums: map[uint64]string{
			1667591779: "cenc",
			1667392305: "cbc1",
			1667591795: "cens",
			1667392371: "cbcs",
		}},
	}

	d.Decode(widevineProtoBufFormat, decode.FormatOptions{InArg: format.ProtoBufIn{Message: widewinePb}})

	return nil
}
