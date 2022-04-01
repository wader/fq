package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var udpPayloadGroup decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.UDP_DATAGRAM,
		Description: "User datagram protocol",
		Groups:      []string{format.IP_PACKET},
		Dependencies: []decode.Dependency{
			{Names: []string{format.UDP_PAYLOAD}, Group: &udpPayloadGroup},
		},
		DecodeFn: decodeUDP,
	})
}

func decodeUDP(d *decode.D, in interface{}) interface{} {
	if ipi, ok := in.(format.IPPacketIn); ok && ipi.Protocol != format.IPv4ProtocolUDP {
		d.Fatalf("incorrect protocol %d", ipi.Protocol)
	}

	sourcePort := d.FieldU16("source_port", format.UDPPortMap)
	destPort := d.FieldU16("destination_port", format.UDPPortMap)
	length := d.FieldU16("length")
	d.FieldU16("checksum", scalar.Hex)

	payloadLen := int64(length-8) * 8
	if dv, _, _ := d.TryFieldFormatLen(
		"payload",
		payloadLen,
		udpPayloadGroup,
		format.UDPPayloadIn{
			SourcePort:      int(sourcePort),
			DestinationPort: int(destPort),
		}); dv == nil {
		d.FieldRawLen("payload", payloadLen)
	}

	// TODO: for checksum need to pass ipv4 pseudo header somehow

	return nil
}
