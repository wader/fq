package protobuf

// TODO: move? make collection on known protobuf messages?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var widevineProtoBufGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.ProtobufWidevine,
		&decode.Format{
			Description: "Widevine protobuf",
			DecodeFn:    widevineDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Protobuf}, Out: &widevineProtoBufGroup},
			},
		})
}

func widevineDecode(d *decode.D) any {
	widewinePb := format.ProtoBufMessage{
		1: {Type: format.ProtoBufTypeEnum, Name: "algorithm", Enums: scalar.UintMapSymStr{
			0: "unencrypted",
			1: "aesctr",
		}},
		2: {Type: format.ProtoBufTypeBytes, Name: "key_id"},
		3: {Type: format.ProtoBufTypeString, Name: "provider"},
		4: {Type: format.ProtoBufTypeBytes, Name: "content_id"},
		6: {Type: format.ProtoBufTypeString, Name: "policy"},
		7: {Type: format.ProtoBufTypeUInt32, Name: "crypto_period_index"},
		8: {Type: format.ProtoBufTypeBytes, Name: "grouped_license"},
		9: {Type: format.ProtoBufTypeUInt32, Name: "protection_scheme", Enums: scalar.UintMapSymStr{
			1667591779: "cenc",
			1667392305: "cbc1",
			1667591795: "cens",
			1667392371: "cbcs",
		}},
	}

	d.Format(&widevineProtoBufGroup, format.ProtoBufIn{Message: widewinePb})

	return nil
}
