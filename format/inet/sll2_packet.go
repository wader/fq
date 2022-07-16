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
	interp.RegisterFormat(decode.Format{
		Name:        format.SLL2_PACKET,
		Description: "Linux cooked capture encapsulation v2",
		Groups:      []string{format.LINK_FRAME},
		Dependencies: []decode.Dependency{
			{Names: []string{format.INET_PACKET}, Group: &sllPacket2InetPacketGroup},
		},
		DecodeFn: decodeSLL2,
	})
}

func decodeSLL2(d *decode.D, in any) any {
	if lfi, ok := in.(format.LinkFrameIn); ok {
		if lfi.Type != format.LinkTypeLINUX_SLL2 {
			d.Fatalf("wrong link type %d", lfi.Type)
		}
	}

	protcolType := d.FieldU16("protocol_type", format.EtherTypeMap, scalar.ActualHex)
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
		_ = d.FieldMustGet("link_address").TryScalarFn(mapUToEtherSym, scalar.ActualHex)
		d.FieldFormatOrRawLen(
			"payload",
			d.BitsLeft(),
			sllPacket2InetPacketGroup,
			format.LinkFrameIn{Type: int(protcolType)},
		)
	default:
		d.FieldRawLen("payload", d.BitsLeft())
	}

	return nil
}
