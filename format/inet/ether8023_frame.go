package inet

// TODO: move to own package?

import (
	"encoding/binary"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var ether8023FrameInetPacketGroup decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ETHER8023_FRAME,
		Description: "Ethernet 802.3 frame",
		Groups:      []string{format.LINK_FRAME},
		Dependencies: []decode.Dependency{
			{Names: []string{format.INET_PACKET}, Group: &ether8023FrameInetPacketGroup},
		},
		DecodeFn: decodeEthernetFrame,
	})
}

// TODO: move to shared?
var mapUToEtherSym = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], s.ActualU())
	s.Sym = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x", b[2], b[3], b[4], b[5], b[6], b[7])
	return s, nil
})

func decodeEthernetFrame(d *decode.D, in interface{}) interface{} {
	if lfi, ok := in.(format.LinkFrameIn); ok {
		if lfi.Type != format.LinkTypeETHERNET {
			d.Fatalf("wrong link type %d", lfi.Type)
		}
	}

	d.FieldU("destination", 48, mapUToEtherSym, scalar.Hex)
	d.FieldU("source", 48, mapUToEtherSym, scalar.Hex)
	etherType := d.FieldU16("ether_type", format.EtherTypeMap, scalar.Hex)

	if dv, _, _ := d.TryFieldFormatLen(
		"payload",
		d.BitsLeft(),
		ether8023FrameInetPacketGroup,
		format.InetPacketIn{EtherType: int(etherType)}); dv == nil {
		d.FieldRawLen("payload", d.BitsLeft())
	}

	return nil
}
