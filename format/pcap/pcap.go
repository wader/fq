package pcap

// https://wiki.wireshark.org/Development/LibpcapFileFormat
// TODO: tshark seems to not support sll2 in pcap, confusing

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/inet/flowsdecoder"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var pcapLinkFrameFormat decode.Group
var pcapTCPStreamFormat decode.Group
var pcapIPv4PacketFormat decode.Group

const (
	bigEndian    = 0xa1b2c3d4
	littleEndian = 0xd4c3b2a1
)

var endianMap = scalar.UToSymStr{
	bigEndian:    "big_endian",
	littleEndian: "little_endian",
}

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PCAP,
		Description: "PCAP packet capture",
		Groups:      []string{format.PROBE},
		Dependencies: []decode.Dependency{
			{Names: []string{format.LINK_FRAME}, Group: &pcapLinkFrameFormat},
			{Names: []string{format.TCP_STREAM}, Group: &pcapTCPStreamFormat},
			{Names: []string{format.IPV4_PACKET}, Group: &pcapIPv4PacketFormat},
		},
		DecodeFn: decodePcap,
	})
}

func decodePcap(d *decode.D, in interface{}) interface{} {
	endian := d.FieldU32("magic", d.AssertU(bigEndian, littleEndian), endianMap, scalar.Hex)
	switch endian {
	case bigEndian:
		d.Endian = decode.BigEndian
	case littleEndian:
		d.Endian = decode.LittleEndian
	default:
		d.Fatalf("unknown endian %d", endian)
	}
	d.FieldU16("version_major")
	d.FieldU16("version_minor")
	d.FieldS32("thiszone")
	d.FieldU32("sigfigs")
	d.FieldU32("snaplen")
	linkType := int(d.FieldU32("network", format.LinkTypeMap))

	fd := flowsdecoder.New()

	d.FieldArray("packets", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("packet", func(d *decode.D) {
				d.FieldU32("ts_sec")
				d.FieldU32("ts_usec")
				inclLen := d.FieldU32("incl_len")
				origLen := d.FieldU32("orig_len")

				// "incl_len: the number of bytes of packet data actually captured and saved in the file. This value should never become larger than orig_len or the snaplen value of the global header"
				// "orig_len: the length of the packet as it appeared on the network when it was captured. If incl_len and orig_len differ, the actually saved packet size was limited by snaplen."

				// TODO: incl_len seems to be larger than snaplen in real pcap files
				// if inclLen > snapLen {
				// 	d.Errorf("incl_len %d > snaplen %d", inclLen, snapLen)
				// }

				if inclLen > origLen {
					d.Errorf("incl_len %d > orig_len %d", inclLen, origLen)
				}

				bb := d.BitBufRange(d.Pos(), int64(inclLen)*8)
				bs, err := bb.Bytes()
				if err != nil {
					// TODO:
					d.IOPanic(err, "BitBufRange: inclLen")
				}

				if fn, ok := linkToDecodeFn[linkType]; ok {
					// TODO: report decode errors
					_ = fn(fd, bs)
				}

				if dv, _, _ := d.TryFieldFormatLen("packet", int64(inclLen)*8, pcapLinkFrameFormat, format.LinkFrameIn{
					Type:         linkType,
					LittleEndian: d.Endian == decode.LittleEndian,
				}); dv == nil {
					d.FieldRawLen("packet", int64(inclLen)*8)
				}
			})
		}
	})
	fd.Flush()

	fieldFlows(d, fd, pcapTCPStreamFormat, pcapIPv4PacketFormat)

	return nil
}
