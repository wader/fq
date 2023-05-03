package inet

// SLL stands for sockaddr_ll
// https://www.tcpdump.org/linktypes/LINKTYPE_LINUX_SLL.html

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var sllPacketInetPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.SLL_Packet,
		&decode.Format{
			Description: "Linux cooked capture encapsulation",
			Groups:      []*decode.Group{format.Link_Frame},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.INET_Packet}, Out: &sllPacketInetPacketGroup},
			},
			DecodeFn: decodeSLL,
		})
}

var sllPacketTypeMap = scalar.UintMap{
	0: {Sym: "to_us", Description: "Sent to us"},
	1: {Sym: "broadcast", Description: "Broadcast by somebody else"},
	2: {Sym: "multicast", Description: "Multicast by somebody else"},
	3: {Sym: "to_other", Description: "Sent to somebody else by somebody else"},
	4: {Sym: "from_us", Description: "Sent by us"},
}

const (
	arpHdrTypeEther    = 1
	arpHdrTypeLoopback = 772
)

// based on https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_arp.h
var arpHdrTypeMAp = scalar.UintMap{
	0:                  {Sym: "netrom", Description: `from KA9Q: NET/ROM pseudo`},
	arpHdrTypeEther:    {Sym: "ether", Description: `Ethernet 10Mbps`},
	2:                  {Sym: "eether", Description: `Experimental Ethernet`},
	3:                  {Sym: "ax25", Description: `AX.25 Level 2`},
	4:                  {Sym: "pronet", Description: `PROnet token ring`},
	5:                  {Sym: "chaos", Description: `Chaosnet`},
	6:                  {Sym: "ieee802", Description: `IEEE 802.2 Ethernet/TR/TB`},
	7:                  {Sym: "arcnet", Description: `ARCnet`},
	8:                  {Sym: "appletlk", Description: `APPLEtalk`},
	15:                 {Sym: "dlci", Description: `Frame Relay DLCI`},
	19:                 {Sym: "atm", Description: `ATM`},
	23:                 {Sym: "metricom", Description: `Metricom STRIP (new IANA id`},
	24:                 {Sym: "ieee1394", Description: `IEEE 1394 IPv4 - RFC 2734`},
	27:                 {Sym: "eui64", Description: `EUI-64`},
	32:                 {Sym: "infiniband", Description: `InfiniBand`},
	256:                {Sym: "slip"},
	257:                {Sym: "cslip"},
	258:                {Sym: "slip6"},
	259:                {Sym: "cslip6"},
	260:                {Sym: "rsrvd", Description: `Notional KISS type`},
	264:                {Sym: "adapt"},
	270:                {Sym: "rose"},
	271:                {Sym: "x25", Description: `CCITT X.25`},
	272:                {Sym: "hwx25", Description: `Boards with X.25 in firmware`},
	280:                {Sym: "can", Description: `Controller Area Network`},
	290:                {Sym: "mctp"},
	512:                {Sym: "ppp"},
	513:                {Sym: "cisco", Description: `Cisco HDLC`},
	516:                {Sym: "lapb", Description: `LAPB`},
	517:                {Sym: "ddcmp", Description: `Digital's DDCMP protocol`},
	518:                {Sym: "rawhdlc", Description: `Raw HDLC`},
	519:                {Sym: "rawip", Description: `Raw IP`},
	768:                {Sym: "tunnel", Description: `IPIP tunnel`},
	769:                {Sym: "tunnel6", Description: `IP6IP6 tunnel`},
	770:                {Sym: "frad", Description: `Frame Relay Access Device`},
	771:                {Sym: "skip", Description: `SKIP vif`},
	arpHdrTypeLoopback: {Sym: "loopback", Description: `Loopback device`},
	773:                {Sym: "localtlk", Description: `Localtalk device`},
	774:                {Sym: "fddi", Description: `Fiber Distributed Data Interface`},
	775:                {Sym: "bif", Description: `AP1000 BIF`},
	776:                {Sym: "sit", Description: `sit0 device - IPv6-in-IPv4`},
	777:                {Sym: "ipddp", Description: `IP over DDP tunneller`},
	778:                {Sym: "ipgre", Description: `GRE over IP`},
	779:                {Sym: "pimreg", Description: `PIMSM register interface`},
	780:                {Sym: "hippi", Description: `High Performance Parallel Interface`},
	781:                {Sym: "ash", Description: `Nexus 64Mbps Ash`},
	782:                {Sym: "econet", Description: `Acorn Econet`},
	783:                {Sym: "irda", Description: `Linux-IrDA`},
	784:                {Sym: "fcpp", Description: `Point to point fibrechannel`},
	785:                {Sym: "fcal", Description: `Fibrechannel arbitrated loop`},
	786:                {Sym: "fcpl", Description: `Fibrechannel public loop`},
	787:                {Sym: "fcfabric", Description: `Fibrechannel fabric`},
	800:                {Sym: "ieee802_tr", Description: `Magic type ident for TR`},
	801:                {Sym: "ieee80211", Description: `IEEE 802.11`},
	802:                {Sym: "ieee80211_prism", Description: `IEEE 802.11 + Prism2 header`},
	803:                {Sym: "ieee80211_radiotap", Description: `IEEE 802.11 + radiotap header`},
	804:                {Sym: "ieee802154"},
	805:                {Sym: "ieee802154_monitor", Description: `IEEE 802.15.4 network monitor`},
	820:                {Sym: "phonet", Description: `PhoNet media type`},
	821:                {Sym: "phonet_pipe", Description: `PhoNet pipe header`},
	822:                {Sym: "caif", Description: `CAIF media type`},
	823:                {Sym: "ip6gre", Description: `GRE over IPv6`},
	824:                {Sym: "netlink", Description: `Netlink header`},
	825:                {Sym: "6lowpan", Description: `IPv6 over LoWPAN`},
	826:                {Sym: "vsockmon", Description: `Vsock monitor header`},
	0xffff:             {Sym: "void", Description: `Void type, nothing is known`},
	0xfffe:             {Sym: "none", Description: `zero header length`},
}

func decodeSLL(d *decode.D) any {
	var lfi format.Link_Frame_In
	if d.ArgAs(&lfi) && lfi.Type != format.LinkTypeLINUX_SLL {
		d.Fatalf("wrong link type %d", lfi.Type)
	}

	d.FieldU16("packet_type", sllPacketTypeMap)
	arpHdrType := d.FieldU16("arphdr_type", arpHdrTypeMAp)
	addressLength := d.FieldU16("link_address_length")
	d.FieldU("link_address", int(addressLength)*8)
	addressDiff := 8 - addressLength
	if addressDiff > 0 {
		d.FieldRawLen("padding", int64(addressDiff)*8)
	}

	// TODO: handle other arphdr types
	switch arpHdrType {
	case arpHdrTypeLoopback, arpHdrTypeEther:
		_ = d.FieldMustGet("link_address").TryUintScalarFn(mapUToEtherSym, scalar.UintHex)
		protcolType := d.FieldU16("protocol_type", format.EtherTypeMap, scalar.UintHex)
		d.FieldFormatOrRawLen(
			"payload",
			d.BitsLeft(),
			&sllPacketInetPacketGroup,
			format.INET_Packet_In{EtherType: int(protcolType)},
		)
	default:
		d.FieldU16LE("protocol_type")
		d.FieldRawLen("payload", d.BitsLeft())
	}

	return nil
}
