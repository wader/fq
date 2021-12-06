package pcap

// https://wiki.wireshark.org/Development/LibpcapFileFormat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/inet/flowsdecoder"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var pcapEther8023Format decode.Group
var pcapSLLPacket decode.Group
var pcapSLL2Packet decode.Group
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
			{Names: []string{format.ETHER8023_FRAME}, Group: &pcapEther8023Format},
			{Names: []string{format.SLL_PACKET}, Group: &pcapSLLPacket},
			{Names: []string{format.SLL2_PACKET}, Group: &pcapSLL2Packet},
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

				bb := d.BitBufRange(d.Pos(), int64(origLen)*8)
				bs, err := bb.Bytes()
				if err != nil {
					// TODO:
					panic(err)
				}

				if fn, ok := linkToDecodeFn[linkType]; ok {
					// TODO: report decode errors
					_ = fn(fd, bs)
				}

				if g, ok := linkToFormat[linkType]; ok {
					d.FieldFormatLen("packet", int64(origLen)*8, *g, nil)
				} else {
					d.FieldRawLen("packet", int64(origLen)*8)
				}
				d.FieldRawLen("capture_padding", int64(inclLen-origLen)*8)
			})
		}
	})
	fd.Flush()

	fieldFlows(d, fd, pcapTCPStreamFormat, pcapIPv4PacketFormat)

	return nil
}
