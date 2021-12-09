package inet

// SLL stands for sockaddr_ll
// https://www.tcpdump.org/linktypes/LINKTYPE_LINUX_SLL2.html

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var sllPacket2Ether8023Format decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.SLL2_PACKET,
		Description: "Linux cooked capture encapsulation v2",
		Dependencies: []decode.Dependency{
			{Names: []string{format.ETHER8023_FRAME}, Group: &sllPacket2Ether8023Format},
		},
		DecodeFn: decodeSLL2,
	})
}

var sllPacket2FrameTypeFormat = map[uint64]*decode.Group{
	format.EtherTypeIPv4: &ether8023FrameIPv4Format,
}

func decodeSLL2(d *decode.D, in interface{}) interface{} {
	protcolType := d.FieldU16("protocol_type", format.EtherTypeMap, scalar.Hex)
	d.FieldU16("reserved")
	d.FieldU32("interface_index")
	arpHdrType := d.FieldU16("arphdr_type", arpHdrTypeMAp)
	d.FieldU8("packet_type", sllPacketTypeMap)
	addressLength := d.FieldU8("link_address_length", d.ValidateURange(0, 8))
	// "If there are more than 8 bytes, only the first 8 bytes are present"
	if addressLength > 8 {
		addressLength = 8
	}
	// TODO: maybe skip padding and always read 8 bytes?
	d.FieldU("link_address", int(addressLength)*8)
	addressDiff := 8 - addressLength
	if addressDiff > 0 {
		d.FieldRawLen("padding", int64(addressDiff)*8)
	}

	// TODO: handle other arphdr types
	switch arpHdrType {
	case arpHdrTypeLoopback, arpHdrTypeEther:
		_ = d.FieldMustGet("link_address").TryScalarFn(mapUToEtherSym, scalar.Hex)
		if g, ok := sllPacket2FrameTypeFormat[protcolType]; ok {
			d.FieldFormatLen("data", d.BitsLeft(), *g, nil)
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	default:
		d.FieldRawLen("data", d.BitsLeft())
	}

	return nil
}
