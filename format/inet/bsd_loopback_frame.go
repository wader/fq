package inet

// TODO: rename NetworkLayer? wireshark calls it "Family", pcap-linktype(7) calls it "network-layer protocol"

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var bsdLoopbackFrameInetPacketGroup decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BSD_LOOPBACK_FRAME,
		Description: "BSD loopback frame",
		Groups:      []string{format.LINK_FRAME},
		Dependencies: []decode.Dependency{
			{Names: []string{format.INET_PACKET}, Group: &bsdLoopbackFrameInetPacketGroup},
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

var bsdLookbackNetworkLayerMap = scalar.UToScalar{
	bsdLoopbackNetworkLayerIPv4: {Sym: "ipv4", Description: `Internet protocol v4`},
	bsdLoopbackNetworkLayerIPv6: {Sym: "ipv6", Description: `Internet protocol v6`},
}

func decodeLoopbackFrame(d *decode.D, in interface{}) interface{} {
	if lfi, ok := in.(format.LinkFrameIn); ok {
		if lfi.Type != format.LinkTypeNULL {
			d.Fatalf("wrong link type %d", lfi.Type)
		}
		// TODO: where is this documented?
		if lfi.IsLittleEndian {
			d.Endian = decode.LittleEndian
		}
	}
	// if no LinkFrameIn assume big endian for now

	networkLayer := d.FieldU32("network_layer", bsdLookbackNetworkLayerMap, scalar.ActualHex)

	d.FieldFormatOrRawLen(
		"payload",
		d.BitsLeft(),
		bsdLoopbackFrameInetPacketGroup,
		// TODO: unknown mapped to ether type 0 is ok?
		format.InetPacketIn{EtherType: bsdLoopbackFrameNetworkLayerEtherType[networkLayer]},
	)

	return nil
}
