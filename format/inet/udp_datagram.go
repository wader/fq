package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var udpPayloadGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.UDP_Datagram,
		&decode.Format{
			Description: "User datagram protocol",
			Groups:      []*decode.Group{format.IP_Packet},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.UDP_Payload}, Out: &udpPayloadGroup},
			},
			DecodeFn: decodeUDP,
		})
}

func decodeUDP(d *decode.D) any {
	var ipi format.IP_Packet_In
	if d.ArgAs(&ipi) && ipi.Protocol != format.IPv4ProtocolUDP {
		d.Fatalf("incorrect protocol %d", ipi.Protocol)
	}

	sourcePort := d.FieldU16("source_port", format.UDPPortMap)
	destPort := d.FieldU16("destination_port", format.UDPPortMap)
	length := d.FieldU16("length")
	d.FieldU16("checksum", scalar.UintHex)

	payloadLen := int64(length-8) * 8
	d.FieldFormatOrRawLen(
		"payload",
		payloadLen,
		&udpPayloadGroup,
		format.UDP_Payload_In{
			SourcePort:      int(sourcePort),
			DestinationPort: int(destPort),
		},
	)

	// TODO: for checksum need to pass ipv4 pseudo header somehow

	return nil
}
