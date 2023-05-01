package inet

// SLL stands for sockaddr_ll
// https://www.tcpdump.org/linktypes/LINKTYPE_LINUX_SLL2.html

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var sllPacket2InetPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.SLL2_Packet,
		&decode.Format{
			Description: "Linux cooked capture encapsulation v2",
			Groups:      []*decode.Group{format.Link_Frame},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.INET_Packet}, Out: &sllPacket2InetPacketGroup},
			},
			DecodeFn: decodeSLL2,
		})
}

func decodeSLL2(d *decode.D) any {
	var lfi format.Link_Frame_In
	if d.ArgAs(&lfi) && lfi.Type != format.LinkTypeLINUX_SLL2 {
		d.Fatalf("wrong link type %d", lfi.Type)
	}

	protcolType := d.FieldU16("protocol_type", format.EtherTypeMap, scalar.UintHex)
	d.FieldU16("reserved")
	d.FieldU32("interface_index")
	arpHdrType := d.FieldU16("arphdr_type", arpHdrTypeMAp)
	d.FieldU8("packet_type", sllPacketTypeMap)
	addressLength := d.FieldU8("link_address_length", d.UintValidateRange(0, 8))
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
		_ = d.FieldMustGet("link_address").TryUintScalarFn(mapUToEtherSym, scalar.UintHex)
		d.FieldFormatOrRawLen(
			"payload",
			d.BitsLeft(),
			&sllPacket2InetPacketGroup,
			format.INET_Packet_In{EtherType: int(protcolType)},
		)
	default:
		d.FieldRawLen("payload", d.BitsLeft())
	}

	return nil
}
