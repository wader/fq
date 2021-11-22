package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var udpDNSFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.UDP,
		Description: "User datagram protocol",
		Dependencies: []decode.Dependency{
			{Names: []string{format.DNS}, Group: &udpDNSFormat},
		},
		DecodeFn: decodeUDP,
	})
}

const (
	udpPortDNS = 53
)

var udpPortFormat = map[uint64]*decode.Group{
	udpPortDNS: &udpDNSFormat,
}

func decodeUDP(d *decode.D, in interface{}) interface{} {
	soucePort := d.FieldU16("source_port", d.MapUToScalar(udpPortMap))
	destPort := d.FieldU16("destination_port", d.MapUToScalar(udpPortMap))
	length := d.FieldU16("length")
	d.FieldU16("checksum", d.Hex)

	// TODO: prio? src/dst map?
	g := udpPortFormat[soucePort]
	if g == nil {
		g = udpPortFormat[destPort]
	}
	dataLen := int64(length-8) * 8
	if g != nil {
		d.FieldFormatLen("data", dataLen, *g, nil)
	} else {
		d.FieldRawLen("data", dataLen)
	}

	return nil
}
