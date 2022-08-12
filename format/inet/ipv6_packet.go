package inet

import (
	"bytes"
	"net"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/bitioex"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var ipv6IpPacketGroup decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.IPV6_PACKET,
		Description: "Internet protocol v6 packet",
		Groups:      []string{format.INET_PACKET},
		Dependencies: []decode.Dependency{
			{Names: []string{format.IP_PACKET}, Group: &ipv6IpPacketGroup},
		},
		DecodeFn: decodeIPv6,
	})
}

const (
	nextHeaderHopByHop                     = 0
	nextHeaderRouting                      = 43
	nextHeaderFragment                     = 44
	nextHeaderEncapsulatingSecurityPayload = 50
	nextHeaderAuthentication               = 51
	nextHeaderDestination                  = 60
	nextHeaderMobility                     = 135
	nextHeaderHostIdentity                 = 139
	nextHeaderShim6                        = 140
)

// TODO:
// 253	Use for experimentation and testing	[RFC3692][RFC4727]
// 254	Use for experimentation and testing	[RFC3692][RFC4727]

var nextHeaderNames = scalar.UToSymStr{
	nextHeaderHopByHop:                     "hop_by_hop",
	nextHeaderRouting:                      "routing",
	nextHeaderFragment:                     "fragment",
	nextHeaderEncapsulatingSecurityPayload: "encapsulating_security_payload",
	nextHeaderAuthentication:               "authentication",
	nextHeaderDestination:                  "destination",
	nextHeaderMobility:                     "mobility",
	nextHeaderHostIdentity:                 "host_identity",
	nextHeaderShim6:                        "shim6",
}

var nextHeaderMap = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	if isIpv6Option(s.ActualU()) {
		return nextHeaderNames.MapScalar(s)
	}
	return format.IPv4ProtocolMap.MapScalar(s)
})

func isIpv6Option(n uint64) bool {
	switch n {
	case nextHeaderHopByHop,
		nextHeaderRouting,
		nextHeaderFragment,
		nextHeaderEncapsulatingSecurityPayload,
		nextHeaderAuthentication,
		nextHeaderDestination,
		nextHeaderMobility,
		nextHeaderHostIdentity,
		nextHeaderShim6:
		return true
	default:
		return false
	}
}

// from https://www.iana.org/assignments/ipv6-parameters/ipv6-parameters.xhtml#ipv6-parameters-2
var hopByHopTypeNames = scalar.UToSymStr{
	0x00: "pad1",
	0x01: "padn",
	0xc2: "jumbo_payload",
	0x23: "rpl_option",
	0x04: "tunnel_encapsulation_limit",
	0x05: "router_alert",
	0x26: "quick_start",
	0x07: "calipso",
	0x08: "smf_dpd",
	0xc9: "home_address",
	0x8b: "ilnp_nonce",
	0x8c: "line_identification_option",
	0x4d: "deprecated",
	0x6d: "mpl_option",
	0xee: "ip_dff",
	0x0f: "performance_and_diagnostin_metrics",
	0x11: "ioam",
	0x31: "ioam",
}

var mapUToIPv6Sym = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	b := &bytes.Buffer{}
	if _, err := bitioex.CopyBits(b, s.ActualBitBuf()); err != nil {
		return s, err
	}
	s.Sym = net.IP(b.Bytes()).String()
	return s, nil
})

func decodeIPv6(d *decode.D, in any) any {
	if ipi, ok := in.(format.InetPacketIn); ok && ipi.EtherType != format.EtherTypeIPv6 {
		d.Fatalf("incorrect ethertype %d", ipi.EtherType)
	}

	d.FieldU4("version")
	d.FieldU6("ds")
	d.FieldU2("ecn")
	d.FieldU20("flow_label")
	dataLength := d.FieldU16("payload_length")
	nextHeader := d.FieldU8("next_header", nextHeaderMap)
	d.FieldU8("hop_limit")
	d.FieldRawLen("source_address", 128, mapUToIPv6Sym)
	d.FieldRawLen("destination_address", 128, mapUToIPv6Sym)

	extStart := d.Pos()
	if isIpv6Option(nextHeader) {
		// TODO: own format?
		d.FieldArray("extensions", func(d *decode.D) {
			for isIpv6Option(nextHeader) {
				d.FieldStruct("extension", func(d *decode.D) {
					currentHeader := nextHeader
					nextHeader = d.FieldU8("next_header", nextHeaderMap)
					extLen := d.FieldU8("length")
					// whole header not including the first 8 octets
					extLen += 6

					d.FramedFn(int64(extLen)*8, func(d *decode.D) {
						switch currentHeader {
						case nextHeaderHopByHop:
							d.FieldArray("options", func(d *decode.D) {
								for !d.End() {
									d.FieldStruct("option", func(d *decode.D) {
										d.FieldU8("type", hopByHopTypeNames)
										l := d.FieldU8("len")
										d.FieldRawLen("data", int64(l)*8)
									})
								}
							})
						default:
							d.FieldRawLen("payload", d.BitsLeft())
						}
					})
				})
			}
		})
	}
	extEnd := d.Pos()
	extLen := extEnd - extStart

	// TODO: jumbo
	// TODO: nextHeader 59 skip

	payloadLen := int64(dataLength)*8 - extLen
	d.FieldFormatOrRawLen(
		"payload",
		payloadLen,
		ipv4IpPacketGroup,
		format.IPPacketIn{Protocol: int(nextHeader)},
	)

	return nil
}
