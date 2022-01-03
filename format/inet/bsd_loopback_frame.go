package inet

// TODO: rename NetworkLayer? wireshark calls it "Family", pcap-linktype(7) calls it "network-layer protocol"

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var bsdLoopbackFrameIPv4Format decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BSD_LOOPBACK_FRAME,
		Description: "BSD loopback frame",
		Groups:      []string{format.LINK_FRAME},
		Dependencies: []decode.Dependency{
			{Names: []string{format.IPV4_PACKET}, Group: &bsdLoopbackFrameIPv4Format},
		},
		DecodeFn: decodeLoopbackFrame,
	})
}

const (
	bsdLoopbackNetworkLayerIPv4 = 2
)

var bsdLoopbackFrameNetworkLayerFormat = map[uint64]*decode.Group{
	bsdLoopbackNetworkLayerIPv4: &bsdLoopbackFrameIPv4Format,
}

var bsdLookbackNetworkLayerMap = scalar.UToScalar{
	bsdLoopbackNetworkLayerIPv4: {Sym: "ipv4", Description: `Internet protocol v4`},
}

func decodeLoopbackFrame(d *decode.D, in interface{}) interface{} {
	lsi, ok := in.(format.LinkFrameIn)
	if ok {
		if lsi.Type != format.LinkTypeNULL {
			d.Fatalf("wrong link type")
		}
		if lsi.LittleEndian {
			d.Endian = decode.LittleEndian
		}
	}
	// if no LinkFrameIn assume big endian for now

	networkLayer := d.FieldU32("network_layer", bsdLookbackNetworkLayerMap, scalar.Hex)
	if g, ok := bsdLoopbackFrameNetworkLayerFormat[networkLayer]; ok {
		d.FieldFormatLen("packet", d.BitsLeft(), *g, nil)
	} else {
		d.FieldRawLen("data", d.BitsLeft())
	}

	return nil
}
