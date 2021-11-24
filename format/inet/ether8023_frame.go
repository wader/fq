package inet

// TODO: move to own package?

import (
	"encoding/binary"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var ether8023FrameIPv4Format decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ETHER8023_FRAME,
		Description: "Ethernet 802.3 frame",
		Dependencies: []decode.Dependency{
			{Names: []string{format.IPV4_PACKET}, Group: &ether8023FrameIPv4Format},
		},
		DecodeFn: decodeEthernet,
	})
}

var ether8023FrameTypeFormat = map[uint64]*decode.Group{
	format.EtherTypeIPv4: &ether8023FrameIPv4Format,
}

// TODO: move to shared?
func mapUToEtherSym(s decode.Scalar) (decode.Scalar, error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], s.ActualU())
	s.Sym = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x", b[2], b[3], b[4], b[5], b[6], b[7])
	return s, nil
}

func decodeEthernet(d *decode.D, in interface{}) interface{} {
	d.FieldU("destination", 48, mapUToEtherSym, d.Hex)
	d.FieldU("source", 48, mapUToEtherSym, d.Hex)
	etherType := d.FieldU16("ether_type", d.MapUToScalar(format.EtherTypeMap), d.Hex)
	if g, ok := ether8023FrameTypeFormat[etherType]; ok {
		d.FieldFormatLen("packet", d.BitsLeft(), *g, nil)
	} else {
		d.FieldRawLen("data", d.BitsLeft())
	}

	return nil
}
