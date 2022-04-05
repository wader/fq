package protobuf

// TODO: move? make collection on known protobuf messages?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var widevineProtoBufFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PROTOBUF_WIDEVINE,
		Description: "Widevine protobuf",
		DecodeFn:    widevineDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROTOBUF}, Group: &widevineProtoBufFormat},
		},
	})
}

func widevineDecode(d *decode.D, in interface{}) interface{} {
	widewinePb := format.ProtoBufMessage{
		1: {Type: format.ProtoBufTypeEnum, Name: "algorithm", Enums: scalar.UToSymStr{
			0: "unencrypted",
			1: "aesctr",
		}},
		2: {Type: format.ProtoBufTypeBytes, Name: "key_id"},
		3: {Type: format.ProtoBufTypeString, Name: "provider"},
		4: {Type: format.ProtoBufTypeBytes, Name: "content_id"},
		6: {Type: format.ProtoBufTypeString, Name: "policy"},
		7: {Type: format.ProtoBufTypeUInt32, Name: "crypto_period_index"},
		8: {Type: format.ProtoBufTypeBytes, Name: "grouped_license"},
		9: {Type: format.ProtoBufTypeUInt32, Name: "protection_scheme", Enums: scalar.UToSymStr{
			1667591779: "cenc",
			1667392305: "cbc1",
			1667591795: "cens",
			1667392371: "cbcs",
		}},
	}

	d.Format(widevineProtoBufFormat, format.ProtoBufIn{Message: widewinePb})

	return nil
}
