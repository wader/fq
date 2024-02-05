package pcap

// https://pcapng.github.io/pcapng/draft-ietf-opsawg-pcapng.html

import (
	"encoding/binary"
	"net"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/inet/flowsdecoder"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var pcapngLinkFrameGroup decode.Group
var pcapngTCPStreamGroup decode.Group
var pcapngIPvPacket4Group decode.Group

func init() {
	interp.RegisterFormat(
		format.PCAPNG,
		&decode.Format{
			Description: "PCAPNG packet capture",
			RootArray:   true,
			Groups:      []*decode.Group{format.Probe},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Link_Frame}, Out: &pcapngLinkFrameGroup},
				{Groups: []*decode.Group{format.TCP_Stream}, Out: &pcapngTCPStreamGroup},
				{Groups: []*decode.Group{format.IPv4Packet}, Out: &pcapngIPvPacket4Group},
			},
			DecodeFn: decodePcapng,
		})
}

const (
	ngBigEndian    = 0x1a2b3c4d
	ngLittleEndian = 0x4d3c2b1a
)

var ngEndianMap = scalar.UintMapSymStr{
	ngBigEndian:    "big_endian",
	ngLittleEndian: "little_endian",
}

const (
	blockTypeSectionHeader        = 0x0a0d0d0a
	blockTypeInterfaceDescription = 0x00000001
	blockTypeNameResolution       = 0x00000004
	blockTypeInterfaceStatistics  = 0x00000005
	blockTypeEnhancedPacketBlock  = 0x00000006
)

// from https://pcapng.github.io/pcapng/draft-ietf-opsawg-pcapng.html#section_block_code_registry
var blockTypeMap = scalar.UintMap{
	blockTypeInterfaceDescription: {Sym: "interface_description"},
	0x00000002:                    {Description: "Packet Block"},
	0x00000003:                    {Description: "Simple Packet Block"},
	blockTypeNameResolution:       {Sym: "name_resolution"},
	blockTypeInterfaceStatistics:  {Sym: "interface_statistics"},
	blockTypeEnhancedPacketBlock:  {Sym: "enhanced_packet"},
	0x00000007:                    {Description: "IRIG Timestamp Block"},
	0x00000008:                    {Description: "ARINC 429 in AFDX Encapsulation Information Block"},
	0x00000009:                    {Description: "systemd Journal Export Block"},
	0x0000000a:                    {Description: "Decryption Secrets Block"},
	0x00000101:                    {Description: "Hone Project Machine Info Block"},
	0x00000102:                    {Description: "Hone Project Connection Event Block"},
	0x00000201:                    {Description: "Sysdig Machine Info Block"},
	0x00000202:                    {Description: "Sysdig Process Info Block, version 1"},
	0x00000203:                    {Description: "Sysdig FD List Block"},
	0x00000204:                    {Description: "Sysdig Event Block"},
	0x00000205:                    {Description: "Sysdig Interface List Block"},
	0x00000206:                    {Description: "Sysdig User List Block"},
	0x00000207:                    {Description: "Sysdig Process Info Block, version 2"},
	0x00000208:                    {Description: "Sysdig Event Block with flags"},
	0x00000209:                    {Description: "Sysdig Process Info Block, version 3"},
	0x00000210:                    {Description: "Sysdig Process Info Block, version 4"},
	0x00000211:                    {Description: "Sysdig Process Info Block, version 5"},
	0x00000212:                    {Description: "Sysdig Process Info Block, version 6"},
	0x00000213:                    {Description: "Sysdig Process Info Block, version 7"},
	0x00000bad:                    {Description: "Custom Block that rewriters can copy into new files"},
	0x40000bad:                    {Description: "Custom Block that rewriters should not copy into new files"},
	blockTypeSectionHeader:        {Sym: "section_header"},
}

