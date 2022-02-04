package inet

import (
	"encoding/binary"
	"net"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/checksum"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var udpPacketFormat decode.Group
var tcpPacketFormat decode.Group
var icmpFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.IPV4_PACKET,
		Description: "Internet protocol v4 packet",
		Dependencies: []decode.Dependency{
			{Names: []string{format.UDP_DATAGRAM}, Group: &udpPacketFormat},
			{Names: []string{format.TCP_SEGMENT}, Group: &tcpPacketFormat},
			{Names: []string{format.ICMP}, Group: &icmpFormat},
		},
		DecodeFn: decodeIPv4,
	})
}

const (
	ipv4OptionEnd = 0
	ipv4OptionNop = 1
)

var ipv4OptionsMap = scalar.UToScalar{
	ipv4OptionEnd: {Sym: "end", Description: "End of options list"},
	ipv4OptionNop: {Sym: "nop", Description: "No operation"},
	2:             {Description: "Security"},
	3:             {Description: "Loose Source Routing"},
	9:             {Description: "Strict Source Routing"},
	7:             {Description: "Record Route"},
	8:             {Description: "Stream ID"},
	4:             {Description: "Internet Timestamp"},
}

var ipv4ProtocolFormat = map[uint64]*decode.Group{
	format.IPv4ProtocolUDP:  &udpPacketFormat,
	format.IPv4ProtocolTCP:  &tcpPacketFormat,
	format.IPv4ProtocolICMP: &icmpFormat,
}

var mapUToIPv4Sym = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(s.ActualU()))
	s.Sym = net.IP(b[:]).String()
	return s, nil
})

func decodeIPv4(d *decode.D, in interface{}) interface{} {
	d.FieldU4("version")
	ihl := d.FieldU4("ihl")
	d.FieldU6("dscp")
	d.FieldU2("ecn")
	totalLength := d.FieldU16("total_length")
	d.FieldU16("identification")
	d.FieldU1("reserved")
	d.FieldBool("dont_fragment")
	moreFragments := d.FieldBool("more_fragments")
	fragmentOffset := d.FieldU13("fragment_offset")
	d.FieldU8("ttl")
	protocol := d.FieldU8("protocol", format.IPv4ProtocolMap)
	checksumStart := d.Pos()
	d.FieldU16("header_checksum", scalar.Hex)
	checksumEnd := d.Pos()
	d.FieldU32("source_ip", mapUToIPv4Sym, scalar.Hex)
	d.FieldU32("destination_ip", mapUToIPv4Sym, scalar.Hex)
	optionsLen := (int64(ihl) - 5) * 8 * 4
	if optionsLen > 0 {
		d.LenFn(optionsLen, func(d *decode.D) {
			d.FieldArray("options", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("option", func(d *decode.D) {
						d.FieldBool("copied")
						d.FieldU2("class")
						kind := d.FieldU5("number", ipv4OptionsMap)
						switch kind {
						case ipv4OptionEnd, ipv4OptionNop:
						default:
							l := d.FieldU8("length")
							d.FieldRawLen("data", (int64(l-2))*8)
						}
					})
				}
			})
		})
	}
	headerEnd := d.Pos()

	ipv4Checksum := &checksum.IPv4{}
	d.MustCopy(ipv4Checksum, bitio.NewIOReader(d.BitBufRange(0, checksumStart)))
	d.MustCopy(ipv4Checksum, bitio.NewIOReader(d.BitBufRange(checksumEnd, headerEnd-checksumEnd)))
	_ = d.FieldMustGet("header_checksum").TryScalarFn(d.ValidateUBytes(ipv4Checksum.Sum(nil)), scalar.Hex)

	dataLen := int64(totalLength-(ihl*4)) * 8
	g, ok := ipv4ProtocolFormat[protocol]
	if !ok || moreFragments || fragmentOffset > 0 {
		d.FieldRawLen("data", dataLen)
	} else {
		d.FieldFormatLen("data", dataLen, *g, nil)
	}

	return nil
}
