package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.ICMPv6,
		&decode.Format{
			Description: "Internet Control Message Protocol v6",
			Groups:      []*decode.Group{format.IP_Packet},
			DecodeFn:    decodeICMPv6,
		})
}

// based on https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol_for_IPv6
var icmpv6TypeMap = scalar.UintMap{
	1:   {Sym: "unreachable", Description: "Destination unreachable"},
	2:   {Sym: "too_big", Description: "Packet too big"},
	3:   {Sym: "time_exceeded", Description: "Time exceeded"},
	4:   {Sym: "parameter_problem", Description: "Parameter problem"},
	100: {Description: "Private experimentation"},
	101: {Description: "Private experimentation"},
	127: {Description: "Reserved for expansion of ICMPv6 error messages"},
	128: {Sym: "echo_reply", Description: "Echo Request"},
	129: {Sym: "echo_request", Description: "Echo Reply"},
	130: {Description: "Multicast Listener Query (MLD)"},
	131: {Description: "Multicast Listener Report (MLD)"},
	132: {Description: "Multicast Listener Done (MLD)"},
	133: {Description: "Router Solicitation (NDP)"},
	134: {Description: "Router Advertisement (NDP)"},
	135: {Description: "Neighbor Solicitation (NDP)"},
	136: {Description: "Neighbor Advertisement (NDP)"},
	137: {Description: "Redirect Message (NDP)"},
	138: {Description: "Router Renumbering	Router Renumbering Command"},
	139: {Description: "ICMP Node Information Query"},
	140: {Description: "ICMP Node Information Response"},
	141: {Description: "Inverse Neighbor Discovery Solicitation Message"},
	142: {Description: "Inverse Neighbor Discovery Advertisement Message"},
	143: {Description: "Multicast Listener Discovery (MLDv2) reports (RFC 3810)"},
	144: {Description: "Home Agent Address Discovery Request Message"},
	145: {Description: "Home Agent Address Discovery Reply Message"},
	146: {Description: "Mobile Prefix Solicitation"},
	147: {Description: "Mobile Prefix Advertisement"},
	148: {Description: "Certification Path Solicitation (SEND)"},
	149: {Description: "Certification Path Advertisement (SEND)"},
	151: {Description: "Multicast Router Advertisement (MRD)"},
	152: {Description: "Multicast Router Solicitation (MRD)"},
	153: {Description: "Multicast Router Termination (MRD)"},
	155: {Description: "RPL Control Message"},
	200: {Description: "Private experimentation"},
	201: {Description: "Private experimentation"},
	255: {Description: "Reserved for expansion of ICMPv6 informational messages"},
}

var icmpv6CodeMapMap = map[uint64]scalar.UintMapDescription{
	1: {
		1: "Communication with destination administratively prohibited",
		2: "Beyond scope of source address",
		3: "Address unreachable",
		4: "Port unreachable",
		5: "Source address failed ingress/egress policy",
		6: "Reject route to destination",
		7: "Error in Source Routing Header",
	},
	3: {
		0: "Hop limit exceeded in transit",
		1: "Fragment reassembly time exceeded",
	},
	4: {
		0: "Erroneous header field encountered",
		1: "Unrecognized Next Header type encountered",
		2: "Unrecognized IPv6 option encountered",
	},
}

func decodeICMPv6(d *decode.D) any {
	var ipi format.IP_Packet_In
	if d.ArgAs(&ipi) && ipi.Protocol != format.IPv4ProtocolICMPv6 {
		d.Fatalf("incorrect protocol %d", ipi.Protocol)
	}

	typ := d.FieldU8("type", icmpv6TypeMap)
	d.FieldU8("code", icmpv6CodeMapMap[typ])
	d.FieldU16("checksum")
	d.FieldRawLen("content", d.BitsLeft())

	return nil
}