const (
	optionEnd     = 0
	optionComment = 1

	sectionHeaderOptionHardware = 2
	sectionHeaderOptionOS       = 3
	sectionHeaderOptionUserAppl = 4

	interfaceDescriptionName        = 2
	interfaceDescriptionDescription = 3
	interfaceDescriptionIPv4addr    = 4
	interfaceDescriptionMACaddr     = 6
	interfaceDescriptionEUIaddr     = 7
	interfaceDescriptionSpeed       = 8
	interfaceDescriptionTsresol     = 9
	interfaceDescriptionTzone       = 10
	interfaceDescriptionFilter      = 11
	interfaceDescriptionOS          = 12
	interfaceDescriptionFcslen      = 13
	interfaceDescriptionTsoffset    = 14

	enhancedPacketFlags     = 2
	enhancedPacketHash      = 3
	enhancedPacketDropcount = 4

	nameResolutionDNSName    = 2
	nameResolutionDNSIP4addr = 3
	nameResolutionDNSIP6addr = 4

	interfaceStatisticsStarttime    = 2
	interfaceStatisticsEndtime      = 3
	interfaceStatisticsIfRecv       = 4
	interfaceStatisticsIfDrop       = 5
	interfaceStatisticsFilterAccept = 6
	interfaceStatisticsOSDrop       = 7
	interfaceStatisticsUsrdeliv     = 8
)

var sectionHeaderOptionsMap = scalar.UintMap{
	optionEnd:                   {Sym: "end", Description: "End of options"},
	optionComment:               {Sym: "comment", Description: "Comment"},
	sectionHeaderOptionHardware: {Sym: "hardware"},
	sectionHeaderOptionOS:       {Sym: "os"},
	sectionHeaderOptionUserAppl: {Sym: "userappl"},
}

var interfaceDescriptionOptionsMap = scalar.UintMap{
	optionEnd:                       {Sym: "end", Description: "End of options"},
	optionComment:                   {Sym: "comment", Description: "Comment"},
	interfaceDescriptionName:        {Sym: "name"},
	interfaceDescriptionDescription: {Sym: "description"},
	interfaceDescriptionIPv4addr:    {Sym: "ipv4addr"},
	interfaceDescriptionMACaddr:     {Sym: "macaddr"},
	interfaceDescriptionEUIaddr:     {Sym: "euiaddr"},
	interfaceDescriptionSpeed:       {Sym: "speed"},
	interfaceDescriptionTsresol:     {Sym: "tsresol"},
	interfaceDescriptionTzone:       {Sym: "tzone"},
	interfaceDescriptionFilter:      {Sym: "filter"},
	interfaceDescriptionOS:          {Sym: "os"},
	interfaceDescriptionFcslen:      {Sym: "fcslen"},
	interfaceDescriptionTsoffset:    {Sym: "tsoffset"},
}

var enhancedPacketOptionsMap = scalar.UintMap{
	optionEnd:               {Sym: "end", Description: "End of options"},
	optionComment:           {Sym: "comment", Description: "Comment"},
	enhancedPacketFlags:     {Sym: "flags"},
	enhancedPacketHash:      {Sym: "hash"},
	enhancedPacketDropcount: {Sym: "dropcount"},
}

var nameResolutionOptionsMap = scalar.UintMap{
	optionEnd:                {Sym: "end", Description: "End of options"},
	optionComment:            {Sym: "comment", Description: "Comment"},
	nameResolutionDNSName:    {Sym: "dnsname"},
	nameResolutionDNSIP4addr: {Sym: "dnsip4addr"},
	nameResolutionDNSIP6addr: {Sym: "dnsip6addr"},
}

