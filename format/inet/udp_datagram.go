package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var udpPayloadFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.UDP_DATAGRAM,
		Description: "User datagram protocol",
		Dependencies: []decode.Dependency{
			{Names: []string{format.UDP_PAYLOAD}, Group: &udpPayloadFormat},
		},
		DecodeFn: decodeUDP,
	})
}

func decodeUDP(d *decode.D, in interface{}) interface{} {
	soucePort := d.FieldU16("source_port", format.UDPPortMap)
	destPort := d.FieldU16("destination_port", format.UDPPortMap)
	length := d.FieldU16("length")
	d.FieldU16("checksum", scalar.Hex)

	dataLen := int64(length-8) * 8
	if dv, _, _ := d.TryFieldFormatLen("data", dataLen, udpPayloadFormat, format.UDPPayloadIn{
		SourcePort:      int(soucePort),
		DestinationPort: int(destPort),
	}); dv == nil {
		d.FieldRawLen("data", dataLen)
	}

	// TODO: for checksum need to pass ipv4 pseudo header somehow

	return nil
}
