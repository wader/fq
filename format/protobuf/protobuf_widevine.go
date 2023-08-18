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
			// FourCC
			0x63_65_6e_63: "cenc",
			0x63_62_63_31: "cbc1",
			0x63_65_6e_73: "cens",
			0x63_62_63_73: "cbcs",
		}},
	}

	d.Format(&widevineProtoBufGroup, format.Protobuf_In{Message: widewinePb})

	return nil
}