var interfaceStatisticsOptionsMap = scalar.UintMap{
	optionEnd:                       {Sym: "end", Description: "End of options"},
	optionComment:                   {Sym: "comment", Description: "Comment"},
	interfaceStatisticsStarttime:    {Sym: "starttime"},
	interfaceStatisticsEndtime:      {Sym: "endtime"},
	interfaceStatisticsIfRecv:       {Sym: "ifrecv"},
	interfaceStatisticsIfDrop:       {Sym: "ifdrop"},
	interfaceStatisticsFilterAccept: {Sym: "filteraccept"},
	interfaceStatisticsOSDrop:       {Sym: "osdrop"},
	interfaceStatisticsUsrdeliv:     {Sym: "usrdeliv"},
}

const (
	nameResolutionRecordEnd  = 0x0000
	nameResolutionRecordIpv4 = 0x0001
	nameResolutionRecordIpv6 = 0x0002
)

var nameResolutionRecordMap = scalar.UintMapSymStr{
	nameResolutionRecordEnd:  "end",
	nameResolutionRecordIpv4: "ipv4",
	nameResolutionRecordIpv6: "ipv6",
}

func decoodeOptions(d *decode.D, opts scalar.UintMap) {
	if d.BitsLeft() < 32 {
		return
	}
	seenEnd := false
	for !seenEnd {
		d.FieldStruct("option", func(d *decode.D) {
			code := d.FieldU16("code", opts)
			length := d.FieldU16("length")
			if code == optionEnd {
				seenEnd = true
				return
			}
			d.FieldUTF8NullFixedLen("value", int(length))
			d.FieldRawLen("padding", int64(d.AlignBits(32)))
		})
	}
}

// TODO: share
var mapUToIPv4Sym = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(s.Actual))
	s.Sym = net.IP(b[:]).String()
	return s, nil
})

var blockFns = map[uint64]func(d *decode.D, dc *decodeContext){
	// TODO: SimplePacket
	// TODO: Packet
	blockTypeSectionHeader: func(d *decode.D, dc *decodeContext) {
		d.FieldU32BE("byte_order_magic", ngEndianMap, scalar.UintHex)
		d.FieldU16("major_version")
		d.FieldU16("minor_version")
		dc.sectionLength = d.FieldS64("section_length")
		d.FieldArray("options", func(d *decode.D) { decoodeOptions(d, sectionHeaderOptionsMap) })

		dc.sectionHeaderFound = true
	},
	blockTypeInterfaceDescription: func(d *decode.D, dc *decodeContext) {
		typ := d.FieldU16("link_type", format.LinkTypeMap)
		d.FieldU16("reserved")
		d.FieldU32("snap_len")
		d.FieldArray("options", func(d *decode.D) { decoodeOptions(d, interfaceDescriptionOptionsMap) })

		dc.interfaceTypes[len(dc.interfaceTypes)] = int(typ)
	},
	blockTypeEnhancedPacketBlock: func(d *decode.D, dc *decodeContext) {
		interfaceID := d.FieldU32("interface_id")
		d.FieldU32("timestamp_high")
		d.FieldU32("timestamp_low")
		capturedLength := d.FieldU32("capture_packet_length")
		d.FieldU32("original_packet_length")

		bs := d.ReadAllBits(d.BitBufRange(d.Pos(), int64(capturedLength)*8))

		linkType := dc.interfaceTypes[int(interfaceID)]

		if fn, ok := linkToDecodeFn[linkType]; ok {
			// TODO: report decode errors
			_ = fn(dc.flowDecoder, bs)
		}

		d.FieldFormatOrRawLen(
			"packet",
			int64(capturedLength)*8,
			&pcapngLinkFrameGroup,
			format.Link_Frame_In{
				Type:           linkType,
				IsLittleEndian: d.Endian == decode.LittleEndian,
			},
		)

		d.FieldRawLen("padding", int64(d.AlignBits(32)))
		d.FieldArray("options", func(d *decode.D) { decoodeOptions(d, enhancedPacketOptionsMap) })
	},
	blockTypeNameResolution: func(d *decode.D, _ *decodeContext) {
		seenEnd := false
		d.FieldArray("records", func(d *decode.D) {
			for !seenEnd {
				d.FieldStruct("record", func(d *decode.D) {
					typ := d.FieldU16("type", nameResolutionRecordMap)
					length := d.FieldU16("length")
					if typ == nameResolutionRecordEnd {
						seenEnd = true
						return
					}
					d.FramedFn(int64(length)*8, func(d *decode.D) {
						switch typ {
						case nameResolutionRecordIpv4:
							d.FieldU32BE("address", mapUToIPv4Sym, scalar.UintHex)
							d.FieldArray("entries", func(d *decode.D) {
								for !d.End() {
									d.FieldUTF8Null("string")
								}
							})
						default:
							d.FieldUTF8NullFixedLen("value", int(d.BitsLeft()/8))
						}
					})
					d.FieldRawLen("padding", int64(d.AlignBits(32)))
				})
			}
		})
		d.FieldArray("options", func(d *decode.D) { decoodeOptions(d, nameResolutionOptionsMap) })
	},
	blockTypeInterfaceStatistics: func(d *decode.D, _ *decodeContext) {
		d.FieldU32("interface_id")
		d.FieldU32("timestamp_high")
		d.FieldU32("timestamp_low")
		d.FieldRawLen("padding", int64(d.AlignBits(32)))
		d.FieldArray("options", func(d *decode.D) { decoodeOptions(d, interfaceStatisticsOptionsMap) })
	},
}

