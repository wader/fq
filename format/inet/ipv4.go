package inet

import (
	"encoding/binary"
	"net"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var udpFormat decode.Group
var tcpFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.IPV4,
		Description: "Internet protocol v4",
		Dependencies: []decode.Dependency{
			{Names: []string{format.UDP}, Group: &udpFormat},
			{Names: []string{format.TCP}, Group: &tcpFormat},
		},
		DecodeFn: decodeIPv4,
	})
}

const (
	ipv4ProtocolTCP = 6
	ipv4ProtocolUDP = 17
)

var ipv4ProtocolFormat = map[uint64]*decode.Group{
	ipv4ProtocolUDP: &udpFormat,
	ipv4ProtocolTCP: &tcpFormat,
}

func mapUToIPv4Sym(s decode.Scalar) (decode.Scalar, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(s.ActualU()))
	s.Sym = net.IP(b[:]).String()
	return s, nil
}

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
	protocol := d.FieldU8("protocol", d.MapUToScalar(ipv4ProtocolMap))
	d.FieldU16("header_checksum", d.Hex)
	d.FieldU32("source_ip", mapUToIPv4Sym, d.Hex)
	d.FieldU32("destination_ip", mapUToIPv4Sym, d.Hex)
	if ihl > 5 {
		d.FieldRawLen("options", (int64(ihl)-5)*8*4)
	}

	dataLen := int64(totalLength-(ihl*4)) * 8
	g, ok := ipv4ProtocolFormat[protocol]
	if !ok || moreFragments || fragmentOffset > 0 {
		d.FieldRawLen("data", dataLen)
	} else {
		d.FieldFormatLen("data", dataLen, *g, nil)
	}

	return nil
}
