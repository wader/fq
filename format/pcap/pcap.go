package pcap

// https://wiki.wireshark.org/Development/LibpcapFileFormat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var pcapEther8023Format decode.Group

const (
	bigEndian    = 0xa1b2c3d4
	littleEndian = 0xd4c3b2a1
)

var endianMap = decode.UToStr{
	bigEndian:    "big_endian",
	littleEndian: "little_endian",
}

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PCAP,
		Description: "PCAP packet capture",
		Groups:      []string{format.PROBE},
		Dependencies: []decode.Dependency{
			{Names: []string{format.ETHER8023}, Group: &pcapEther8023Format},
		},
		DecodeFn: decodePcap,
	})
}

func decodePcap(d *decode.D, in interface{}) interface{} {
	endian := d.FieldU32("magic", d.AssertU(bigEndian, littleEndian), d.MapUToStrSym(endianMap), d.Hex)
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
	linkType := int(d.FieldU32("network", d.MapUToScalar(linkTypeMap)))

	d.FieldArray("packets", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("packet", func(d *decode.D) {
				d.FieldU32("ts_sec")
				d.FieldU32("ts_usec")
				inclLen := d.FieldU32("incl_len")
				origLen := d.FieldU32("orig_len")
				if g, ok := linkToFormat[linkType]; ok {
					d.FieldFormatLen("packet", int64(origLen)*8, *g, nil)
				} else {
					d.FieldRawLen("packet", int64(origLen)*8)
				}
				d.FieldRawLen("capture_padding", int64(inclLen-origLen)*8)
			})
		}
	})

	return nil
}
