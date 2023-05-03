package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.ICMP,
		&decode.Format{
			Description: "Internet Control Message Protocol",
			Groups:      []*decode.Group{format.IP_Packet},
			DecodeFn:    decodeICMP,
		})
}

// based on https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol
var icmpTypeMap = scalar.UintMap{
	0:  {Sym: "echo_reply", Description: "Echo reply"},
	3:  {Sym: "unreachable", Description: "Destination network unreachable"},
	4:  {Sym: "source_quench", Description: "Source quench (congestion control)"},
	5:  {Sym: "redirect", Description: "Redirect Datagram for the Network"},
	6:  {Description: "Alternate Host Address"},
	8:  {Sym: "echo_request", Description: "Echo request"},
	9:  {Sym: "router_advertisement", Description: "Router Advertisement"},
	10: {Sym: "router_solicitation", Description: "Router discovery/selection/solicitation"},
	11: {Sym: "time_exceeded", Description: "TTL expired in transit"},
	12: {Sym: "parameter_problem", Description: "Pointer indicates the error"},
	13: {Sym: "timestamp", Description: "Timestamp"},
	14: {Sym: "timestamp_reply", Description: "Timestamp reply"},
	15: {Sym: "information_request", Description: "Information Request"},
	16: {Sym: "information_reply", Description: "Information Reply"},
	17: {Sym: "address_mask_request", Description: "Address Mask Request"},
	18: {Sym: "address_mask_reply", Description: "Address Mask Reply"},
	30: {Sym: "traceroute", Description: "Information Request"},
	31: {Description: "Datagram Conversion Error"},
	32: {Description: "Mobile Host Redirect"},
	33: {Description: "Where-Are-You (originally meant for IPv6)"},
	34: {Description: "Here-I-Am (originally meant for IPv6)"},
	35: {Description: "Mobile Registration Request"},
	36: {Description: "Mobile Registration Reply"},
	37: {Description: "Domain Name Request"},
	38: {Description: "Domain Name Reply"},
	39: {Description: "Simple Key-Management for Internet Protocol"},
	40: {Sym: "photuris"},
	41: {Description: "Experimental icmp for experimental mobility protocols"},
	42: {Sym: "extended_echo_request", Description: "Request Extended Echo"},
	43: {Sym: "extended_echo_reply", Description: "No Error"},
}

var icmpCodeMapMap = map[uint64]scalar.UintMapDescription{
	3: {
		1:  "Destination host unreachable",
		2:  "Destination protocol unreachable",
		3:  "Destination port unreachable",
		4:  "Fragmentation required, and DF flag set",
		5:  "Source route failed",
		6:  "Destination network unknown",
		7:  "Destination host unknown",
		8:  "Source host isolated",
		9:  "Network administratively prohibited",
		10: "Host administratively prohibited",
		11: "Network unreachable for ToS",
		12: "Host unreachable for ToS",
		13: "Communication administratively prohibited",
		14: "Host Precedence Violation",
		15: "Precedence cutoff in effect",
	},
	5: {
		0: "Redirect Datagram for the Network",
		1: "Redirect Datagram for the Host",
		2: "Redirect Datagram for the ToS & network",
		3: "Redirect Datagram for the ToS & host",
	},
	11: {
		0: "TTL expired in transit",
		1: "Fragment reassembly time exceeded",
	},
	12: {
		0: "Pointer indicates the error",
		1: "Missing a required option",
		2: "Bad length",
	},
	43: {
		0: "No Error",
		1: "Malformed Query",
		2: "No Such Interface",
		3: "No Such Table Entry",
		4: "Multiple Interfaces Satisfy Query",
	},
}

func decodeICMP(d *decode.D) any {
	var ipi format.IP_Packet_In
	if d.ArgAs(&ipi) && ipi.Protocol != format.IPv4ProtocolICMP {
		d.Fatalf("incorrect protocol %d", ipi.Protocol)
	}

	typ := d.FieldU8("type", icmpTypeMap)
	d.FieldU8("code", icmpCodeMapMap[typ])
	d.FieldU16("checksum")
	d.FieldRawLen("content", d.BitsLeft())

	return nil
}
