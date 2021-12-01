package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ICMP,
		Description: "Internet Control Message Protocol",
		DecodeFn:    decodeICMP,
	})
}

// based on https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol
var icmpTypeMap = scalar.UToScalar{
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

var icmpCodeMapMap = map[uint64]scalar.UToScalar{
	3: {
		1:  {Description: "Destination host unreachable"},
		2:  {Description: "Destination protocol unreachable"},
		3:  {Description: "Destination port unreachable"},
		4:  {Description: "Fragmentation required, and DF flag set"},
		5:  {Description: "Source route failed"},
		6:  {Description: "Destination network unknown"},
		7:  {Description: "Destination host unknown"},
		8:  {Description: "Source host isolated"},
		9:  {Description: "Network administratively prohibited"},
		10: {Description: "Host administratively prohibited"},
		11: {Description: "Network unreachable for ToS"},
		12: {Description: "Host unreachable for ToS"},
		13: {Description: "Communication administratively prohibited"},
		14: {Description: "Host Precedence Violation"},
		15: {Description: "Precedence cutoff in effect"},
	},
	5: {
		0: {Description: "Redirect Datagram for the Network"},
		1: {Description: "Redirect Datagram for the Host"},
		2: {Description: "Redirect Datagram for the ToS & network"},
		3: {Description: "Redirect Datagram for the ToS & host"},
	},
	11: {
		0: {Description: "TTL expired in transit"},
		1: {Description: "Fragment reassembly time exceeded"},
	},
	12: {
		0: {Description: "Pointer indicates the error"},
		1: {Description: "Missing a required option"},
		2: {Description: "Bad length"},
	},
	43: {
		0: {Description: "No Error"},
		1: {Description: "Malformed Query"},
		2: {Description: "No Such Interface"},
		3: {Description: "No Such Table Entry"},
		4: {Description: "Multiple Interfaces Satisfy Query"},
	},
}

func decodeICMP(d *decode.D, in interface{}) interface{} {
	typ := d.FieldU8("type", icmpTypeMap)
	d.FieldU8("code", icmpCodeMapMap[typ])
	d.FieldU16("checksum")
	d.FieldRawLen("content", d.BitsLeft())

	return nil
}
