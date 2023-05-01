package inet

// TODO: rename NetworkLayer? wireshark calls it "Family", pcap-linktype(7) calls it "network-layer protocol"

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var bsdLoopbackFrameInetPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.BSD_Loopback_Frame, &decode.Format{
			Description: "BSD loopback frame",
			Groups:      []*decode.Group{format.Link_Frame},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.INET_Packet}, Out: &bsdLoopbackFrameInetPacketGroup},
			},
			DecodeFn: decodeLoopbackFrame,
		})
}

const (
	bsdLoopbackNetworkLayerIPv4 = 0x2
	bsdLoopbackNetworkLayerIPv6 = 0x1e
)

var bsdLoopbackFrameNetworkLayerEtherType = map[uint64]int{
	bsdLoopbackNetworkLayerIPv4: format.EtherTypeIPv4,
	bsdLoopbackNetworkLayerIPv6: format.EtherTypeIPv6,
}

var bsdLookbackNetworkLayerMap = scalar.UintMap{
	bsdLoopbackNetworkLayerIPv4: {Sym: "ipv4", Description: `Internet protocol v4`},
	bsdLoopbackNetworkLayerIPv6: {Sym: "ipv6", Description: `Internet protocol v6`},
}

func decodeLoopbackFrame(d *decode.D) any {
	var lfi format.Link_Frame_In
	if d.ArgAs(&lfi) {
		if lfi.Type != format.LinkTypeNULL {
			d.Fatalf("wrong link type %d", lfi.Type)
		}
		// TODO: where is this documented?
		if lfi.IsLittleEndian {
			d.Endian = decode.LittleEndian
		}
	}
	// if no LinkFrameIn assume big endian for now

	networkLayer := d.FieldU32("network_layer", bsdLookbackNetworkLayerMap, scalar.UintHex)

	d.FieldFormatOrRawLen(
		"payload",
		d.BitsLeft(),
		&bsdLoopbackFrameInetPacketGroup,
		// TODO: unknown mapped to ether type 0 is ok?
		format.INET_Packet_In{EtherType: bsdLoopbackFrameNetworkLayerEtherType[networkLayer]},
	)

	return nil
}
