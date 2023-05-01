package inet

// TODO: move to own package?

import (
	"encoding/binary"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var ether8023FrameInetPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Ether_8023_Frame,
		&decode.Format{
			Description: "Ethernet 802.3 frame",
			Groups:      []*decode.Group{format.Link_Frame},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.INET_Packet}, Out: &ether8023FrameInetPacketGroup},
			},
			DecodeFn: decodeEthernetFrame,
		})
}

// TODO: move to shared?
var mapUToEtherSym = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], s.Actual)
	s.Sym = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x", b[2], b[3], b[4], b[5], b[6], b[7])
	return s, nil
})

func decodeEthernetFrame(d *decode.D) any {
	var lfi format.Link_Frame_In
	if d.ArgAs(&lfi) {
		if lfi.Type != format.LinkTypeETHERNET {
			d.Fatalf("wrong link type %d", lfi.Type)
		}
	}

	d.FieldU("destination", 48, mapUToEtherSym, scalar.UintHex)
	d.FieldU("source", 48, mapUToEtherSym, scalar.UintHex)
	etherType := d.FieldU16("ether_type", format.EtherTypeMap, scalar.UintHex)

	d.FieldFormatOrRawLen(
		"payload",
		d.BitsLeft(),
		&ether8023FrameInetPacketGroup,
		format.INET_Packet_In{EtherType: int(etherType)},
	)

	return nil
}
