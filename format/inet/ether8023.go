package inet

// TODO: move to own package?

import (
	"encoding/binary"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var ipv4Format decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.ETHER8023,
		Description: "Ethernet 802.3",
		Dependencies: []decode.Dependency{
			{Names: []string{format.IPV4}, Group: &ipv4Format},
		},
		DecodeFn: decodeEthernet,
	})
}

const (
	etherTypeIPv4 = 0x0800
)

// from https://en.wikipedia.org/wiki/EtherType
// TODO: cleanup
var etherTypeMap = decode.UToScalar{
	etherTypeIPv4: {Sym: "ipv4", Description: `Internet Protocol version 4`},
	0x0806:        {Sym: "arp", Description: `Address Resolution Protocol`},
	0x0842:        {Sym: "wake", Description: `Wake-on-LAN[9]`},
	0x22f0:        {Sym: "audio", Description: `Audio Video Transport Protocol`},
	0x22f3:        {Sym: "trill", Description: `IETF TRILL Protocol`},
	0x22ea:        {Sym: "srp", Description: `Stream Reservation Protocol`},
	0x6002:        {Sym: "dec", Description: `DEC MOP RC`},
	0x6003:        {Sym: "decnet", Description: `DECnet Phase IV, DNA Routing`},
	0x6004:        {Sym: "declat", Description: `DEC LAT`},
	0x8035:        {Sym: "Reverse", Description: `Reverse Address Resolution Protocol`},
	0x809b:        {Sym: "appletalk", Description: `AppleTalk`},
	0x80f3:        {Sym: "appletalk_arp", Description: `AppleTalk Address Resolution Protocol`},
	0x8100:        {Sym: "vlan", Description: `VLAN-tagged (IEEE 802.1Q)`},
	0x8102:        {Sym: "slpp", Description: `Simple Loop Prevention Protocol`},
	0x8103:        {Sym: "vlacp", Description: `Virtual Link Aggregation Control Protocol`},
	0x8137:        {Sym: "ipx", Description: `IPX`},
	0x8204:        {Sym: "qnx", Description: `QNX Qnet`},
	0x86dd:        {Sym: "ipv6", Description: `Internet Protocol Version 6`},
	0x8808:        {Sym: "flow_control", Description: `Ethernet flow control`},
	0x8809:        {Sym: "lacp", Description: `Ethernet Slow Protocols] such as the Link Aggregation Control Protocol`},
	0x8819:        {Sym: "cobranet", Description: `CobraNet`},
	0x8847:        {Sym: "mpls", Description: `MPLS unicast`},
	0x8848:        {Sym: "mpls", Description: `MPLS multicast`},
	0x8863:        {Sym: "pppoe_discovery", Description: `PPPoE Discovery Stage`},
	0x8864:        {Sym: "pppoe_session", Description: `PPPoE Session Stage`},
	0x887b:        {Sym: "homeplug", Description: `HomePlug 1.0 MME`},
	0x888e:        {Sym: "eap", Description: `EAP over LAN (IEEE 802.1X)`},
	0x8892:        {Sym: "profinet", Description: `PROFINET Protocol`},
	0x889a:        {Sym: "hyperscsi", Description: `HyperSCSI (SCSI over Ethernet)`},
	0x88a2:        {Sym: "ata", Description: `ATA over Ethernet`},
	0x88a4:        {Sym: "ethercat", Description: `EtherCAT Protocol`},
	0x88a8:        {Sym: "service", Description: `Service VLAN tag identifier (S-Tag) on Q-in-Q tunnel.`},
	0x88ab:        {Sym: "ethernet", Description: `Ethernet Powerlink`},
	0x88b8:        {Sym: "goose", Description: `GOOSE (Generic Object Oriented Substation event)`},
	0x88b9:        {Sym: "gse", Description: `GSE (Generic Substation Events) Management Services`},
	0x88ba:        {Sym: "sv", Description: `SV (Sampled Value Transmission)`},
	0x88bf:        {Sym: "mikrotik", Description: `MikroTik RoMON (unofficial)`},
	0x88cc:        {Sym: "link", Description: `Link Layer Discovery Protocol (LLDP)`},
	0x88cd:        {Sym: "sercos", Description: `SERCOS III`},
	0x88e1:        {Sym: "homeplug", Description: `HomePlug Green PHY`},
	0x88e3:        {Sym: "media", Description: `Media Redundancy Protocol (IEC62439-2)`},
	0x88e5:        {Sym: "ieee", Description: `IEEE 802.1AE MAC security (MACsec)`},
	0x88e7:        {Sym: "provider", Description: `Provider Backbone Bridges (PBB) (IEEE 802.1ah)`},
	0x88f7:        {Sym: "precision", Description: `Precision Time Protocol (PTP) over IEEE 802.3 Ethernet`},
	0x88f8:        {Sym: "nc", Description: `NC-SI`},
	0x88fb:        {Sym: "parallel", Description: `Parallel Redundancy Protocol (PRP)`},
	0x8902:        {Sym: "ieee", Description: `IEEE 802.1ag Connectivity Fault Management (CFM) Protocol / ITU-T Recommendation Y.1731 (OAM)`},
	0x8906:        {Sym: "fibre", Description: `Fibre Channel over Ethernet (FCoE)`},
	0x8914:        {Sym: "fcoe", Description: `FCoE Initialization Protocol`},
	0x8915:        {Sym: "rdma", Description: `RDMA over Converged Ethernet (RoCE)`},
	0x891d:        {Sym: "ttethernet", Description: `TTEthernet Protocol Control Frame (TTE)`},
	0x893a:        {Sym: "1905", Description: `1905.1 IEEE Protocol`},
	0x892f:        {Sym: "high", Description: `High-availability Seamless Redundancy (HSR)`},
	0x9000:        {Sym: "ethernet", Description: `Ethernet Configuration Testing Protocol[12]`},
	0xf1c1:        {Sym: "redundancy", Description: `Redundancy Tag (IEEE 802.1CB Frame Replication and Elimination for Reliability)`},
}

var etherTypeFormat = map[uint64]*decode.Group{
	etherTypeIPv4: &ipv4Format,
}

func mapUToEtherSym(s decode.Scalar) (decode.Scalar, error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], s.ActualU())
	s.Sym = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x", b[2], b[3], b[4], b[5], b[6], b[7])
	return s, nil
}

func decodeEthernet(d *decode.D, in interface{}) interface{} {
	d.FieldU("destination", 48, mapUToEtherSym, d.Hex)
	d.FieldU("source", 48, mapUToEtherSym, d.Hex)
	etherType := d.FieldU16("ether_type", d.MapUToScalar(etherTypeMap), d.Hex)
	if g, ok := etherTypeFormat[etherType]; ok {
		d.FieldFormatLen("packet", d.BitsLeft(), *g, nil)
	} else {
		d.FieldRawLen("data", d.BitsLeft())
	}

	return nil
}