func decodeBlock(d *decode.D, dc *decodeContext) {
	d.Endian = dc.endian

	typ := d.FieldU32("type", blockTypeMap, scalar.UintHex)
	// section header is special as it stores byte order marker
	if typ == blockTypeSectionHeader {
		d.SeekRel(32)       // skip length
		endian := d.U32BE() // force BE
		switch endian {
		case ngBigEndian:
			dc.endian = decode.BigEndian
		case ngLittleEndian:
			dc.endian = decode.LittleEndian
		default:
			d.Fatalf("unknown endian %d", endian)
		}
		d.Endian = dc.endian
		d.SeekRel(-64)
	}
	length := d.FieldU32("length") - 8
	const footerLengthSize = 32
	blockLen := int64(length)*8 - footerLengthSize
	if blockLen <= 0 {
		d.Fatalf("%d blockLen < 0", blockLen)
	}
	d.FramedFn(blockLen, func(d *decode.D) {
		if fn, ok := blockFns[typ]; ok {
			fn(d, dc)
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	})
	d.FieldU32("footer_length")
}

func decodeSection(d *decode.D, dc *decodeContext) {
	d.FieldArray("blocks", func(d *decode.D) {
		sectionStart := d.Pos()
		// assume and read first section header
		d.FieldStruct("block", func(d *decode.D) { decodeBlock(d, dc) })
		for (dc.sectionLength == -1 && !d.End()) ||
			(dc.sectionLength != -1 && d.Pos()-sectionStart < dc.sectionLength*8) {
			d.FieldStruct("block", func(d *decode.D) { decodeBlock(d, dc) })
		}
	})
}

type decodeContext struct {
	endian             decode.Endian
	sectionLength      int64
	sectionHeaderFound bool
	interfaceTypes     map[int]int
	flowDecoder        *flowsdecoder.Decoder
}

func decodePcapng(d *decode.D) any {
	sectionHeaders := 0
	for !d.End() {
		fd := flowsdecoder.New(flowsdecoder.DecoderOptions{CheckTCPOptions: false})
		dc := decodeContext{
			interfaceTypes: map[int]int{},
			flowDecoder:    fd,
		}

		d.FieldStruct("section", func(d *decode.D) {
			decodeSection(d, &dc)
			fd.Flush()
			fieldFlows(d, dc.flowDecoder, pcapngTCPStreamGroup, pcapngIPvPacket4Group)
		})
		if dc.sectionHeaderFound {
			sectionHeaders++
		}
	}

	if sectionHeaders == 0 {
		d.Fatalf("no section headers found")
	}

	return nil
}
