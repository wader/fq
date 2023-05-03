package inet

import (
	"encoding/binary"
	"net"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/checksum"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var ipv4IpPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.IPv4Packet,
		&decode.Format{
			Description: "Internet protocol v4 packet",
			Groups: []*decode.Group{
				format.INET_Packet,
				format.Link_Frame,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.IP_Packet}, Out: &ipv4IpPacketGroup},
			},
			DecodeFn: decodeIPv4,
		})
}

const (
	ipv4OptionEnd = 0
	ipv4OptionNop = 1
)

var ipv4OptionsMap = scalar.UintMap{
	ipv4OptionEnd: {Sym: "end", Description: "End of options list"},
	ipv4OptionNop: {Sym: "nop", Description: "No operation"},
	2:             {Description: "Security"},
	3:             {Description: "Loose Source Routing"},
	9:             {Description: "Strict Source Routing"},
	7:             {Description: "Record Route"},
	8:             {Description: "Stream ID"},
	4:             {Description: "Internet Timestamp"},
}

var mapUToIPv4Sym = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(s.Actual))
	s.Sym = net.IP(b[:]).String()
	return s, nil
})

func decodeIPv4(d *decode.D) any {
	var ipi format.INET_Packet_In
	var lfi format.Link_Frame_In
	if d.ArgAs(&ipi) && ipi.EtherType != format.EtherTypeIPv4 {
		d.Fatalf("incorrect ethertype %d", ipi.EtherType)
	} else if d.ArgAs(&lfi) && lfi.Type != format.LinkTypeIPv4 && lfi.Type != format.LinkTypeRAW {
		d.Fatalf("incorrect linktype %d", lfi.Type)
	}

	d.FieldU4("version", d.UintAssert(4))
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
	d.FieldU16("header_checksum", scalar.UintHex)
	checksumEnd := d.Pos()
	d.FieldU32("source_ip", mapUToIPv4Sym, scalar.UintHex)
	d.FieldU32("destination_ip", mapUToIPv4Sym, scalar.UintHex)
	optionsLen := (int64(ihl) - 5) * 8 * 4
	if optionsLen > 0 {
		d.FramedFn(optionsLen, func(d *decode.D) {
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
	d.Copy(ipv4Checksum, bitio.NewIOReader(d.BitBufRange(0, checksumStart)))
	d.Copy(ipv4Checksum, bitio.NewIOReader(d.BitBufRange(checksumEnd, headerEnd-checksumEnd)))
	_ = d.FieldMustGet("header_checksum").TryUintScalarFn(d.UintValidateBytes(ipv4Checksum.Sum(nil)), scalar.UintHex)

	dataLen := int64(totalLength-(ihl*4)) * 8

	if moreFragments || fragmentOffset > 0 {
		d.FieldRawLen("payload", dataLen)
	} else {
		d.FieldFormatOrRawLen(
			"payload",
			dataLen,
			&ipv4IpPacketGroup,
			format.IP_Packet_In{Protocol: int(protocol)},
		)
	}

	return nil
}
